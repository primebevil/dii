package router

import (
	"context"
	"errors"
	"io"

	"dii/internal/modelserver"
)

// Door is how a request entered the node. The token check in ingress picks it.
type Door int

const (
	DoorLocal    Door = iota // trusted owner: serve on our own backend
	DoorConsumer             // stranger with the stub token: hand off to a peer
)

// ErrNoCapacity is the honest failure when no backend can serve the request.
var ErrNoCapacity = errors.New("no node in the pod can serve this request")

// Router decides where a request runs. M1 is deliberately dumb and has no
// manifest yet: the local door serves on our own backend, the consumer door
// hands off to the first peer. Local-first fallback and manifest-based
// capability matching arrive in M2/M3.
type Router struct {
	local modelserver.Backend
	peers []modelserver.Backend
}

func New(local modelserver.Backend, peers []modelserver.Backend) *Router {
	return &Router{local: local, peers: peers}
}

// Route picks a backend for the door and returns its raw OpenAI SSE stream.
func (r *Router) Route(ctx context.Context, door Door, req modelserver.ChatRequest) (io.ReadCloser, error) {
	switch door {
	case DoorConsumer:
		// The consumer has nothing to serve locally; go straight to a peer.
		if len(r.peers) == 0 {
			return nil, ErrNoCapacity
		}
		return r.peers[0].ChatCompletionStream(ctx, req)
	default: // DoorLocal
		return r.local.ChatCompletionStream(ctx, req)
	}
}
