# DII Architecture Overview

## Local-First, Network-as-Overflow

> Status: RFC v0.1 — 2026-07-02
> Type: Architecture overview
> Depends on: [Why DII Exists](../docs/Why_DII_Exists.md), [The Case Against DII](../research/The-Case-Against-DII.md)

---

## Framing note: what changed

DII is now scoped as a **free, open-source RFC project** with no commercial model
behind it. Its goal is a dependable baseline of open-model intelligence that
anyone can run and that stays beyond the reach of any single off-switch.
Monetization is out of scope for the foreseeable future; if a *sustainable*
funding path appears later it can be reviewed on its own merits, and no part of
this architecture assumes or requires one.

This reframing removes a whole problem class. There is no token, tradable credit,
marketplace to clear, or compute to price. "The scheduler is the product" now
means only that the router/scheduler is the core engineering contribution.

---

## The one-sentence architecture

> Every participant runs a **self-sufficient local node** that serves capable
> open-weight AI offline; nodes optionally federate into **trusted pods**; and a
> **capability router** spills work outward — local first, then pod, then
> federation — only when needed and only within policy the user controls.

The network is an enhancement layer. Remove it entirely and every node still
works. That property is the whole point of the project, so it sits at the top of
every design choice below.

---

## Design principles

1. **Local-first, always.** A node with zero connectivity still delivers useful
   AI. The network amplifies a node that already works on its own.
2. **Graceful degradation, never zero.** Network down → local. Pod down →
   federation. Peer fails mid-request → reroute. The floor is always "still works
   locally."
3. **No single point of control.** No global registry, no mandatory hub, nothing
   that a single directive could switch off. Discovery and identity are
   federated, like email or Matrix.
4. **Open weights only.** The guarantee only holds for models users may legally
   hold and run. This also keeps DII clear of the frontier-proliferation
   objection.
5. **Trust before proof.** Prefer running work inside a trust boundary over
   proving correctness on strangers' hardware. Heavy verification is a last
   resort (see Trust, below).
6. **Sovereignty is a first-class routing input.** "Never leaves my machines" or
   "stays in-country" are constraints the router must honor from the start.

---

## Topology: a federation of trusted pods (key decision)

The instinct behind the original vision — "millions of strangers' machines become
one computer" — is the version the red-team teardown hit hardest on trust and
sustainability. This architecture makes a different bet:

> **The unit is the *pod*: a small group of nodes that already trust each
> other** — your own machines, a lab, a hackerspace, a company, a family, a
> town's library co-op. Pods are self-sufficient. Pods may *federate* with other
> pods to share overflow capacity, under explicit policy.

Why this shape:

- **It defuses the trust tax.** Most work runs inside a boundary where the nodes
  are already trusted, so expensive verification stays unnecessary for the common
  case.
- **It makes participation durable.** Reciprocity within a real community is a
  stronger, more human incentive than altruism toward anonymous strangers — the
  failure mode that plateaued volunteer computing.
- **It has no chokepoint.** A federation of independent pods has no switch to
  throw. Kill any pod and the rest keep running.

Closer to the Fediverse for inference than to SETI@home for LLMs.

---

## Components

**1. Node runtime (the atom).**
Runs one or more open-weight models locally, exposes them as capabilities, and
functions fully offline. This is the star of the whole system: a person who only
ever runs a single node has already received the core promise of DII. Everything
else is optional scaling on top of this.

**2. Capability abstraction.**
Requests target *capabilities* — reasoning, coding, vision, speech, retrieval,
embedding — at a required quality/latency tier. The caller names the capability
and tier, and the router selects the hardware and model. Each node publishes a
**capability manifest**: which models it serves, context length, throughput,
modalities, and current load. This lets a heterogeneous fleet look uniform to a
caller.

**3. Discovery (federated).**
Within a pod: a simple local registry / LAN discovery. Across pods: signed,
gossiped peer lists and a DNS-like federated directory, deliberately avoiding a
single global index — that would reintroduce the chokepoint the project exists to
eliminate.

**4. Router / scheduler (the coordination core).**
Decides where a request executes and enforces the overflow order: **local → pod
→ federation.** Routing inputs: capability match, sovereignty/policy constraints,
trust level of candidate nodes, latency, availability, and load. Each hop outward
passes a policy gate — work only spills across a boundary the user has authorized.

**5. Trust & reputation.**
With no payment, trust rests on **identity** (signed, persistent node identity),
**reputation** (track record within and across pods), and **policy** (who you'll
accept work from, who you'll send work to, where data may travel). Verification is
tiered: negligible inside a trusted pod; optional for cross-pod work; and for the
rare high-stakes cross-boundary case, borrow an existing protocol (TOPLOC/Verde-
style) rather than inventing one. Work inside a trusted pod skips the redundant-
execution tax.

**6. Participation model (reciprocity).**
Contribution earns **access.** Run a node and contribute idle capacity, and you
gain priority when *you* need to spill work onto the federation. This access is
deliberately **non-transferable, non-tradable, and expiring** — a fair-use signal
that keeps DII clear of the crypto category. Pure voluntary contribution is
welcome too; the reciprocity layer just makes sustained participation rational.

**7. Privacy-preserving execution.**
Because sovereignty is the reason-to-exist, the default keeps data inside the
user's trust boundary. Cross-boundary execution is opt-in and policy-gated.
Techniques for protecting data on semi-trusted hops (and whether to allow them at
all) are an open design area, flagged here for later.

**8. Fault tolerance.**
The router treats failure as normal: peers vanish, networks partition, pods go
dark. Every failure mode resolves *downward* toward local execution.

---

## What runs where (the overflow path)

```
Request for a capability
   │
   ▼
[1] Local node can serve it?  ──yes──►  run locally (always preferred)
   │ no / over capacity
   ▼
[2] A trusted pod peer can?   ──yes──►  run in-pod (policy gate)
   │ no
   ▼
[3] A federated pod can?      ──yes──►  run cross-pod (policy + trust gate)
   │ no
   ▼
Degrade gracefully: queue, use a smaller local model, or report honestly.
```

The final step always returns the best the local node can do, which is the
resilience guarantee made concrete.

---

## Build order

- **Phase 0 — The atom.** A solid local-first node runtime: capable open model,
  offline-capable, capability manifest. This alone delivers the core promise to a
  single user and is worth shipping on its own.
- **Phase 1 — Two nodes, one pod.** The Week-3 proof of concept: node A routes a
  capability request to node B and back. Measure the numbers the teardown demands
  (tokens/sec over a residential link vs. local; overflow overhead).
- **Phase 2 — Pod overflow + policy.** Router honors sovereignty/trust policy;
  reciprocity accounting; graceful degradation paths.
- **Phase 3 — Federation.** Signed discovery across pods; cross-pod trust and
  optional verification; the "Fediverse for inference" topology.

---

## Open decisions (to record as ADRs)

- **Reciprocity vs. pure-voluntary** as the default participation model. *(Leaning
  reciprocity, illiquid — see Component 6.)*
- **Discovery protocol** — gossip vs. DNS-like federation vs. Matrix-style
  servers.
- **Cross-boundary verification** — when (if ever) redundant execution or a
  borrowed proof is worth its cost.
- **Data protection on untrusted hops** — allowed at all, and if so how.
- **Node identity** — self-sovereign keys vs. pod-issued identity.

---

## What this deliberately excludes

- A global marketplace or a swarm of anonymous strangers.
- Any token, credit-currency, or compute-for-pay scheme.
- Serving frontier, dual-use capability — DII targets a dependable floor.
- Any hub, registry, or provider that could be switched off.
