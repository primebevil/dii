# Variety Aggregation Experiment Plan

Status: Plan, 2026-07-26. Operationalizes the one open empirical question left by
docs/Pod_Aggregation_Red_Team.md: does aggregation for *variety* actually pay on
DII's target workloads, or is it an intuition. Depends on ADR-0006 (reliable
floor), ADR-0009 (public and private pods), and Who_DII_Is_For. This document is
the pre-registration: the baseline rule, the metrics, the sample sizes, and the
pass and kill thresholds are fixed here, before any run, exactly as the Week-3
kill criteria were signed off before the M4 measurement. No code and no run in
this session; this is the protocol for review.

## Why this experiment exists

The red team already conceded the size claim and does not ask us to defend it.
Self-hosting is the entry ramp, pods are structurally poor at large-model
serving, and "reach a bigger model than your card holds" erodes as local models
improve. What the whole outward pitch now rests on is the variety claim: that a
pod reaching many diverse models produces something no single strong model can.
The red team is blunt that this claim currently has no evidence behind it and
names an experiment as the thing to build before leaning on it in messaging.

So the experiment has one job. On the work DII's target users actually do, does a
pool of diverse floor-class models beat the single best model in that pool. If
yes, the durable leg has a floor under it and the messaging can lead with
variety. If no, the leg is not durable and that has to be said plainly now,
while it is cheap to say, rather than after the pitch is built on it.

## The claim, split so it can be measured

Variety buys two goods that get said in one breath, and they are different bets:

Variety as specialization. Different models are better at different tasks, so
routing each task to its best-suited model beats any single generalist. This is
closest to what the prototype's capability router already does structurally.

Variety as diversity. Running several different models on the same task and
aggregating their answers beats any single model, because independent models
make independent errors and disagree in useful ways.

Testing both head to head against one honest baseline is the point. The red team
leaves open which form pays, and they can come apart: routing can win while
ensembling does not, or the reverse.

## Hypotheses, stated so they can fail

- H-route: a practical router over the diverse pool scores higher than the single
  best model, on held-out tasks, by a margin that exceeds measurement noise, net
  of its added cost.
- H-ensemble: a diversity ensemble over the pool scores higher than the single
  best model *run at the same sampling budget*, by a margin that exceeds noise,
  net of cost.
- Null: neither arm beats the single best model once compute and noise are
  accounted for. The variety leg does not pay on these workloads.

The phrase "at the same sampling budget" in H-ensemble is load-bearing and is
defended under Controls below. It is the difference between measuring diversity
and measuring "we spent more compute."

## The baseline, chosen honestly

The bar is not an average model, it is the strongest single model in the pool.
A skeptic will hold us to that, so we fix the selection rule before the run:

- The pool's models are each scored on a development split of the task suite,
  separate from the held-out evaluation split.
- The single model with the highest overall development score becomes the
  baseline. It is chosen once, on the development split, and then frozen. It is
  never re-picked per task or after seeing the evaluation results.
- The largest model in the pool is always a baseline candidate, so a variety win
  can never be dismissed as "you just avoided using the big model."

This split-and-freeze rule is what stops the baseline from being quietly
weakened to make variety look good.

## The comparison matrix

Every system is run on the same held-out tasks so the comparison is paired.
Five systems, two of them ceilings that bound what is possible rather than what
is deliverable:

- S1, naive baseline: the frozen single best model, one greedy sample. The floor
  of the comparison and the least interesting number.
- S2, strong same-model baseline: the frozen single best model run N times and
  self-aggregated (self-consistency or best-of-N under the same selector the
  ensemble uses). This is the real competitor to diversity, because it spends the
  same sampling budget on one model instead of many. Beating S1 but not S2 means
  more compute helped, not variety.
- S3, diversity ensemble: N different models from the pool, one sample each,
  aggregated by the selector. The treatment for H-ensemble. Compared against S2,
  not S1.
- S4, practical routing: a router that picks one model per task using only
  information available before the answer is graded (task features, a cheap
  classifier, or a prompted router), then runs that one model. The treatment for
  H-route.
- S5, oracle ceilings: oracle routing (per task, pick whichever single model
  turned out best) and oracle selection (from the ensemble's N answers, pick
  whichever turned out best). These use the grade to choose and therefore cannot
  be deployed. They are reported only as the upper bound on exploitable variety.
  If even the oracle barely beats S2, there is no variety to capture and S3 and
  S4 cannot rescue it.

The decisive contrasts are S3 versus S2 and S4 versus S1, with S5 as the ceiling
that tells us whether any practical method left value on the table.

## The model pool

Floor-class open-weight models per ADR-0006, deliberately spread across base
families and specializations so the diversity is real and not five fine-tunes of
one base. The bundle mirrors the reliable-floor bundle: general reasoning models
in the 14B-to-30B band, at least one dedicated code model, and enough family
spread (for example a Qwen-family model, a Gemma-family model, a Mistral-family
model, a Llama-family model, and a code specialist) that inter-model answer
disagreement is genuine. All at the Q4_K_M default. The exact roster is picked
from the current exemplars in Reliable_Floor_Definition.md at run time and
recorded with the results, since the roster ages but the protocol does not.

A manipulation check on diversity: measure how often the models actually disagree
on the evaluation tasks. If they rarely disagree, the pool is not diverse and any
null result is about the roster, not the thesis, and the roster is what changes.

## The task suite

Anchored to load-bearing cognitive work for the target user in Who_DII_Is_For,
the independent professional and small shop, weighted toward objective grading so
the result does not rest entirely on a judge model:

- Coding. Real tasks graded by unit tests. Fully objective. The aggregated answer
  either passes the tests or it does not.
- Multi-step reasoning. Tasks with checkable answers: quantitative word problems,
  multi-hop logic, small data-analysis questions with a numeric or exact result.
  Objective.
- Long-context document work. Extraction, faithful summarization, and question
  answering over long documents, with reference answers and checkable facts.
  Mixed objective and judged.
- A privacy-sensitive domain slice. Question answering in a field that matches the
  exposure profile (a data-residency-constrained professional domain), graded
  against a rubric and reference answers. Mostly judged, and the hardest to score,
  so it is the smallest weighted slice and never the sole basis for a verdict.

Tasks favor recent or transformed items over famous public benchmarks to blunt
contamination, and every task is fixed and versioned before the run.

## Metrics and the decision rule

Primary quality metric per domain: task success rate where a ground truth exists
(tests pass, answer matches), and a rubric score from a judge where it does not.
Systems are compared paired, per task, and reported as win/loss/tie counts and
mean score difference with a bootstrap confidence interval on the difference.

The noise floor is measured, not assumed. The single best model is rerun at
generation temperature enough times to estimate within-system variance, and any
between-system margin must exceed that noise floor to count. A variety arm that
beats the baseline by less than the model's own run-to-run wobble has not beaten
it.

Cost is a first-class axis, because a pod win only matters if it is cheaper to
the member than a metered call for the latency-tolerant tail (the red team's
economics, ADR-0006's best-effort ceiling). Every system reports quality per unit
of compute and per unit of wall-clock latency. A method that needs five models to
gain two points over one model is recorded as such, and whether that trade is
worth it is judged explicitly, not hidden inside an average.

## The judging protocol, and its own controls

Ground truth is used wherever it exists and the judge is reserved for the
open-ended remainder. Where a judge is used: a strong model that is not in the
candidate pool, to avoid a model scoring its own family; pairwise comparison with
answer positions randomized and swapped, discarding any pair where the verdict
flips on swap; a fixed rubric rather than a bare "which is better"; and blind
scoring with model identities stripped. Judge trust is itself validated: a small
subset is scored by hand and the judge is only relied on for the rest if its
agreement with the human subset is high. If the judge cannot be trusted on a
domain, that domain falls back to objective-only tasks or is dropped, rather than
resting the verdict on an unvalidated judge.

## Confounds and how each is killed

- Baseline weakened after the fact. Killed by the split-and-freeze selection rule
  and the held-out evaluation split.
- Oracle leakage dressed up as a result. Killed by labeling S5 as a ceiling and
  never reporting it as achievable.
- Size masquerading as variety. Killed by S2, the same-model best-of-N control,
  and by keeping the largest model as a baseline candidate. Variety must beat
  equal compute spent on one model, not just beat one cheap sample.
- Judge bias, verbosity, self-preference, position. Killed by the judging
  controls above and by grounding as much of the suite as possible in objective
  tests.
- Data contamination. Blunted by preferring recent or transformed tasks and by
  including novel items, with public-benchmark leakage risk noted per domain.
- Fake diversity. Caught by the inter-model disagreement manipulation check; a
  null over a pool that never disagreed is a roster finding, not a thesis finding.

## Pre-registered thresholds

Fixed now, before the run:

- Sample size: on the order of 50 to 100 tasks per domain across four domains, for
  a few hundred paired tasks total, enough that a real effect clears the noise
  floor and small enough to run on a pod.
- The "meaningful" margin is fixed at 0.05 (five points of task success) above the
  measured noise floor: a contrast counts only if its paired mean difference beats
  noise_floor + 0.05 and its bootstrap CI excludes zero. Written down here before
  results are seen, not chosen to make the outcome land a particular way. The margin
  may be set per domain (config `margin_by_domain`); absent an override the 0.05
  default applies to every domain, and the noise floor is measured per domain.
- Primary verdict is on the objective-heavy domains (coding, reasoning,
  long-context); the privacy slice informs but cannot by itself decide.

## Success and kill criteria

Read against the red team's own kill test that "model variety stops paying":

- Variety pays. At least one of S3 (over S2) or S4 (over S1) clears the margin on
  a majority of the objective-heavy domains, net of cost. The durable leg is
  supported; messaging can lead with variety; record which form won, in which
  domains, and at what cost multiplier, since a domain-specific win is still real
  and tells the pod which tail it owns.
- Variety is marginal. The arms beat S1 but not S2, or beat S2 only within the
  noise floor. The gain is compute, not diversity. This is a serious finding
  against the current pitch: variety demotes from thesis to a weak, best-effort
  feature and the outward messaging must stop leaning on it.
- Variety is uncaptured. Neither practical arm (S3 over S2, S4 over S1) clears the
  margin, but the S5 oracle does beat S2 by the margin: the pool genuinely holds
  exploitable variety that the router and ensemble failed to capture. This is not a
  loss — it says the diversity is real and the work is a better selector or router,
  not that variety is absent. Reported as UNCAPTURED, kept distinct from LOSES so a
  method failure is never mistaken for an absence of variety. (Added to the
  pre-registration before the full run; the rig implements it directly.)
- Variety loses. Neither arm beats S2 outside noise, and even the S5 oracle is
  close to S2 (within the margin). The durable leg is not durable on these
  workloads. Per the red team, this is one of the three conditions that put the
  project's reason to exist in real question, and it gets said plainly rather than
  managed. The oracle-is-close condition is what separates LOSES from UNCAPTURED.

## What this settles and what it does not

It settles whether variety pays on these workloads, at floor scale, with mid-2026
open-weight models. It is one dated data point. It does not settle durability over
time, since the model population changes, which is why the red team tracks leading
indicators rather than a single verdict; the experiment is meant to be rerun as
the roster ages. It does not touch the size leg (already conceded) or the
access-under-conditions and sovereignty legs, which are separate arguments and do
not rise or fall with this result. A variety win does not prove the sovereignty
thesis; a variety loss does not sink it, it removes one of two surviving legs and
throws more weight on access-under-conditions.

## How it runs later, not built here

When approved, the run extends the existing measurement harness (prototype's
cmd/harness) from latency into quality scoring and executes on the real pod
substrate that M4 already stood up, floor-class models on owned hardware, no
hosted stand-ins. That build is a separate session; this document is only the
protocol it will follow.
