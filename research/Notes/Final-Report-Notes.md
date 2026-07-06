# Final Report Notes

A running capture of insights surfaced while discussing the prior-art batches, to fold into [../Prior-Art/Final-Report.md](../Prior-Art/Final-Report.md) when it is finalized. Nothing here is final; it is the raw material the discussions are producing.

## Demand (from the batch-3 discussion)

The single hardest open question is demand, and the discussion sharpened it well past the teardown's wording.

Demand for DII is latent and trigger-dependent, not absent. Right now AI access is abundant and mostly affordable, so the pain is not acute and demand looks low. It spikes when a blocker fires: price becomes prohibitive for some, access gets gated or rate-limited, or a data-residency or privacy rule bites. The reason this is not hand-waving is that the project's own founding event was exactly that trigger firing. The June 2026 export-control directive that disabled Fable 5 and Mythos 5 worldwide (see TIMELINE and docs/Why_DII_Exists.md) is a real precedent, not a hypothetical. The trigger has a track record.

Removing monetization clears a repellent but is necessary, not sufficient. The token and payout dynamics of the DePIN projects genuinely repelled mission-aligned people and drew speculators, so dropping them clears the ground. But the research contains its own counter-example: Petals had no monetization, no token, pure mission, and still decayed. So "no token, therefore demand" is too strong. An empty, clean field still grows nothing without the other ingredients, which are real exposure and a community.

This makes the demand probe two-part, and both parts matter. The first is the exposed-now minority, people whose alternative is already nothing: data that cannot go to a hyperscaler, a jurisdiction or profession with a residency rule, someone already priced out or cut off. This is the teardown's "find five" test, with a sharper bar that came out of the discussion: the test is "the alternative is nothing," not "this would be nice to have." That bar also tightens the exposure criterion in ADR-0007. The second part is preparedness, which is the Tor pattern exactly. Most people do not need Tor until the moment they suddenly do, and its value is that it already exists when the moment arrives. Part of DII's reason to exist is to be there and be ready when the next blocker fires, rather than starting the day after.

Build-to-learn is the right move here, and it is safe for a reason worth stating. io.net built enormous supply ahead of demand on a token treadmill and hollowed out. DII builds pod-zero, one machine and a personal network, small, cheap, and non-commercial with no emissions to fund, so building before demand is certain is nearly free, and the pod itself becomes the demand-discovery instrument. The discipline that makes building ahead of demand reckless for a commercial DePIN is exactly what makes it safe for DII: it costs almost nothing.

Positioning line worth keeping: bottom-up versus top-down. Governments and institutions are moving toward sovereign AI from the top down; DII is the individual-first, bottom-up version of the same instinct. AI for All, in the spirit of Tor.

## Where these feed the Final Report

The open-work section on demand should carry the two-part probe and the "alternative is nothing" bar. The latent-and-triggered model, with the founding event as proof the trigger fires, belongs in the batch-3 synthesis and the demand discussion. The Petals "necessary not sufficient" caveat guards against overclaiming that a token-free design creates demand. The bottom-up positioning is a candidate line for the verdict or for outward-facing messaging.

## Batch 4 recap, for discussion (not yet discussed)

The architecture influences, read for what DII borrows rather than as competitors. Unlike the first three batches, this one mostly confirms the build is sound rather than raising new risks, so it may go quickly. The four in one line each:

BitTorrent, won. Trackerless DHT discovery and swarm resilience have endured for about 25 years with no central server for the data and no single off-switch. DII borrows the discovery and resilience layer almost directly, which validates the plan for DHT-like federated discovery. What does not transfer is the incentive: tit-for-tat works only because a file piece is non-rivalrous, cacheable, and verifiable by hash, and live inference is none of those, so BitTorrent solves DII's discovery problem and gives no help at all on the incentive problem.

Federation, the Fediverse and Matrix, won as architecture but hard on sustainability. Signed, DNS-like federation with autonomous instances, defederation as a trust tool, and no global index is exactly DII's intended topology, and it runs in production at scale, so the model is proven. But the failure mode is the one that keeps recurring: operator burnout and steward funding. The Matrix Foundation has had repeated funding crises and is still not at break-even, and Mastodon's founder stepped down citing burnout. This is the strongest outside validation of ADR-0008 and of the founder-sustainability concern.

Ray, won its niche, which is the point. Capability-aware scheduling inside a trusted cluster is already commoditized, and Ray's own security docs say it expects a safe network and trusted code, with isolation enforced outside the cluster. The ShadowRay incident proved the model breaks the moment it leaves a trusted perimeter. So Ray is the clean answer to the teardown's Objection 5: DII must not be a better scheduler, it must be the trust and sovereignty layer Ray refuses to be. Borrow the task and actor vocabulary and Ray Serve's OpenAI-compatible serving patterns; do not try to out-Ray Ray inside a cluster.

IETF DIN, not a product but standards context. This is the IRTF DINRG research group. The specific "DIN draft" is minor and dormant, but the surrounding work is worth aligning with: RFC 9518's definition of centralization as a design acceptance test, the documented re-centralization failure modes as a checklist, and W3C DIDs and DNSSEC or DANE as identity and naming standards to borrow rather than reinvent.

The batch-4 finding: DII's topology is proven and its differentiation, the trust and sovereignty layer rather than the scheduler, is defensible, which is the answer Objection 5 demanded. Two cross-cutting lessons point back at decisions already made. The incentive problem is explicitly DII's own to invent, because no prior art solves incentives for rivalrous compute. And steward funding plus operator burnout is the recurring failure mode of real federations, which validates ADR-0008 and the founder-sustainability amendment.

Open questions to react to when back:

1. The incentive problem is DII's to invent. Is cost-offset to break-even (ADR-0008) enough on its own, or does DII also need the non-monetary priority or reputation layer left open as the reciprocity signal in ADR-0005?

2. Re-centralization gravity is real: the Fediverse and Matrix both funnel users onto a few big instances despite being decentralized. How does DII keep pod-zero-plus-a-few-big-pods from quietly becoming the center? This is a defaults-and-onboarding design question.

3. Standards alignment: is aligning pod identity with W3C DIDs and DNSSEC worth doing early, or is that premature for the two-node prototype and better deferred to Week 3 architecture?

4. Positioning: are we comfortable stating plainly, in outward-facing material, that DII is not a scheduler but a trust and sovereignty layer? That is the crisp Objection-5 answer the Final Report needs, and Ray is the evidence for it.

## Item 1 resolution: the contribution incentive, and mutual benefit vs the Reciprocity Signal

Discussed and resolved. The incentive to contribute is a three-part stack, and only one part is anything new to build:

1. Mission, primary. You run a node so that people who cannot self-host keep a reliable floor of intelligence that cannot be switched off from above. This is the reason, on the Tor model (ADR-0005).

2. Mutual benefit, which comes along with the mission rather than being engineered. By being in a pod you get more than your own machine gives: continuity, since the pod serves you when your node is busy, asleep, updating, dead, or when you are on a phone; and capability amplification, since a pod can run several models of varied use, so each node reaches the union of what the pod runs. One consumer box cannot hold a strong general model plus a coder plus embeddings plus something bigger all at once, but a pod of specialized nodes can, and because the router speaks in capabilities this is automatic. This is real self-interested value even for a fully capable self-hoster, and it is math, not mission.

3. Cost-offset (ADR-0008), which removes the money barrier so none of the above costs the operator out of pocket.

Naming decision: call point 2 "mutual benefit." Do not call it "reciprocity," which is reserved for the Reciprocity Signal, the built priority-credit knob that ADR-0005 demoted and deferred. They are different mechanisms. Mutual benefit is emergent, unledgered, non-gameable, and free; the Reciprocity Signal is a built priority credit that risks drifting toward a currency and tripping Objection 4. "Pod dividend" was considered and rejected, because "dividend" implies a payout, the exact connotation to avoid.

Key consequence: because mutual benefit already does the "what do I get" work structurally, the Reciprocity Signal stays deferred and is probably unnecessary. This is a candidate future decision (an ADR or an Overview note that mutual benefit supersedes the need for a built reciprocity credit), not to be finalized here.

The rivalry worry, that sharing blocks my own time, is answered by a load-bearing design guarantee: local-first, idle-only overflow, operator-controlled policy, and your own requests preempt a peer's. That guarantee has to be airtight and obvious, or contributors quit the moment their box feels slow because of a stranger. It is a hard requirement, not a nicety.

Honest limit retained: for a fully self-sufficient power user in steady state, the personal need for the network is small, so they contribute mostly for mission with mutual benefit as the sweetener. The real need lives with the person who cannot self-host. DII serves the dependent, backstopped by the capable-and-mission-aligned. Owner's own read: with cost-offset plus the pod's multi-model mutual benefit, even a capable self-hoster is happy to join.

Design requirement this surfaced: pod composition should deliberately spread models of varied use across nodes, so the capability-amplification benefit is real rather than incidental. A pod-composition note for the Week-3 prototype and architecture.

Where it feeds the Final Report and beyond: sharpen the contributor value prop into mission plus mutual benefit plus cost-offset, and keep the honest limit that the need lives with the dependent user. Candidate glossary term "Mutual Benefit," to sit beside "Reciprocity Signal" and prevent the two from being conflated.

## Item 2 resolution: re-centralization, public vs private pods, and how funding qualifies

Discussed at length via a worked scenario (Pod Vancouver, three personal nodes; Pod Ontario, a library or university; Jerry in Alberta, a broke consumer with no local pod). Two findings.

Finding one: re-centralization is resisted structurally, for a reason the prior art did not have. The Fediverse and Matrix funnel users onto a few big instances because hosting a user is nearly free, so gravity piles everyone onto the largest. DII cannot do this, because sponsoring a consumer costs real, rivalrous GPU time. A pod that tried to serve everyone would exhaust its members' idle capacity and degrade its own people, who would leave. So rivalrous cost caps pod size and there is no such thing as an everyone-pod; it collapses under its own load. The same rivalrousness that makes contribution hard is here a feature: it caps concentration automatically. The design goal is therefore not uniformity but "no single pod, seed, client, or directory is load-bearing," monitored the way Tor caps and tracks exit concentration rather than pretending it is flat. Re-centralization also creeps in at the deployment layer (the IPFS "Cloud Strikes Back" lesson): bootstrap and discovery seeds, a trust root, a dominant client, common cloud hosting, and the steward itself. Defense: multiple independent replicable seeds, pods admitting pods pairwise rather than a central registry, an open protocol with more than one client, hosting diversity, and above all keeping the steward out of the data path. The keeper principle: separate the steward (a legitimate center for governance and funding) from the operation (which must have no chokepoint). Kill-the-steward test, on the Tor model: if removing the nonprofit stops the network, it has re-centralized; if the network keeps serving, it is fine.

Finding two: AI-for-All for the unaffiliated, plus a fair funding rule, comes from a public-versus-private pod distinction. The earlier "come in through a community" answer was wrong as the primary path, because it leaves the isolated, broke, institution-less person, exactly who the mission is most about, out in the cold. The correct model is the Tor exit-relay analog: some pods opt to serve the public (strangers who give nothing back), and Jerry's request is matched to whichever public-serving pod has spare capacity and fits his jurisdiction, out of many spread across the network. The load distributes, no single pod carries all the Jerrys, and killing one just re-matches him to another. The common path (come in through your library, employer, town) still exists and is easy; the public-serving pods are the guaranteed path that makes "for All" real. The matching directory is kept decentralized (signed, gossiped, replicable, more than one), a discovery aid rather than a chokepoint.

Funding follows mission, and the clean proxy is: does the pod carry cost for people beyond its own members? Public-serving pods do, so they are the funding priority, and this is a concrete on-mission use of ADR-0008 money, seeding regional public pods (a library-of-compute) where communities cannot, since public-serving pods are the hardest to recruit (same as Tor exit operators). Private pods serve only their own trusted members and are fully legitimate and valuable (a lab or firm forming a trusted zone to share models and capacity); they self-fund, because they already get their return as mutual benefit (item 1), so subsidy is inversely proportional to how much a pod already benefits itself. Closed pods are allowed (pod autonomy); the fix for the "why fund a closed node" discomfort is not to force them open but to let them pay their own way. A private pod also cannot freeload on the commons, because federation is bilateral and opt-in, so a pod that gives nothing to others gets no overflow from them. The think-tank case refines the rule: funding eligibility is "does this pod serve the mission," where serving the unaffiliated public is the main path to qualify and public-interest work (research, civic, nonprofit) is a second path. Funding ladder: public-serving pods first, public-interest pods (think tank, library) case-by-case, private self-interested pods self-funded and welcome but not subsidized.

Honest limit retained: the public-serving model only works if enough public capacity exists. If it does not, Jerry honestly queues or gets a smaller model rather than the network centralizing to serve him, consistent with the reliable-floor honesty. "Not enough public capacity yet" is a recruiting and funding problem to grind on, not a reason to exclude Jerry by architecture.

Candidate additions when finalizing (not now): glossary terms "Public-Serving Pod" and "Private Pod"; a possible ADR on the funding-eligibility ladder (public-serving, public-interest, private-self-funded) as it extends ADR-0008; and an architecture note that the consumer on-ramp is a distributed directory of public pods, never a central public pod. The consumer discovery and abuse/identity surface remains open under ADR-0002.

## Item 3 resolution: identity-standards timing

Deferred to Week 3. Whether to align pod identity with W3C DIDs or DNSSEC and DANE is an architecture decision the prototype should inform, and there is no reason to settle it before we know what identity actually needs to do. It stays on the IETF DIN one-pager's borrow list as a "when the prototype needs it" item.

## Item 4 resolution: positioning, and where the moat actually is

Resolved. DII's answer to the teardown's Objection 5 is that it is not a scheduler, so "the scheduler is just a feature" does not land. Scheduling mechanics are table stakes, since Ray and Kubernetes own capability-aware placement inside a trusted cluster. What DII is, smallest to largest: the router itself (commoditized, not the moat); sovereignty, jurisdiction, privacy, and consent as first-class routing inputs; local-first plus honest degradation to a reliable floor; and, the largest and the real point, a federation across untrusted, multi-owner nodes with no central admin, ungated, held together by trust and mission.

The correction that matters: sovereignty routing by itself is absorbable, because a hyperscaler can add data-residency as a feature. The un-absorbable part is the substrate, a non-commercial federation of untrusted volunteer pods, because copying it means contradicting an incumbent's two foundational assumptions, one trusted administrative domain and a paying-customer relationship. Ray's own security documentation (it assumes trusted code and a safe network) is the evidence that flipping this is a different architecture, not a toggle. So the moat is the substrate, and sovereignty routing is the visible product riding on it.

Reconciliation: the Overview's "router and scheduler are the core engineering contribution" still holds. What DII builds is the router; what defends it is the trust-federation-mission substrate. Different things, both true.

Honest residual: this positioning closes Objection 5 (defensibility) but hands the weight to Objection 2 (demand). Uncopyable is worthless if unwanted, so item 4 and the demand probe are linked, and the strong defensibility answer must not paper over the open demand question.

Plain-language version for outward messaging: the clever routing is table stakes; keeping data in-country is a checkbox anyone can add; the only thing no company can copy is that DII is a free volunteer network, because being that means giving up being a business. Candidate one-liner, draft only: "DII is sovereign, resilient AI access for the people the cloud won't or can't serve, a federation you can run yourself, not a compute marketplace."

## Open question: sponsorship from private pods to fund public access

Raised at the end of the 07-05 session; resolved 07-06 in favor of voluntary donation (ADR-0010). A private pod (a law firm, a medical-research group, a corporate team) gets real value from the software and network DII built, even though it shares only within itself. Could DII ask such pods for money, structured so it funds the public open resources rather than becoming a business? The network stays free and ungated either way, because a private pod runs its own hardware in its own closed zone and consumes no public capacity, so any charge is for the software and support, not for access to the public floor. That distinction is what keeps this from contradicting ADR-0002 and ADR-0005.

Two versions, very different:

1. Voluntary sponsorship (clean, and already basically permitted). Invite well-resourced private pods to sponsor the steward or fund public-serving pods, framed as private beneficiaries funding public access, a cross-subsidy for the mission. This is the Tor, Mozilla, and Wikipedia corporate-sponsor model, and it is just an explicit application of what ADR-0008 (institutional sponsors) and ADR-0009 (private pods self-fund) already allow. Could be normed as a "suggested sustaining contribution" for commercial private pods, expected socially but not enforced. Probably captures most of the money with no cost on principle.

2. Required fee or license for private or commercial use (a real departure). Does not violate "the network is free," since it charges for the software, not public access, but it introduces three costs: paying customers exert pull and pull drifts a project off mission (the batch-2 and batch-3 lesson); enforcement requires either trust or a non-commercial / source-available license, which turns "free software for all" into "free for some"; and it complicates the clean non-commercial story.

Recommendation to pick up later: start with the voluntary, normed sponsorship, keep the required-license option in reserve if voluntary funding does not sustain, and do not commit to the hard version yet. Candidate to become an ADR (extending ADR-0008 and ADR-0009) once decided.
