# ADR-0006: What the Reliable Floor Actually Is

Status: Accepted, 2026-07-04. Makes concrete the "reliable floor" named in the Vision, Charter, and Why DII Exists. Depends on ADR-0002 and ADR-0004.

## Context

The project's guarantee is a "reliable floor" of capable AI, but until now the floor was defined only qualitatively, as "an advanced-enough model you can rely on." That phrasing is fine for the pitch and useless for engineering. The Week-3 prototype needs a specific model to run and measure, the kill criteria need a specific target to compare against centralized serving, and the consumer promise needs a concrete quality level. This record pins the floor.

Two facts shape it. First, in the surviving envelope the model must fit on one node, because split inference across a residential network is not viable (The Case Against DII, Objection 1). So the floor is bounded by single-machine memory, which makes "hardware people own" the real variable. Second, and more important, the floor is defined by a capability cliff, not by whatever technically loads. Below roughly the 14B class, an open model can hold a conversation but cannot reliably do multi-step reasoning, tool use, longer-context work, or real coding. A model that hands someone a toy is not a floor. Usefulness, not loadability, sets the line.

The hardware reality, as of mid-2026 and confirmed by the project owner's own machine, an EVO-X2 with 64GB unified memory that runs the 30B class comfortably: a 12 to 16GB gaming GPU runs the 14B class well, a 24 to 32GB GPU or a 64GB unified-memory machine runs the 30B class, and 70B needs prosumer hardware and runs only at latency-tolerant speeds. 30B is therefore the realistic ceiling on consumer-class hardware.

## Decision

The reliable floor is the 14B-to-30B class of open-weight models at Q4_K_M quantization, delivered as a capability bundle, with a hard useful-line below which the network degrades honestly rather than serving sub-useful output.

Three numbers make it concrete.

The node-entry floor is the ~14B class on a 12 to 16GB gaming GPU. This is the minimum to run a node worth having, and it is deliberately set at the useful line rather than lower. A node that can only serve sub-14B work adds little, so it is treated as a consumer rather than a serving node.

The promise floor is the ~30B class on committed hardware: a 24 to 32GB GPU or a 64GB unified-memory machine such as the reference EVO-X2. This is the headline capability a pod offers, including to consumers who cannot self-host.

The guaranteed floor is the useful line itself, ~14B. When capable nodes are offline or saturated, the network degrades toward 14B, and below that it queues or honestly reports that no capable node is available. It never serves a sub-useful model as if it were the floor. This is the correct reading of "graceful degradation, never to zero": never to zero usefulness, not never to zero tokens.

The floor is a bundle, not a single model, because the router speaks in capabilities: a general reasoning and chat model in the band, a coding model, and a small embedding model. Current mid-2026 exemplars are Qwen3 30B and Gemma 4 31B for general use, a Qwen-Coder-class model for code, and a small embedding model, with a 14B-class model at the entry tier. These names will age out; the durable definition is the class, the hardware tier, and the workloads, not the specific model.

## Consequences

The prototype is now fully concrete. The reference node runs a 30B general model, the kill criteria measure that specific model local versus pod, and the owner's EVO-X2 is node A.

Raising the node-entry bar to the 14B useful line is a deliberate trade of node quantity for floor quality, and it is another place DII diverges from Tor. Any spare Tor bandwidth helps, but sub-useful inference does not, so the project prefers fewer nodes that each genuinely help over many that do not. This is consistent with a project whose entire promise is a floor that can be relied on.

Because 30B-capable nodes are a committed minority and their presence is intermittent, the floor is defined by what a pod can reliably keep available, not by peak capability when everything is online. This raises the importance of availability and redundancy within a pod, and it makes the honest-degradation behavior a first-class part of the guarantee rather than an error path.

## Open questions

Which exact models fill the bundle at any given time is a rolling operational choice, not an architecture decision, and should be tracked in the Reliable Floor Definition doc rather than here. Whether mixture-of-experts models with a small active-parameter count should be preferred in the bundle, since decode latency tracks active rather than total parameters and they bend Objection 1 slightly in our favor, is worth a dedicated look. How a pod measures and advertises the floor it can currently guarantee, given intermittent nodes, connects to the capability-manifest and discovery work still to be specified.
