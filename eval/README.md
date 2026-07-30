# eval — variety aggregation experiment rig

A disposable quality-benchmark harness. It answers one question:
docs/Variety_Experiment_Plan.md's open one. On the work DII's target users do,
does a pool of diverse floor-class models beat the single best model in that pool?

This is a measuring instrument, not product code. It ships nothing in the node,
lives walled off in its own directory, and can be deleted without touching the
prototype. It talks to a pod exactly like `scripts/ask.sh` does: an
OpenAI-compatible POST to `/v1/chat/completions` with a model name, and the
router decides which node serves it.

## What it measures

Five systems over the same held-out tasks:

- `S1_naive` — the frozen single best model, one greedy sample. The baseline.
- `S2_same_N` — the same best model run N times and aggregated. The compute
  control: it spends the ensemble's budget on one model instead of many.
- `S3_ensemble` — one sample from each of N diverse models, same aggregator.
- `S4_route` — a practical per-domain router (domain label only, frozen on dev).
- `S5a/S5b` — oracle ceilings (pick the best after grading). Not deployable;
  they bound what any method could have captured.

The decisive contrasts are **S3 vs S2** (does diversity beat equal compute) and
**S4 vs S1** (does routing beat one generalist). A contrast counts only if its
bootstrap CI excludes zero and it clears the measured noise floor. The verdict is
PAYS, MARGINAL (the gain was compute, not diversity), or LOSES, exactly per the
plan.

## Setup

    pip install -r requirements.txt        # pyyaml, the only dependency
    cp config.example.yaml config.yaml     # then edit for your pod's roster

Pull the pool models on your nodes first, and set `judge_model` to a strong model
that is NOT in the pool (no model grades its own family).

## Run it

    python3 run.py --config config.yaml --out results.json

The scoreboard prints to stdout; `results.json` holds every per-task number.

## Prove it works without real models

    python3 tools/mock_server.py --tasks tasks --port 8091 &
    # point config.node_url at http://127.0.0.1:8091, pool = any fake names, judge = judge:32b
    python3 run.py --config config.smoke.yaml

The mock simulates a pod of uneven models so the whole pipeline exercises end to
end. Its numbers are meaningless as evidence; it only proves the plumbing.
`scripts/smoke.sh` does this in one command.

## Security

The `unit_test` grader executes model-produced code in a subprocess. Run the rig
on a throwaway machine or container, never on anything you care about.

## The task bank

`tasks/<domain>/*.yaml`, one exam question each. Objective graders
(`numeric`/`exact`/`mcq`/`unit_test`) check ground truth; `judge` tasks carry a
rubric and are scored by the independent judge. The seeds here are small and, for
long-context, deliberately short so the rig runs cheaply; **expand the bank and
swap in genuinely long documents before any run you would quote.** The seed set is
enough to shake out the rig, not to draw a conclusion from.

## What a run does and does not settle

See the plan. Short version: a run is one dated data point at floor scale with
today's models. It does not settle durability over time (rerun as the roster
ages), and it does not touch the size or sovereignty legs, which are separate
arguments.
