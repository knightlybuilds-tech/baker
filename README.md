# baker

**Bake projects every AI agent understands on day one.**

```bash
baker new shop --recipe laravel
```

One command creates the project and sets up everything AI coding agents need: one `AGENTS.md` that Claude Code, Codex, Gemini CLI, Cursor, Copilot, OpenCode and Aider all read, plain-markdown memory, shared skills, MCP config, and **continuity**: when you hit a usage limit, `baker resume codex` picks up the task in another agent or account.

No database. No lock-in. Everything is markdown you can `cat`, `grep`, commit, and open in Obsidian.

> Status: **v0.1, early.** The standard and the CLI are young; expect changes. Feedback and recipes welcome.

---

## Why

If you work on many projects with AI agents, you know these chores:

- **Every new project needs the AI setup again:** `CLAUDE.md`, `AGENTS.md`, `GEMINI.md`, Cursor rules, MCP config, skills… copied by hand, drifting apart.
- **Every tool wants its own file**, and they disagree on names, imports and nesting.
- **Memory is a mess:** one giant `CLAUDE.md` that loads every session, goes stale, and nobody trusts.
- **Usage limits kill your flow:** you switch accounts, or from Claude to Codex, and the new session knows nothing.

baker fixes these with a small, open standard and a single binary.

## What you get

```
shop/
├── AGENTS.md              ← the one file you edit: your notes + stack + memory rules
├── CLAUDE.md              ← bridge: @AGENTS.md            (generated)
├── GEMINI.md              ← bridge: @AGENTS.md            (generated)
├── .agents/skills/        ← session-close, consolidate, ingest + stack skills
├── .claude/skills/        ← mirror for Claude Code        (generated)
├── .claude/settings.json  ← checkpoint hook               (generated)
├── .mcp.json …            ← MCP servers for every agent   (generated)
├── memory/
│   ├── MEMORY.md          ← index: one line per note
│   ├── hot.md             ← live handoff note: goal, done, next (gitignored)
│   ├── log/               ← one file per day
│   ├── decisions/         ← atomic notes: why we chose X
│   ├── lessons/           ← atomic notes: what we learned the hard way
│   └── knowledge/         ← atomic notes: domain facts
├── baker.toml             ← manifest: recipe, agents, MCP servers
└── .baker/lock.toml       ← hashes so baker never overwrites your edits
```

Codex, Cursor, Copilot, OpenCode, Windsurf, Amp, Zed and Cline read `AGENTS.md` natively, so they need no extra file. baker never writes legacy files like `.cursorrules`, which hide `AGENTS.md` from some tools.

## Install

baker is a single Go binary.

```bash
go install github.com/roguenet/baker/cmd/baker@latest
```

Or build from source:

```bash
git clone https://github.com/roguenet/baker && cd baker
make install        # builds and copies baker to $GOPATH/bin
```

Make sure `baker` is on your `PATH`: agent hooks call it to save checkpoints.

## Quick start

```bash
baker global init          # once: creates ~/.baker (your personal rules + cross-project memory)
baker new shop -r laravel  # new project: runs the Laravel installer, then adds the AI layer
cd shop && claude          # or codex, gemini, cursor… they all read AGENTS.md
```

Already have a project?

```bash
cd my-project
baker init                 # detects the stack, shows a plan, asks before writing
```

Existing `CLAUDE.md` / `GEMINI.md` instructions are moved into `AGENTS.md`, never thrown away.

## Commands

| Command | What it does |
|---|---|
| `baker new <name> -r <recipe>` | Create a folder, run the recipe's installer, add the AI layer, `git init` |
| `baker init` | Add the AI layer to the current project (plan → confirm → apply) |
| `baker sync` | Regenerate bridges, skills, MCP config, hooks and managed blocks. Your edits are protected (`--force` overrides, `--dry-run` previews) |
| `baker doctor` | Check bridges, file sizes, legacy files, hooks and memory health |
| `baker checkpoint` | Snapshot git state + latest transcript (run automatically by hooks) |
| `baker handoff` | Print a resume brief (`--copy` to clipboard, `-o file`) |
| `baker resume <agent>` | Launch another agent with the resume brief as its first prompt |
| `baker global init` | Create `~/.baker` |
| `baker global sync` | Push your personal rules into every installed agent's user-level file |
| `baker list` | Projects baker created or adopted |
| `baker recipes` | Available recipes and agents |

Run `baker <command> --help` for flags.

## Continuity: switch agents mid-task

You're halfway through a task and Claude hits its usage limit:

```bash
baker resume codex
```

Codex starts with a brief built from three layers. Any one of them is enough to continue:

1. **`memory/hot.md`**, written by the agent after every finished step (Goal, Done, Next, Gotchas).
2. **Checkpoints**, saved by a hook after every agent turn: branch, uncommitted changes, todo list, last messages.
3. **Transcript recovery**: baker reads the previous agent's local transcript (Claude Code, Codex, Gemini CLI), including from other accounts' config folders.

Everything stays on your machine, and secrets are redacted. → [docs/continuity.md](docs/continuity.md)

## Memory: small notes, not a giant file

Memory is **atomic notes**: one idea per file, the title is the claim, under 200 words, linked with `[[wikilinks]]`, indexed one line each in `memory/MEMORY.md`. Agents read the short index and open only what they need. Every note carries `as_of` / `last_verified` / `source`, and `baker doctor` flags stale, oversized, orphaned or broken notes.

Linked notes already form a graph, and Obsidian's graph view shows it with no extra tooling. → [docs/memory.md](docs/memory.md)

## Two layers: personal and project

| | `~/.baker/` (personal) | project (team) |
|---|---|---|
| Rules | `AGENTS.md`: your taste, every project | `AGENTS.md`: this project |
| Memory | cross-project lessons | `memory/` |
| Skills | `skills/` | `.agents/skills/` |
| Recipes | `recipes/`: your own or overrides | `.baker/recipes/` |

`baker global sync` writes your personal rules into each tool's *own* user-level file (`~/.claude/CLAUDE.md`, `~/.codex/AGENTS.md`, `~/.gemini/GEMINI.md`…). Projects stay team-safe: nothing committed references your home folder. → [docs/configuration.md](docs/configuration.md)

## Recipes

A recipe is a folder: `recipe.toml` (installer, detection, facts, MCP servers), an `agents.md` template, and optional skills.

| Recipe | Installer | Extras |
|---|---|---|
| `generic` | none | starter Commands/Conventions sections |
| `laravel` | `composer create-project laravel/laravel` | detects Laravel/PHP version and Pest vs PHPUnit; registers Laravel Boost MCP when installed; `laravel-migration` skill |

Override any recipe by putting your own in `~/.baker/recipes/<name>/`. → [docs/recipes.md](docs/recipes.md)

## Supported agents

| Agent | Reads AGENTS.md via | MCP file | Checkpoint hook | Transcript recovery |
|---|---|---|---|---|
| Claude Code | `CLAUDE.md` → `@AGENTS.md` | `.mcp.json` | ✅ Stop + PreCompact | ✅ |
| Codex | native | `.codex/config.toml` | ✅ `notify` (global sync) | ✅ |
| Gemini CLI | `GEMINI.md` → `@AGENTS.md` | `.gemini/settings.json` | — | ✅ (best effort) |
| Cursor | native | `.cursor/mcp.json` | — | — |
| GitHub Copilot | native | `.vscode/mcp.json` | — | — |
| OpenCode | native | `opencode.json` | — | — |
| Aider | `.aider.conf.yml` → `read: AGENTS.md` | — | — | — |
| Windsurf, Amp, Zed, Cline | native | — | — | — |

Details and file paths to re-verify: [docs/agents.md](docs/agents.md)

## Principles

1. **One source, many agents.** `AGENTS.md` is the only instruction file you edit.
2. **Plain markdown, no database.** Obsidian is a nice viewer, never a dependency.
3. **The repo is team-safe; the global layer is personal.**
4. **baker owns marked blocks; you own the rest.** It only touches text between `baker:begin` / `baker:end` markers, and a lock file stops it overwriting your edits.
5. **Any session can die; the work survives.**
6. **Plan, then apply.**

The full design is in [docs/spec.md](docs/spec.md), and the research behind it is in [docs/research/](docs/research/landscape-report.md).

## Roadmap

- [ ] `baker graph`: optional [graphify](https://github.com/safishamsi/graphify) code-graph adapter (designed in the spec)
- [ ] Remote recipes: `--recipe github:user/repo`
- [ ] More recipes: Next.js, Django, Rails, Go
- [ ] Gemini CLI checkpoint hook
- [ ] Homebrew tap and prebuilt release binaries
- [ ] Organization profiles for the vault (PARA, Johnny Decimal)

## Development

```bash
make test     # go test ./...
make build    # ./bin/baker
make lint     # go vet + gofmt check
```

Code map: `cmd/baker` (entry) · `internal/cli` (commands) · `internal/scaffold` (plan/apply engine, build, doctor, global) · `internal/blocks` (managed regions) · `internal/recipe` · `internal/memory` (note checks) · `internal/continuity` (transcripts, checkpoints, handoff) · `internal/agents` (agent registry) · `assets/` (embedded core block, skills, recipes).

## License

To be decided before the first public release.
