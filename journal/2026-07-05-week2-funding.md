# Week 2, Volunteer Infrastructure and the Funding Decision, 2026-07-05

The second Week-2 research batch and the decision it forced. Batch 1 had covered the projects closest to DII's substrate; this batch asked a different question, the one DII's own model rests on: can mission-driven volunteer infrastructure sustain, especially when the thing being donated is expensive?

## The batch

Four one-pagers, files 7 through 10 in research/Prior-Art: SETI@home, BOINC, Folding@home, and Tor, read as a build from three volunteer-computing cases to Tor as the exception. The scorecard and the independent verification pass carried over from batch 1, and the verification earned its place again, catching an unverified SETI throughput figure, an over-precise BOINC token number, and two Tor specifics that could not be cleanly sourced. Those were corrected rather than shipped.

The finding is clean and it cuts both ways. Mission-driven volunteer infrastructure can be durable: Tor and BOINC are both past twenty years. But every case that sustained did so on cheap, non-rivalrous donation, spare CPU cycles or spare bandwidth that cost the volunteer nothing. SETI@home and BOINC plateaued or declined even at that near-zero cost. Folding@home is the one case that mobilized GPU donors at scale, and it did so only under acute crisis, then receded about 98 percent from its exascale peak once urgency faded. Tor lasted, but through a funded nonprofit with a paid core team and a backbone of institutional operators, not through mass free volunteering. DII's ask, expensive rivalrous GPU inference as an everyday baseline, is the one variable none of the durable cases had to overcome.

## The decision

The owner's read was immediate and, I think, correct: it is unlikely DII will see large growth of free GPU resources, so relying on free contribution alone is the assumption most likely to break the project. Covering at least the bare-minimum expenses a node imposes on its operator is not a later refinement, it is the precondition for anyone coming on board at all. You cannot ask people to run expensive, rivalrous hardware for the mission and also leave them out of pocket for it.

This became ADR-0008, accepted the same day. The care in it is the line it draws. The network stays non-commercial exactly as ADR-0005 set: no business model, no token, no tradable credit, no paying the individual to profit. What is now permitted is funding that flows into a non-commercial steward and cost-offset that returns an operator to break-even. The test that keeps this consistent with ADR-0005 is direction and destination: money out of the network, or into a contributor's pocket beyond cost, stays forbidden; money into a nonprofit steward, or into offsetting a real out-of-pocket cost, is allowed. Tor is the proof that non-commercial and funded coexist. Gridcoin is the named anti-model for why paying individuals to profit fails. Cost-offset restores the volunteer to break-even; it does not pay them to participate, which is the distinction that keeps the motivation mission-based.

Three mechanisms are now in scope to develop: a grant-and-donation-funded steward paying a small core team, institutional operators as the backbone (DII's existing pod examples, universities, labs, hackerspaces, libraries, companies), and sponsor cost-offset for hardware and electricity. The open work is the legal vehicle, how to keep direct cost-offset from drifting into a wage, and funding diversification to avoid Tor's single-source dependence.

## Documents touched

ADR-0008 written and accepted. The Architecture Overview line that had said no part of the architecture assumes or requires a funding path was corrected to distinguish the network, which assumes none, from the steward, which now does. ADR-0005 received a dated update note pointing to ADR-0008 and clarifying that its retired "managed-federation funding path" was a commercial revenue path, not steward funding, so the two records do not read as contradicting each other. No code this session; the output is research, a decision, and its reconciliation.

## Follow-ups

Turning ADR-0008 from principle into practice is now real work: choosing a legal vehicle, designing cost-offset that cannot become a wage, and a funding-diversification plan. The demand probe should widen to probe for willing institutional operators and sponsors, not only end users. Two prior-art batches remain, the DePIN markets and the coordination and topology influences, when the owner is ready.
