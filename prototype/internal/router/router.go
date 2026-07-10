package router

import (
	"context"
	"errors"
	"io"

	"dii/internal/manifest"
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

// Router decides where a request runs. M2 makes the local door manifest-gated:
// serve locally only if our own manifest has the requested model, else fail
// honestly. Local-then-peer overflow and manifest-aware consumer routing arrive
// in M3; the consumer door here still forwards to the first peer as in M1.
type Router struct {
	local    modelserver.Backend
	peers    []modelserver.Backend
	manifest *manifest.Store
}

func New(local modelserver.Backend, peers []modelserver.Backend, store *manifest.Store) *Router {
	return &Router{local: local, peers: peers, manifest: store}
}

// Route picks a backend for the door and returns its raw OpenAI SSE stream.
func (r *Router) Route(ctx context.Context, door Door, req modelserver.ChatRequest) (io.ReadCloser, error) {
	switch door {
	case DoorConsumer:
		// The consumer has nothing to serve locally; go straight to a peer.
		// (M3 makes this a manifest-aware selection.)
		if len(r.peers) == 0 {
			return nil, ErrNoCapacity
		}
		return r.peers[0].ChatCompletionStream(ctx, req)
	default: // DoorLocal
		// Local-first: only serve if our own model server actually has the model.
		if !r.manifest.LocalCanServe(req.Model) {
			return nil, ErrNoCapacity
		}
		return r.local.ChatCompletionStream(ctx, req)
	}
}
