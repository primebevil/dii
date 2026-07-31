# Running a node as a service

Two ways to keep a node running. Both stop with `SIGTERM`, which is what the node's
graceful drain expects, and both point at the config through `DII_CONFIG` so nothing
depends on a working directory.

| | Use it when | Files | Status |
|---|---|---|---|
| **Docker** | the default, on any platform | `Dockerfile`, `docker-compose.yaml`, `config-docker.example.yaml` | in use on Linux and macOS |
| **systemd** | Linux without a container runtime | `dii-node.service` | **written but never run** |

Take that status column literally. Docker is the path with evidence behind it: it is
what runs today on both a Linux server and a Mac. The systemd unit is offered
because a Linux box without a container runtime is a perfectly reasonable way to run
a node, but nobody has yet started it in anger, so treat it as a starting point
rather than a tested recipe.

**On macOS**, use Docker. It is verified there, and it changes nothing about
Ollama, which has to run natively either way because a container cannot reach the
GPU. If you specifically want the node itself native on a Mac, `launchd` is the
mechanism — a user agent in `~/Library/LaunchAgents` invoking the binary with
`DII_CONFIG` set. `launchctl unload` sends `SIGTERM`, so the drain behaves exactly
as it does everywhere else. We do not ship a plist for this, because it would be one
more untested file to keep honest.

## The one hard rule

**Inference stays out of the node image.** The node is control plane only — it
never runs a model, it brokers to a model server over a URL. So Ollama runs
natively on the host where the GPU is, and the node reaches it at
`http://host.docker.internal:11434/v1`. On Linux you could later containerize
Ollama with GPU passthrough if you want uniformity; on macOS the GPU is not
reachable from a container, so Ollama stays native there regardless.

## Docker

```sh
cd prototype/deploy
cp config-docker.example.yaml config-docker.yaml   # then edit node_id, advertise, peers
docker compose up -d
docker compose logs -f
```

Three things in the compose file are load-bearing, not incidental:

- **`stop_grace_period: 30s`.** `docker stop` sends `SIGTERM`, waits, then
  `SIGKILL`s. The node drains in-flight streams for up to `shutdown_timeout`
  (default 15s), so the grace period has to be comfortably larger or a drain gets
  killed part-way through. Docker's own default is 10s, which is too tight.
- **`extra_hosts: host.docker.internal:host-gateway`.** Inside a container,
  `localhost` is the container. Without this line the node cannot see the host's
  Ollama on Linux.
- **`ports` must match `advertise`.** Peers and consumers dial the address in
  `advertise`; that has to be a host address Docker actually publishes. Bind to one
  interface (`"100.118.77.40:8090:8090"`) to keep the node off the others.

The health check runs `node -healthcheck`, which probes the node's own `/readyz`
and exits 0 or 1. The image is distroless — no shell, no curl — so the binary does
its own probing rather than the image carrying a tool just for this.

That has one sharp edge worth knowing, because it fails confusingly. **A
shell-form health check cannot work in this image**, since there is no shell to run
it. In particular `docker run --health-cmd "..."` always wraps the command in
`/bin/sh -c`, so it reports the container `unhealthy` even while the probe run by
hand succeeds. The healthcheck therefore lives in the Dockerfile in JSON exec form,
where it is correct however the container is launched. If you override it in
compose, keep the exec form: `test: ["CMD", "/usr/local/bin/node", "-healthcheck"]`.

### The optional GUI

`docker-compose.yaml` has a commented-out `open-webui` service. It is worth
understanding *why* it is only a comment away and needs no node changes: the node
is an OpenAI-compatible endpoint, so a web GUI is not an integration, it is just
another consumer. Point its `OPENAI_API_BASE_URL` at the node's `/v1`, give it a
consumer token as the API key, and it lists models and chats through the node.
The embeddings endpoint is what powers its document and retrieval features.

One caveat: the GUI's own login is a layer above the node, separate from
per-consumer node identity. Behind one shared token, every GUI user looks like a
single consumer to the node. That is fine for personal use; mapping GUI users to
distinct node consumers is M2 work.

## systemd (untested)

Header comments in `dii-node.service` carry the install commands. Two settings
matter: `KillSignal=SIGTERM` and `TimeoutStopSec=30s`, for the same drain reason as
the Docker grace period. `After=ollama.service` is ordering only — the node starts
fine without its backend and reports not-ready on `/readyz` until it answers.

To be plain about it: this unit has never been started. It is written from the same
understanding of the node's lifecycle that the container path proves, so the drain
and readiness behaviour should carry over — but the install steps, the `dii` user,
and the sandboxing directives have not been exercised. Run `m1-check.sh` against it
before trusting it, and fix this file if anything is off.

## Verifying either of them

```sh
curl -s localhost:8090/healthz    # liveness: the process is up
curl -s localhost:8090/readyz     # readiness: the model server is reachable too
curl -s localhost:8090/manifest   # what this node can serve
```

`/readyz` is the one that tells you something you could not otherwise see: a node
that is up but whose backend is down answers `200` on `/healthz` and `503` on
`/readyz`. That distinction is the whole point of having both.

Every request also writes one JSON accounting line to stdout — `docker compose logs`
or `journalctl -u dii-node`. See the usage-accounting section in `../README.md` for
the fields and what the token counts mean.
