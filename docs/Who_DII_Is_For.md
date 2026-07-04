# Who DII Is For

Status: Living reference, v0.1 (2026-07-04). The decision behind this is ADR-0007; this doc holds the working detail for targeting and the demand probe.

## The definition, in one line

DII is for people who depend on capable AI for real cognitive work and would be materially harmed by losing it. The target is defined by dependency, not by demographic, which is why it is not "everyone."

## The line: dependency, not casual use

Banana bread is casual. "I could not do this project without it, the research alone would take weeks" is load-bearing. The target lives on the load-bearing side of that line. This is the same line as the reliable floor in ADR-0006: the floor is 14B-to-30B and not a 7B toy precisely because the target does real thinking, not trivia.

## Two axes

Dependency: how load-bearing AI is for someone's work or life. High for researchers, analysts, developers, independent professionals, and anyone who has folded AI into how they think and produce.

Exposure: how likely and how unable they are to tolerate losing centralized access. Driven by present-tense reasons, cost (priced out of frontier tiers), privacy and data-residency (cannot send client, patient, or sensitive data to a hyperscaler), and continuity risk (cannot build on a service that might vanish), as well as the founding case of being cut off by directive.

The target is the high-dependency, high-exposure corner. Pure dependency is not enough, because a dependent person with cheap, available cloud access will just use the cloud; the exposure is what creates demand for DII specifically.

## The three layers

Mission beneficiary, broad. The individual who depends on capable AI and would be harmed if it vanished and they could not self-host. AI for All. The north star, deliberately not narrowed.

Recruiting wedge, narrow. The independent professional and the small business or research shop, where dependency meets a present-tense exposure. Concrete harm (a shop that goes under without the tool), findable, and probe-able. This is who the demand probe and the first deployments target.

Pod-zero, narrowest. The owner and their network: capable, high-dependency, mission-aligned people who can be reached directly. How Tor and most infrastructure started.

The individual and the small business are one target at two scales, a continuum unified by load-bearing dependency. The individual is the base case; the small business is the same person where the stakes are legible.

## The honest tension

Need and ability are inversely correlated. Those who most need DII, cut off or phone-only or in a restricted place, are least able to contribute and hardest to reach. Those easiest to recruit can already self-host and need it least. The Tor model is what bridges the gap: the capable and motivated serve the constrained out of mission. The target holds both the node-runner and the served consumer on purpose (ADR-0002, ADR-0005).

## For the demand probe

Kill-criterion #3 becomes concrete: find five dependency-plus-exposure users, an independent professional or small shop that already feels a cost, privacy, or continuity problem, plus the community that would serve them. Probe the exposure the person already feels today, not the hypothetical future off-switch, and keep the off-switch as the mission that gives the wedge meaning.

## What to revisit here

The specific first wedge segment (a particular profession or data-residency-constrained field) is a go-to-market choice to refine as the demand probe runs. Update this doc with what the probe finds; the layered definition and the dependency-plus-exposure principle come from ADR-0007 and should not drift without a new decision record.
