# Agent support

What baker writes for each agent, and why. File locations come from vendor docs as of October 2026. Agent tools change fast, so items marked **verify** should be re-checked before relying on them. The registry lives in `internal/agents/agents.go`.

## Instructions

`AGENTS.md` is the single source. It is stewarded by the Linux Foundation's Agentic AI Foundation and read natively by most tools.

| Agent | How it gets AGENTS.md | Notes |
|---|---|---|
| Claude Code | `CLAUDE.md` containing `@AGENTS.md` | Claude Code skips AGENTS.md once a CLAUDE.md exists, so the import is required. Imports work better than symlinks: Claude's edit tools don't write through symlinks, and symlinks break on Windows. |
| Gemini CLI | `GEMINI.md` containing `@AGENTS.md` | Gemini can also be configured via `context.fileName`; the import works without config. |
| Aider | `.aider.conf.yml` with `read: AGENTS.md` | Native AGENTS.md support is disputed between sources, so the config is the safe choice. |
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
| OpenCode | `opencode.json` | `{"mcp": {name: {type: "local", command: [...]}}}` |
| Codex | `.codex/config.toml` | `[mcp_servers.name]` in a managed block. **Verify** that your Codex version reads project-level config; `~/.codex/config.toml` always works. |

baker merges into existing JSON files and only adds or removes servers it registered (tracked in `.baker/lock.toml`).

## User-level instructions (`baker global sync`)

| Agent | File | |
|---|---|---|
| Claude Code | `~/.claude/CLAUDE.md` (+ skills → `~/.claude/skills/`) | |
| Codex | `~/.codex/AGENTS.md` | |
| Gemini CLI | `~/.gemini/GEMINI.md` | |
| OpenCode | `~/.config/opencode/AGENTS.md` | **verify** |
| Cursor | none: user rules live in the settings UI | |

Only agents whose config folder exists are touched.

## Checkpoint hooks

| Agent | Hook | Status |
|---|---|---|
| Claude Code | `Stop` + `PreCompact` in `.claude/settings.json` | ✅ |
| Codex | `notify` in `~/.codex/config.toml` (via `baker global sync`) | ✅. Skipped if you already have a `notify` |
| Gemini CLI | hooks in `.gemini/settings.json` | not wired yet; **verify** event names |
| others | none | hot.md and transcript recovery still work for Claude/Codex/Gemini transcripts |

## Launch commands (`baker resume`)

| Agent | Command |
|---|---|
| claude | `claude "<brief>"` |
| codex | `codex "<brief>"` |
| gemini | `gemini -i "<brief>"` |
| opencode | `opencode --prompt "<brief>"` (**verify**) |
| aider | `aider --message "<brief>"` |

For anything else, use `baker handoff --copy` and paste.

## Adding an agent

Add an entry to `internal/agents/agents.go` with its bridge (if it doesn't read AGENTS.md), MCP file and format, global file, and launch command. Then add a row to the tables above and to the README.
