# 30-Day Research Plan

Week 1, Define. Nail the proposition: reliable floor, local-first, federated pods, non-commercial. Done; see Why_DII_Exists and the Architecture Overview.

Week 2, Research. Prior-art teardown of distributed compute, done in The Case Against DII. Next, one-page write-ups of why each prior project won or stalled.

Week 3, Prototype. Two-node, one-pod proof of concept: route a capability request from node A to node B and back, measuring the kill-criteria numbers. Done; built in Go (prototype/) and validated on a live three-node pod, with all four kill-criteria passing (journal/2026-07-12-week3-m4-findings.md): overflow throughput ~100% of the peer's local, time-to-first-token overhead +20-43ms against a ~200ms budget, consumer path within the overflow envelope, and an honest 503 on no-capacity. Inter-node transport recorded in ADR-0011; the identity note (docs/Identity_Note_From_Prototype.md) seeds the next ADR.

Week 4, Pressure Test. Run the kill criteria, revise the architecture, and record decisions as ADRs. Done. Kill criteria passed in M4; the variety experiment tested the durability question and came back inconclusive-to-null, which reframed the project (ADR-0015, access distribution rather than aggregation; docs/Variety_Experiment_Findings.md); and the decisions are recorded as ADRs 0011 through 0016, including consumer admission and identity (ADR-0016). The identity-mechanism prototype, denylist governance, a legal review, and node/member identity carry forward as gated follow-ons.
