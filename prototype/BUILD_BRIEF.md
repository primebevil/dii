# Prototype Build Brief

Handoff doc for building the Week-3 prototype in Go. Read this first, then build in your editor with a tight run-test loop against a real Ollama.

## Read these first (sources of truth)

- architecture/Sketchbook.md, the settled mental model: model server vs. node, the router, two doors, the manifest, the full workflow, the Go package layout, and backend portability.
- docs/Week_3_Prototype_Plan.md, the scope, milestones, and kill-criteria measurements.

This brief does not restate those; it pins the constraints and gives a starting task.

## What we are building

A DII node in Go that sits in front of a local OpenAI-compatible model server (Ollama for now) and routes capability requests: serve locally when it can, hand off to a peer node when it can't, and serve a consumer that has no node of its own. Two nodes plus a consumer form the one-pod proof of concept.

## Hard constraints (KISS)

- The node is control plane only. It never runs inference; it brokers to the model server.
- OpenAI-compatible on both sides: the node's own endpoint, and the calls it makes to the model server and to peers.
- Backend access goes through a modelserver.Backend interface using only the standard OpenAI subset (/v1/chat/completions with streaming, /v1/models). No Ollama-native calls. Ollama is just a base URL behind the interface.
- Inter-node transport is the OpenAI HTTP call reused between nodes. Node A calls node B's /v1/chat/completions like any client. Do not build a bespoke RPC yet.
- One binary is a node. Two nodes is the same binary run twice with different config files.
- Stream tokens through end to end (SSE passthrough), for both the local and the peer path.
- Capability match is model-name match: a request names a model; a node can serve it if that model is in its manifest.

## Package layout

Per the sketchbook:

```
prototype/
  go.mod
  cmd/node/main.go        - load config, wire packages, start the node
  internal/
    config/               - parse config (listen addr, node_id, model_server URL, consumer_token, peers)
    ingress/              - OpenAI-compatible HTTP server; the two doors; token check selects the door
    router/               - decide where a request runs: local-first, else a peer, else honest fail
    modelserver/          - Backend interface (ListModels, ChatCompletionStream) + OpenAI-compatible client
    manifest/             - build own manifest from modelserver; serve /manifest; fetch + cache peers'
    peer/                 - call another node's OpenAI endpoint
  config.example.yaml
```

## Config shape

Two files, one per node (config-a.yaml, config-b.yaml), same binary.

```yaml
listen: ":8080"
node_id: "A"
model_server: "http://localhost:11434/v1"   # Ollama's OpenAI-compatible base URL
consumer_token: "dev-secret"                 # stub auth for the consumer door
peers:
  - "http://192.168.1.50:8080"               # node B's node endpoint
```

## Build order (milestones from the plan)

Build one milestone at a time and stop for review between each.

1. M1, walking skeleton. Two node processes. Node A receives an OpenAI request, forwards to node B over the OpenAI HTTP call, streams a canned or mocked response back. No real inference. Proves the loop, the streaming, and the two doors.
2. M2, real local inference and the manifest. Each node serves a real model through the modelserver interface; each publishes and fetches manifests; local-first works.
3. M3, overflow routing and the consumer door. Local-then-peer routing, the consumer door with the stub token, and honest degradation when nobody can serve.
4. M4, measurement. Move node B to a second machine over a residential link; run the curl harness across local, overflow, and consumer paths; record the kill-criteria numbers.

## Non-goals (do NOT build in the POC)

- Real identity. The consumer door is a pre-shared token only.
- Fair-use quotas, rate limits, priority, or preemption.
- Federation, gossip, or discovery beyond the static peer list in config.
- Clever load balancing. Selection is: has the model and has capacity, pick one.
- Sovereignty or policy gates.
- Any frontend. Use curl and Open WebUI (both are just OpenAI clients).

## Definition of done (code)

The loop works over two machines: local-first serves on the owner's node, overflow serves on the peer, a consumer with only the token borrows from the pod, and an unservable request returns an honest no-capacity response rather than hanging. Then the measurements from the plan are collected.

## Paste-ready kickoff prompt for Claude Code

```
Read architecture/Sketchbook.md and prototype/BUILD_BRIEF.md in this repo.

We're building the DII two-node prototype in Go under prototype/. Do milestone
M1 ONLY: two node processes where node A receives an OpenAI-compatible request,
forwards it to node B over the reused OpenAI HTTP call, and streams a mocked
response back to the caller. Prove the loop, streaming, and the two doors (local
vs. consumer, selected by the stub token). No real inference yet.

Follow the package layout and config shape in BUILD_BRIEF.md. KISS: control plane
only, OpenAI-compatible on both sides, backend behind a modelserver.Backend
interface, no identity beyond the stub token, no federation/quotas/load-balancing/
frontend.

First show me the package layout and the main.go wiring plus the config parsing,
and stop for my review. Then implement M1. I handle all git commits.
```
