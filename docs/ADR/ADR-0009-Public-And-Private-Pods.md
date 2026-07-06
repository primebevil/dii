# ADR-0009: Public and Private Pods, and Funding Eligibility

Status: Accepted, 2026-07-06, from the Week-2 batch-4 design review. Builds on ADR-0002, ADR-0005, and ADR-0008. May be refined on the owner's fresh-eyes review of the prior-art Final Report.

## Context

The batch-4 research (research/Prior-Art, the coordination and topology influences) raised a re-centralization worry: the Fediverse and Matrix are decentralized in protocol yet funnel users onto a few large instances. ADR-0008 opened funding into a steward and cost-offset for operators, but it did not say which pods qualify for funding, nor how a person with no pod of their own gets served without a central "everyone-pod" forming. This record settles the pod taxonomy, the rule for who gets funded, and how the unaffiliated are served while staying decentralized. It was worked out through a concrete scenario (a personal pod, an institutional pod, and a broke consumer named Jerry with no local pod).

## Decision

Two structural facts frame everything. First, rivalrous cost caps pod size. Sponsoring a stranger costs real GPU time, so a pod that tried to serve everyone would exhaust its members' idle capacity and degrade its own people, who would leave. There is therefore no viable everyone-pod; it collapses under its own load. The same rivalrousness that makes contribution hard is here a feature, because it resists re-centralization by construction. Second, the goal is not uniformity but that no single pod, seed, client, or directory is load-bearing, monitored the way Tor caps and tracks exit-capacity concentration rather than pretending the network is flat.

Pods come in three kinds, and a pod may be closed; pod autonomy is respected. A public-serving pod chooses to serve unaffiliated consumers, strangers who give nothing back, as an act of mission; it is the DII analog of a Tor exit relay. A private pod serves only its own trusted members, such as a lab or a firm sharing models in a trusted zone. A public-interest pod may be closed to the general public but does work that is itself a public good, such as a research institute or a library.

The unaffiliated are served without a central pod. A person with no community pod is matched to whichever public-serving pod has spare capacity and fits their jurisdiction, out of many spread across the network, through a decentralized directory that is signed, gossiped, replicable, and allowed to exist in more than one, so it is a discovery aid rather than a chokepoint. The common path remains "come in through a community you belong to," a library, an employer, a town; the public-serving pods are the guaranteed path that makes access universal. DII-the-project does not run one big public pod, because that would be the re-centralization; it fosters many public-serving pods and a directory.

Funding follows mission, and the clean proxy is whether a pod carries cost for people beyond its own members. This gives a ladder. Public-serving pods are the funding priority, and seeding regional public pods where communities cannot stand up their own is a primary, on-mission use of ADR-0008 money. Public-interest pods are funded case by case, because public access is the main path to qualify but not the only one. Private pods are self-funded and welcome but not subsidized, because their members already get their return as mutual benefit, so subsidy is inversely proportional to how much a pod benefits itself. A private pod also cannot freeload on the commons, because federation is bilateral and opt-in, so a pod that gives nothing to others gets no overflow from them.

Finally, the steward and the operation are separate. The steward is a legitimate center for governance and funding, but it must never sit in the data path. The test, on the Tor model, is that removing the steward should not stop the network from serving; if it would, the design has re-centralized.

## Consequences

Recruiting public-serving pods is the hardest and most mission-critical supply task, the same problem Tor has with exit operators, and it is the reason steward and sponsor funding exists. The consumer on-ramp is a directory of public pods, not a central pod, which keeps the "no single off-switch" property honest at the deployment layer, not just in the protocol. Mutual benefit (see Glossary) is why private pods need no subsidy, and it is likely why the reciprocity signal left open in ADR-0005 is unnecessary. Access stays ungated (ADR-0002), motivation stays mission (ADR-0005), and money still only flows into the non-commercial steward and cost-offset, never to a contributor as a wage (ADR-0008).

## Open questions

Whether enough public-serving capacity will exist to serve the unaffiliated is a recruiting and funding grind, not settled by this record; until it does, the honest behavior is to queue or serve a smaller model rather than centralize. The directory and matchmaker design, and the consumer abuse-surface and identity questions it raises, remain open under ADR-0002. The exact criteria for public-interest funding eligibility are left for when the funding vehicle exists.
