# Continuity

Keep working when a session ends: usage limits, switching accounts, or moving from one provider to another.

## The problem

A usage limit ends a session with no warning. You switch to another account, or from Claude Code to Codex, and the new session knows nothing about the task. You can't count on the old agent writing a goodbye note.

## Three layers

| Layer | Written by | When | Survives a hard cutoff? |
|---|---|---|---|
| `memory/hot.md` | the agent | after every finished step (rule in the AGENTS.md core block) | yes, up to the last step |
| `.baker/checkpoints/*.json` | `baker checkpoint`, run by agent hooks | after every agent turn | yes, up to the last turn |
| Transcript recovery | `baker handoff` | when you run it | yes, everything said |

None of them needs an LLM to read, and none is tied to one vendor.

### 1. `memory/hot.md`

A fixed structure any agent can write and read:

```markdown
---
updated: 2026-10-18T14:32
agent: claude-code
branch: feat/tenant-billing
---
# Goal
# Done
# Next
# Gotchas
# Open questions
```

The core block in `AGENTS.md` tells every agent to update it after each finished step and before risky operations. The `session-close` skill rewrites it at the end of a session. It is gitignored because it's personal scratch.

### 2. Checkpoints

`baker init` / `sync` install a hook so that `baker checkpoint` runs after every turn:

- **Claude Code:** `Stop` and `PreCompact` hooks in `.claude/settings.json`. Claude passes the transcript path on stdin.
- **Gemini CLI:** `AfterAgent` and `PreCompress` hooks in `.gemini/settings.json`. Gemini passes the transcript path on stdin, and asks you to approve a project hook before it first runs.
- **Codex:** `baker global sync` adds a `notify` program to `~/.codex/config.toml`, but only if you don't already have one. Codex passes the turn's `cwd` as a JSON argument, so the checkpoint lands in the right project. (Codex ignores `notify` in project config.)

A checkpoint is mechanical: branch, last commit, `git status`, diff stats, the agent's todo list, and the last 8 messages (redacted, truncated). The newest 20 are kept in `.baker/checkpoints/` (gitignored). With `--quiet`, the command never prints and always exits 0, so a hook can never break your agent. Outside a baker project it does nothing.

### 3. Transcript recovery

Agents already keep full transcripts locally. baker finds the newest one for the current project folder:

| Agent | Location | Match |
|---|---|---|
| Claude Code | `<config>/projects/<path-with-dashes>/*.jsonl` | folder name, falls back to the `cwd` field |
| Codex | `<config>/sessions/YYYY/MM/DD/rollout-*.jsonl` | `cwd` in the session metadata; plan (`update_plan`) and files (`apply_patch`) are picked up too |
| Gemini CLI | `<config>/tmp/<project-id>/chats/session-*.jsonl` | `<project-id>` from `<config>/projects.json` or the `.project_root` marker; older versions used `sha256(path)` and `.json` files, which still work |

**Several accounts?** Each account usually has its own config dir (for example a second `CLAUDE_CONFIG_DIR`). List them all in `~/.baker/config.toml`:

```toml
[transcripts]
claude = ["~/.claude", "~/.claude-work"]
codex  = ["~/.codex"]
gemini = ["~/.gemini"]
```

## Using it

```bash
baker resume codex            # build the brief and launch Codex with it
baker resume claude --from codex
baker handoff                 # print the brief
baker handoff --copy          # copy it, paste into any agent or web chat
baker handoff -o brief.md     # save it
```

The brief contains:

- the hot.md note
- live repo state: branch, last commit, uncommitted files
- the previous agent's todo list
- the files it edited
- the last exchange, always including your latest request
- a "first move" telling the new agent to confirm the plan before continuing

Every section is labeled with its source, so the new agent knows that hot.md was written on purpose and transcript lines are raw.

## Privacy

- Checkpoints and briefs never leave your machine and are gitignored.
- Before anything is written, text that looks like a credential is replaced with `[REDACTED]`: API keys (`sk-…`, `ghp_…`, `AKIA…`, `xox…`, `AIza…`), JWTs, bearer tokens, private keys, and `SECRET=` / `password:` style assignments.
- Redaction is pattern-based and best effort. Don't paste production secrets into agent chats.
