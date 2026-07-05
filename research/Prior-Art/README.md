# Prior-Art One-Pagers

Week-2 research deliverable: for each prior project in distributed and decentralized AI compute, a one-page write-up of what it does, whether it won or stalled and why, and how it compares to DII. The teardown in [../The-Case-Against-DII.md](../The-Case-Against-DII.md) supplies the falsification frame; these pages supply the evidence behind it.

This first batch covers the six projects closest to DII's substrate, decentralized LLM inference and training, because they most directly inform the Week-3 prototype. Later batches will cover the DePIN markets (Akash, io.net), the volunteer-computing precedent (BOINC, SETI@home, Folding@home), and the coordination and resilience models (Tor, BitTorrent, the Fediverse and Matrix, Ray, the IETF DIN draft).

## Template

Every page follows the same structure so they stay comparable:

1. **What it does.** Plain description of the system and the problem it attacks.
2. **Verdict.** Won, Stalled, or Mixed, in one line.
3. **Why.** The reasoning, with specific figures and dates.
4. **How it compares to DII.** What DII offers that this project does not or cannot, and where the two diverge on purpose.
5. **Borrow / Avoid.** The concrete lessons for DII.
6. **Sources.**

## Reading order

The files are numbered. The first two are decentralized inference, DII's actual substrate, read closest to home first. The next four are the training and verification world, which DII borrows from but does not itself do. [Petals](1-Petals.md) is the foundation everything else reacts to, so start there.

| # | Project | Does | Verdict | What DII offers that it does not |
| --- | --- | --- | --- | --- |
| 1 | [Petals](1-Petals.md) | BitTorrent-style LLM inference across volunteer GPUs | Stalled | A participation model that survives past novelty (mission and community, per Tor), plus pod trust and honest degradation |
| 2 | [Parallax](2-Parallax.md) | Decentralized LLM inference with a global two-phase scheduler | Early traction, unproven at scale | Sustainable non-token participation and access not gated on airdrop farming |
| 3 | [FusionAI](3-FusionAI.md) | Research paper on consumer-GPU decentralized training | Stalled as a system, won as a research thread | An actual network with an economic, trust, and community layer the paper never had |
| 4 | [Prime Intellect](4-Prime-Intellect.md) | Decentralized training and a GPU marketplace, VC and token backed | Mixed (tech won, thesis unproven) | Non-commercial serving for the individual, not a compute market for AI labs; no token |
| 5 | [Gensyn / Verde](5-Gensyn-Verde.md) | Trustless decentralized-training market with on-chain verification | Stalled on thesis (pivoted to prediction markets) | A live purpose today, serving people, without a token-first market; pod trust instead of trustless proofs |
| 6 | [DisTrO](6-DisTrO.md) | Low-communication optimizers for training over the internet | Partial win (bandwidth solved, competitiveness unproven) | A different goal entirely: serving a reliable floor now, not training a frontier model later |

The [Validation-Checklist](Validation-Checklist.md) is a companion, not part of the sequence: keep it open while reading to confirm each load-bearing figure against its source.

The through-line: the physics and verification research is real and improving, but every project that tried to build a market on it is either stalled, pivoted, or unproven on demand. DII's divergence is to drop the market and the token and borrow Tor's motivation model instead, which is the one axis the evidence leaves standing.
