# Legal Review Scope: Denylist and Data Handling

Status: Draft for legal-counsel scoping, 2026-07-30. Companion to ADR-0016 (Consumer Admission and Identity) and docs/Governance_And_Abuse_Resistance.md, both of which named a mandatory legal review as a gate before public serving ships. This document turns that gate into a concrete list of questions.

This is not legal advice. It is written by the project to be handed to counsel, and to separate the questions we can answer ourselves by design from the ones that need a lawyer.

## How to read this

Each item below is tagged one of two ways:

- **[Design-addressable]**: we can shrink or remove this exposure by our own decision, before spending counsel time on it. A recommendation is stated.
- **[Needs a legal answer]**: this survives any reasonable design change and requires counsel. The question is framed for them, with our current stance noted so they can react to something concrete.

The point of the split is to keep the counsel engagement short and pointed at real unknowns, not at features we could have dropped.

## The inherent tension, stated plainly

Moderating bad actors and holding no personal data are not separable. The moment a ban can follow an actor across pods, the network has created a durable record about a person, and someone becomes the controller of that record with the accuracy, appeal, and erasure duties that follow. ADR-0016 already names this in its third consequence: the smallest store in the design carries the largest legal weight. So the goal here is not to escape the obligation, which is inherent to the function, but to concentrate it where it is smallest and cheapest to hold, and to be honest about what remains.

The load-bearing observation for the whole review: almost all of the legal weight sits in one optional feature, the shared cross-pod denylist. Remove or defer it and what remains is local moderation, which is comparatively clean. The sections below build on that.

## 1. The denylist as regulated personal data

A denylist of stable cross-pod pseudonyms tied to real bad actors is personal data under GDPR-style regimes, and whoever publishes it is a data controller. This is the single heaviest object in the design.

The largest lever is a design one. **[Design-addressable]** The shared cross-pod denylist is the source of the durable PII store, the controller question, the joint-controllership question, and the erasure-versus-propagation conflict below. Local moderation does not need it: a pod can revoke a token, kill the in-flight request, and report where the law already requires, all without a shared store. ADR-0015 already made cross-pod overflow optional, and the governance notes already accept pod-hopping as the cost of a no-chokepoint system. Recommendation: defer the shared denylist entirely for the first public door. Ship local-only banning, and treat cross-pod propagation as a later feature that turns on only once its governance and legal mechanism exist. This removes most of what follows in this section from the near-term critical path.

If and when the shared denylist does ship, these questions are for counsel and do not go away by design:

**[Needs a legal answer]** Who is the controller of the published denylist, and can there be joint controllership among the steward that publishes it, the pods that contribute listings, and the pods that subscribe? The design keeps the steward off the inference data path, but publishing a list of persons is a controller act regardless of the data path. We need the controllership map pinned before anyone publishes.

**[Needs a legal answer]** A denylist whose entries relate to criminal conduct (for example CSAM-related bans) may be criminal-offence data, which several regimes restrict more tightly than ordinary personal data. Does maintaining such a list require a specific legal basis or official authority, and does that change who is allowed to hold it?

**[Needs a legal answer]** What is the lawful basis for computing and checking the pseudonym at admission for every consumer, including the innocent ones who are never listed? Our stance is that abuse prevention is a legitimate interest, and that the pseudonym is computed transiently and persisted only for actual bad actors. We need confirmation that the transient computation itself clears the bar, and what transparency it obliges us to give.

## 2. Where data is actually collected, and the access-versus-privacy conflict

There are three collection points, each with a different owner.

The external identity anchor (a library patron system, an employer SSO) holds the real personal data. **[Design-addressable]** This is the cleanest of the three because it reuses someone else's existing, already-compliant system rather than building a new one. Recommendation: keep leaning on external anchors and hold no anchor PII inside DII, as ADR-0016 already commits.

The per-pod pseudonym for ordinary users is unlinkable across pods by design, so the innocent leave no cross-pod footprint. **[Design-addressable]** Preserve this property in any concrete credential mechanism chosen later; it is the main privacy win and it is free to keep.

The privacy-preserving denylist design creates its own conflict, and this one I do not think the docs have fully surfaced. The design wants membership to be checkable but the list not browsable, so a pod can test an identity without downloading a roster of the banned. That is a real privacy gain, and it fights the data-subject rights of access, rectification, and erasure directly. **[Needs a legal answer]** How does a listed person learn they are listed, contest an entry, or obtain removal, when the list is built to be un-enumerable? Is a non-browsable denylist compatible with subject-access and rectification duties at all, or does compliance force a queryable record that partly defeats the privacy design?

## 3. Erasure and correction across a decentralized feed

**[Needs a legal answer]** The propagation model that makes bans stick is opt-in, multi-publisher, and cached by subscribing pods. That is also what makes removal unreliable: if a listing is wrong, no single party can guarantee deletion everywhere it has propagated. Erasure and rectification regimes assume a controller who can actually delete. What does a defensible erasure and correction process look like when the architecture deliberately has no central deletion point? This is another reason the section 1 recommendation to defer the shared denylist is attractive, since local bans have a single owner who can unwind them.

## 4. Retention and logging defaults for public-serving pods

The node plan emits one structured log line per request (node_id, door, consumer handle, model, status, latency, token counts), and the data-handling section commits to ephemeral by default with no prompt retention beyond what liability and active moderation require. That last clause is where minimization and moderation collide.

**[Design-addressable]** We can set concrete, conservative defaults ourselves: what fields are logged, for how long, and with what separation between the content stream (what was asked) and the identity stream (who asked), which the governance notes already say must never be joined. Recommendation: draft default retention and logging values as a design artifact, then have counsel confirm them rather than invent them.

**[Needs a legal answer]** Some obligations may require retention that conflicts with our minimization default (see section 5 on CSAM). Where a reporting duty forces retention, what is the minimum we must keep, for how long, and how do we reconcile that with the ephemeral-by-default promise we make to consumers?

## 5. Operator liability and mandatory reporting

This is the part that cannot be designed away, because it comes with being a public door at all. A public-serving pod runs inference on a volunteer's hardware, and serving strangers means carrying liability for what they generate.

**[Needs a legal answer]** What liability does an individual or small group running a public-serving pod actually take on for content that strangers generate through it, and does any intermediary-liability protection apply to an inference provider of this kind? We should not assume the protections that apply to hosting platforms transfer to this setting.

**[Needs a legal answer]** CSAM is the reportable hard line, and reporting mechanics differ by jurisdiction. In the US, running a public door may pull the operator into provider-style reporting obligations, which carry their own retention requirements and can conflict with ephemeral-by-default. What are the operator's affirmative duties on detection, what must be preserved and for how long, and to whom must it be reported, per jurisdiction we expect operators to run in?

**[Design-addressable, with a legal check]** We can reduce, though not remove, this exposure: serve a safety-tuned model by default, run input and output classifiers at the remote ingress, and revoke fast, all of which the governance notes already propose. Recommendation: make the safety model plus classifier a default for public-serving pods rather than opt-in, and ask counsel whether these measures materially affect the operator's liability position or merely the practical risk.

**[Needs a legal answer]** Should the steward be a formal legal entity, and should public-serving pod operators run under an operator agreement, terms of service for consumers, and some indemnification or insurance structure? This shapes whether volunteering to be a public door is survivable for an ordinary person, which ADR-0009 already flags as the hardest supply problem.

## 6. Jurisdiction and cross-border reach

**[Needs a legal answer]** Pods, consumers, external anchors, and any denylist publisher can each sit in a different country. Which jurisdiction's law governs a listing, a retention default, or an operator's reporting duty, and does DII need to constrain where public-serving pods or denylist publishers may operate to keep the obligations tractable? We would rather learn the constraints now than discover them per incident.

## What we are deciding ourselves versus asking counsel

Design decisions we can make without counsel, and that shrink the surface before the engagement:

- Defer the shared cross-pod denylist for the first public door; ship local-only banning (section 1).
- Keep all real PII at external anchors and hold none inside DII (section 2).
- Preserve unlinkable per-pod pseudonyms for ordinary users (section 2).
- Draft conservative default retention and logging values, content stream separated from identity stream (section 4).
- Make the safety model plus classifier a default, not opt-in, for public-serving pods (section 5).

Questions that survive design and need a lawyer:

- Controllership and joint-controllership for any published denylist, and whether criminal-offence-data rules apply (section 1).
- Whether a non-browsable denylist can coexist with subject-access, rectification, and erasure rights (sections 2 and 3).
- Erasure and correction across a decentralized, multi-publisher feed (section 3).
- Retention forced by reporting duties versus ephemeral-by-default (sections 4 and 5).
- Operator liability for stranger-generated content, and whether any intermediary protection applies (section 5).
- CSAM detection duties, preservation, and reporting mechanics per jurisdiction (section 5).
- Steward legal entity, operator agreements, consumer terms, and indemnification or insurance (section 5).
- Governing jurisdiction and any constraints on where pods and publishers may operate (section 6).

## Recommended sequence

Make the design decisions above first, since they cost nothing and remove several of the counsel questions from the near-term path (most importantly, deferring the shared denylist takes sections 1 through 3 largely off the table for the first public door). Then take the surviving list to counsel. The near-term unblock for running a node for yourself and for trusted members is unaffected by any of this; only the unaffiliated public door is gated, exactly as the node plan already states.
