"""Rig configuration: the pod endpoint, the model pool, the judge, and budgets.

Everything that ages (which models are in the pool, which judge) lives in config,
so the protocol in docs/Variety_Experiment_Plan.md stays fixed while the roster
changes. YAML for authoring comfort; pyyaml is the rig's only dependency.
"""

from __future__ import annotations

from dataclasses import dataclass, field

import yaml


@dataclass
class Config:
    node_url: str = "http://localhost:8080"
    token: str = ""
    # The diverse pool under test: floor-class open-weight model names as the pod
    # advertises them. Deliberately spread across base families (ADR-0006).
    pool: list[str] = field(default_factory=list)
    # The judge must NOT be in the pool, so no model grades its own family.
    judge_model: str = ""
    # Where the judge runs. Empty => same endpoint as the pool (node_url/token).
    # Set these to point the judge at a SEPARATE OpenAI-compatible endpoint (e.g. a
    # hosted frontier model): the plan allows this because the judge is measurement,
    # not part of the network under test. A strong hosted judge sidesteps both the
    # self-family-bias problem and the "no fast, strong pod judge available" problem.
    judge_url: str = ""
    judge_token: str = ""
    # The judge must complete its verdict within this budget. Kept separate from the
    # pool's max_tokens; a reasoning judge needs room to think AND emit the verdict.
    judge_max_tokens: int = 1024
    # Sampling budget N: how many samples the ensemble draws (one per model) and,
    # to match it, how many the same-model best-of-N baseline draws. Equal by
    # construction so S2 vs S3 isolates diversity from compute.
    samples: int = 0  # 0 means "use len(pool)", enforced in load()
    temperature: float = 0.7
    # Fraction of tasks held out for evaluation; the rest is the dev split used to
    # pick the frozen baseline and build the router table.
    eval_fraction: float = 0.5
    # Deterministic split/seed control so a run is reproducible.
    split_seed: int = 1
    # Repeats of the baseline used to estimate the run-to-run noise floor.
    noise_repeats: int = 5
    max_tokens: int = 1024
    # The pre-registered "meaningful margin": the effect size a contrast must clear
    # ABOVE the measured noise floor to count, per docs/Variety_Experiment_Plan.md.
    # This is a scientific choice that must be written down BEFORE a run, not tuned
    # to make a result land. Default 0.0 == "noise floor only", which is NOT a real
    # pre-registration; set a genuine value (and per-domain overrides) before any
    # run you would quote.
    margin: float = 0.0
    margin_by_domain: dict = field(default_factory=dict)

    def margin_for(self, domain: str) -> float:
        return float(self.margin_by_domain.get(domain, self.margin))

    @staticmethod
    def load(path: str) -> "Config":
        with open(path) as f:
            raw = yaml.safe_load(f) or {}
        cfg = Config(**{k: v for k, v in raw.items() if k in Config.__dataclass_fields__})
        if not cfg.pool:
            raise ValueError("config: 'pool' must list at least two models")
        if len(cfg.pool) < 2:
            raise ValueError("config: 'pool' needs at least two models to test variety")
        if cfg.judge_model and cfg.judge_model in cfg.pool:
            raise ValueError(
                f"config: judge_model '{cfg.judge_model}' is in the pool; the judge "
                "must be independent of the candidates"
            )
        if cfg.samples == 0:
            cfg.samples = len(cfg.pool)
        return cfg
