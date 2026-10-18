# Configuration

## Project: `baker.toml`

Created by `baker new` / `baker init`. Commit it.

```toml
standard = '0.1'              # baker standard version the project follows
name = 'shop'
recipe = 'laravel'
recipe_version = 1
agents = ['claude', 'codex', 'gemini', 'cursor', 'copilot', 'opencode']

[facts]                       # detected from project files; refreshed on every sync
laravel_version = '13'
php_version = '8.3'
test_runner = 'Pest'

[[mcp]]                       # add your own servers here, then `baker sync`
name = 'laravel-boost'
command = 'php'
args = ['artisan', 'boost:mcp']
```

- **Change agents:** edit `agents`, run `baker sync`. All ids: `baker recipes`.
- **Add an MCP server:** add an `[[mcp]]` table, run `baker sync`. It's written to every enabled agent's MCP file.
- **Remove an MCP server:** delete its table and run `baker sync`. baker removes it from the files it manages.

## Lock: `.baker/lock.toml`

Short hashes of every file and block baker generated. `sync` compares them with what's on disk: if you edited something, baker leaves it alone and says so. `--force` overwrites. Commit the lock so teammates get the same protection.

## Managed blocks

Inside shared files, baker only touches text between its markers:

```markdown
<!-- baker:begin core@0.1 -->
…baker owns this…
<!-- baker:end core -->
```

YAML, TOML and `.gitignore` use `# baker:begin …` / `# baker:end …`. Everything outside the markers is yours.

## Global: `~/.baker/`

Created by `baker global init`. You can keep it in your dotfiles repo.

```
~/.baker/
├── config.toml      defaults (below)
├── AGENTS.md        your personal rules; pushed to every agent by `baker global sync`
├── skills/          personal skills; copied to ~/.claude/skills
├── recipes/         your recipes and overrides
├── memory/          cross-project memory (open as an Obsidian vault)
├── projects.toml    registry used by `baker list`
└── lock.toml        protects edits in agents' user-level files
```

`config.toml`:

```toml
agents = ['claude', 'codex', 'gemini', 'cursor', 'copilot', 'opencode']   # default for new projects
default_recipe = 'generic'

[transcripts]        # config dirs to search for transcripts, one per account
claude = ['~/.claude', '~/.claude-work']
codex = ['~/.codex']
gemini = ['~/.gemini']
```

## Environment variables

| Variable | Effect |
|---|---|
| `BAKER_HOME` | Use another folder instead of `~/.baker` |
| `BAKER_USER_HOME` | Treat another folder as your home (agents' config dirs, transcripts). Used in tests. |
