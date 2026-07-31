# DII Architecture Overview: Local-First, Network-as-Overflow

Status: RFC v0.1 (2026-07-02), updated 2026-07-03 with the router ingress model and prototype scope from ADR-0003 and ADR-0004, 2026-07-27 with the access-distribution reframe (ADR-0015), and 2026-07-30 with the functional-node M1 build. Depends on Why DII Exists and The Case Against DII.

## What changed

DII is now scoped as a free, open-source RFC project with no commercial model behind it. Its goal is a dependable baseline of open-model intelligence that anyone can run and that stays beyond the reach of any single off-switch. It is a mission-driven volunteer effort on the model of Tor, non-commercial by commitment: monetization is retired, not paused, and the network itself assumes no funding path, since there is nothing to price (ADR-0005). Sustaining the nonprofit steward and offsetting operators' bare-minimum costs to break-even does assume a funding path, held strictly distinct from monetizing the network (ADR-0008). This reframing removes a whole problem class, since there is no token, tradable credit, marketplace to clear, or compute to price. The old slogan "the scheduler is the product" now means only that the router and scheduler are the core engineering contribution.

## What changed, 2026-07-27 (ADR-0015)

The thesis is access distribution, not capability aggregation. DII exists to put a capable floor model (ADR-0006) into the hands of people who lack access, under conditions they can accept: no counterparty who can switch them off, no gating beyond the hard legal lines, and no data leaving a boundary they trust. That is the load-bearing reason the network exists, alongside the mission motivation (ADR-0005). Model diversity, a pod reaching many models, is a feature and not the reason to exist. The Week-4 variety experiment came back inconclusive-to-null on that question and it is settled enough to act on (docs/Variety_Experiment_Findings.md, ADR-0015).

This reframes the "Network-as-Overflow" emphasis in the title above. The load-bearing primitive is the consumer ingress (ADR-0004): a person with no node reaching a willing node that serves them. Node-to-node aggregation, the overflow path where one node borrows a model a peer holds, is reclassified from core mechanism to optional optimization, validated and available (ADR-0011) but not something the network's guarantee depends on. A node serving a single model within its capability is a full participant. What the coordination layer still owes is availability, not variety: enough nodes serving the floor, and a directory that finds an available one (ADR-0009). The next load-bearing decision is therefore consumer admission and abuse resistance, who can use them, promoted from a loose end to the center (docs/Identity_Note_From_Prototype.md, docs/Governance_And_Abuse_Resistance.md).

The sections below still describe overflow as the spine; read them through this reframe until they are rewritten. The engineering is unchanged and validated. What changed is which part is the reason to exist.

## What changed, 2026-07-30 (functional node M1)

The node runtime stopped being a demonstration. The Week-3 prototype proved the loop but, by its own non-goals, skipped everything that makes a node dependable rather than illustrative. M1 of the functional-node plan (docs/Functional_Node_Plan.md) closes the first part of that gap, and three of its results matter at this level rather than only in the code.

The reliable-floor bundle is now actually served, not just specified. ADR-0006 and the Reliable Floor Definition describe the floor as a bundle of a general model, a coder, and an embedding model; until M1 the node served only the chat shape, so the third of the bundle existed on paper. It is now a real endpoint on the same backend seam. A useful detail fell out of building it: the standard model-list call the capability manifest is built from reports every model a backend serves and carries no type field, so an embedding model is simply another name in the same list and the single model-name capability match covers both. Capability tags remain a later layer, and this is one more piece of evidence that they are not yet needed.

Degrading honestly now includes degrading *visibly*. "Graceful degradation, never to zero" was stated as a routing property — fall back to local, report honestly when nobody can serve. M1 adds the operational half: a node reports liveness and readiness separately, so a node that is up but whose model server has died is detectable rather than looking healthy to every supervisor and directory that asks. That distinction is what a directory promising availability (ADR-0009) will eventually need to consume, since "a node is reachable" and "a node can serve" are different facts.

Fair-use has a substrate rather than a plan. The participation model notes that because access does not require contribution, fair-use quotas are what keep a single actor from posing as many consumers and draining a pod's donated capacity. Quotas need accounting to enforce, and accounting is now built: one structured record per request, carrying which node served it and what it cost, keyed on an opaque handle and never on prompt text. That last property is not incidental — it is the data-minimization seam ADR-0016 describes, and building the counters this way now is what keeps a later fair-use mechanism from quietly becoming a behavioral reputation store.

What M1 deliberately does not do: per-consumer identity, fair-use enforcement, and moderation are M2 through M4, and serving the unaffiliated public stays gated on the ADR-0016 legal review regardless of what the code can do.

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

The shape earns its keep in three ways. It defuses the trust tax, because most work runs inside a boundary where the nodes are already trusted, so expensive verification stays unnecessary for the common case. It makes participation durable, because belonging to a real community is a stronger and more human motivation than altruism toward anonymous strangers, which is the failure mode that plateaued volunteer computing, and it is what lets a mission carry the cost of contributing a GPU. And it has no chokepoint, because a federation of independent pods has no switch to throw; kill any pod and the rest keep running. The mental model is closer to the Fediverse for inference than to SETI@home for LLMs.

Pods come in kinds, and the distinction carries the access and funding model (ADR-0009). A public-serving pod chooses to serve unaffiliated strangers as an act of mission, the analog of a Tor exit relay, and is the funding priority; a private pod serves only its own trusted members and self-funds, because it already gets its return as mutual benefit; a public-interest pod may be closed yet do public-good work and can be funded case by case. A person with no pod of their own is matched across many public-serving pods by a decentralized directory, never to one central public pod, which is how "AI for All" is delivered without rebuilding a chokepoint. Rivalrous cost helps here rather than hurts, because sponsoring strangers consumes real GPU time, so no pod can grow into an everyone-pod without collapsing, which caps concentration by construction. One principle keeps the whole thing honest: the steward is a legitimate center for governance and funding, but it must never sit in the data path, and the test is that removing it should not stop the network from serving.

## Components

The node runtime is the atom. It runs one or more open-weight models locally, exposes them as capabilities, and functions fully offline. This is the star of the whole system, since a person who only ever runs a single node has already received the core promise of DII, and everything else is optional scaling on top of it.

The consumer client is the counterpart for people who cannot self-host. It is a lightweight role, a phone app or a browser, that requests capabilities without hosting a model of its own. A consumer owns no serving hardware and is sponsored by a pod that agrees to serve it, which is how the reliable floor reaches people who have no capable device. Architecturally the consumer is not a separate system but a request that enters the overflow path past the local hop: a node-runner's request begins at the local stage, while a consumer, hosting no model, has no local stage and enters the router at the pod stage. Because a node exposes an OpenAI-compatible endpoint, any existing OpenAI-compatible application pointed at a serving node is already a consumer client, so the role needs no bespoke software. This is recorded in ADR-0004.

The capability abstraction lets requests target capabilities such as reasoning, coding, vision, speech, retrieval, or embedding, at a required quality and latency tier. The caller names the capability and tier, and the router selects the hardware and model. Each node publishes a capability manifest describing which models it serves, its context length, throughput, modalities, and current load, which is what lets a heterogeneous fleet look uniform to a caller. The concrete quality target is the reliable floor: the 14B-to-30B open-weight class at Q4, delivered as a bundle of a general model, a coder, and an embedding model, with node-entry at ~14B and the promise at ~30B (ADR-0006 and the Reliable Floor Definition). Below the ~14B useful line the router degrades honestly, which is what the final step of the overflow path means.

Discovery is federated. Within a pod it can be a simple local registry or LAN discovery. Across pods it uses signed, gossiped peer lists and a DNS-like federated directory, deliberately avoiding a single global index that would reintroduce the chokepoint the project exists to eliminate.

The router and scheduler are the coordination core. They decide where a request executes and enforce the overflow order of local, then pod, then federation. Routing weighs capability match, sovereignty and policy constraints, the trust level of candidate nodes, latency, availability, and load. Each hop outward passes a policy gate, so work only spills across a boundary the user has authorized.

The router has two ingress types feeding the same logic. A local, trusted ingress serves the node's own user and starts each request at the local stage. A remote, authenticated ingress serves a sponsored consumer and starts the request at the pod stage, since a consumer has no local model to try first. There is one router, not two, and the consumer path is the ordinary overflow path entered one step in (ADR-0004). The remote ingress is authenticated because a consumer takes without giving and is a higher abuse surface than a trusted peer; the identity mechanism behind that authentication is still an open decision, stubbed in the prototype.

Trust and reputation carry the weight that payment would carry elsewhere. With no money in the system, trust rests on identity through signed and persistent node identities, reputation through track record within and across pods, and policy covering who you accept work from, who you send work to, and where data may travel. Verification is tiered: negligible inside a trusted pod, optional for cross-pod work, and for the rare high-stakes cross-boundary case, borrowing an existing protocol in the TOPLOC or Verde family rather than inventing one. Work inside a trusted pod skips the redundant-execution tax.

The participation model is mission, not reciprocity (ADR-0005). People run nodes for the same reason people run Tor relays: because they want a reliable floor of intelligence to exist for their community and for those who cannot self-host, not because they earn anything. Because ADR-0002 already guarantees access without a contribution gate, reciprocity can never be the gate and never the engine. At most it is an optional, non-gating priority signal that decides who is served first when a federation is congested, and if it is built at all it must be illiquid: non-transferable, non-tradable, expiring, never a currency. Whether to build it is a Phase-2 question and may prove unnecessary, since within a small trusted pod social accountability does the fairness work a ledger otherwise would. Because access does not require contribution, a pod can also hold donated or sponsored capacity, a shared pool that contributors and institutions feed so consumers have something to draw on. Since contribution no longer guards that pool, fair-use quotas and a way to vouch for real people keep a single actor from posing as many consumers and draining it.

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

- Phase 0, the atom. A solid local-first node runtime with a capable open model, offline capability, and a capability manifest. This alone delivers the core promise to a single user and is worth shipping on its own. Underway as the functional node (docs/Functional_Node_Plan.md): M1 completed the floor bundle including embeddings, graceful lifecycle, separate liveness and readiness, per-request usage accounting, and a container image, which is what turns the atom into something a person can actually depend on daily. Per-consumer identity, fair-use, and the moderation seam are M2 through M4.
- Phase 1, two nodes and a consumer in one pod. The Week-3 proof of concept, written in Go (ADR-0003), where node A routes a capability request to node B and back, and a consumer C with no local model borrows from the pod through the remote ingress (ADR-0004). Each node exposes an OpenAI-compatible endpoint so existing clients work unchanged. Built and validated (prototype/, extended to a live three-node pod): all four kill-criteria passed, with overflow throughput about 100 percent of the peer's own local and only 20 to 43 milliseconds of added time-to-first-token, so borrowing a peer's capability is effectively free, and the overhead is transport-bound and model-independent so it carries to the 30B floor (journal/2026-07-12-week3-m4-findings.md). The inter-node transport is the reused OpenAI-compatible HTTP call (ADR-0011). What remains unsettled is durability, not feasibility (docs/Pod_Aggregation_Red_Team.md).
- Phase 2, pod overflow and policy. The router honors sovereignty and trust policy, reciprocity accounting comes online, and graceful degradation paths are exercised.
- Phase 3, federation. Signed discovery across pods, cross-pod trust with optional verification, and the "Fediverse for inference" topology.

## Open decisions to record as ADRs

- Whether to build a reciprocity fairness knob at all, and its exact illiquid form if so. The principle is settled (mission-primary, non-gating, never a currency, per ADR-0005); only the Phase-2 mechanism is open.
- Discovery protocol: gossip, DNS-like federation, or Matrix-style servers.
- Cross-boundary verification: when, if ever, redundant execution or a borrowed proof is worth its cost.
- Data protection on untrusted hops: whether to allow it at all, and if so how.
- Node identity: self-sovereign keys versus pod-issued identity. The Week-3 prototype stubbed this and produced a requirements note (docs/Identity_Note_From_Prototype.md) surfacing three concrete needs the ingress exposed, node admission, per-consumer credentials, and caller attribution across hops, to drive the identity ADR from observed behavior rather than a paper comparison.
- Serving consumers: how pods decide whom to serve and at what fair-use limit, how donated capacity is accounted for without becoming a currency, and how vouching resists abuse without becoming an exclusionary gate.

## What this deliberately excludes

DII is not a global marketplace or a swarm of anonymous strangers, carries no token, credit-currency, or compute-for-pay scheme, does not serve frontier dual-use capability since it targets a dependable floor, and depends on no hub, registry, or provider that could be switched off.
