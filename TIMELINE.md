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
