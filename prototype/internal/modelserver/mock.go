package modelserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Mock is the M1 backend: it streams canned OpenAI chat-completion chunks so the
// loop, streaming, and doors can be proven with no real inference. Each message
// is tagged with the node id, so on the peer path you can see which node
// actually served the request. In M2 this is replaced by the real OpenAI client
// to Ollama, behind the same Backend interface.
type Mock struct {
	nodeID string
}

var _ Backend = (*Mock)(nil)

func NewMock(nodeID string) *Mock { return &Mock{nodeID: nodeID} }

func (m *Mock) ListModels(ctx context.Context) ([]string, error) {
	return []string{"mock-model"}, nil
}

func (m *Mock) ChatCompletionStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error) {
	model := req.Model
	if model == "" {
		model = "mock-model"
	}
	message := fmt.Sprintf("[mock] served by node %s — no real inference yet (M1).", m.nodeID)

	// A pipe lets the goroutine push SSE chunks while the reader (ingress)
	// pulls and flushes them, so backpressure and streaming are real.
	pr, pw := io.Pipe()
	go func() {
		var err error
		defer func() { pw.CloseWithError(err) }()

		created := time.Now().Unix()
		id := "chatcmpl-mock-" + m.nodeID

		// First chunk carries the assistant role, like OpenAI.
		if err = writeChunk(pw, id, model, created, delta{Role: "assistant"}, nil); err != nil {
			return
		}
		// Then stream the message word by word, with a small pause so the
		// streaming is actually observable on the wire.
		for _, word := range strings.Fields(message) {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case <-time.After(40 * time.Millisecond):
			}
			if err = writeChunk(pw, id, model, created, delta{Content: word + " "}, nil); err != nil {
				return
			}
		}
		// Final chunk sets finish_reason, then the SSE terminator.
		stop := "stop"
		if err = writeChunk(pw, id, model, created, delta{}, &stop); err != nil {
			return
		}
		_, err = io.WriteString(pw, "data: [DONE]\n\n")
	}()
	return pr, nil
}

// OpenAI chat.completion.chunk shapes, minimal subset.
type chunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []chunkChoice `json:"choices"`
}

type chunkChoice struct {
	Index        int     `json:"index"`
	Delta        delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

type delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func writeChunk(w io.Writer, id, model string, created int64, d delta, finish *string) error {
	payload, err := json.Marshal(chunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []chunkChoice{{Index: 0, Delta: d, FinishReason: finish}},
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
	return err
}
