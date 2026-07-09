# io.net

## What it does

io.net is a DePIN on Solana that aggregates GPUs from data centers, crypto miners, and individuals into a compute marketplace for AI training and inference, coordinated by the IO token and pitched as centralized-cloud compute at a 50-to-90 percent discount. Its surfaces are IO Cloud (rentable clusters) and IO Intelligence (a low-cost inference platform). It launched in 2024 with a 30-million-dollar Series A led by Hack VC at a one-billion-dollar token valuation. Despite the branding, the core is closed-source with no on-chain governance; it is closer to a centralized GPU reseller with a token attached.

## Verdict

Stalled, in exactly the way the demand objection predicts. Abundant idle supply chasing paying demand that did not show up.

## Why

io.net won the easy side and lost the hard one. It registered roughly 327,000 GPUs by early 2025, but the daily average of verified active GPUs was about 6,720, near 2 percent of the registered base. The decisive figure is on the demand side: over a sampled fortnight in June 2026, clients bought about 35,700 IO while suppliers were paid about 1.37 million IO, so client fees ran at roughly 2.6 percent of supplier payouts. The network pays its suppliers overwhelmingly with token emissions, not client money, which means the advertised revenue-backed token burn is really emission-funded. Third-party trackers put annualized revenue near 12.5 million dollars in mid-2026, down from the 20-million-plus run-rate that Q1 2025 growth had implied. The IO token trades near 0.17 dollars against a June 2024 high of 6.44, a fall of about 97 percent.

The supply numbers were also literally inflated. In April 2024 io.net was hit by a Sybil attack in which it acknowledged around 400,000 spoofed workers, after exposed API tokens let attackers forge device metadata to farm rewards. The low real-utilization figure is partly the hangover from that. There is a thin genuine business underneath, including named enterprise deals and confidential-compute offerings, but it is a small reseller with a token, not a thriving decentralized market.

## How it compares to DII

io.net is the cleanest evidence for the teardown's Objection 2, and it flips the usual worry about DII. The fear that a distributed network cannot muster enough compute is backwards: registering hundreds of thousands of idle GPUs is trivial, and supply was never the bottleneck. Durable paying demand is, and a cheaper-than-AWS pitch on a commodity attracts token farmers and tire-kickers rather than sustained workloads.

What DII offers that io.net does not is an escape from the mechanism that sank it. Because DII is non-commercial with no token and no marketplace, there is no emissions treadmill paying suppliers against thin demand, and no pressure to manufacture network revenue, which is the loop that inflated io.net's supply and then collapsed its token. But io.net also sets DII's hardest homework honestly: the reason anyone shows up cannot be that idle GPUs are cheap. DII's demand has to be mission-shaped, the privacy, sovereignty, and access cases that only a network like this can serve, because the price story is necessary-but-worthless on its own.

## Borrow / Avoid

Borrow: almost nothing operationally, but take the diagnosis. Treat supply as the easy side and demand as the scarce side, and design the demand thesis first. Note that io.net eventually reached for confidential compute and privacy-sensitive local handling, which is a tacit admission that the non-price axis is where the real demand lives.

Avoid: the entire emissions-subsidy model, which manufactures supply the market does not need and prices a token on it; equating registered supply with a working network, since 2 percent utilization is the reality behind a large headline; and any revenue framing that cannot survive an honest client-to-payout ratio.

## Sources

- [io.net deep-dive (utilization, 2.6% client-to-payout, $12.5M), Own Your Mind](https://ownyourmind.ai/projects/io-net/)
- [IO.net: does it have what it takes, Nansen](https://research.nansen.ai/articles/ionet-does-it-have-what-it-takes)
- [DePIN revenue pivot, BlockEden](https://blockeden.xyz/blog/2026/04/12/depin-revenue-pivot-token-subsidies-ai-compute-akash-render-ionet/)
- [io.net price, ATH $6.44 Jun 2024, CoinMarketCap](https://coinmarketcap.com/currencies/io-net/)
- [io.net April 2024 incident report](https://ionet.medium.com/25th-april-incident-report-176e5fb5c576)

Note: the ~$12.5M figure is from an aggregator republishing io.net's own reporting, not an independent audit. The $20M-annualized and revenue-backed-burn claims are io.net marketing, contradicted by on-chain data. The utilization, ~2.6% ratio, and Sybil-attack figures are third-party or self-acknowledged.
