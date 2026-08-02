// Package usage is the node's per-request accounting: one structured JSON line
// per request, saying which door it came in, which node served it, and how many
// tokens it cost.
//
// It is deliberately dual-purpose. It is the substrate M3 fair-use will read its
// budgets from, and it is what offload sizing needs today — what got routed,
// where, which model, how big. Records key on an opaque consumer handle and never
// carry prompt or completion text, which is the data-minimization seam ADR-0016
// describes: the pod holds token to handle to counters, and never joins the
// content stream to the identity stream.
package usage

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

// Record is one request's accounting line.
type Record struct {
	Time     time.Time `json:"time"`
	NodeID   string    `json:"node_id"`
	Door     string    `json:"door"`     // "local" | "consumer"; empty when the token was rejected
	Consumer string    `json:"consumer"` // opaque handle; always empty for the owner
	Endpoint string    `json:"endpoint"` // "/v1/chat/completions" | "/v1/embeddings"
	Model    string    `json:"model"`
	ServedBy string    `json:"served_by"` // "local" | the peer endpoint; empty if nothing served it

	Status    int   `json:"status"`
	LatencyMS int64 `json:"latency_ms"`
	TTFTMS    int64 `json:"ttft_ms,omitempty"` // first byte from upstream; streaming only

	PromptTokens     int  `json:"prompt_tokens"`
	CompletionTokens int  `json:"completion_tokens"`
	Estimated        bool `json:"tokens_estimated"` // true => counts are a proxy, see Meter

	Error string `json:"error,omitempty"`
}

// Recorder writes Records as one JSON object per line. Safe for concurrent use
// by multiple request handlers.
type Recorder struct {
	nodeID string

	mu  sync.Mutex
	enc *json.Encoder
}

// NewRecorder writes records to w, stamping each with nodeID. Writing to stdout
// is the intended default: a JSON line per request is what both a container log
// driver and a human reading `docker logs` want.
func NewRecorder(w io.Writer, nodeID string) *Recorder {
	return &Recorder{nodeID: nodeID, enc: json.NewEncoder(w)}
}

// Write emits one record, filling in the node id and a timestamp. Accounting must
// never take down a request, so a write error is dropped rather than returned.
func (r *Recorder) Write(rec Record) {
	rec.NodeID = r.nodeID
	if rec.Time.IsZero() {
		rec.Time = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	_ = r.enc.Encode(rec)
}
