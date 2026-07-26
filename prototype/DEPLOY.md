# Deploying a multi-machine pod

One binary is a node; a pod is several nodes that know each other. This guide
covers running nodes on separate machines. For a single-machine two-node run, see
`README.md`.

## Build

Native:

```sh
go build -o node ./cmd/node
```

Cross-compile for another machine (e.g. a Linux x86-64 server, built from a Mac):

```sh
GOOS=linux GOARCH=amd64 go build -o node-linux-amd64 ./cmd/node
```

The binary is statically linked (no libc/runtime deps) — copy it over and run.
Common targets: `linux/amd64`, `linux/arm64`, `darwin/arm64`, `windows/amd64`.

## Topology: reachability is the only rule

Each node fetches its peers' manifests at startup and forwards requests to them —
both are plain OpenAI HTTP calls. So the only requirement is:

> Every address in a node's `peers:` list must be reachable **from that node**, and
> each node's `listen`/`advertise` address must be reachable by the peers that call it.

Nodes can sit on different networks. A node reachable on two networks (say a LAN
and a Tailscale mesh) can act as a **hub** that bridges peers which can't reach
each other directly. Every request enters the hub, which routes it to whichever
node holds the model.

### Worked example: our 3-node test pod

| Node | Machine | Reachable at | Serves |
|------|---------|--------------|--------|
| **laptop** (hub) | macOS | LAN `10.0.0.248` **and** Tailscale `100.94.5.71` | `llama3.2:1b` |
| **atlas** | Linux server | Tailscale `100.118.77.40` | big qwen/deepseek models |
| **sirius** | Linux laptop | LAN `10.0.0.216` | qwen2.5-coder / phi models |

`atlas` and `sirius` are on different networks and can't reach each other, but both
reach the laptop — so the laptop is the hub. Configs (one per node; swap the IPs
for your own):

**config-laptop.yaml** (hub — binds all interfaces to answer on both networks)
```yaml
listen: ":8080"
node_id: "laptop"
model_server: "http://localhost:11434/v1"
consumer_token: "dev-secret"
advertise: "http://100.94.5.71:8080"
peers:
  - "http://100.118.77.40:8090"       # atlas over Tailscale
  - "http://10.0.0.216:8080"          # sirius over the LAN
```

**config-atlas.yaml** (spoke — binds to its Tailscale IP only)
```yaml
listen: "100.118.77.40:8090"
node_id: "atlas"
model_server: "http://localhost:11434/v1"
consumer_token: "dev-secret"
advertise: "http://100.118.77.40:8090"
peers:
  - "http://100.94.5.71:8080"         # laptop over Tailscale
```

**config-sirius.yaml** (spoke — binds to its LAN IP only)
```yaml
listen: "10.0.0.216:8080"
node_id: "sirius"
model_server: "http://localhost:11434/v1"
consumer_token: "dev-secret"
advertise: "http://10.0.0.216:8080"
peers:
  - "http://10.0.0.248:8080"          # laptop over the LAN
```

Bind `listen` to a specific interface IP (as the spokes do) to keep a node off
public interfaces; the hub binds `:8080` (all interfaces) because it must answer on
two networks.

## Run a node (Linux, detached)

Copy the binary + its config over, then start it so it survives your SSH session:

```sh
mkdir -p ~/dii
# scp the binary and config-<name>.yaml into ~/dii
chmod +x ~/dii/node
setsid nohup ~/dii/node -config ~/dii/config-<name>.yaml > ~/dii/node.log 2>&1 < /dev/null &
tail ~/dii/node.log
```

**Start order:** bring up the spoke nodes first, then the hub — a node fetches peer
manifests only at **startup**, so peers should be listening when it boots. To pick
up a newly pulled model or a peer that came up late, restart the node.

## Verify the pod

From the hub:

```sh
scripts/pod-status.sh http://localhost:8080 http://100.118.77.40:8090 http://10.0.0.216:8080
scripts/ask.sh llama3.2:1b       # served locally on the hub
scripts/ask.sh qwen2.5:7b        # overflows to the node that has it
scripts/ask.sh nope:1b           # -> honest 503, no node can serve it
```

## Stop a node

```sh
ssh <host> 'pkill -f "dii/node"'
```

## Known limits (POC)

- Manifests are built once at **startup**; restart to refresh models or peers.
- A node advertises a single `advertise` endpoint even when it's reachable at
  several addresses. Fine for hub-and-spoke (spokes don't call each other), but a
  multi-homed node can't yet advertise different addresses to different peers.
- `consumer_token` is a shared stub, not real identity.
- No discovery: peers are the static list in config.
