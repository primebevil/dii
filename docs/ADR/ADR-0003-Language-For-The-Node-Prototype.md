# ADR-0003: Language for the Node Prototype

Status: Accepted, 2026-07-03.

## Context

The Week-3 proof of concept (see the 30-Day Research Plan and the Architecture Overview build order) needs a node daemon: a long-lived process that wraps a local open-weight model, exposes an HTTP API, publishes a capability manifest, and routes requests to pod peers. It must ship as something a non-expert can run on Linux, macOS, and Windows, because the prototype's job is partly to recruit a pod, and a person who cannot install it will not join one.

The prototype is explicitly a measurement artifact. The 30-Day Research Plan frames Week 3 as producing the kill-criteria numbers, and the Charter accepts that no further system may emerge. So the code may be thrown away, and the decision should optimize for learning and reversal, not for production durability.

Two candidate languages were considered seriously: Go and Rust. A colleague recommended Rust. PHP, the project owner's home language, and Node were ruled out earlier on distribution grounds, since "install a runtime and its dependencies" is where a volunteer pod dies before it forms; both Go and Rust instead compile to a single static binary with no runtime.

A fact specific to how this project is built shapes the decision. The work is AI-driven with the human in the loop: the AI writes the code and the owner reviews it. That moves the scarce resource from human authoring effort to human review bandwidth.

## Decision

Use Go for the prototype. Treat the choice as reversible, and revisit Rust if and when the prototype survives the kill-criteria and moves toward a production daemon.

## Rationale

The forces separate cleanly once the AI-driven model is taken into account.

Distribution is a tie. Both Go and Rust produce a single cross-platform static binary, so both clear the bar that eliminated PHP and Node.

Control-plane performance does not matter here. The bottleneck is model inference and a residential network, which are orders of magnitude slower than anything the router does, so Rust's runtime-speed advantage is invisible in this prototype. Naming this matters, because raw performance is the usual reason to reach for Rust and it does not apply.

Reviewability is the operative criterion. Because the AI writes the code, human learning curve and human iteration speed stop being the owner's concern, and the real constraint becomes how quickly the owner can review a change. Go is more readable for someone from a backend PHP background than Rust is, so it optimizes for the resource that is actually scarce in this workflow.

The honest argument for Rust is that its compiler acts as a second reviewer, catching classes of bug, such as memory-safety errors and data races, that a human skim will miss. That argument is weak for a throwaway prototype, where those bugs are rare and low-stakes, and it becomes strong exactly at the point where the daemon starts accepting untrusted cross-pod traffic. That point is the prototype-to-production boundary, which is also where a rebuild is already on the table.

The colleague's ability to help review Rust is real but not decisive, since the owner knows code and can carry the review himself.

## Consequences

The prototype is written in Go, wrapping a local model runtime (Ollama or llama.cpp) rather than implementing inference, and exposing an OpenAI-compatible HTTP endpoint so existing clients work against it without integration.

The decision is low-regret because the node-to-node protocol is HTTP and JSON, which is language-agnostic. A future Rust node can speak the same protocol and federate with Go nodes, so switching does not require a coordinated rewrite of the whole system, only of a node. Under the AI-driven model, regenerating a small daemon is cheap.

Rust remains the leading candidate for the eventual production daemon, on memory-safety grounds, and this record should be revisited at the prototype-to-production boundary rather than treated as settled for all time.

## Open questions

Whether the production node daemon is written in Rust, kept in Go, or split, is deferred to the point where the prototype has passed the kill-criteria. Whether any part of the control plane benefits from the owner's PHP experience, for instance a management or dashboard layer that is not the daemon, is also left for later.
