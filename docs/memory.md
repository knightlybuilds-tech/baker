# Memory

Local memory works without Bakery or approval. Use `baker memory list --json`
to inventory notes, `baker memory recall --query "topic" --json` to retrieve
matching notes, and `baker memory check` for index, link, and freshness checks.
Recall flags missing or stale metadata and keeps those notes available as leads
to verify. No account, network service, or experimental review record is needed.

The earlier [review lifecycle prototype](reviewed-memory.md) is isolated under
`baker experimental crumb` and is not part of the default note workflow.

Project memory lives in `memory/` as small markdown notes. Any agent can read and write it, and humans can browse it in an editor or Obsidian.

## Layout

```
memory/
├── MEMORY.md        index: one line per note; agents read it at session start
├── hot.md           live handoff note (see continuity.md), gitignored
├── log/             one file per day, append-only (no merge conflicts in teams)
├── decisions/       why we chose X
├── lessons/         what we learned the hard way
└── knowledge/       domain facts ("a Tenant owns many Stores")
```

Everything except `hot.md` is committed. Memory is team knowledge.

## Atomic note rules

1. **One idea per note.** If it makes two claims, make two notes.
2. **The title is the claim**, written as a full sentence. The filename is its slug: `lessons/queue-workers-must-restart-after-deploy.md`.
3. **Small:** aim for under 200 words. `baker doctor` warns above 400.
4. **Link, don't repeat:** `[[other-slug]]`.
5. **Indexed:** one line in `MEMORY.md`, or listed in a linked hub note.
6. **Fresh:** frontmatter with `type`, `as_of`, `last_verified`, `source`.

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

Workers keep old code in memory, so new jobs ran stale logic after the 10/12 deploy.

**Apply:** the deploy script runs `php artisan queue:restart`. See [[deploys-use-envoyer]].
```

## Why small notes

- **Cheaper context:** the agent reads a short index and opens 2–3 notes, not one big file.
- **Trustworthy:** each claim has its own date and source, so it can be checked, updated or deleted on its own.
- **A free graph:** wikilinks connect notes, and Obsidian's graph view shows them.

## Hub notes

When an index section grows past ~30 lines, the `consolidate` skill proposes a hub note such as `knowledge/_billing.md` that lists that topic's notes. `MEMORY.md` then links the hub instead. Notes listed in a linked hub count as indexed. Hubs (filenames starting with `_`) are exempt from the word limit.

## What `baker doctor` checks

| Check | Level |
|---|---|
| `MEMORY.md` missing | error |
| Index links to a file that doesn't exist | error |
| Note not in the index or a linked hub | warning |
| Missing `type`, `as_of`, `last_verified` or `source` | warning |
| Over 400 words | warning |
| `last_verified` older than 90 days | warning |
| `[[link]]` to a note that doesn't exist | warning |

## Saving a crumb: `baker memory add`

```bash
baker memory add --kind decision \
  --title "Account codes must remain strings" \
  --body "007 and 7 are different accounts; exports must keep leading zeros." \
  --source "tests/export_test.go" --tags billing --link money-in-minor-units
```

| Flag | |
|---|---|
| `--kind` | `decision`, `lesson` or `knowledge` (required) |
| `--title` | the claim as one sentence; becomes the slug (required) |
| `--body` / `--body-file <path or ->` | detail and why it matters; warns over 200 words, refuses over 400 |
| `--source` | file, commit, issue, doc or `session YYYY-MM-DD` (required) |
| `--tags`, `--link` | tags; related note slugs (written as `[[slug]]`) |
| `--as-of`, `--confidence` | date the fact describes (default today); `high`/`medium`/`low` (default medium) |
| `--hook` | the short phrase after the `MEMORY.md` link (default: the body's first sentence) |
| `--dry-run`, `--json` | validate without writing; machine-readable output |

It refuses duplicate slugs and anything that looks like a secret, warns about
near-duplicate titles and links to notes that don't exist, writes the file with
full frontmatter, and adds the index line under the right heading. Ordinary
notes need no approval.

## Getting agents to save

Left to "the end of the session", agents often save nothing. So:

- The `AGENTS.md` memory block tells agents to save **when it happens**: a
  decision and its reason, a bug's cause or a trap, a project fact the code
  doesn't show, a correction from the user.
- **Claude Code** also gets a Stop hook, `baker nudge`. After a session edits 3+
  files with no crumb saved (and again after every 8 more), it blocks the stop
  once with a reason asking the agent to save crumbs or reply "No new crumbs."
  It never fires while a previous nudge is being answered (`stop_hook_active`),
  keeps its per-session state in `.baker/nudge/` (gitignored), and never fails.
  Turn it off with `nudge = false` in `baker.toml` and `baker sync`.
- Codex and Gemini have no equivalent hook reply yet; they rely on the
  instructions and the `session-close` skill.

## Skills that maintain memory

| Skill | Use it when |
|---|---|
| `session-close` | Ending a session or switching agents: updates hot.md, writes the daily log, saves crumbs the session earned |
| `consolidate` | Weekly cleanup: merge duplicates, re-verify stale notes, build hubs, propose promotions |
| `ingest` | Turning a doc, URL or PDF into notes with sources |
| `learn` | Seeding memory in a project with history: proposes cited crumbs from git, comments and docs (and past agent sessions if you agree), saves the ones you keep |

Agents save ordinary notes directly. Merges, deletions and promotions are proposed, and you approve them. A lesson becomes a rule in `AGENTS.md`, or moves to `~/.baker/memory/` (global), only with your OK.

## Global memory

`~/.baker/memory/` follows the same rules, for lessons that hold across projects. `baker global sync` tells every agent where it is.
