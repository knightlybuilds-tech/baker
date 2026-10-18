# baker

## Project notes
baker is an open-source, AI-agnostic CLI (Go) that scaffolds AI-ready projects: one AGENTS.md bridged to every agent, plain-markdown atomic memory, shared skills, MCP config, and continuity across agents/accounts (`baker resume <agent>`).
The design lives in `docs/spec.md`; the research behind it in `docs/research/`.

## Commands
- Build: `make build` (→ `./bin/baker`)
- Test: `make test` (`go test ./...`)
- Lint / format: `make lint` (go vet + gofmt check); fix formatting with `gofmt -w .`
- Try it safely: set `BAKER_HOME` and `BAKER_USER_HOME` to temp dirs so you never touch the real `~/.baker` or `~/.claude`.

## Conventions
- **Document as you build:** any new command, flag, file layout or recipe change updates `README.md` and the matching `docs/*.md` in the same change.
- Everything baker writes is computed into a `scaffold.Plan` first; only `Plan.Apply` touches disk. Never write files directly from commands.
- Shared files are edited only through managed blocks (`internal/blocks`); whole generated files go through `Plan.Managed` so the lock protects user edits.
- Templates, core block, skills and built-in recipes live in `assets/` (embedded). Agent-specific paths live only in `internal/agents/agents.go`.
- Hooks (`baker checkpoint --quiet`) must never fail or print: an error there breaks the user's agent.
- Keep dependencies minimal (cobra, go-toml). Standard library first.

<!-- baker:begin recipe:generic@1 -->
## Stack
No stack recipe is applied. The stack, commands and conventions are described in Project notes above.
If they're empty, ask the user (or inspect the repo) and propose filling them in.
<!-- baker:end recipe:generic -->

<!-- baker:begin core@0.1 -->
## Memory
This project keeps its memory in `memory/` as small markdown notes. Any agent can read and write it.

- **Session start:** read `memory/hot.md` and `memory/MEMORY.md`. Open other memory notes only when relevant to the task.
- **Continuing work:** if `hot.md` has unchecked items under Next, you are picking up someone else's task (another agent, account, or provider). Say the plan back to the user in one line, then continue.
- **Checkpoint as you go:** after each finished step, before any long or risky operation, and whenever the user changes direction, update `memory/hot.md` (Goal, Done, Next, Gotchas, Open questions). Your session can end at any moment (usage limits) — never keep progress only in your head.
- **Session end:** follow the `session-close` skill: append to `memory/log/<YYYY-MM-DD>.md`, update `hot.md`, propose lessons.
- **Atomic notes:** one idea per file in `memory/decisions/`, `memory/lessons/` or `memory/knowledge/`. The title is the claim as a full sentence, the filename is its slug, under 200 words. Link related notes with `[[slug]]`. Add one line per note to `memory/MEMORY.md`.
- **Freshness:** every note has `type`, `as_of`, `last_verified`, `source` frontmatter. Treat a note not verified in 90 days as a hint, not a fact; re-verify before relying on it.
- **Promotion:** you may propose turning a lesson into a rule in this file, but only the user approves it.
- **Never** write secrets, tokens, passwords or customer data into memory.
<!-- baker:end core -->
