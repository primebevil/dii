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
