# Reliable Floor Definition

Status: Living reference, v0.1 (2026-07-04). The decision behind this is ADR-0006; this doc holds the concrete, rolling detail that will change as models do.

## The definition, in one line

The reliable floor is the 14B-to-30B class of open-weight models at Q4_K_M quantization, run on hardware people own, delivered as a capability bundle, with a hard useful-line below which the network degrades honestly rather than serving a toy.

## Why defined this way

Models age out in months. Rather than name a model, the floor is defined by three durable things: the capability class, the reference hardware tier, and the workloads it must serve. The specific models below are current examples that meet it, not the definition itself.

Two constraints set the shape. The model must fit on one node, because split inference over a residential network is not viable (The Case Against DII, Objection 1). And the floor is set by a capability cliff, not by what technically loads: below roughly 14B, a model can chat but cannot reliably reason across steps, use tools, handle longer context, or do real coding. Usefulness, not loadability, is the line.

## The three tiers

Node-entry floor. The ~14B class on a 12 to 16GB gaming GPU. The minimum to run a node worth having, set at the useful line on purpose. Below this a participant is treated as a consumer, not a serving node.

Promise floor. The ~30B class on committed hardware: a 24 to 32GB GPU, or a 64GB unified-memory machine such as the reference EVO-X2. The headline capability a pod offers, including to consumers who cannot self-host. 30B is the realistic ceiling on consumer-class hardware; 70B needs prosumer hardware and runs only at latency-tolerant speeds.

Guaranteed floor. The useful line itself, ~14B. When capable nodes are offline or saturated, the network degrades toward 14B, then queues or honestly reports "no capable node available." It never serves a sub-useful model as if it were the floor. This is "graceful degradation, never to zero" read correctly: never to zero usefulness, not never to zero tokens.

## The bundle

The floor is a set of capabilities, not one model, because the router requests capabilities rather than models:

- General reasoning and chat: a model in the 14B-to-30B band.
- Coding: a dedicated code model.
- Embedding and retrieval: a small embedding model.

## Reference hardware

- Base gaming rig, 12 to 16GB VRAM: runs the ~14B class comfortably at Q4, up to ~20B tight. The node-entry tier.
- Committed enthusiast, 24 to 32GB GPU or 64GB unified memory: runs the ~30B class. Reference machine: EVO-X2, 64GB unified, runs Qwen3 30B and a Qwen-Coder-class model comfortably. The promise tier.
- Institutional or prosumer, 48 to 128GB+: runs 70B+ and can serve many consumers. Above the consumer floor; useful but not assumed.

## Current model exemplars (mid-2026, will change)

- General, promise tier: Qwen3 30B, Gemma 4 31B, Mistral Small 4.
- General, entry tier: a 14B-class model, for example a Qwen3 14B.
- Coding: a Qwen-Coder-class 32B, or DeepSeek-Coder-class.
- Embedding: a small open embedding model.

Quantization default is Q4_K_M, roughly half the memory of FP16 for about a 3 to 5% quality cost. A 30B at Q4 runs around 40 to 55 tokens per second on a 32GB consumer GPU, and a 70B at Q4 (~40GB) runs 8 to 15 tokens per second on a 64GB Apple Silicon machine, which is why 70B is treated as latency-tolerant only.

Worth watching: mixture-of-experts models with a small active-parameter count, such as a 35B-total, 3B-active model. Decode latency tracks active rather than total parameters, so these can deliver more capability per unit of latency on ownable hardware, bending Objection 1 slightly in the project's favor. A candidate for the bundle once evaluated.

## What to revisit here

This doc is expected to change. Update the exemplars as better open models ship, revise the token-per-second figures as hardware improves, and record any decision to prefer MoE models in the bundle. The tiers and the useful-line principle come from ADR-0006 and should not drift without a new decision record.

## Sources

- [Best Open Source LLMs in 2026 — AceCloud](https://acecloud.ai/blog/best-open-source-llms/)
- [Local LLM Hardware Requirements 2026 — PromptQuorum](https://www.promptquorum.com/local-llms/local-llm-hardware-guide-2026)
- [Local LLM Deployment on 24GB GPUs — IntuitionLabs](https://intuitionlabs.ai/articles/local-llm-deployment-24gb-gpu-optimization)
- [Best Mac for Local AI 2026 (M1–M5) — Local AI Master](https://localaimaster.com/blog/apple-silicon-ai-buying-guide)
- [Apple Silicon LLM Benchmarks 2026 — llmcheck.net](https://llmcheck.net/benchmarks)
