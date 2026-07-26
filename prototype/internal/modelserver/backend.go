package modelserver

import (
	"context"
	"io"
)

// ChatRequest is an OpenAI /v1/chat/completions call. The raw request Body is
// kept as bytes so the node stays a transparent broker: it forwards the payload
// verbatim and only peeks at Model, which routing needs for capability matching
// (arrives in M2).
type ChatRequest struct {
	Model string // the "model" field, pulled out for routing
	Body  []byte // the untouched request body, forwarded as-is
}

// Backend is the seam that keeps the node backend-portable. Anything that can
// list models and stream an OpenAI chat completion is a Backend:
//   - the mock (M1),
//   - the real OpenAI client to Ollama (M2),
//   - a peer node, which from our side is just another OpenAI-compatible endpoint.
//
// The node core depends only on this interface, never on Ollama directly.
//
// ChatCompletionStream returns the raw OpenAI SSE byte stream. The caller relays
// those bytes verbatim to its own client, so tokens stream through end to end.
type Backend interface {
	ListModels(ctx context.Context) ([]string, error)
	ChatCompletionStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error)
}
