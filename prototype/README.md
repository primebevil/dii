# Prototype

Goal: two-node proof of concept. A DII node is a Go control-plane broker in front
of an OpenAI-compatible model server. See `../architecture/Sketchbook.md` and
`BUILD_BRIEF.md` for the design.

## Status: M1 — walking skeleton

Two node processes, no real inference. Node A takes an OpenAI-compatible request
and either serves it from a mock backend (local door) or forwards it to peer node
B over the reused OpenAI HTTP call (consumer door), streaming a mocked response
back either way. Proves the loop, the streaming, and the two doors.

The mock backend stands in for the model server behind `modelserver.Backend`;
real inference (Ollama) and the manifest arrive in M2.

## Layout

```
cmd/node/main.go        load config, wire packages, start the node
internal/config         parse the YAML config
internal/ingress        OpenAI endpoint; token -> door; SSE relay with flush
internal/router         M1: local door -> local backend, consumer door -> peer
internal/modelserver    Backend interface + mock backend (M1)
internal/peer           call a peer node's OpenAI endpoint (implements Backend)
```

## Run it (one machine, two ports)

Requires Go 1.22+. First time only, fetch the one dependency:

```sh
cd prototype
go mod tidy      # resolves gopkg.in/yaml.v3
```

Start both nodes in separate terminals:

```sh
go run ./cmd/node -config config-a.yaml   # node A on :8080, peer -> :8081
go run ./cmd/node -config config-b.yaml   # node B on :8081, no peers
```

## Drive the two doors

Local door — no token, node A serves from its own mock backend:

```sh
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"mock-model","messages":[{"role":"user","content":"hi"}]}'
```

You should see SSE chunks arrive one at a time, ending in `data: [DONE]`, with the
content tagged **served by node A**.

Consumer door — stub token, node A forwards to node B (watch it stream through):

```sh
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer dev-secret' \
  -d '{"model":"mock-model","messages":[{"role":"user","content":"hi"}]}'
```

Same stream, but the content is tagged **served by node B** — proof the request
went A -> B and the tokens came back. A bad token returns `401`; the consumer door
with no reachable peer returns `503` (honest no-capacity).

Health check:

```sh
curl http://localhost:8080/healthz    # -> node A ok
```
