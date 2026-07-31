package modelserver

import (
	"context"
	"io"
)

// ChatRequest is an OpenAI /v1/chat/completions call. The raw request Body is
// kept as bytes so the node stays a transparent broker: it forwards the payload
// verbatim and only peeks at the two fields the node itself needs — Model for
// capability matching, and Stream to know the response shape (OpenAI defaults
// stream to false, in which case the reply is a single JSON object, not SSE).
type ChatRequest struct {
	Model  string // the "model" field, pulled out for routing
	Stream bool   // the "stream" field, pulled out to label the response
	Body   []byte // the untouched request body, forwarded as-is
}

// EmbeddingsRequest is an OpenAI /v1/embeddings call: the third of the
// reliable-floor bundle (ADR-0006) alongside general chat and a coder. Same
// transparent-broker shape as ChatRequest.
type EmbeddingsRequest struct {
	Model string // the "model" field, pulled out for routing
	Body  []byte // the untouched request body, forwarded as-is
}

// Backend is the seam that keeps the node backend-portable. Anything that can
// list models, stream an OpenAI chat completion, and return embeddings is a
// Backend:
//   - the mock (M1),
//   - the real OpenAI client to Ollama (M2),
//   - a peer node, which from our side is just another OpenAI-compatible endpoint.
//
// The node core depends only on this interface, never on Ollama directly, and
// only ever on the standard OpenAI subset.
//
// ChatCompletionStream returns the raw OpenAI response body. The caller relays
// those bytes verbatim to its own client, so tokens stream through end to end.
// Embeddings is a single request/response call with no streaming, so it returns
// the whole body rather than a reader.
type Backend interface {
	ListModels(ctx context.Context) ([]string, error)
	ChatCompletionStream(ctx context.Context, req ChatRequest) (io.ReadCloser, error)
	Embeddings(ctx context.Context, req EmbeddingsRequest) ([]byte, error)
}
