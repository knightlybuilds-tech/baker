---
name: session-close
description: Wrap up a work session so the next agent (any provider, any account) can continue. Use when the user says they're done, stopping, switching agents, or near a usage limit.
---

# Session close

Do these steps in order. Keep every write short.

1. **Update `memory/hot.md`** with this exact structure:
   - frontmatter: `updated` (ISO date-time), `agent` (your tool name), `branch` (current git branch)
   - `# Goal` — one or two sentences
   - `# Done` — checked items `- [x]` finished this session
   - `# Next` — unchecked items `- [ ]`, in order, each specific enough to start without asking (name files)
   - `# Gotchas` — traps the next agent would hit
   - `# Open questions` — things waiting on the user
   Remove Done items older than this session so the file stays small.

2. **Append to `memory/log/<YYYY-MM-DD>.md`** (create if missing) one entry:
   ```
   ## HH:MM — <agent>
   - what changed, in 1–5 bullets
   - decisions made (link the note if you wrote one)
   ```

3. **Propose lessons, don't just write them.** If something was learned the hard way, or a decision was made, list the proposed atomic notes (title as a full-sentence claim + one-line reason) and ask the user which to keep. For each approved one:
   - create `memory/lessons/<slug>.md`, `memory/decisions/<slug>.md` or `memory/knowledge/<slug>.md`
   - frontmatter: `type`, `scope: project`, `tags`, `as_of`, `last_verified` (today), `source` (this session), `confidence`
   - body under 200 words; link related notes with `[[slug]]`
   - add one line to `memory/MEMORY.md` under the right heading

4. Tell the user in one line that the handoff is written and they can continue with any agent (`baker resume <agent>`).

Never write secrets, tokens, passwords or customer data.
