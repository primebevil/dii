package ingress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dii/internal/config"
	"dii/internal/manifest"
	"dii/internal/modelserver"
	"dii/internal/router"
	"dii/internal/usage"
)

// readinessTimeout bounds the backend ping behind /readyz. The standard
// /v1/models call is cheap and does not load a model, so a readiness probe should
// not need a config knob to wait longer than this.
const readinessTimeout = 2 * time.Second

const (
	routeChat       = "/v1/chat/completions"
	routeEmbeddings = "/v1/embeddings"
)

var errUnauthorized = errors.New("unauthorized")

// Server is the node's OpenAI-compatible front door. The token check turns each
// request into one of the two doors, then the router decides where it runs.
type Server struct {
	cfg     *config.Config
	router  *router.Router
	store   *manifest.Store
	backend modelserver.Backend // pinged by /readyz, through the same seam
	usage   *usage.Recorder
}

func New(cfg *config.Config, rt *router.Router, store *manifest.Store, backend modelserver.Backend, rec *usage.Recorder) *Server {
	return &Server{cfg: cfg, router: rt, store: store, backend: backend, usage: rec}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/chat/completions", s.handleChat)
	mux.HandleFunc("POST /v1/embeddings", s.handleEmbeddings)
	mux.HandleFunc("GET /v1/models", s.handleModels)
	mux.HandleFunc("GET /manifest", s.handleManifest)
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /readyz", s.handleReady)
	return mux
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := usage.Record{Endpoint: routeChat}
	status := http.StatusOK

	// Exactly one accounting line per request, on every exit path.
	defer func() {
		rec.Status = status
		rec.LatencyMS = time.Since(start).Milliseconds()
		s.usage.Write(rec)
	}()

	door, err := s.door(r)
	if err != nil {
		status, rec.Error = http.StatusUnauthorized, err.Error()
		writeError(w, status, "invalid or missing token for this door")
		return
	}
	rec.Door, rec.Consumer = doorName(door), consumerHandle(door)

	body, err := io.ReadAll(io.LimitReader(r.Body, s.cfg.MaxBodyBytes))
	if err != nil {
		status, rec.Error = http.StatusBadRequest, err.Error()
		writeError(w, status, "could not read request body")
		return
	}

	// Best-effort peek; the body is still forwarded verbatim. Model drives
	// capability matching, and stream tells us the response shape — OpenAI
	// defaults it to false, in which case the reply is one JSON object, not SSE.
	var probe struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
	}
	_ = json.Unmarshal(body, &probe)
	rec.Model = probe.Model

	stream, decision, err := s.router.Route(r.Context(), door, modelserver.ChatRequest{
		Model:  probe.Model,
		Stream: probe.Stream,
		Body:   body,
	})
	rec.ServedBy = decision.ServedBy
	if err != nil {
		status, rec.Error = upstreamStatus(err), err.Error()
		writeError(w, status, err.Error())
		return
	}

	meter := usage.NewMeter(stream, probe.Stream, start)
	relay(w, meter, probe.Stream)
	rec.PromptTokens, rec.CompletionTokens, rec.Estimated = meter.Usage()
	rec.TTFTMS = meter.TTFT().Milliseconds()
}

// handleEmbeddings completes the reliable-floor bundle (ADR-0006). Embeddings are
// a single request/response with no streaming, so the whole body is brokered in
// one piece and its usage block gives exact token counts.
func (s *Server) handleEmbeddings(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := usage.Record{Endpoint: routeEmbeddings}
	status := http.StatusOK

	defer func() {
		rec.Status = status
		rec.LatencyMS = time.Since(start).Milliseconds()
		s.usage.Write(rec)
	}()

	door, err := s.door(r)
	if err != nil {
		status, rec.Error = http.StatusUnauthorized, err.Error()
		writeError(w, status, "invalid or missing token for this door")
		return
	}
	rec.Door, rec.Consumer = doorName(door), consumerHandle(door)

	body, err := io.ReadAll(io.LimitReader(r.Body, s.cfg.MaxBodyBytes))
	if err != nil {
		status, rec.Error = http.StatusBadRequest, err.Error()
		writeError(w, status, "could not read request body")
		return
	}

	var probe struct {
		Model string `json:"model"`
	}
	_ = json.Unmarshal(body, &probe)
	rec.Model = probe.Model

	out, decision, err := s.router.RouteEmbeddings(r.Context(), door, modelserver.EmbeddingsRequest{
		Model: probe.Model,
		Body:  body,
	})
	rec.ServedBy = decision.ServedBy
	if err != nil {
		status, rec.Error = upstreamStatus(err), err.Error()
		writeError(w, status, err.Error())
		return
	}
	rec.PromptTokens, rec.CompletionTokens, rec.Estimated = usage.FromResponse(out)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

// handleManifest publishes this node's self-description for peers to fetch.
func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.store.Own())
}

// handleModels answers "what can I ask this node for?" — which is the whole pod's
// capability, this node's own models plus any its cached peers hold, not just what
// sits on this node's shelf. A caller cannot discover overflow any other way: an
// OpenAI client builds its model picker from this list, so a model the node would
// happily borrow is simply invisible unless the node admits to it.
//
// owned_by carries where each model would actually run. Clients ignore the field,
// but it makes `curl /v1/models` explain the pod without a second call.
//
// Embedding models sit in this list alongside chat models: the standard /v1/models
// call the manifest is built from reports every model a backend serves and carries
// no type field, so capability stays a single model-name match for both endpoints.
func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	type model struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	}
	out := struct {
		Object string  `json:"object"`
		Data   []model `json:"data"`
	}{Object: "list"}

	for _, m := range s.store.RoutableModels() {
		servedBy := router.ServedLocal
		if m.Endpoint != "" {
			servedBy = m.Endpoint
		}
		out.Data = append(out.Data, model{ID: m.Model, Object: "model", OwnedBy: servedBy})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// handleHealth is liveness: the process is up and answering.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "node %s ok\n", s.cfg.NodeID)
}

// handleReady is readiness: can this node actually serve? It pings the model
// server through the Backend seam, so a live-but-degraded node — process up,
// backend down — answers 503 and a supervisor, a container health check, or the
// owner can tell the two apart.
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")
	if _, err := s.backend.ListModels(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"node_id": s.cfg.NodeID,
			"ready":   false,
			"reason":  "model server unreachable: " + err.Error(),
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"node_id": s.cfg.NodeID,
		"ready":   true,
	})
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

// relay sends the chosen backend's response through to the caller, bytes
// untouched. A streaming request gets SSE, flushed after every read so tokens
// arrive as they are produced; a non-streaming request gets the single JSON object
// it asked for, correctly labelled — the POC announced every response as SSE,
// which only confused strict clients on the non-streaming path.
func relay(w http.ResponseWriter, body io.ReadCloser, stream bool) {
	defer body.Close()

	if stream {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(http.StatusOK)

	flusher, canFlush := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := body.Read(buf)
		if n > 0 {
			if _, werr := w.Write(buf[:n]); werr != nil {
				return
			}
			if canFlush && stream {
				flusher.Flush()
			}
		}
		if err != nil {
			return
		}
	}
}

// upstreamStatus maps a routing failure to the status the caller sees: no node in
// the pod has the model is an honest, immediate 503; anything else failed talking
// to a backend we did pick.
func upstreamStatus(err error) int {
	if errors.Is(err, router.ErrNoCapacity) {
		return http.StatusServiceUnavailable
	}
	return http.StatusBadGateway
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

// consumerHandle is the opaque identity the usage log keys on. The owner is
// always empty. There are no per-consumer credentials yet, so every request
// through the one shared consumer_token logs as "shared"; M2 replaces this with a
// real handle per issued token, which is why the field exists now.
func consumerHandle(d router.Door) string {
	if d == router.DoorConsumer {
		return "shared"
	}
	return ""
}
