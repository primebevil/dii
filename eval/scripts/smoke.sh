#!/usr/bin/env bash
# One-command smoke test: stand up the mock pod, run the rig against it, tear down.
# Proves the pipeline end to end with no real models. Numbers are meaningless.
set -euo pipefail
cd "$(dirname "$0")/.."

PORT="${PORT:-8091}"

cat > config.smoke.yaml <<YAML
node_url: "http://127.0.0.1:${PORT}"
token: ""
pool: ["alpha:14b", "beta:12b", "gamma:24b", "delta:8b", "coder:14b"]
judge_model: "judge:32b"
samples: 0
temperature: 0.7
eval_fraction: 0.5
split_seed: 1
noise_repeats: 4
max_tokens: 1024
YAML

python3 tools/mock_server.py --tasks tasks --port "$PORT" >/tmp/dii_mock.log 2>&1 &
MOCK=$!
trap 'kill $MOCK 2>/dev/null || true' EXIT
sleep 1

python3 run.py --config config.smoke.yaml --tasks tasks --out /tmp/dii_smoke_results.json
echo
echo "smoke run complete. full results at /tmp/dii_smoke_results.json"
