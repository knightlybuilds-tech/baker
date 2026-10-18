---
type: decision
scope: project
tags: [standards, agents]
as_of: 2026-10-18
last_verified: 2026-10-18
source: design session 2026-10-18
confidence: high
---
# AGENTS.md is the single source of instructions

Every agent-specific file (CLAUDE.md, GEMINI.md, .aider.conf.yml) is a generated bridge that imports AGENTS.md. AGENTS.md is stewarded by the Linux Foundation's Agentic AI Foundation and read natively by most tools, so it is the most portable choice.

Bridges use `@AGENTS.md` imports, not symlinks: Claude Code's edit tools don't write through symlinks and symlinks break on Windows clones. See [[legacy-rules-files-hide-agents-md]].
