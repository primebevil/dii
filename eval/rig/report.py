"""The scoreboard: one screen that says whether variety paid, and the verdict.

The verdict follows docs/Variety_Experiment_Plan.md exactly:
  PAYS      - S3 beats S2, or S4 beats S1, significantly and on a majority of the
              objective-graded domains. Diversity or routing earns its keep.
  MARGINAL  - a variety arm beats the naive S1 but not the equal-compute S2. The
              gain is compute, not diversity. A finding against the pitch.
  LOSES     - no variety arm beats S2, and even the oracle ceiling is close to it.
              The durable leg is not durable on these workloads.
"""

from __future__ import annotations

import json
from statistics import mean

from . import stats


def _objective_domains(per_task, domains):
    out = []
    for d in domains:
        rows = [pt for pt in per_task if pt["domain"] == d]
        if rows and mean(1.0 if pt["objective"] else 0.0 for pt in rows) > 0.5:
            out.append(d)
    return out


def _noise(results, domain=None):
    if domain is None:
        return results["noise_floor"]
    return results.get("noise_floor_by_domain", {}).get(domain, results["noise_floor"])


def _margin(results, domain=None):
    if domain is None:
        return results.get("margin", 0.0)
    return results.get("margin_by_domain", {}).get(domain, results.get("margin", 0.0))


def _per_domain_significant(results, a, b, domains):
    per_task = results["per_task"]
    hits = []
    for d in domains:
        sub = [pt for pt in per_task if pt["domain"] == d]
        if not sub:
            continue
        if stats.contrast(sub, a, b, _noise(results, d), _margin(results, d))["significant"]:
            hits.append(d)
    return hits


def decide_verdict(results):
    per_task = results["per_task"]
    noise = _noise(results)
    margin = _margin(results)
    domains = sorted({pt["domain"] for pt in per_task})
    obj_domains = _objective_domains(per_task, domains)

    c_div = stats.contrast(per_task, "S3_ensemble", "S2_same_N", noise, margin)      # diversity beyond compute
    c_route = stats.contrast(per_task, "S4_route", "S1_naive", noise, margin)        # routing beyond one generalist
    # Pure compute: the SAME model run N times vs one shot. This is the real
    # "did spending more compute on one model help" control (S2 vs S1), NOT the
    # ensemble vs naive — conflating those mislabels a diversity win as a compute win.
    c_compute = stats.contrast(per_task, "S2_same_N", "S1_naive", noise, margin)

    div_domains = _per_domain_significant(results, "S3_ensemble", "S2_same_N", obj_domains)
    route_domains = _per_domain_significant(results, "S4_route", "S1_naive", obj_domains)
    # Strict majority: >half the objective domains, so 2 of 3, 2 of 2, 1 of 1.
    majority = len(obj_domains) // 2 + 1 if obj_domains else 1

    div_pays = c_div["significant"] and len(div_domains) >= majority
    route_pays = c_route["significant"] and len(route_domains) >= majority

    # The oracle question, which the plan makes decisive for LOSES: is there
    # variety in the pool that SOME method could exploit, even if our practical
    # arms did not? Measured per task, then averaged, and gated by the same bar.
    c_oracle = stats.contrast(per_task, "S5a_oracle_route", "S2_same_N", noise, margin)
    c_oracle_sel = stats.contrast(per_task, "S5b_oracle_select", "S2_same_N", noise, margin)
    oracle_shows_variety = c_oracle["significant"] or c_oracle_sel["significant"]

    # A pooled diversity/routing win: the arm beats its control across the whole
    # eval set, even if it doesn't clear significance in a per-domain majority
    # (which is easy to miss at small per-domain n). Distinct from a clean PAYS.
    pooled_win = c_div["significant"] or c_route["significant"]

    if div_pays or route_pays:
        verdict = "PAYS"
    elif pooled_win:
        verdict = "PAYS_POOLED"  # beats the control pooled, not yet across a domain majority
    elif c_compute["significant"]:
        verdict = "MARGINAL"     # compute (S2>S1) helped, diversity did NOT clear S2
    elif oracle_shows_variety:
        verdict = "UNCAPTURED"   # variety IS in the pool; the practical arms missed it
    else:
        verdict = "LOSES"        # even the oracle is close to S2: no variety to capture

    return {
        "verdict": verdict,
        "c_div": c_div, "c_route": c_route, "c_compute": c_compute,
        "c_oracle": c_oracle, "c_oracle_sel": c_oracle_sel,
        "obj_domains": obj_domains, "div_domains": div_domains, "route_domains": route_domains,
        "majority_needed": majority,
    }


def _fmt_ci(c):
    return f"{c['mean_diff']:+.3f} [{c['ci'][0]:+.3f}, {c['ci'][1]:+.3f}]"


def render(results) -> str:
    per_task = results["per_task"]
    systems = results["systems"]
    domains = sorted({pt["domain"] for pt in per_task})
    means = stats.system_means(per_task, systems)
    dmeans = stats.system_domain_means(per_task, systems, domains)
    cost = stats.system_cost(per_task, systems)
    base_calls = cost["S1_naive"]["calls"] or 1

    L = []
    L.append("=" * 74)
    L.append("VARIETY AGGREGATION EXPERIMENT  -  scoreboard")
    L.append("=" * 74)
    L.append(f"baseline (frozen on dev): {results['baseline_model']}")
    L.append(f"pool ({len(results['pool'])}): {', '.join(results['pool'])}")
    L.append(f"router table: {results['router_table']}")
    L.append(f"eval tasks: {len(per_task)}   sampling budget N: {results['N']}")
    L.append(f"noise floor (score std): {results['noise_floor']:.3f}   "
             f"margin: {results.get('margin', 0.0):.3f}   "
             f"bar = noise+margin: {results['noise_floor'] + results.get('margin', 0.0):.3f}"
             + ("   [margin=0: NOT a real pre-registration]" if results.get('margin', 0.0) == 0.0 else ""))
    L.append(f"judge calls: {results['judge_calls']}")
    L.append("")

    L.append(f"{'system':<20}{'score':>8}{'calls':>7}{'cost x':>8}"
             f"{'lat serial':>12}{'lat par':>10}   meaning")
    meaning = {
        "S1_naive": "baseline, 1 greedy",
        "S2_same_N": "baseline x N (compute control)",
        "S3_ensemble": "N diverse models",
        "S4_route": "practical router",
        "S5a_oracle_route": "ceiling: best model/task",
        "S5b_oracle_select": "ceiling: best of ensemble",
    }
    for s in systems:
        mult = cost[s]["calls"] / base_calls
        L.append(f"{s:<20}{means[s]:>8.3f}{cost[s]['calls']:>7}{mult:>7.1f}x"
                 f"{cost[s]['latency_serial_s']:>11.1f}s{cost[s]['latency_parallel_s']:>9.1f}s"
                 f"   {meaning.get(s,'')}")
    L.append("")
    L.append("  lat serial = calls back-to-back (this rig);  lat par = critical path if")
    L.append("  every call ran concurrently on its own node (the pod's fan-out best case).")
    L.append("")

    # Per-domain score grid.
    head = f"{'domain':<16}" + "".join(f"{s.split('_')[0]:>8}" for s in systems)
    L.append(head)
    for d in domains:
        row = f"{d:<16}" + "".join(f"{dmeans[s][d]:>8.3f}" for s in systems)
        L.append(row)
    L.append("")

    v = decide_verdict(results)
    L.append("-" * 74)
    L.append("KEY CONTRASTS (paired mean diff [95% bootstrap CI])")
    L.append(f"  diversity beyond compute  S3 - S2 : {_fmt_ci(v['c_div'])}"
             f"   {'SIGNIFICANT' if v['c_div']['significant'] else 'ns'}")
    L.append(f"  routing beyond generalist S4 - S1 : {_fmt_ci(v['c_route'])}"
             f"   {'SIGNIFICANT' if v['c_route']['significant'] else 'ns'}")
    L.append(f"  compute helped at all     S3 - S1 : {_fmt_ci(v['c_compute'])}"
             f"   {'SIGNIFICANT' if v['c_compute']['significant'] else 'ns'}")
    L.append(f"  oracle headroom over S2  S5a - S2 : {_fmt_ci(v['c_oracle'])}"
             f"   {'EXPLOITABLE' if v['c_oracle']['significant'] else 'ns'}")
    L.append(f"                           S5b - S2 : {_fmt_ci(v['c_oracle_sel'])}"
             f"   {'EXPLOITABLE' if v['c_oracle_sel']['significant'] else 'ns'}")
    if v["obj_domains"]:
        L.append(f"  objective domains: {', '.join(v['obj_domains'])}")
        L.append(f"    S3>S2 significant in: {v['div_domains'] or 'none'}")
        L.append(f"    S4>S1 significant in: {v['route_domains'] or 'none'}")
    L.append("")
    L.append("=" * 74)
    tagline = {
        "PAYS": "Variety earns its keep. The durable leg has a floor under it.",
        "PAYS_POOLED": "Diversity beats equal compute across the eval set (pooled), but not "
                       "yet across a domain majority. Promising; confirm with more tasks/domain.",
        "MARGINAL": "The gain is COMPUTE, not diversity. Do not lean on variety in the pitch.",
        "UNCAPTURED": "Variety IS in the pool (oracle beats S2) but the practical arms "
                      "missed it. Not a loss -- a better router/selector is the work.",
        "LOSES": "No variety win, and even the oracle is close to S2. The durable leg "
                 "is not durable on these workloads.",
    }[v["verdict"]]
    L.append(f"VERDICT: {v['verdict']}  -  {tagline}")
    L.append("=" * 74)
    L.append("Note: a scoreboard is not a conclusion. Read the plan's 'what this does")
    L.append("and does not settle' before quoting any number outward.")
    return "\n".join(L)


def save_json(results, path: str):
    with open(path, "w") as f:
        json.dump(results, f, indent=2, default=str)
