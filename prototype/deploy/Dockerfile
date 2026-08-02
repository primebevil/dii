# The DII node as a container image.
#
# One hard rule, straight from the control-plane-only constraint: inference stays
# out of this image. The node brokers to a model server over a URL; Ollama runs
# natively where the GPU is. See docker-compose.yaml.
#
# Build from the prototype/ directory, not from here:
#   docker build -f deploy/Dockerfile -t dii-node:latest .

# --- build: a fully static binary, so the runtime image can be minimal ---
FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/node ./cmd/node

# --- runtime: distroless static — no shell, no package manager, nonroot ---
# CA certificates come with this base, so a peer or backend reached over HTTPS
# works without adding anything.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/node /usr/local/bin/node

# The listen port comes from the mounted config, so there is deliberately no
# EXPOSE here to go stale — publish whatever port the config binds.
ENV DII_CONFIG=/etc/dii/config.yaml

# The node probes its own /readyz, so a container whose model server has died shows
# as unhealthy instead of merely running. It has to be the JSON exec form: there is
# no shell in this image, so a shell-form healthcheck cannot execute at all. That
# also rules out `docker run --health-cmd`, which always wraps in /bin/sh -c — which
# is why the healthcheck lives here in the image rather than at the call site.
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 \
  CMD ["/usr/local/bin/node", "-healthcheck"]

USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/node"]
