# baker — scaffolding spec (draft v0.1)

> `baker new my-app --recipe laravel` → a project folder that every AI agent understands on day one.

Status: design document for standard v0.1. Core decisions locked (section 6).

> **Implementation notes (v0.1.0):**
> - `baker upgrade` is folded into `baker sync`: sync re-renders core/recipe blocks and uses the lock to protect your edits.
> - `baker graph` (graphify adapter, section 3.5) is designed but **not implemented yet**.
> - Remote recipes (`--recipe github:user/repo`) are **not implemented yet**.
> - Recipe facts use a `[[facts]]` array — see [recipes.md](recipes.md) for the real format.
> - Gemini CLI checkpoint hooks are not wired yet; Claude Code (Stop/PreCompact) and Codex (`notify`) are.

---

## 1. Principles

1. **One source, many agents.** `AGENTS.md` is the only instruction file a human edits. Every vendor file is generated from it.
2. **Plain markdown, no database.** Everything works with `cat` and `grep`. Obsidian is a nice viewer, never a dependency.
3. **The repo is team-safe; the global layer is personal.** Nothing committed to a project references your home folder or personal preferences.
4. **Baker owns marked blocks, you own the rest.** Generated content lives between `baker:begin` / `baker:end` markers. Baker never touches text outside them.
5. **Any session can die; the work survives.** Progress is written down continuously so any agent, on any account or provider, can pick it up (section 3.6).
6. **Plan, then apply.** Any command that changes an existing folder shows what it will write and waits for a yes (`--yes` skips).

---

## 2. Two layers

| Layer | Where | Committed? | Holds |
|---|---|---|---|
| **Global (personal)** | `~/.baker/` | No (your own dotfiles repo if you like) | Your preferences, your recipes, cross-project lessons, project registry |
| **Project (team)** | `my-app/` | Yes | Stack instructions, project memory, skills, MCP config |

### 2.1 Global layer: `~/.baker/`

```
~/.baker/
├── config.toml        # defaults: agents to target, default recipe, author
├── AGENTS.md          # YOUR rules ("I use Pest", "prefer Actions over fat controllers")
├── skills/            # personal skills, available in every project
├── recipes/           # your own recipes, or overrides of built-in ones
├── memory/            # cross-project knowledge (open as an Obsidian vault)
│   ├── MEMORY.md      # index
│   └── lessons/       # lessons that proved true in more than one project
└── projects.toml      # registry of every project baker created or adopted
```

**How global rules reach agents:** not by copying them into projects. `baker global sync` writes a managed block into each tool's *own* user-level file:

| Agent | User-level file baker writes |
|---|---|
| Claude Code | `~/.claude/CLAUDE.md` (+ skills to `~/.claude/skills/`) |
| Codex | `~/.codex/AGENTS.md` |
| Gemini CLI | `~/.gemini/GEMINI.md` |
| OpenCode | `~/.config/opencode/AGENTS.md` |
| Cursor | none (user rules live in the settings UI) — baker prints them for you to paste |

The managed block contains your `~/.baker/AGENTS.md` plus one pointer:
"Cross-project lessons live in `~/.baker/memory/`; read `MEMORY.md` when a task may benefit."

So every project, even ones baker never touched, gets your personal rules and memory.

*(Exact paths per tool need verifying against current docs before v1.)*

### 2.2 Project layer

```
my-app/
├── AGENTS.md                 # source of truth (core + recipe blocks + your notes)
├── CLAUDE.md                 # bridge: @AGENTS.md
├── GEMINI.md                 # bridge: @AGENTS.md
├── .aider.conf.yml           # bridge: read: AGENTS.md   (only if aider enabled)
├── .agents/skills/           # canonical skills
├── .claude/skills/           # copy of .agents/skills (Claude Code doesn't read .agents/)
├── .mcp.json                 # + other tools' MCP files, generated from baker.toml
├── memory/                   # committed (except hot.md)
│   ├── MEMORY.md             # index: one line per note
│   ├── hot.md                # handoff state: goal, done, next — gitignored (see 3.6)
│   ├── log/2026-10-18.md     # one file per day (avoids merge conflicts in teams)
│   ├── decisions/            # atomic notes: why we chose X
│   ├── lessons/              # atomic notes: what we learned the hard way
│   └── knowledge/            # atomic notes: domain facts ("a Tenant owns many Stores")
├── baker.toml                # manifest: recipe, standard version, agents, MCP servers, graph
└── .baker/
    ├── lock.toml             # hashes of generated files (for safe upgrades)
    ├── checkpoints/          # gitignored: automatic session snapshots (see 3.6)
    └── graph/                # derived, gitignored: graphify output (optional)
```

**Agents that read `AGENTS.md` natively get no extra file:** Codex, Cursor, Copilot, OpenCode, Amp, Windsurf, Cline, Zed.
Baker **never** writes legacy root files (`.cursorrules`, `.windsurfrules`, `.rules`) — they hide `AGENTS.md` from Zed.

---

## 3. How the initial files are built

### 3.1 AGENTS.md is composed from blocks

```
core block      (from baker itself: memory protocol, session rituals)
recipe block    (from the recipe: stack, commands, conventions — templated with detected facts)
project notes   (empty section for you; baker never edits it)
```

Example output for `baker new my-app --recipe laravel`:

```markdown
# my-app

## Project notes
<!-- Yours. Baker never edits outside its markers. -->

<!-- baker:begin recipe:laravel@1 -->
## Stack
Laravel 13 on PHP 8.4. Tests: Pest. Frontend: Vite + Tailwind.

## Commands
- Dev server: `composer run dev`
- Tests: `php artisan test`
- Format: `./vendor/bin/pint`
- Migrations: `php artisan migrate`

## Conventions
- Business logic in `app/Actions`, not controllers.
- Every migration has a working `down()`.
<!-- baker:end recipe:laravel -->

<!-- baker:begin core@0.1 -->
## Memory
- Session start: read `memory/hot.md` and `memory/MEMORY.md`. Open other memory pages only when relevant. If `hot.md` has unfinished Next steps, you are continuing someone else's work: confirm the plan with the user, then continue.
- Checkpoint as you go: after each finished step, before any long or risky operation, and whenever the user changes direction, update `memory/hot.md` (Goal, Done, Next, Gotchas). Your session can end at any moment (usage limits) — never hold progress only in your head.
- Session end: run the `session-close` skill — append to `memory/log/<today>.md`, update `hot.md`, propose lessons.
- Write memory as atomic notes: one idea per file, title is the claim, under 200 words, link related notes with `[[slug]]`, add one line to `MEMORY.md`.
- Notes carry `as_of`, `last_verified`, `source` frontmatter. Treat anything older than 90 days unverified as a hint, not a fact.
- Never write secrets, tokens, or customer data to memory.
<!-- baker:end core -->
```

Facts like "PHP 8.4" and "Pest" are **detected** after the installer runs (from `composer.json`, `package.json`, etc.), not hardcoded.

### 3.2 Bridges are tiny and fully generated

`CLAUDE.md`:
```markdown
<!-- baker:begin bridge@1 — generated; edit AGENTS.md instead -->
@AGENTS.md
<!-- baker:end bridge -->
```
You may add Claude-only notes below the block. Same pattern for `GEMINI.md`.

Why an import and not a symlink: Claude Code's edit tools don't write through symlinks, and committed symlinks break on Windows clones. Claude Code also stops reading `AGENTS.md` once a `CLAUDE.md` exists, so the import is required.

### 3.3 Memory is atomic notes

Every memory entry is a **small, single-idea note**. Bite-sized notes mean an agent reads the one-line index, opens only the 2–3 notes it needs, and never loads a giant file. They're also easy to verify, update, or delete one at a time.

**Rules (enforced by `baker doctor`):**

1. **One idea per note.** If it has two claims, it's two notes.
2. **Title is the claim**, as a full sentence. Filename is its slug:
   `lessons/queue-workers-must-restart-after-deploy.md`
3. **Small.** Aim for under 200 words; `doctor` warns above 400.
4. **Link, don't repeat.** Reference other notes with `[[slug]]` wikilinks (Obsidian-native, plain text for agents).
5. **Every note is in the index.** One line in `MEMORY.md`. `doctor` flags orphans and broken links.
6. **Every note has freshness frontmatter.**

Example note:
```markdown
---
type: lesson              # lesson | decision | knowledge
scope: project            # project | global
tags: [queues, deploy]
as_of: 2026-10-18
last_verified: 2026-10-18
source: session 2026-10-18
confidence: high          # high | medium | low
---
# Queue workers must restart after every deploy

Workers keep the old code in memory, so new jobs ran stale logic after the 10/12 deploy.

**Apply:** the deploy script runs `php artisan queue:restart`. See [[deploys-use-envoyer]].
```

`memory/MEMORY.md` (the index agents read at session start):
```markdown
# Memory index

## Decisions
- [Deploys use Envoyer](decisions/deploys-use-envoyer.md) — zero-downtime, not Forge scripts

## Lessons
- [Queue workers must restart after every deploy](lessons/queue-workers-must-restart-after-deploy.md) — stale code otherwise

## Knowledge
- [A Tenant owns many Stores](knowledge/a-tenant-owns-many-stores.md) — billing is per Tenant
```

When a section passes ~30 lines, `consolidate` splits it into a topic hub note (`knowledge/_billing.md`) and the index links to the hub instead. The index stays short no matter how big memory grows.

Because notes link to each other, **memory is already a graph** — Obsidian's graph view shows it with zero extra tooling.

### 3.4 Built-in skills (in `.agents/skills/`)

| Skill | Does |
|---|---|
| `session-close` | Writes today's log entry, updates `hot.md`, proposes lessons |
| `consolidate` | Merges duplicates, flags stale facts, proposes promoting lessons to rules or to global memory (human approves) |
| `ingest` | Turns a doc/URL/PDF dropped in the chat into a memory page with source |

Recipes can add stack skills (e.g. `laravel-migration`).

### 3.5 Code graph (optional: graphify)

Two graphs, two jobs:

| | Notes graph | Code graph |
|---|---|---|
| Built from | `memory/` wikilinks | source code (+ optionally `memory/` and docs) |
| Made by | humans + agents, by writing notes | [graphify](https://github.com/safishamsi/graphify), automatically |
| Stored | committed markdown | `.baker/graph/`, gitignored, rebuildable anytime |
| Answers | "why / what did we learn" | "what calls what, where does X live" |

Baker treats graphify as an **optional adapter, never the source of truth**. Enable with `--graph` on `new`/`init`, or in `baker.toml`:

```toml
[graph]
engine = "graphify"       # "graphify" | "none" (default: none)
mode = "code-only"        # code-only: local tree-sitter parse, no LLM calls, nothing leaves the machine
                          # full: also sends docs/PDFs to an LLM for extraction
include_memory = true     # add memory/ notes so code and knowledge link up
output = ".baker/graph"
mcp = true                # register graphify's MCP server in every agent's MCP config
refresh = "post-commit"   # "post-commit" git hook | "manual"
```

What baker does when enabled:
1. Checks graphify is installed (`pipx install graphifyy`) and prints the command if not — baker doesn't bundle it.
2. Runs the first build into `.baker/graph/`, with query logging turned off.
3. Adds `.baker/graph/` to `.gitignore`.
4. Adds the MCP server to `.mcp.json` and the other agents' MCP files.
5. Adds one line to the core block: "For code navigation in a large repo, query the code graph before grepping."
6. `baker graph` rebuilds manually.

Default is **off**: on a small fresh project the graph saves little (graphify's own README admits ~1x on small corpora). Turn it on when a repo gets big. Before shipping, re-check graphify's license (the repo showed both MIT and Apache-2.0) and its paid "Cloud" upsell.

### 3.6 Continuity: switch accounts or providers mid-task

**The problem:** you hit a usage limit halfway through a task and switch to another account, or from Claude to Codex. The new session knows nothing. Worse, a limit cuts the session off with no warning, so you can't count on the old agent writing a goodbye note.

**The answer: three layers, from "agent wrote it" to "baker reconstructed it".** Any one of them is enough to resume; together they make resuming reliable.

| Layer | Written by | When | Survives a hard cutoff? |
|---|---|---|---|
| 1. `memory/hot.md` | the agent | after every finished step (core-block rule) | yes, up to the last step |
| 2. `.baker/checkpoints/` | baker, via agent hooks | after every agent turn | yes, up to the last turn |
| 3. Transcript recovery | baker, reading the old tool's local logs | when you run `baker handoff` | yes, everything said |

No layer needs an LLM to read it, and none is tied to a vendor.

#### Layer 1: `hot.md` is a living handoff note

Fixed structure, so any agent can write and read it:

```markdown
---
updated: 2026-10-18T14:32
agent: claude-code
branch: feat/tenant-billing
---
# Goal
Add per-Tenant Stripe billing with a monthly invoice job.

# Done
- [x] Migration + Tenant::subscription() relation
- [x] Stripe webhook controller, tests passing

# Next
- [ ] Monthly invoice job (`app/Jobs/InvoiceTenants.php`) — started, not tested
- [ ] Add queue:restart to deploy script, see [[queue-workers-must-restart-after-deploy]]

# Gotchas
- Stripe test keys live in `.env.testing`, not `.env`
- Don't touch `LegacyBillingService`, still used by old tenants

# Open questions
- Proration on mid-month plan change? (asked user, no answer yet)
```

The core block tells agents to update it **after every finished step**, not only at the end. That's the main defense against limits: progress is written down before it can be lost.

#### Layer 2: automatic checkpoints via hooks

Agents forget instructions; hooks don't. `baker init` installs a tiny hook in every agent that supports one, running `baker checkpoint` after each turn:

| Agent | Hook used |
|---|---|
| Claude Code | `Stop` and `PreCompact` hooks in `.claude/settings.json` |
| Codex | `notify` command in config |
| Gemini CLI | hooks in `.gemini/settings.json` |
| Others | none yet — layers 1 and 3 still work |

*(Hook names and config paths to be verified per tool before v1.)*

`baker checkpoint` is mechanical and fast (no LLM, under 100 ms). It saves to `.baker/checkpoints/<timestamp>.json`:
- branch, last commit, `git status`, diff stats of uncommitted changes
- which agent and account it came from
- the agent's own todo list, if its transcript has one
- the last few user and agent messages

Only the newest 20 are kept.

#### Layer 3: transcript recovery

Every major CLI agent already saves full transcripts locally. `baker handoff` knows where each one keeps them, finds the latest session **for this project folder**, and pulls out what matters:

| Agent | Transcript location (to verify) |
|---|---|
| Claude Code | `~/.claude/projects/<folder-slug>/*.jsonl` (or `$CLAUDE_CONFIG_DIR/projects/…` per account) |
| Codex | `~/.codex/sessions/YYYY/MM/DD/*.jsonl` |
| Gemini CLI | `~/.gemini/tmp/<project-hash>/` |
| OpenCode | its local storage dir |

**Switching accounts on the same provider** often means a different config dir (for example a second `CLAUDE_CONFIG_DIR`). List them all in `~/.baker/config.toml` and baker searches every one:

```toml
[transcripts]
claude = ["~/.claude", "~/.claude-work"]
codex  = ["~/.codex"]
```

#### The resume brief

`baker handoff` merges all three layers into one short prompt:

```markdown
# Continuing work on my-app

You're picking up a task another agent (claude-code, 14:32) was doing when its session ended.

## From memory/hot.md
<goal, done, next, gotchas, open questions>

## Repo state (checkpoint 14:35)
Branch feat/tenant-billing · 3 files changed, uncommitted:
  app/Jobs/InvoiceTenants.php (+84)  tests/Feature/InvoiceTest.php (+40) ...

## Last exchange before the cutoff
User: "make the invoice job skip tenants on a free plan"
Agent: "Adding the filter in handle() now, then I'll write the test…"   ← cut off here

## Your first move
Read AGENTS.md and the files above, confirm the plan with the user in one line, then continue from "Next".
```

Then either:
```bash
baker resume codex          # launches codex with the brief as its first prompt
baker handoff --copy        # copies it, paste into any agent or web chat
```

#### Rules

- **Local only.** Checkpoints and briefs are gitignored and never leave the machine.
- **Redacted.** Before writing anything, baker strips anything that looks like an API key, token, or password from transcript excerpts.
- **Shows its sources.** The brief says which layer each part came from, so the new agent knows `hot.md` was intended and transcript lines were raw.
- **The new agent confirms first.** It doesn't blindly continue; it states the plan in one line and waits for a go.

---

## 4. Recipes

```
recipes/laravel/
├── recipe.toml       # metadata, detection, installer, variables, MCP servers
├── agents.md         # recipe block template
├── skills/           # stack skills
└── files/            # extra files to copy (optional)
```

`recipe.toml`:
```toml
name = "laravel"
version = 1
description = "Laravel application"

[install]                                   # skipped with --bare or on `baker init`
run = "composer create-project laravel/laravel {{name}}"

[detect]                                    # used by `baker init` on existing folders
file = "composer.json"
contains = "laravel/framework"

[facts]                                     # read after install, fed to templates
php_version  = { from = "composer.json", path = "require.php" }
test_runner  = { if_package = "pestphp/pest", then = "Pest", else = "PHPUnit" }

[[mcp]]
name = "laravel-boost"
command = "php"
args = ["artisan", "boost:mcp"]
```

**Lookup order:** `./.baker/recipes` → `~/.baker/recipes` → built-in → remote (`--recipe github:user/repo`).
So you can override the built-in Laravel recipe with your own taste without forking baker.

---

## 5. Commands

| Command | What happens |
|---|---|
| `baker new <name> --recipe <r>` | mkdir → run recipe installer → detect facts → render AGENTS.md → bridges → skills → memory skeleton → MCP files → `.gitignore` entries → `git init` → write `baker.toml` + lock → register in `~/.baker/projects.toml` |
| `baker init` | Same, minus mkdir/installer, on an existing folder. Detects recipe. Shows a plan; if an `AGENTS.md`/`CLAUDE.md` already exists, wraps existing content as project notes instead of overwriting. |
| `baker sync` | Re-renders bridges, skill mirrors and MCP files from `AGENTS.md` + `baker.toml`. Run after editing. |
| `baker global sync` | Pushes `~/.baker/AGENTS.md` and skills into each tool's user-level files. |
| `baker doctor` | Checks: bridges present and correct, `AGENTS.md` under 32 KiB (Codex cap), no legacy root files, stale memory, broken links, orphan notes, notes over 400 words or missing frontmatter. |
| `baker graph` | Rebuilds the code graph in `.baker/graph/` (only if `[graph]` is enabled). |
| `baker checkpoint` | Takes a mechanical snapshot (git state + latest transcript tail) into `.baker/checkpoints/`. Usually run by agent hooks, not by hand. |
| `baker handoff` | Builds a resume brief from `hot.md` + the latest checkpoint + the last session transcript of any agent. Prints it or copies it to the clipboard. |
| `baker resume <agent>` | `baker handoff`, then launches `<agent>` (claude, codex, gemini, opencode…) with the brief as its first prompt. |
| `baker upgrade` | Updates core/recipe blocks to newer versions. Uses the lock: untouched files are replaced; edited files only get their managed blocks refreshed; edited blocks show a diff and ask. |
| `baker list` | Lists registered projects with recipe and standard version. `baker upgrade --all` walks them. |

---

## 6. Decisions (locked 2026-10-18)

1. **Memory folder:** `memory/` at the project root (visible to humans and Obsidian).
2. **Commit memory:** yes — `MEMORY.md`, `decisions/`, `lessons/`, `log/` are committed; only `hot.md` is gitignored.
3. **Config format:** TOML (`baker.toml`, `recipe.toml`, `~/.baker/config.toml`).
4. **CLI language:** Go. Single static binary, ~instant startup, no runtime needed (works the same for PHP, Python, JS devs). Recipes embedded with `embed.FS`, templates with `text/template`. Distribute via GitHub Releases, Homebrew tap, Scoop, and `go install`.
5. **Name:** `baker`. Checked 2026-10-18:
   - Homebrew: no `baker` formula or cask — command name is free.
   - GitHub: no dominant project named baker in dev tooling (largest is a 2012 HTML5 ebook framework).
   - npm `baker`: taken by an abandoned 2022 placeholder; irrelevant for a Go binary. PyPI `baker`: an old Python library.
   - Domains: `baker.dev` and `baker.sh` are taken; `getbaker.dev` looked unregistered (not confirmed).
   - Risk: a common word, so search/SEO is harder — brand as "baker — scaffolding for AI-ready projects".
