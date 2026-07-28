"""Grading: objective graders against ground truth, and the blind judge.

Objective tasks grade themselves (a number matches, tests pass) and need no
model opinion. Judged tasks are scored by an independent judge as a blind,
position-swapped pairwise comparison against the frozen baseline's answer, so a
judged system's per-task score is simply "did it beat the baseline here": 1 win,
0.5 tie, 0 loss. Ground truth is preferred everywhere it exists; the judge is the
fallback, not the default.

SECURITY: the unit_test grader executes model-produced code in a subprocess.
Run the rig on a throwaway machine or container, never on anything you care about.
"""

from __future__ import annotations

import os
import re
import subprocess
import sys
import tempfile

from .tasks import Task

_NUM = re.compile(r"-?\d+(?:,\d{3})*(?:\.\d+)?")
_CODE_FENCE = re.compile(r"```(?:python)?\s*\n(.*?)```", re.DOTALL | re.IGNORECASE)


def _last_number(text: str) -> str | None:
    matches = _NUM.findall(text.replace(",", ""))
    return matches[-1] if matches else None


def _extract_code(text: str) -> str:
    fences = _CODE_FENCE.findall(text)
    if fences:
        return "\n\n".join(f.strip() for f in fences)
    return text  # assume the whole answer is code if unfenced


def grade_objective(task: Task, text: str, unit_test_timeout_s: float = 20.0) -> float:
    if task.grader == "numeric":
        got = _last_number(text)
        if got is None:
            return 0.0
        try:
            return 1.0 if abs(float(got) - float(task.answer)) <= 1e-6 else 0.0
        except ValueError:
            return 0.0
    if task.grader == "exact":
        return 1.0 if _exact_ok(text, task.answer) else 0.0
    if task.grader == "mcq":
        letter = _first_choice(text)
        return 1.0 if letter and letter == task.answer.strip().upper()[:1] else 0.0
    if task.grader == "unit_test":
        return _run_unit_test(task, text, unit_test_timeout_s)
    raise ValueError(f"grade_objective called on non-objective task {task.id}")


def _norm(s: str) -> str:
    return " ".join(s.strip().lower().split())


def _extract_exact(text: str) -> str:
    """Best-guess the model's intended exact answer from a reasoning-then-answer
    response. Tasks instruct 'give the final answer on its own line', so prefer the
    last non-empty line's final token (stripped of '=', trailing punctuation)."""
    lines = [ln for ln in (text or "").splitlines() if ln.strip()]
    if not lines:
        return _norm(text)
    toks = _norm(lines[-1]).replace("=", " ").split()
    return toks[-1].strip(".,;:!") if toks else _norm(lines[-1])


def _exact_ok(text: str, answer: str) -> bool:
    """True if `answer` matches the whole response, appears as its own line, or is
    the final token of the last line. Robust to models that show work first."""
    target = _norm(answer)
    if not target:
        return False
    if _norm(text) == target:
        return True
    if any(_norm(ln) == target for ln in (text or "").splitlines()):
        return True
    return _extract_exact(text) == target


def _first_choice(text: str) -> str | None:
    m = re.search(r"\b([A-E])\b", text.strip().upper())
    return m.group(1) if m else None


def _run_unit_test(task: Task, text: str, timeout_s: float) -> float:
    code = _extract_code(text)
    program = (
        code
        + "\n\n# ---- test harness (task-supplied) ----\n"
        + task.test
        + "\nprint('__RIG_PASS__')\n"
    )
    with tempfile.TemporaryDirectory() as d:
        path = os.path.join(d, "candidate.py")
        with open(path, "w") as f:
            f.write(program)
        try:
            proc = subprocess.run(
                [sys.executable, path],
                capture_output=True,
                text=True,
                timeout=timeout_s,
                cwd=d,
            )
        except subprocess.TimeoutExpired:
            return 0.0
    return 1.0 if "__RIG_PASS__" in proc.stdout else 0.0


class Judge:
    """Independent judge for open-ended tasks. Pairwise, blind, position-swapped."""

    def __init__(self, client, model: str, max_tokens: int = 512):
        self.client = client
        self.model = model
        self.max_tokens = max_tokens
        self.calls = 0  # judge calls are measurement overhead, tracked but not charged to any system

    def _ask_winner(self, task: Task, first: str, second: str) -> str:
        self.calls += 1
        rubric = task.rubric or "overall correctness, usefulness, and adherence to the request"
        prompt = (
            "You are judging two answers to the same task. Decide which better "
            f"satisfies this rubric: {rubric}\n\n"
            f"TASK:\n{task.prompt}\n\n"
            f"ANSWER 1:\n{first}\n\n"
            f"ANSWER 2:\n{second}\n\n"
            "Reply with exactly one token: 1 if answer 1 is better, 2 if answer 2 "
            "is better, or TIE if they are equivalent. No other text."
        )
        out = self.client.chat(self.model, prompt, temperature=0.0, max_tokens=self.max_tokens).text
        u = out.strip().upper()
        if u.startswith("1"):
            return "first"
        if u.startswith("2"):
            return "second"
        return "tie"

    def compare(self, task: Task, a: str, b: str) -> str:
        """Return 'a', 'b', or 'tie', requiring agreement across a position swap."""
        r1 = self._ask_winner(task, a, b)   # a is first
        r2 = self._ask_winner(task, b, a)   # a is second
        first_says = {"first": "a", "second": "b", "tie": "tie"}[r1]
        second_says = {"first": "b", "second": "a", "tie": "tie"}[r2]
        if first_says == second_says:
            return first_says
        return "tie"  # verdict flipped on swap -> not trustworthy, call it a tie

    def pick_best(self, task: Task, candidates: list[str]) -> int:
        """Index of the best candidate by pairwise knockout (stable, swap-checked)."""
        best = 0
        for i in range(1, len(candidates)):
            winner = self.compare(task, candidates[best], candidates[i])
            if winner == "b":
                best = i
        return best


def score_judged(task: Task, system_text: str, baseline_text: str, judge: Judge) -> float:
    winner = judge.compare(task, system_text, baseline_text)
    return {"a": 1.0, "tie": 0.5, "b": 0.0}[winner]
