# ADR-0011: Inter-Node Transport Is the Reused OpenAI-Compatible HTTP Call

Status: Accepted, 2026-07-12. Builds on ADR-0003 and ADR-0004. Validated by the Week-3 prototype (journal/2026-07-12-week3-m4-findings.md).

## Context

The Week-3 plan named one real design question: how node A carries a request to node B and streams the result back. Two candidates. The simple path reuses the OpenAI-compatible HTTP call, where node A acts as an ordinary client to node B's own OpenAI endpoint and relays the streamed tokens to its caller, using one wire format for both the caller-facing API and the node-to-node hop. The alternative is a thin internal RPC carrying a capability-request envelope (manifest match, policy stub, payload) with the OpenAI shape kept only at the edges. The plan recommended starting simple, getting the numbers, and recording whether the conflation of public API and internal transport bit us.

## Decision

Use the reused OpenAI-compatible HTTP call as the inter-node transport for the prototype and Phase 1. Node A calls node B's `/v1/chat/completions` like any OpenAI client, and tokens stream back over SSE passthrough. In the code this falls out cleanly: the `peer` package implements the same `modelserver.Backend` interface as the local model server, so to the router a peer is just another backend behind the seam.

## Consequences

One wire format, minimal code, and full interchangeability: existing clients, nodes, and model servers all speak the same protocol, and a peer is not a special case in the router.

The choice was validated, not just assumed. M4 measured overflow at about 100 percent of the peer's own local throughput and only 20 to 43 milliseconds of added time-to-first-token, so the reused-HTTP hop is effectively free and the feared conflation did not bite for the POC.

The identity note records the cost that comes with it: because a peer forward is an ordinary OpenAI call, it carries no node-to-node identity and the original caller's identity is lost across the hop (docs/Identity_Note_From_Prototype.md). That is an identity gap deferred to the identity ADR, not a transport defect, but the transport is where any attribution or policy metadata would eventually have to travel.

The decision is reversible. A thin internal RPC remains the candidate if a later phase needs to separate the public API from the internal transport, or to carry caller identity, admission proof, or policy fields the OpenAI shape has no place for. Recorded, not adopted.

## Open questions

Whether Phase 2 or production keeps the reused HTTP path or introduces an internal envelope, driven by what identity and policy turn out to need to carry across the hop. Tracked with the forthcoming identity ADR and docs/Identity_Note_From_Prototype.md.
