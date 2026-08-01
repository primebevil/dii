# Functional Node Plan

Status: Plan, 2026-07-30 (revised same day to add the Packaging section and M5, optional escalation fallback). The build target after Week 4. Depends on ADR-0004
(consumer as first-class ingress), ADR-0006 (reliable floor bundle), ADR-0013
(pod is the accountability boundary), ADR-0014 (consumer work preemptible,
per-identity fair-use), ADR-0015 (access distribution, overflow optional), and
ADR-0016 (consumer admission and identity, direction with gates). It continues
docs/Week_3_Prototype_Plan.md and prototype/BUILD_BRIEF.md, which built and
validated the skeleton. This plan promotes `prototype/` from a disposable proof
of concept to the functional node, so it is no longer throwaway.

This is a handoff document. It is written to be built from directly in Claude
Code with a tight run-test loop against a real Ollama on atlas.

## Why this exists

Week 4 reframed the project: access distribution is the thesis, the consumer
ingress is the load-bearing primitive, and node-to-node overflow is an optional
optimization (ADR-0015). The prototype proved the skeleton but, by its own
BUILD_BRIEF non-goals, skipped everything that makes a node functional rather
than demonstrative: real identity, fair-use, moderation, and operational
durability. This plan closes that gap.

It serves two goals at once, which is deliberate. It is the DII reference node,
and it is the owner's personal daily driver on atlas (pod-zero). The early
milestones make the daily driver dependable; the later ones add the
public-serving capability the mission needs. Building one artifact for both is
the point.

## What exists today (the starting point)

The node in `prototype/` is a control plane in front of a local
OpenAI-compatible model server (Ollama). It has:

- An OpenAI-compatible front door: `POST /v1/chat/completions`, `GET /v1/models`,
  `GET /manifest`, `GET /healthz`.
- Two doors selected by bearer token: no token is the trusted owner (local door),
  a token equal to the single `consumer_token` is a guest (consumer door),
  anything else is 401.
- Routing: owner is local-first then overflow to a capable peer; guest is
  peer-first then local; an unservable model returns an honest 503
  (`ErrNoCapacity`).
- Capability matching by model-name against cached peer manifests.
- A `modelserver.Backend` interface over the standard OpenAI subset
  (`/v1/chat/completions` streaming, `/v1/models`); Ollama is just a base URL.
- SSE streaming passthrough end to end.
- Config: `listen`, `node_id`, `model_server`, `consumer_token`, `advertise`,
  `peers`, and operational tunables (`startup_timeout`,
  `response_header_timeout`, `max_body_bytes`).

What it deliberately lacks (BUILD_BRIEF non-goals): real identity, fair-use,
moderation, persistence, lifecycle hardening, and the embeddings half of the
floor bundle.

## What "functional" means (definition of done for the whole plan)

A functional node:

- serves the reliable-floor bundle locally and reliably enough to be a daily
  driver (general chat, a coder, and embeddings), with graceful lifecycle and
  visibility into what it served;
- admits and revokes individual consumers, each with their own credential, not a
  shared secret;
- bounds each consumer with fair-use limits and keeps the owner's own work from
  being starved;
- has a moderation seam at the consumer door, so it is structurally ready to
  serve strangers once the ADR-0016 legal gate clears;
- degrades honestly and never hangs.

Cross-pod and federated identity, discovery, and a real policy engine are out of
scope here (see Deferred).

## Hard constraints (KISS, carried from BUILD_BRIEF)

- Control plane only. The node never runs inference; it brokers to the model
  server.
- OpenAI-compatible on both sides: the node's own endpoint and every call it
  makes.
- Backend access stays behind `modelserver.Backend`, standard OpenAI subset only.
  No Ollama-native calls.
- Inter-node calls reuse the OpenAI HTTP call (ADR-0011). No bespoke RPC yet.
- One binary is a node. Keep the single-binary, config-per-node shape.
- Stream tokens through end to end (SSE passthrough).
- Capability match stays model-name match.
- Prefer the standard library; add a dependency only with a clear reason (a small
  embedded store is the one likely exception, see M2).
- Optimize for reviewability and reversibility: small milestones, each landing
  behind a flag or default that leaves existing behavior intact until turned on.
  The owner reviews and handles all git commits.
- Keep `prototype/` as the location for now; do not rename the directory or the
  `dii/` module path in this plan, since that churns imports for no functional
  gain. Treat "the node" and `prototype/` as the same thing.

## Milestones

Build one milestone at a time and stop for review between each. Each milestone
should leave the node runnable and the existing owner/consumer/overflow behavior
intact.

Two deployment paths, checked at different rates. Day to day the pod runs on
containers, because that is the lowest-friction way to iterate: atlas as the
serving node and the laptop as a spoke. But the node must also run as a plain
supervised service with no container runtime at all — that is a real way people
will run one, and it is the only thing keeping the node honest about being a single
static binary that depends on nothing.

So at each milestone boundary, before review, run the node under systemd once and
put `prototype/scripts/m1-check.sh` against it. It takes a few minutes and it
catches the class of regression a container hides: a new dependency on a mounted
path, an env var only compose sets, a signal only Docker delivers. Steps are in
`prototype/deploy/README.md`. Whichever path is not in daily use is the one that
rots, which is exactly why this is a gate and not a good intention.

### M1 — Floor bundle and reliability

Status: built and verified 2026-07-30, awaiting review. All five tasks landed, plus
one addition the build forced (see below). Running as containers on atlas and the
laptop; `prototype/scripts/m1-check.sh` covers the definition of done and passes
12/12 against both.

Addition not in the task list: `/v1/models` now advertises the whole pod rather
than only the node's own models. A node would otherwise overflow models it never
advertised, leaving overflow undiscoverable to any client that builds a model
picker from that list. Recorded in architecture/APIs/README.md.

Goal: a node the owner can run on atlas as a dependable daily driver and see what
it is doing.

Tasks:

- Add `POST /v1/embeddings` passthrough through the `modelserver.Backend`
  interface, so the floor bundle (general, coder, embeddings per ADR-0006) is
  complete. Extend the manifest and `/v1/models` to include embedding models the
  backend reports. Embeddings are a single request/response (no streaming), so
  handle the non-streaming path cleanly.
- Graceful lifecycle: catch SIGINT/SIGTERM and `http.Server.Shutdown` with a
  drain deadline so in-flight streams finish or cancel cleanly.
- Readiness: keep `/healthz` as liveness; add `/readyz` that pings the backend
  `/v1/models` and reports not-ready (503) when the backend is unreachable, so a
  supervisor and the owner can tell a live-but-degraded node apart.
- Observability: emit one structured (JSON) log line per request with timestamp,
  node_id, door, consumer handle (empty for owner), model, status, latency, and
  prompt/completion token counts when the backend reports usage. Best-effort on
  tokens: read usage from a non-streaming response, or from the final SSE chunk
  if present, else count content deltas and mark it estimated.
- Deployment: a `systemd` unit for atlas (Linux) and a note on the launchd
  equivalent for macOS, plus a config-path flag/env so the service is not tied to
  the working directory. A container image and compose file are the packaging
  companion to this (see Packaging below); they can land with M1 or immediately
  after, and are the preferred way to run the node on atlas.

Design notes:

- Usage accounting here is dual-purpose: it is the substrate for M3 fair-use and
  it is exactly what the owner's Atlas offload sizing needs (what got routed,
  which model, how many tokens). Keep the per-request record in a shape that both
  a human and later code can read.

Definition of done: the node serves chat, coder, and embeddings against atlas's
Ollama; survives a restart cleanly under systemd; `/readyz` flips when the
backend is stopped; and every request produces a usage log line. Stop for review.

### M2 — Real consumer identity

Goal: replace the single shared secret with per-consumer credentials that can be
issued and revoked independently. This is ADR-0016's sponsored-consumer path,
which is not gated.

Tasks:

- Introduce a consumer store: a token maps to a `Consumer` record with a stable
  opaque `handle` (a random id), a human label, issued-at, and a revoked flag.
  Logs and counters key on the handle, never the raw token.
- Change the door logic: no token is the owner (unchanged); a token that matches a
  live (non-revoked) consumer opens the consumer door and carries that handle
  into routing and logging; an unknown or revoked token is 401. Killing a token
  cancels its in-flight requests.
- Management surface, kept minimal: issue, list (handles and labels, never raw
  tokens), and revoke. Prefer a CLI subcommand on the same binary, or an admin
  HTTP listener bound to localhost only. Do not expose consumer management on the
  public listener.
- Persistence: the store must survive restart. A single-file store (JSON with an
  atomic write, or SQLite if concurrent counter updates in M3 argue for it) is
  enough at single-node scale. Choose the simpler option that M3 can build on.

Design notes:

- The opaque handle is the seam ADR-0016 described: the pod holds token to handle
  to counters, and never joins the content stream to the identity stream. Do not
  log prompts against handles.
- Keep the old `consumer_token` working as a deprecated single-consumer shortcut
  if it eases migration, or migrate it into the store as one seeded consumer.
  Owner behavior (no token) must not change.

Definition of done: the owner can issue two distinct consumer tokens, watch each
appear in usage logs under its own handle, revoke one and see it 401 immediately
while the other keeps working, and restart the node without losing the store.
Stop for review.

### M3 — Fair-use and preemption

Goal: bound each consumer so one guest cannot drain the pod, and keep the owner's
own work first (ADR-0014).

Tasks:

- Per-consumer limits, checked at ingress before routing: max concurrent
  requests, a request rate (per minute), and a rolling token budget over a window,
  read from the M1 usage accounting. Defaults in config, per-consumer overrides in
  the store.
- Backpressure: over a limit returns HTTP 429 with `Retry-After`, which every
  OpenAI-compatible client already understands. Never silently drop or hang.
- Owner-first: the owner (local door) is never rate-limited by consumer budgets.
  Consumer work is best-effort. Because the node is control plane only, true
  mid-flight preemption of a running inference is limited by the backend, so the
  honest mechanism here is low consumer concurrency caps plus admission control
  (refuse or queue new consumer work when the backend is saturated by the owner),
  not interrupting a call already in Ollama. State that limit plainly rather than
  implying hard preemption.

Definition of done: a consumer looping requests hits its concurrency and rate
caps and gets 429s while the owner's requests always go through; budgets reset on
the window; limits are configurable per consumer. Stop for review.

### M4 — The safety seam

Goal: make the node structurally ready to serve strangers, with the moderation
seam in place, without turning on public serving (which the ADR-0016 legal gate
still blocks).

Tasks:

- A `Moderator` interface applied only on the consumer door: check input before
  routing and, where feasible, check output. Default implementation is a no-op
  that allows everything, so owner and trusted use are unaffected.
- A pluggable classifier implementation: call an OpenAI-compatible safety model or
  an external classifier endpoint named in config, fail closed on the consumer
  door if the classifier is configured but unreachable.
- Wire enforcement to M2: a request that fails moderation is refused and logged
  under its handle, and repeated failures can trigger revoke. Keep the classifier
  itself pluggable and deferred; build the seam and the no-op default now.

Design notes:

- This is the point where the ADR-0016 concerns become operational. The seam is
  built here; actually serving the unaffiliated public is still gated on the legal
  review and on denylist governance, neither of which this milestone closes.

Definition of done: with the no-op moderator the node behaves exactly as after
M3; with a classifier configured, a consumer request that trips it is refused with
a clear error and logged, and the owner path is never moderated. Stop for review.

### M5 — Optional escalation fallback

Goal: let an identity optionally escalate to a configured upstream when the pod
cannot serve a request, while leaving the default behavior (an honest 503) exactly
as it is when no fallback is configured. This is the node-side form of the owner's
local-first, escalate-to-paid offload pattern.

Tasks:

- Add an optional fallback as the last hop in the router: after local and peers
  both miss, if the requesting identity has an escalation target configured, try it
  through the same `modelserver.Backend` interface; if not, return `ErrNoCapacity`
  as today. No configured fallback means no change in behavior.
- Owner escalation: the owner door can name a personal upstream (an
  OpenAI-compatible endpoint plus credential, for example Claude via its
  OpenAI-compatible endpoint, OpenRouter, or another pod). Off by default.
- Consumer escalation: off by default and never a pod-wide default. It may be
  enabled only for a specific consumer that supplies its own upstream credential,
  and never for an anonymous or sovereignty-promised consumer, because escalation
  sends the prompt off-boundary to a third party.
- Off-boundary logging: whenever a request escalates to an external upstream, log
  it explicitly (identity or handle, model, upstream target) and mark it as a
  trust-boundary crossing, so it is auditable and never silent. Consider a response
  header noting the request was served by an upstream rather than the pod.

Design notes:

- Architecturally this is the stubbed policy gate made real for exactly one policy:
  on no-capacity, escalate this identity to its target. It reuses ADR-0012 (the
  fallback is just another OpenAI-compatible backend), so it is a small addition to
  the router, not a new subsystem.
- This is capacity fallback only: it triggers when no node can serve the requested
  model. Quality escalation (a node answered but the answer is not good enough) is a
  different, harder layer and is out of scope here (see Deferred).
- The honest 503 stays the default and the guarantee (ADR-0006, ADR-0014).
  Escalation is an opt-in override, never the pod's baseline promise, so the
  no-counterparty property holds for everyone who has not chosen otherwise.

Definition of done: with no escalation configured, behavior is identical to today,
an honest 503 when no node can serve. With an owner upstream configured, a request
for a model the pod cannot serve is answered by the upstream and the escalation is
logged as an off-boundary crossing. Consumer escalation stays off unless explicitly
configured for a named consumer with its own credential. Stop for review.

## Config evolution

The config grows additively; existing fields keep their meaning. Sketch:

```yaml
listen: ":8090"
node_id: "atlas"
model_server: "http://localhost:11434/v1"
advertise: "http://100.118.77.40:8090"
peers: []                      # overflow is optional now; empty is fine

# M2: consumer store location (tokens are managed via CLI/admin, not hand-edited)
consumer_store: "./consumers.db"
admin_listen: "127.0.0.1:8091" # localhost-only management surface

# M3: fair-use defaults, overridable per consumer in the store
fair_use:
  max_concurrent: 2
  requests_per_min: 60
  token_budget_per_hour: 200000

# M4: moderation (empty => no-op allow-all; owner door is never moderated)
moderation:
  classifier_url: ""           # OpenAI-compatible safety endpoint, optional
  fail_closed: true

# M5: optional escalation (empty => honest 503 on no-capacity, unchanged)
escalation:
  owner:
    upstream_url: ""           # an OpenAI-compatible endpoint (Claude/OpenRouter/another pod)
    upstream_token: ""
  # per-consumer escalation is off by default; enable it only in the consumer
  # store, and only for a consumer that brings its own upstream credential

# operational tunables (existing)
startup_timeout: "5s"
response_header_timeout: "30s"
max_body_bytes: 1048576
```

## Packaging: Docker and an optional GUI

Additive, not a milestone. It changes how the node is run, not what it does, so it
slots under M1's deployment work and can land with M1 or right after. It does not
alter any milestone's definition of done.

Containerize the node. The Go static binary drops into a minimal distroless or
scratch image, so the node ships as a small container and `docker compose up`
stands up a node without a from-scratch setup. This is the easiest on-ramp for
anyone running a node, which is directly the access-distribution supply problem of
lowering the barrier to run one, and it turns atlas into pull-and-run rather than
build-and-configure, with no per-change binary builds. One hard rule, straight
from the control-plane-only constraint: inference stays out of the node image.
Ollama runs where the GPU is, natively on the host, and the node container points
at it by URL (`http://host:11434/v1`, or `host.docker.internal` on Docker
Desktop). On atlas (Linux) you can later containerize Ollama with GPU passthrough
if you want uniformity; on macOS the GPU is not reachable from a container, so
Ollama stays native there. Keep the node image backend-agnostic; it already talks
to a URL.

An optional GUI as a consumer. Because the node is an OpenAI-compatible endpoint,
a web GUI such as Open WebUI is not an integration, it is just another consumer.
Point its OpenAI base URL at the node's `/v1` and give it a consumer token as the
API key, and it lists models and chats through the node with no node changes. Ship
it as an optional compose service alongside the node. Two notes: M1's embeddings
endpoint also enables the GUI's document and retrieval features, so there is real
synergy; and the GUI's own login is a separate layer above the node, distinct from
the per-consumer identity in M2, so behind one shared token every GUI user looks
like a single consumer to the node, which is fine for personal use and only needs
wiring if you later want GUI users mapped to distinct node consumers.

Suggested shape: a `Dockerfile` for the node, a `docker-compose.yaml` with the
node service and a commented-optional `open-webui` service, and a note in the
compose that Ollama is expected on the host. Keep these under `prototype/`
(a `prototype/deploy/` folder is fine) next to the code.

## Deferred (gated or later), with why

- Delegated federated login and the cross-pod denylist (ADR-0016). Gated: they
  need a delegated-admission prototype and a legal review before shipping. M2's
  local per-consumer credentials are the non-gated subset and are enough to run a
  node for yourself and trusted pod members.
- Node-to-node authentication and cross-hop caller attribution. Overflow is
  optional now (ADR-0015), so these drop out of the load-bearing path. Add them
  only if and when overflow between untrusted pods is turned on.
- Discovery and federation beyond the static peer list. A later phase.
- A real sovereignty and policy engine. Keep it a hook; do not build a policy DSL
  in this plan. M5 builds one concrete policy (escalation on no-capacity); the
  general engine stays deferred.
- Quality-based escalation, where a node answered but the answer is judged not good
  enough so the request escalates. Distinct from M5's capacity fallback, which
  fires only when no node can serve the model at all. Quality escalation needs a
  verify step the node cannot do alone (the router-and-verify layer from the
  offload discussion), and is later.

## The legal gate, stated plainly

Running this node for yourself and for trusted, affiliated members is fine now.
Serving the unaffiliated public is gated by the ADR-0016 mandatory legal review
(the denylist as regulated personal data, retention, subject rights) and by
denylist governance. M4 builds the seam so you are ready; it does not open the
door. Do not point the public directory at a node serving strangers until that
gate clears.

## Build order and a paste-ready kickoff for Claude Code

Do one milestone at a time, stop for review, keep the owner path unchanged. Start
with M1. Paste-ready kickoff:

```
Read docs/Functional_Node_Plan.md, prototype/BUILD_BRIEF.md, and
architecture/Sketchbook.md in this repo. We are promoting the DII node in
prototype/ from a proof of concept to a functional node, in Go, KISS, control
plane only, OpenAI-compatible on both sides, backend behind the
modelserver.Backend interface.

Do milestone M1 ONLY from the Functional Node Plan:
  1. Add POST /v1/embeddings passthrough via modelserver.Backend, and include
     embedding models in the manifest and /v1/models.
  2. Graceful shutdown on SIGINT/SIGTERM with a drain deadline.
  3. Add /readyz that pings the backend and reports not-ready when it is down;
     keep /healthz as liveness.
  4. Emit one structured JSON log line per request: node_id, door, consumer
     handle (empty for owner), model, status, latency, and prompt/completion
     tokens (best-effort, mark estimated when counted).
  5. A systemd unit for atlas and a config-path flag.

Do not touch identity, fair-use, or moderation yet (those are M2-M4). Keep the
owner (no token) and existing consumer/overflow behavior unchanged. First show me
the interface changes and the main.go wiring and stop for my review, then
implement. I handle all git commits.
```

## References

- architecture/Overview.md and architecture/Sketchbook.md (the node model).
- docs/Reliable_Floor_Definition.md (the bundle M1 completes).
- docs/ADR/ADR-0004, ADR-0013, ADR-0014, ADR-0015, ADR-0016.
- docs/Identity_Note_From_Prototype.md and docs/Governance_And_Abuse_Resistance.md
  (why M2-M4 are shaped as they are, and what stays gated).
- prototype/BUILD_BRIEF.md and docs/Week_3_Prototype_Plan.md (the skeleton this
  builds on).
