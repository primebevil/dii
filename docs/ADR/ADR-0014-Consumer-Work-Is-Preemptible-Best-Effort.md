# ADR-0014: Consumer Work Is Preemptible Best-Effort; the Guarantee Is Floor-Access, Not Unlimited Throughput

Status: Accepted, 2026-07-12. Builds on ADR-0002, ADR-0004, ADR-0006, and ADR-0009. Reasoning in docs/Governance_And_Abuse_Resistance.md.

## Context

A consumer, even a well-meaning one running an agentic loop, can consume enormous capacity on donated hardware, and the consumer is the higher abuse surface because it takes without giving (ADR-0004). The question is how a pod serves consumers without letting one starve the pod's finite donated budget or its other consumers, and without a human policing every request. The worry is real but the resolution turns out to make the heavy consumer the easy case, so it is worth pinning.

## Decision

Consumer work is preemptible and best-effort by default. The node-runner's own work and interactive requests take priority; a consumer's work is served from slack capacity and yields under contention. Each consumer is bounded by per-identity fair-use limits (in GPU-seconds, tokens, or concurrency) enforced through ordinary rate-limit responses, which an OpenAI-compatible client already understands as backoff. The guarantee DII makes to a consumer is floor-access to a capable model (ADR-0006), not unlimited throughput and not frontier-scale compute. Past the floor the correct behavior is to queue or to report honestly, never to silently starve other callers.

## Consequences

The heavy consumer becomes the easy case rather than the frightening one, because latency-tolerant agentic work tolerates preemption, so a pod can be generous with idle capacity and strict the moment anything higher-priority appears. The operator is already protected by the local-first, idle-only, preemptible-overflow guarantee, so a consumer cannot degrade the host's own work. Rivalrous cost does the fairness work: a finite donated budget forces per-consumer limits into existence rather than leaving them optional. The honest dependency is that fair-use accounting is per-identity, so its strength rests on the identity layer, which is deferred; until then, limits are enforced per stub credential.

## Open questions

The concrete units and defaults for consumer budgets, how donated capacity is accounted for without becoming a currency (see ADR-0005 and the Reciprocity Signal), and how per-identity enforcement hardens once the identity mechanism is chosen.
