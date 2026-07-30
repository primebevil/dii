#!/usr/bin/env python3
"""Variety experiment runner. Points at a DII pod (or any OpenAI-shaped endpoint),
runs the five systems over the task bank, and prints the scoreboard.

    python3 run.py --config config.yaml
    python3 run.py --config config.yaml --tasks tasks --out results.json

See docs/Variety_Experiment_Plan.md for the protocol this executes.
"""

from __future__ import annotations

import argparse
import sys

from rig.client import NodeClient
from rig.config import Config
from rig.grading import Judge
from rig.report import render, save_json
from rig.systems import run_experiment
from rig.tasks import load_tasks, split


def main() -> int:
    ap = argparse.ArgumentParser(description="DII variety aggregation experiment")
    ap.add_argument("--config", required=True, help="path to config.yaml")
    ap.add_argument("--tasks", default="tasks", help="task bank directory")
    ap.add_argument("--out", default="", help="optional path to write full results JSON")
    ap.add_argument("--quiet", action="store_true", help="suppress progress lines")
    args = ap.parse_args()

    cfg = Config.load(args.config)
    tasks = load_tasks(args.tasks)
    if not tasks:
        print(f"no tasks found under {args.tasks}", file=sys.stderr)
        return 2

    dev, ev = split(tasks, cfg.eval_fraction, cfg.split_seed)
    if not ev or not dev:
        print("split produced an empty dev or eval set; add tasks or adjust eval_fraction",
              file=sys.stderr)
        return 2

    client = NodeClient(cfg.node_url, cfg.token)
    # The judge may live on a separate endpoint (e.g. a hosted frontier model);
    # falls back to the pod endpoint when judge_url is unset.
    judge_client = NodeClient(cfg.judge_url or cfg.node_url, cfg.judge_token or cfg.token)
    judge = Judge(judge_client, cfg.judge_model, max_tokens=cfg.judge_max_tokens) if cfg.judge_model else None

    def log(s):
        if not args.quiet:
            print(f"  ... {s}", file=sys.stderr)

    print(f"tasks: {len(tasks)}  (dev {len(dev)} / eval {len(ev)})  "
          f"pool {len(cfg.pool)}  N {cfg.samples}  judge {cfg.judge_model or 'none'}",
          file=sys.stderr)
    if judge is None:
        print("  warning: no judge_model set; judged tasks will all score 0.5 (tie). "
              "Set judge_model for open-ended domains.", file=sys.stderr)
        if any(t.grader == "unit_test" for t in ev):
            print("  warning: unit_test (coding) tasks present but no judge; the "
                  "ensemble/best-of-N selector cannot rank code without one and will "
                  "fall back to the FIRST candidate. S2/S3 coding results will be "
                  "degenerate. A judge is required for coding, not just open-ended tasks.",
                  file=sys.stderr)

    results = run_experiment(cfg, dev, ev, client, judge, log=log)
    print()
    print(render(results))
    if args.out:
        save_json(results, args.out)
        print(f"\nfull results written to {args.out}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
