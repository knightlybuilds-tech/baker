# Agent support

What baker writes for each agent, and why. File locations, hook events and transcript formats were checked against vendor docs and source in October 2026 (Codex 0.115 transcripts also against real files). Agent tools change fast; anything still marked **verify** is unconfirmed. The registry lives in `internal/agents/agents.go`.

## Instructions

`AGENTS.md` is the single source. It is stewarded by the Linux Foundation's Agentic AI Foundation and read natively by most tools.

| Agent | How it gets AGENTS.md | Notes |
|---|---|---|
| Claude Code | `CLAUDE.md` containing `@AGENTS.md` | Claude Code skips AGENTS.md once a CLAUDE.md exists, so the import is required. Imports work better than symlinks: Claude's edit tools don't write through symlinks, and symlinks break on Windows. |
| Gemini CLI | `GEMINI.md` containing `@AGENTS.md` | Gemini can also be configured via `context.fileName`; the import works without config. |
| Aider | `.aider.conf.yml` with `read: AGENTS.md` | Aider (0.86) doesn't read AGENTS.md on its own; `read:` loads it read-only every session. |
| Codex | native | Concatenates AGENTS.md files root → cwd, capped at 32 KiB by default. `doctor` warns. |
| Cursor, Copilot, OpenCode, Windsurf, Amp, Zed, Cline | native | |

**Legacy files hide AGENTS.md.** Zed reads only the first match from a list where `.rules`, `.cursorrules`, `.windsurfrules` and `.clinerules` come before AGENTS.md. baker never writes them, and `doctor` warns when they exist.

## Skills

Canonical: `.agents/skills/<name>/SKILL.md` ([Agent Skills](https://agentskills.io) standard). Claude Code doesn't read `.agents/`, so baker mirrors skills to `.claude/skills/`. The mirror is checked by `doctor` and refreshed by `sync`.

## MCP

Servers are declared once in `baker.toml` (or by a recipe) and written in each tool's format:

| Agent | File | Shape |
|---|---|---|
| Claude Code | `.mcp.json` | `{"mcpServers": {name: {command, args, env}}}` |
| Cursor | `.cursor/mcp.json` | same as Claude |
| Gemini CLI | `.gemini/settings.json` | `mcpServers` key, merged with your other settings |
| Copilot (VS Code) | `.vscode/mcp.json` | `{"servers": {name: {type: "stdio", …}}}` |
| OpenCode | `opencode.json` | `{"mcp": {name: {type: "local", command: [...], environment}}}` |
| Codex | `.codex/config.toml` | `[mcp_servers.name]` in a managed block. Codex loads project config only once you trust the project (it asks the first time); `~/.codex/config.toml` always works. |

baker merges into existing JSON files and only adds or removes servers it registered (tracked in `.baker/lock.toml`).

## User-level instructions (`baker global sync`)

| Agent | File | |
|---|---|---|
| Claude Code | `~/.claude/CLAUDE.md` (+ skills → `~/.claude/skills/`) | |
| Codex | `~/.codex/AGENTS.md` | |
| Gemini CLI | `~/.gemini/GEMINI.md` | |
| OpenCode | `~/.config/opencode/AGENTS.md` | Without it, OpenCode falls back to `~/.claude/CLAUDE.md` |
| Cursor | none: user rules live in the settings UI | |

Only agents whose config folder exists are touched.

## Checkpoint hooks

| Agent | Hook | Status |
|---|---|---|
| Claude Code | `Stop` + `PreCompact` in `.claude/settings.json` | ✅ `{transcript_path, cwd}` on stdin |
| Gemini CLI | `AfterAgent` + `PreCompress` in `.gemini/settings.json` | ✅ `{transcript_path, cwd}` on stdin. Gemini asks before running a project hook the first time, and again if the command changes |
| Codex | `notify` in `~/.codex/config.toml` (via `baker global sync`) | ✅ JSON payload (with `cwd`) as the last argument. Codex ignores `notify` in project config. Skipped if you already have a `notify` |
| others | none | hot.md and transcript recovery still work for Claude/Codex/Gemini transcripts |

The hook entries come from `HooksFile` / `HookEvents` in the registry. Gemini parses hook stdout as JSON, so the command must print nothing; `--quiet` guarantees that.

## Launch commands (`baker resume`)

| Agent | Command |
|---|---|
| claude | `claude "<brief>"` |
| codex | `codex "<brief>"` |
| gemini | `gemini -i "<brief>"` |
| opencode | `opencode --prompt "<brief>"` (**verify** whether it submits or only pre-fills; the flag exists) |
| aider | `aider --message "<brief>"` |

For anything else, use `baker handoff --copy` and paste.

## Adding an agent

Add an entry to `internal/agents/agents.go` with its bridge (if it doesn't read AGENTS.md), MCP file and format, global file, hook file and events (if it has Claude-style JSON hooks), and launch command. Then add a row to the tables above and to the README.
