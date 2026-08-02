# Project Context

> Living document, updated at the end of each coding session.
> Tool-agnostic and readable by any AI assistant.

## Read these first, and do not duplicate them here

This project already records its thinking properly. This file is only for what those
do not carry: the live environment, and gotchas that cost real time.

| Source | What it holds |
|---|---|
| `TIMELINE.md` | Dated running log of what happened and why. The project's narrative memory. |
| `docs/ADR/` | Numbered decision records (ADR-0001 … ADR-0017). Load-bearing decisions live here, not in this file. |
| `architecture/Overview.md` | The RFC: topology, principles, build phases. |
| `architecture/APIs/README.md` | The node's interface **as built** — endpoints, doors, failure responses. |
| `architecture/Sketchbook.md` | Design record that fed the Week-3 build, kept as history. |
| `docs/Functional_Node_Plan.md` | The live build plan: M1 (done) through M4. |
| `CONTRIBUTING.md` | The project is AI-drafted and human-directed, stated plainly. Claims should trace to sources. |

**The owner handles all git commits** unless they explicitly ask otherwise. No AI
attribution, signatures, or emojis in commits or PRs.

---

## The live pod (as of 2026-08-01)

One Go binary is a node; a pod is nodes that list each other as peers. The node is
control plane only — it never runs inference, it brokers to a model server.

| Node | Host | Address | Runs as | Serves |
|---|---|---|---|---|
| **atlas** | Linux, AMD GPU (ROCm/Vulkan) | `100.118.77.40:8090` (Tailscale) | Docker, `network_mode: host` | 10 models incl. `all-minilm:l6-v2` |
| **laptop** | macOS | `localhost:8080`, Tailscale `100.94.5.71` | Docker Desktop, **bridged** | `llama3.2:1b`, `all-minilm:l6-v2` |
| **Open WebUI** | macOS | `127.0.0.1:3000` | Docker, opt-in compose overlay | a consumer of the laptop node |

- SSH to atlas is **`ssh atlas@atlas`**. Plain `ssh atlas` fails with a Tailscale
  user-lookup error. Tailscale SSH may prompt for browser re-auth mid-command.
- Atlas has **no Go toolchain** and **no passwordless sudo**. Cross-compile locally
  (`GOOS=linux GOARCH=amd64`) and copy over; anything needing root needs the owner.
- Deployment source lives at `~/dii-node/` on atlas; compose runs from
  `~/dii-node/deploy/`. `~/dii/` holds the pre-container binary and config, kept for
  rollback.
- Verify anything with `prototype/scripts/m1-check.sh <url>`; it exits non-zero on
  failure. Atlas needs `DII_CHAT_MODEL=qwen2.5:7b` since it has no `llama3.2:1b`.

---

## Architecture Decisions

Only decisions made in-session and not yet promoted to an ADR. Anything load-bearing
should graduate to `docs/ADR/`.

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-07-30 | `/v1/models` advertises the **pod** (own + cached peers', with `owned_by` naming where each runs); `/manifest` stays own-only | A node would overflow models it never advertised, so overflow was undiscoverable to any client that builds a picker from the list — ADR-0004's "any OpenAI client is a DII client" was true for chat, false for discovery. `/manifest` answers a *peer's* question; republishing borrowed capability would have peers advertising each other's models back and forth. |
| 2026-07-30 | Widen `modelserver.Backend` with `Embeddings` rather than a separate optional `Embedder` interface | Embeddings are a required third of the floor bundle (ADR-0006), not an optional capability, and all three implementations serve it trivially. Keeps "one seam" literally true. |
| 2026-07-30 | Do **not** inject `stream_options.include_usage` to force exact token counts | The node promises to forward request bodies verbatim, and injection appends a usage chunk strict clients do not expect. Cost: default streams report estimated completion tokens and no prompt tokens. |
| 2026-07-30 | Response content type follows the request's own `stream` field | Authoritative per the OpenAI spec, versus sniffing the response. Fixed the POC announcing every response as SSE. |
| 2026-07-30 | Usage records key on an opaque handle and carry **no** prompt/completion text | The ADR-0016 data-minimisation seam, built in from the start so a later fair-use mechanism cannot drift into a behavioural reputation store. |
| 2026-07-31 | Linux containers use `network_mode: host`; macOS bridges via an overlay | Host networking makes the container config identical to a host config, lets `listen` bind a specific interface (a bridged container can only *publish* to one), and sidesteps the firewall/subnet problem below. Docker Desktop has no real host networking. |
| 2026-07-31 | Dropped `com.dii.node.plist`; the plan asked for a *note*, not a file | macOS runs the node fine in Docker (verified), so a native plist served only "macOS user who refuses Docker" — an untested file that looked supported. Unnecessary code that *looks* tested is worse than none. |
| 2026-08-01 | Docker is the daily path; the systemd path is a **gate run at each milestone boundary** | Whichever path is not in daily use is the one that rots. Catches what a container hides: a dependency on a mounted path, an env var only compose sets, a signal only Docker delivers, a libc dependency in a supposedly static binary. Steps in `prototype/deploy/README.md`. |

---

## Known Gotchas

### Environment and tooling

- **`pkill -f "<pattern>"` over SSH matches its own command string** and kills the
  session. It took atlas down for ~90s. Use the PID, or a pattern that cannot appear
  in your own command line.
- **macOS `bsdtar --exclude='./node'` also matches `./cmd/node`** — patterns match any
  path component, not just anchored paths. This silently shipped an empty `cmd/` to
  atlas and the Docker build failed on `stat /src/cmd/node`. `prototype/.dockerignore`
  now defines the build context so it does not depend on the copy command.
- **`gh` is authenticated as `bryanbeginly` (READ on `primebevil/dii`), but git's SSH
  alias authenticates as `primebevil`.** `git push` works; `gh pr create` fails with
  "must be a collaborator". Either `gh auth login` as `primebevil` or use the browser
  link that `git push` prints.
- **Ollama does not auto-start after a laptop reboot.** The node then comes up healthy
  but with an empty model list. `/readyz` reports it correctly.
- The owner's `~/.zshrc` `ask()` defaults `DII_HUB` to `http://localhost:8080` and
  **swallows errors** — a blank line means the target is down, not that the model
  replied emptily. `export DII_HUB=http://100.118.77.40:8090` to hit atlas.

### Docker

- **A distroless image has no shell**, so a shell-form `HEALTHCHECK` cannot execute at
  all. `docker run --health-cmd` *always* wraps in `/bin/sh -c`, so it reports the
  container unhealthy while the identical probe run by hand exits 0. The healthcheck
  lives in the Dockerfile in JSON exec form for this reason.
- **Compose creates its own network subnet** (`172.19.0.0/16`), not the default bridge
  (`172.17.0.0/16`). Atlas's firewall admits only the default bridge, so a bridged
  compose node could not reach the host's Ollama — packets dropped, presenting as a
  timeout, and the node started with an empty model list. This is why Linux uses host
  networking.
- **`stop_grace_period` must exceed the node's `shutdown_timeout`** (15s) or the
  runtime SIGKILLs mid-drain. Docker's 10s default is too tight; compose sets 30s.

### Open WebUI (a consumer, not part of the node)

- **Settings are persisted to its database on first run, after which environment
  changes are silently ignored.** `ENABLE_PERSISTENT_CONFIG: "false"` makes compose the
  source of truth — at the cost that UI changes to those settings do not survive.
- **RAG uses its own bundled embedding model by default** and would never touch the
  node. `RAG_EMBEDDING_ENGINE=openai` plus `RAG_OPENAI_API_BASE_URL` points it at the
  node — which is the only thing that exercises `/v1/embeddings` realistically.
- **Title, tag and follow-up generation default to the chat model**, turning one
  question into a burst of 5–7 completions on a 30B. `TASK_MODEL_EXTERNAL` points them
  at a small local model.
- An **empty** `OPENAI_API_KEY` means no auth header, so the node's *local* door
  (local-first). The consumer token opens the *consumer* door, which is peer-first and
  routes work away on purpose. Any other value is a 401 — there is no third option.
- `:main` is the **development branch**, not a release. Pin a version tag.
- Notes-as-a-chat-surface is a real v0.11.0 feature ("The Interface, Reorganized"), not
  a malfunction. `ENABLE_NOTES: "false"` removes it.

### The node itself

- **Manifests are fetched once at startup.** Restart a node to pick up a newly pulled
  model or a peer that came up late. Now that `/v1/models` advertises peers, the list
  can also name a model whose holder has since gone down — that fails as a `502` from
  the dead peer rather than an immediate `503`.
- `prototype/deploy/config-docker.yaml` is **gitignored** — it holds a real address and
  token. `config-docker.example.yaml` is the one that ships.

---

## Established Patterns

- **Verify by deploying, not just by testing.** Four M1 defects were found only by
  standing the real thing up: the distroless healthcheck, the compose subnet, the
  `/v1/models` discovery gap, and the GUI task-model churn. None would have surfaced in
  a unit test.
- **Say what is unverified.** `deploy/README.md` carries a status column and states
  exactly which parts of the systemd unit were exercised and which were not. Config
  that looks supported but has never run is worse than absent — for a human or an
  assistant reading the repo.
- **Milestones stop for review.** Build one at a time; leave the node runnable and the
  owner/consumer/overflow behaviour intact at each boundary.
- **`m1-check.sh` is the arbiter.** When the GUI or anything else misbehaves, run it
  first: if it passes, the node is not the problem.

---

## Open Questions

- [ ] **Exact token accounting.** Default streams report a content-chunk proxy for
  completion tokens and **no prompt tokens** — and prompt tokens dominate long-context
  cost. The chunk proxy tracks Ollama closely (14 chunks vs 14 tokens observed), but
  the gap is real. Likely fix is an opt-in config flag injecting
  `stream_options.include_usage`. Blocks precise M3 budgets. Owner has flagged usage as
  important.
- [ ] **Five error-handling defects carried from the M1 review into M2**, listed in
  full in `docs/Functional_Node_Plan.md` under "Carried into M2 from the M1 code
  review": mid-stream failures logged as `200`, oversized request bodies truncated
  into a wrong `503` instead of a `413`, oversized embeddings responses truncated and
  served as `200`, upstream errors leaking internal topology to the consumer door, and
  `ServedBy` attributed to peers that never ran the work. Three of the five corrupt
  the accounting record, which M3 budgets depend on, so they should land before M3.
- [ ] **Periodic manifest refresh.** Would fix both "restart a node to notice its peer"
  and the stale-model-list staleness above. Currently the highest-value small addition.
- [ ] **M2 shape, to settle before starting:** consumer store as a single JSON file with
  atomic write vs. SQLite (pick whichever M3's concurrent counters build on), and
  management as a CLI subcommand vs. a localhost-only admin listener.
- [ ] **systemd root install path still unexercised** — `User=`/`Group=`,
  `ProtectHome=true`, `After=ollama.service`, and `/opt/dii` + `/etc/dii`. Do the real
  install once at the M2 boundary and the gap closes permanently.
- [ ] **MoE vs dense for the floor bundle.** Atlas holds `qwen3:30b` and
  `qwen3-coder:30b` (MoE, small active parameter count) alongside dense 24–32B models —
  the exact comparison `docs/Reliable_Floor_Definition.md` flagged as "worth watching".
  M1's per-request `ttft_ms`/`latency_ms` makes this measurable as a side effect of
  normal use.
- [ ] **`gh` auth for `primebevil`**, so PRs can be opened from the CLI.
- [ ] Gated on the ADR-0016 legal review regardless of code: serving the unaffiliated
  public. M4 builds the moderation seam; it does not open the door.
