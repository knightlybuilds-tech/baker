---
type: decision
scope: project
tags: [memory, capture, hooks]
as_of: 2026-11-17
last_verified: 2026-11-17
source: session 2026-11-17; internal/cli/nudge.go, internal/memory/add.go
confidence: high
---
# Agents save crumbs through baker memory add, nudged by a Claude Stop hook

Agents often skipped saving lessons when it was left to the end of a session or to a propose-then-ask flow. So capture is immediate: the AGENTS.md memory block lists concrete triggers and points at `baker memory add`, which validates format, source, duplicates and secrets before writing and indexing. For Claude Code, `baker nudge` runs on the Stop hook and, after a stretch of edits with no crumb saved, blocks the stop once with a reason asking the agent to save or say "No new crumbs." It never fires while stop_hook_active is set, and nudge = false in baker.toml removes it.

Related: [[bakery-local-is-a-read-only-view-served-by-baker]]
