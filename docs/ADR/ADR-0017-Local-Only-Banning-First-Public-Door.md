# ADR-0017: Local-Only Banning for the First Public Door

Status: Committed as a direction, 2026-07-30. A near-term scoping decision that narrows ADR-0016. Rests on docs/Legal_Review_Scope_Denylist_And_Data.md and docs/Governance_And_Abuse_Resistance.md, and builds on ADR-0013, ADR-0014, ADR-0015, and ADR-0016. It commits a direction for the first public door, not a final mechanism and not a permanent stance on cross-pod bans. It is supersedable once the build and a legal pass teach more.

## Context

ADR-0016 committed the admission and identity shape and named three gates before shipping: the credential mechanism, denylist governance, and a mandatory legal review. A review of the legal exposure (docs/Legal_Review_Scope_Denylist_And_Data.md) found that almost all of the legal weight sits in one optional feature, the shared cross-pod denylist. That list is the durable personal-data store, the controller and joint-controller question, and the erasure-versus-propagation conflict. Local moderation, revoking the token, killing the in-flight request, and reporting where the law already requires, carries very little of that weight. This record decides what to build first in light of that.

It also fits the sequencing already in the node plan. M4 builds the moderation seam but keeps the public door shut, and the legal gate applies to serving the unaffiliated public. Running a node for yourself and for trusted, affiliated members is not gated. So there is a non-gated subset to build and test now, and a gated superset to defer.

## Decision

For the first public door, moderation is local to the pod and there is no shared cross-pod denylist. A pod revokes a consumer's token, kills any in-flight request, and reports where legally required, all without persisting a cross-pod identifier or contributing to any shared list. Cross-pod ban propagation is deferred, not rejected. It is the heavier feature that turns on only once its governance and a legal pass exist.

This concentrates accountability where ADR-0013 already puts it, at the pod, and keeps the durable personal-data object out of the design for now. A banned actor can move to another pod and must be banned there too. Gradual, pod-by-pod banning is the accepted behavior of a no-chokepoint network, consistent with the pod-hopping the governance notes already accept and with ADR-0015 making cross-pod overflow optional, rather than a gap to be closed by rebuilding a center.

The decision commits the direction, not the mechanism. The concrete credential and revocation technique is still to be chosen and validated by the build, on the same principle that let ADR-0011 follow the prototype rather than precede it.

## Consequences

On the positive side, the first public door can be built and tested without the shared denylist, so most of the legal exposure in the review is off the near-term critical path, and the counsel engagement, when it comes, is shorter and pointed at what survives. Local token revocation creates almost no durable personal data. The stance composes cleanly with the pod as accountability boundary (ADR-0013) and with preemptible best-effort consumer work bounded by per-identity fair-use (ADR-0014).

The residuals are accepted as real rather than smoothed over:

First, a determined actor pod-hops until banned everywhere, one pod at a time. This is slower than a shared denylist would be, and it is the accepted cost of the design.

Second, the federation edge is unresolved and parked on purpose. A local ban revokes standing at one pod only, and federation moves compute overflow rather than admission or bans, so a ban at Pod A does not ban the actor at Pod B directly. The actor loses only the indirect path that routes through Pod A. This is the exact pressure that tempts the shared denylist back in, and it is left open here.

Third, deferring the shared denylist does not remove the irreducible legal items, which are independent of it: operator liability for stranger-generated content, and CSAM detection, preservation, and reporting mechanics per jurisdiction, which may force retention that conflicts with the ephemeral-by-default commitment.

## Open questions

The concrete credential and local-revocation mechanism is unchosen and is to be validated by the build.

Whether and when to add cross-pod ban propagation, together with the federation-versus-ban interaction parked above, is revisited only when the shared denylist comes off the shelf.

The irreducible legal items above are carried to counsel after the build, per docs/Legal_Review_Scope_Denylist_And_Data.md, which is a living document to be refined against the real implementation.
