# Why DII Exists: The Fable 5 Episode and the Case for a Reliable Floor

Status: Founding rationale, v0.1 (2026-07-02). Companion to the Project Charter and The Case Against DII.

## The event that made this concrete

On June 12, 2026, Anthropic disabled its two most capable models, Claude Fable 5 and Claude Mythos 5, for every customer worldwide, roughly four days after releasing Fable 5. It did so to comply with an emergency US Commerce Department export-control directive issued that evening, following a June 2, 2026 executive order titled "Promoting Advanced Artificial Intelligence Innovation and Security." The directive barred access by any foreign national, whether outside the US, inside the US, or working for Anthropic as a foreign-national employee. Because the company could not verify every requester's nationality in real time, it could not comply selectively, so it shut the models off for everyone. Access was restored roughly two weeks later, around June 26 to July 1, 2026.

This already happened. A frontier model was switched off for the entire planet, by directive, in a matter of hours, and then switched back on.

## What the episode proves

It is tempting to read this as "the government took AI away from the public." That framing is emotionally true but strategically imprecise. The sharper and more useful lesson is that a centralized model provider is a single point of control: access to it can be revoked globally, in hours, by a party that is neither the user nor the provider.

The fact that access came back strengthens the point rather than weakening it. The episode showed that frontier access is now revocable at will. Something that can be switched off in an afternoon and switched on again two weeks later is a permission granted for now, which is a shaky foundation to build a life, a business, or a country's capability on. DII exists because "permitted for now" should not be the only option. Weights running on hardware you own cannot be recalled by a directive, and that is the whole thesis. The Fable 5 episode is its founding case study.

## The concession we make honestly

We accept that some control is legitimate. The trigger in this case was real: a reported method to bypass Fable 5's safeguards in a way that could help identify vulnerabilities in critical infrastructure, alongside documented efforts by foreign labs to distill the models' capabilities. Some AI capabilities genuinely resemble dangerous dual-use technology, and "give the most powerful possible system to absolutely everyone, no questions asked" is a poor position to defend. We will not stake our project on it. DII therefore makes no demand that frontier, weapons-relevant capability be distributed to all. That is the strongest objection to our value, and we concede it up front rather than fight on that ground.

## The concern that motivates us

This next part is a pattern-based inference, and we label it as a hypothesis in keeping with the project's principles. The worry is about what tends to follow the first control. A restriction introduced for a specific, defensible reason such as "this exploit is dangerous" establishes both the mechanism and the precedent for control. Once the mechanism exists, meaning a switch that one party can throw, the set of reasons for throwing it tends to widen over time, from acute safety to strategic advantage, then to controlling who may use the technology and for what, and eventually to shaping which uses and narratives are permitted at all. History with other powerful, centrally gated technologies suggests the ratchet more often tightens than loosens. We may be wrong about the trajectory, but a resilient design costs little if the concern is overblown and matters enormously if it holds. That asymmetry is why we build.

## What DII actually guarantees

The promise is deliberately modest, and that modesty is what makes it credible. DII does not promise the most advanced model in existence. It aims to guarantee that you always retain access to an advanced-enough model you can rely on, one that runs on hardware people own and stays beyond the reach of any single off-switch.

A reliable floor of capability, permanently available, serves people better than intermittent access to the ceiling. Open-weight models that are already broadly published, the kind the June 2026 episode did not target, are more than capable enough to keep individuals, businesses, researchers, and communities functioning. The goal is continuity of useful intelligence. This is "AI for everyone" in the one form that survives the strongest objection: everyday capability that stays switched on, like electricity.

## What this is not

To be clear about scope, DII is not a claim that all AI governance is illegitimate, not a demand for unrestricted access to frontier dual-use capability, not an attempt to evade export controls or safety law, and not dependent on beating centralized providers on cost or at the frontier. It is a resilience project. The aim is to keep a dependable floor of open, self-hostable intelligence available so that ordinary capability cannot be revoked from above.

## Why this shapes the architecture

The positioning dictates the build order. Local execution is the foundation for anyone who can self-host, because if the value is that it cannot be switched off, capable AI has to run on hardware someone controls. The guarantee itself, though, is access to a capable model, and for people who cannot self-host the network delivers it through a pod that serves them, since an exclusive network would strand exactly the people with no fallback when the cloud goes dark (see ADR-0002). The distributed network is therefore both the overflow layer for nodes and the delivery path for consumers. Open weights are the substrate, since the guarantee only holds for models that can be legally held and run by their users, which also keeps us clear of the proliferation objection above. The target workloads are latency-tolerant and good-enough, the same envelope the physics of home hardware can actually serve, as discussed in The Case Against DII under Objection 1. And sovereignty and data-residency are first-class routing inputs, decided at the start of the design rather than added later, because that is the demand the evidence leaves standing.

## Sources

- [Statement on the US government directive to suspend access to Fable 5 and Mythos 5, Anthropic](https://www.anthropic.com/news/fable-mythos-access)
- [Anthropic Pulls Its Most Powerful AI Models After U.S. Bars Foreign Access, Time](https://time.com/article/2026/06/13/anthropic-fable-mythos-ban-US-security/)
- [Anthropic Disabled Fable 5 And Mythos 5 After A U.S. Export-Control Order, Forbes](https://www.forbes.com/sites/anishasircar/2026/06/16/anthropic-disabled-fable-5-and-mythos-5-after-a-us-export-control-order-heres-what-happened/)
- [US lifts restrictions on Anthropic's powerful AI models Fable and Mythos, Al Jazeera](https://www.aljazeera.com/economy/2026/7/1/us-lifts-restrictions-on-powerful-ai-models-fable-mythos-anthropic-says)
- [License to model: Emerging US rules impact global access to frontier AI, HSF Kramer](https://www.hsfkramer.com/insights/2026-07/license-to-model-emerging-us-rules-impact-global-access-to-frontier-ai)
