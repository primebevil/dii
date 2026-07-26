# ADR-0012: The Backend Is Any OpenAI-Compatible Model Server, Behind a Thin Interface

Status: Accepted, 2026-07-12. Builds on ADR-0003 and ADR-0004. Realized in the Week-3 prototype (the `modelserver.Backend` seam).

## Context

The node never runs inference; it brokers to a model server, which is Ollama in the prototype because Ollama is the easiest to run. The question this records is whether the node couples to Ollama or stays portable across inference backends. It matters on mission, not just on engineering: DII opposes lock-in and serves open weights, so marrying the node to one vendor's server would contradict the project's own thesis. The prototype also proved the boundary is cheap to keep clean.

## Decision

The node reaches its model server only through a thin `modelserver.Backend` interface, essentially two operations, list models and stream a chat completion, using the standard OpenAI subset (`/v1/models`, and `/v1/chat/completions` with streaming) and never a server's native calls. Ollama is one implementation of that interface, selected by a base URL in config. Any other OpenAI-compatible server (llama.cpp's server, vLLM, LM Studio, LocalAI) drops in by changing the URL; a backend that speaks something other than the OpenAI API gets a small adapter behind the same interface.

## Consequences

The node core stays ignorant of which backend is running, so operators are not locked to Ollama and the "no lock-in" commitment holds in code rather than only in intent. The capability manifest is built from the standard `/v1/models` call, not an Ollama-native one, which was the single place portability could have leaked and was corrected.

The cost is discipline: the prototype must use only the common OpenAI subset and must not lean on any one server's proprietary extras, because that is exactly what would re-weld the node to a single backend. The decision is cheaply reversible, since swapping backends is swapping an implementation or a config URL, and it keeps the production-language reversibility of ADR-0003 intact.

## Open questions

Whether a future capability manifest needs richer backend introspection (context length, modalities, quantization, live load) than the standard model-list call exposes, and if so whether that is read through the interface or declared in config. Left open until a real need appears.
