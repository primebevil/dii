#!/usr/bin/env python3
"""Verify the task bank before trusting a run.

A wrong answer key or a vacuous unit test corrupts the benchmark silently. This
checks every task mechanically:

  numeric   : the reference worked-solution (mock_solution) actually grades 1.0
              against the answer key, and a deliberately-wrong number grades 0.0.
  exact/mcq : the answer key is self-consistent (grades itself 1.0) and a wrong
              value grades 0.0.
  unit_test : the reference solution (mock_solution) PASSES the test, and a do-
              nothing stub FAILS it (so the test is not vacuous).
  judge     : schema only (rubric present) -- correctness can't be auto-checked.

`mock_solution` is a verification-only field the rig itself ignores (tasks.py
drops unknown fields on load), so we read the raw YAML to reach it.

    python3 tools/verify_tasks.py --tasks tasks
    python3 tools/verify_tasks.py --tasks tasks/coding      # one domain

Exit code is nonzero if any task fails, so it can gate a run in a script.

SECURITY: unit_test verification executes the reference solutions in a subprocess,
same as grading. Run on a throwaway machine.
"""

from __future__ import annotations

import argparse
import glob
import os
import re
import sys

import yaml

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from rig.grading import grade_objective  # noqa: E402
from rig.tasks import JUDGED_GRADERS, OBJECTIVE_GRADERS, Task  # noqa: E402


def _wrong_number(answer: str) -> str:
    try:
        return str(float(answer) + 1)
    except ValueError:
        return "definitely-not-the-answer-zzz"


def _stub(entrypoint: str) -> str:
    return f"```python\ndef {entrypoint or 'solve'}(*a, **k):\n    return None\n```"


def _num_appears(text: str, answer: str) -> bool:
    """True if the numeric answer appears somewhere in the worked solution — a soft
    correctness signal (a right solution should reach and state the answer)."""
    try:
        a = float(answer)
    except (ValueError, TypeError):
        return True  # non-numeric: not our check to make
    for n in re.findall(r"-?\d+(?:\.\d+)?", (text or "").replace(",", "")):
        try:
            if abs(float(n) - a) <= 1e-6:
                return True
        except ValueError:
            continue
    return False


def verify_task(raw: dict) -> list[str]:
    """Return a list of failure strings; empty means the task is sound."""
    fails = []
    known = {k: v for k, v in raw.items() if k in Task.__dataclass_fields__}
    try:
        t = Task(**known)
    except TypeError as e:
        return [f"cannot construct Task: {e}"]
    mock = raw.get("mock_solution", "")

    if t.grader not in OBJECTIVE_GRADERS | JUDGED_GRADERS:
        fails.append(f"unknown grader '{t.grader}'")
        return fails

    if t.grader == "numeric":
        # Grader self-consistency, on the reasoning-then-answer format models produce.
        if grade_objective(t, f"...work...\nThe answer is {t.answer}") != 1.0:
            fails.append(f"answer key '{t.answer}' does not grade 1.0 on a "
                         "reasoning-then-answer response")
        if grade_objective(t, f"...work...\nThe answer is {_wrong_number(t.answer)}") != 0.0:
            fails.append("a wrong number still grades 1.0 (grader not discriminating)")
        # Soft correctness: the worked solution should actually reach the answer value.
        if mock and not _num_appears(mock, t.answer):
            fails.append(f"answer '{t.answer}' never appears in the worked mock_solution "
                         "(likely key/solution mismatch — check by hand)")

    elif t.grader in ("exact", "mcq"):
        if grade_objective(t, t.answer) != 1.0:
            fails.append(f"answer key '{t.answer}' does not grade itself 1.0")
        if t.grader == "exact" and grade_objective(t, f"...reasoning...\n{t.answer}") != 1.0:
            fails.append(f"exact answer '{t.answer}' not matched when given on its own line")
        wrong = "z z z wrong" if t.grader == "exact" else "Z"
        if grade_objective(t, wrong) != 0.0:
            fails.append("a wrong value still grades 1.0 (grader not discriminating)")

    elif t.grader == "unit_test":
        if not t.entrypoint:
            fails.append("unit_test task missing entrypoint")
        if not t.test:
            fails.append("unit_test task missing test")
        if not mock:
            fails.append("no mock_solution reference to validate the test")
        else:
            if grade_objective(t, mock) != 1.0:
                fails.append("reference solution FAILS its own test (bad test or bad solution)")
            if grade_objective(t, _stub(t.entrypoint)) != 0.0:
                fails.append("a do-nothing stub PASSES the test (test is vacuous)")

    elif t.grader == "judge":
        if not (t.rubric or "").strip():
            fails.append("judge task has no rubric")

    return fails


def main() -> int:
    ap = argparse.ArgumentParser(description="Verify the task bank")
    ap.add_argument("--tasks", default="tasks", help="task dir (recursed) or a single .yaml")
    args = ap.parse_args()

    if os.path.isfile(args.tasks):
        paths = [args.tasks]
    else:
        paths = sorted(glob.glob(os.path.join(args.tasks, "**", "*.yaml"), recursive=True))

    seen_ids: dict[str, str] = {}
    total = ok = 0
    bad: list[tuple[str, str, list[str]]] = []
    by_domain: dict[str, int] = {}

    for path in paths:
        with open(path) as f:
            raw = yaml.safe_load(f)
        if not raw:
            continue
        for item in (raw if isinstance(raw, list) else [raw]):
            if not item:
                continue
            tid = item.get("id", "<no-id>")
            total += 1
            fails = []
            if tid in seen_ids:
                fails.append(f"duplicate id (also in {seen_ids[tid]})")
            else:
                seen_ids[tid] = os.path.basename(path)
            fails += verify_task(item)
            if fails:
                bad.append((tid, path, fails))
            else:
                ok += 1
                by_domain[item.get("domain", "?")] = by_domain.get(item.get("domain", "?"), 0) + 1

    print(f"verified {ok}/{total} tasks sound")
    if by_domain:
        print("sound tasks by domain: " + ", ".join(f"{d}={n}" for d, n in sorted(by_domain.items())))
    if bad:
        print(f"\n{len(bad)} PROBLEM task(s):")
        for tid, path, fails in bad:
            print(f"  [{tid}]  ({os.path.basename(path)})")
            for fmsg in fails:
                print(f"      - {fmsg}")
        return 1
    print("\nall tasks sound.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
