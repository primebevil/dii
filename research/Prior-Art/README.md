# Prior-Art One-Pagers

Week-2 research deliverable: for each prior project in distributed and decentralized AI compute, a one-page write-up of what it does, whether it won or stalled and why, and how it compares to DII. The teardown in [../The-Case-Against-DII.md](../The-Case-Against-DII.md) supplies the falsification frame; these pages supply the evidence behind it.

Batch 1 covers the six projects closest to DII's substrate, decentralized LLM inference and training. Batch 2 covers the volunteer-computing precedent, the projects that test whether mission-driven volunteer infrastructure can sustain, which is the assumption DII's participation model rests on. Still to come: the DePIN markets (Akash, io.net) and the coordination and topology influences (BitTorrent, the Fediverse and Matrix, Ray, the IETF DIN draft).

## Template

Every page follows the same structure so they stay comparable:

1. **What it does.** Plain description of the system and the problem it attacks.
2. **Verdict.** Won, Stalled, or Mixed, in one line.
3. **Why.** The reasoning, with specific figures and dates.
4. **How it compares to DII.** What DII offers that this project does not or cannot, and where the two diverge on purpose.
5. **Borrow / Avoid.** The concrete lessons for DII.
6. **Sources.**

## Batch 1: decentralized LLM compute

The files are numbered. The first two are decentralized inference, DII's actual substrate, read closest to home first. The next four are the training and verification world, which DII borrows from but does not itself do. [Petals](1-Petals.md) is the foundation everything else reacts to, so start there.

| # | Project | Does | Verdict | What DII offers that it does not |
| --- | --- | --- | --- | --- |
| 1 | [Petals](1-Petals.md) | BitTorrent-style LLM inference across volunteer GPUs | Stalled | A participation model that survives past novelty (mission and community, per Tor), plus pod trust and honest degradation |
| 2 | [Parallax](2-Parallax.md) | Decentralized LLM inference with a global two-phase scheduler | Early traction, unproven at scale | Sustainable non-token participation and access not gated on airdrop farming |
| 3 | [FusionAI](3-FusionAI.md) | Research paper on consumer-GPU decentralized training | Stalled as a system, won as a research thread | An actual network with an economic, trust, and community layer the paper never had |
| 4 | [Prime Intellect](4-Prime-Intellect.md) | Decentralized training and a GPU marketplace, VC and token backed | Mixed (tech won, thesis unproven) | Non-commercial serving for the individual, not a compute market for AI labs; no token |
| 5 | [Gensyn / Verde](5-Gensyn-Verde.md) | Trustless decentralized-training market with on-chain verification | Stalled on thesis (pivoted to prediction markets) | A live purpose today, serving people, without a token-first market; pod trust instead of trustless proofs |
| 6 | [DisTrO](6-DisTrO.md) | Low-communication optimizers for training over the internet | Partial win (bandwidth solved, competitiveness unproven) | A different goal entirely: serving a reliable floor now, not training a frontier model later |

Read [Synthesis](Synthesis.md) after these six: why the prior art clusters into split-inference and distributed-training, why those are complex, and why DII, which does neither, is a different bet rather than a simpler one.

## Batch 2: does volunteer infrastructure sustain

The question DII's participation model depends on. Read as a build: three volunteer-computing cases that show the plateau, then Tor as the exception that shows how sustained volunteer infrastructure actually works.

| # | Project | Does | Verdict | Lesson for DII |
| --- | --- | --- | --- | --- |
| 7 | [SETI@home](7-SETI-at-home.md) | Screensaver donating idle CPU to scan for alien signals | Won proof of concept, stalled on durability | Mass scale came from near-zero donation cost; ended when useful work and staff ran out |
| 8 | [BOINC](8-BOINC.md) | Shared middleware for volunteer science projects | Won durability, stalled on scale | Durable for 23 years but declined to a small core; money (Gridcoin) failed to grow it |
| 9 | [Folding@home](9-Folding-at-home.md) | Volunteer GPU/CPU for protein-folding disease research | Won scale and science, stalled on retention | GPU donors can be mobilized by crisis but receded ~98%; plan for the core, not the spike |
| 10 | [Tor](10-Tor.md) | Volunteer relay network for anonymity and censorship circumvention | Won | Sustained 20+ years via nonprofit steward plus institutional operators, but on cheap non-rivalrous bandwidth |

The batch-2 finding: mission-driven volunteer infrastructure can be durable, but every case that sustained did so on cheap, non-rivalrous donation (spare cycles or spare bandwidth). DII's ask is expensive, rivalrous GPU inference, the one variable none of the durable cases had to overcome. Folding@home shows GPU donors can be moved by mission but do not retain against real cost. The implication points DII toward Tor's model more heavily than Tor itself needed it: a nonprofit steward, a backbone of committed institutional operators, a small trusted core rather than mass scale, and everyday self-interested utility (the node serves its own operator) as the durable fuel the pure-donation projects lacked. This aligns with DII's pod and community-scope design rather than against it.

The [Validation-Checklist](Validation-Checklist.md) is a companion to both batches: keep it open while reading to confirm each load-bearing figure against its source.

The through-line across both batches: the market-based efforts stalled on demand, and the pure-donation efforts stalled on scale or retention. DII's bet is the narrow path between them, mission and community at pod scale, with benefit ungated on the Tor model, which is the one combination the evidence has not already killed.
