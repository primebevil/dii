# Distributed Intelligence Initiative
## Project Charter & Research Prospectus

> Version 0.2 — 2026-07-02 (non-commercial resilience pivot)

## Executive Summary

The Distributed Intelligence Initiative (DII) is a **free, open-source research
initiative** (RFC-style) exploring how to keep a dependable baseline of capable,
open-weight AI available to everyone — running on hardware people own, and
impossible to switch off from above.

The goal is **not** to compete with frontier model providers, and **not** to
build a business. It is to guarantee a **reliable floor** of useful intelligence:
not the most advanced model in existence, but an advanced-enough one people can
always rely on. See [Why DII Exists](./Why_DII_Exists.md).

---

# Problem Statement

Access to the most capable AI is centralized, and therefore revocable. In June
2026 a frontier model was disabled worldwide by government directive within hours,
then restored two weeks later — proving that centralized access is something users
are *permitted to use for now*, not something they *have*.

The central research question is:

> Can participant-owned hardware provide a resilient, uncontrollable floor of
> capable open-weight AI that keeps working when centralized access is
> restricted, repriced, or switched off?

---

# Vision

Users request **capabilities**, not hardware — reasoning, coding, vision, speech,
retrieval. A router determines where work executes based on capability, trust,
sovereignty/privacy, latency, and availability, always preferring the user's own
machine first.

The measure of success is **continuity**, not frontier parity.

---

# Working Hypotheses

- People are the infrastructure.
- Intelligence should be requested as a capability.
- Local-first, network-as-overflow is the resilient shape.
- A federation of trusted **pods** beats a global swarm of strangers.
- Contribution should earn **access** (reciprocity), not credits or currency.
- A reliable floor for everyone beats an intermittent ceiling for a few.

These are hypotheses, not conclusions.

---

# Research Principles

- Research before implementation.
- Credit all prior work.
- Separate facts from hypotheses.
- Record architecture decisions (ADRs) and explain *why*.
- Validate through prototypes and honest kill-criteria.
- Let evidence change conclusions.
- Treat the repository as a living engineering knowledge base.

---

# Repository Philosophy

This repository should become a research library, an engineering notebook, an
architecture handbook, a prototype log, and a historical record of how the ideas
evolve. Superseded ideas are marked, not deleted. Every paper should be
summarized; every important decision should explain *why*.

---

# Initial Research Areas

- Distributed and split inference
- Local-first / on-device AI
- Volunteer and federated compute
- Capability routing and scheduling
- Trust, identity, and reputation
- Privacy-preserving and sovereign execution
- Fault tolerance and graceful degradation
- Federation protocols (discovery without a central hub)

---

# 30-Day Exploration Plan

**Week 1 — Define** the proposition (done).
**Week 2 — Research** existing work; teardown of prior art (in progress).
**Week 3 — Prototype** a two-node, one-pod proof of concept.
**Week 4 — Pressure test** against kill-criteria and revise the architecture.

See [30-Day Research Plan](./30_Day_Research_Plan.md).

---

# Architecture Themes

Capability routing · local-first execution · federated pod topology · node/pod
discovery without a central hub · trust and reputation · privacy-preserving
execution · fault tolerance and graceful degradation.

See [Architecture Overview](../architecture/Overview.md).

---

# On Sustainability (Non-Commercial)

DII is non-commercial. There is no business model, token, tradable credit, or paid
marketplace in scope, and nothing in the architecture depends on one. If a
*sustainable* funding path (e.g. optional managed federation for organizations) is
ever worth revisiting, it will be evaluated separately, on its own merits, and
never at the expense of the free reliable-floor mission. The earlier commercial
framing is retired — see [Startup Concept (superseded)](./Startup_Concept.md).

---

# What This Is Not

- Not a business, and not monetized.
- Not a GPU rental marketplace.
- Not a cryptocurrency, token, or tradable-credit project.
- Not an attempt to serve frontier, dual-use capability — a floor, not the ceiling.
- Not dependent on any hub, registry, or provider that could be switched off.
- Not a race to build frontier models.

---

# Success Criteria

Success means developing enough understanding to make informed engineering
decisions — and, ideally, a working local-first node plus a two-node pod that
demonstrate the reliable floor is real.

Even if no further system emerges, the project should stand as a thoughtful,
honest contribution to resilient, sovereign AI infrastructure.

---

# Long-Term Philosophy

Treat this as a field of study. Preserve every experiment. Credit every
influence. Question every assumption. Build carefully. Learn continuously.
