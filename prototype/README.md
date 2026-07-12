# Prototype

Goal: a proof-of-concept DII pod. A DII node is a Go control-plane broker in front
of an OpenAI-compatible model server (Ollama); it never runs inference itself, it
routes each request to whichever node in the pod can serve the model. See
`../architecture/Sketchbook.md` and `BUILD_BRIEF.md` for the design, and
`DEPLOY.md` to run a pod across multiple machines.

## Status: M1–M3 complete

- **M1 — walking skeleton.** Two nodes, SSE streaming end to end, the two doors,
  A→B forwarding over the reused OpenAI HTTP call.
- **M2 — real inference + manifest.** Each node serves a real model through the
  `modelserver.Backend` seam (Ollama behind a URL); nodes build, publish, fetch,
  and cache manifests; the local door is manifest-gated.
- **M3 — overflow + honest degradation.** Capability-aware routing off the cached
  manifests: local-first, then overflow to a peer that has the model, else an
  honest 503. The two doors differ — local door is local-first; consumer door
  (stub token) prefers peers.

Verified live on a 3-node pod spanning two networks (LAN + Tailscale). Next up is
M4: measurement against the plan's kill-criteria.

## Layout

```
cmd/node/main.go        load config, build+exchange manifests, wire packages, start
internal/config         parse the YAML config
internal/ingress        OpenAI endpoint; token -> door; SSE relay; /manifest, /v1/models
internal/router         local-first + capability-aware overflow to peers
internal/modelserver    Backend interface + Ollama client (+ mock, from M1)
internal/manifest       node self-description; build own, fetch/cache peers'
internal/peer           call a peer node's OpenAI endpoint + /manifest (implements Backend)
```

## Run it (one machine, two nodes)

Requires Go 1.22+ and a running Ollama with at least one model pulled
(`ollama pull llama3.2:1b`). First time only:

```sh
cd prototype
go mod tidy      # resolves gopkg.in/yaml.v3
```

Start both nodes in separate terminals (both point at the local Ollama):

```sh
go run ./cmd/node -config config-a.yaml   # node A on :8080, peer -> :8081
go run ./cmd/node -config config-b.yaml   # node B on :8081, no peers
```

For the interesting overflow behavior (a node borrowing a model it doesn't have),
you need nodes with *different* models — see `DEPLOY.md` for a real multi-machine
pod.

## Drive it

```sh
# Local door (no token): served by node A's own Ollama
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"llama3.2:1b","stream":true,"messages":[{"role":"user","content":"hi"}]}'

# Consumer door (stub token): routed into the pod
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer dev-secret' \
  -d '{"model":"llama3.2:1b","stream":true,"messages":[{"role":"user","content":"hi"}]}'

# A model no node has -> honest 503
curl -s -w '\n%{http_code}\n' http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"nope:1b","messages":[{"role":"user","content":"hi"}]}'

# Inspect a node
curl -s http://localhost:8080/manifest    # what this node can serve
curl -s http://localhost:8080/v1/models   # OpenAI-shaped model list
curl    http://localhost:8080/healthz     # -> node A ok
```

## Helper scripts

```sh
scripts/ask.sh <model> [prompt]                 # ask the hub; router picks the node
scripts/pod-status.sh <node-url> [<node-url>…]  # each node's manifest + reachability
```
`DII_HUB` overrides the target node (default `http://localhost:8080`); `DII_TOKEN`
sends the consumer-door bearer token.
