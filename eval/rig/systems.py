"""The five systems, and the run that produces every number.

Two-phase generation, on purpose, so cost is honest and the baseline is chosen
without peeking:

  Phase A: one greedy answer from every model on every task.
           The frozen baseline and the router table are chosen here, on the dev
           split, using objective ground truth only (the most defensible basis).
  Phase B: N stochastic answers from the frozen baseline on the eval tasks, which
           feed the same-model best-of-N baseline (S2) and the noise floor.

Systems assembled on the eval split:

  S1 naive       : baseline, one greedy sample.
  S2 same-N      : baseline, N stochastic samples, aggregated by the selector.
  S3 ensemble    : one greedy sample from each pool model, same selector.
  S4 route       : practical per-domain router (domain label only, frozen on dev).
  S5a oracle-route  : per task, the best single model (uses the grade; a ceiling).
  S5b oracle-select : the best of S3's candidates (uses the grade; a ceiling).

The contrasts that matter: S3 vs S2 (diversity beyond compute) and S4 vs S1
(routing beyond one generalist), with the oracles bounding what any method could
have captured.
"""

from __future__ import annotations

from statistics import mean, pstdev

from .grading import Judge, grade_objective, score_judged
from .selectors import select
from .tasks import Task

SYSTEMS = ["S1_naive", "S2_same_N", "S3_ensemble", "S4_route", "S5a_oracle_route", "S5b_oracle_select"]


def _greedy(client, model, task: Task, cfg):
    return client.chat(
        model, task.prompt, temperature=0.0, seed=0,
        max_tokens=cfg.max_tokens, system=task.system or None,
    )


def _dev_model_score(dev, greedy, model) -> float:
    obj = [t for t in dev if t.is_objective]
    if not obj:
        return 0.0
    return mean(grade_objective(t, greedy[t.id][model].text) for t in obj)


def run_experiment(cfg, dev, ev, client, judge: Judge | None, log=lambda s: None):
    pool = cfg.pool
    all_tasks = dev + ev

    # Phase A: greedy from every model on every task.
    greedy: dict[str, dict[str, object]] = {}
    for i, t in enumerate(all_tasks):
        greedy[t.id] = {}
        for m in pool:
            greedy[t.id][m] = _greedy(client, m, t, cfg)
        log(f"phase A: {i + 1}/{len(all_tasks)} tasks generated")

    # Freeze the baseline and the router table on the dev split (objective only).
    baseline = max(pool, key=lambda m: _dev_model_score(dev, greedy, m))
    domains = sorted({t.domain for t in all_tasks})
    router_table: dict[str, str] = {}
    for dom in domains:
        dev_dom_obj = [t for t in dev if t.domain == dom and t.is_objective]
        if dev_dom_obj:
            router_table[dom] = max(
                pool, key=lambda m: mean(grade_objective(t, greedy[t.id][m].text) for t in dev_dom_obj)
            )
        else:
            router_table[dom] = baseline
    log(f"baseline frozen: {baseline}; router table: {router_table}")

    # Phase B: N stochastic samples from the frozen baseline on eval tasks.
    N = cfg.samples
    stoch: dict[str, list] = {}
    for i, t in enumerate(ev):
        stoch[t.id] = [
            client.chat(baseline, t.prompt, temperature=cfg.temperature, seed=r + 1,
                        max_tokens=cfg.max_tokens, system=t.system or None)
            for r in range(N)
        ]
        log(f"phase B: {i + 1}/{len(ev)} baseline sample sets")

    per_task = []
    for t in ev:
        g = greedy[t.id]
        baseline_text = g[baseline].text

        s1 = g[baseline].text
        s2_cands = [c.text for c in stoch[t.id][:N]]
        s2 = select(t, s2_cands, judge)
        s3_cands = [g[m].text for m in pool]
        s3 = select(t, s3_cands, judge)
        s4_model = router_table[t.domain]
        s4 = g[s4_model].text

        def grade(text: str) -> float:
            if t.is_objective:
                return grade_objective(t, text)
            return score_judged(t, text, baseline_text, judge) if judge else 0.5

        scores = {
            "S1_naive": grade(s1),
            "S2_same_N": grade(s2),
            "S3_ensemble": grade(s3),
            "S4_route": grade(s4),
            "S5a_oracle_route": max(grade(g[m].text) for m in pool),
            "S5b_oracle_select": max(grade(c) for c in s3_cands),
        }
        tokens = {
            "S1_naive": g[baseline].completion_tokens,
            "S2_same_N": sum(c.completion_tokens for c in stoch[t.id][:N]),
            "S3_ensemble": sum(g[m].completion_tokens for m in pool),
            "S4_route": g[s4_model].completion_tokens,
            "S5a_oracle_route": sum(g[m].completion_tokens for m in pool),
            "S5b_oracle_select": sum(g[m].completion_tokens for m in pool),
        }
        calls = {
            "S1_naive": 1, "S2_same_N": N, "S3_ensemble": len(pool),
            "S4_route": 1, "S5a_oracle_route": len(pool), "S5b_oracle_select": len(pool),
        }
        # Per-call latencies, so the report can show both serial cost and the
        # parallel critical path (the pod's real latency when calls fan out).
        latencies = {
            "S1_naive": [g[baseline].latency_s],
            "S2_same_N": [c.latency_s for c in stoch[t.id][:N]],
            "S3_ensemble": [g[m].latency_s for m in pool],
            "S4_route": [g[s4_model].latency_s],
            "S5a_oracle_route": [g[m].latency_s for m in pool],
            "S5b_oracle_select": [g[m].latency_s for m in pool],
        }
        per_task.append({
            "id": t.id, "domain": t.domain, "objective": t.is_objective,
            "scores": scores, "tokens": tokens, "calls": calls, "latencies": latencies,
        })

    # Noise floor: reuse the baseline's N stochastic samples as `noise_repeats`
    # independent single-sample runs; the spread of their eval-set mean scores is
    # the run-to-run wobble any real effect must clear. Computed per domain as well
    # as globally, because domains differ in intrinsic variance and the plan sets a
    # meaningful margin PER DOMAIN, so a single global floor would mis-gate them.
    #
    # CAVEAT (documented, not silently hidden): this is single-sample dispersion,
    # while S2/S3 are N-sample aggregates with lower variance. So this floor is a
    # CONSERVATIVE bar for the aggregate contrasts -- it errs toward false nulls,
    # not false wins. That direction is the safe one for a skeptic's benchmark, but
    # it means a near-threshold null should be read as "not proven", not "disproven".
    reps = min(cfg.noise_repeats, N)

    def _score_rep(t, r):
        text = stoch[t.id][r].text
        return (grade_objective(t, text) if t.is_objective
                else (score_judged(t, text, greedy[t.id][baseline].text, judge) if judge else 0.5))

    def _noise_over(tasks) -> float:
        if not tasks:
            return 0.0
        run_means = [mean(_score_rep(t, r) for t in tasks) for r in range(reps)]
        return pstdev(run_means) if len(run_means) > 1 else 0.0

    noise_floor = _noise_over(ev)
    domains = sorted({t.domain for t in ev})
    noise_floor_by_domain = {d: _noise_over([t for t in ev if t.domain == d]) for d in domains}

    return {
        "baseline_model": baseline,
        "router_table": router_table,
        "systems": SYSTEMS,
        "per_task": per_task,
        "noise_floor": noise_floor,
        "noise_floor_by_domain": noise_floor_by_domain,
        "margin": cfg.margin,
        "margin_by_domain": dict(cfg.margin_by_domain),
        "judge_calls": judge.calls if judge else 0,
        "N": N,
        "pool": pool,
    }
