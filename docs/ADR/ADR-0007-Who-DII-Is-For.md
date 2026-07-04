# ADR-0007: Who DII Is For

Status: Accepted, 2026-07-04. Sharpens the "everyone" in the Vision, Charter, and Why DII Exists. Depends on ADR-0002 and ADR-0005, and answers Objection 2 of The Case Against DII.

## Context

The project's beneficiary was described as "everyone," or as "individuals, communities, and organizations." That is the exact framing the teardown's Objection 2 attacks: aggregating undifferentiated demand is not a market, and "everyone" is not a target anyone can build for or recruit. Week 1 needs a target sharp enough to run a demand probe against.

Two ideas sharpen it. First, the real defining variable is dependency, not demographic. The target is anyone for whom AI has become load-bearing for real cognitive work, professionally or personally, and who would be materially harmed by losing it. This excludes casual use by definition, so it is not "everyone," and it is the same line drawn in ADR-0006: the floor is 14B-to-30B and not a 7B toy precisely because the target does load-bearing thinking rather than trivia. Second, dependency alone is not recruitable demand, because a dependent person with cheap, available cloud access will simply use the cloud. Recruitable demand is dependency plus a present-tense reason the cloud underserves them: cost, privacy or data-residency, or intolerance of the service being switched off.

## Decision

Define the target in three layers rather than as one audience, because the mission, the go-to-market, and the first pod are different jobs.

The mission beneficiary is broad and stays broad: the individual who depends on capable AI and would be harmed if it vanished and they could not self-host. This is "AI for All," the north star, and it is deliberately not narrowed. Everyone who depends on AI for real work deserves a floor that cannot be switched off, the way everyone deserves electricity.

The recruiting wedge is narrow: the same individual at the scale where the stakes are legible, the independent professional and the small business or research shop, where dependency meets a present-tense exposure such as cost, data-residency, or continuity risk. This is who the demand probe targets and who the first real deployments serve, because here the harm is concrete and the demand is findable.

Pod-zero is narrower still and honest: the project owner and their network, technically capable, high-dependency, mission-aligned people who can be reached directly. This is not a compromise; it is how Tor and most infrastructure began, built by and for the community that believed in it before it spread.

The individual and the small business are not two targets but one target at two scales, a continuum from individual to solo professional to small shop to small organization, unified by load-bearing dependency. The mission speaks to the whole line; building and probing start where the stakes are sharpest; and it stays honest to the individual because the individual is the base case.

## Consequences

This makes the demand probe actionable. Kill-criterion #3 stops being "find five people who would use this" and becomes "find five dependency-plus-exposure users, an independent professional or small shop who already feels the cost, privacy, or continuity problem, and the community that would serve them."

It also names a tension the design must carry rather than hide: need and ability are inversely correlated. The people who most need DII, cut off or phone-only or in a restricted jurisdiction, are the least able to contribute and the hardest to reach, while the people easiest to recruit can already self-host and need it least. This is exactly the gap the Tor model exists to bridge, with the capable and motivated serving the constrained out of mission (ADR-0002, ADR-0005). The target definition holds both the node-runner and the served consumer on purpose; dropping either collapses the project into hobbyists serving themselves or a charity with no supply.

Because the wedge is defined by exposure the user already feels rather than by a future off-switch, the present-tense reasons, cost, privacy, data-residency, and continuity risk, are the messaging and the probe, while the founding off-switch scenario remains the mission that gives them meaning.

## Open questions

The exact first wedge segment to probe, for example a specific profession or a specific data-residency-constrained field, is a go-to-market choice for the demand probe, not an architecture decision. How to reach and vouch for served consumers who cannot self-host connects to the open abuse-surface and identity questions under ADR-0002.
