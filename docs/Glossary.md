# Glossary

**Reliable Floor**: the project's goal, a dependable baseline of capable AI that is always available, distinct from intermittent access to the most advanced ("ceiling") models. Concretely, the 14B-to-30B class of open-weight models at Q4, delivered as a capability bundle: node-entry at ~14B on a gaming GPU, promise at ~30B on committed hardware, degrading to honest queuing below the useful line rather than to a toy. Defined in ADR-0006 and docs/Reliable_Floor_Definition.md.

**Useful Line**: the capability cliff, around the 14B class, below which a model can chat but cannot reliably reason across steps, use tools, handle longer context, or code. The floor is set by usefulness, not by what technically loads; the network degrades toward this line and then queues honestly rather than serving sub-useful output.

**Local-First**: every node is fully functional offline; the network only adds to what a node can already do on its own.

**Node**: a participant's machine running one or more open-weight models, exposing them as capabilities.

**Node Runtime (the atom)**: the daemon on a node that serves local models, publishes a capability manifest, and routes requests. Self-sufficient offline; the smallest thing that delivers the core promise to a single user.

**Consumer (Consumer Client)**: a participant with no serving hardware who only borrows, such as a phone or browser, sponsored by a pod that agrees to serve it. Architecturally a request that enters the overflow path at the pod stage, having no local stage of its own. See ADR-0002 and ADR-0004.

**Ingress**: how a request enters the router. A local, trusted ingress serves the node's own user and starts at the local stage; a remote, authenticated ingress serves a sponsored consumer and starts at the pod stage. One router, two ingress types.

**Capability Manifest**: what a node publishes about itself: which models it serves, context length, throughput, modalities, and current load. Lets a heterogeneous fleet look uniform to a caller.

**Pod**: a small group of nodes that already trust each other, such as personal machines, a lab, an org, or a community. The core unit of the system.

**Federation**: trusted pods sharing overflow capacity under explicit policy, with no central hub. Closer to the Fediverse than to a single cloud.

**Capability Routing**: requesting intelligence as a capability (reasoning, coding, vision, retrieval) at a required tier; the router then selects the hardware and model.

**Overflow Path**: the order work spills outward, local then pod then federation, only when needed and only within user policy.

**Mission-Driven Participation**: the reason people run nodes, on the Tor model. You contribute because you want a reliable floor of intelligence to exist for your community and for those who cannot self-host, not to earn anything. Access is never gated on contribution (ADR-0002, ADR-0005).

**Reciprocity Signal**: an optional, non-gating priority hint that may let a contributor be served first under congestion. If built at all it is illiquid: non-transferable, non-tradable, expiring, never a currency. A Phase-2 fairness detail, not the reason to participate, and not the earlier "proof-of-contribution credits," which are retired.

**Sovereignty Constraint**: a first-class routing input, for example "data stays on my machines" or "data stays in-country."
