# `baker learn` — design (not built)

> Status: proposal, 2026-11-13. Nothing here is implemented. It records the
> intended shape, the known holes, and the order to build it in.

## Goal

Give a project that has never used Baker a useful first set of crumbs, so the
first agent session after adoption doesn't start from zero.

`baker learn` reads the project's history (and, if the developer opts in, past
agent sessions in this repo), asks an agent the developer already uses to
propose crumbs **with evidence**, and puts them in an **inbox** for a human to
keep or drop. Nothing reaches `memory/` without that triage.

## Principles

- **The why, not the what.** Agents can read the code. A crumb is worth keeping
  when it records a decision, a gotcha or a lesson the code doesn't show.
  "Uses PostgreSQL" isn't a crumb.
- **No source, no crumb.** Every draft cites a commit, `file:line`, doc or
  transcript excerpt.
- **Honest freshness.** `as_of` is the evidence's date. `last_verified` stays
  empty until a person (or a later check) confirms it.
- **Baker writes, the agent reads.** The learning agent runs read-only and
  returns structured drafts. Baker validates them and writes files through the
  usual `scaffold.Plan` → `Apply` path.
- **Bulk drafts get triage. Ordinary notes still don't.** The no-approval rule
  for hand-written notes stays. The inbox exists because machine-written notes
  arrive in bulk and can be wrong or planted.
- **No new AI dependency.** Baker doesn't call model APIs or hold keys. It
  drives an agent CLI the developer already has installed and signed in.

## Flow

```
$ baker learn
This project isn't set up with Baker yet.
Plan: create memory/ (MEMORY.md, decisions/, lessons/, knowledge/, inbox/) and baker.toml
Apply? [y/N] y

Agents found on this machine:
  1) Claude Code   installed · 14 past sessions in this repo
  2) Codex         installed · 3 past sessions in this repo
  3) Gemini        installed · no sessions here
  4) Aider         not installed
Who should do the learning? [1] 1

Read past agent sessions for lessons? (stays on this machine, secrets redacted)
  [ ] Claude Code (14)   [ ] Codex (3)
Sources: git history (last 12 months, 1,240 commits) · 38 HACK/NOTE comments · docs/ · 2 ADRs
Continue? [y/N] y

Learning with Claude Code… 22 drafts in memory/inbox/.
Triage them:  baker learn review           (terminal)
         or:  baker learn review --web     (local Bakery, Keep / Edit / Drop)
```

Non-interactive: `baker learn --agent claude --transcripts claude --since 12m --yes`.
Without a terminal, `--agent` is required and Baker never guesses.

### Projects not set up with Baker

`learn` creates a **memory-only** setup (`memory/`, `MEMORY.md`, `baker.toml`)
through a normal plan, then suggests `baker init` for bridges, skills and hooks.
It doesn't install hooks or touch agent config on its own.

## Steps

### 1. Harvest evidence (deterministic, no AI)

Baker builds a bounded evidence pack in `.baker/learn/` (gitignored):

| Source | Signals |
|---|---|
| Git history | reverts, "fix the fix" chains, files changed again and again, long commit bodies, commits touching config/CI/migrations |
| Pull requests and issues | via `gh` if installed and signed in; titles, descriptions, review threads on hot files |
| Code comments | `HACK`, `NOTE`, `WORKAROUND`, `FIXME`, `XXX` with surrounding lines |
| Docs | README, `docs/`, ADRs, CHANGELOG, CONTRIBUTING |
| Past agent sessions (opt-in) | user corrections, errors that were resolved, "don't do X" moments; parsed with the existing transcript readers and redacted |

Bounds: a time window (`--since`, default 12 months), top N files by churn, a
total size cap. Before running, Baker shows what it will read and a rough size.

### 2. Agent pass (read-only)

Baker launches the chosen agent non-interactively with a `learn` prompt and the
evidence pack. The agent may read repository files to check claims against
`HEAD`. It returns JSON drafts on stdout:

```json
{"kind": "lesson", "title": "Retry payments only on 429 and 503",
 "body": "…", "why": "…",
 "evidence": [{"type": "commit", "ref": "a1b2c3d", "date": "2026-03-02"},
              {"type": "file", "ref": "billing/retry.go:41"}],
 "confidence": "inferred"}
```

Each agent's non-interactive and read-only modes differ (Claude Code `-p` with
read-only tools, Codex `exec` with a read-only sandbox, Gemini `-p`, OpenCode
`run`, …). They belong in `internal/agents` as a per-agent `Learn` adapter and
must be checked against the installed version. An agent without a reliable
read-only mode isn't offered.

### 3. Validate and stage

Baker drops or flags drafts that:
- cite no evidence, or cite commits or files that don't exist;
- duplicate an existing crumb (slug or title match);
- read like commands ("always…", "never run…", "disable…"). These are flagged
  as possible instruction injection, not silently kept;
- exceed the cap (default 25 per run, ranked by evidence strength).

Survivors are written to `memory/inbox/<slug>.md` with frontmatter:

```yaml
type: lesson
origin: learned
learned_by: claude-code
confidence: inferred
as_of: 2026-03-02          # date of the newest cited evidence
source: commit a1b2c3d; billing/retry.go:41
evidence_kinds: [git, transcript]
```

`memory/inbox/` isn't indexed in `MEMORY.md`, so agents don't treat drafts as
memory. It can be committed (to triage as a team in a PR) or gitignored.

### 4. Triage

**Terminal:** `baker learn review` walks the inbox one draft at a time:
Keep / Edit (opens `$EDITOR`) / Drop / Skip.

**Browser (opt-in):** `baker learn review --web` serves the local Bakery with an
**Inbox** view: each draft as a card with its evidence, flags and the existing
crumbs it might duplicate, plus Keep / Edit / Drop.

**Keep** moves the draft into `memory/<kind>s/`, removes `origin`/`confidence`
flags the person accepted, sets `last_verified` to today, and adds the line to
`MEMORY.md`, all through a plan. **Drop** removes it (or moves it to
`.baker/learn/dropped/` with `--keep-dropped`).

The words are deliberately **Keep / Drop**, not "approve". This is a local
triage decision by whoever runs the command, not authenticated team approval.
Team approval still happens in Git pull requests until a hosted Bakery exists.

## Writable local Bakery: the exception

`baker bakery` stays read-only. `baker learn review --web` is a separate,
short-lived mode with the smallest write surface that works:

- **Scope:** three actions only (keep, edit, drop) on files inside
  `memory/inbox/`. Paths are resolved and confined; anything else is refused.
- **CSRF:** a loopback bind and `Host` check aren't enough once writes exist,
  because any website can send a simple `POST` to `127.0.0.1`. Every write
  requires a random per-session token in a custom header (the custom header
  forces a CORS preflight, which the server refuses), plus an `Origin` check.
  The token reaches the page only through the URL Baker prints and opens.
- **Lifetime:** the server exits when the inbox is empty or on Ctrl+C. It's
  never left running in the background.
- **Same writer:** browser actions call the same plan code as the terminal
  review, so both paths behave the same and are tested once.

When this ships, the decision note
`memory/decisions/bakery-local-is-a-read-only-view-served-by-baker.md` should
be updated: still true for `baker bakery`, with the inbox exception named.

## Known holes

1. **Wrong crumbs are worse than none.** Future agents trust memory. Mitigated
   by required evidence, checking against `HEAD`, `confidence: inferred`, and
   triage, but not eliminated.
2. **Persistent prompt injection.** Repo text, issues or old transcripts can
   contain instructions that would become "trusted" crumbs. Mitigated by the
   read-only agent, flagging command-like drafts, and never accepting
   automatically.
3. **Transcripts are shorter-lived than expected.** Agents prune old local
   sessions, and sessions are keyed by folder path. Moved or renamed repos and
   worktrees can be missed unless matched by git remote too.
4. **Privacy.** History and transcripts can hold secrets, customer data and
   coworkers' names. Transcript sources are opt-in, redaction is reused, and
   drafts note which ones came from transcripts.
5. **Cost and time.** It uses the developer's agent quota. Big repos need the
   bounds above and a preview before running.
6. **"Installed" ≠ "usable".** Baker can't cheaply tell whether an agent is
   signed in or out of quota. Failures must be clear and let the user pick
   another agent.
7. **Repeat runs.** A cursor in `.baker/learn/state.toml` (last commit and
   transcript timestamps processed) keeps later runs incremental. Dedupe
   prevents clobbering crumbs a person has edited since.
8. **Unproven benefit.** Needs an evaluation (below) before it's advertised.

## Evaluation

Reuse `scripts/eval_agents.py`. On a few real repos, ask a fresh agent a fixed
set of project questions and tasks:

- with no memory;
- with a hand-written HANDOFF.md;
- with crumbs kept from `baker learn`.

Measure corrections needed, missed constraints, repeated mistakes, time to a
correct action, and triage effort (drafts kept vs dropped). Run each several
times.

## Build order

1. **`learn` skill only.** A prompt devs trigger inside their own agent, writing
   cited drafts to `memory/inbox/`. Almost no Go code, and it tells us quickly
   whether the drafts are any good.
2. **Evidence harvest** in Go (`internal/learn`), with fixture repos and tests.
3. **`baker learn`:** agent detection and chooser, per-agent read-only `Learn`
   adapters, validation and staging, cursor.
4. **`baker learn review`** in the terminal.
5. **`baker learn review --web`** with the scoped, token-protected writes above.
6. **Evaluation** before any public claims.

## Open questions

- Commit `memory/inbox/` by default (team triage in a PR) or gitignore it
  (solo triage)?
- Should Keep require at least one evidence item that still exists at `HEAD`?
- Default time window and cap: 12 months / 25 drafts?
- Include `gh` pull requests and issues by default when available, or opt-in
  like transcripts?
