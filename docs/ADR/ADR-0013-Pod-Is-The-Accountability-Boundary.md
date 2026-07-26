# ADR-0013: The Pod Is the Accountability Boundary

Status: Accepted, 2026-07-12. Builds on ADR-0002 and ADR-0009. Reasoning in docs/Governance_And_Abuse_Resistance.md.

## Context

Abuse prevention, fair-use, banning, content moderation, and admission all need an accountable party, and the reflex is a central authority that can approve, police, and revoke. But DII forbids a chokepoint: a central authority over access is the exact off-switch the project exists to remove, and it is a breach target and a subpoena target as well. Something has to carry these responsibilities without becoming that switch. Across every scenario worked through (denying a consumer, isolating a bad pod, serving the public, moderating content) the same answer recurred, so it is worth recording as a principle rather than rediscovering each time.

## Decision

The pod is the unit of accountability. Admission (who may join as a member node or be served as a consumer), moderation (what a pod's own hardware is asked to produce for strangers), fair-use, and revocation are the pod operator's responsibility, scoped to that pod. There is no global authority over these. Concretely: a bad consumer is refused pod by pod; a malicious pod is handled by defederation and directory delisting, not a central kill switch; and content inspection is legitimate only at a public-serving ingress, where an operator watches their own liability surface, never over a member's own use or a private pod's internal traffic. The steward is a legitimate center for funding and norms but must never sit in the data path, per the kill-the-steward test.

## Consequences

Every enforcement mechanism becomes local and federated, on the Fediverse model, so no single actor can switch the network off. The accepted cost is that there is no global ban and a determined bad actor can pod-hop or self-host; DII guarantees only that no operator is ever forced to serve them and that refusing is cheap and fast. When cross-pod propagation of a ban is wanted, it is delivered by opt-in reputation and denylist feeds, never a central registry (see the identity work in docs/Governance_And_Abuse_Resistance.md). This boundary is the frame that the forthcoming identity ADR, the moderation approach, and consumer admission all sit inside; it fixes where responsibility lives, not the mechanisms, which remain open.

## Open questions

The mechanisms that hang off this boundary are deferred: node and consumer admission, Sybil-resistant identity, and the exact form of opt-in reputation feeds. This ADR records only that the pod, not a center, is where accountability sits.
