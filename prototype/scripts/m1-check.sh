#!/usr/bin/env bash
# Check everything M1 added, against any running node: the floor bundle including
# embeddings, both doors, honest failures, and liveness vs. readiness.
#
# Usage:   scripts/m1-check.sh [<node-url>]
#          defaults to http://localhost:8080
# Env:     DII_TOKEN        consumer-door token          (default dev-secret)
#          DII_CHAT_MODEL   a chat model the pod serves  (default llama3.2:1b)
#          DII_EMBED_MODEL  an embedding model it serves (default all-minilm:l6-v2)
#
# Exits 0 if every check passed, 1 otherwise. Does not test graceful shutdown —
# that one needs to stop the node, so it is a manual step (see prototype/README.md).
set -uo pipefail

url="${1:-http://localhost:8080}"
token="${DII_TOKEN:-dev-secret}"
chat="${DII_CHAT_MODEL:-llama3.2:1b}"
embed="${DII_EMBED_MODEL:-all-minilm:l6-v2}"

pass=0; fail=0

# check <name> <got> <want>
check() {
  if [ "$2" = "$3" ]; then
    printf '  \033[32mPASS\033[0m  %s\n' "$1"; pass=$((pass+1))
  else
    printf '  \033[31mFAIL\033[0m  %s\n          got  %s\n          want %s\n' "$1" "$2" "$3"; fail=$((fail+1))
  fi
}

# post <path> <json> [token] -> "<status>|<content-type>"
post() {
  local args=(-s -o /dev/null -w '%{http_code}|%{content_type}' -X POST "$url$1" -H 'Content-Type: application/json')
  [ -n "${3:-}" ] && args+=(-H "Authorization: Bearer $3")
  curl -m 300 "${args[@]}" -d "$2" 2>/dev/null
}

status() { curl -s -o /dev/null -m 15 -w '%{http_code}' "$url$1" 2>/dev/null; }

echo "Checking $url"
echo "  chat model: $chat        embedding model: $embed"
echo

echo "the endpoints exist and answer"
check "GET /healthz  (liveness)"  "$(status /healthz)"  "200"
check "GET /manifest"             "$(status /manifest)" "200"
check "GET /v1/models"            "$(status /v1/models)" "200"
check "the embedding model is in /v1/models" \
  "$(curl -s -m 15 "$url/v1/models" | grep -c -- "$embed")" "1"

echo
echo "the floor bundle: chat and embeddings, on the owner door (no token)"
check "chat, streaming     -> SSE" \
  "$(post /v1/chat/completions "{\"model\":\"$chat\",\"stream\":true,\"max_tokens\":8,\"messages\":[{\"role\":\"user\",\"content\":\"hi\"}]}")" \
  "200|text/event-stream"
check "chat, non-streaming -> JSON" \
  "$(post /v1/chat/completions "{\"model\":\"$chat\",\"max_tokens\":8,\"messages\":[{\"role\":\"user\",\"content\":\"hi\"}]}")" \
  "200|application/json"
check "embeddings          -> JSON" \
  "$(post /v1/embeddings "{\"model\":\"$embed\",\"input\":\"the reliable floor\"}")" \
  "200|application/json"

echo
echo "the consumer door (token) — prefers a peer, falls back to local"
check "chat via consumer door" \
  "$(post /v1/chat/completions "{\"model\":\"$chat\",\"stream\":true,\"max_tokens\":8,\"messages\":[{\"role\":\"user\",\"content\":\"hi\"}]}" "$token")" \
  "200|text/event-stream"
check "embeddings via consumer door" \
  "$(post /v1/embeddings "{\"model\":\"$embed\",\"input\":\"x\"}" "$token")" \
  "200|application/json"

echo
echo "honest failure, not hanging"
check "a model nobody serves -> 503" \
  "$(post /v1/chat/completions '{"model":"nope:1b","messages":[{"role":"user","content":"hi"}]}')" \
  "503|application/json"
check "a bad token           -> 401" \
  "$(post /v1/embeddings "{\"model\":\"$embed\",\"input\":\"x\"}" "definitely-not-the-token")" \
  "401|application/json"

echo
echo "readiness is separate from liveness"
ready="$(status /readyz)"
check "GET /readyz (backend up -> 200)" "$ready" "200"
echo "        to see the two diverge, stop the model server and re-run:"
echo "        /healthz stays 200, /readyz turns 503 with a reason."

echo
if [ "$fail" -eq 0 ]; then
  printf '\033[32m%s passed, 0 failed\033[0m\n' "$pass"
else
  printf '\033[31m%s passed, %s FAILED\033[0m\n' "$pass" "$fail"
fi
echo
echo "Now look at what the node logged: one JSON line per request, with"
echo "served_by = local or the peer that ran it. A peer-served request appears"
echo "on both nodes."

[ "$fail" -eq 0 ]
