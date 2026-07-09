# Week 2 Research, First Batch, 2026-07-04

Later the same day the Week-1 tidy closed (see journal/2026-07-04-week1-definition.md), Week 2 began. The plan names Week 2 as research: the adversarial teardown was already done, and the next deliverable was one-page write-ups of why each prior project won or stalled. This session produced the first batch and, more importantly, a way of trusting it.

## What was built

Six one-pagers, in research/Prior-Art, covering the projects closest to DII's own substrate: Petals, Parallax, FusionAI, Prime Intellect, Gensyn with Verde, and DisTrO. This group was chosen first, ahead of the DePIN markets, the volunteer-computing precedent, and the coordination models, because decentralized LLM inference and training is what the Week-3 prototype actually has to reason about. The pages use a fixed scorecard so they stay comparable, and at the owner's request each one leads with a plain statement of what the project does and a section on how it compares to DII, meaning what DII offers that the project does not or cannot. That contrast is the point of the exercise, not a decoration on it.

Research was gathered by parallel agents pulling current 2026 figures, then written up in the repository's voice. Two glossary terms surfaced as worth defining while writing and were added: latency bound, which is the physics the whole batch turns on, and DHT discovery, which is how the swarm designs find peers without a central index.

## The finding that matters

Read together, the batch says one thing. The physics and verification research is real and improving, TOPLOC and Verde in particular are strong, but every project that tried to build a market on top of it is stalled, pivoted, or unproven on demand. Petals decayed for lack of any reason to keep GPUs online. Gensyn shipped a prediction market instead of a training market. Prime Intellect's flagship quietly ran centralized. This is the same conclusion the teardown reached from first principles, now met from the evidence side: the cost-and-market axis is where these efforts die, and the sovereignty-and-mission axis is the one left standing. DII's decision to drop the token and borrow Tor's motivation model is not a preference, it is the move the wreckage points to.

## Validation, not trust

The owner framed his workflow as AI validation: he will not repeat a figure publicly without confirming it maps to a real source, so the deliverable had to be checkable, not just plausible. Two things came out of that. A validation checklist now maps every load-bearing claim to its source and the exact spot to confirm it, ordered by leverage. And an independent verification pass re-fetched every source with fresh agents that never saw the write-ups, checking roughly thirty claims against the live pages.

The pass earned itself. Almost everything confirmed verbatim, but two claims did not survive and were corrected. The Gensyn token was reported as down about seventy percent from its high; the token is not actually trading on public trackers yet, so the figure had no support and was removed. Parallax was credited with about 1.6 million nodes; that number is not in the cited source, where the real nearby figure is 1.6 billion peer connections, a different thing, so the count was dropped and the surviving point, that the headline participation is browser Sentry Nodes rather than inference nodes, was kept. Both are exactly the kind of confident, wrong number that would have cost the project credibility if it had gone out unchecked. The corrections are noted at the top of the checklist so the catch is on record.

## Housekeeping and handoff

The six files were numbered to fix a reading order: inference first, since it is DII's substrate, then the training and verification world it borrows from but does not do. The work was pushed to a review branch so the owner could read it away from the desk and return with notes. No code was written this session; the output is research and a method for trusting it.

## Follow-ups

The remaining prior-art batches are staged but deliberately not started, since the owner wants to be present for the next one: the DePIN markets (Akash, io.net), the volunteer-computing precedent (BOINC, SETI@home, Folding@home), and the coordination and resilience models (Tor, BitTorrent, the Fediverse and Matrix, Ray, the IETF DIN draft). The demand probe from the teardown, five dependency-plus-exposure users, remains the highest-lead-time open item and is still the objection most likely to kill the project. The node-runtime component spec and the two-node-plus-consumer workflow remain the next code-adjacent documents when prototype work resumes.
