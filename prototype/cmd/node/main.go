package main

import (
	"flag"
	"log"
	"net/http"

	"dii/internal/config"
	"dii/internal/ingress"
	"dii/internal/modelserver"
	"dii/internal/peer"
	"dii/internal/router"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to the node's YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// M1: no real inference. The backend is a mock that streams canned OpenAI
	// chunks. In M2 this becomes the real OpenAI client to Ollama. The seam is
	// modelserver.Backend; nothing above the interface depends on the mock.
	var backend modelserver.Backend = modelserver.NewMock(cfg.NodeID)

	// A peer node is, from our side, just another OpenAI-compatible backend: we
	// call its /v1/chat/completions exactly like any client.
	var peers []modelserver.Backend
	for _, endpoint := range cfg.Peers {
		peers = append(peers, peer.NewClient(endpoint))
	}

	rt := router.New(backend, peers)
	srv := ingress.New(cfg, rt)

	log.Printf("node %s listening on %s (peers: %v)", cfg.NodeID, cfg.Listen, cfg.Peers)
	if err := http.ListenAndServe(cfg.Listen, srv.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
