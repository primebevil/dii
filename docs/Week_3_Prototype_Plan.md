# Week 3 Prototype Plan

Status: Draft, 2026-07-08. Depends on the Architecture Overview, ADR-0003 (Go for the node prototype), ADR-0004 (consumer as a first-class ingress), and ADR-0006 (the reliable floor). Weeks 1 and 2 are complete; this is the Phase-1 build from the Architecture Overview build order.

## Goal

Build a two-node, one-pod proof of concept that routes a capability request from node A to node B and back, and lets a consumer C with no local model borrow from the pod through the remote ingress. Each node exposes an OpenAI-compatible endpoint, so any existing OpenAI-compatible client is a consumer client with no bespoke software. The point of the week is to produce the kill-criteria numbers the teardown demands and to learn what the identity layer actually needs to do before we commit to a design.

The prototype is disposable. It exists to answer questions and to give us something real to recruit a pod with, not to become the production node runtime. Rust stays reversible at the production boundary (ADR-0003); nothing this week should assume Go is permanent.

## Definition of done

The week is done when all of the following hold, on a two-machine setup over a real residential link:

1. Node A receives a request for a capability it cannot serve (or is over capacity for), routes it to node B, streams the result back to the caller, and node A's own user sees a normal OpenAI-compatible streaming response.
2. Consumer C, hosting no model, sends the same OpenAI-compatible request to a serving node and gets a streamed answer through the remote ingress, entering the overflow path at the pod stage.
3. When no peer can serve, a node degrades to its local model and a consumer gets an honest "no capacity" response rather than a hang.
4. We have recorded the measurement set below for local, in-pod overflow, and consumer paths, and compared them against the proposed thresholds.
5. A short findings write-up exists (a journal entry plus a numbers table) stating whether each kill-criterion passed, and an identity note capturing what the remote ingress actually needed.

Anything past this (federation, real identity, reciprocity accounting, sovereignty policy enforcement) is out of scope for Week 3 and belongs to later phases.

## Scope

In scope: one pod, two full nodes plus one consumer, local-then-pod routing only, OpenAI-compatible ingress and inter-node transport, a stubbed remote-ingress auth token, and the measurement harness.

Out of scope, deferred with a reason:

- Federation and cross-pod discovery (Phase 3). Within a two-node pod, discovery is a static config file or LAN lookup; no gossip or DNS-like directory this week.
- Real node identity. Stubbed per the decision below; the prototype should tell us the requirements before we pick DIDs, DNSSEC, or pod-issued keys.
- Sovereignty and trust policy gates (Phase 2). The router will have the seam for a policy gate but will pass everything through for now.
- Reciprocity accounting (Phase 2, and likely unnecessary given mutual benefit). Not built.
- Data protection on untrusted hops. Not relevant inside a single trusted pod.

## Topology and roles

Three roles, one pod, matching ADR-0004. Node A and node B are full nodes that each serve a local model and can overflow to the other. Consumer C is a device with no local model that only borrows. The consumer is not a separate subsystem; it is a second ingress into the same router. There is one router with two ingress types: a local trusted ingress that starts a request at the local stage, and a remote authenticated ingress that skips the local stage and enters at the pod stage.

```
Consumer C (no model) ──remote ingress──┐
                                        ▼
                              [ Node A: router + local model ]
                                        │  local miss / over capacity
                                        ▼
                              [ Node B: local model ]  ──stream──► back to caller
```

## Decisions carried into the build

Go for the node prototype (ADR-0003), Rust reversible at the production boundary. Each node exposes an OpenAI-compatible HTTP endpoint, which is both the caller-facing API and the shape the consumer client already speaks. One router, two ingress types (ADR-0004). The reliable floor is the 14B-to-30B open-weight class at Q4 (ADR-0006): node-entry around 14B on a gaming GPU, the promise around 30B on committed hardware. For the prototype we serve the node-entry class so it runs on the hardware we actually have, and we note where a 30B result would differ.

Node identity is stubbed this week and decided later. The remote ingress authenticates with a pre-shared token or a vouched allowlist entry in configuration, standing in for the real node-identity mechanism (this is exactly what ADR-0004 anticipated). We do not research or pick DIDs versus DNSSEC versus pod-issued keys now. Instead we capture, as we build, what the ingress actually needed from identity: what it had to verify, what it had to persist, and where the stub felt wrong. That note becomes the input to the identity ADR after the prototype, so the decision is driven by observed requirements rather than a paper comparison.

## Inter-node transport, the one real design question

The interesting question for the week is how node A carries a request to node B and streams the result back. Two candidates, to be settled early in the build rather than in advance:

The simple path is node-to-node OpenAI-compatible HTTP: node A acts as a client to node B's own OpenAI endpoint, and streams the tokens back to its caller as they arrive. This reuses one wire format for both the caller-facing API and the inter-node hop, and it is the fastest thing to stand up. The risk is that it conflates the public API shape with the internal transport, which we may not want long term.

The alternative is a thin internal RPC between nodes (a small capability-request envelope carrying the manifest match, the policy stub, and the payload) with the OpenAI shape kept only at the edges. Cleaner separation, more code this week.

Recommendation: start with the simple HTTP path to get the walking skeleton and the numbers, and note in the findings whether the conflation bit us. Do not spend Week 3 building the internal RPC; record it as a candidate if the simple path shows strain.

## Milestones

The build goes in four steps, each one leaving something runnable so review stays incremental and reversible.

M1, walking skeleton. Two Go processes on one machine. Node A receives an OpenAI-compatible request, forwards it to node B over the chosen transport with the model call hardcoded or mocked, and streams a canned response back. No real inference yet. This proves the loop, the streaming, and the process boundaries before any model is involved.

M2, real local inference and the manifest. Each node serves a real node-entry-class model (14B Q4) behind its OpenAI endpoint and publishes a minimal capability manifest: which model, context length, modality, and a current-load signal. Node A's router reads B's manifest to decide it can serve. Local-first works: a request node A can serve runs on A.

M3, overflow routing and the consumer ingress. The router implements local-then-pod: A serves locally when it can, forwards to B on a miss or over capacity, and degrades honestly when neither can. The remote ingress goes in with the stub token, and consumer C (just an OpenAI-compatible client pointed at a node) borrows from the pod. Graceful degradation paths are exercised on purpose.

M4, measurement over a residential link. Move node B to a second machine on a real residential connection, run the harness across local, in-pod overflow, and consumer paths, and collect the numbers. Write the findings and the identity note.

M1 and M2 are the bulk of the engineering; M3 is small because ADR-0004 made the consumer a branch rather than a subsystem; M4 is setup and measurement, not new code.

## Kill criteria and what we measure

The teardown put numbers ahead of belief, so the week has to produce them. For each of the three paths (local baseline, in-pod overflow between two machines, consumer through the remote ingress) we record: time to first token, streaming throughput in tokens per second, total wall time for a fixed prompt and output length, and the routing and network overhead, meaning overflow time minus the peer's own local time for the same request.

An observation that shapes the thresholds: token streaming is tiny in bandwidth terms. A model producing 20 to 40 tokens per second emits a few kilobytes per second, far under any residential upload link, so raw bandwidth should not be the binding constraint. The real risks are time to first token (connection setup, routing, and prompt processing on the far node), residential reachability (NAT traversal and holding a stable connection between two home machines), and per-request overhead that a short interactive request would feel. The harness should therefore stress short interactive prompts, not just long batch generations, because the short case is where overhead hurts most.

Proposed pass thresholds, for the owner to confirm before M4 (these are candidates, not settled):

- In-pod overflow throughput is at least roughly 90 percent of node B's own local throughput for the same request. Rationale: once the first token is out, streaming is bandwidth-cheap, so overflow should cost almost nothing in tokens per second. A large gap means the transport is the problem.
- Overflow adds no more than a small fixed penalty to time to first token beyond node B's local time to first token, on the order of a couple hundred milliseconds over a residential link. This is the number that decides whether remote overflow feels usable interactively.
- The consumer path is within the same envelope as in-pod overflow, since architecturally it is the same path entered one stage in.
- Degradation is honest and bounded: when no peer can serve, the caller gets a local answer or a clear no-capacity response within a small timeout, never a hang.

If overflow throughput collapses or time to first token balloons over the residential link, that is a real finding against the overflow thesis and gets recorded as such, not smoothed over. The point of the week is to find out.

## Test harness

A small script (Go or shell plus a load tool) that fires a fixed prompt set at a target endpoint, records the four measurements per request, and runs the same set against each of the three paths so the numbers are comparable. Keep the prompt set in the repo so the run is reproducible. The output is a numbers table that goes straight into the findings journal entry. Two physical machines are needed for M4; M1 through M3 run on one.

## Open decisions to resolve during or after the week

- Inter-node transport: reuse the OpenAI HTTP shape internally, or a thin internal RPC. Resolve early in M1 by trying the simple path first.
- Node identity: settled after the prototype, driven by the identity note, not before. Feeds the identity ADR.
- Whether the capability manifest needs more than model, context, modality, and load for a two-node pod to route correctly. Let M2 tell us.
- Model and quantization actually used, pinned to the hardware on hand, with a note on how a 30B result would differ from the 14B one we measure.

## Risks

The residential-to-residential connection is the most likely thing to eat time, because two home machines behind NAT do not connect as easily as two cloud hosts; budget for this in M4 and treat reachability itself as a finding. The simple HTTP transport may blur the public API and the internal hop in a way we regret, which is why the alternative is written down rather than discarded. And the honest residual from the whole research arc still stands: uncopyable is worthless if unwanted, so a working prototype proves feasibility, not demand. The demand probe remains the highest-lead-time open item and is not what this week is for.

## Task breakdown

1. Stand up the Go module and two node processes; pick and wire the inter-node transport (M1 walking skeleton).
2. Add real local inference behind the OpenAI endpoint and the capability manifest; make local-first work (M2).
3. Implement local-then-pod routing, the remote consumer ingress with the stub token, and the graceful-degradation paths (M3).
4. Build the measurement harness and the reproducible prompt set.
5. Set up the two-machine residential run, collect the numbers, and write the findings journal entry plus the identity note (M4).
6. Record any resulting decisions as ADRs (transport, and after the week, identity).
