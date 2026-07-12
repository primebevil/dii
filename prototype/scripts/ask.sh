#!/usr/bin/env bash
# Ask the pod's hub for a model; the router picks which node actually serves it.
#
# Usage:   scripts/ask.sh <model> [prompt]
# Env:     DII_HUB    node door to hit (default http://localhost:8080)
#          DII_TOKEN  if set, sent as the consumer-door bearer token
#
# Prints the streamed reply text, or the error JSON verbatim on failure.
set -euo pipefail

model="${1:?usage: ask.sh <model> [prompt]}"
prompt="${2:-say hello in five words}"
hub="${DII_HUB:-http://localhost:8080}"

args=(-sN "$hub/v1/chat/completions" -H 'Content-Type: application/json')
[ -n "${DII_TOKEN:-}" ] && args+=(-H "Authorization: Bearer $DII_TOKEN")

resp="$(curl "${args[@]}" \
  -d "{\"model\":\"$model\",\"stream\":true,\"messages\":[{\"role\":\"user\",\"content\":\"$prompt\"}]}")"

# An error response has no content deltas; show it plainly instead of swallowing it.
if printf '%s' "$resp" | grep -q '"error"'; then
  printf '%s\n' "$resp"
  exit 1
fi

printf '%s' "$resp" | grep -o '"content":"[^"]*"' | sed 's/"content":"//; s/"$//' | tr -d '\n'
echo
