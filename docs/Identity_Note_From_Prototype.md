# Identity note from the prototype

Status: Note, 2026-07-12. The Week-3 plan deferred node identity and asked the
prototype to tell us what the ingress *actually needed* before we pick a
mechanism (ADR-0004 anticipated this). This captures what we observed while
building M1–M4, so the identity ADR is driven by requirements rather than a
paper comparison of DIDs vs DNSSEC vs pod-issued keys. Companion to
journal/2026-07-12-week3-m4-findings.md.

## What the ingress actually did

- **Local door (owner):** no token. Any request with no bearer is treated as the
  trusted local user and enters at the local-first stage.
- **Consumer door (stranger):** a bearer token compared for equality against one
  `consumer_token` string in config. Match → treated as a pod guest.
- **Node-to-node (peer forward):** a node calls its peer's OpenAI endpoint with
  **no token at all**, so the peer serves the forwarded request on its own
  *local* (trusted) door.

What identity had to **verify**: string equality of one shared secret, on the
consumer door only. Nothing cryptographic, nothing per-user. What it had to
**persist**: nothing — the check is stateless. That is the whole of it.

## Where the stub felt wrong (the requirements that surfaced)

1. **There is no node-to-node authentication.** Because peer forwards carry no
   identity, a node serves *any* peer that can reach it. "Pod membership" is
   currently just network reachability, not an agreed trust relationship. The
   real requirement: a node should be able to prove it is an admitted pod member
   before another node serves its forwarded request. This is the load-bearing
   gap — the consumer token gates one door on one node but does nothing to gate
   the pod.

2. **One shared consumer token can't identify or revoke anyone.** Every guest
   presents the same secret, so individuals are indistinguishable and a single
   leak compromises all access with no way to revoke one holder. The real
   requirement: per-consumer credentials that can be issued and revoked
   independently.

3. **Caller identity is lost across the hop.** When the hub forwards a guest's
   request to a serving node, the serving node sees *the hub*, not the original
   caller. Attribution, accounting, and any per-caller policy are therefore
   impossible at the node that does the work. The real requirement: the
   originating caller's identity must survive (or be vouched for) across
   forwards.

4. **Identity is the missing input to the policy seam.** The router passes
   everything through (the sovereignty/policy gate is a deliberate stub). Any
   future policy decision — who may use what, under which terms — needs an
   authenticated identity to decide on. Without it there is nothing to gate on.

## What we did NOT need (so the design doesn't over-build)

Inside a trusted pod, the loop worked correctly with zero identity. Identity is
needed at the **boundaries** — admitting a node to the pod, and admitting a
consumer at the remote ingress — not as a per-hop cryptographic tax on every
request between already-trusted members. The prototype suggests identity should
be an edge/admission concern, verified once and cached, not re-proven on the hot
path (which is where the ~20–43 ms overflow budget lives; see the findings).

## Input to the identity ADR

The three concrete needs to evaluate any mechanism against:

- **Node admission / membership:** can node B confirm that node A is a member of
  this pod before serving A's forwarded request?
- **Per-consumer credentials:** can the pod issue and revoke individual consumer
  access without a shared secret?
- **Caller attribution across hops:** can the serving node learn (or verifiably
  be told) who originated a forwarded request?

DIDs, DNSSEC, and pod-issued keys should be judged by how cleanly they satisfy
these three — with the constraint that verification belongs at admission time,
not on every request, so identity does not erode the overflow performance M4
just demonstrated.
