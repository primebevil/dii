# Tor

## What it does

Tor is free software plus a volunteer-operated network of relays that anonymizes internet traffic through onion routing, giving privacy and censorship circumvention to roughly two to two and a half million daily users, almost none of whom run a relay. It is stewarded by The Tor Project, a US nonprofit founded in 2006; the software's alpha launched in 2002, out of onion-routing work at the US Naval Research Laboratory. It is DII's foundational influence, the source of its motivation, access, and resilience models, and this page examines the one question that matters most for DII: how Tor sustained volunteer infrastructure for more than twenty years.

## Verdict

Won. The rare, verifiable success case of volunteer-run infrastructure lasting past two decades, and still growing in 2026. The mechanism was not volunteers alone.

## Why

Four mechanisms did it, and none is "lots of volunteers." First, a nonprofit steward with a paid core team: salaried staff funded by grants and donations maintain the software, network health, and operator support, so continuity does not depend on volunteer availability. This is the direct answer to the labor problem that helped end SETI@home. Second, a backbone of high-capacity institutional operators: capacity is heavily concentrated in a small number of high-capacity operators, so the network survives casual-volunteer churn because a committed few carry most of the load. Third, a mission with continuous urgency: censorship and surveillance renew the reason to participate every year, unlike a finite goal. Fourth, benefit is never gated on contribution, so a large served-user base creates the legitimacy and donor constituency that funds the steward.

The network is healthy in 2026, at roughly 9,600 to 9,900 relays and about 2,700 bridges, up from the historical six-to-seven-thousand range. Funding is the visible tension and the visible discipline: the US government share fell from about 80 percent of a much smaller budget in 2012 to 53 percent in 2021-2022 to 35 percent of roughly eight million dollars in 2023-2024, with the rest from corporations, foundations, and individual donations. The concentration that gives durability is also a risk, shown when a single malicious operator ran more than 23 percent of exit capacity in May 2020, which is why Tor now caps any operator's share.

## How it compares to DII

Tor is the model DII copies, and the comparison is bounded on purpose. The transferable lessons carry cleanly: DII should have a nonprofit steward with paid maintainers rather than volunteer-maintained infrastructure, recruit high-capacity institutional operators as the backbone and design for a skewed capacity distribution rather than a flat volunteer base, keep benefit ungated so the served majority builds legitimacy and funding, and provide legal and operational scaffolding for whoever bears external-facing risk.

The lesson that does not transfer is the one DII most needs to notice. A Tor relay forwards cheap, non-rivalrous spare bandwidth, idle uplink the operator was not otherwise using, and that low marginal cost is a structural reason volunteers tolerate the mostly legal rather than financial burden. A DII node runs expensive, rivalrous GPU inference, which breaks the donate-your-idle-scraps economics that underpin Tor, BOINC, SETI@home, and Folding@home alike. The implication is that DII will likely need the nonprofit-steward and institutional-operator model even more heavily than Tor did, because the pool willing to donate expensive rivalrous compute for free is smaller than the pool willing to share spare bandwidth. Expect capacity to be even more concentrated in a few well-resourced operators, and consider cost-offset mechanisms Tor never needed. What DII adds that Tor does not have is that a node also serves its own operator, so contribution and benefit partly coincide, which is a second motivation Tor relays lack.

## Borrow / Avoid

Borrow: a nonprofit steward with a paid core team; institutional and high-capacity operators cultivated as the backbone; benefit not gated on contribution; legal and operator-support scaffolding built early, on the model of the EFF's relay-operator legal FAQ; and deliberate, transparent funding diversification, with published financials, as Tor has done in moving off government dependence.

Avoid: funding concentration, especially heavy single-source government dependence, which is both a strategic and a credibility risk for a trust-critical tool; operator concentration as a centralization and attack surface, which needs caps and continuous monitoring; and underestimating the legal and liability burden on whatever plays the exit-equivalent role. Above all, do not assume Tor's free-spare-capacity economics, since DII's rivalrous cost removes exactly the thing that made volunteering cheap.

## Sources

- [The Tor Project 2023-2024 financials](https://blog.torproject.org/financials-blog-post-2023-2024/)
- [Tor Metrics, network size](https://metrics.torproject.org/networksize.html)
- [The Tor Project, Wikipedia](https://en.wikipedia.org/wiki/The_Tor_Project)
- [Tor history (2002 alpha, Naval Research Laboratory origin)](https://www.torproject.org/about/history/)
- [EFF legal FAQ for Tor relay operators](https://www.eff.org/pages/legal-faq-tor-relay-operators)
- [How malicious Tor relays are exploiting users in 2020, nusenu](https://nusenu.medium.com/how-malicious-tor-relays-are-exploiting-users-in-2020-part-i-1097575c0cac)

Note: relay and bridge counts and the FY2023-2024 funding mix are from primary Tor sources. The count of unique relay operators (as distinct from relays) and the precise current daily-user figure could not be authoritatively verified for 2026.
