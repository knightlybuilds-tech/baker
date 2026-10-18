---
type: lesson
scope: project
tags: [agents, zed]
as_of: 2026-10-18
last_verified: 2026-10-18
source: design session 2026-10-18
confidence: high
---
# Legacy rules files hide AGENTS.md from some tools

Zed reads only the first match from a list where `.rules`, `.cursorrules`, `.windsurfrules` and `.clinerules` come before AGENTS.md. baker never writes those files and `baker doctor` warns when they exist. See [[agents-md-is-the-single-source]].
