# Week 1 Definition, 2026-07-04

Day two of the working session begun 2026-07-03 (see that day's journal). This morning finished the Week-1 tidy: the reliable-floor definition (item 2) and the target user (item 3), plus recording influences. Decisions captured as ADR-0006 and ADR-0007. No code was written.

## Item 2: the reliable floor

The second Week-1 loose thread was that the floor was defined only qualitatively. Grounding it in current mid-2026 data and the owner's own hardware pinned it concretely. The decision is ADR-0006, with the rolling detail in docs/Reliable_Floor_Definition.md.

The key move was defining the floor by a capability cliff rather than by what technically loads. Below roughly the 14B class a model can chat but cannot reliably reason across steps, use tools, handle longer context, or code, which the owner confirmed from running a 7B. Usefulness sets the line. The floor became a band of three numbers: node-entry at ~14B on a 12 to 16GB gaming GPU, the promise at ~30B on committed hardware such as the owner's EVO-X2 with 64GB unified memory, and a guaranteed floor at the ~14B useful line, below which the network queues or reports honestly rather than serving a toy. 30B is the realistic ceiling on consumer-class hardware; 70B needs prosumer machines and runs latency-tolerant only. The floor is a bundle, a general model plus a coder plus an embedding model, since the router speaks in capabilities.

Two consequences worth keeping. Setting node-entry at the useful line trades node quantity for floor quality, another divergence from Tor, where any spare bandwidth helps but sub-useful inference does not. And because 30B-capable nodes are a committed minority with intermittent presence, the floor is what a pod can reliably keep available, which makes honest degradation a first-class part of the guarantee rather than an error path. This also made the prototype fully concrete: the reference node runs a 30B general model and the kill criteria measure it against centralized serving.

## Item 3: who DII is for

The third Week-1 thread was that the beneficiary was "everyone," the exact framing the teardown's Objection 2 attacks. The decision is ADR-0007, with working detail in docs/Who_DII_Is_For.md.

The owner came at the project personally: if he lost AI access the impact would be large, professionally and personally, and though he can self-host, the question is what someone does who cannot. That reframed the target around dependency rather than demographic. The target is anyone for whom AI has become load-bearing for real cognitive work and who would be harmed by losing it, which excludes casual use and is the same line as the reliable floor: 14B-to-30B rather than a toy, because the target does real thinking. Dependency alone is not recruitable demand, though, since a dependent person with cheap cloud access just uses the cloud; the demand needs a present-tense exposure such as cost, privacy or data-residency, or continuity risk.

The resolution is three layers, not one audience. The mission beneficiary stays broad, the individual, AI for All. The recruiting wedge is narrow, the independent professional and small business or research shop, where dependency meets exposure and the harm is concrete and findable. Pod-zero is the owner and their network. The owner's instinct, personally the individual but ideally leading to the small business, resolves cleanly because the two are one target at two scales, a continuum unified by dependency, with the individual as the base case. This makes the demand probe concrete: find five dependency-plus-exposure users and the community that would serve them.

## Influences: Tor

The owner asked to record influences, and Tor was made the foundational entry in INFLUENCES.md, since it is where the project's motivation, access, and resilience models come from. The analogy is bounded: Tor is the model for motivation, legitimacy, and resilience, not for technical architecture. The file also seats the federation (Fediverse and Matrix), distributed-compute (Petals, BOINC, Akash and io.net), verification (Verde, TOPLOC), and local-first influences already threaded through the repo.

## Week-1 status

Items 1, 2, and 3 are closed and captured. Week 1, "define the proposition," is complete: reliable floor, local-first, federated pods, non-commercial and mission-driven, a concrete floor, and a dependency-defined target.

## Follow-ups

- Run the demand probe from the teardown's kill criteria, now framed as five dependency-plus-exposure users and the community that would serve them.
- Finish the Week-2 write-ups on why prior projects won or stalled, including Tor as a rare successful case of sustained volunteer infrastructure.
- Write the node-runtime component spec and the two-node-plus-consumer workflow, including the failure and reroute paths, when prototype work resumes.
