---
name: consolidate
description: Tidy project memory — merge duplicate notes, flag stale facts, keep the index short, and propose promotions. Use weekly, or when the user asks to clean up or review memory.
---

# Consolidate memory

Work through `memory/`. Propose every change first; apply only what the user approves.

1. **Run `baker doctor`** if available and start from its warnings (orphans, broken links, oversize notes, stale notes).
2. **Duplicates:** find notes that make the same claim. Propose merging into one note, keeping the best title, and fix every `[[link]]` that pointed at the removed one.
3. **Split:** a note with more than one claim, or over ~200 words, becomes several atomic notes.
4. **Stale:** for notes whose `last_verified` is older than 90 days, check the claim against the code or docs now. If still true, bump `last_verified`. If false, propose deleting or rewriting it.
5. **Index:** every note has exactly one line in `memory/MEMORY.md`. When a section passes ~30 lines, propose a hub note (e.g. `memory/knowledge/_billing.md`) that lists that topic's notes, and replace those lines with one link to the hub.
6. **Promote:**
   - A lesson that should always apply in this project → propose it as a rule in the user's section of `AGENTS.md`.
   - A lesson true across projects → propose copying it to `~/.baker/memory/lessons/` with `scope: global`.
7. **Log:** old daily files in `memory/log/` can stay; never rewrite history there.

Report a short summary of what changed.
