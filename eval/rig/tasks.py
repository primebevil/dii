"""Task bank loading and the frozen dev/eval split.

A task is one exam question with a way to grade it. Objective tasks
(numeric/exact/mcq/unit_test) grade themselves against ground truth; judged tasks
carry a rubric and are scored by the independent judge. The dev/eval split is
deterministic (hash of the task id), so the baseline and router are always chosen
on the same held-out-independent half across runs.
"""

from __future__ import annotations

import glob
import hashlib
import os
from dataclasses import dataclass, field

import yaml

OBJECTIVE_GRADERS = {"numeric", "exact", "mcq", "unit_test"}
JUDGED_GRADERS = {"judge"}


@dataclass
class Task:
    id: str
    domain: str
    prompt: str
    grader: str
    answer: str = ""            # numeric/exact/mcq
    entrypoint: str = ""        # unit_test: function the candidate must define
    test: str = ""              # unit_test: python asserting on the candidate
    rubric: str = ""            # judge: what a good answer must do
    system: str = ""            # optional system prompt
    meta: dict = field(default_factory=dict)

    @property
    def is_objective(self) -> bool:
        return self.grader in OBJECTIVE_GRADERS


def load_tasks(tasks_dir: str) -> list[Task]:
    tasks: list[Task] = []
    seen: set[str] = set()
    for path in sorted(glob.glob(os.path.join(tasks_dir, "**", "*.yaml"), recursive=True)):
        with open(path) as f:
            raw = yaml.safe_load(f)
        if not raw:
            continue
        items = raw if isinstance(raw, list) else [raw]
        for item in items:
            known = {k: v for k, v in item.items() if k in Task.__dataclass_fields__}
            t = Task(**known)
            if t.grader not in OBJECTIVE_GRADERS | JUDGED_GRADERS:
                raise ValueError(f"task {t.id}: unknown grader '{t.grader}'")
            if t.id in seen:
                raise ValueError(f"duplicate task id: {t.id}")
            seen.add(t.id)
            tasks.append(t)
    return tasks


def split(tasks: list[Task], eval_fraction: float, seed: int) -> tuple[list[Task], list[Task]]:
    """Deterministic per-task assignment to dev or eval by hashing (id, seed)."""
    dev, ev = [], []
    for t in tasks:
        h = hashlib.sha256(f"{seed}:{t.id}".encode()).hexdigest()
        frac = int(h[:8], 16) / 0xFFFFFFFF
        (ev if frac < eval_fraction else dev).append(t)
    return dev, ev
