# DisTrO

## What it does

DisTrO (Distributed Training Over-the-Internet), from Nous Research, is a family of optimizers that cut the inter-node bandwidth training normally requires, so large models can be trained across ordinary internet links instead of datacenter interconnects. Its seed algorithm, DeMo (Decoupled Momentum Optimization, co-authored with Adam co-inventor Diederik Kingma), decouples momentum across nodes, applies a discrete cosine transform, and communicates only the top few highest-energy frequency components, reconstructing an approximate update on the other side, an approach it likens to JPEG compression. The production system adds overlapped communication and a one-bit trick that sends only the sign of the coefficients. DisTrO is the optimizer layer; Psyche is the network built on it (Rust, peer-to-peer, with the training coordinator living in a Solana smart contract), and Consilience is the model, a 40B dense model trained over Psyche on 20 trillion tokens.

## Verdict

Partial win. The bandwidth barrier is genuinely lowered and independently corroborated, but no decentralized-trained model has matched a comparable centralized one, and sustained mass participation is unproven.

## Why

The core claim is real and confirmed by outside parties (Epoch AI, Import AI) and by parallel efforts: low-communication distributed pretraining works at billion-parameter scale. Nous completed a 15B run in December 2024 on donated H100s, launched Psyche in May 2025, and ran Consilience 40B, which it calls the largest distributed pretraining run ever by parameters and dataset. The bandwidth reduction is roughly 1,000x in the repeatable figure (74.4 GB down to 86.8 MB on a Llama-2 architecture test), which makes training feasible on about 100 Mbps down and 10 Mbps up.

The unproven half is competitiveness and participation. Epoch AI's independent analysis puts the largest decentralized runs at roughly 1,000x less compute than frontier models and judges that low-communication methods moderately harm performance (on the order of a 1.5x effective-compute loss going from one to eight nodes), concluding that decentralized training is feasible even at frontier scale but unlikely to catch the frontier this decade because the binding constraint is amassing compute and participants, not the optimizer. Nous itself reframed Consilience as a provability milestone and moved on to performance, which is the tell: they proved it can be done, not that the result competes. The number of genuinely independent participating nodes is undisclosed, and practical verification of honest work remains an open problem solved only heuristically.

## How it compares to DII

DisTrO is the most orthogonal project in this batch, and usefully so, because it clarifies what DII is not doing. DisTrO trains; DII serves. DisTrO's goal is to eventually train a frontier-competitive model over the internet; DII's goal is to serve a reliable floor of existing open weights to real people now. That difference means most of DisTrO's hard problems (competitiveness with centralized training, mustering enough compute to matter) simply are not DII's problems. DII does not need to train anything to deliver its guarantee.

What DII offers that DisTrO does not is present-tense usefulness without a token and without waiting on a frontier-scale outcome. DisTrO's value is contingent on a future where decentralized training catches up, which Epoch judges unlikely this decade; DII's value is available the moment a pod can keep a 14B-to-30B model online. If DII ever did explore distributed training, DisTrO's comms-reduction ideas would be the first thing to borrow, but that is explicitly outside the current scope.

## Borrow / Avoid

Borrow, mostly as future reference: frequency-domain sparsification of momentum (DCT plus top-k), the best-validated comms-reduction primitive; one-bit sign quantization of the coefficients; overlapped communication and computation so link latency stops being the bottleneck; epoch-based join and leave with checkpoint resync, which is directly relevant to DII's serving churn even though it comes from a training context; and reliable off-the-shelf peer-to-peer networking (QUIC with UDP hole-punching) rather than building NAT traversal by hand.

Avoid: marketing the best-case "10,000x" when the repeatable figure is about 1,000x; assuming quality parity with centralized training, since there is a real if modest penalty that grows with node count; treating on-chain coordination as required, since it adds a dependency and verification of honest work is still unsolved; and betting on catching the frontier soon, since the binding constraint is participation, not the algorithm.

## Sources

- [The Psyche Network architecture (Nous Research)](https://nousresearch.com/nous-psyche)
- [The next phase of Psyche](https://nousresearch.com/the-next-phase-of-psyche)
- [DisTrO (GitHub, NousResearch)](https://github.com/NousResearch/DisTrO)
- [How far can decentralized training over the internet scale? (Epoch AI, Dec 2025)](https://epoch.ai/gradient-updates/how-far-can-decentralized-training-over-the-internet-scale)
- [Nous Research trains a model across the internet (VentureBeat)](https://venturebeat.com/ai/nous-research-is-training-an-ai-model-using-machines-distributed-across-the-internet)
- [Import AI 413: 40B distributed training run](https://importai.substack.com/p/import-ai-413-40b-distributed-training)

Note: completion of Consilience to a full released base model, the participating-node count, and Nous's reported ~$50M funding at a ~$1B token valuation could not be confirmed from primary sources.
