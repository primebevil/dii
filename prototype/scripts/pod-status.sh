#!/usr/bin/env bash
# Show each node's manifest (node_id + models) and whether it's reachable.
#
# Usage: scripts/pod-status.sh <node-url> [<node-url> ...]
#        defaults to the local hub if no URLs are given.
set -uo pipefail

urls=("$@")
[ "${#urls[@]}" -eq 0 ] && urls=("http://localhost:8080")

for u in "${urls[@]}"; do
  m="$(curl -s -m 5 "$u/manifest" 2>/dev/null || true)"
  if [ -z "$m" ]; then
    printf '%-34s UNREACHABLE\n' "$u"
  else
    printf '%-34s %s\n' "$u" "$m"
  fi
done
