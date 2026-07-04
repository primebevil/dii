# Timeline

## 2026-06-12 (context, pre-project)

- US Commerce Dept export-control directive; Anthropic disables Fable 5 and Mythos 5 worldwide. Access restored around June 26 to July 1. Founding case study in docs/Why_DII_Exists.md.

## 2026-07-02

- Project started.
- Core hypotheses established.
- Adversarial teardown written (research/The-Case-Against-DII.md).
- Founding rationale captured (docs/Why_DII_Exists.md): resilience and reliable-floor positioning.
- Decided access is not gated on contribution; the network serves non-contributing consumers (docs/ADR/ADR-0002).

## 2026-07-03

- Architecture working session, focused on the path to a Week-3 prototype (journal/2026-07-03-architecture-session.md).
- Chose Go for the node prototype, Rust reversible at the production boundary (docs/ADR/ADR-0003).
- Made the consumer a first-class ingress in the prototype: one router, two ingress types, consumer as a request past the local hop (docs/ADR/ADR-0004).
- Folded the router ingress model and updated prototype topology into the Architecture Overview.
- Settled the participation model: motivation is mission and community, not reciprocity or money; monetization retired, not paused; positioning made explicit as AI for All, in the spirit of Tor (docs/ADR/ADR-0005). Scrubbed remaining money-making hedges and reframed participation across the docs.
- Pinned a concrete reliable-floor definition: the 14B-to-30B open-weight class at Q4, defined by a usefulness cliff rather than raw hardware; node-entry ~14B on a gaming GPU, promise ~30B on committed hardware, degrading to honest queuing below the useful line (docs/ADR/ADR-0006, docs/Reliable_Floor_Definition.md).
