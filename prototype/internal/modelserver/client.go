package modelserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is the real backend: an OpenAI-compatible HTTP client to a model
// server (Ollama for the POC). It uses only the standard OpenAI subset
// (/v1/models, /v1/chat/completions with streaming), so any server that speaks
// that API drops in behind the same base URL. This is the M2 replacement for
// the mock, wired in through the Backend interface.
type Client struct {
	baseURL string // e.g. http://localhost:11434/v1
	http    *http.Client
}

var _ Backend = (*Client)(nil)

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		// No client-level timeout: chat streams are long-lived; cancellation
		// rides on the request context.
		http: &http.Client{},
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
