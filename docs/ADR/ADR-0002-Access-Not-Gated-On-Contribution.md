# ADR-0002: Access Is Not Gated on Contribution

Status: Accepted, 2026-07-02.

## Context

DII exists to keep capable AI available to everyone when centralized access is restricted or withdrawn. An earlier version of the architecture put local-first execution at the center and let access follow from contribution, through the reciprocity model where running a node earns priority on the federation.

That framing has a blind spot. The people most exposed when a hyperscaler cuts off access are the ones who cannot run a capable model at all, because their only device is a phone or a modest laptop. They used the cloud because it was the one way they could reach a capable model, and when it disappears they have nothing. A design that answers "run it locally" or "contribute to earn access" serves the already-equipped and abandons exactly the people the mission is meant to protect. An exclusive network does not solve the consumer problem.

## Decision

Access to the reliable floor does not require the ability to run a model or the act of contributing compute. The network serves non-contributing consumers as a first-class case.

Local execution remains the foundation for people who can self-host, and contribution through reciprocity remains the default way full nodes earn priority. But the guarantee the project makes is access to a capable model, and local execution is only one way to deliver it. For someone who cannot self-host, the network is the delivery path, through a pod that chooses to serve them.

The organizing model is the pod as community infrastructure rather than a strict tit-for-tat market. A library, a school, a mutual-aid group, or a few people with capable hardware run nodes and offer access to people around them who cannot self-host, the way a public library's computers serve people who do not own one. Contributors underwrite consumers.

## Consequences

This is the harder and truer version of the mission, and it reopens problems a local-only framing had set aside. Serving a consumer means real inference runs on someone else's hardware over a network, which brings the physics of distributed inference and the question of who provides capacity back into scope (see The Case Against DII, Objections 1 and 2). It also breaks a clean reciprocity story, because consumers take without giving, which is the free-rider dynamic that has limited volunteer computing.

Three capabilities become required rather than optional. There has to be a lightweight consumer or client role, distinct from a full node, so a phone or browser can request work without hosting a model. There has to be a notion of donated or sponsored capacity, a shared pool that contributors and institutions feed for consumers to draw on. And because contribution is no longer the gate that protects that pool, something else has to prevent its exhaustion: fair-use quotas together with a way to vouch for real people, so a single actor cannot pose as many consumers and drain it.

This raises the stakes on the routing, capacity-sharing, and fairness layer, which is now the core of the engineering rather than the local runtime. It also strengthens the project's answer to "why not just run a local model yourself," because serving people who cannot do that is precisely what a bare local runtime does not offer.

## Open questions

How pods decide whom to serve and at what fair-use limit. How donated capacity is accounted for without turning into a currency. How vouching resists abuse without becoming an exclusionary gate of its own. These are left for later decision records.
