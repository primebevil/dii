package usage

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

// maxBufferedBytes caps what a Meter holds in memory. For a streaming response
// that is only the partial trailing line; for a non-streaming one it is the whole
// body, which has to be buffered to find the usage block at the end. Past the cap
// the Meter stops accumulating and reports its counts as estimated rather than
// growing without bound.
const maxBufferedBytes = 8 << 20 // 8 MiB

// Meter wraps a chat response body, relays every byte through untouched, and
// reads token usage off the response as it passes.
//
// Where the numbers come from, best first:
//
//  1. A non-streaming response (`"stream": false`) carries a real usage block.
//     Exact.
//  2. A streaming response carries one too, but only when the caller asked for it
//     with `stream_options: {"include_usage": true}`. Exact.
//  3. Otherwise — the common streaming case — no OpenAI-compatible server reports
//     usage at all. The completion count falls back to a proxy, the number of
//     content chunks seen; prompt tokens are unknown and reported as 0; and Usage
//     marks the result estimated.
//
// The node deliberately does not inject `stream_options` into the caller's request
// to force case 2. That would mutate a request the node promises to forward
// verbatim and append a usage chunk some strict clients do not expect. Exact
// accounting on the default streaming path is a known gap, called out in the M1
// notes, and is the thing to close when M3's token budgets need to be precise.
//
// A Meter is used by one goroutine (the relay loop), and Usage is read after that
// loop and Close have finished.
type Meter struct {
	src    io.ReadCloser
	stream bool
	start  time.Time

	firstByte time.Time
	buf       []byte // partial trailing SSE line; or the whole body when !stream

	prompt     int
	completion int
	chunks     int  // content chunks seen, the fallback completion proxy
	exact      bool // a real usage block was found
}

var _ io.ReadCloser = (*Meter)(nil)

// NewMeter meters src. stream says which response shape to expect (taken from the
// request's own "stream" field, which is what determines it). start is when the
// request arrived, so TTFT measures the whole path the caller experienced.
func NewMeter(src io.ReadCloser, stream bool, start time.Time) *Meter {
	return &Meter{src: src, stream: stream, start: start}
}

// Read is a pure passthrough: the caller's bytes are never altered, only observed.
func (m *Meter) Read(p []byte) (int, error) {
	n, err := m.src.Read(p)
	if n > 0 {
		if m.firstByte.IsZero() {
			m.firstByte = time.Now()
		}
		m.consume(p[:n])
	}
	return n, err
}

// Close finishes accounting and closes the underlying body.
func (m *Meter) Close() error {
	if !m.stream {
		// The whole non-streaming body is in hand now, so parse it for usage.
		m.apply(m.buf)
	}
	m.buf = nil
	return m.src.Close()
}

// Usage reports the token counts read off the response. When estimated is true
// the completion count is a content-chunk proxy and prompt tokens are unknown.
func (m *Meter) Usage() (prompt, completion int, estimated bool) {
	if m.exact {
		return m.prompt, m.completion, false
	}
	return 0, m.chunks, true
}

// TTFT is the time from the request arriving to the first byte of the response,
// which for a streaming request is time to first token. Zero if nothing arrived.
func (m *Meter) TTFT() time.Duration {
	if m.firstByte.IsZero() {
		return 0
	}
	return m.firstByte.Sub(m.start)
}

func (m *Meter) consume(b []byte) {
	if len(m.buf) > maxBufferedBytes {
		return
	}
	m.buf = append(m.buf, b...)
	if !m.stream {
		return
	}

	// SSE frames arrive split across arbitrary read boundaries, so parse whole
	// lines and keep the remainder for the next read.
	for {
		i := bytes.IndexByte(m.buf, '\n')
		if i < 0 {
			return
		}
		line := m.buf[:i]
		m.buf = m.buf[i+1:]
		m.parseSSELine(line)
	}
}

func (m *Meter) parseSSELine(line []byte) {
	line = bytes.TrimSpace(line)
	if !bytes.HasPrefix(line, []byte("data:")) {
		return
	}
	payload := bytes.TrimSpace(line[len("data:"):])
	if len(payload) == 0 || bytes.Equal(payload, []byte("[DONE]")) {
		return
	}
	m.apply(payload)
}

// apply reads one JSON payload — an SSE chunk or a whole response — taking its
// usage block if it has one, and counting content otherwise.
func (m *Meter) apply(payload []byte) {
	var resp chatResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		return
	}
	if resp.Usage != nil {
		m.prompt, m.completion, m.exact = resp.Usage.PromptTokens, resp.Usage.CompletionTokens, true
		return
	}
	for _, c := range resp.Choices {
		if c.Delta.Content != "" {
			m.chunks++
		}
	}
}

// FromResponse reads the usage block out of a single non-streaming OpenAI
// response body. This is how the embeddings path gets its counts, which are
// exact: an embeddings response always carries usage.
func FromResponse(body []byte) (prompt, completion int, estimated bool) {
	var resp chatResponse
	if err := json.Unmarshal(body, &resp); err != nil || resp.Usage == nil {
		return 0, 0, true
	}
	return resp.Usage.PromptTokens, resp.Usage.CompletionTokens, false
}

// chatResponse is the minimal shared subset of an OpenAI chat.completion.chunk, a
// whole chat.completion, and an embeddings response: enough to find a usage block,
// and enough to count streamed content when there isn't one.
type chatResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}
