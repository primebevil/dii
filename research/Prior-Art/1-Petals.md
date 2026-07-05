# Petals

## What it does

Petals runs a very large language model across a swarm of volunteer GPUs spread over the public internet, BitTorrent-style. The model is split into consecutive transformer blocks; each volunteer server hosts a slice of the layers, and a client stitches together a pipeline that covers the whole model for a forward pass. A distributed hash table maps blocks to peers, and clients cache intermediate activations so they can reroute around servers that drop out. Because activations cross the network between layers, clients can also inject adapters for lightweight fine-tuning without touching the hosted weights. It came out of the BigScience workshop (Yandex Research, HSE, University of Washington, Hugging Face and others), with an ACL 2023 demo paper and a NeurIPS 2023 paper.

## Verdict

Stalled. A real research success and an influential proof of concept that never reached sustained adoption.

## Why

The engineering worked. Petals demonstrated that 50B-to-176B-plus models can run interactively over consumer, geo-distributed hardware with fault tolerance, validated across two continents. But development wound down after the last feature release in September 2023, with commit activity trailing off through mid-2024, and the public swarm has since decayed. On the live health monitor the advertised flagship models (Llama-3.1-405B, Mixtral-8x22B) now show as broken or served by a single partial contributor, so the headline models cannot be run end to end on the public network anymore. The repo keeps roughly ten thousand stars against near-zero development, the classic stalled-but-visible signature.

Three forces killed the momentum. There was no incentive layer, so purely altruistic GPU-sharing eroded once the novelty faded, which is the flaw later work cites most. The physics did not cooperate: pipeline-parallel inference over the internet is latency-bound, generation requires a sequential round trip through the whole server chain, and single-digit tokens per second cannot compete with a centralized GPU or even a good local quantized model. And demand collapsed from both sides as quantized local models (llama.cpp, MLX) let people run capable 7B-to-70B models on one machine while cheap hosted APIs covered the frontier, squeezing the exact niche Petals occupied.

## How it compares to DII

Petals is the honest ancestor and the clearest evidence for the constraints DII designs around. Its Objection-1 lesson (the memory-bandwidth and network wall) is why DII targets latency-tolerant, node-local workloads and a reliable floor of 14B-to-30B models that fit on one node, rather than splitting a frontier model mid-inference across the internet.

What DII offers that Petals does not: a reason for the supply side to persist. Petals had no answer to the question of why anyone keeps a GPU online after their own curiosity is satisfied; DII takes Tor's answer, mission and community rather than payment, and makes benefit never gated on contribution so the network serves people whether or not they run a node. DII also adds a trust model (federated, semi-trusted pods with identity and reputation) where Petals exposed activations to anonymous strangers, and it makes honest degradation, queuing or reporting truthfully below the useful line, a first-class guarantee rather than the silent decay Petals now shows. Where Petals raised no bar on contribution, DII sets node-entry at the useful line, trading node quantity for floor quality.

## Borrow / Avoid

Borrow: fault-tolerant rerouting via client-cached activations and DHT discovery, latency-aware shortest-path routing, a near drop-in `from_pretrained`/`generate` client API, adapter injection between layers, and a public health monitor for swarm transparency. Optimize for interactivity, not batch throughput.

Avoid: shipping without any participation model (the number-one lesson), leaning on pure altruism and a single research team with no operating entity, naive scheduling with no continuous batching or paged attention, marketing cross-internet pipeline parallelism as competitive on tokens per second, and letting privacy leakage and no-custom-code restrictions push every serious user into private swarms that starve the public one. Track live swarm health, not stars.

## Sources

- [Petals (GitHub, bigscience-workshop)](https://github.com/bigscience-workshop/petals)
- [Petals releases (v2.2.0, Sep 2023 latest)](https://github.com/bigscience-workshop/petals/releases)
- [Live swarm health monitor](https://health.petals.dev/)
- [Petals, Yandex Research blog](https://research.yandex.com/blog/petals-decentralized-inference-and-finetuning-of-large-language-models)
- [Petals: collaborative inference and fine-tuning of large models (NeurIPS 2023, arXiv 2312.08361)](https://arxiv.org/abs/2312.08361)
- [Petals ACL 2023 demo (arXiv 2209.01188)](https://arxiv.org/abs/2209.01188)
- [An explorative study of distributed LLM inference (arXiv 2510.11211)](https://arxiv.org/html/2510.11211v1)

Note: peak-swarm figures sometimes cited (hundreds of nodes, thousands of clients) could not be confirmed against a primary source and are omitted.
