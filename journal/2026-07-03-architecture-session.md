# Architecture Session, 2026-07-03

A working discussion about the path to a Week-3 prototype. No code was written. The output is decisions and their reasoning, captured as ADR-0003 and ADR-0004 and folded into the Architecture Overview and Glossary.

## Starting point

The session began as a "where do I even start building" question, including whether nodes should coordinate through a ledger or blockchain, and whether to write a PRD. Reading the existing corpus reframed all of it: most of these questions were already answered by the project's own documents.

The ledger and blockchain question is settled by the existing architecture, which rules out any token, tradable credit, marketplace, and any switchable hub, and specifies federated discovery through signed, gossiped peer lists in the manner of email or Matrix. A blockchain exists to coordinate distrusting strangers who exchange value; DII has deliberately removed value exchange and replaced strangers with trusted pods, so the mechanism does not apply. Within a pod, discovery is a simple registry or LAN discovery; across pods it is gossip plus a DNS-like federated directory. Reopening blockchain would fight the project's own design.

A PRD is not the missing artifact either. The project already has a Charter. What is missing is one level down: a component spec for the node runtime and a step-by-step workflow for the two-node proof of concept. Those, not a PRD, are the next documents worth writing.

## What the prototype has to be

The purpose chosen for the network is resilient, decentralized inference. Combined with the teardown's verdict, that narrows the surviving envelope to latency-tolerant, node-local work, where each request fits on one machine and nodes do not talk mid-inference. For the prototype this is a simplification rather than a limitation: Phase 0 and Phase 1 need no split inference, no mid-inference peer chatter over a residential link, and no heavy cross-boundary verification. The first thing to build is close to mundane, a local inference server with a capability manifest and an overflow router.

The recruiting demo is not "it routes." It is that a normal AI application, pointed at localhost, works fully offline, and transparently borrows a peer's GPU when the local machine cannot serve a request, degrading to local rather than to zero when the network is gone.

## Decisions taken

Language. Go for the prototype, recorded in ADR-0003. The deciding factor is that the work is AI-driven with the human in the loop, which makes the owner's review bandwidth the scarce resource, and Go is more readable to review from a backend PHP background than Rust. Rust's runtime-speed advantage is invisible here because the bottleneck is inference and the network. Rust's compiler-as-second-reviewer argument is real but bites at the production boundary, where a rebuild is already likely. The choice is low-regret because the node protocol is HTTP and JSON, so a future Rust node can federate with Go nodes.

Consumer as first-class ingress. Recorded in ADR-0004, building on ADR-0002. The two target use cases are a node-runner who both shares and borrows, and a pure consumer with no hardware who only borrows. The consumer is not a separate system; it is a request that enters the overflow path past the local hop. One router, two ingress types: a local trusted ingress that starts at the local stage, and a remote authenticated ingress that starts at the pod stage. Because a node already exposes an OpenAI-compatible endpoint, any existing OpenAI-compatible client is a consumer client. The abuse-surface, identity, and fair-use questions are pulled forward but only stubbed in the prototype; their substance stays open under ADR-0002.

## On sequencing

Working through the prototype felt like jumping to Week 3 before Week 2's research was finished, but the session's output was architecture and decisions with reasons, which is current-phase work. Designing a prototype is a research act; it reveals which questions the research still has to answer. The trap to avoid is letting prototype-thinking substitute for the demand probe, which the teardown names as the objection most likely to kill the project. The highest-leverage next moves are code-free: capture these decisions, which this session did, and run the demand probe to find people who would use this for sovereignty or privacy reasons, which has the longest lead time and can validate or kill the rest before any Go is written.

## Return to Week 1: item 1, the participation model

After parking the prototype work, the session went back to Week 1 to tighten the proposition in order. Reviewing it surfaced that the participation model was inconsistent: the Manifesto and FAQ presented reciprocity ("contribution earns priority") as the reason to take part, while the Overview listed it as an open decision and the teardown's Objection 4 showed reciprocity-as-motivation stalls once a contributor's own demand is met.

The resolution, now ADR-0005, is that ADR-0002 had already settled the principle. Once access is not gated on contribution, reciprocity cannot be the gate or the engine; it can only be an optional, non-gating, illiquid priority signal, and whether to build it is a Phase-2 question. The real motivation is mission and community, on the model of Tor: you run a node because you want a reliable floor of intelligence to exist for your community and for those who cannot self-host, the way people run Tor relays for privacy without being paid, and the way most who benefit from Tor never run a relay.

Adopting the Tor framing also weakens the teardown's Objection 2. That objection assumed a market with thin paying demand; with no price and nothing for sale, the question shifts from "will anyone pay" to "will anyone use," a much lower bar. It also explains why community scope is load-bearing: a Tor relay forwards cheap spare bandwidth, while a DII node runs expensive, rivalrous GPU inference, so the cost is only payable where the benefit is visible and personal.

The owner decided to make "AI for All, in the spirit of Tor" explicit positioning and to scrub the money-making ideas that kept leaking in, since monetization is not a viable route at this time. Monetization is now retired rather than paused. The Tor framing was added to the Vision, Manifesto, README, and FAQ; the funding hedges were removed from the Charter and Overview; participation language was moved from reciprocity-primary to mission-primary across the repository; and the Glossary replaced "Contribution-Earns-Access" with mission-driven participation plus a demoted reciprocity signal.

## Return to Week 1: item 2, the reliable floor

The second Week-1 loose thread was that the floor was defined only qualitatively. Grounding it in current mid-2026 data and the owner's own hardware pinned it concretely. The decision is ADR-0006, with the rolling detail in docs/Reliable_Floor_Definition.md.

The key move was defining the floor by a capability cliff rather than by what technically loads. Below roughly the 14B class a model can chat but cannot reliably reason across steps, use tools, handle longer context, or code, which the owner confirmed from running a 7B. Usefulness sets the line. The floor became a band of three numbers: node-entry at ~14B on a 12 to 16GB gaming GPU, the promise at ~30B on committed hardware such as the owner's EVO-X2 with 64GB unified memory, and a guaranteed floor at the ~14B useful line, below which the network queues or reports honestly rather than serving a toy. 30B is the realistic ceiling on consumer-class hardware; 70B needs prosumer machines and runs latency-tolerant only. The floor is a bundle, a general model plus a coder plus an embedding model, since the router speaks in capabilities.

Two consequences worth keeping. Setting node-entry at the useful line trades node quantity for floor quality, another divergence from Tor, where any spare bandwidth helps but sub-useful inference does not. And because 30B-capable nodes are a committed minority with intermittent presence, the floor is what a pod can reliably keep available, which makes honest degradation a first-class part of the guarantee rather than an error path. This also made the prototype fully concrete: the reference node runs a 30B general model and the kill criteria measure it against centralized serving.

## Follow-ups

- Continue the Week-1 tidy in order: sharpen the target user beyond "everyone" (item 3).
- Write the node-runtime component spec and the two-node-plus-consumer workflow, including the failure and reroute paths.
- Run the demand probe from the teardown's kill criteria, now framed as "will anyone use this," and study the Tor Project as prior art for sustaining volunteer infrastructure.
- Finish the Week-2 write-ups on why prior projects won or stalled.
