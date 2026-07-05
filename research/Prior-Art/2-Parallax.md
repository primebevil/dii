# Parallax

## What it does

Parallax, built by Gradient (a Singapore AI lab, part of the Solana-based Gradient Network), is a decentralized LLM inference engine: it pools heterogeneous GPUs and Apple Silicon into one virtual serving platform so you can build an AI cluster anywhere. The core mechanism is pipeline parallelism, sharding a model's layers across devices, plus a two-phase scheduler. First a placement phase uses dynamic programming with a water-filling heuristic to lay each replica's layers across diverse devices under memory and bandwidth limits, minimizing pipeline stages and maximizing replicas. Then, at request time, it treats the placement as a graph and stitches layers from different replicas into per-request execution chains that balance load and adapt to live conditions. It wraps mature backends (SGLang and vLLM for GPU, MLX for Mac), runs on Windows, macOS, and Linux, and is Apache-2.0. It was open-sourced around October 2025 and formally launched in November 2025.

## Verdict

Early traction, not yet won. Better engineered than the stalled Petals lineage, but its scale is unproven and its supply is bootstrapped by airdrop farming.

## Why

Parallax is real and moving: steady releases into 2026, an arXiv paper, broad support for 40-plus open models up to very large mixture-of-experts, and a working Apple Silicon path with paged KV cache and continuous batching, which is genuinely differentiated and hard. Its two-phase scheduler is a direct, credible answer to Petals' core weakness, greedy swarm heuristics with no global optimization. Backing is credible, a 10-million-dollar seed co-led by Multicoin and Pantera, with HSG.

The caution flags are about proof and incentives. The published benchmarks are small and self-reported: up to 3.6x throughput and lower latency, but on five-to-fourteen co-located datacenter GPUs at roughly 10 milliseconds latency, with the head-to-head against HexGen rather than Petals, and the "up to 3.6x" headline sitting well above the roughly 1.6x average. Real consumer-broadband, high-churn conditions are under-demonstrated. The GitHub footprint is modest (around 1.3k stars) relative to the marketing. And the large headline participation numbers Gradient cites belong to its browser-based Sentry Node program, an edge-compute and point-farming layer for the broader network, not to Parallax inference nodes, whose actual count is not disclosed. There is no token yet and no disclosed paying-inference workload.

## How it compares to DII

Parallax is DII's closest technical cousin: decentralized inference, pipeline-parallel, heterogeneity-aware, local-first, with first-class Mac support. The divergence is entirely in motivation and governance. Parallax is a commercial DePIN project whose supply is seeded by airdrop-farming incentives and whose eventual coordination runs through a token on Solana. DII is non-commercial and token-free, its supply motivated by mission and community following Tor, and its access never gated on contribution or point-farming.

What DII offers that Parallax does not is a participation model that does not depend on a speculative reward, and honest reporting of what the network actually is. Parallax's headline participation optics come from point-farming Sentry Nodes; DII counts nodes that clear the useful line and serve real users, and would publish real inference capacity and churn rather than headline node counts. Technically the two are close enough that DII should borrow heavily from Parallax; on incentives they are opposites by design.

## Borrow / Avoid

Borrow: pipeline parallelism as the base strategy over slow links (only activations cross the wire); the two-phase scheduler that separates static placement from dynamic request-time routing, the strongest idea here and the fix for Petals' weakness; first-class Apple Silicon and MLX support with paged KV cache and continuous batching; wrapping mature backends rather than rewriting kernels; and straggler-aware incremental re-allocation so churn only disturbs affected nodes. Shipping a paper plus Apache-2.0 code plus a polished launch is the credibility combination worth copying.

Avoid: conflating point-farming node counts with real inference capacity; best-case benchmark framing where the headline is far above the average and figures differ between paper and press; presenting small co-located datacenter testbeds as decentralized; and token or airdrop dependency for bootstrapping, which seeds idle supply but not genuine demand.

## Sources

- [Parallax (GitHub, GradientHQ)](https://github.com/GradientHQ/parallax)
- [Parallax: efficient LLM inference over a decentralized environment (arXiv 2509.26182)](https://arxiv.org/html/2509.26182v1)
- [Gradient launches Parallax (GlobeNewswire, Nov 2025)](https://www.globenewswire.com/news-release/2025/11/06/3182435/0/en/Gradient-Launches-Parallax-a-Sovereign-AI-Operating-System-for-an-Open-Source-Future.html)
- [Gradient $10M seed on Solana (Decrypt)](https://decrypt.co/326140/ai-infrastructure-firm-gradient-bags-10-million-to-develop-protocols-on-solana)
- [Gradient Network and airdrop explainer (CoinGecko)](https://www.coingecko.com/learn/what-is-gradient-network-and-gradient-airdrop)

Note: Parallax's inference-network size is not separately disclosed. Independent re-checking did not find a "1.6 million nodes" figure in the cited source (the nearby real number is 1.6 billion peer connections, not a node count), so no specific participation number is asserted; the point that stands is that the headline counts are Sentry Nodes, not Parallax inference nodes. Benchmarks are self-reported on small testbeds, with HexGen (not Petals) as the direct baseline.
