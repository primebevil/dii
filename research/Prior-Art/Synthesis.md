# Synthesis: How DII Differs From the Prior Art

The "so what" of the first prior-art batch, to be read after the six one-pagers. It answers a question the owner raised after reading them: these systems are all doing something adjacent to DII, yet none is what DII proposes, and they look far more complex than what DII is building. Both halves of that are true, and the reason is worth stating plainly.

## The prior art clusters into two hard problems, and DII is in neither

Every project in the batch lives in one of two problems. Split one model across many GPUs, which is Petals and Parallax. Or train one model across many sources, which is FusionAI, Prime Intellect, Gensyn, and DisTrO. DII does neither. It runs a whole model on a single node and routes a whole request to a peer that also runs a whole model. Even Parallax, the closest cousin, is built around splitting a model across devices with pipeline parallelism. DII never splits a model mid-inference and never trains anything. That single structural choice is where most of the prior art's complexity does not apply to DII at all.

## Their complexity is the price of three things DII refuses to do

The machinery in those systems is not incidental. It buys three capabilities, and each maps onto something DII has deliberately given up.

Splitting a model across machines forces pipeline parallelism, activation transfer, and a losing fight with the latency wall. That is most of what Petals and Parallax are. DII declines it by defining a reliable floor that fits on one node (ADR-0006).

Trusting strangers forces bitwise-reproducible verification or cryptographic proofs. That is Verde and TOPLOC. DII declines it by making the pod, a group that already trusts each other, the unit, and keeping heavy verification as a last resort for the rare cross-boundary case (Overview, "trust before proof").

Paying strangers forces a token and a market, with pricing, speculation, and thin demand. That is the DePIN layer under all of them. DII declines it by retiring monetization and running on mission instead (ADR-0005).

Take those three requirements away and the engineering that remains is, as the 2026-07-03 architecture session already concluded, close to mundane: a local inference server, a capability manifest, and an overflow router.

## DII is not a simpler version of these systems, it is a different bet

The honest framing is not that DII solves the same problem more simply. It stays inside the envelope where their hard problems never arise. Their complexity is the cost of running models bigger than one machine, or trusting anyone, or paying anyone. DII pays a different cost instead: a smaller guarantee, a reliable floor rather than the frontier, and a narrower supply, people who share a mission and some trust rather than an open market. That trade is the whole thesis, and the prior art is the evidence that the other side of it is where projects get expensive and stall.

## The caveat to keep in view

Because the engineering is mundane, DII's contribution is not a systems breakthrough. It is the architecture and the trust-and-motivation model, Tor for inference. That is a strength for buildability, but it is exactly the soft spot the teardown's Objection 5 names: if the router is a thin layer, why can't an incumbent absorb it? The answer DII needs to be able to give crisply is what makes pod federation, sovereignty-as-a-routing-input, and honest degradation worth more than "run a local model and let a friend hit your endpoint." The simplicity is the point. Naming the part that is not trivial is the open work.

## References

- The six one-pagers: [01-Petals](01-Petals.md), [02-Parallax](02-Parallax.md), [03-FusionAI](03-FusionAI.md), [04-Prime-Intellect](04-Prime-Intellect.md), [05-Gensyn-Verde](05-Gensyn-Verde.md), [06-DisTrO](06-DisTrO.md).
- [The Case Against DII](../The-Case-Against-DII.md), especially Objection 5.
- ADR-0002 (access before contribution), ADR-0005 (mission, not monetization), ADR-0006 (the reliable floor), and the Architecture Overview.
