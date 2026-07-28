"""OpenAI-compatible client for a DII node (or any OpenAI-shaped endpoint).

Talks to the pod hub exactly like the week-3 harness and scripts/ask.sh do: a
POST to /v1/chat/completions with the model name, and the router decides which
node actually serves it. Non-streaming, because here we want the whole answer
plus token counts, not time-to-first-token.

Stdlib only, so the rig runs on a pod machine with nothing but Python 3.
"""

from __future__ import annotations

import json
import time
import urllib.error
import urllib.request
from dataclasses import dataclass

# Transient network hiccups (a cold mesh path, a momentary refusal) should not
# abort a run that is hundreds of calls long. Retry a few times with backoff on
# connection-level errors only; a real HTTP error (4xx/5xx from the node) is not
# retried, because that is the node telling us something true.
_RETRIES = 3
_BACKOFF_S = 1.5


@dataclass
class Completion:
    text: str
    prompt_tokens: int
    completion_tokens: int
    latency_s: float
    model: str

    @property
    def total_tokens(self) -> int:
        return self.prompt_tokens + self.completion_tokens


class NodeClient:
    def __init__(self, base_url: str, token: str = "", timeout_s: float = 300.0):
        self.base_url = base_url.rstrip("/")
        self.token = token
        self.timeout_s = timeout_s

    def chat(
        self,
        model: str,
        prompt: str,
        temperature: float = 0.0,
        seed: int | None = None,
        max_tokens: int | None = None,
        system: str | None = None,
    ) -> Completion:
        messages = []
        if system:
            messages.append({"role": "system", "content": system})
        messages.append({"role": "user", "content": prompt})

        payload: dict = {
            "model": model,
            "stream": False,
            "temperature": temperature,
            "messages": messages,
        }
        if seed is not None:
            payload["seed"] = seed
        if max_tokens is not None:
            payload["max_tokens"] = max_tokens

        body = json.dumps(payload).encode("utf-8")
        req = urllib.request.Request(
            self.base_url + "/v1/chat/completions",
            data=body,
            headers={"Content-Type": "application/json"},
            method="POST",
        )
        if self.token:
            req.add_header("Authorization", "Bearer " + self.token)

        start = time.monotonic()
        raw = None
        for attempt in range(_RETRIES):
            try:
                with urllib.request.urlopen(req, timeout=self.timeout_s) as resp:
                    raw = resp.read()
                break
            except urllib.error.HTTPError as e:
                detail = e.read().decode("utf-8", "replace")[:500]
                # 502/504 are gateway timeouts (a peer/model server too slow on a
                # cold load) -- transient, so retry. Other codes are the node
                # telling us something true (bad request, no capacity): do not retry.
                if e.code in (502, 504) and attempt < _RETRIES - 1:
                    time.sleep(_BACKOFF_S * (attempt + 1))
                    continue
                raise RuntimeError(f"node returned {e.code} for model {model}: {detail}") from e
            except urllib.error.URLError as e:
                if attempt == _RETRIES - 1:
                    raise RuntimeError(f"node unreachable at {self.base_url} after "
                                       f"{_RETRIES} tries: {e.reason}") from e
                time.sleep(_BACKOFF_S * (attempt + 1))
        latency = time.monotonic() - start

        data = json.loads(raw)
        choice = (data.get("choices") or [{}])[0]
        text = (choice.get("message") or {}).get("content", "") or ""
        usage = data.get("usage") or {}
        return Completion(
            text=text,
            prompt_tokens=int(usage.get("prompt_tokens", 0)),
            completion_tokens=int(usage.get("completion_tokens", 0)),
            latency_s=latency,
            model=model,
        )
