"""Answer selectors, applied identically to the same-model and diversity systems.

The selector is a function of the task type, never of which system called it, so
S2 (same model, N samples) and S3 (N different models) are aggregated the exact
same way and the only thing that differs between them is model diversity. That
identity is what makes S3-vs-S2 a clean test of diversity rather than of compute.

For tasks with a comparable answer key (numeric/exact/mcq) the selector is a
majority vote, the natural self-consistency aggregator. For code and open-ended
tasks, where you cannot vote and cannot peek at ground truth before answering,
the selector defers to the blind judge picking the best candidate.
"""

from __future__ import annotations

import re
from collections import Counter

from .grading import _extract_exact, _first_choice, _last_number
from .tasks import Task


def _vote_key(task: Task, text: str) -> str | None:
    if task.grader == "numeric":
        return _last_number(text)
    if task.grader == "exact":
        return _extract_exact(text)  # the answer token, so reasoning prose doesn't split the vote
    if task.grader == "mcq":
        return _first_choice(text)
    return None


def select(task: Task, candidates: list[str], judge) -> str:
    """Return the aggregated answer text for a set of candidate answers."""
    if not candidates:
        return ""
    if len(candidates) == 1:
        return candidates[0]

    if task.grader in ("numeric", "exact", "mcq"):
        keyed = [(cand, _vote_key(task, cand)) for cand in candidates]
        counts = Counter(k for _, k in keyed if k is not None)
        if counts:
            winning_key, _ = counts.most_common(1)[0]
            for cand, k in keyed:
                if k == winning_key:
                    return cand  # a representative candidate carrying the modal answer
        return candidates[0]

    # unit_test / judge: no vote key and no peeking at ground truth -> blind judge.
    if judge is None:
        return candidates[0]
    return candidates[judge.pick_best(task, candidates)]
