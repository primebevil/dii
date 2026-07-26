package ingress

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"dii/internal/config"
	"dii/internal/manifest"
	"dii/internal/modelserver"
	"dii/internal/router"
)

var errUnauthorized = errors.New("unauthorized")

// Server is the node's OpenAI-compatible front door. The token check turns each
// request into one of the two doors, then the router decides where it runs.
type Server struct {
	cfg    *config.Config
	router *router.Router
	store  *manifest.Store
}

func New(cfg *config.Config, rt *router.Router, store *manifest.Store) *Server {
	return &Server{cfg: cfg, router: rt, store: store}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/chat/completions", s.handleChat)
	mux.HandleFunc("GET /v1/models", s.handleModels)
	mux.HandleFunc("GET /manifest", s.handleManifest)
	mux.HandleFunc("GET /healthz", s.handleHealth)
	return mux
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	door, err := s.door(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or missing token for this door")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, s.cfg.MaxBodyBytes))
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read request body")
		return
	}

	// Best-effort peek at the model name; the body is still forwarded verbatim.
	var probe struct {
		Model string `json:"model"`
	}
	_ = json.Unmarshal(body, &probe)

	log.Printf("node %s: %s door, model=%q", s.cfg.NodeID, doorName(door), probe.Model)

	stream, err := s.router.Route(r.Context(), door, modelserver.ChatRequest{Model: probe.Model, Body: body})
	if err != nil {
		if errors.Is(err, router.ErrNoCapacity) {
			writeError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	streamSSE(w, stream)
}

// handleManifest publishes this node's self-description for peers to fetch.
func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.store.Own())
}

// handleModels exposes the node's own models in the OpenAI /v1/models shape,
// served straight from the cached manifest.
func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	type model struct {
		ID     string `json:"id"`
		Object string `json:"object"`
	}
	out := struct {
		Object string  `json:"object"`
		Data   []model `json:"data"`
	}{Object: "list"}
	for _, id := range s.store.Own().Models {
		out.Data = append(out.Data, model{ID: id, Object: "model"})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "node %s ok\n", s.cfg.NodeID)
}

// door maps the bearer token to a door: no token is the trusted local owner, a
// token matching consumer_token is the consumer door, anything else is rejected.
func (s *Server) door(r *http.Request) (router.Door, error) {
	switch token := bearerToken(r); {
	case token == "":
		return router.DoorLocal, nil
	case s.cfg.ConsumerToken != "" && token == s.cfg.ConsumerToken:
		return router.DoorConsumer, nil
	default:
		return 0, errUnauthorized
	}
}

// streamSSE relays the backend's raw SSE bytes to the client, flushing after
// each read so tokens reach the caller as they arrive rather than buffered.
func streamSSE(w http.ResponseWriter, body io.ReadCloser) {
	defer body.Close()
	flusher, canFlush := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	buf := make([]byte, 4096)
	for {
		n, err := body.Read(buf)
		if n > 0 {
			if _, werr := w.Write(buf[:n]); werr != nil {
				return
			}
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			return
		}
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"message": msg,
			"type":    "dii_error",
		},
	})
}

func bearerToken(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[len("bearer "):])
	}
	return h
}

func doorName(d router.Door) string {
	if d == router.DoorConsumer {
		return "consumer"
	}
	return "local"
}
