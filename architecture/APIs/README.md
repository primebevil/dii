# APIs

Status: describes the node interface as built, 2026-07-30 (functional node M1). This
was a placeholder until M1; the surface below is implemented in `prototype/` and
exercised on a live pod, so it is a record of what exists rather than a proposal.
Federation interfaces beyond the static peer list are still unwritten.

![Inside a functional node](../../diagrams/Functional_Node_M1.svg)

## The design rule

The node speaks the standard OpenAI HTTP API and nothing bespoke, in both
directions: the endpoint it serves, the calls it makes to a model server, and the
calls it makes to peers (ADR-0011). Two consequences follow, and both are load-bearing
rather than conveniences.

Any OpenAI-compatible client is already a DII client. A consumer needs a base URL
and a token, which every such client already has fields for, so the consumer role
needs no bespoke software (ADR-0004). A web GUI is not an integration either — it is
just another consumer.

And any OpenAI-compatible server is already a usable backend. The node reaches its
model server through one narrow interface using only the common subset, so
llama.cpp, vLLM, LM Studio, or LocalAI drop in behind a URL. Leaning on one
server's proprietary extras is exactly what would weld the project to a single
backend, so the node does not do it.

## Node endpoints

| Endpoint | Purpose |
|---|---|
| `POST /v1/chat/completions` | Chat, streaming (SSE) or single-response. The main serving path. |
| `POST /v1/embeddings` | Embeddings, always a single request and response. Completes the reliable-floor bundle (ADR-0006). |
| `GET /v1/models` | Everything this node can get served — its own models plus its peers' — in the OpenAI shape. See below. |
| `GET /manifest` | This node's self-description, for peers to fetch. Not an OpenAI route. |
| `GET /healthz` | Liveness: the process is up and answering. |
| `GET /readyz` | Readiness: the model server is reachable too. `503` when it is not. |

### The two doors

The bearer token selects which of two doors a request came in, and the router
treats them differently (ADR-0004):

- **No token — the local door.** The node's own trusted user. Local-first, then
  overflow to a capable peer.
- **A token matching a known consumer — the consumer door.** Someone with no node
  of their own. Peers first, keeping the owner's node free, then local.
- **Anything else — `401`.**

Today the consumer door accepts one shared `consumer_token`. Per-consumer
credentials that can be issued and revoked independently are M2 (ADR-0016's
non-gated sponsored path).

### What `/v1/models` advertises, and why it is the pod

`/v1/models` answers a *caller's* question — "what can I ask you for?" — so it lists
the whole pod: this node's own models plus any its cached peers hold, deduped with
the local copy winning. Each entry's `owned_by` carries where it would actually run,
`local` or the peer's endpoint. Clients ignore that field; it means a single call
explains the pod's shape.

This was not the original behavior and the change came out of use rather than
design review. The node would happily overflow a model it did not hold, but never
said so, so `curl` worked if you already knew the model's name while a GUI — which
builds its picker from this list — could not discover overflow at all. That made
ADR-0004's claim that any OpenAI-compatible client is already a DII client true for
chat mechanics and false for discovery. For a consumer with no node of their own,
who ADR-0015 puts at the center, the pod *is* the unit of service; which member
holds a model is an implementation detail they have no way to learn.

`/manifest` deliberately did **not** change: it answers a *peer's* question, "what
do you yourself hold?" If nodes republished models they had merely borrowed, peers
would advertise each other's capacity back and forth and the capability table would
stop meaning anything. Two questions, two answers.

The cost is honesty about staleness. Manifests are fetched once at startup, so this
list can name a model whose holder has since gone down, and that request fails as a
`502` from the dead peer rather than an immediate `503`. It is the same staleness
the router already carried, now visible in one more place — and the clearest
argument for a periodic manifest refresh.

### Capability matching

An OpenAI request already names a model, so "can I serve this?" reduces to "is
that model in my manifest?" That single model-name match covers both chat and
embeddings, because the standard `/v1/models` call the manifest is built from
reports every model a backend serves and carries no type or capability field. An
embedding model is therefore just a name in the same list. Richer capability tags
(chat, code, vision, quality tiers) are a later layer.

### Liveness versus readiness

Worth stating plainly, because having both only pays off if the difference is
understood. `/healthz` answers "is the process alive?" `/readyz` answers "can this
node actually serve?" by pinging the model server through the same backend
interface. A node whose Ollama has stopped answers `200` on the first and `503` on
the second. Without that split, a live-but-degraded node looks healthy to every
supervisor and directory that asks.

## Inter-node calls

A node calls a peer exactly as any OpenAI client would, on the peer's own
`/v1/chat/completions` or `/v1/embeddings`, with no consumer token — so the peer
serves the request on its own local door and does not forward again. That is also
what lets a peer be treated as just another backend behind the same interface.

Peers exchange `/manifest` once at startup and cache the result. There is no
gossip, no discovery, and no re-fetch yet; a node is restarted to pick up a newly
pulled model or a peer that came up late.

Two things this transport deliberately does not carry, both dropping out of the
load-bearing path now that overflow is optional (ADR-0015): node-to-node
authentication, and caller attribution across hops. They come back only if overflow
between untrusted pods is ever turned on.

## Failure responses

Failures are stated rather than hidden, which is the graceful-degradation principle
made concrete at the API surface:

- **`503`, immediately**, when no node in the pod has the model — returned before
  any upstream call, never as a hang or a timeout.
- **`502`** when a backend we did select failed, carrying the upstream error text.
- **`401`** for an unknown or missing token on the consumer door.

## Accounting

Every request writes one JSON line: door, opaque consumer handle, endpoint, model,
which node served it, status, latency, time to first token, and token counts. It
never carries prompt or completion text. That separation is the data-minimization
seam ADR-0016 describes — the pod holds token to handle to counters, and never
joins the content stream to the identity stream.

Token counts are exact when the response reports usage (any non-streaming
response, or a stream where the caller asked for `stream_options.include_usage`).
On a default stream no OpenAI-compatible server reports usage at all, so the count
falls back to a content-chunk proxy and the line is marked `tokens_estimated`. The
node does not silently inject `stream_options` to force exact numbers, because it
promises to forward request bodies verbatim. Closing that gap properly is open work
for when M3's token budgets need precision.

## Not yet specified

- An admin API for issuing, listing, and revoking consumer credentials (M2).
- Fair-use limits and the `429` backpressure shape (M3).
- A moderation refusal shape on the consumer door (M4).
- Any federation interface beyond the static peer list: discovery, signed peer
  lists, and the directory (Phase 3).
