# ADR-0005: Motivation Is Mission, Not Money

Status: Accepted, 2026-07-03. Builds on ADR-0002. Resolves the participation-model open decision in the Architecture Overview and answers Objection 4 of The Case Against DII.

Update (2026-07-05): the funding question deferred in this record's open questions was opened and accepted in ADR-0008. Funding a nonprofit steward and offsetting operators' bare-minimum costs to break-even are permitted as distinct from monetization; the network commitments below, no business model, token, tradable credit, or paid marketplace, stand unchanged. The "managed-federation funding path" retired below referred to a commercial revenue path taken from the network, which remains retired; ADR-0008 permits only funding flowing into the non-commercial steward.

## Context

Two related things were left unsettled. The participation model was described as reciprocity in the Manifesto, FAQ, Charter, and Architecture Overview, where "contribution earns priority" read as the reason people take part. At the same time the Overview listed "reciprocity versus pure-voluntary" as an open decision, and the teardown's Objection 4 showed that reciprocity-as-motivation has a known failure mode: a non-transferable credit is worth only the compute the holder will personally consume, so once a contributor's own demand is met the incentive to keep contributing falls to zero, and making the credit transferable to fix that re-creates the currency the project rejected.

Separately, a monetization idea kept leaking back into the documents even after the commercial framing was retired, in the form of a hedge about an optional managed-federation funding path that could be revisited later.

Both threads point to the same decision about why anyone participates.

## Decision

The motivation for participating in DII is mission and community, not self-interest, and not money. The model is Tor: people run a Tor relay because they believe in privacy for everyone, and most people who benefit from Tor never run a relay at all. DII is the same shape for capable AI. You run a node so that your community, your family, your town, your lab, and people who cannot self-host retain access to a reliable floor of intelligence that cannot be switched off from above. AI for All, in the spirit of Tor.

Two consequences follow, and this record commits to both.

Reciprocity is demoted from the reason to participate to an optional fairness mechanism. Because ADR-0002 already guarantees access without a contribution gate, reciprocity can never be a gate and can never be the engine. At most it is a non-gating priority signal that decides who is served first when a federation is congested. If it is built at all, it must be illiquid: non-transferable, non-tradable, expiring, never a currency. Whether to build it is a Phase-2 question and may prove unnecessary, since within a small trusted pod social accountability does the fairness work that a ledger would otherwise do.

Monetization is retired, not paused. There is no business model, token, tradable credit, paid marketplace, or managed-federation funding path in scope, and the documents should not hold the door open for one. The forward-looking hedges are removed. The earlier commercial framing remains preserved as history in the superseded Startup Concept note. If a funding question ever genuinely arises it will be opened fresh as its own decision, rather than lingering as a standing hedge that quietly shapes the design.

## Consequences

Adopting the Tor framing does more than settle motivation. It weakens the teardown's Objection 2, the demand objection, which assumed a market: abundant supply, thin paying demand, idle compute worth little. None of that binds when there is no price and nothing for sale. Abundant volunteer supply is what a mission network wants, and "nobody will pay for idle compute" is irrelevant when nobody is paying. The open question shifts from "will anyone pay for this" to "will anyone use this," which is a far lower and more plausible bar, and which the demand probe in the kill criteria should now be framed around.

The framing also sharpens why the pod and community scope is load-bearing rather than optional. Running a Tor relay is cheap and forwards spare bandwidth; running a DII node is expensive and rivalrous, because a token generated for a peer is one not generated for yourself. Global anonymous altruism can sustain a relay but will strain to sustain a GPU, so contribution is sustained at community scale, where the benefit is visible and personal. The overflow design bounds the rivalry by serving peers mainly from idle capacity.

The public-facing documents adopt the AI-for-All framing explicitly, and the participation language across the repository is moved from reciprocity-primary to mission-primary. Tor also becomes a prior-art case worth studying in the Week-2 research, and a rare successful one, since it shows mission-driven volunteer infrastructure sustaining for two decades through a nonprofit and institutional operators rather than payment.

The analogy is bounded. Tor is the model for motivation and legitimacy, not for technical architecture; onion routing has nothing to contribute to the capability router.

## Open questions

Whether a reciprocity fairness knob is worth building at all, and if so its exact illiquid form, remains a Phase-2 decision. How pods sustain themselves in practice, following the Tor Project's nonprofit-and-institutional model, is a research and governance question for later. The consumer abuse-surface, identity, and fair-use questions remain open under ADR-0002.
