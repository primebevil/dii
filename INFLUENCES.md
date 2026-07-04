# Influences

Track papers, companies, books and people that influenced the project. Each entry says what it is and how it shaped DII, credited in keeping with the project's principle of crediting prior work.

## Tor (foundational)

The single clearest influence on the project's shape. Tor is privacy for all, sustained by volunteers who run relays because they believe in it, and most people who benefit from Tor never run a relay. DII takes that whole model and points it at capable AI: intelligence for all, sustained by volunteers who run nodes, free to use whether or not you contribute. Tor is where the project gets its motivation model (mission and community, not payment; monetization retired, per ADR-0005), its access model (benefit is never gated on contribution; the many are served by the few, per ADR-0002), and its resilience thesis (a distributed, volunteer network has no single off-switch to throw).

The analogy is bounded on purpose. Tor is the model for motivation, legitimacy, and resilience, not for technical architecture; onion routing has nothing to contribute to the capability router. DII also diverges where the substrate differs: a Tor relay forwards cheap, non-rivalrous spare bandwidth, while a DII node runs expensive, rivalrous GPU inference, which is why DII leans harder on community scope and why it raises the node-entry bar to the useful line rather than accepting any contribution (ADR-0006). Tor is also a rare successful case of sustaining volunteer infrastructure for two decades through a nonprofit and institutional operators, and is worth studying as prior art for how DII pods sustain themselves.

## Federation: the Fediverse and Matrix

The model for topology and discovery without a central hub. "A federation of trusted pods" is deliberately closer to the Fediverse than to a single cloud, and cross-pod discovery borrows the signed, gossiped, DNS-like approach that email and Matrix use, avoiding a global index that would reintroduce the chokepoint the project exists to remove.

## Distributed and volunteer compute prior art

- Petals: the honest proof of concept for decentralized LLM inference on consumer GPUs, and the clearest evidence of the physics constraints DII designs around (The Case Against DII, Objection 1).
- BOINC, SETI@home, and volunteer computing broadly: the precedent for volunteer-donated compute, and a cautionary tale about how participation plateaus, which informed the mission-first, community-scoped participation model.
- Akash and io.net: DePIN compute markets, studied largely as what DII is not. Their token and marketplace dynamics, and the gap between registered supply and paying demand, motivated the non-commercial stance and the demand-side focus.

## Verification research (to borrow, not reinvent)

- Gensyn's Verde and Prime Intellect's TOPLOC: trustless and tamper-evident verification for computation on untrusted nodes. DII's stance is to borrow from this family for the rare high-stakes cross-boundary case rather than invent its own, and to avoid needing it at all inside a trusted pod.

## Local-first software

The principle that a node is fully functional offline and the network only adds to it, rather than the network being required for basic function. This underpins the local-first design commitment and the graceful-degradation guarantee.
