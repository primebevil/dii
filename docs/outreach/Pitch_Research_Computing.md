# For University Research Computing and ML Labs

Status: v0.1 (2026-07-04). One-page pitch aimed at HPC / research-computing directors and sympathetic lab PIs. Background in Outreach_Guide.md; hardware in Node_Requirements.md.

## The short version

DII is a free, open-source, non-commercial research initiative in the RFC tradition, studying whether participant-owned hardware can provide a resilient floor of capable open-weight AI that keeps working when centralized access is restricted, repriced, or switched off. The motivating event is concrete and documented: in June 2026 a frontier model was disabled worldwide by a US export-control directive within hours, then restored about two weeks later. The research question that follows is a clean one, and it is the kind your group is well placed to help answer: can a federation of trusted pods, running quantized open models on hardware people already own, deliver a dependable baseline of everyday intelligence with no central point of control.

## Why you, and why it is low-friction

Research-computing clusters and lab machines run well below capacity at night, between jobs, and between grant cycles. DII targets latency-tolerant, good-enough workloads on quantized 14B-to-30B open-weight models, which is exactly the kind of load that fits idle capacity without disturbing priority work. The ask can be scoped entirely inside your existing policy: contribute idle cycles, on machines you already control, within limits you set, to a public-good floor. There is no data leaving your boundary unless you allow it, since sovereignty and data-residency are first-class routing inputs, and cross-boundary execution is opt-in and policy-gated.

## The ask, pick the tier that fits

Host an anchor node: one always-on machine (a 24 GB GPU or equivalent) that serves the ~30B floor to a pod and, if you choose, to consumers who have no capable hardware. Or donate idle cycles from existing cluster or lab nodes as best-effort capacity the router uses when local resources are free. Or, lightest of all, let a student group or lab run a node under your umbrella with your blessing. Any one of these makes you a founding pod.

## Why it is worth a researcher's attention

The engineering problems are real and publishable-adjacent: capability routing and scheduling across heterogeneous hardware, federated discovery with no central index, tiered trust and verification, and graceful degradation under partition. The project credits prior work, records its decisions as ADRs, and treats the repository as a living research library, so contributions slot into an honest engineering record rather than a marketing effort. For a lab working on distributed systems, split inference, or local-first AI, a live pod is a testbed as much as a donation.

## The honest caveats

This is deliberately non-commercial and deliberately modest. There is no token, marketplace, or funding path, and none is planned. It does not serve frontier or dual-use capability; it targets a dependable floor of already-published open weights, which also keeps it clear of the proliferation and export-control concerns. Several design areas, including node identity and cross-boundary verification, are open and flagged as such. The near-term deliverable is a working local node and a two-node pod with measured benchmarks, not a large network.

## Next step

A short conversation to see which tier fits your policy, followed by one node stood up on idle hardware and measured. The numbers it produces (throughput over a residential and a campus link, time-to-first-token, overflow overhead) are useful to the project and to your own group's understanding of the workload.
