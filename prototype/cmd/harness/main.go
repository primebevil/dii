// Command harness measures one request path against a node's OpenAI endpoint:
// time to first token, streaming throughput, and total wall time, over N runs.
//
// It sends a fixed prompt with a fixed max_tokens so runs are comparable, warms
// the model first (to measure steady state, not cold load), and reports medians.
// Run it against each path (local, direct-to-peer, overflow, consumer) and the
// overhead is overflow-minus-the-peer's-own-local for the same request.
//
//	go run ./cmd/harness -url http://localhost:8080 -model llama3.2:1b -label local
package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

//go:embed prompt.txt
var defaultPrompt string

type sample struct {
	ttft   time.Duration // request start -> first content token
	total  time.Duration // request start -> stream end
	genDur time.Duration // first token -> stream end (the streaming window)
	tokens int           // content chunks received
}

func main() {
	url := flag.String("url", "http://localhost:8080", "node base URL to hit")
	model := flag.String("model", "", "model name (required)")
	token := flag.String("token", "", "bearer token for the consumer door (optional)")
	promptFile := flag.String("prompt-file", "", "prompt file (default: embedded prompt.txt)")
	maxTokens := flag.Int("max-tokens", 64, "fixed output length")
	n := flag.Int("n", 5, "measured iterations")
	warmup := flag.Int("warmup", 1, "warmup iterations (not measured)")
	label := flag.String("label", "", "label for the output row (default: the model)")
	flag.Parse()

	if *model == "" {
		fmt.Fprintln(os.Stderr, "error: -model is required")
		os.Exit(2)
	}
	prompt := defaultPrompt
	if *promptFile != "" {
		b, err := os.ReadFile(*promptFile)
		if err != nil {
			fatal(err)
		}
		prompt = string(b)
	}
	prompt = strings.TrimSpace(prompt)
	label2 := *label
	if label2 == "" {
		label2 = *model
	}

	for i := 0; i < *warmup; i++ {
		if _, err := measure(*url, *model, *token, prompt, *maxTokens); err != nil {
			fatal(fmt.Errorf("warmup: %w", err))
		}
	}

	samples := make([]sample, 0, *n)
	for i := 0; i < *n; i++ {
		s, err := measure(*url, *model, *token, prompt, *maxTokens)
		if err != nil {
			fatal(fmt.Errorf("run %d: %w", i+1, err))
		}
		samples = append(samples, s)
	}

	ttft := medianDur(mapDur(samples, func(s sample) time.Duration { return s.ttft }))
	total := medianDur(mapDur(samples, func(s sample) time.Duration { return s.total }))
	tps := medianFloat(mapFloat(samples, func(s sample) float64 {
		if s.genDur <= 0 {
			return 0
		}
		return float64(s.tokens) / s.genDur.Seconds()
	}))
	tokens := medianInt(mapInt(samples, func(s sample) int { return s.tokens }))

	fmt.Printf("%-20s %-20s n=%-3d ttft=%7.1fms  tok/s=%6.1f  total=%8.1fms  tokens=%d\n",
		label2, *model, *n, ms(ttft), tps, ms(total), tokens)
}

func measure(base, model, token, prompt string, maxTokens int) (sample, error) {
	payload, _ := json.Marshal(map[string]any{
		"model":       model,
		"stream":      true,
		"max_tokens":  maxTokens,
		"temperature": 0,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return sample{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return sample{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return sample{}, fmt.Errorf("status %d: %s", resp.StatusCode, bytes.TrimSpace(b))
	}

	var firstToken time.Time
	tokens := 0
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 || chunk.Choices[0].Delta.Content == "" {
			continue
		}
		if firstToken.IsZero() {
			firstToken = time.Now()
		}
		tokens++
	}
	if err := sc.Err(); err != nil {
		return sample{}, err
	}
	end := time.Now()
	if firstToken.IsZero() {
		return sample{}, fmt.Errorf("no content tokens received")
	}
	return sample{
		ttft:   firstToken.Sub(start),
		total:  end.Sub(start),
		genDur: end.Sub(firstToken),
		tokens: tokens,
	}, nil
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }

func mapDur(s []sample, f func(sample) time.Duration) []time.Duration {
	out := make([]time.Duration, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}
func mapFloat(s []sample, f func(sample) float64) []float64 {
	out := make([]float64, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}
func mapInt(s []sample, f func(sample) int) []int {
	out := make([]int, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

func medianDur(v []time.Duration) time.Duration {
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	return v[len(v)/2]
}
func medianFloat(v []float64) float64 {
	sort.Float64s(v)
	return v[len(v)/2]
}
func medianInt(v []int) int {
	sort.Ints(v)
	return v[len(v)/2]
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
