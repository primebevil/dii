# Week 3 — M4 findings: overflow measurement

Status: Findings, 2026-07-12. Closes M4 of the Week-3 prototype plan
(docs/Week_3_Prototype_Plan.md). Records the kill-criteria numbers and whether
each passed. Companion to docs/Identity_Note_From_Prototype.md.

## Setup

A 3-node pod, measured live (more than the plan's two-machine minimum):

- **laptop** (hub) — macOS, reachable on the LAN and on Tailscale; serves `llama3.2:1b`.
- **atlas** — Linux server, reached over **Tailscale**; serves `qwen2.5:7b`.
- **sirius** — Linux laptop, reached over the **LAN**; serves `qwen2.5-coder:1.5b`.

All requests enter the laptop hub. `atlas` and `sirius` are on different
networks and can't reach each other; the laptop bridges them (see DEPLOY.md).

Measured with `cmd/harness`: a fixed embedded prompt, fixed `max_tokens: 64`,
`temperature: 0`, 2 warmup runs (discarded) then 5 measured, medians reported.
Time to first token is measured to the first *content* token. Each peer's local
baseline was measured **on that peer** (loopback), so the overflow overhead
includes the real network hop, not just hub routing.

Models are the node-entry class (1.5B/7B) that runs on the hardware on hand, not
the 30B reliable floor (ADR-0006). See "Generalization" below for why the
conclusion still carries.

## Numbers (median of n=5, 64-token output, warm)

| Path | Measured on | Model | TTFT | tok/s | Total |
|------|-------------|-------|------|-------|-------|
| laptop-local            | laptop              | llama3.2:1b        | 158 ms | 81.2 | 947 ms |
| atlas-local (baseline)  | on atlas            | qwen2.5:7b         | 126 ms | 44.8 | 1555 ms |
| atlas-overflow          | laptop→atlas (TS)   | qwen2.5:7b         | 146 ms | 44.7 | 1579 ms |
| atlas-consumer          | laptop→atlas (TS)   | qwen2.5:7b         | 147 ms | 44.8 | 1571 ms |
| sirius-local (baseline) | on sirius           | qwen2.5-coder:1.5b | 147 ms | 34.2 | 2022 ms |
| sirius-overflow         | laptop→sirius (LAN) | qwen2.5-coder:1.5b | 191 ms | 34.8 | 2068 ms |
| sirius-consumer         | laptop→sirius (LAN) | qwen2.5-coder:1.5b | 184 ms | 35.0 | 2014 ms |

Derived overhead (overflow minus the peer's own local):

| Path | TTFT overhead | Throughput vs local |
|------|---------------|---------------------|
| atlas-overflow  | +20 ms | 99.8% |
| atlas-consumer  | +21 ms | 100% |
| sirius-overflow | +43 ms | 102% |
| sirius-consumer | +37 ms | 102% |

## Kill criteria

1. **Overflow throughput ≥ ~90% of the peer's local.** PASS — Atlas 99.8%,
   sirius 102%. Streaming rate is effectively unchanged by the hop, as expected:
   once tokens flow, the wire cost is negligible.
2. **Overflow adds ≤ ~a couple hundred ms to TTFT.** PASS — +20 ms (Tailscale),
   +43 ms (LAN). Comfortably interactive.
3. **Consumer path within the overflow envelope.** PASS — consumer and overflow
   are within a few ms of each other on both nodes, confirming they are the same
   path entered one stage in.
4. **Honest, bounded degradation, never a hang.** PASS — a model no node has
   returns HTTP 503 immediately (the router refuses before any upstream call).

All four pass. The technical overflow thesis is validated on real hardware over
two different network types.

## Findings

- **Overflow is nearly free.** Borrowing a peer's model costs ~0% throughput and
  tens of milliseconds of first-token latency. The dominant cost of a request is
  the inference itself, which happens on the peer either way.
- **The overhead is model-independent.** It is transport plus hub routing;
  neither scales with model size. A 30B model would raise absolute TTFT and lower
  tok/s on *both* the local and overflow sides, so the *overhead* — the number
  under test — would be the same tens of ms. Measuring with small models
  therefore validates the transport for the reliable-floor class too.
- **Tailscale (+20 ms) beat bare LAN (+43 ms).** Not a network effect: sirius is
  a weaker machine adding a little more per-request node latency. Both are
  trivial; the point is that an encrypted mesh hop across the internet was no
  worse than a LAN hop.
- **Residential reachability was a non-issue.** The plan's #1 risk was two home
  machines failing to connect through NAT. Tailscale dissolved it — nodes reach
  each other by stable mesh IPs with no port-forwarding. Finding: a mesh VPN is
  the pragmatic answer to residential-to-residential, and it costs nothing
  measurable in overhead.

## Limitations (so the numbers aren't over-read)

- Small sample (n=5), one session, warm models, short generations. The margins
  are large enough that the verdicts are safe, but these are not rigorous
  benchmarks.
- Node-entry-class models, not the 30B floor (hardware on hand). Overhead
  generalizes (above); absolute latency/throughput would differ.
- One physical LAN and one tailnet, both low-latency. A worse residential uplink
  (higher RTT, more jitter) would raise the TTFT overhead — but the budget is
  ~200 ms and we spent ~20–43 ms, so there is substantial headroom.

## What this does and does not settle

It settles the **technical** claim: a node can transparently borrow capability
from a peer pod, over real links, at negligible cost. It does not settle the
**durability** claim (is aggregation still worth it as local models improve) —
that is a market/time question, argued in docs/Pod_Aggregation_Red_Team.md, and
no latency measurement can close it. M4's job was the former, and it passed.
