package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dii/internal/config"
	"dii/internal/ingress"
	"dii/internal/manifest"
	"dii/internal/modelserver"
	"dii/internal/peer"
	"dii/internal/router"
	"dii/internal/usage"
)

// healthcheckTimeout bounds the self-probe behind -healthcheck. It has to exceed
// the readiness handler's own backend ping (2s) to distinguish "backend slow" from
// "node not answering at all".
const healthcheckTimeout = 5 * time.Second

func main() {
	configPath := flag.String("config", envOr("DII_CONFIG", "config.yaml"),
		"path to the node's YAML config file (or set DII_CONFIG)")
	healthcheck := flag.Bool("healthcheck", false,
		"probe this node's own /readyz, print the result, and exit 0 (ready) or 1 (not ready)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Container health checks need something to run, and the node's image is
	// distroless: no shell, no curl. So the binary can probe itself.
	if *healthcheck {
		os.Exit(probeReady(cfg))
	}

	// M2: real inference. The local backend is an OpenAI-compatible client to
	// the model server (Ollama). The seam is modelserver.Backend, unchanged
	// above this line since M1.
	var backend modelserver.Backend = modelserver.NewClient(cfg.ModelServer, cfg.ResponseHeaderTimeout)

	// Build our own manifest by asking the model server what it can serve. If
	// the model server is down we still start (empty model list) so the node
	// and manifest exchange stay up; local serving returns 503 until it's back,
	// and /readyz reports not-ready in the meantime.
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

	// One JSON accounting line per request on stdout, which is what a container
	// log driver and a human reading the logs both want.
	recorder := usage.NewRecorder(os.Stdout, cfg.NodeID)

	rt := router.New(backend, peers, store)
	srv := ingress.New(cfg, rt, store, backend, recorder)

	// requestCtx is the parent of every request context (via BaseContext), so the
	// hard deadline below can cancel streams that outlive the graceful drain.
	requestCtx, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()

	httpSrv := &http.Server{
		Addr:        cfg.Listen,
		Handler:     srv.Handler(),
		BaseContext: func(net.Listener) context.Context { return requestCtx },
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("node %s listening on %s (peers: %v)", cfg.NodeID, cfg.Listen, cfg.Peers)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	select {
	case err := <-serveErr:
		log.Fatalf("server: %v", err)
	case <-signalCtx.Done():
		// Restore default signal handling, so a second Ctrl-C or SIGTERM during a
		// slow drain kills the process immediately instead of being swallowed.
		stopSignals()
	}

	// A zero shutdown_timeout means "no deadline", matching what the other timeouts
	// in the config mean by 0. Taken literally as a deadline it would be a deadline
	// already in the past, which would sever every in-flight stream — the exact
	// opposite of what someone writing "0s" is asking for. A supervisor still bounds
	// it (Docker's stop_grace_period, systemd's TimeoutStopSec), and a second signal
	// still exits immediately, so an unbounded drain cannot hang forever.
	shutdownCtx := context.Background()
	cancelShutdown := context.CancelFunc(func() {})
	if cfg.ShutdownTimeout > 0 {
		log.Printf("node %s: shutting down, draining in-flight requests for up to %s", cfg.NodeID, cfg.ShutdownTimeout)
		shutdownCtx, cancelShutdown = context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	} else {
		log.Printf("node %s: shutting down, draining in-flight requests with no deadline", cfg.NodeID)
	}
	defer cancelShutdown()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		// Drain deadline hit: cancel the request contexts so in-flight streams end
		// rather than hang, then close the listener and remaining connections.
		log.Printf("node %s: drain deadline exceeded (%v); cancelling in-flight requests", cfg.NodeID, err)
		cancelRequests()
		_ = httpSrv.Close()
	}
	log.Printf("node %s: stopped", cfg.NodeID)
}

// envOr lets the config path come from the environment, which is how the
// container image and a systemd unit point at a mounted config without depending
// on the working directory.
func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// probeReady asks a locally running node whether it is ready and returns the
// process exit code: 0 ready, 1 not. It derives the URL from the same config the
// node listens with, so there is nothing extra to keep in sync.
func probeReady(cfg *config.Config) int {
	host, port, err := net.SplitHostPort(cfg.Listen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: cannot parse listen address %q: %v\n", cfg.Listen, err)
		return 1
	}
	// A wildcard bind is not an address to dial; talk to the loopback instead.
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}

	url := "http://" + net.JoinHostPort(host, port) + "/readyz"
	client := &http.Client{Timeout: healthcheckTimeout}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	fmt.Fprintf(os.Stderr, "healthcheck: %s -> %d %s\n", url, resp.StatusCode, bytes.TrimSpace(body))
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
