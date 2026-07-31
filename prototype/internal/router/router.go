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

// ServedLocal is the ServedBy label for work this node's own backend ran; any
// other value is the endpoint of the peer that served it.
const ServedLocal = "local"

// ErrNoCapacity is the honest failure when no backend can serve the request.
var ErrNoCapacity = errors.New("no node in the pod can serve this request")

// Peer is a peer node paired with the endpoint that identifies it in the
// manifest table, so a capability hit (an endpoint) maps back to the client to call.
type Peer struct {
	Endpoint string
	Backend  modelserver.Backend
}

// Decision records where the router actually sent a request. It rides back out
// to the usage log so a request can be attributed to the node that ran it, which
// is what offload sizing needs: what got routed, where, and how big it was.
type Decision struct {
	ServedBy string // ServedLocal, or the peer endpoint that served it
}

// Router decides where a request runs, using the cached manifests for
// capability matching (model-name match). The two doors differ in preference:
//   - Local door (owner): local-first, then overflow to a capable peer.
//   - Consumer door (guest): "skip the local step" — try the pod's peers first,
//     then fall back to local.
//
// Either way, if nobody can serve the model, fail honestly with ErrNoCapacity.
type Router struct {
	local    modelserver.Backend
	peers    []Peer
	manifest *manifest.Store
}

func New(local modelserver.Backend, peers []Peer, store *manifest.Store) *Router {
	return &Router{local: local, peers: peers, manifest: store}
}

// Route picks where a chat request runs and returns the chosen backend's raw
// OpenAI response stream.
func (r *Router) Route(ctx context.Context, door Door, req modelserver.ChatRequest) (io.ReadCloser, Decision, error) {
	t, ok := r.pick(door, req.Model)
	if !ok {
		return nil, Decision{}, ErrNoCapacity
	}
	stream, err := t.backend.ChatCompletionStream(ctx, req)
	return stream, Decision{ServedBy: t.servedBy}, err
}

// RouteEmbeddings picks where an embeddings request runs and returns the chosen
// backend's raw response body. Capability matching and door preference are the
// same as for chat: an embeddings model is just a model name in the manifest.
func (r *Router) RouteEmbeddings(ctx context.Context, door Door, req modelserver.EmbeddingsRequest) ([]byte, Decision, error) {
	t, ok := r.pick(door, req.Model)
	if !ok {
		return nil, Decision{}, ErrNoCapacity
	}
	body, err := t.backend.Embeddings(ctx, req)
	return body, Decision{ServedBy: t.servedBy}, err
}

// target is a chosen place to run a request: the backend to call, plus the label
// that identifies it in the usage log.
type target struct {
	backend  modelserver.Backend
	servedBy string
}

// pick applies the door's preference order. Both endpoints share it, so the
// two-doors rule lives in exactly one place.
func (r *Router) pick(door Door, model string) (target, bool) {
	localCanServe := r.manifest.LocalCanServe(model)

	if door == DoorConsumer {
		// Guest: prefer the pod's shared capacity, keeping the owner's node free.
		if p, ok := r.peerFor(model); ok {
			return target{backend: p.Backend, servedBy: p.Endpoint}, true
		}
		if localCanServe {
			return target{backend: r.local, servedBy: ServedLocal}, true
		}
		return target{}, false
	}

	// Owner: local-first, then overflow to a peer that has the model.
	if localCanServe {
		return target{backend: r.local, servedBy: ServedLocal}, true
	}
	if p, ok := r.peerFor(model); ok {
		return target{backend: p.Backend, servedBy: p.Endpoint}, true
	}
	return target{}, false
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
