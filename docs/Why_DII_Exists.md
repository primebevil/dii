# Why DII Exists

## The Fable 5 Episode and the Case for a Reliable Floor

> Status: Founding rationale, v0.1 — 2026-07-02
> Type: Positioning / motivation
> Companion to: [Project Charter](./Project_Charter.md), [The Case Against DII](../research/The-Case-Against-DII.md)

---

## The event that made this concrete

On **June 12, 2026**, Anthropic disabled its two most capable models — Claude
Fable 5 and Claude Mythos 5 — for **every customer worldwide**, roughly four
days after releasing Fable 5. It did so to comply with an emergency US Commerce
Department export-control directive issued that evening, following a June 2, 2026
executive order ("Promoting Advanced Artificial Intelligence Innovation and
Security"). The directive barred access by any foreign national — outside the US,
inside the US, and including Anthropic's own foreign-national employees. Because
the company could not verify every requester's nationality in real time, it
could not comply selectively, so it shut the models off for everyone. Access was
restored roughly two weeks later, around June 26–July 1, 2026.

This is not a hypothetical about some future clampdown. A frontier model was
switched off for the entire planet, by directive, in a matter of hours — and
then switched back on.

---

## What the episode actually proves

It is tempting to read this as "the government took AI away from the public."
That framing is emotionally true but strategically imprecise. The sharper, more
useful lesson is this:

> **A centralized model provider is a single point of control. Access to it can
> be revoked — globally, in hours — by a party that is not you and not the
> provider.**

The fact that access came *back* does not weaken this; it strengthens it. The
episode demonstrated that frontier access is now **revocable at will**. What can
be switched off in an afternoon and on again two weeks later is not something you
*have* — it is something you are *permitted to use for now*. DII exists because
"permitted for now" is not a foundation anyone should have to build their life,
work, or country's capability on.

Weights running on hardware you own cannot be recalled by a directive. That is
the entire thesis, and the Fable 5 episode is its founding case study.

---

## The concession we make honestly

We do not pretend all control is illegitimate. The trigger in this case was
real: a reported method to bypass Fable 5's safeguards in a way that could help
identify vulnerabilities in critical infrastructure, alongside documented
efforts by foreign labs to distill the models' capabilities. There are AI
capabilities that genuinely resemble dangerous dual-use technology, and "give the
most powerful possible system to absolutely everyone, no questions asked" is not
a serious position. We will not stake our project on it.

So DII is **not** a demand that frontier, weapons-relevant capability be
distributed to all. That is the strongest objection to our value, and we concede
it up front rather than fight on that ground.

---

## The concern that motivates us (stated as a hypothesis)

*This is a pattern-based inference, not a proven trajectory. We label it as a
hypothesis, per the project's principles.*

The concern is not the first control. It is what tends to follow it. A
restriction introduced for a specific, defensible reason ("this exploit is
dangerous") establishes the *mechanism* and the *precedent* for control. Once the
mechanism exists — a switch that one party can throw — the set of reasons for
throwing it tends to widen over time: from acute safety, to strategic advantage,
to controlling who may use the technology and for what, to shaping which uses and
narratives are permitted at all. History with other powerful, centrally-gated
technologies suggests the ratchet more often tightens than loosens.

We may be wrong about the trajectory. But a resilient design costs little if the
concern is overblown, and matters enormously if it is not. That asymmetry is why
we build.

---

## What DII actually guarantees

The promise is deliberately modest, and that modesty is what makes it credible:

> **DII will not necessarily give you the most advanced model in existence. It
> aims to guarantee that you always retain access to an *advanced-enough* model
> you can rely on — one that runs on hardware people own and cannot be switched
> off from above.**

A reliable *floor* of capability, permanently available, beats intermittent
access to the ceiling. Open-weight models that are already broadly published —
the kind the June 2026 episode did **not** target — are more than capable enough
to keep individuals, businesses, researchers, and communities functioning. The
goal is continuity of *useful* intelligence, not parity with the frontier.

This is "AI for everyone" in the one form that survives the strongest objection:
not the bomb, but the electricity.

---

## What this is not

- Not a claim that all AI governance is illegitimate.
- Not a demand for unrestricted access to frontier, dual-use capability.
- Not an attempt to evade export controls or safety law.
- Not dependent on beating centralized providers on cost or on the frontier.

It is a resilience project: keep a dependable floor of open, self-hostable
intelligence available so that ordinary capability cannot be revoked from above.

---

## Why this shapes the architecture

The positioning is not just rhetoric; it dictates the build order:

- **Local-first is the star, not a feature.** If the value is "can't be switched
  off," step one is that capable AI runs on hardware you control, offline if
  needed. The distributed network is the *overflow and scaling* layer, not the
  foundation.
- **Open weights are the substrate.** The guarantee only holds for models that
  can be legally held and run by their users. That also keeps us clear of the
  proliferation objection above.
- **Latency-tolerant, good-enough workloads are the target envelope** — which is
  exactly what the physics of home hardware can actually serve (see
  [The Case Against DII](../research/The-Case-Against-DII.md), Objection 1).
- **Sovereignty and data-residency are first-class routing inputs**, not
  afterthoughts — the demand that the evidence leaves standing.

---

## Sources

- [Statement on the US government directive to suspend access to Fable 5 and Mythos 5 — Anthropic](https://www.anthropic.com/news/fable-mythos-access)
- [Anthropic Pulls Its Most Powerful AI Models After U.S. Bars Foreign Access — Time](https://time.com/article/2026/06/13/anthropic-fable-mythos-ban-US-security/)
- [Anthropic Disabled Fable 5 And Mythos 5 After A U.S. Export-Control Order — Forbes](https://www.forbes.com/sites/anishasircar/2026/06/16/anthropic-disabled-fable-5-and-mythos-5-after-a-us-export-control-order-heres-what-happened/)
- [US lifts restrictions on Anthropic's powerful AI models Fable and Mythos — Al Jazeera](https://www.aljazeera.com/economy/2026/7/1/us-lifts-restrictions-on-powerful-ai-models-fable-mythos-anthropic-says)
- [License to model: Emerging US rules impact global access to frontier AI — HSF Kramer](https://www.hsfkramer.com/insights/2026-07/license-to-model-emerging-us-rules-impact-global-access-to-frontier-ai)
