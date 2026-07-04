# DII Architecture Overview: Local-First, Network-as-Overflow

Status: RFC v0.1 (2026-07-02), updated 2026-07-03 with the router ingress model and prototype scope from ADR-0003 and ADR-0004. Depends on Why DII Exists and The Case Against DII.

## What changed

DII is now scoped as a free, open-source RFC project with no commercial model behind it. Its goal is a dependable baseline of open-model intelligence that anyone can run and that stays beyond the reach of any single off-switch. Monetization is out of scope for the foreseeable future. If a sustainable funding path appears later it can be reviewed on its own merits, and no part of this architecture assumes or requires one. This reframing removes a whole problem class, since there is no token, tradable credit, marketplace to clear, or compute to price. The old slogan "the scheduler is the product" now means only that the router and scheduler are the core engineering contribution.

## The idea in one sentence

Participants who can self-host run a self-sufficient local node that serves capable open-weight AI offline. Nodes optionally federate into trusted pods. People who cannot self-host join as consumers and reach the same capable models through a pod that chooses to serve them. A capability router places each request along the shortest viable path, local first where possible, then pod, then federation, within policy the user controls. For a node-runner the network is an enhancement, and removing it leaves the node still working; for a consumer the network is the delivery itself. Either way the guarantee is the same: access to a capable model that stays beyond the reach of any single off-switch. That guarantee sits at the top of every design choice below (see ADR-0002).

## Design principles

These are the commitments the rest of the design has to satisfy:

- Local-first, always. A node with zero connectivity still delivers useful AI, and the network only amplifies a node that already works on its own.
- Access before contribution. Serving people who cannot self-host, and who may never contribute, is a first-class case rather than an afterthought (ADR-0002).
- Graceful degradation, never to zero. Network down falls to local, pod down falls to the federation, and a peer failing mid-request reroutes. The floor is always "still works locally."
- No single point of control. No global registry, no mandatory hub, and nothing a single directive could switch off. Discovery and identity are federated the way email or Matrix are.
- Open weights only. The guarantee holds only for models users may legally hold and run, which also keeps DII clear of the frontier-proliferation objection.
- Trust before proof. Prefer running work inside a trust boundary over proving correctness on strangers' hardware, keeping heavy verification as a last resort.
- Sovereignty as a first-class routing input. Constraints such as "never leaves my machines" or "stays in-country" are honored from the start.

## Topology: a federation of trusted pods

This is the load-bearing decision. The instinct behind the original vision, that millions of strangers' machines become one computer, is the version the red-team teardown hit hardest on trust and sustainability. This architecture makes a different bet. The unit is the pod, a small group of nodes that already trust each other: your own machines, a lab, a hackerspace, a company, a family, or a town's library co-op. Pods are self-sufficient, and they may federate with other pods to share overflow capacity under explicit policy.

The shape earns its keep in three ways. It defuses the trust tax, because most work runs inside a boundary where the nodes are already trusted, so expensive verification stays unnecessary for the common case. It makes participation durable, because reciprocity within a real community is a stronger and more human incentive than altruism toward anonymous strangers, which is the failure mode that plateaued volunteer computing. And it has no chokepoint, because a federation of independent pods has no switch to throw; kill any pod and the rest keep running. The mental model is closer to the Fediverse for inference than to SETI@home for LLMs.

## Components

The node runtime is the atom. It runs one or more open-weight models locally, exposes them as capabilities, and functions fully offline. This is the star of the whole system, since a person who only ever runs a single node has already received the core promise of DII, and everything else is optional scaling on top of it.

The consumer client is the counterpart for people who cannot self-host. It is a lightweight role, a phone app or a browser, that requests capabilities without hosting a model of its own. A consumer owns no serving hardware and is sponsored by a pod that agrees to serve it, which is how the reliable floor reaches people who have no capable device. Architecturally the consumer is not a separate system but a request that enters the overflow path past the local hop: a node-runner's request begins at the local stage, while a consumer, hosting no model, has no local stage and enters the router at the pod stage. Because a node exposes an OpenAI-compatible endpoint, any existing OpenAI-compatible application pointed at a serving node is already a consumer client, so the role needs no bespoke software. This is recorded in ADR-0004.

The capability abstraction lets requests target capabilities such as reasoning, coding, vision, speech, retrieval, or embedding, at a required quality and latency tier. The caller names the capability and tier, and the router selects the hardware and model. Each node publishes a capability manifest describing which models it serves, its context length, throughput, modalities, and current load, which is what lets a heterogeneous fleet look uniform to a caller.

Discovery is federated. Within a pod it can be a simple local registry or LAN discovery. Across pods it uses signed, gossiped peer lists and a DNS-like federated directory, deliberately avoiding a single global index that would reintroduce the chokepoint the project exists to eliminate.

The router and scheduler are the coordination core. They decide where a request executes and enforce the overflow order of local, then pod, then federation. Routing weighs capability match, sovereignty and policy constraints, the trust level of candidate nodes, latency, availability, and load. Each hop outward passes a policy gate, so work only spills across a boundary the user has authorized.

The router has two ingress types feeding the same logic. A local, trusted ingress serves the node's own user and starts each request at the local stage. A remote, authenticated ingress serves a sponsored consumer and starts the request at the pod stage, since a consumer has no local model to try first. There is one router, not two, and the consumer path is the ordinary overflow path entered one step in (ADR-0004). The remote ingress is authenticated because a consumer takes without giving and is a higher abuse surface than a trusted peer; the identity mechanism behind that authentication is still an open decision, stubbed in the prototype.

Trust and reputation carry the weight that payment would carry elsewhere. With no money in the system, trust rests on identity through signed and persistent node identities, reputation through track record within and across pods, and policy covering who you accept work from, who you send work to, and where data may travel. Verification is tiered: negligible inside a trusted pod, optional for cross-pod work, and for the rare high-stakes cross-boundary case, borrowing an existing protocol in the TOPLOC or Verde family rather than inventing one. Work inside a trusted pod skips the redundant-execution tax.

The participation model is reciprocity. Contribution earns access: run a node and contribute idle capacity, and you gain priority when you later need to spill work onto the federation. This access is deliberately non-transferable, non-tradable, and expiring, which makes it a fair-use signal rather than a currency and keeps DII clear of the crypto category. Pure voluntary contribution is welcome too; the reciprocity layer just makes sustained participation rational. Because access does not require contribution (ADR-0002), a pod can also hold donated or sponsored capacity, a shared pool that contributors and institutions feed so consumers have something to draw on. Since contribution no longer guards that pool, fair-use quotas and a way to vouch for real people keep a single actor from posing as many consumers and draining it.

Privacy-preserving execution follows from sovereignty being the reason to exist. The default keeps data inside the user's trust boundary, and cross-boundary execution is opt-in and policy-gated. Techniques for protecting data on semi-trusted hops, and whether to allow them at all, are an open design area flagged here for later.

Fault tolerance treats failure as normal. Peers vanish, networks partition, and pods go dark, and every failure mode resolves downward toward local execution.

## The overflow path

```
Request for a capability
   │
   ▼
[1] Local node can serve it?  ──yes──►  run locally (always preferred)
   │ no / over capacity
   ▼
[2] A trusted pod peer can?   ──yes──►  run in-pod (policy gate)
   │ no
   ▼
[3] A federated pod can?      ──yes──►  run cross-pod (policy + trust gate)
   │ no
   ▼
Degrade gracefully: queue, use a smaller local model, or report honestly.
```

The final step always returns the best the local node can do, which is the resilience guarantee made concrete. A node-runner's request enters at [1]. A sponsored consumer, having no local model, enters at [2] instead; the rest of the path is identical.

## Build order

- Phase 0, the atom. A solid local-first node runtime with a capable open model, offline capability, and a capability manifest. This alone delivers the core promise to a single user and is worth shipping on its own.
- Phase 1, two nodes and a consumer in one pod. The Week-3 proof of concept, written in Go (ADR-0003), where node A routes a capability request to node B and back, and a consumer C with no local model borrows from the pod through the remote ingress (ADR-0004). Each node exposes an OpenAI-compatible endpoint so existing clients work unchanged. Measure the numbers the teardown demands, such as tokens per second over a residential link versus local, time-to-first-token, and overflow overhead.
- Phase 2, pod overflow and policy. The router honors sovereignty and trust policy, reciprocity accounting comes online, and graceful degradation paths are exercised.
- Phase 3, federation. Signed discovery across pods, cross-pod trust with optional verification, and the "Fediverse for inference" topology.

## Open decisions to record as ADRs

- Reciprocity versus pure-voluntary as the default participation model. Current lean is reciprocity, illiquid, as described under the participation model above.
- Discovery protocol: gossip, DNS-like federation, or Matrix-style servers.
- Cross-boundary verification: when, if ever, redundant execution or a borrowed proof is worth its cost.
- Data protection on untrusted hops: whether to allow it at all, and if so how.
- Node identity: self-sovereign keys versus pod-issued identity.
- Serving consumers: how pods decide whom to serve and at what fair-use limit, how donated capacity is accounted for without becoming a currency, and how vouching resists abuse without becoming an exclusionary gate.

## What this deliberately excludes

DII is not a global marketplace or a swarm of anonymous strangers, carries no token, credit-currency, or compute-for-pay scheme, does not serve frontier dual-use capability since it targets a dependable floor, and depends on no hub, registry, or provider that could be switched off.
