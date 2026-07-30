# Variety Aggregation Experiment — Findings and Validation Assessment

Status: Findings, 2026-07-27. Companion to `docs/Variety_Experiment_Plan.md` (the
pre-registration) and the `eval/` rig that executed it. This document records what
was actually run, what came back, and — more importantly than either — an honest
assessment of what these results can and cannot support. It is written to be argued
with. The headline: the experiment as it stands **does not settle the variety
question**, because two of its four measurement instruments are not yet trustworthy;
what it does provide is one real signal, one robust secondary finding, and a clear
list of what has to be fixed before a verdict means anything.

## The claim under test

From the plan: on the work DII's target users do, does a pool of diverse
floor-class models beat the single best model in that pool? Split into two bets:

- **H-route** (specialization): routing each task to its best-suited model beats
  any single generalist.
- **H-ensemble** (diversity): aggregating several models' answers to the same task
  beats any single model run at the same sampling budget.

The decisive contrasts are **S3 vs S2** (diversity beyond equal compute) and **S4 vs
S1** (routing beyond one generalist), with oracle ceilings (S5) bounding what any
method could have captured.

## What was built

- A disposable measurement rig (`eval/`) that talks to the pod exactly like a
  consumer would: OpenAI-compatible calls to a node, no privileged access.
- A task bank of **102 tasks** — coding 33 (unit-tested), reasoning 33
  (numeric/exact), long-context 24 (9 objective + 3 judged per half), privacy 12
  (judged). Every task was mechanically verified (answer key grades itself, grader
  discriminates, unit tests execute against a reference); the reasoning and
  long-context objective keys were each independently re-solved by a second pass.
- Pool under test: `gemma3:12b` (Google), `llama3.1:8b` (Meta), `qwen3:30b` (Qwen,
  reasoning), `qwen2.5-coder:14b` (code specialist). Judge: `mistral-small:24b`,
  held out of the pool. All floor-class open weights on owned hardware (atlas), no
  hosted stand-ins.

## What was run

Two runs, both paired, both under the split-and-freeze rule (baseline and router
chosen on a held-out-independent dev split, using objective ground truth only):

### Run A — 54-task pilot (31 eval tasks)

Baseline frozen on dev: `qwen2.5-coder:14b`.

| system | overall | coding | longctx | privacy | reasoning |
|---|---|---|---|---|---|
| S1 naive | 0.726 | 0.857 | 0.786 | 0.500 | 0.778 |
| S2 same-N | 0.726 | 0.857 | 0.643 | 0.500 | 0.889 |
| S3 ensemble | 0.871 | 0.857 | 0.929 | 0.688 | 1.000 |
| S4 route | 0.823 | 0.857 | 0.929 | 0.500 | 1.000 |
| S5 oracle | 0.919–0.935 | 1.000 | 0.929 | 0.75–0.81 | 1.000 |

Key contrasts: **S3−S2 = +0.145 [+0.048, +0.274], significant**; S2−S1 = 0.000;
S4−S1 = +0.097 ns; oracle over S2 exploitable. Corrected verdict (see Deviations):
**PAYS_POOLED** — diversity beat equal compute across the eval set, though not
across a per-domain majority.

### Run B — 102-task full run (59 eval tasks)

Baseline frozen on dev: `gemma3:12b` (note: the baseline **changed** from Run A).

| system | overall | coding | longctx | privacy | reasoning |
|---|---|---|---|---|---|
| S1 naive | 0.805 | 0.833 | 0.750 | 0.500 | 0.947 |
| S2 same-N | 0.805 | 0.889 | 0.679 | 0.500 | 0.947 |
| S3 ensemble | 0.839 | 0.833 | 0.893 | 0.500 | 0.947 |
| S4 route | 0.805 | 0.722 | 0.893 | 0.500 | 0.947 |
| S5 oracle | 0.873 | 0.944 | 0.893 | 0.500 | 0.947 |

Key contrasts: **S3−S2 = +0.034 [−0.034, +0.102], not significant**; S2−S1 = 0.000;
S4−S1 = 0.000; oracle over S2 = +0.068 [0.000, 0.153] ns. Verdict: **LOSES** by the
decision rule.

## The central finding: the pilot signal did not replicate

Run A's significant diversity win (+0.145) fell to a non-significant +0.034 in Run
B. Two mechanisms explain the drop, and only one of them is about the thesis:

1. **The pilot baseline was weak.** On the smaller dev split, coding tasks
   dominated the baseline selection and `qwen2.5-coder:14b` — a specialist that is
   mediocre off-coding — was frozen as the bar. On the larger split, the strong
   generalist `gemma3:12b` was selected. A tougher baseline is harder to beat, and
   a meaningful share of the pilot's "variety win" was beating a soft baseline.
   This is the split-and-freeze rule working as designed; it is why the run was
   scaled up rather than quoted.
2. **Two domains measured nothing** (next section), diluting any real signal into a
   pooled null.

The pilot is therefore best read as a small-sample, weak-baseline fluke rather than
evidence for variety. This is the intended function of the pre-registered larger
run, and it did its job.

## Why the LOSES verdict cannot carry weight yet — validity threats

The plan frames LOSES as a result that "puts the project's reason to exist in real
question." That framing assumes the instruments are sound. Here, half of them are
not.

**Construct validity — the judge does not discriminate.** Privacy scored exactly
0.500 for every system including both oracles in both runs. `mistral-small:24b`,
under the position-swap agreement check, tied every long-prose answer against the
baseline. The privacy slice measured the judge's indecision, not answer quality.
The plan explicitly required judge validation against a hand-scored subset before
trusting it; that validation was not performed. **Any judged-domain result is
currently uninterpretable.**

**Internal validity — reasoning is saturated.** Reasoning scored 0.947 for all six
systems in Run B: every model solves ~18 of 19 tasks. With no variance between
systems, no effect of any kind can be detected. The "harder" task expansion was not
hard enough for these models. A ceiling this high makes the domain a null by
construction.

**The disagreement check was skipped.** The plan requires a manipulation check —
measure how often the pool's models actually disagree — precisely so that a null can
be attributed correctly. It was not run. The reasoning saturation is itself indirect
evidence of low inter-model disagreement on that domain: if all models agree and all
are right, aggregation has nothing to work with, and the null is a roster/task
finding, not a thesis finding.

**Statistical validity.** The bank is below the pre-registered 50–100 tasks/domain
(eval per domain: coding 18, long-context 14, reasoning 19, privacy 8). Two
non-informative domains drag the pooled contrast toward zero regardless of the
thesis. The noise floor is a single-sample dispersion applied to aggregate
contrasts — conservative, biased toward false nulls.

**External validity.** One dated roster, four models, mid-2026 open weights, one
hardware substrate. Nothing here speaks to durability over time or to larger pools.

Net: Run B is **inconclusive with broken instruments**, not a clean refutation of
variety.

## What the experiment does support

Two things survive the caveats:

1. **Diversity helps on long-context — the one domain where the instrument works.**
   In both runs, the ensemble beat the baseline on long-context by a wide margin
   (Run B: S3 0.893 vs S1 0.750), and notably beat same-model best-of-N (0.679) by
   +0.21. This is objective grading over hard tasks, so it is not a judge artifact.
   It is the clearest evidence in the study that a diverse pool can produce
   something a single model does not — plausibly because independent models make
   independent errors on long-document extraction. This is the lead worth pulling.

2. **Pure compute does not pay — and sometimes hurts.** S2 (best-of-N on the single
   best model) equaled S1 in both runs (compute alone bought nothing) and actively
   *underperformed* the baseline on long-context (0.679 vs 0.750) — the judge-based
   selector picked worse answers than one greedy shot. This robust secondary
   finding matters for messaging: "spend more on one model" is not a free win, which
   makes the diversity-vs-compute comparison the right frame.

## Deviations from the pre-registration

Recorded plainly, since a pre-registration is only worth what its deviations
disclose:

1. **Verdict-logic bug fix (correctness, but applied after seeing Run B).** The
   "compute helped" contrast was coded as S3−S1 (ensemble vs naive) rather than
   S2−S1 (the real compute control), which had mislabeled Run A's diversity win as
   "MARGINAL / the gain is compute." Corrected to S2−S1. The underlying contrasts
   are unaffected; only the verdict label logic changed.
2. **`PAYS_POOLED` verdict added post-hoc** to represent "beats the control pooled
   but not across a domain majority," a state no existing label fit. Added to the
   code, deliberately **not** written into the pre-registration pending explicit
   ratification. Moot for Run B (LOSES), but the taxonomy question stands.
3. **Judge not human-validated** (plan required it). Consequence: judged domains
   untrustworthy, as observed.
4. **Sample size below the pre-registered 50–100/domain.**
5. **Pool is four models, not the five-family roster** — `mistral-small` was held
   out as the judge because no fast, strong, out-of-family judge was otherwise
   available on the pod (a reasoning judge is too slow across hundreds of calls).
6. **Inter-model disagreement manipulation check not performed.**

## What has to happen before a real verdict

Not "abandon variety" — the experiment cannot currently test it fairly:

1. **Harder reasoning tasks** (or a harder objective domain). The current set is
   below these models' ceiling; there is no variance to measure.
2. **Fix the judge**: validate `mistral-small` against hand-scored answers, replace
   it (a hosted frontier judge via `judge_url` is supported), or drop judged
   domains to objective-only. Judged results are dead weight until then.
3. **Run the disagreement check** so any future null can be correctly attributed.
4. **Expand and exploit the long-context lead** — the one place variety demonstrably
   pays; understand the mechanism and test whether it generalizes.
5. **Reach the pre-registered per-domain sample size** on the fixed instruments.

## What this settles and what it does not

It does not settle whether variety pays; two of four instruments were not measuring.
It does establish, on this dated roster: (a) a real, objective-grading diversity
benefit on long-context; (b) that spending equal compute on one model does not pay
and can hurt; and (c) that the pilot's apparent variety win was fragile to baseline
strength and sample size. It leaves the durable-leg question open — pending a
benchmark that can actually see the effect it is trying to measure. The size and
sovereignty legs are untouched, as the plan notes; they do not rise or fall with
this.
