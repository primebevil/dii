# Validation Checklist

A scan-and-confirm map for the Prior-Art batches. Each row is a claim from the one-pagers, what to confirm, the source, and where in the source to look. Purpose is validation, not verbatim reading: land on the right spot, confirm the data point maps, move on.

Priority key: **[H]** high-leverage, a DII decision leans on it, confirm first. **[F]** flagged as shaky in the write-up, confirm the caveat is fair. **[S]** supporting, spot-check only if a number smells wrong.

**How the boxes work.** Each claim is a checkbox. An unchecked box is not yet owner-verified; a checked box marked **Owner-verified** means the owner personally confirmed it against the source. This is a separate and stronger tier than the independent AI passes summarized below, which were done by fresh agents that did not see the write-ups.

**Owner-verification progress: 8 of 36** (28 core claims plus 8 secondary figures).

<details>
<summary>Independent AI passes (background, already done)</summary>

Verification status, batch 1 (2026-07-04): an independent pass re-fetched every source with fresh agents that did not see the write-ups. Of roughly 30 load-bearing claims, all confirmed verbatim except two, which were corrected: the Gensyn token price (no live trading data found, "~70% drop" removed) and the Parallax "1.6 million nodes" count (not in the source; corrected to a Sentry-Nodes-vs-inference distinction without the number). Precision note: Petals' "single-digit tokens/sec" is corroborated by the NeurIPS paper and release notes, not the Yandex blog specifically.

Verification status, batch 2 (2026-07-04): a second independent pass checked the sixteen load-bearing volunteer-infrastructure claims. The core numbers held (SETI's 5.2M lifetime and March 2020 hibernation, BOINC's 700K-to-60K decline, Folding@home's 768 native petaflops and 12.9 as of Oct 2025, Tor's ~9,600 relays and 35% government funding). Corrections to the one-pagers: SETI's "27 teraflops in 2001" was removed (replaced with the confirmed 668 teraflops in 2013); BOINC's precise "98% token drop" was softened to "collapsed" and the direct Anderson "no commercial market" quote was reattributed as a paraphrase; Tor's "about 55 staff" and "five networks control half of exit capacity" were removed as not cleanly sourceable, keeping the cited >23% May 2020 exit-concentration example.

Verification status, batch 3 (2026-07-05): a third independent pass checked the sixteen load-bearing DePIN claims. The core structural figures held (io.net's ~327k GPUs at ~2% utilization and ~2.6% client-to-payout, Aethir's ~$128M FY2025 revenue and 150+ clients, Akash's ~$4M revenue against a ~$227M market cap, the token drawdowns). Corrections: Aethir's token figures (June 2024 high ~$0.106, current ~$0.006, down ~94%, market cap ~$120M); Render's frame count (above 65M, not 67-69M), the "35% of all-time frames in 2025" removed as a misread, "5,600 active nodes" corrected to nodes since inception; io.net's "1.8 million fake GPUs" spike dropped, keeping the acknowledged ~400k spoofed workers. General caveat: most revenue and utilization figures are vendor-self-reported and unaudited.

Verification status, batch 4 (2026-07-05): a fourth independent pass checked the fourteen load-bearing architecture claims. All confirmed in substance, no corrections. The only flags were facts on a sibling page (BitTorrent's 2018 TRON acquisition is on the Variety source, not Wikipedia; the Fediverse's defederation is on the Fediverse page, not the ActivityPub page) and minor wording nuances (DHT is one of several discovery mechanisms; the Matrix SRV record is a deprecated fallback). Ray's security quotes were confirmed via the Oligo writeup because the primary docs page is JavaScript-rendered.
</details>

## Highest-leverage first

If you validate nothing else, validate these three; a DII decision leans on each.

- [x] **[H] TOPLOC verification is cheap and accurate** (Prime Intellect). Confirm: 100% detection in their evals, validation up to 100x faster than inference, 258 bytes per 32 tokens. Source: [arXiv 2501.16007](https://arxiv.org/html/2501.16007v1) (abstract and results). Watch for: these are the authors' own evals, not independent. — **Owner-verified 2026-07-06** (abstract quote confirms 100x, 258 bytes/32 tokens, 1000x memory reduction).

- [x] **[H] Verde verifies at under one order of magnitude overhead** (Gensyn). Confirm: refereed delegation vs roughly 4 orders for SNARK proofs; RepOps overhead 30% on large matmul, ~60% on Llama-3.1-1B; reproducibility guaranteed single-node only, multi-node is future work. Source: [arXiv 2502.19405](https://arxiv.org/html/2502.19405v1) (overhead in evaluation; single-node limit in limitations/future-work). Watch for: the single-node limit is the reason the full training vision has not shipped, so confirm it is stated, not inferred. — **Owner-verified 2026-07-06** (all four confirmed verbatim: "less than an order of magnitude overhead" vs "4 orders of magnitude" for crypto; "less than 30% overhead" on matmul; "around 60% overhead" on Llama-3.1-1B; single-GPU-only stated under "Future work: Model parallelism").

- [x] **[H] Decentralized training stays ~1000x behind the frontier and won't catch up this decade** (DisTrO, and the batch's overall conclusion). Confirm: the ~1000x-less-compute estimate, the "feasible but unlikely to catch the frontier this decade" judgment, and that the constraint is amassing compute and participants, not the algorithm. Source: [Epoch AI, How far can decentralized training scale](https://epoch.ai/gradient-updates/how-far-can-decentralized-training-over-the-internet-scale) (summary and conclusion). The only independent, non-vendor voice in the batch. — **Owner-verified 2026-07-06**.

## Flagged claims, confirm the caveat is fair

- [x] **[F] Prime Intellect funding is ~$20M confirmed, not ~$70M** (Prime Intellect). Confirm: the ~$50M December 2025 round appears only in funding databases with redacted valuation and no primary announcement. Source: [Sacra profile](https://sacra.com/c/prime-intellect/) for the confirmed ~$20M ($5.5M seed + $15M round); the unconfirmed round is on Tracxn/PitchBook/CB Insights only. The point: there is no primary press release for the larger round. — **Owner-verified 2026-07-06**.

- [x] **[F] INTELLECT-3 ran centralized, not distributed** (Prime Intellect). Confirm: 512 H200s over 400 Gbps InfiniBand, and it is a post-train of GLM-4.5-Air, not from scratch. Source: [Implicator analysis](https://www.implicator.ai/prime-intellects-intellect-3-open-source-ambition-meets-centralized-reality/) (stated up top). The "pitch vs execution" gap. — **Owner-verified 2026-07-06**.

- [x] **[F] Gensyn's live product is a prediction market** (Gensyn). Confirm: what shipped to mainnet April 2026 is Delphi, a prediction market, not decentralized training; RL Swarm hosted nodes paused. Source: [Delphi launch](https://news.bitcoin.com/gensyn-network-debuts-delphi-a-permissionless-ai-prediction-market-platform-on-mainnet/). Note: an earlier "token down ~70%" claim was removed as unverifiable; do not cite a Gensyn token price without a live source. — **Owner-verified 2026-07-06**.

- [x] **[F] Petals flagship models are broken on the live swarm** (Petals). Confirm: the advertised large models show as broken or single-partial-contributor now. Source: [health.petals.dev](https://health.petals.dev/) (directly checkable in a browser right now, the easiest validation in the batch). Also confirm last feature release was September 2023 on the [releases page](https://github.com/bigscience-workshop/petals/releases). — **Owner-verified 2026-07-06** (Sept 2023 last release confirmed on the releases page; the live health monitor did not load at check time, itself consistent with a wound-down project, so the "flagships broken" detail rests on the earlier independent pass rather than a fresh owner view).

- [x] **[F] FusionAI is a paper, not a system; its headline number is modeled** (FusionAI). Confirm: the "50 RTX 3080s ≈ 4 H100s" figure is a performance analysis, not a measured run, and there is no code or network under the name. Source: [arXiv 2309.01172](https://arxiv.org/abs/2309.01172) (abstract says "we envision"; the 50-vs-4 figure is in the performance-analysis section). The measured follow-on numbers live in [FusionLLM](https://arxiv.org/html/2410.12707v1). — **Owner-verified 2026-07-06**.

- [ ] **[F] Parallax's headline participation is Sentry Nodes, not inference nodes** (Parallax). Confirm: the large headline counts belong to browser Sentry Nodes for the broader Gradient Network, not Parallax inference capacity (undisclosed); benchmarks are self-reported with HexGen (not Petals) as baseline. Source: [CoinGecko explainer](https://www.coingecko.com/learn/what-is-gradient-network-and-gradient-airdrop) for what Sentry Nodes are; [arXiv 2509.26182](https://arxiv.org/html/2509.26182v1) for the benchmark and baseline. Note: the "1.6 million nodes" figure was removed (nearby real figure is 1.6 billion peer connections, not a node count).

- [ ] **[F] DisTrO's headline reduction is ~1000x, not 10000x** (DisTrO). Confirm: the repeatable figure is ~857x (74.4 GB to 86.8 MB); "10000x" is best-case marketing; participating-node count and full-run completion are undisclosed. Source: [DisTrO GitHub](https://github.com/NousResearch/DisTrO) README for the conservative "three to four orders of magnitude," [VentureBeat](https://venturebeat.com/ai/nous-research-is-training-an-ai-model-using-machines-distributed-across-the-internet) for the 857x and the 15B run.

## Batch 1 supporting numbers, spot-check only if something smells wrong

- [ ] INTELLECT-1 specifics (10B, 1T tokens, up to 112 H100s, ~400x comms reduction): [INTELLECT-1 release](https://www.primeintellect.ai/blog/intellect-1-release).
- [ ] Gensyn $43M Series A led by a16z, 2023: [CoinDesk](https://www.coindesk.com/business/2023/06/11/blockchain-based-ai-compute-protocol-gensyn-closes-43m-series-a-funding-round-led-by-a16z).
- [ ] Parallax two-phase scheduler and MLX support: [Parallax GitHub](https://github.com/GradientHQ/parallax). Petals throughput and the memory-bandwidth wall: already grounded in the teardown, Objection 1.

## Batch 2: volunteer infrastructure

- [ ] **[H] Every durable volunteer-infrastructure case sustained on cheap, non-rivalrous donation** (synthesis of batch 2). Validate the pattern across the four pages: SETI@home and BOINC donate spare CPU cycles, Tor forwards spare bandwidth, Folding@home's GPU surge (the one rivalrous-cost case) receded ~98%. The claim DII rests on: none of the durable cases faced DII's expensive, rivalrous GPU cost.

- [ ] **[F] Folding@home GPU surge receded ~98%** (Folding@home). Confirm: 768 native petaflops at the March 2020 peak versus 12.9 native petaflops as of Oct 31 2025. Source: [Folding@home, Wikipedia](https://en.wikipedia.org/wiki/Folding@home) for both; [Bowman's 2020 review](https://foldingathome.org/2021/01/05/2020-in-review-and-happy-new-year-2021/) for 280K GPUs / 4.8M CPU cores and 400K-new-devices-in-two-weeks. The 98% is arithmetic (12.9 vs 768).

- [ ] **[F] BOINC declined from ~700K to ~60K, and money did not help** (BOINC). Confirm the decline in [Anderson's history essay](https://continuum-hypothesis.com/boinc_history.php) (search "700K" and "diehards"); the Nov 2021 snapshot (34,236 participants / 136,341 hosts / 20 PF) is on [Wikipedia](https://en.wikipedia.org/wiki/Berkeley_Open_Infrastructure_for_Network_Computing). Note: the precise "98% Gridcoin drop" and the verbatim "no commercial market" quote were softened; the substance (commercialization failed to grow the base) stands.

- [ ] **[F] Tor sustained via nonprofit steward plus concentrated operators, on cheap bandwidth** (Tor). Confirm funding and scale in the [FY2023-2024 financials](https://blog.torproject.org/financials-blog-post-2023-2024/) (revenue ~$8M, US gov 35% down from 53.5%) and relay counts via [Tor Metrics](https://metrics.torproject.org/networksize.html) (use the CSV export; the HTML is a JS-rendered graph). The >23% May 2020 exit-concentration incident is in the [nusenu writeup](https://nusenu.medium.com/how-malicious-tor-relays-are-exploiting-users-in-2020-part-i-1097575c0cac). The 2002 alpha and Naval Research Laboratory origin are on the [Tor history page](https://www.torproject.org/about/history/).

- [ ] **[S] SETI@home scale and hibernation**: 5.2M lifetime participants and the March 31 2020 hibernation with its two stated reasons are on [Wikipedia](https://en.wikipedia.org/wiki/SETI@home) and the [SETI Institute announcement](https://www.seti.org/news/seti-at-home-going-into-hibernation/). The "27 teraflops in 2001" figure was removed; the confirmed throughput is 668 teraflops in 2013.

## Batch 3: DePIN compute markets

- [ ] **[H] Supply is abundant; paying demand is the scarce side, and it does not come from cheap GPUs** (synthesis of batch 3, ties to Objections 2 and 5). Validate across the four pages: io.net registered ~327k GPUs at ~2% utilization with client fees ~2.6% of payouts; where demand is real it attaches to enterprise datacenter GPUs (Aethir, Akash) or a pre-existing latency-tolerant workload (Render), driven by price competition with hyperscalers. The claim DII rests on: its consumer substrate is the hard case and its demand must be non-price.

- [ ] **[F] io.net: supply glut, thin demand** (io.net). Confirm ~327k registered GPUs at ~2% utilization and ~2.6% client-to-payout in the [Own Your Mind deep-dive](https://ownyourmind.ai/projects/io-net/) and [Nansen](https://research.nansen.ai/articles/ionet-does-it-have-what-it-takes); token ATH $6.44 / current ~$0.17 on [CoinMarketCap](https://coinmarketcap.com/currencies/io-net/). Caveats: revenue (~$12.5M) and utilization are vendor-reported via aggregators, not audited. The ~400k spoofed-workers figure is io.net's own acknowledgment in its [incident report](https://ionet.medium.com/25th-april-incident-report-176e5fb5c576).

- [ ] **[F] Aethir: real enterprise revenue, collapsed token** (Aethir). Confirm ~$127.8M FY2025 revenue and 150+ clients via [Aethir's Q3 report](https://ecosystem.aethir.com/blog-posts/aethirs-record-breaking-q3) and [BlockEden](https://blockeden.xyz/blog/2026/03/12/depin-compute-revenue-pivot-akash-ionet-aethir/); token figures on [CoinMarketCap](https://coinmarketcap.com/currencies/aethir/) (June 2024 high ~$0.106, current ~$0.006, down ~94%, market cap ~$120M). Revenue is self-reported and unaudited; the ~$130M Checker Node sale is from Aethir/Chainwire.

- [ ] **[F] Render: genuine workload-fit demand** (Render). Confirm the token drawdown (~$13.53 ATH, ~$1.60 now) on [CoinGecko](https://www.coingecko.com/en/coins/render) and the ~$38M monthly revenue via [BlockEden](https://blockeden.xyz/blog/2026/04/12/depin-revenue-pivot-token-subsidies-ai-compute-akash-render-ionet/). Notes: frames above 65M (sources vary), "35% of all-time frames in 2025" removed (it was job volume), ~5,600 is nodes since inception, not active. The Solana migration (Nov 2023) is in the [Render Medium post](https://medium.com/render-token/render-network-completes-successful-upgrade-to-solana-3d7947b60aed).

- [ ] **[F] Akash: won on supply and price, not the scheduler** (Akash, ties to Objection 5). Confirm ~$4M 2025 revenue against ~$227M market cap, and the GB200/Nodekeeper professional-supply move, via the [Akash Q1 2026 report](https://akash.network/blog/akash-network-q1-2026-report/) and [BlockEden](https://blockeden.xyz/blog/2026/03/12/depin-compute-revenue-pivot-akash-ionet-aethir/); AKT market cap and $8.07 ATH on [CoinGecko](https://www.coingecko.com/en/coins/akash-network). Homenode's consumer-card onboarding, justified by sovereignty, is in the Q1 2026 report. Utilization figures conflict (~34% independent vs 80%+ vendor).

## Batch 4: coordination and topology influences

- [ ] **[H] DII's topology is proven and its scheduler is not the moat** (synthesis of batch 4, ties to Objection 5). Validate across the pages: trackerless DHT discovery and swarm resilience endure in BitTorrent; signed, DNS-like federation with no central index runs in the Fediverse and Matrix; Ray, by its own security docs, disclaims the untrusted, federated, sovereignty-constrained regime that is DII's reason to exist.

- [ ] **[S] BitTorrent** protocol mechanics (Bram Cohen, 2001, hash-verified pieces, Kademlia Mainline DHT, magnet links, tit-for-tat) are on [Wikipedia](https://en.wikipedia.org/wiki/BitTorrent) and [Mainline DHT](https://en.wikipedia.org/wiki/Mainline_DHT). The 2018 TRON acquisition is on the [Variety](https://variety.com/2018/digital/news/bittorrent-acquisition-tron-justin-sun-1202841793/) source, not the Wikipedia page.

- [ ] **[F] Federation: the hard part is the steward, not the topology** (Federation). Topology (ActivityPub W3C Rec 2018; Matrix public-key-signed federation, .well-known discovery) is on [ActivityPub Wikipedia](https://en.wikipedia.org/wiki/ActivityPub) and the [Matrix federation spec](https://spec.matrix.org/v1.16/server-server-api/). The load-bearing part is sustainability: the Matrix Foundation's funding warning in the [crossroads post](https://matrix.org/blog/2025/02/crossroads/), and Rochko's burnout departure in [TechCrunch](https://techcrunch.com/2025/11/18/mastodon-ceo-steps-down-as-the-social-network-restructures/). Validates ADR-0008.

- [ ] **[H] Ray: the Objection-5 answer** (Ray). The decisive fact is Ray's own security stance: it "expects to run in a safe network environment and to act upon trusted code," with isolation enforced outside the cluster. Confirm via the [Ray security docs](https://docs.ray.io/en/latest/ray-security/index.html) (JavaScript-rendered; the quote is also in the [Oligo ShadowRay writeup](https://www.oligo.security/blog/shadowray-attack-ai-workloads-actively-exploited-in-the-wild)). Ray joined the PyTorch Foundation Oct 2025 ([blog](https://pytorch.org/blog/pytorch-foundation-welcomes-ray-to-deliver-a-unified-open-source-ai-compute-stack/)).

- [ ] **[S] IETF DIN**: DINRG is active, chartered 2017, no adopted drafts ([datatracker](https://datatracker.ietf.org/rg/dinrg/about/)). RFC 9518's centralization definition is the usable takeaway ([RFC 9518](https://www.rfc-editor.org/rfc/rfc9518.html)). The specific "DIN draft" is minor and dormant.

## Secondary figures: sourced, but not independent-verification targets

Concrete numbers that appear in the one-pagers without their own row above. They come from the sources each one-pager lists and were not re-checked by the independent passes; treat them as sourced-but-not-independently-verified. Tick a box only if you personally confirm it.

- [ ] Parallax's $10M seed (Decrypt).
- [ ] io.net's $30M Series A at a ~$1B token valuation, and the ~6,720 daily-active-GPU figure.
- [ ] Gensyn's ~12,000-node RL Swarm peak (self-reported testnet).
- [ ] Prime Intellect's ~$150K INTELLECT-3 post-train cost (Implicator) and the 32B/106B model sizes.
- [ ] DisTrO's Consilience 40B / 20T-token figures (Nous's own announcement; completion unconfirmed).
- [ ] Aethir's ~$9M venture raise and the Respeecher 66%-vs-AWS claim (vendor).
- [ ] Tor's 2-to-2.5M daily users (directional, not verified for 2026) and the 80%-government-funding-in-2012 figure.
- [ ] SETI@home's ~91k-active / 145k-host 2020 snapshot.
