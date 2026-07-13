# Glossary

**Reliable Floor**: the project's goal, a dependable baseline of capable AI that is always available, distinct from intermittent access to the most advanced ("ceiling") models. Concretely, the 14B-to-30B class of open-weight models at Q4, delivered as a capability bundle: node-entry at ~14B on a gaming GPU, promise at ~30B on committed hardware, degrading to honest queuing below the useful line rather than to a toy. Defined in ADR-0006 and docs/Reliable_Floor_Definition.md.

**Useful Line**: the capability cliff, around the 14B class, below which a model can chat but cannot reliably reason across steps, use tools, handle longer context, or code. The floor is set by usefulness, not by what technically loads; the network degrades toward this line and then queues honestly rather than serving sub-useful output.

**Latency Bound**: a workload whose speed is limited by waiting, not by compute. In distributed inference the wait is the round trip across the network between machines, so adding faster GPUs does not help; the wall is the link. LLM token generation is also memory-bandwidth bound on a single machine, since each token requires reading the whole model from memory. Both are why DII targets latency-tolerant, node-local work that fits on one node rather than splitting a model mid-inference across the internet. See research/The-Case-Against-DII.md, Objection 1.

**Local-First**: every node is fully functional offline; the network only adds to what a node can already do on its own.

**Node**: a participant's machine running one or more open-weight models, exposing them as capabilities.

**Node Runtime (the atom)**: the daemon on a node that serves local models, publishes a capability manifest, and routes requests. Self-sufficient offline; the smallest thing that delivers the core promise to a single user.

**Model Server (Backend)**: the inference engine a node runs in front of, such as Ollama, llama.cpp's server, or vLLM. It loads the open-weight model and generates tokens over an OpenAI-compatible API. The node never infers itself; it brokers to the model server through a thin interface, so the backend is swappable. Distinct from a node: the node is the control-plane broker, the model server is the engine behind it. See architecture/Sketchbook.md and ADR-0011.

**Consumer (Consumer Client)**: a participant with no serving hardware who only borrows, such as a phone or browser, sponsored by a pod that agrees to serve it. Architecturally a request that enters the overflow path at the pod stage, having no local stage of its own. Not a node: there is no "consumer node," because a consumer hosts no model. See ADR-0002 and ADR-0004.

**Ingress**: how a request enters the router. A local, trusted ingress serves the node's own user and starts at the local stage; a remote, authenticated ingress serves a sponsored consumer and starts at the pod stage. One router, two ingress types.

**Capability Manifest**: what a node publishes about itself: which models it serves, context length, throughput, modalities, and current load. Lets a heterogeneous fleet look uniform to a caller.

**DHT Discovery**: finding who serves what without a central index. A distributed hash table spreads the who-has-what map across the participating nodes themselves, so a caller can look up a capability or peer without a single directory that could become a chokepoint or an off-switch. The mechanism Petals uses to map model layers to peers, and the kind of decentralized, no-central-hub discovery DII favors for cross-pod federation.

**Pod**: a small group of nodes that already trust each other, such as personal machines, a lab, an org, or a community. The core unit of the system.

**Member Node (Pod Member, Peer)**: a node that has been admitted to a pod. It both serves its own user and shares overflow with the pod, as a trusted peer of the other members. This is the precise term for a node in a pod; a consumer is not a member node, since a consumer is not a node at all.

**Pod Operator**: the person or group that runs a pod and is its admission authority, deciding which nodes and consumers it accepts and can revoke. Local to a single pod and distinct from the Steward, the network-level nonprofit that funds and governs but never sits in the data path.

**Admission** (proposed, not yet built): a pod approving who is allowed in, as opposed to any global registration. Two kinds: node admission, where a machine becomes a member node, and consumer sponsorship, where a device becomes a sponsored consumer. Knowing a pod's identifier is not enough to join; admission is a separate, revocable step. The prototype stubs this by treating network reachability as membership; the real mechanism is the subject of the forthcoming identity ADR. See docs/Identity_Note_From_Prototype.md.

**Consumer Sponsorship** (proposed, not yet built): the admission path for a consumer, where a pod agrees to serve a specific borrow-only participant and issues it a per-consumer credential presented at the remote ingress. The lower-trust tier of admission, distinct from node admission: a sponsored consumer only borrows and never serves. See ADR-0002 and ADR-0004.

**Federation**: trusted pods sharing overflow capacity under explicit policy, with no central hub. Closer to the Fediverse than to a single cloud. Bilateral and opt-in: a pod federates with pods that reciprocate, so a pod that gives nothing gets no overflow, which is how freeloading is handled without a rule.

**Public-Serving Pod**: a pod that chooses to serve unaffiliated consumers, strangers with no pod of their own, as an act of mission. The DII analog of a Tor exit relay: the hardest role to recruit and the one that makes "AI for All" real. The unaffiliated are matched across many public-serving pods by jurisdiction and capacity, so no single one becomes a center. Funding priority, since it carries cost for people who give nothing back. See ADR-0009.

**Private Pod**: a pod that serves only its own trusted members, such as a lab, a firm, or a family sharing models and capacity in a trusted zone. Fully legitimate and self-funded, because its members already get their return as mutual benefit; not subsidized from shared funds, and it cannot draw on the commons without reciprocal federation. Closed is allowed (pod autonomy). See ADR-0009.

**Public-Interest Pod**: a pod that may not offer general public access but whose work is itself a public good, such as a research institute, library, or nonprofit. Eligible for funding case by case even when closed, because funding follows mission and public access is the main path to qualify but not the only one. See ADR-0009.

**Capability Routing**: requesting intelligence as a capability (reasoning, coding, vision, retrieval) at a required tier; the router then selects the hardware and model.

**Overflow Path**: the order work spills outward, local then pod then federation, only when needed and only within user policy.

**Mission-Driven Participation**: the reason people run nodes, on the Tor model. You contribute because you want a reliable floor of intelligence to exist for your community and for those who cannot self-host, not to earn anything. Access is never gated on contribution (ADR-0002, ADR-0005).

**Mutual Benefit**: the value a node operator gets simply by being in a pod, beyond what their own machine gives: continuity, since the pod serves them when their node is busy, asleep, or offline; and capability amplification, since a pod can run several models of varied use and each node reaches the union of what the pod runs. Emergent, unledgered, non-gameable, and free. It is the self-interested half of why anyone contributes, alongside mission. Distinct from the Reciprocity Signal, and likely makes that built signal unnecessary. See ADR-0009.

**Reciprocity Signal**: an optional, non-gating priority hint that may let a contributor be served first under congestion. If built at all it is illiquid: non-transferable, non-tradable, expiring, never a currency. A Phase-2 fairness detail, not the reason to participate, and not the earlier "proof-of-contribution credits," which are retired. Distinct from Mutual Benefit, which is emergent and unledgered; because mutual benefit already does the "what do I get" work, this built signal is probably unnecessary.

**Sovereignty Constraint**: a first-class routing input, for example "data stays on my machines" or "data stays in-country."

**Target User (dependency-defined)**: not "everyone," but anyone for whom AI has become load-bearing for real cognitive work and who would be materially harmed by losing it. Defined by dependency, not demographic, and sharpened by exposure (cost, privacy or data-residency, continuity risk). See ADR-0007 and docs/Who_DII_Is_For.md.

**Wedge**: the narrow, recruitable slice of the target used to build and demand-probe first: the independent professional and small business or research shop, where dependency meets a present-tense exposure. The individual and the small business are one target at two scales, not two audiences.

**Pod-Zero**: the first pod, the owner and their network of capable, high-dependency, mission-aligned people reachable directly. How Tor and most infrastructure began.

**DePIN Project**: Decentralized Physical Infrastructure Network. A project that uses a token to crowdsource physical hardware, here GPUs, from many owners into a paid market. Akash, io.net, Gensyn, and Prime Intellect's marketplace are examples. DII is deliberately not a DePIN: it has no token and nothing for sale, so the supply comes from mission rather than payment and access is never bought. See research/Prior-Art.

**TOPLOC**: Prime Intellect's method for verifying that an untrusted node ran the inference it claimed. Locality-sensitive hashing over a model's activations detects a swapped model, prompt, or precision, cheaply (258 bytes per 32 tokens) and quickly (up to 100x faster than the original inference). The kind of borrowed check DII would reserve for the rare high-stakes cross-boundary case rather than run inside a trusted pod. See research/Prior-Art/04-Prime-Intellect.md.

**OpenDiLoCo**: Prime Intellect's open implementation of DeepMind's DiLoCo (Distributed Low-Communication) training. Workers take many local optimizer steps and only sync periodically, cutting inter-node communication by roughly 500x so training can run over ordinary internet links. A training technique; DII serves existing open weights and does not train, so this is reference context, not a DII component. See research/Prior-Art/04-Prime-Intellect.md.

**H100**: NVIDIA's datacenter GPU, the reference point for frontier-scale training and serving, with about 3.35 TB/s of memory bandwidth. It is the centralized baseline the distributed projects measure against, and the contrast that defines DII's envelope: DII targets consumer cards (roughly a 4090 or 5090, at a quarter to a third of that bandwidth) and a reliable floor that fits on them, not H100-class frontier inference.
