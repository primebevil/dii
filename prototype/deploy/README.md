# Running a node as a service

Two ways to keep a node running. Both stop with `SIGTERM`, which is what the node's
graceful drain expects, and both point at the config through `DII_CONFIG` so nothing
depends on a working directory.

| | Use it when | Files | Status |
|---|---|---|---|
| **Docker** | the default, on any platform | `Dockerfile`, `docker-compose.yaml`, `config-docker.example.yaml` | in use on Linux and macOS |
| **systemd** | Linux without a container runtime | `dii-node.service` | lifecycle verified; root install path not |

Take that status column literally. Docker is what runs today on both a Linux server
and a Mac. The systemd path was validated on atlas by running this unit's directives
as a user unit against the real binary and config: the service started and served
(`m1-check.sh` 12/12), restarted cleanly, came back by itself after the process was
`SIGKILL`ed, and drained on `systemctl stop` rather than being killed mid-stream.
`ProtectSystem=strict` and the other sandboxing directives did not interfere.

What that run could **not** cover, because it was unprivileged: `User=`/`Group=`,
`ProtectHome=true` (the test config sat in `$HOME`, where this unit uses
`/etc/dii`), the `After=ollama.service` ordering, and the `/opt/dii` + `/etc/dii`
install steps. Those are the parts to watch the first time you install it for real.

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
  interface (`"100.100.100.100:8080:8080"`) to keep the node off the others.

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

`docker-compose.gui.yaml` is an overlay you opt into with a second `-f`:

```sh
docker compose -f docker-compose.yaml -f docker-compose.gui.yaml up -d
# on macOS, add the bridged override too:
docker compose -f docker-compose.yaml -f docker-compose.macos.yaml -f docker-compose.gui.yaml up -d
```

It is an overlay rather than a service in the base file because a GUI is not part
of a node, it is a *consumer* of one. The node is an OpenAI-compatible endpoint, so
Open WebUI needs no node changes at all — a base URL and an API key, the two fields
every OpenAI client already has. The embeddings endpoint is what powers its document
and retrieval features.

Two things the overlay gets right that are easy to get wrong by hand:

- **It reaches the node through the host's published port** (`DII_NODE_URL`, default
  `http://host.docker.internal:8080`), not a compose service name. Service-name DNS
  does not resolve under host networking, and the node is not on the GUI's bridge
  network.
- **It sends an empty API key**, so the GUI uses the node's *local* door, which is
  local-first. Sending the `consumer_token` instead would open the consumer door —
  peer-first — and route your own work off this machine on purpose. Any other value
  is a 401; there is no third option.

One caveat: the GUI's own login is a layer above the node, separate from
per-consumer node identity. Behind one connection, every GUI user looks like a
single consumer to the node. That is fine for personal use; mapping GUI users to
distinct node consumers is M2 work.

## systemd

Header comments in `dii-node.service` carry the install commands. Two settings
matter: `KillSignal=SIGTERM` and `TimeoutStopSec=30s`, for the same drain reason as
the Docker grace period. `After=ollama.service` is ordering only — the node starts
fine without its backend and reports not-ready on `/readyz` until it answers.

The lifecycle half of this unit was exercised on atlas (see the status column
above): start, serve, restart, recover from a `SIGKILL`, and drain on stop all
behave. The install half — the `dii` user, `/opt/dii`, `/etc/dii`, and
`ProtectHome=true` — has not been, since that needs root. Run `m1-check.sh` against
the node after installing and fix this file if anything is off.

If you would rather not install system-wide, the same unit works as a user unit in
`~/.config/systemd/user/` with `User=`/`Group=` dropped and paths under `%h`; pair
it with `loginctl enable-linger` so it survives logout and starts at boot.

### The per-milestone check

The pod runs on containers day to day, so this path is the one that will rot. At
each milestone boundary, run it once and check it. On a Linux node:

```sh
# 1. build the binary the unit runs (from prototype/, on any machine)
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o node-linux-amd64 ./cmd/node
scp node-linux-amd64 <host>:/tmp/node

# 2. swap the container out for the service
ssh <host> 'cd ~/dii-node/deploy && docker compose stop'
ssh <host> 'sudo install -o dii -g dii /tmp/node /opt/dii/node && sudo systemctl restart dii-node'

# 3. check it, from anywhere that can reach the node
DII_CHAT_MODEL=<a model it serves> prototype/scripts/m1-check.sh http://<host>:8080

# 4. exercise the things a container otherwise hides
ssh <host> 'sudo systemctl restart dii-node'                    # restarts clean
ssh <host> 'sudo kill -9 $(systemctl show dii-node -p MainPID --value)'  # Restart=on-failure
ssh <host> 'sudo systemctl stop dii-node'                       # drains, not killed
ssh <host> 'sudo journalctl -u dii-node -n 5'                   # "draining..." then "stopped"

# 5. put the container back
ssh <host> 'sudo systemctl stop dii-node && cd ~/dii-node/deploy && docker compose start'
```

What this catches that the container hides: a new dependency on a path only a
volume mount provides, an environment variable only compose sets, a signal only
Docker delivers, or a libc dependency creeping into what is supposed to be a static
binary.

## Verifying either of them

```sh
curl -s localhost:8080/healthz    # liveness: the process is up
curl -s localhost:8080/readyz     # readiness: the model server is reachable too
curl -s localhost:8080/manifest   # what this node can serve
```

`/readyz` is the one that tells you something you could not otherwise see: a node
that is up but whose backend is down answers `200` on `/healthz` and `503` on
`/readyz`. That distinction is the whole point of having both.

Every request also writes one JSON accounting line to stdout — `docker compose logs`
or `journalctl -u dii-node`. See the usage-accounting section in `../README.md` for
the fields and what the token counts mean.
