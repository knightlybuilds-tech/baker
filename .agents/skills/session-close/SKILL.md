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

3. **Save the crumbs this session earned.** Look back for decisions (and why), lessons learned the hard way, project facts the code doesn't show, and corrections from the user. Skip anything already saved, obvious from the code, or only about task progress. For each one:
   ```
   baker memory add --kind <decision|lesson|knowledge> --title "<the claim as one sentence>" \
     --body "<detail and why it matters, under 200 words>" --source "session <YYYY-MM-DD>" [--link <slug>]
   ```
   Without Baker, write the note by hand (one claim per file in `memory/<kind>s/`, frontmatter `type`, `as_of`, `last_verified`, `source`, a line in `memory/MEMORY.md`). Ordinary notes need no approval; list what you saved in one line. If nothing qualifies, say so.

4. Tell the user in one line that the handoff is written and they can continue with any agent (`baker resume <agent>`).

Never write secrets, tokens, passwords or customer data.
