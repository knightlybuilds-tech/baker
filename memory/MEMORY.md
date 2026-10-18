# Memory index

<!-- One line per note: - [Title](path) — hook
     Agents read this at session start, then open only the notes they need.
     When a section passes ~30 lines, move it into a hub note (e.g. knowledge/_billing.md) and link the hub. -->

## Decisions
- [AGENTS.md is the single source of instructions](decisions/agents-md-is-the-single-source.md) — bridges import it; no symlinks
- [baker is a single Go binary](decisions/baker-is-a-single-go-binary.md) — no runtime, instant start for hooks
- [Project memory is committed, except hot.md](decisions/project-memory-is-committed-except-hot-md.md) — daily log files avoid conflicts

## Lessons
- [Legacy rules files hide AGENTS.md from some tools](lessons/legacy-rules-files-hide-agents-md.md) — Zed reads .cursorrules etc. first
- [Laravel 13 ships CLAUDE.md as a copy of AGENTS.md](lessons/laravel-13-ships-claude-md-as-a-copy-of-agents-md.md) — baker dedupes it

## Knowledge
