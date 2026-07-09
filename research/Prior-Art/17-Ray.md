# Ray

## What it does

Ray is an open-source distributed-computing framework from UC Berkeley's RISELab, commercialized by Anyscale, that scales Python and AI workloads from a laptop to a cluster of thousands of nodes. Its core abstractions are tasks, stateful actors, and a distributed object store, with libraries on top for training, serving, tuning, RL, and data. It does resource- and capability-aware scheduling against typed resources (CPUs, GPUs, memory, and custom labels), with autoscaling and lineage-based fault tolerance, and it is heavily adopted, joining the PyTorch Foundation in October 2025 and used by OpenAI, Uber, and others. It is here as the direct test of the teardown's Objection 5, whether a scheduler is a defensible product.

## Verdict

Won its niche, which is exactly the point. Ray dominates capability-aware scheduling inside a trusted cluster, so DII must not position itself as a better cluster scheduler.

## Why

Ray has thoroughly solved and commoditized capability-aware placement within a cluster. It schedules against typed and custom resources, supports gang and topology-aware placement, claims millions of tasks per second, recovers from node failure, and has a deep Kubernetes ecosystem through KubeRay, Kueue, Volcano, and the NVIDIA GPU operator. That is the home turf, and it has won it.

The decisive fact for DII is what Ray assumes, stated plainly in its own security documentation: Ray expects to run in a safe network environment and to act upon trusted code, and security and isolation must be enforced outside the Ray cluster. In other words, Ray assumes one administrative domain, a trusted and isolated low-latency network, and trusted code, and it does not differentiate between a legitimate job and a rootkit. This is not academic: the ShadowRay campaign disclosed in 2024 exploited thousands of internet-exposed Ray clusters precisely because Ray has no built-in trust boundary. Ray's model breaks the moment it leaves a single trusted domain.

## How it compares to DII

Ray is the reason the teardown's Objection 5 is half right, and it clarifies exactly which half. The half it gets right is that capability-aware placement inside a trusted cluster is commoditized, so if DII pitched itself as smarter GPU matching in a datacenter it would lose to Ray outright. The half it gets wrong is the conclusion that DII therefore has no product. Ray, by its own stated assumptions, does nothing for the regime that is DII's entire reason to exist: verifying an untrusted node's claimed capabilities, placing work under sovereignty and jurisdiction rules, routing privately to consumer machines, and tolerating residential-network churn and latency across many owners. Those are the constraints Ray assumes away by requiring a trusted, isolated network.

So DII's defensible framing writes itself. DII's router is not a better scheduler; it is a trust, sovereignty, and federation layer that Ray explicitly does not attempt. The scheduling mechanics are table stakes, and the constraints the scheduler runs under are the product. Where a DII node happens to be a real cluster, DII should sit above Ray or KubeRay and delegate to it locally, not try to replace it. Ray's dominance is a positioning gift, not a threat, as long as DII does not miscast itself as a cluster scheduler.

## Borrow / Avoid

Borrow: the task and actor vocabulary as a clean mental model; resource-typed scheduling as the right primitive to extend, with DII's twist that capabilities are attested and cross-owner rather than self-declared inside a trusted cluster; and Ray Serve's OpenAI-compatible serving patterns, which DII should interoperate with rather than reinvent.

Avoid: trying to out-Ray Ray inside a cluster, a losing position against a dominant, foundation-governed framework; and, most of all, letting DII be described as a better scheduler at all, when its value is the trust and sovereignty layer Ray disclaims.

## Sources

- [PyTorch Foundation welcomes Ray, Oct 2025](https://pytorch.org/blog/pytorch-foundation-welcomes-ray-to-deliver-a-unified-open-source-ai-compute-stack/)
- [Ray security documentation (trust-domain assumptions)](https://docs.ray.io/en/latest/ray-security/index.html)
- [ShadowRay, Oligo Security](https://www.oligo.security/blog/shadowray-attack-ai-workloads-actively-exploited-in-the-wild)
- [Ray Serve documentation](https://docs.ray.io/en/latest/serve/index.html)

Note: the claim that OpenAI used Ray for large-model training is widely reported but not confirmed from a primary OpenAI statement naming specific models. The security quotes are verbatim from Ray's own docs.
