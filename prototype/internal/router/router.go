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

// Peer is a peer node paired with the endpoint that identifies it in the
// manifest table, so a capability hit (an endpoint) maps back to the client to call.
type Peer struct {
	Endpoint string
	Backend  modelserver.Backend
}

// Router decides where a request runs, using the cached manifests for
// capability matching (model-name match). The two doors differ in preference:
//   - Local door (owner): local-first, then overflow to a capable peer.
//   - Consumer door (guest): "skip the local step" — try the pod's peers first,
//     then fall back to local.
// Either way, if nobody can serve the model, fail honestly with ErrNoCapacity.
type Router struct {
	local    modelserver.Backend
	peers    []Peer
	manifest *manifest.Store
}

func New(local modelserver.Backend, peers []Peer, store *manifest.Store) *Router {
	return &Router{local: local, peers: peers, manifest: store}
}

// Route picks where the request runs and returns the chosen backend's raw
// OpenAI SSE stream.
func (r *Router) Route(ctx context.Context, door Door, req modelserver.ChatRequest) (io.ReadCloser, error) {
	switch door {
	case DoorConsumer:
		// Guest: prefer the pod's shared capacity, keeping the owner's node free.
		if p, ok := r.peerFor(req.Model); ok {
			return p.Backend.ChatCompletionStream(ctx, req)
		}
		if r.manifest.LocalCanServe(req.Model) {
			return r.local.ChatCompletionStream(ctx, req)
		}
		return nil, ErrNoCapacity
	default: // DoorLocal
		// Owner: local-first, then overflow to a peer that has the model.
		if r.manifest.LocalCanServe(req.Model) {
			return r.local.ChatCompletionStream(ctx, req)
		}
		if p, ok := r.peerFor(req.Model); ok {
			return p.Backend.ChatCompletionStream(ctx, req)
		}
		return nil, ErrNoCapacity
	}
}

// peerFor returns the first peer whose cached manifest can serve the model.
func (r *Router) peerFor(model string) (Peer, bool) {
	for _, endpoint := range r.manifest.PeersForModel(model) {
		for _, p := range r.peers {
			if p.Endpoint == endpoint {
				return p, true
			}
		}
	}
	return Peer{}, false
}
