# BitTorrent

## What it does

BitTorrent is a peer-to-peer file-distribution protocol, designed by Bram Cohen in 2001, where a file is split into hash-verified pieces and peers download pieces from each other while uploading the pieces they hold, forming a swarm. Peer discovery began with a central tracker and moved to a trackerless Mainline DHT, a Kademlia-based distributed hash table, so a client can find a swarm from a content hash with no central index, and magnet links let it join with only that hash. Incentives use tit-for-tat: each peer preferentially uploads to the peers uploading fastest to it. It is included here as an architecture influence, not a competitor.

## Verdict

Won. One of the very few no-central-hub architectures that worked at internet scale and endured, roughly 25 years, which makes it the existence proof for DII's discovery and resilience layer.

## Why

BitTorrent's durability comes from where it put the data and how peers find each other. Because the data lives in the peers and peers are discovered peer-to-peer through the DHT, there is no central server holding the content and no single node whose removal stops distribution; trackers and index sites have been shut down repeatedly without killing DHT-discoverable swarms. It was a dominant share of internet traffic in the 2000s, often cited around a third in 2004, and although that share fell to low single digits by the 2020s as streaming took over, it remains in active use in 2026 for large open-source datasets, Linux images, and game updates. A separate point worth keeping straight: the 2018 acquisition by TRON's Justin Sun and the BTT token are a commercial and crypto layer bolted on top, not the protocol, whose endurance predates and does not depend on them.

The incentive detail matters for DII. Tit-for-tat with optimistic unchoking works because a file piece is non-rivalrous and cacheable: once you hold a piece, uploading a copy costs almost nothing and earns you the pieces you still need, so sharing is positive-sum. It only partially resists strategic free-riders, but it is stable because the thing being shared is not consumed by sharing it.

## How it compares to DII

BitTorrent's discovery and resilience layer ports to DII almost directly, and validates a choice DII already made. DII plans DHT-like federated discovery through signed, gossiped peer lists, and BitTorrent is the proof that trackerless, hash-addressed discovery and swarm resilience work at scale and survive the loss of any node. The content-addressed, self-describing identifier (the magnet link) is a good pattern for naming DII resources, and optimistic-unchoke-style bootstrapping is a good pattern for letting a new node join before it has any reputation.

The incentive layer does not port, and it is important to be clear about why, because it is the same rivalrousness that runs through the whole project. BitTorrent's reciprocity is stable only because file pieces are non-rivalrous, cacheable, and verifiable by hash. Live inference is none of those. A GPU serving one request cannot serve another at the same time, so sharing compute has a real marginal cost, and an inference response is dynamic and request-specific, so it cannot be pinned to a fixed hash and verified the way a file piece can. So DII gets BitTorrent's discovery and resilience for free, conceptually, and gets no help at all on the incentive problem, which it must solve through mission, the node serving its own operator, and cost-offset under ADR-0008 rather than through tit-for-tat.

## Borrow / Avoid

Borrow: DHT-based discovery with no central index, swarm resilience and graceful handling of peer churn, content-addressed self-describing identifiers on the magnet-link model, and optimistic-unchoke-style bootstrapping for new nodes.

Avoid: assuming any part of tit-for-tat solves DII's incentive problem, since it depends on non-rivalrous cacheable content; and confusing the BTT token layer with the protocol, which endured without it.

## Sources

- [BitTorrent, Wikipedia](https://en.wikipedia.org/wiki/BitTorrent)
- [Mainline DHT, Wikipedia](https://en.wikipedia.org/wiki/Mainline_DHT)
- [BitTorrent still dominates, then no longer king of upstream, TorrentFreak](https://torrentfreak.com/bittorrent-is-no-longer-the-king-of-upstream-internet-traffic-240315/)
- [BitTorrent acquired by TRON's Justin Sun, Variety](https://variety.com/2018/digital/news/bittorrent-acquisition-tron-justin-sun-1202841793/)

Note: protocol mechanics and dates are well corroborated. Traffic-share figures are directional and vary by methodology, and current user-count figures come from statistics aggregators rather than primary measurement.
