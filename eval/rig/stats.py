"""Paired statistics: mean scores, cost, and bootstrap CIs on the contrasts.

Systems are compared paired (same tasks), so the unit is the per-task score
difference. A contrast counts only if its confidence interval excludes zero and
its mean clears the measured noise floor. That double bar is deliberate: it keeps
a two-point win that is really run-to-run wobble from being read as a result.
"""

from __future__ import annotations

import random
from statistics import mean


def system_means(per_task: list[dict], systems: list[str]) -> dict[str, float]:
    return {s: mean(pt["scores"][s] for pt in per_task) for s in systems}


def system_domain_means(per_task, systems, domains) -> dict[str, dict[str, float]]:
    out = {}
    for s in systems:
        out[s] = {}
        for d in domains:
            rows = [pt["scores"][s] for pt in per_task if pt["domain"] == d]
            out[s][d] = mean(rows) if rows else float("nan")
    return out


def system_cost(per_task, systems) -> dict[str, dict]:
    out = {}
    for s in systems:
        toks = sum(pt["tokens"][s] for pt in per_task)
        calls = sum(pt["calls"][s] for pt in per_task)
        # Two honest latency numbers per system, summed over tasks:
        #   serial   = the calls run back-to-back (what this rig actually did).
        #   parallel = the critical path if every call ran concurrently on its own
        #              node (the pod's best case). For a diversity system this is
        #              the single slowest model, not the sum -- which is exactly the
        #              cost dimension where a pod earns its keep.
        lat = [pt.get("latencies", {}).get(s, []) for pt in per_task]
        serial = sum(sum(x) for x in lat)
        parallel = sum(max(x) if x else 0.0 for x in lat)
        out[s] = {"tokens": toks, "calls": calls,
                  "latency_serial_s": serial, "latency_parallel_s": parallel}
    return out


def paired_diffs(per_task, a: str, b: str) -> list[float]:
    return [pt["scores"][a] - pt["scores"][b] for pt in per_task]


def bootstrap_ci(values: list[float], iters: int = 5000, alpha: float = 0.05, seed: int = 7):
    if not values:
        return (0.0, 0.0, 0.0)
    rng = random.Random(seed)
    n = len(values)
    means = []
    for _ in range(iters):
        resample = [values[rng.randrange(n)] for _ in range(n)]
        means.append(sum(resample) / n)
    means.sort()
    lo = means[int((alpha / 2) * iters)]
    hi = means[int((1 - alpha / 2) * iters)]
    return (mean(values), lo, hi)


def win_loss_tie(per_task, a: str, b: str) -> tuple[int, int, int]:
    w = l = t = 0
    for pt in per_task:
        d = pt["scores"][a] - pt["scores"][b]
        if d > 1e-9:
            w += 1
        elif d < -1e-9:
            l += 1
        else:
            t += 1
    return w, l, t


def contrast(per_task, a: str, b: str, noise_floor: float, margin: float = 0.0) -> dict:
    diffs = paired_diffs(per_task, a, b)
    m, lo, hi = bootstrap_ci(diffs)
    w, l, t = win_loss_tie(per_task, a, b)
    ci_excludes_zero = lo > 0 or hi < 0
    # The bar is the noise floor PLUS the pre-registered meaningful margin: a real
    # effect must clear both run-to-run wobble and the effect size we declared
    # worth caring about before seeing results.
    bar = noise_floor + margin
    beats_noise = m > bar
    return {
        "a": a, "b": b, "mean_diff": m, "ci": (lo, hi), "bar": bar,
        "win": w, "loss": l, "tie": t,
        "ci_excludes_zero": ci_excludes_zero, "beats_noise": beats_noise,
        "significant": ci_excludes_zero and beats_noise and m > 0,
    }
