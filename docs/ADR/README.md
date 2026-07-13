# Architecture Decision Records

One file per significant decision.

- ADR-0001: Research before implementation.
- ADR-0002: Access is not gated on contribution; the network serves non-contributing consumers.
- ADR-0003: Language for the node prototype; Go now, Rust reversible at the production boundary.
- ADR-0004: The consumer is a first-class ingress in the prototype; one router, two ingress types.
- ADR-0005: Motivation is mission, not money; reciprocity demoted, monetization retired (AI for All, in the spirit of Tor).
- ADR-0006: What the reliable floor is; the 14B-to-30B open-weight class at Q4, defined by a usefulness line, not just hardware.
- ADR-0007: Who DII is for; dependency-defined, not "everyone"; broad mission, narrow recruiting wedge, personal pod-zero.
- ADR-0008: Funding the stewards; the network stays non-commercial, but funding into a nonprofit steward and cost-offset for operators to break-even is accepted as distinct from monetization.
- ADR-0009: Public and private pods, and funding eligibility; rivalrous cost caps pod size, the unaffiliated are served by many public-serving pods (Tor exit-relay model), and funding follows mission (public-serving first, public-interest case-by-case, private self-funded).
- ADR-0010: Voluntary sponsorship, not paid private use; private pods are invited (not required) to sponsor public access, and a required fee or license for private/commercial use is reserved, not adopted.
- ADR-0011: Inter-node transport is the reused OpenAI-compatible HTTP call; a peer is just another backend behind the interface, validated in the Week-3 prototype as effectively free, with a thin internal RPC reserved for when identity or policy must travel across the hop.
- ADR-0012: The backend is any OpenAI-compatible model server behind a thin interface; Ollama is one implementation selected by a URL, so the node core never learns which backend runs and the no-lock-in commitment holds in code.
- ADR-0013: The pod is the accountability boundary; admission, moderation, fair-use, and revocation live with the pod operator, never a central authority, so no single actor can switch the network off.
- ADR-0014: Consumer work is preemptible best-effort, bounded by per-identity fair-use, and the guarantee is floor-access to a capable model, not unlimited throughput.
