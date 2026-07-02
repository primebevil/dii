# The Case Against DII

## A Red-Team Teardown of the Core Hypotheses

> Status: Draft v0.1 — 2026-07-02
> Type: Adversarial review / falsification exercise
> Author: Red-team pass (assistant), commissioned by project owner

---

## Why this document exists

The Charter says success means *understanding*, and the Research Principles say
"let evidence change conclusions." Both commit the project to trying to kill its
own idea before investing in it. This document is that attempt.

The rule here is inverted from the rest of the repo: **nothing is credited to
DII, everything is doubted.** Where the Charter builds the idea up, this builds
the strongest honest case that it will not work — grounded in what actually
happened to the projects that tried, as of mid-2026. If the concept survives
this, that is worth something. If it does not, the project has succeeded on its
own terms: it produced understanding cheaply.

Each objection below follows the same shape: **the claim**, **the evidence**,
**why it bites DII specifically**, and **what would have to be true for DII to
survive it.** That last part matters — a good adversary is falsifiable too.

---

## Objection 1 — The physics of distributed inference is hostile

**Claim.** The bottleneck in modern LLM inference is not compute; it is memory
bandwidth. Distributing across home hardware attacks the one axis that is
already the constraint, and adds a second bottleneck (the internet) on top.

**Evidence.** LLM token generation (the decode phase) is memory-bandwidth-bound,
not FLOPS-bound. To generate one token from a 70B-parameter model at FP16, the
machine must read the entire ~140 GB of weights from memory *per token*. On a
datacenter H100 (3.35 TB/s HBM3) that read alone is a ~42 ms floor per token —
and *adding more compute cannot lower it*, because the wall is bandwidth, not
math. Consumer cards are far behind on the exact axis that matters: an RTX 4090
has 24 GB of VRAM, a 5090 has 32 GB, at roughly one-third to one-quarter of the
H100's bandwidth. A large model does not even *fit* on a single consumer card,
so it must be split across machines — and the moment activations cross the
public internet, they travel over a channel that is two to three orders of
magnitude slower than a datacenter interconnect.

Petals is the honest proof-of-concept here. It works — consumer GPUs
BitTorrent-style, each holding a slice of the layers — but it runs at roughly
**one step per second**, is throughput-poor compared to centralized serving,
and (as of the current builds) is Linux-only. It is genuinely useful for
latency-tolerant, interactive, single-stream use. It is not competitive for
production throughput.

**Why it bites DII.** "The scheduler is the product" moves the project up an
abstraction layer, but the scheduler still has to place work onto physical
machines connected by real networks. Abstraction does not repeal the memory
wall or the speed of light. Any capability that needs a large model, low
latency, *and* high throughput at once is structurally hard to serve on
distributed home hardware — and that is precisely the profile of the most
valuable workloads.

**What would have to be true for DII to survive.** The viable envelope is
narrower than "intelligence infrastructure": latency-tolerant, embarrassingly
parallel, or small-model workloads where each request fits on one node and
nodes don't need to talk mid-inference (batch reasoning, async agents, bulk
embedding/retrieval, fine-tuning jobs, evals). If DII's target is *that* set,
the physics is survivable. If the target is interactive frontier-scale
inference, it is not.

---

## Objection 2 — The wall is demand, not idle supply

**Claim.** The project's founding image — "millions of machines sit idle, let's
use them" — assumes the scarce resource is supply. In this market the scarce
resource is *paying demand*. Aggregating more idle supply makes the real problem
worse.

**Evidence.** io.net registered ~327,000 GPUs, but in Q1 2025 the daily average
of verified active GPUs was ~6,720 — about **2% utilization** of the registered
base. Client purchases ran around 2.6% of supplier payouts in a sampled period,
and third-party trackers put annualized revenue (~$12.5M in mid-2026) well below
the run-rate its earlier self-reported figures implied. The supply showed up;
the demand did not. Historically, every attempt to commercialize volunteer
computing — paying participants, lotteries, reselling cycles — has failed, and
BOINC's volunteer population has sat flat around 500,000 for years. Idle
capacity is abundant and nearly worthless; what's scarce is someone willing to
pay to use it reliably.

**Why it bites DII.** DII's entire on-ramp is contribution of idle hardware.
That is the easy side of the market to build and the side that is already
oversupplied everywhere it's been tried. The hard side — durable, paying,
production demand that tolerates the reliability profile of consumer hardware —
is the side DII has said the least about.

**What would have to be true for DII to survive.** There must be a demand
population that specifically *prefers* distributed/home compute for a reason
other than price: privacy/data-residency (compute stays on hardware you
control), sovereignty (the Charter's own founding concern about
government-gated frontier AI), or access when centralized providers are
unavailable/rate-limited. If the only pitch is "cheaper GPUs," DII is competing
on the one axis where hyperscalers have structural advantages, and loses.

---

## Objection 3 — The verification tax may eat the savings

**Claim.** You cannot trust an untrusted node to have actually done the work. The
cheap defense — run everything twice — doubles cost and erases the economic
premise. The sophisticated defenses are real but young, and they impose their
own constraints.

**Evidence.** The "lazy node" problem is fundamental: a node can return
plausible garbage (or a cached lie) and collect the reward, and re-running the
whole job to check defeats the purpose. Standard volunteer-computing practice is
to compute every task at least twice for validation — an explicit, permanent
tax on efficiency. The frontier defenses are moving fast but reveal how hard
this is: Gensyn's **Verde** achieves trustless verification only by requiring
*bitwise-reproducible* operators (RepOps) across heterogeneous hardware, plus a
dispute-resolution referee that recomputes the first disagreeing node in the
graph; Prime Intellect's **TOPLOC** uses locality-sensitive hashing to detect
tampering with model/prompt/computation. These are impressive — and they exist
precisely *because* naive verification is unaffordable.

**Why it bites DII.** Bitwise reproducibility across a zoo of consumer GPUs,
drivers, and OS versions is a demanding invariant; the more heterogeneous your
node pool (DII's stated strength), the harder it is to guarantee. And this
objection **compounds with Objection 1**: if you defend correctness by redundant
execution, you pay the memory-bandwidth-bound cost two or more times, on
hardware that was already slower per token. Trust and economics pull against
each other.

**What would have to be true for DII to survive.** Either adopt an existing
verification protocol rather than inventing one (Verde/TOPLOC-style), *or*
restrict to workloads where the result is cheaply self-verifying (e.g. outputs
checkable far faster than they're produced), *or* build trust from identity and
reputation rather than cryptographic proof and accept a smaller, semi-trusted
node pool. Each choice shrinks the "anyone can join" vision.

---

## Objection 4 — Proof-of-Contribution has a known failure mode

**Claim.** Rejecting tokens/crypto was the right instinct, but "contribute
compute, earn credits, spend later" is a closed-loop barter economy, and those
have a predictable way of stalling.

**Evidence & reasoning.** A non-transferable credit is only worth the compute
*you personally* will later consume. The moment a contributor's own demand is
satisfied, their incentive to keep contributing goes to zero, and supply dries
up exactly when the network needs depth. Make the credit transferable or
tradable to fix that, and you have re-created a currency — with pricing,
speculation, and arbitrage, i.e. the crypto dynamics the project explicitly
rejected. This is the same reef that sank paid/lottery volunteer-computing
experiments. There is a genuine tension between "not financial" and "liquid
enough to sustain participation," and it is not obviously resolvable.

**Why it bites DII.** PoC is load-bearing for the whole participation story. If
credits don't circulate, the flywheel doesn't spin; if they do circulate freely,
the project drifts into the exact category ("cryptocurrency project") its own
"What This Is Not" section rules out.

**What would have to be true for DII to survive.** The credit system probably
needs a deliberately *illiquid* design — credits that expire, are
non-transferable, and function more like a fair-use / reciprocity mechanism than
a currency — accepting that this caps growth. Or the primary motivation isn't
the credit at all (mission, research access, community), with credits as a
secondary fairness knob. This deserves its own ADR early, because it shapes who
participates and why.

---

## Objection 5 — "The scheduler is the product" may describe a feature, not a company

**Claim.** Orchestrating capability across a heterogeneous fleet is a real
problem — but it may be a thin layer that the compute providers, model runtimes,
and existing schedulers absorb, rather than a defensible independent product.

**Evidence.** The GPU networks that are actually earning money in 2025–26
(Akash, and the broader DePIN-compute set) got there by pivoting *away* from
token-subsidized idealism toward professional, largely datacenter-grade supply
and real cash flow — Akash reporting ~50% utilization and single-digit-millions
revenue, trending toward high utilization into 2026. Note what that supply *is*:
it is not idle home machines, it is underused professional GPUs. Meanwhile
mature orchestration (Ray, Kubernetes-based schedulers, model-serving stacks)
already handles capability-aware placement, and Petals shows the routing logic
itself is not the moat — the hard parts are the physics, the trust, and the
demand. A scheduler with no unique substrate underneath it is a coordination
convenience that others can and do replicate.

**Why it bites DII.** If the substrate (home hardware) is uncompetitive for the
valuable workloads (Objection 1) and the demand isn't there (Objection 2), then
"the scheduler is the product" risks being a well-designed control plane sitting
on top of a market that doesn't need it. The most successful analog didn't win
by scheduling idle consumer GPUs; it won by sourcing better supply and chasing
real demand.

**What would have to be true for DII to survive.** The scheduler needs a
defensible reason to exist that incumbents *won't* replicate: hybrid placement
across owned + rented + community fleets with privacy/sovereignty constraints as
first-class routing inputs (not just cost/latency). That is a real enterprise
pain point and is arguably where the genuine, fundable idea hides — but it's a
different, narrower product than "internet-native infrastructure for
intelligence."

---

## How the objections compound

The objections are not independent, and that's the most important finding. They
form a vise:

- Defending **trust** (Obj. 3) by redundant execution multiplies the
  **memory-bandwidth cost** (Obj. 1).
- The **physics** (Obj. 1) narrows viable workloads to latency-tolerant ones,
  which are also the *least* differentiated and easiest for cheap centralized
  batch compute to serve — worsening the **demand problem** (Obj. 2).
- Fixing the **incentive stall** (Obj. 4) by making credits liquid pushes the
  project into the **crypto category** it rejected, which is also where the
  **demand for decentralized compute has been most speculative and least real**.

Solving any one in isolation tends to make another worse. A viable DII is not
one that beats each objection separately; it's one that finds a single point in
the design space where all of them are simultaneously survivable. That point
exists, but it is *small*, and it is not the maximal vision in the Charter.

---

## Where the concept is NOT dead

Steelmanning the survivor: there is a coherent, defensible version of DII that
lives inside the intersection of every "what would have to be true" above:

> A trust-and-privacy-first scheduler that routes **latency-tolerant,
> node-local** intelligence workloads across **hybrid** (owned + community +
> rented) hardware, motivated primarily by **sovereignty and data-residency**
> rather than price, using **borrowed** verification (Verde/TOPLOC-style) and a
> deliberately **illiquid** reciprocity-credit model, targeting users who
> *cannot* or *will not* send their data to a hyperscaler.

That is smaller than "one computer made of millions of machines," but it is
directly downstream of the Charter's *actual* founding question — what happens
to everyone else when frontier AI becomes government-gated infrastructure. The
sovereignty angle, not the cost angle, is the one the evidence leaves standing.

---

## Kill criteria — how to settle this in the Week 3 prototype

The point of the two-node PoC should not be "make it work." It should be to
measure the numbers that decide whether the small viable envelope is real:

1. **Tokens/sec and time-to-first-token** for a target model split across two
   consumer nodes over a *residential* network — not LAN. Compare to a single
   centralized instance. If distributed is >10x worse on the workloads you
   actually care about, the physics objection is confirmed for that use case.
2. **Verification overhead** — implement the cheapest credible check and measure
   the real multiplier on cost/latency. If it's ≥2x, model whether any workload
   still pencils out.
3. **A demand probe, not a supply probe** — before building more, find 5 people
   or orgs who would *pay or commit* to run something on this specifically for
   privacy/sovereignty reasons. If you can't find 5, Objection 2 stands.
4. **Write the PoC credit ADR** and stress-test it: does participation survive a
   contributor whose own demand is already met?

Each is a cheap experiment that can *falsify* a load-bearing hypothesis. That is
the fastest path to the Charter's definition of success.

---

## Verdict

The maximal vision — distributed, home-hardware-powered, general intelligence
infrastructure that competes as a cost-driven compute layer — is contradicted by
current evidence on all three of physics, demand, and trust-economics, and its
compounding interactions make the maximal version very unlikely to work.

A *narrow* version — a sovereignty/privacy-first scheduler for latency-tolerant
workloads across hybrid fleets — is not refuted by the evidence and is well
aligned with the project's real founding question. It is a smaller and different
thing than the Charter's headline, and it should be treated as the hypothesis
worth prototyping.

The project's own framing — "I don't want to own the idea, I want to say I did
it" — means this verdict is not bad news. Narrowing to the survivable core, and
knowing *why* the rest doesn't hold, is exactly the understanding the Charter set
out to produce.

---

## Sources

- [Petals (GitHub, bigscience-workshop)](https://github.com/bigscience-workshop/petals)
- [Petals: decentralized inference and finetuning of LLMs — Yandex Research](https://research.yandex.com/blog/petals-decentralized-inference-and-finetuning-of-large-language-models)
- [AI's Memory Wall Problem: Why More GPUs Don't Fix Inference Latency (2026) — Spheron](https://www.spheron.network/blog/ai-memory-wall-inference-latency-guide-2026/)
- [NVIDIA H100 GPU: Specs, VRAM, Price — RunPod](https://www.runpod.io/articles/guides/nvidia-h100)
- [DePIN's Revenue Reckoning: Akash, io.net, Aethir — BlockEden](https://blockeden.xyz/blog/2026/03/12/depin-compute-revenue-pivot-akash-ionet-aethir/)
- [State of Akash Q2 2025 — Messari](https://messari.io/report/state-of-akash-q2-2025)
- [io.net Review: Solana GPU Marketplace, tokenomics — Own Your Mind](https://ownyourmind.ai/projects/io-net/)
- [Verde: a verification system for machine learning over untrusted nodes — Gensyn](https://blog.gensyn.ai/verde-a-verification-system-for-machine-learning-over-untrusted-nodes/)
- [Optimistic Verifiable Training by Controlling Hardware Nondeterminism (arXiv)](https://arxiv.org/html/2403.09603v3)
- [Our approach to decentralized training — Prime Intellect](https://www.primeintellect.ai/blog/our-approach-to-decentralized-training)
- [Volunteer computing: the ultimate cloud — BOINC (Berkeley)](https://boinc.berkeley.edu/boinc_papers/crossroads.pdf)
- [Distributed Deep Learning in Open Collaborations (arXiv)](https://arxiv.org/pdf/2106.10207)
