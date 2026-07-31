# Timeline

## 2026-06-12 (context, pre-project)

- US Commerce Dept export-control directive; Anthropic disables Fable 5 and Mythos 5 worldwide. Access restored around June 26 to July 1. Founding case study in docs/Why_DII_Exists.md.

## 2026-07-02

- Project started.
- Core hypotheses established.
- Adversarial teardown written (research/The-Case-Against-DII.md).
- Founding rationale captured (docs/Why_DII_Exists.md): resilience and reliable-floor positioning.
- Decided access is not gated on contribution; the network serves non-contributing consumers (docs/ADR/ADR-0002).

## 2026-07-03

- Architecture working session, focused on the path to a Week-3 prototype (journal/2026-07-03-architecture-session.md).
- Chose Go for the node prototype, Rust reversible at the production boundary (docs/ADR/ADR-0003).
- Made the consumer a first-class ingress in the prototype: one router, two ingress types, consumer as a request past the local hop (docs/ADR/ADR-0004).
- Folded the router ingress model and updated prototype topology into the Architecture Overview.
- Settled the participation model: motivation is mission and community, not reciprocity or money; monetization retired, not paused; positioning made explicit as AI for All, in the spirit of Tor (docs/ADR/ADR-0005). Scrubbed remaining money-making hedges and reframed participation across the docs.

## 2026-07-04

- Continued the working session; completed the Week-1 definition (journal/2026-07-04-week1-definition.md).
- Pinned a concrete reliable-floor definition: the 14B-to-30B open-weight class at Q4, defined by a usefulness cliff rather than raw hardware; node-entry ~14B on a gaming GPU, promise ~30B on committed hardware, degrading to honest queuing below the useful line (docs/ADR/ADR-0006, docs/Reliable_Floor_Definition.md).
- Sharpened the target user from "everyone" to a dependency-defined target: broad mission (the individual, AI for All), narrow recruiting wedge (independent professional and small business where dependency meets present-tense exposure), personal pod-zero (docs/ADR/ADR-0007, docs/Who_DII_Is_For.md).
- Recorded influences, with Tor as the foundational one (INFLUENCES.md).
- Week 1, "define the proposition," complete.
- Began Week 2, research (journal/2026-07-04-week2-research.md). Wrote the first prior-art batch, the six projects closest to DII's substrate: Petals, Parallax, FusionAI, Prime Intellect, Gensyn/Verde, and DisTrO (research/Prior-Art). Each leads with what it does and how it compares to DII.
- Added a validation checklist mapping every load-bearing claim to its source, then ran an independent verification pass. Roughly thirty claims confirmed; two were corrected as unverifiable (the Gensyn token price and the Parallax node count).
- Added latency bound and DHT discovery to the Glossary. Pushed the batch to a review branch. No code written.
- Wrote a synthesis of batch 1 (research/Prior-Art/Synthesis.md): the prior art splits into split-inference and distributed-training, both complex; DII does neither and is a different bet, not a simpler one. Added glossary terms DePIN, TOPLOC, OpenDiLoCo, H100.

## 2026-07-05

- Wrote Week-2 prior-art batch 2, the volunteer-infrastructure sustainability question: SETI@home, BOINC, Folding@home, Tor (research/Prior-Art, files 7-10, journal/2026-07-05-week2-funding.md). Ran a second independent verification pass; corrected several unverifiable figures.
- Finding: mission-driven volunteer infrastructure can be durable, but every sustained case ran on cheap, non-rivalrous donation; DII's expensive, rivalrous GPU ask is the variable none overcame.
- Accepted ADR-0008, funding the stewards. The network stays non-commercial (ADR-0005 holds), but funding a nonprofit steward and offsetting operators' bare-minimum costs to break-even is permitted as distinct from monetization. Reconciled the Architecture Overview and ADR-0005 wording accordingly.
- Wrote prior-art batch 3, the DePIN compute markets: io.net, Aethir, Render, Akash (files 11-14). Finding: supply is abundant and demand is the scarce side; where demand is real it attaches to enterprise hardware or price competition DII avoids, so DII needs a non-price thesis. Render validates the latency-tolerant workload choice; Akash confirms Objection 5.
- Wrote prior-art batch 4, the coordination and topology influences: BitTorrent, the Fediverse and Matrix, Ray, the IETF DIN work (files 15-18). Finding: DII's federated topology is proven and its differentiation (the trust and sovereignty layer, not the scheduler) is defensible; the incentive problem and the funded steward are DII's own to solve.
- Ran independent verification passes on both batches (journal/2026-07-05-week2-research-complete.md). Week-2 prior-art research complete: eighteen one-pagers across four batches. Noted authorship transparency in CONTRIBUTING; zero-padded the one-pager filenames.
- Talked the four batches through and turned findings into decisions (journal/2026-07-05-batch-review-design.md; working notes in research/Notes/Final-Report-Notes.md).
- Contribution incentive resolved as mission plus mutual benefit plus cost-offset; named "mutual benefit" to keep it distinct from the Reciprocity Signal, which is now likely unnecessary. Added the term to the Glossary and a note to ADR-0005.
- Accepted ADR-0009, public and private pods and funding eligibility: rivalrous cost caps pod size, the unaffiliated are served by many public-serving pods (Tor exit-relay model) via a decentralized directory rather than a central pod, and funding follows mission (public-serving first, public-interest case-by-case, private self-funded). Added the pod terms to the Glossary and the model to the Architecture Overview.
- Settled the Objection-5 positioning: the moat is the untrusted-volunteer-pod substrate, not the sovereignty-routing feature; deferred the identity-standards choice to Week 3.
- Updated the system diagram to place the consumer outside every pod, reaching in through a public-serving pod's door. Rewrote the prior-art Final Report to integrate all of the above. Week-2 research and its review are complete, nearly a week early.

## 2026-07-06

- Accepted ADR-0010, voluntary sponsorship not paid private use: private pods are invited (not required) to sponsor public access, and a required fee or license for private or commercial use is reserved, not adopted.

## 2026-07-08

- Pivoted from Week 2 (research complete) to Week 3, the prototype (journal/2026-07-08-week3-kickoff.md). Wrote the Week-3 prototype plan (docs/Week_3_Prototype_Plan.md): a two-node, one-pod proof of concept in Go where node A routes a capability request to node B and back, and consumer C borrows through the remote ingress, each node exposing an OpenAI-compatible endpoint.
- Scoping calls: this session produces the plan only, no code yet; node identity stays stubbed with a pre-shared token, and the DIDs-versus-DNSSEC-versus-pod-issued-keys decision is driven by what the prototype's ingress actually needs rather than a paper comparison.
- Plan sets four milestones (walking skeleton, real local inference plus manifest, overflow routing and consumer ingress, residential-link measurement), flags the inter-node transport as the one real design question (start with the OpenAI HTTP shape, keep an internal RPC as a candidate), and proposes kill-criteria thresholds for sign-off before the measurement run. No code written.

## 2026-07-12

- Week 3 prototype nearing completion. Built and validated the concept on the week-3-prototype branch: M1 walking skeleton, M2 real Ollama inference plus manifest build/exchange, M3 local-then-peer overflow with honest degradation, and M4 measurement. Extended the two-node plan to a live three-node pod — laptop hub, atlas over Tailscale, sirius over the LAN.
- M4 kill-criteria all pass (journal/2026-07-12-week3-m4-findings.md): overflow throughput ~100% of the peer's own local, time-to-first-token overhead +20-43ms against a ~200ms budget, consumer path within the overflow envelope, and an honest immediate 503 when no node can serve. The overhead is transport-bound and model-independent, so it generalizes to the 30B reliable floor. Residential reachability, the plan's top risk, was dissolved by Tailscale.
- The overflow thesis is proven technically; durability as local models improve stays a separate market question (docs/Pod_Aggregation_Red_Team.md).
- Identity note captured from the build (docs/Identity_Note_From_Prototype.md): the ingress surfaced three concrete needs — node admission, per-consumer credentials, and caller attribution across hops — to drive the identity ADR from observed requirements rather than a paper comparison.
- Remaining to close the week: record the inter-node transport decision (reused OpenAI HTTP) as an ADR; the identity ADR is the next phase.

## 2026-07-14

- Week 3 marked complete. All five definition-of-done items met and every kill-criterion passed; set docs/Week_3_Prototype_Plan.md to Complete.
- Recorded the build and design decisions as ADRs: ADR-0011 (inter-node transport is the reused OpenAI HTTP call), ADR-0012 (backend is any OpenAI-compatible model server behind a thin interface), ADR-0013 (the pod is the accountability boundary), ADR-0014 (consumer work is preemptible best-effort, floor-access not unlimited).
- Worked the consumer-access and identity-at-scale thread and captured it in docs/Governance_And_Abuse_Resistance.md with four diagrams (diagrams/Pod_Admission, Consumer_Access, Delegated_Admission, Data_Minimization): delegated admission for affiliated consumers, the passport-versus-decentralized fork and preferred middle path, and the data-minimization model (cross-pod revocation equals a stable pseudonym plus a shared denylist; keep the linkage only for the guilty; refuse any behavioral reputation store). Added glossary terms model server, member node, pod operator, admission, and consumer sponsorship, and pinned canonical role terminology.
- Week 4 (pressure test) is underway rather than untouched: the kill criteria were run in M4 and the Architecture Overview was revised to mark Phase 1 validated. What remains in Week 4 is the identity ADR, now fed by the identity note and the governance section, and the durability question of pod aggregation as local models improve (docs/Pod_Aggregation_Red_Team.md).

## 2026-07-27

- Ran the Week-4 variety experiment on the pod (atlas as the main node) using the eval/ rig against the pre-registration (docs/Variety_Experiment_Plan.md). A 54-task pilot showed a significant diversity win that did not replicate at full scale: once the split-and-freeze rule froze a stronger baseline, the 102-task run came back inconclusive-to-null (docs/Variety_Experiment_Findings.md). Two of four instruments were not measuring (the judge tied every open-ended answer, reasoning was saturated at the models' ceiling), the one clean signal on long-context reads as specialization rather than diversity, and equal compute on a single model did not pay and sometimes hurt.
- Reframed the project on the result: accepted ADR-0015, DII is access distribution, not capability aggregation. Diversity demotes from thesis to feature; node-to-node overflow reclassifies from core to optional with the consumer ingress (ADR-0004) as the spine; the directory's job narrows from variety to availability; and consumer admission and abuse resistance, who can use them, is promoted to the next load-bearing decision. Closed the red team's open question 1 and updated the Architecture Overview emphasis to match.
- Kept the eval/ rig as a good-enough-capability measuring tool rather than a thesis prover, including for sizing which real workloads a floor node can serve.
- Accepted ADR-0016, consumer admission and identity: pod-level admission, three-tier (sponsored, delegated federated login, unaffiliated middle path), denylist-not-reputation with linkage kept only for the guilty; committed as a direction, with the credential mechanism, denylist governance, and a legal review named as gates before shipping.
- Week 4 (pressure test) complete. Kill criteria passed (M4), the architecture was revised (ADR-0015, access distribution), and the decisions are recorded as ADRs 0011 through 0016. The remaining threads (the identity-mechanism prototype, denylist governance, legal review, and node/member identity) carry forward as gated follow-ons rather than Week-4 blockers.

## 2026-07-30

- Built M1 of the functional node (docs/Functional_Node_Plan.md), promoting prototype/ from a disposable proof of concept to the node the project actually runs. Five things landed: `POST /v1/embeddings` through the same `modelserver.Backend` seam, completing the reliable-floor bundle of general chat, a coder, and embeddings (ADR-0006); a graceful drain on SIGINT/SIGTERM so in-flight streams finish rather than being cut; `/readyz` alongside `/healthz`, so a node that is up but whose model server has died is visible instead of looking healthy; one JSON accounting line per request, including which node actually served it; and packaging — a container image and compose files, plus a `DII_CONFIG` env var so nothing depends on a working directory.
- Verified end to end on a two-node local pod against real Ollama: both doors, chat streaming and non-streaming, embeddings served locally and overflowed to a peer, honest 503 on an unservable model, 401 on a bad token, readiness flipping to 503 with the backend stopped and back on recovery, and a SIGTERM mid-stream draining to a complete response with its accounting line written. The container image itself is written but unbuilt — the local Docker daemon was unresponsive.
- Three decisions worth keeping. Capability stays a single model-name match: the standard model-list call reports every model a backend serves with no type field, so an embedding model is just another name in the list, and the tempting name heuristic fails on real models (`all-minilm:l6-v2` contains no "embed"). The node does not inject `stream_options.include_usage` to force exact token counts, because it promises to forward request bodies verbatim; on a default stream the completion count is therefore a content-chunk proxy and the record says so, which leaves exact accounting as named open work for when M3's budgets need precision. And the accounting record carries an opaque handle and never prompt or completion text, which is the ADR-0016 data-minimization seam built in from the start rather than retrofitted.
- Also fixed an inherited wart: the node announced every response as `text/event-stream`, including non-streaming ones. The response shape is determined by the request's own `stream` field per the OpenAI spec, so it is now labelled from that rather than guessed.
- Brought the docs into agreement with the build: filled in the architecture/APIs stub, which was a placeholder until M1 gave it something real to describe; recorded the M1 changes and the Phase-0 progress in architecture/Overview.md; noted in architecture/Sketchbook.md that the sketched mental model survived the promotion unchanged, with the manifest shape correction and the missing observability item added to its parking lot; marked prototype/BUILD_BRIEF.md superseded, distinguishing its still-binding hard constraints from its now-retiring non-goals; rewrote prototype/README.md and prototype/DEPLOY.md; and drew diagrams/Functional_Node_M1.svg.
- Deployed it: both nodes now run as containers, atlas as the always-on serving node and the laptop as a spoke that overflows to it, with Open WebUI on the laptop as an ordinary consumer. Standing the real thing up produced four defects that isolated testing had not: a health check that could not execute because a distroless image has no shell; a compose network on a subnet the host firewall dropped, which left a node running with an empty model list; the model-list gap below; and a GUI running every message's title, tags and follow-ups through the 30B chat model. Deployment, not unit tests, is what found all four.
- Changed what `/v1/models` means, and it is the one design decision of the day. It listed only the node's own models, so a node would happily overflow a model it had never advertised: fine for `curl` if you knew the name, impossible for a GUI, which builds its model picker from that list. That made ADR-0004's "any OpenAI-compatible client is already a DII client" true for chat and false for discovery — and discovery is the consumer's whole problem (ADR-0015). It now lists the pod, own models plus peers', with `owned_by` naming where each runs. `/manifest` deliberately did not change: it answers a peer's question rather than a caller's, and republishing borrowed capability would have peers advertising each other's models back and forth. The cost is that a startup-cached list can name a model whose holder has since gone down, which is the clearest argument yet for a periodic manifest refresh.
- Confirmed the GUI is a consumer and nothing more: Open WebUI reached the pod through the node with no node changes, and pointing its retrieval at the node exercises M1's embeddings endpoint. Its own login is not node identity, which is the M2 seam.
- M2 (per-consumer identity), M3 (fair-use and preemption), and M4 (the moderation seam) are next. Serving the unaffiliated public remains gated on the ADR-0016 legal review regardless of what the code can do.
