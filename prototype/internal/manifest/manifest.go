// Package manifest is how the router answers "who can serve this?". Each node
// publishes a small self-description built from its model server's model list,
// and caches its peers' descriptions. Capability match is model-name match (the
// POC shortcut); richer capability tags are a later layer.
package manifest

import "sync"

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
