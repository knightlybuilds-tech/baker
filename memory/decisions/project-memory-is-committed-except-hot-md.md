---
type: decision
scope: project
tags: [memory, git]
as_of: 2026-10-18
last_verified: 2026-10-18
source: design session 2026-10-18
confidence: high
---
# Project memory is committed, except hot.md

`memory/` is team knowledge, so decisions, lessons, knowledge and daily logs are committed. Logs are one file per day to avoid merge conflicts. `hot.md` is personal, in-flight handoff state and is gitignored.
