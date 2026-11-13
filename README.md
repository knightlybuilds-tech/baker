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
├── .gemini/settings.json  ← checkpoint hook + MCP         (generated)
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

baker is a single Go binary. Download it for your OS from
[Releases](https://github.com/knightlybuilds-tech/baker/releases) (Linux, macOS,
Windows; amd64 and arm64), unpack it, and put `baker` on your PATH. Or with Go:

```bash
go install github.com/knightlybuilds-tech/baker/cmd/baker@latest
```

Or build from source:

```bash
git clone https://github.com/knightlybuilds-tech/baker && cd baker
make install        # builds, copies baker to $GOPATH/bin, and checks your PATH
```

### Put baker on your PATH

`go install` puts binaries in `~/go/bin` (or `$GOBIN`), which often isn't on your `PATH`. Then typing `baker` says "command not found", and agent hooks can't save checkpoints. Add the folder once:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

Use `~/.bashrc` for bash, or `fish_add_path ~/go/bin` for fish. You don't have to remember this: `make install`, and baker itself whenever you run it by its full path (`~/go/bin/baker`), detect the problem and print the right command for your shell. Set `BAKER_NO_PATH_HINT=1` to silence it.

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
| `baker memory list / recall / check` | Discover, retrieve, and check ordinary local notes without an account or approval |
| `baker bakery` | Browse this project's memory as cards, a list and a connected web on `127.0.0.1` (read-only, live-updating) |
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

Codex starts with a brief built from three complementary recovery layers:

1. **`memory/hot.md`**, written by the agent after every finished step (Goal, Done, Next, Gotchas).
2. **Checkpoints**, saved by a hook after every agent turn: branch, uncommitted changes, todo list, last messages.
3. **Transcript recovery**: baker reads the previous agent's local transcript (Claude Code, Codex, Gemini CLI), including from other accounts' config folders.

Baker builds the brief locally and applies best-effort secret redaction. Passing
it to a hosted agent sends its contents to that provider. → [docs/continuity.md](docs/continuity.md)

## Memory: small notes, not a giant file

Memory is **atomic notes**: one idea per file, the title is the claim, under 200 words, linked with `[[wikilinks]]`, indexed one line each in `memory/MEMORY.md`. Agents read the short index and open only what they need. Every note carries `as_of` / `last_verified` / `source`, and `baker doctor` flags stale, oversized, orphaned or broken notes.

Linked notes already form a graph, and Obsidian's graph view shows it with no extra tooling. → [docs/memory.md](docs/memory.md)

### Standalone memory retention

Baker works without Bakery: agents retain ordinary Markdown notes, recover
unfinished work, and switch tools locally. Use `baker memory list --json` to
inventory project knowledge, `baker memory recall --query "topic" --json` to
retrieve it, and `baker memory check` to check note health. Stale or incomplete
notes are flagged and remain available as leads to verify.

Team-first means project knowledge is portable and safe to commit, while
personal handoffs and checkpoints remain local. Bakery will add collaboration;
it is not required for local capture, retention, retrieval, or continuation.
See [the standalone scope](docs/implementation-plan.md).

### See your memory: `baker bakery`

```bash
baker bakery            # opens http://127.0.0.1:4410/
baker bakery --no-open --port 4500
```

A local, read-only Bakery view of `memory/`: cards, a list and a connected web
built from your `[[links]]`, with `baker memory check` findings at the top and
freshness flags on each note. It refreshes within a few seconds of a note
changing. It binds to loopback only, never writes files, and doesn't show
`hot.md`, checkpoints or transcripts. Review, branch and team screens belong to
the future hosted Bakery and are marked as not yet available. See
[docs/bakery.md](docs/bakery.md).

The earlier review-lifecycle prototype is retained under `baker experimental
crumb`; it is not part of the default memory workflow. Its evaluations and
limitations are in [reviewed-memory.md](docs/reviewed-memory.md).

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
| Codex | native | `.codex/config.toml` (trusted projects) | ✅ `notify` (global sync) | ✅ |
| Gemini CLI | `GEMINI.md` → `@AGENTS.md` | `.gemini/settings.json` | ✅ AfterAgent + PreCompress | ✅ |
| Cursor | native | `.cursor/mcp.json` | — | — |
| GitHub Copilot | native | `.vscode/mcp.json` | — | — |
| OpenCode | native | `opencode.json` | — | — |
| Aider | `.aider.conf.yml` → `read: AGENTS.md` | — | — | — |
| Windsurf, Amp, Zed, Cline | native | — | — | — |

Details, file paths and what's still unverified: [docs/agents.md](docs/agents.md)

## Principles

1. **One source, many agents.** `AGENTS.md` is the only instruction file you edit.
2. **Plain markdown, no database.** Obsidian is a nice viewer, never a dependency.
3. **The repo is team-safe; the global layer is personal.**
4. **baker owns marked blocks; you own the rest.** It only touches text between `baker:begin` / `baker:end` markers, and a lock file stops it overwriting your edits.
5. **Any session can die; the work survives.**
6. **Plan, then apply.**

The full design is in [docs/spec.md](docs/spec.md), and the research behind it is in [docs/research/](docs/research/landscape-report.md).

## Roadmap
- [ ] `baker learn`: propose cited crumbs for an existing project from git history and (opt-in) past agent sessions, triaged in an inbox ([design](docs/learn.md))

- [ ] `baker graph`: optional [graphify](https://github.com/safishamsi/graphify) code-graph adapter (designed in the spec)
- [ ] Remote recipes: `--recipe github:user/repo`
- [ ] More recipes: Next.js, Django, Rails, Go
- [x] Prebuilt release binaries
- [ ] Homebrew tap
- [ ] Organization profiles for the vault (PARA, Johnny Decimal)

## Development

```bash
make test     # go test ./...
make build    # ./bin/baker
make lint     # go vet + gofmt check
make bakery-ui  # sync the Bakery UI from ../bakery-workspace/dist into internal/bakery/ui
```

Code map: `cmd/baker` (entry) · `internal/cli` (commands) · `internal/scaffold` (plan/apply engine, build, doctor, global) · `internal/blocks` (managed regions) · `internal/recipe` · `internal/memory` (note checks) · `internal/bakery` (local Bakery server + embedded UI) · `internal/continuity` (transcripts, checkpoints, handoff) · `internal/agents` (agent registry) · `assets/` (embedded core block, skills, recipes).

## Releases

Every merge to `main` runs the tests and publishes a release with binaries and
checksums ([release.yml](.github/workflows/release.yml), [GoReleaser](.goreleaser.yaml)).
The version is the next patch by default. Put `#minor` or `#major` in the merge
commit message to bump further, or `#skip-release` to skip. The first release
is `v0.1.0`.

## License

[MIT](LICENSE)
