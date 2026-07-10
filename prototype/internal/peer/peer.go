package peer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"dii/internal/manifest"
	"dii/internal/modelserver"
)

// Client talks to a peer node's OpenAI-compatible endpoint. Inter-node transport
// is deliberately just the OpenAI HTTP call reused between nodes (BUILD_BRIEF):
// we call the peer exactly like any other OpenAI client, with no consumer token,
// so the peer serves the request on its own local door and does not forward
// again. That is also what keeps a peer usable as a plain modelserver.Backend.
type Client struct {
	endpoint string // e.g. http://host:8080
	http     *http.Client
}

var _ modelserver.Backend = (*Client)(nil)

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		// No client-level timeout: chat streams are long-lived; cancellation
		// rides on the request context instead.
		http: &http.Client{},
	}
}

// Endpoint returns the peer's normalized base URL.
func (c *Client) Endpoint() string { return c.endpoint }

// FetchManifest reads the peer's self-description from its /manifest route.
func (c *Client) FetchManifest(ctx context.Context) (manifest.Manifest, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/manifest", nil)
	if err != nil {
		return manifest.Manifest{}, err
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("peer %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return manifest.Manifest{}, fmt.Errorf("peer %s: /manifest status %d", c.endpoint, resp.StatusCode)
	}
	var m manifest.Manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return manifest.Manifest{}, err
	}
	return m, nil
}

func (c *Client) ChatCompletionStream(ctx context.Context, req modelserver.ChatRequest) (io.ReadCloser, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/v1/chat/completions", bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("peer %s: %w", c.endpoint, err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		resp.Body.Close()
		return nil, fmt.Errorf("peer %s: unexpected status %d: %s", c.endpoint, resp.StatusCode, bytes.TrimSpace(body))
	}
	return resp.Body, nil
}

func (c *Client) ListModels(ctx context.Context) ([]string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/v1/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("peer %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("peer %s: /v1/models status %d", c.endpoint, resp.StatusCode)
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
