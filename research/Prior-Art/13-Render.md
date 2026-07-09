# Render

## What it does

Render Network, created by OTOY and its founder Jules Urbach, is a decentralized GPU-rendering marketplace that matches people needing 3D graphics, VFX, and motion-graphics rendering through OTOY's OctaneRender with GPU owners who run those jobs for token payment. The token whitepaper dates to 2017, wider live use to around 2019 to 2020, and it is one of the oldest and most-cited DePIN projects. It migrated from Ethereum (RNDR) to Solana (RENDER) in November 2023, and in December 2025 launched Dispersed, an AI-compute subnet extending it from pure rendering into AI inference.

## Verdict

Won on demand, stalled on token price. Arguably the DePIN with the most genuine, non-speculative demand, because the workload existed before the token did.

## Why

Render's strength is that rendering is a real job that studios already paid for, and it is exactly the shape a distributed network can serve well: GPU-bound, latency-tolerant, and embarrassingly parallel. The usage is real and recurring. Cumulative frames rendered sit above 65 million, a meaningful share of them in 2025, at a monthly rate around 1.5 million frames, across about 5,600 GPU nodes since inception, so the demand is growing rather than a legacy total. The named jobs are checkable and non-speculative: Hollywood-studio VFX, a Coca-Cola activation on the Las Vegas Sphere, NASA content for the ISS, and an A$AP Rocky music video.

The token is the weak axis. RENDER trades near 1.60 dollars against a March 2024 high above 13.50, down roughly 88 percent, though its Burn-Mint Equilibrium ties token burn to actual job payments rather than pure emissions. The newer AI-compute leg, reported around 38 million dollars of monthly revenue in early 2026, is materializing but young, and how much of it is rendering versus new AI inference is not cleanly broken out, so that part rides the same GPU-shortage wave as everyone else. The durable moat is the rendering demand that was never speculative.

## How it compares to DII

Render is the batch's positive lesson, and it validates DII's central design instinct. It succeeded where general GPU-compute DePINs struggled because it targeted a real, pre-existing, latency-tolerant, embarrassingly-parallel workload instead of trying to be a general compute market. That is precisely the viable envelope the teardown identified and the one DII narrowed itself to: latency-tolerant, node-local inference, where each request fits on one machine and nodes do not talk mid-job. Render is evidence that when you pick that workload shape, distributed GPU works.

The contrast to hold is the incentive. Render still runs on a token and is a commercial venture; its providers render because they are paid in RENDER. DII is neither commercial nor tokenized, so it validates Render's technical thesis while getting no help from Render's motivation mechanism. Render simply paid people; DII has to supply provider motivation by other means, the mission and the fact that a node serves its own operator, with cost-offset to break-even under ADR-0008 rather than a wage. So Render tells DII it has chosen the right kind of work, and leaves the harder question, why anyone runs a node without payment, to DII's own participation model.

## Borrow / Avoid

Borrow: the workload-selection discipline above all, target real, pre-existing, latency-tolerant, embarrassingly-parallel demand rather than a general compute market; and the lesson that demand which predates the token is the durable kind.

Avoid: leaning on the token and burn-mint mechanism, which DII rules out; and over-reading the new AI-compute revenue, which is young and rides the same shortage that inflated the rest of the sector.

## Sources

- [Can Render ride the AI wave in 2026, Disruption Banking](https://www.disruptionbanking.com/2026/02/19/can-render-ride-the-ai-wave-in-2026/)
- [Render statistics (frames, nodes), CoinLaw](https://coinlaw.io/render-statistics/)
- [DePIN revenue pivot, BlockEden](https://blockeden.xyz/blog/2026/04/12/depin-revenue-pivot-token-subsidies-ai-compute-akash-render-ionet/)
- [Render migration to Solana, Nov 2023](https://medium.com/render-token/render-network-completes-successful-upgrade-to-solana-3d7947b60aed)
- [Render price and market cap, CoinGecko](https://www.coingecko.com/en/coins/render)

Note: market data and the Solana migration date are corroborated. Cumulative frame counts vary across sources (around 65M on one tracker, higher on Render's own dashboard), and the ~5,600 figure is GPU nodes since inception, not active (active is lower). The rendering-versus-AI split of the ~$38M monthly revenue is not disclosed. Treat frame and growth figures as directional.
