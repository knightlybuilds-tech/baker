# Memory

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

## Skills that maintain memory

| Skill | Use it when |
|---|---|
| `session-close` | Ending a session or switching agents: updates hot.md, writes the daily log, proposes lessons |
| `consolidate` | Weekly cleanup: merge duplicates, re-verify stale notes, build hubs, propose promotions |
| `ingest` | Turning a doc, URL or PDF into notes with sources |

Agents **propose** notes and promotions, and you approve them. A lesson becomes a rule in `AGENTS.md`, or moves to `~/.baker/memory/` (global), only with your OK.

## Global memory

`~/.baker/memory/` follows the same rules, for lessons that hold across projects. `baker global sync` tells every agent where it is.
