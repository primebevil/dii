# Pod Aggregation: A Red Team Against the Thesis

Status: Working notes, 2026-07-10. An adversarial pass, not an ADR. Depends on ADR-0006 (reliable floor), ADR-0009 (public and private pods), ADR-0002 (access not gated on contribution), and the Governance_And_Abuse_Resistance notes. The job here is to argue, as hard as possible, that pod aggregation is redundant within roughly eighteen months, and then say what survives the argument and what would actually break. Written because the worry is real: if everyone has a capable model, a project whose pitch is "help people run a model" has no reason to exist.

## What is actually under test

Two claims travel together in casual conversation and have to be pried apart, because they age very differently. The first is the self-hosting claim: that DII's value is helping people run inference they could not otherwise run. The second is the aggregation claim: that pooling idle capacity across a pod gives a member access to something no single member could reach alone. The self-hosting claim is the one at risk of redundancy, and it should be conceded early. If open-weight models keep improving and a strong one runs on a single consumer GPU, then "we help you self-host" is a shrinking value proposition with a clock on it. That is fine, because self-hosting was always the entry ramp and never the thesis. The load-bearing claim is aggregation, and the rest of this note attacks it directly.

## Splitting aggregation into two, because only one is fragile

Pooling buys two different goods that are usually said in one breath: aggregation for size, reaching a larger or more capable model than a member's hardware can hold, and aggregation for variety, reaching many models at once for diversity, ensembling, and creative range. These are not the same bet and they do not face the same threat. Size is the fragile one. Variety is the durable one. Keeping them separate is what stops the red team from either over-claiming or over-conceding.

## The steelman: why pods are redundant in eighteen months

Take the strongest form of the redundancy case, granting every trend its most favorable reading.

Frontier capability keeps getting cheaper to rent, and the vendors close the one gap that mattered to sovereignty-minded buyers. Zero-retention terms, in-VPC and on-prem deployment, and contractual non-cutoff become normal enterprise offerings rather than exotic ones. Once a member can rent frontier capability under terms that are good enough on data handling and continuity, the sovereignty premium that justified the pod thins out, because most people do not want sovereignty as a principle, they want it as insurance, and a contract is cheaper insurance than a pod.

Local models close the gap for the median request. A single strong open-weight model on one consumer card handles the everyday distribution of work, so the "reach a bigger model" motive fades for most members most of the time. The size claim erodes from below.

Aggregation is structurally worse at serving than a datacenter, and that gap widens rather than closes. Modern serving economics reward concentration: high utilization, mixture-of-experts routing, disaggregated prefill and decode, big batches. A loose pod of heterogeneous, preemptible, idle GPUs is close to the worst possible substrate for exactly the large-model inference the size claim promises. So the size claim asks the network to do the one thing its architecture is least suited to.

Coordination has costs a rented endpoint does not. Latency, scheduling, trust, and the unsolved Sybil problem are all overhead the member pays for pooling. A metered API call has none of it. Convenience is a real force and it usually wins, so even a member who agrees with DII in principle reaches for the endpoint that answers in 200 milliseconds with no identity dance.

Put together, the steelman is tidy: rented frontier eats the sovereignty motive, local models eat the size motive from below, serving economics say the network is bad at the size motive from above, and convenience eats whatever is left. On this reading pods are a transitional artifact of a moment when models were expensive and hard to run, and that moment is closing.

## What survives, and why the strongest leg was never size

The steelman lands cleanly on size and barely touches the rest. Take back the ground it does not actually hold.

Aggregation for variety is not a compute argument and does not fall to serving economics. A pod that reaches many different models, families, fine-tunes, and specializations produces diversity, ensembling, and creative range that a single vendor endpoint cannot, because a vendor serves its own family and has a commercial interest in funneling you onto it. Model variety is a property of the population of models, not of any one model's size, so "models got bigger and cheaper" does not erase it. This is the leg the earlier conversation kept returning to, and it is correctly the least commoditizable one. If DII has a durable reason to exist in a world where everyone has a good local model, this is it.

The value of aggregation was never to match frontier capability, it was to provide access under conditions frontier will not sell. No per-call meter, no content gate imposed by a vendor's terms, and no dependency on a party that can change price, policy, or availability under you. This is the point ADR-0009 already makes structurally: DII's guarantee is universal floor access with no single off-switch, not the biggest model on the market. Even a genuinely good vendor contract is still a contract with a counterparty, and the whole project exists for the people and institutions for whom that counterparty risk is the thing they cannot accept. A better vendor offering raises the bar the floor has to clear; it does not remove the reason a floor beholden to no one exists.

The serving-efficiency gap is real but it is measured on the wrong axis. Pod compute is donated and idle, so its marginal cost basis is near zero. The pod does not have to be as efficient as a datacenter, it has to be cheaper to the member than a metered call, and near-zero marginal cost gives it enormous room to be worse at utilization and still win on price for the bursty, latency-tolerant tail. ADR-0006's honest ceiling already fits this: the pod promises a capable floor and best-effort access to more, not unlimited frontier-scale compute on demand.

The local-good-enough trend clarifies pods rather than killing them. If local handles the median request, the pod is no longer competing for the median, it owns the tail: the occasional larger model, the longer context, the ensemble, the creative reach beyond one card. Local floor plus pod tail is a cleaner division of labor than the current story, not a refutation of it. The thing that shrinks is self-hosting-as-value, which was already conceded.

## The honest scorecard

Size aggregation is genuinely exposed. If it is load-bearing anywhere in the pitch or the roadmap, that is the weakest leg and it should be demoted from thesis to feature, framed as a best-effort tail benefit rather than a reason the network exists. Variety aggregation and access-under-conditions are the legs that carry weight, and they are the ones the outward messaging should lead with. The recurring lesson is that every time the argument leans on capability, DII loses, and every time it leans on independence and diversity, DII holds. The sovereignty framing is not decoration here, it is the load path.

## Kill tests: what would actually falsify the thesis

Stated as concrete conditions so they can be watched rather than debated.

- A frontier vendor ships a sovereignty-acceptable offering at commodity price: verifiable zero-retention, in-boundary deployment, no content gating beyond the illegal hard lines, and a credible non-cutoff commitment. If this becomes normal, the access-under-conditions leg weakens sharply, because the pod's insurance is now available as a cheaper contract.
- Model variety stops paying. If diverse ensembles and multi-model routing show no measurable advantage over a single best model on the tasks DII members actually care about, the durable leg is not durable, and there is little left.
- Coordination overhead makes pods worse than a metered call even for the tail. If Sybil, latency, and reliability costs are bad enough that members prefer to pay and meter even when they hold the sovereignty view, the aggregation leg breaks on user experience regardless of the economics.

If any one of these comes true, the response is not to defend the whole thesis but to fall back to the legs still standing. If two come true at once, the project's reason to exist is in real question and that should be said plainly rather than managed.

## Leading indicators worth tracking

The kill tests are lagging. These are the early signals that tell you which way the kill tests are trending: the open-weight versus frontier capability gap over successive releases; whether frontier vendors move toward zero-retention, in-boundary, no-gating terms as standard rather than enterprise-only; consumer GPU memory growth against model size growth, which sets where the local floor lands; and any real measurement of whether model-diversity ensembles beat the single best model on DII's target tasks, which is the one indicator worth building an experiment around because the durable leg rests on it.

## Open questions

- Whether the variety advantage is real and measurable on DII's target workloads, or an intuition. This is the load-bearing empirical claim now and it has no evidence behind it yet; it deserves a small experiment before it is leaned on in outward messaging.
- Where size aggregation should sit once demoted: cut from the pitch entirely, or kept as an explicitly best-effort tail feature under ADR-0006's ceiling.
- How much of the access-under-conditions leg depends on the Sybil-resistant identity problem staying tractable, since coordination overhead is one of the kill tests and identity is its largest component (see Governance_And_Abuse_Resistance).
- Whether any of this changes the outward messaging priority, given that the argument consistently holds on independence and diversity and consistently loses on capability.
