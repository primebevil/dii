// Package manifest is how the router answers "who can serve this?". Each node
// publishes a small self-description built from its model server's model list,
// and caches its peers' descriptions. Capability match is model-name match (the
// POC shortcut); richer capability tags are a later layer.
package manifest

import (
	"sort"
	"sync"
)

// Manifest is a node's self-description, served at /manifest and exchanged
// between peers at startup.
type Manifest struct {
	NodeID   string   `json:"node_id"`
	Endpoint string   `json:"endpoint"`
	Models   []string `json:"models"`
	Busy     bool     `json:"busy"` // always false in M2; a real load signal is parked
}

// Store holds this node's own manifest plus a cached table of its peers'.
type Store struct {
	mu    sync.RWMutex
	own   Manifest
	peers map[string]Manifest // keyed by peer endpoint
}

func NewStore(own Manifest) *Store {
	return &Store{own: own, peers: make(map[string]Manifest)}
}

func (s *Store) Own() Manifest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.own
}

func (s *Store) SetPeer(endpoint string, m Manifest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers[endpoint] = m
}

func (s *Store) Peers() []Manifest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Manifest, 0, len(s.peers))
	for _, m := range s.peers {
		out = append(out, m)
	}
	return out
}

// ModelSource is a model this node can get served, paired with where it would run.
type ModelSource struct {
	Model    string
	Endpoint string // empty for this node's own backend; otherwise the peer's endpoint
}

// RoutableModels returns every model this node can have served — its own plus its
// cached peers' — deduped, with the local copy winning when both hold a model.
//
// This is deliberately a wider answer than Own(): what a caller can ask for is the
// whole pod, not just this node's shelf. A client that populates a model picker
// from /v1/models can only discover overflow if the node admits to it.
//
// It is NOT what /manifest publishes. A manifest says what a node itself holds; if
// nodes republished models they had merely borrowed, peers would advertise each
// other's capacity back and forth and the capability table would stop meaning
// anything.
func (s *Store) RoutableModels() []ModelSource {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[string]bool, len(s.own.Models))
	out := make([]ModelSource, 0, len(s.own.Models))
	for _, m := range s.own.Models {
		if !seen[m] {
			seen[m] = true
			out = append(out, ModelSource{Model: m})
		}
	}

	// Walk peers in a fixed order: map iteration is randomized, and a model list
	// that reshuffles between calls makes a client's picker jump around.
	endpoints := make([]string, 0, len(s.peers))
	for endpoint := range s.peers {
		endpoints = append(endpoints, endpoint)
	}
	sort.Strings(endpoints)

	for _, endpoint := range endpoints {
		for _, m := range s.peers[endpoint].Models {
			if !seen[m] {
				seen[m] = true
				out = append(out, ModelSource{Model: m, Endpoint: endpoint})
			}
		}
	}
	return out
}

// LocalCanServe reports whether this node's own model server has the model.
func (s *Store) LocalCanServe(model string) bool {
	if model == "" {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return contains(s.own.Models, model)
}

// PeersForModel returns the endpoints of peers whose cached manifest lists the
// model and that aren't marked busy. This is the capability lookup the router
// uses to pick an overflow target (M3).
func (s *Store) PeersForModel(model string) []string {
	if model == "" {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []string
	for endpoint, m := range s.peers {
		if !m.Busy && contains(m.Models, model) {
			out = append(out, endpoint)
		}
	}
	return out
}

func contains(models []string, model string) bool {
	for _, m := range models {
		if m == model {
			return true
		}
	}
	return false
}
