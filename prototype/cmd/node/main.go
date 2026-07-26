package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"dii/internal/config"
	"dii/internal/ingress"
	"dii/internal/manifest"
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

	// M2: real inference. The local backend is an OpenAI-compatible client to
	// the model server (Ollama). The seam is modelserver.Backend, unchanged
	// above this line since M1.
	var backend modelserver.Backend = modelserver.NewClient(cfg.ModelServer, cfg.ResponseHeaderTimeout)

	// Build our own manifest by asking the model server what it can serve. If
	// the model server is down we still start (empty model list) so the node
	// and manifest exchange stay up; local serving returns 503 until it's back.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.StartupTimeout)
	models, err := backend.ListModels(ctx)
	cancel()
	if err != nil {
		log.Printf("node %s: model server %s unavailable (%v); starting with no local models", cfg.NodeID, cfg.ModelServer, err)
		models = nil
	}
	store := manifest.NewStore(manifest.Manifest{
		NodeID:   cfg.NodeID,
		Endpoint: cfg.Advertise,
		Models:   models,
	})
	log.Printf("node %s: local models %v", cfg.NodeID, models)

	// A peer is just another OpenAI-compatible backend for serving, plus a
	// /manifest endpoint we fetch once at startup. We pair each peer with its
	// endpoint so the router can map a manifest hit back to the client to call.
	var peers []router.Peer
	for _, endpoint := range cfg.Peers {
		pc := peer.NewClient(endpoint, cfg.ResponseHeaderTimeout)
		peers = append(peers, router.Peer{Endpoint: pc.Endpoint(), Backend: pc})

		pctx, pcancel := context.WithTimeout(context.Background(), cfg.StartupTimeout)
		m, err := pc.FetchManifest(pctx)
		pcancel()
		if err != nil {
			log.Printf("node %s: could not fetch manifest from peer %s (%v); serving without it", cfg.NodeID, pc.Endpoint(), err)
			continue
		}
		store.SetPeer(pc.Endpoint(), m)
		log.Printf("node %s: peer %s serves %v", cfg.NodeID, pc.Endpoint(), m.Models)
	}

	rt := router.New(backend, peers, store)
	srv := ingress.New(cfg, rt, store)

	log.Printf("node %s listening on %s (peers: %v)", cfg.NodeID, cfg.Listen, cfg.Peers)
	if err := http.ListenAndServe(cfg.Listen, srv.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
