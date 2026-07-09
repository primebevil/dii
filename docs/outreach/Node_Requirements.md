# Node Requirements: What Counts as "Enough to Run the Floor"

Status: v0.1 (2026-07-04). Practical hardware guidance for hosts, derived from the Reliable Floor Definition and the Architecture Overview (ADR-0006). The numbers below are engineering estimates for outreach conversations. They should be validated against the Phase-1 proof-of-concept benchmarks (tokens per second over a residential link, time-to-first-token, overflow overhead) and updated once real measurements exist.

## The floor, in one paragraph

The reliable floor is the 14B-to-30B open-weight model class, quantized to roughly 4-bit (Q4), delivered as a small bundle: a general model, a coder, and an embedding model. Entry to the network is around 14B, and the promise the network tries to hold is around 30B. Below roughly 14B the router degrades honestly rather than pretending, so sub-14B hardware can still contribute light work but does not by itself deliver the floor. This is what "enough" means: enough VRAM and memory bandwidth to serve a Q4 model in that class at usable speed.

## Two roles, two specs

A host is usually filling one of two roles, and the hardware bar is different for each.

An edge node runs the entry-class model for its own user and adds spare capacity to a pod. It does not have to be always on, and losing it does not break the pod. This is the friends-and-members machine.

An anchor node is the always-on box that makes a pod real for a whole community. It serves the ~30B promise, stays up, and answers requests from lighter nodes and from consumers who host nothing. When you ask a hackerspace, lab, or library to "host one node," this is the one you mean.

| | Edge node (entry, ~14B Q4) | Anchor node (promise, ~30B Q4) |
|---|---|---|
| GPU VRAM | 12 GB usable | 24 GB usable, or 2 cards |
| Example GPUs | RTX 3060 12GB, 4060 Ti 16GB, A2000 12GB | RTX 3090 24GB, RTX 4090 24GB, or 2x 16GB |
| Apple Silicon | M-series, 16 GB+ unified memory | M-series Max/Ultra, 32-64 GB unified memory |
| System RAM | 16 GB+ | 32 GB+ |
| Storage | ~30 GB free (model bundle) | ~60 GB+ SSD (bundle plus headroom) |
| Uptime | occasional, best-effort | always-on |
| Network | any home connection | stable link, ideally with decent upload |

These are starting points, not hard cutoffs. A 16 GB card can run a 14B model comfortably and a smaller 30B-class model with tight context; a 24 GB card runs the 30B promise with room for a few concurrent requests. More VRAM mostly buys you longer context and more simultaneous users, which is exactly what an anchor serving a pod needs.

## Apple Silicon is a strong quiet anchor

Worth calling out because libraries and hackerspaces often have Macs. Apple's unified memory lets an M-series machine run these models through llama.cpp or MLX, and a Mac Studio or a Max/Ultra laptop with 32 to 64 GB makes an excellent always-on anchor: low power draw, near-silent, and easy to leave running in a shared space. If a host says "we don't have a gaming GPU but we have a Mac," that is often good news, not a fallback.

## What sub-floor hardware can still do

A machine with an 8 GB GPU, or a recent CPU with no capable GPU, sits below the entry line for the full floor. It can still run a smaller model for light tasks and act as an edge node that the router uses when nothing better is free, and it can always serve as a consumer client that borrows capability from the pod. So no one is turned away for having modest hardware; they just are not the anchor. This matters in outreach because it lets you say yes to almost anyone while still being clear about what a pod actually needs to stand up.

## Power, network, and the honest caveats

An anchor node is one machine drawing on the order of a few hundred watts under load (less for Apple Silicon, more for a dual-GPU rig), plus a network drop. Residential connections work for latency-tolerant workloads, which is the envelope DII targets on purpose, but an anchor that serves remote consumers benefits from a stable connection with reasonable upload. A node that loses the network keeps working locally for its own users, so a host is never left with a dead box.

Two honest caveats to carry into any conversation. First, throughput on a single consumer GPU serving a 30B Q4 model is measured in a handful of tokens per second per request, not the speed of a cloud endpoint; the workloads DII serves are chosen to tolerate that. Second, all of the numbers above are pre-prototype estimates. Once the Phase-1 pod is measured, replace them here with real figures, and prefer citing those to citing these.
