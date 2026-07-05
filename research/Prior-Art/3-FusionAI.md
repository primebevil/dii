# FusionAI

## What it does

FusionAI is a 2023 research paper, not a running system: "FusionAI: Decentralized Training and Deploying LLMs with Massive Consumer-Level GPUs" (arXiv 2309.01172), from a group led by Xiaowen Chu and Shaohuai Shi across Hong Kong Baptist University, HKUST(GZ), and others. It proposes a design for using idle consumer GPUs, such as RTX 3080s, to pre-train, fine-tune, and serve LLMs while protecting privacy, and it attacks the three classic problems of volunteer compute: device heterogeneity, unreliable nodes that join and quit, and low bandwidth. Its design pillars are a broker with a backup pool for churn, hardware-aware scheduling, and abstracting ML tasks into directed acyclic graphs so heterogeneous devices and framework versions can interoperate. The same authors later implemented the vision as FusionLLM (arXiv 2410.12707, 2024), which is the artifact with real measurements.

Note on the name: an unrelated crypto token also uses "Fusion AI." It is not this.

## Verdict

Stalled as a named system, won as a research thread. There is no FusionAI network, code release, or user base under that name, but its ideas were validated by the follow-on FusionLLM and are cited as prior art in the field's 2025 survey of decentralized training.

## Why

FusionAI is a design-and-feasibility paper; its wording is explicitly aspirational ("we envision"), and its headline claim, that 50 RTX 3080s can approach the throughput of 4 H100s at far lower cost, is an analytical projection, not an end-to-end measured run. It never shipped as infrastructure, built no community, and added no economic or trust layer, so it lost the running-network race to Petals and later Prime Intellect. What it did establish is a credible blueprint and a clean statement of the hard problems, which FusionLLM then made real: on 48 heterogeneous GPUs across networks from 8 Mbps to 10 Gbps, it reported 1.45x to 9.39x speedups while preserving convergence, using an operator-level DAG with remote autodiff (so nodes need not share framework versions), a bandwidth-clustering scheduler, and adaptive Top-K compression applied only at the slowest links.

## How it compares to DII

FusionAI is useful to DII precisely as a systems-design source with none of the parts DII most needs to get right. It is a pure compute-systems proposal: no incentive mechanism, no Sybil resistance, no payment, no verification of honest work, no community. DII is the opposite emphasis. Its hardest problems are the participation model, trust, and demand, and its systems layer serves a modest reliable floor rather than chasing large-model training throughput.

So what DII offers that FusionAI does not is the entire economic, trust, and community layer that turns a design into a network that persists, plus a concrete, narrow target (latency-tolerant inference for dependency-plus-exposure users) rather than an open-ended feasibility argument. DII should mine FusionAI and FusionLLM for engineering tactics and supply the layers they omit.

## Borrow / Avoid

Borrow: the operator-level DAG with remote autodiff, which lets nodes run different framework versions, a major practical win for a heterogeneous volunteer fleet; the bandwidth-clustering scheduler that cuts the pipeline where links are best; selective compression sized by measured bandwidth rather than uniform compression that hurt convergence; a broker plus backup pool for churn; and honest evaluation on genuinely heterogeneous testbeds with realistic bad links.

Avoid: citing "50 3080s equal 4 H100s" as proven, since it is modeled, not measured; treating FusionAI as shippable infrastructure, since there is no code or network under that name (use FusionLLM for running code); over-promising near-datacenter efficiency, since even the measured system tops out around 9x on carefully partitioned work; and assuming the systems design is the whole stack, since it solves none of Sybil resistance, payment, verification, or Byzantine nodes.

## Sources

- [FusionAI (arXiv 2309.01172)](https://arxiv.org/abs/2309.01172)
- [FusionLLM (arXiv 2410.12707)](https://arxiv.org/html/2410.12707v1)
- [Beyond a Single AI Cluster: a survey of decentralized LLM training (arXiv 2503.11023)](https://arxiv.org/html/2503.11023v2)
- [FusionAI, Semantic Scholar](https://www.semanticscholar.org/paper/FusionAI:-Decentralized-Training-and-Deploying-LLMs-Tang-Wang/660380b17d3a37d8132f2e6dcb5cb47092e5b7d1)

Note: no public FusionAI source code or live network exists under that name; the "50 3080s equal 4 H100s" figure is an analytical projection, not a measured result.
