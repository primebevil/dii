#!/usr/bin/env python3
"""A mock OpenAI-compatible endpoint that simulates a pod of uneven models.

For smoke-testing the rig only. It reads the task bank so it can return a
correct-ish answer for "strong" models and a wrong one for "weak" ones (strength
is a stable hash of the model name), which gives the scoreboard a real spread and
exercises voting, judging, selection, oracle ceilings, and stats end to end. The
numbers it produces are meaningless as evidence; it exists to prove the pipeline
runs before you point the rig at real models.

    python3 tools/mock_server.py --tasks tasks --port 8080

Coding tasks may carry a `mock_solution:` field, read only here, so the mock can
return code that passes the task's real unit tests.
"""

from __future__ import annotations

import argparse
import glob
import hashlib
import json
import os
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

import yaml

TASKS: dict[str, dict] = {}


def _norm(s: str) -> str:
    return " ".join((s or "").strip().lower().split())


def load_tasks(tasks_dir: str):
    for path in glob.glob(os.path.join(tasks_dir, "**", "*.yaml"), recursive=True):
        with open(path) as f:
            raw = yaml.safe_load(f)
        for item in (raw if isinstance(raw, list) else [raw]):
            if item and item.get("prompt"):
                TASKS[_norm(item["prompt"])] = item


def _strong(model: str) -> bool:
    return int(hashlib.sha256(model.encode()).hexdigest(), 16) % 3 != 0  # ~2/3 of models are "strong"


def _is_judge(model: str) -> bool:
    return "judge" in model.lower() or model.lower().endswith(":32b")


def _wrong_number(ans: str) -> str:
    try:
        return str(float(ans) + 7)
    except ValueError:
        return "0"


def _answer(model: str, prompt: str, seed: int) -> str:
    task = TASKS.get(_norm(prompt))
    if task is None:
        return "I don't know."
    strong = _strong(model)
    # a little per-seed wobble so the stochastic baseline and voting have work to do
    if (seed % 5) == 0:
        strong = not strong
    grader = task.get("grader")
    if grader == "numeric":
        return f"The answer is {task['answer'] if strong else _wrong_number(str(task['answer']))}."
    if grader in ("exact", "mcq"):
        return str(task["answer"]) if strong else "none"
    if grader == "unit_test":
        if strong and task.get("mock_solution"):
            return "```python\n" + task["mock_solution"] + "\n```"
        return "```python\ndef " + (task.get("entrypoint") or "solve") + "(*a, **k):\n    return None\n```"
    # judged/open-ended: strong models write more, and the mock judge prefers length
    return ("A thorough, well-reasoned answer. " * 8) if strong else "Short answer."


def _judge_reply(prompt: str) -> str:
    # The mock judge prefers the longer of the two answers, position-independently,
    # so the swap-agreement check yields real winners instead of all ties.
    a1 = prompt.split("ANSWER 1:", 1)[-1].split("ANSWER 2:", 1)[0] if "ANSWER 1:" in prompt else ""
    a2 = prompt.split("ANSWER 2:", 1)[-1] if "ANSWER 2:" in prompt else ""
    return "1" if len(a1) >= len(a2) else "2"


class Handler(BaseHTTPRequestHandler):
    def log_message(self, *a):
        pass

    def do_POST(self):
        n = int(self.headers.get("Content-Length", 0))
        req = json.loads(self.rfile.read(n) or b"{}")
        model = req.get("model", "")
        seed = int(req.get("seed", 0) or 0)
        messages = req.get("messages", [])
        prompt = next((m["content"] for m in reversed(messages) if m.get("role") == "user"), "")

        content = _judge_reply(prompt) if _is_judge(model) else _answer(model, prompt, seed)
        pt = max(1, len(prompt) // 4)
        ct = max(1, len(content) // 4)
        body = json.dumps({
            "id": "mock", "object": "chat.completion", "model": model,
            "choices": [{"index": 0, "message": {"role": "assistant", "content": content},
                         "finish_reason": "stop"}],
            "usage": {"prompt_tokens": pt, "completion_tokens": ct, "total_tokens": pt + ct},
        }).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--tasks", default="tasks")
    ap.add_argument("--port", type=int, default=8080)
    args = ap.parse_args()
    load_tasks(args.tasks)
    print(f"mock pod up on :{args.port} with {len(TASKS)} tasks indexed")
    ThreadingHTTPServer(("127.0.0.1", args.port), Handler).serve_forever()


if __name__ == "__main__":
    main()
