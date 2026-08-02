package modelserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// maxEmbeddingsBytes caps an embeddings response we buffer into memory. Vectors
// are verbose as JSON text (a 384-dimension vector is roughly 8 KB), so the cap
// is generous; the request-body cap in ingress is what really bounds batch size.
const maxEmbeddingsBytes = 32 << 20 // 32 MiB

// Client is the real backend: an OpenAI-compatible HTTP client to a model
// server (Ollama for the POC). It uses only the standard OpenAI subset
// (/v1/models, /v1/chat/completions with streaming, /v1/embeddings), so any
// server that speaks that API drops in behind the same base URL. This is the M2
// replacement for the mock, wired in through the Backend interface.
type Client struct {
	baseURL string // e.g. http://localhost:11434/v1
	http    *http.Client
}

var _ Backend = (*Client)(nil)

// NewClient builds a client to an OpenAI-compatible model server.
// responseHeaderTimeout bounds the wait for the first response byte (0 = no
// limit); there is deliberately no overall client timeout, since chat streams
// are long-lived and cancellation rides on the request context.
func NewClient(baseURL string, responseHeaderTimeout time.Duration) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = responseHeaderTimeout
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Transport: transport},
	}
}

func (c *Client) ChatCompletionStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("model server %s: %w", c.baseURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		resp.Body.Close()
		return nil, fmt.Errorf("model server %s: unexpected status %d: %s", c.baseURL, resp.StatusCode, bytes.TrimSpace(body))
	}
	return resp.Body, nil
}

// Embeddings performs a single OpenAI /v1/embeddings call and returns the raw
// response body. There is no streaming here, so unlike ChatCompletionStream the
// whole body is read before returning, which also means a non-200 carries the
// server's own error text back to the caller.
func (c *Client) Embeddings(ctx context.Context, req EmbeddingsRequest) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embeddings", bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("model server %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxEmbeddingsBytes))
	if err != nil {
		return nil, fmt.Errorf("model server %s: reading embeddings response: %w", c.baseURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("model server %s: /v1/embeddings status %d: %s", c.baseURL, resp.StatusCode, bytes.TrimSpace(body))
	}
	return body, nil
}

func (c *Client) ListModels(ctx context.Context) ([]string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("model server %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("model server %s: /v1/models status %d", c.baseURL, resp.StatusCode)
	}
	var out struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	models := make([]string, 0, len(out.Data))
	for _, d := range out.Data {
		models = append(models, d.ID)
	}
	return models, nil
}
