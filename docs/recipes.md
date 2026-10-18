# Recipes

A recipe teaches baker a stack: how to install it, how to recognize it, what agents should know about it, and which MCP servers and skills it needs.

## Where recipes come from

Lookup order, first match wins:

1. `<project>/.baker/recipes/<name>/`: a project-specific recipe
2. `~/.baker/recipes/<name>/`: your own recipes, or overrides of built-in ones
3. Built-in (embedded in the binary): `generic`, `laravel`

So to make the Laravel recipe match your taste, copy `assets/recipes/laravel/` to `~/.baker/recipes/laravel/` and edit it.

## Folder layout

```
laravel/
├── recipe.toml      metadata, installer, detection, facts, MCP servers
├── agents.md        template for the managed recipe block in AGENTS.md
├── notes.md         optional: extra starter text for the user-owned Project notes
└── skills/          optional: skills/<name>/SKILL.md (Agent Skills format)
```

## recipe.toml

```toml
name = "laravel"
version = 1                    # bump when agents.md changes; shown in the block marker
description = "Laravel application"

[install]                      # optional; skipped by --bare and by `baker init`
run = "composer create-project laravel/laravel {{.Name}}"   # runs in the parent dir, must create {{.Name}}

[detect]                       # used by `baker init` to pick a recipe
file = "composer.json"
contains = "laravel/framework"

# Facts are read from project files after install, then used in templates.
[[facts]]
name = "php_version"
file = "composer.json"
json_path = "require.php"      # dot path into JSON
format = "version"             # "version": ^8.4 → 8.4   "major": ^13.0 → 13

[[facts]]
name = "test_runner"
file = "composer.json"
if_contains = "pestphp/pest"   # or: test whether the file contains a string
then = "Pest"
else = "PHPUnit"

# MCP servers are registered for every enabled agent, optionally only when a package is present.
[[mcp]]
name = "laravel-boost"
command = "php"
args = ["artisan", "boost:mcp"]
when_file = "composer.json"
when_contains = "laravel/boost"

gitignore = []                 # extra lines for the baker .gitignore block
```

## Templates

`agents.md`, `notes.md` and `install.run` are Go [text/template](https://pkg.go.dev/text/template)s. They can use:

| Field | Value |
|---|---|
| `{{.Name}}` | project name |
| `{{.Facts.<name>}}` | a resolved fact (empty string if missing) |

```
Laravel{{with .Facts.laravel_version}} {{.}}{{end}} on PHP {{.Facts.php_version}}. Tests: {{.Facts.test_runner}}.
{{- if not .Facts.boost}}
Tip: install Laravel Boost…
{{- end}}
```

## What goes in agents.md

The recipe block is loaded in **every** session, so keep it short and specific:

- **Stack:** versions and key packages (from facts)
- **Commands:** exactly how to run, test, format and migrate
- **Conventions:** the 5–10 rules an agent would otherwise get wrong

Don't put long docs here; link them or make a skill. Codex truncates `AGENTS.md` past 32 KiB.

`agents.md` is managed by baker. If a user edits the block in their project, `baker sync` leaves it alone. Put anything users should customize in `notes.md`, which is written once and then belongs to them.

## Skills

Recipe skills are copied to `.agents/skills/` and mirrored to `.claude/skills/`. Use the [Agent Skills](https://agentskills.io) format:

```markdown
---
name: laravel-migration
description: Create or change database tables safely. Use for any schema change.
---
# Steps…
```

## Testing a recipe

```bash
mkdir /tmp/try && cd /tmp/try
baker new demo --recipe myrecipe --bare   # skip the installer
cat demo/AGENTS.md
```
