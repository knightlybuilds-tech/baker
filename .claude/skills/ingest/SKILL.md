---
name: ingest
description: Turn a document, URL, PDF, meeting note or pasted text into atomic memory notes with sources. Use when the user shares material and asks to remember, learn, or save it.
---

# Ingest into memory

1. **Read the source fully.** Note where it came from (URL, file path, "pasted by user on <date>").
2. **Extract claims** that matter for this project: decisions, domain facts, constraints, gotchas. Skip anything generic or already in memory (check `memory/MEMORY.md` first).
3. **Propose** the list of notes to the user: one line each, title as a full-sentence claim. Wait for approval.
4. **Write each approved note** as `memory/knowledge/<slug>.md` (or `decisions/`, `lessons/`):
   ```markdown
   ---
   type: knowledge
   scope: project
   tags: []
   as_of: <date the source describes>
   last_verified: <today>
   source: <url or path>
   confidence: high | medium | low
   ---
   # <The claim as a full sentence>

   <Under 200 words: the detail, and why it matters here. Link related notes with [[slug]].>
   ```
5. **Index:** add one line per note to `memory/MEMORY.md`.
6. If the source is large and will be needed again, save it (or a link to it) under `memory/sources/` and reference that path in `source`.

Never copy secrets or personal data from the source.
