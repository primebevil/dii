# The Case Against DII: A Red-Team Teardown of the Core Hypotheses

Status: Draft v0.1 (2026-07-02). An adversarial review and falsification exercise, commissioned by the project owner.

## Why this document exists

The Charter says success means understanding, and the research principles commit the project to letting evidence change conclusions. Both oblige us to try to kill the idea before investing in it, and this document is that attempt. The rule here is inverted from the rest of the repository: nothing is credited to DII and everything is doubted. Where the Charter builds the idea up, this builds the strongest honest case that it will not work, grounded in what actually happened to the projects that tried, as of mid-2026. If the concept survives this, that is worth something. If it does not, the project has still succeeded on its own terms by producing understanding cheaply.

Each objection below states a claim, gives the evidence, explains why it bites DII specifically, and says what would have to be true for DII to survive it. That last part matters, because a good adversary is falsifiable too.

## Objection 1: the physics of distributed inference is hostile

The bottleneck in modern LLM inference is not compute, it is memory bandwidth. Distributing work across home hardware attacks the one axis that is already the constraint, and it adds a second bottleneck, the internet, on top.

The evidence is concrete. Token generation, the decode phase, is bound by memory bandwidth rather than FLOPS. To generate one token from a 70B-parameter model at FP16, the machine must read the entire model, roughly 140 GB of weights, from memory for every token. On a datacenter H100 with 3.35 TB/s of bandwidth that read alone is a floor of about 42 milliseconds per token, and adding more compute cannot lower it because the wall is bandwidth. Consumer cards are far behind on exactly that axis: an RTX 4090 has 24 GB of VRAM and a 5090 has 32 GB, at roughly a third to a quarter of the H100's bandwidth. A large model does not even fit on a single consumer card, so it has to be split across machines, and the moment activations cross the public internet they travel over a channel two to three orders of magnitude slower than a datacenter interconnect. Petals is the honest proof of concept here. It works, with consumer GPUs each holding a slice of the layers in a BitTorrent style, but it runs at roughly one step per second, is throughput-poor compared to centralized serving, and in current builds is Linux-only. It is genuinely useful for latency-tolerant single-stream use, and it is not competitive for production throughput.

This bites DII because moving up a layer of abstraction does not repeal the memory wall or the speed of light. "The scheduler is the product" still has to place work onto physical machines connected by real networks. Any capability that needs a large model, low latency, and high throughput at once is structurally hard to serve on distributed home hardware, and that is precisely the profile of the most valuable workloads.

DII survives this only if the target envelope is narrower than "intelligence infrastructure." The viable set is latency-tolerant, embarrassingly parallel, or small-model work where each request fits on one node and nodes do not need to talk mid-inference: batch reasoning, async agents, bulk embedding and retrieval, fine-tuning jobs, and evals. If the target is that set, the physics is survivable. If the target is interactive frontier-scale inference, it is not.

## Objection 2: the wall is demand, not idle supply

The founding image, that millions of machines sit idle so we should use them, assumes the scarce resource is supply. In this market the scarce resource is paying demand, and aggregating more idle supply makes the real problem worse.

The evidence is stark. io.net registered about 327,000 GPUs, but in Q1 2025 the daily average of verified active GPUs was about 6,720, roughly 2 percent utilization of the registered base. Client purchases ran around 2.6 percent of supplier payouts in a sampled period, and third-party trackers put annualized revenue near 12.5 million dollars in mid-2026, well below the run-rate its earlier self-reported figures implied. The supply showed up and the demand did not. Historically, every attempt to commercialize volunteer computing, whether by paying participants, running lotteries, or reselling cycles, has failed, and BOINC's volunteer population has sat flat around 500,000 for years. Idle capacity is abundant and nearly worthless; what is scarce is someone willing to pay to use it reliably.

This bites DII because its entire on-ramp is contribution of idle hardware, which is the easy side of the market to build and the side that is already oversupplied everywhere it has been tried. The hard side, durable and paying production demand that tolerates the reliability profile of consumer hardware, is the side DII has said the least about.

DII survives this only if there is a demand population that specifically prefers distributed or home compute for a reason other than price: privacy and data-residency, where compute stays on hardware you control; sovereignty, which is the Charter's own founding concern about government-gated frontier AI; or access when centralized providers are unavailable or rate-limited. If the only pitch is cheaper GPUs, DII competes on the one axis where hyperscalers have structural advantages, and loses.

## Objection 3: the verification tax may eat the savings

You cannot trust an untrusted node to have actually done the work. The cheap defense, running everything twice, doubles cost and erases the economic premise. The sophisticated defenses are real but young, and they impose their own constraints.

The evidence sits in the state of the art. The lazy-node problem is fundamental: a node can return plausible garbage or a cached lie and collect the reward, and re-running the whole job to check defeats the purpose. Standard volunteer-computing practice computes every task at least twice for validation, an explicit and permanent tax on efficiency. The frontier defenses show how hard this is. Gensyn's Verde achieves trustless verification only by requiring bitwise-reproducible operators across heterogeneous hardware, plus a dispute-resolution referee that recomputes the first disagreeing node in the graph. Prime Intellect's TOPLOC uses locality-sensitive hashing to detect tampering with model, prompt, or computation. These are impressive, and they exist precisely because naive verification is unaffordable.

This bites DII because bitwise reproducibility across a zoo of consumer GPUs, drivers, and OS versions is a demanding invariant, and the more heterogeneous the node pool, which is DII's stated strength, the harder it is to guarantee. It also compounds with Objection 1: if you defend correctness by redundant execution, you pay the memory-bandwidth-bound cost two or more times, on hardware that was already slower per token. Trust and economics pull against each other.

DII survives this by adopting an existing verification protocol rather than inventing one, by restricting to workloads whose results are cheaply self-verifying because they are checkable far faster than they are produced, or by building trust from identity and reputation rather than cryptographic proof and accepting a smaller, semi-trusted node pool. Each choice shrinks the "anyone can join" vision.

## Objection 4: proof-of-contribution has a known failure mode

Rejecting tokens and crypto was the right instinct, but "contribute compute, earn credits, spend later" is a closed-loop barter economy, and those have a predictable way of stalling.

The reasoning is straightforward. A non-transferable credit is worth only the compute you personally will later consume, so the moment a contributor's own demand is satisfied, the incentive to keep contributing falls to zero and supply dries up exactly when the network needs depth. Make the credit transferable to fix that, and you have re-created a currency, with pricing, speculation, and arbitrage, which is the crypto dynamic the project rejected. This is the same reef that sank paid and lottery volunteer-computing experiments. There is a real tension between "not financial" and "liquid enough to sustain participation," and it is not obviously resolvable.

This bites DII because proof-of-contribution is load-bearing for the whole participation story. If credits do not circulate, the flywheel does not spin. If they do circulate freely, the project drifts into the exact category its own "What This Is Not" section rules out.

DII survives this with a deliberately illiquid design: credits that expire, cannot be transferred, and function as a fair-use and reciprocity mechanism rather than a currency, accepting that this caps growth. Alternatively, the primary motivation is not the credit at all but mission, research access, or community, with credits as a secondary fairness knob. This deserves its own ADR early, because it shapes who participates and why.

## Objection 5: "the scheduler is the product" may describe a feature, not a company

Orchestrating capability across a heterogeneous fleet is a real problem, but it may be a thin layer that compute providers, model runtimes, and existing schedulers absorb, rather than a defensible independent product.

The evidence is in who is earning money. The GPU networks actually taking revenue in 2025 and 2026, led by Akash, got there by pivoting away from token-subsidized idealism toward professional, largely datacenter-grade supply and real cash flow, with Akash reporting around 50 percent utilization, single-digit-millions revenue, and utilization trending high into 2026. That supply is underused professional GPUs, not idle home machines. Meanwhile mature orchestration such as Ray and Kubernetes-based schedulers already handles capability-aware placement, and Petals shows the routing logic itself is not the moat; the hard parts are the physics, the trust, and the demand. A scheduler with no unique substrate underneath it is a coordination convenience that others can replicate.

This bites DII because if the substrate, home hardware, is uncompetitive for the valuable workloads per Objection 1, and the demand is not there per Objection 2, then "the scheduler is the product" risks being a well-designed control plane on top of a market that does not need it. The most successful analog did not win by scheduling idle consumer GPUs; it won by sourcing better supply and chasing real demand.

DII survives this if the scheduler has a defensible reason to exist that incumbents will not replicate: hybrid placement across owned, rented, and community fleets with privacy and sovereignty constraints as first-class routing inputs rather than just cost and latency. That is a real enterprise pain point and is arguably where the genuine idea hides, but it is a narrower product than "internet-native infrastructure for intelligence."

## How the objections compound

The objections are not independent, and that is the most important finding. Defending trust in Objection 3 by redundant execution multiplies the memory-bandwidth cost in Objection 1. The physics in Objection 1 narrows viable workloads to latency-tolerant ones, which are also the least differentiated and easiest for cheap centralized batch compute to serve, worsening the demand problem in Objection 2. Fixing the incentive stall in Objection 4 by making credits liquid pushes the project into the crypto category it rejected, which is also where demand for decentralized compute has been most speculative and least real. Solving any one in isolation tends to make another worse. A viable DII is therefore one that finds a single point in the design space where all of them are simultaneously survivable. That point exists, but it is small, and it is not the maximal vision in the Charter.

## Where the concept is not dead

Steelmanning the survivor, there is a coherent and defensible version of DII that lives inside the intersection of every "what would have to be true" above. It is a trust-and-privacy-first scheduler that routes latency-tolerant, node-local intelligence workloads across hybrid hardware combining owned, community, and rented machines, motivated primarily by sovereignty and data-residency rather than price, using borrowed verification in the TOPLOC or Verde family and a deliberately illiquid reciprocity-credit model, and targeting users who cannot or will not send their data to a hyperscaler. That is smaller than "one computer made of millions of machines," but it is directly downstream of the Charter's actual founding question about what happens to everyone else when frontier AI becomes government-gated infrastructure. The sovereignty angle, rather than the cost angle, is the one the evidence leaves standing.

## Kill criteria: how to settle this in the Week 3 prototype

The two-node proof of concept should not aim to make it work. It should measure the numbers that decide whether the small viable envelope is real.

1. Tokens per second and time-to-first-token for a target model split across two consumer nodes over a residential network, not a LAN, compared to a single centralized instance. If distributed is more than ten times worse on the workloads you actually care about, the physics objection is confirmed for that use case.
2. Verification overhead. Implement the cheapest credible check and measure the real multiplier on cost and latency. If it is two times or more, model whether any workload still pencils out.
3. A demand probe, not a supply probe. Before building more, find five people or organizations who would pay or commit to run something on this specifically for privacy or sovereignty reasons. If you cannot find five, Objection 2 stands.
4. Write the proof-of-concept credit ADR and stress-test it. Does participation survive a contributor whose own demand is already met?

Each is a cheap experiment that can falsify a load-bearing hypothesis, which is the fastest path to the Charter's definition of success.

## Verdict

The maximal vision, distributed home-hardware-powered general intelligence infrastructure that competes as a cost-driven compute layer, is contradicted by current evidence on all three of physics, demand, and trust-economics, and the compounding interactions make it very unlikely to work. The narrow version, a sovereignty and privacy-first scheduler for latency-tolerant workloads across hybrid fleets, is not refuted by the evidence and is well aligned with the project's real founding question. It is a smaller and different thing than the Charter's headline, and it should be treated as the hypothesis worth prototyping.

The project's own framing, that the owner wants to say he did it rather than own the idea, means this verdict is good news. Narrowing to the survivable core, and knowing why the rest does not hold, is exactly the understanding the Charter set out to produce.

## Sources

- [Petals (GitHub, bigscience-workshop)](https://github.com/bigscience-workshop/petals)
- [Petals: decentralized inference and finetuning of LLMs, Yandex Research](https://research.yandex.com/blog/petals-decentralized-inference-and-finetuning-of-large-language-models)
- [AI's Memory Wall Problem, Spheron (2026)](https://www.spheron.network/blog/ai-memory-wall-inference-latency-guide-2026/)
- [NVIDIA H100 GPU: Specs, VRAM, Price, RunPod](https://www.runpod.io/articles/guides/nvidia-h100)
- [DePIN's Revenue Reckoning: Akash, io.net, Aethir, BlockEden](https://blockeden.xyz/blog/2026/03/12/depin-compute-revenue-pivot-akash-ionet-aethir/)
- [State of Akash Q2 2025, Messari](https://messari.io/report/state-of-akash-q2-2025)
- [io.net Review: Solana GPU Marketplace, Own Your Mind](https://ownyourmind.ai/projects/io-net/)
- [Verde: a verification system for machine learning over untrusted nodes, Gensyn](https://blog.gensyn.ai/verde-a-verification-system-for-machine-learning-over-untrusted-nodes/)
- [Optimistic Verifiable Training by Controlling Hardware Nondeterminism (arXiv)](https://arxiv.org/html/2403.09603v3)
- [Our approach to decentralized training, Prime Intellect](https://www.primeintellect.ai/blog/our-approach-to-decentralized-training)
- [Volunteer computing: the ultimate cloud, BOINC (Berkeley)](https://boinc.berkeley.edu/boinc_papers/crossroads.pdf)
- [Distributed Deep Learning in Open Collaborations (arXiv)](https://arxiv.org/pdf/2106.10207)
