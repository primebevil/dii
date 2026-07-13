# Distributed Intelligence Initiative

A free, open-source research initiative (RFC-style) exploring how to keep a dependable baseline of capable, open-weight AI available to everyone, running on hardware people own and impossible to switch off from above.

Think of it as AI for All, in the spirit of Tor: where Tor is privacy for all, sustained by volunteers who run relays, DII is intelligence for all, sustained by volunteers who run nodes. It is free and non-commercial, with no business model, token, or marketplace. The aim is a reliable floor of capable, everyday intelligence that stays available to everyone, while the frontier labs keep building the ceiling.

## Start here

For the plainest, non-technical explanation, the one-page summary is the fastest read: [DII One-Pager (PDF)](./DII_One_Pager.pdf).

If you want the proposal in more depth, two documents give you the whole concept:

- [Why DII Exists](./docs/Why_DII_Exists.md), the motivation and the founding case study, the June 2026 episode where a frontier model was switched off worldwide by directive. This is the why.
- [Architecture Overview](./architecture/Overview.md), how the system works, starting from the idea in one sentence. This is the how.

For a faster pass, the one-page [system diagram](./diagrams/DII_System_Overview.html) shows the whole shape at a glance, and the [Manifesto](./docs/Manifesto.md) is the stance in ten lines.

If you want to know whether it holds up, two more:

- [The Case Against DII](./research/The-Case-Against-DII.md), the adversarial teardown that tried to kill the idea from first principles.
- [Prior-Art Final Report](./research/Prior-Art/Final-Report.md), what eighteen comparable projects teach and why the surviving thesis is the narrow one. The capstone of the Week-2 research.

## The full map

For due diligence, the repository is organized so each layer stands on its own:

- [docs/](./docs), the proposition: the [Project Charter](./docs/Project_Charter.md) (full prospectus), [Project Vision](./docs/Project_Vision.md), [Who DII Is For](./docs/Who_DII_Is_For.md), the [Reliable Floor Definition](./docs/Reliable_Floor_Definition.md), the [FAQ](./docs/FAQ.md), and the [Glossary](./docs/Glossary.md).
- [architecture/](./architecture), how it works, starting with the [Overview](./architecture/Overview.md).
- [prototype/](./prototype), the Week-3 proof of concept: a Go pod that routes each request to whichever node can serve it and overflows to a peer when it can't. See the [build brief](./prototype/BUILD_BRIEF.md), the [deploy guide](./prototype/DEPLOY.md), and the [M4 findings](./journal/2026-07-12-week3-m4-findings.md) where all four kill-criteria passed on a live three-node pod.
- [research/](./research), the evidence: the [teardown](./research/The-Case-Against-DII.md), the [prior-art one-pagers](./research/Prior-Art) and their [Final Report](./research/Prior-Art/Final-Report.md), and the [Validation Checklist](./research/Prior-Art/Validation-Checklist.md) that maps every load-bearing figure to a source.
- [docs/ADR/](./docs/ADR), the decisions and their reasoning, ADR-0001 through ADR-0010.
- [TIMELINE.md](./TIMELINE.md), [INFLUENCES.md](./INFLUENCES.md), and [journal/](./journal), how the project got here.

## Core idea

Local-first nodes that work offline, optionally federating into trusted pods, with a capability router that spills work outward, local then pod then federation, only when needed and only within policy the user controls. People who cannot self-host join as consumers and are served by a pod that chooses to serve them; access is never gated on contribution.

## Principles

- Research before implementation
- Credit prior work
- Facts versus hypotheses
- Local-first, with no single point of control
- Build small prototypes

This repository is AI-drafted and human-directed; see [CONTRIBUTING.md](./CONTRIBUTING.md).
