# Governance and Abuse Resistance: Design Notes

Status: Working notes, 2026-07-08. Reasoning and open questions, not settled ADRs. Depends on the Architecture Overview, ADR-0002 (access not gated on contribution), ADR-0004 (consumer as a first-class ingress), ADR-0006 (reliable floor), and ADR-0009 (public, private, and public-interest pods). Captures a design conversation about abuse, fairness, data handling, and content moderation. The load-bearing open problem, Sybil-resistant identity without a central authority, is parked at the end and feeds the deferred identity ADR.

## The two promises, held apart

Most of the confusion in this space comes from blurring two different guarantees. DII promises that no single actor can switch off the whole network for everyone. It does not promise that any individual can force a given pod to serve a given request. Keeping these apart resolves the apparent contradiction between "beyond the reach of any single off-switch" and "we need to deny some users." A pod can always refuse a request, and users will sometimes be denied. What no one can do is kill the network. A bad actor losing access pod by pod is decentralized moderation working, not the chokepoint the project exists to remove. Email and Tor and the Fediverse already live here: no one has a right to a specific server, and no one can shut the system down.

## The pod is the accountability boundary

Abuse, fairness, banning, and content moderation all reduce to one problem, accountability for a shared rivalrous resource, and the answer to all of them is the same: push it to the pod, never to a center. The pod is the trust boundary, the legal boundary, and the moderation boundary. There is no global moderator because there is no global anything. This is the Fediverse pattern, each instance governs its own door and federation is bilateral, applied to inference.

A consequence worth stating plainly: this puts expected review and enforcement inside DII, but only at the public-serving boundary. When a stranger reaches into a public-serving pod, that pod's operator is looking at what their own hardware is being asked to do, on their own liability. That is not surveillance of the network; it is an operator watching their own liability surface. A node-runner's own use, and a private pod's internal traffic, stays private and uninspected. Holding that line bright is what lets the sovereignty promise and the moderation duty coexist.

## Denying and banning a single consumer (Fred)

A public-serving pod reviews and regulates the consumers it sponsors, blocks a bad actor, and reports where legally required. Blocking is fast and local: revoke the token or allowlist entry, and kill any in-flight request. There is no global ban, because a global ban list is a global control point, so a determined actor can seek another pod or self-host. That is the accepted cost of every no-chokepoint system, and it is honest to say so rather than pretend otherwise. What DII guarantees is that no operator is ever forced to serve a bad actor and that refusing is cheap.

## Isolating a bad pod

No one holds a button that terminates another pod; such a button would be the chokepoint the design forbids. A pod that turns malicious is handled by defederation and directory delisting. Peers stop routing to it, the decentralized directory stops matching consumers to it, and traffic redistributes exactly as the overflow design intends. The pod can keep serving its own members on its own hardware, and that cannot be stopped any more than self-hosting can. The network cuts off what a bad pod does to others, not its existence. A defederated pod that tries to return under a fresh identity is the Sybil problem at pod scale.

## Fair use and the heavy consumer (Jerry)

A legitimate consumer running an agentic loop can consume enormous capacity without malice. The operator is already protected by the local-first, idle-only, preemptible overflow guarantee, so the heavy consumer cannot degrade the node-runner's own work. The remaining exposure is the pod's finite donated budget and its other consumers. The handling: class consumer work as preemptible, best-effort by default, which fits because heavy agentic workloads are latency-tolerant and yield gracefully; bound each consumer with per-identity fair-use limits in GPU-seconds, tokens, or concurrency; apply backpressure through ordinary rate-limit responses, which an OpenAI-compatible client already understands; and spread one large consumer across many public-serving pods through the directory, on the Tor exit-relay model, so no single pod carries all of them. The honest ceiling is that DII promises floor-access to a capable model, not unlimited free frontier-scale compute (ADR-0006). Past the floor, the correct behavior is to queue or report honestly, not to starve others silently.

## Data handling and privacy

What a request touches: the prompt going in, the tokens coming out, a transient KV cache during generation, and whatever the server chooses to log. The principle follows sovereignty. Inside your own node or private pod, data stays in your boundary and nothing needs to leave. At a public-serving pod, the operator necessarily processes a stranger's prompt on their hardware, so the commitment there is minimization: ephemeral by default, no prompt retention beyond what liability and active moderation require, and clear terms so a consumer knows a sponsored request is not private the way local execution is. Privacy is absolute for your own hardware and conditional the moment you ask someone else's pod to compute for you. Those are consistent as long as the same party is never promised both.

## Content moderation

Moderation lives where the content is produced. Self-hosted and private use has no DII enforcement point and should not pretend to, since policing it would require the surveillance backdoor the project exists to avoid; it is no different from running an open-weight model on a laptop. Public-serving pods, where a stranger produces content on an operator's liable hardware, are where moderation is legitimate and necessary. The layered tools: serve a safety-tuned open-weight model by default as a cheap first line, bypassable by clever prompts and therefore not relied on alone; run an input and output classifier at the remote ingress; and revoke and kill fast. Sexual content involving minors is not a content preference to balance but an illegal, reportable hard line, and it is the strongest reason public-serving pods should run a safety model plus a classifier by default rather than opt-in.

## Fast propagation without a chokepoint

To make a block stick across pods quickly without rebuilding a center, the steward publishes voluntary, opt-in blocklists and reputation feeds, and subscribing pods act on them immediately. Because subscription is voluntary, no single directive silences the network, so the kill-the-steward test holds; the steward publishes signals and norms and never sits in the data path. In practice most public-serving pods would subscribe, so a flagged identity is refused nearly everywhere within minutes.

## The central open problem: Sybil-resistant identity

Every mechanism above, banning Fred, isolating a bad pod, per-consumer fair-use, cross-pod blocklists, is accounted per identity, and each is only as strong as the assurance that an identity cannot be recreated for free. Ban a consumer and they return under a new key; defederate a pod and it reappears as a new one. This is the spine under all of DII's governance, not a detail, and it has no clean answer that avoids a central authority. Vouching, a real person in a pod staking reputation on a newcomer, is the federated candidate, but it trades openness for accountability and risks becoming an exclusionary gate, which is the tension ADR-0002 already flags. The stance here is to resist settling identity on paper and let the prototype's consumer ingress reveal what identity actually has to enforce, then decide (DIDs, DNSSEC, or pod-issued keys) in the identity ADR.

## Open questions

- Identity: the mechanism for Sybil-resistant, portable, per-identity accountability without a central registry. The blocking question; deferred to the identity ADR after the prototype.
- Vouching: how it resists abuse without becoming an exclusionary gate (inherited from ADR-0002).
- Fair-use accounting: the units and defaults for consumer budgets, and how donated capacity is tracked without becoming a currency.
- Blocklist and reputation feeds: their exact form, who curates them, and how a pod appeals a listing.
- Data minimization: concrete retention and logging defaults for public-serving pods.
- Whether any of this warrants its own ADRs (candidates: pod as accountability boundary, consumer work as preemptible best-effort) once identity is settled enough to anchor them.

A visual of the two proposed admission paths, member node versus consumer, is in diagrams/Pod_Admission.svg.
