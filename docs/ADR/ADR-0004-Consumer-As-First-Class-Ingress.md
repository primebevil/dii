# ADR-0004: The Consumer Is a First-Class Ingress in the Prototype

Status: Accepted, 2026-07-03. Builds on ADR-0002.

## Context

ADR-0002 established that access is not gated on contribution, and that the network serves non-contributing consumers as a first-class case. It defined the consumer as a lightweight role, a phone or browser that requests work without hosting a model, served by a pod that chooses to serve it. That decision settled the policy question of whether consumers are served. It did not settle how a consumer request is handled inside a node, or whether the prototype should implement the consumer role at all or defer it to a later phase.

The prototype exists in part to recruit a pod. Two use cases now drive the demo:

1. A participant who runs a model, and both shares spare capacity and borrows from others. This is the node runtime, the atom of the system.
2. A participant with no ability to run a model, who only borrows. This is the consumer of ADR-0002.

The question is whether the second case is a separate system to build later, or something the prototype can carry cheaply now. It is worth resolving because the second case is the more human version of the founding story: a person with only a phone still reaches a capable model through their community when centralized access is switched off.

## Decision

Make the consumer a first-class case in the prototype, implemented as a second ingress into the same router rather than as a separate subsystem.

A consumer request is a request that enters the overflow path past the local hop. A node-runner's request begins at the local stage of the overflow order, local then pod then federation. A consumer has no local stage, because it hosts no model, so every consumer request is overflow by definition and enters the router at the pod stage. There is no second router and no separate consumer service. There is one router with two ingress types: a local, trusted ingress for the node's own user, and a remote, authenticated ingress for a sponsored consumer, which skips the local stage.

Because the node already exposes an OpenAI-compatible HTTP endpoint (ADR-0003), the consumer needs no bespoke client. Any existing OpenAI-compatible application, pointed at a serving node's URL instead of localhost, is a consumer client.

## Consequences

The architectural delta for the prototype is small: an authenticated remote ingress, and a router branch that starts a consumer request at the pod stage. Everything else in the router is shared.

The prototype topology becomes three roles rather than two: node A and node B as full nodes that both serve and overflow to each other, and consumer C, a device with no local model, borrowing from the pod. The demo shows peer overflow and pure consumption through one system, and it degrades to local for the nodes and reports honestly for the consumer when capacity is gone.

Making the consumer first-class pulls forward the problems ADR-0002 left open, because a consumer takes without giving and is therefore a higher abuse surface than a trusted peer, not a lower one. The prototype does not solve these; it stubs them and marks the stubs. Consumer identity is a pre-shared token or a vouched allowlist entry in configuration, standing in for the node-identity decision still to be made. Fair-use and Sybil resistance are represented only by the fact that the pod serves a known, configured consumer rather than the open internet. The real mechanisms remain open questions under ADR-0002.

## Open questions

The substance is inherited from ADR-0002: how a pod decides whom to serve and at what fair-use limit, how donated capacity is accounted for without becoming a currency, and how vouching resists abuse without becoming an exclusionary gate. Added here: what consumer identity looks like beyond the prototype stub, which connects to the deferred node-identity decision.
