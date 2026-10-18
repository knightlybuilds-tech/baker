# Cross-Vendor Standards and File Conventions for Agent Context, Instructions, Skills, and Memory (as of Oct 2026)

Research date: 2026-10-18. Primary vendor docs are preferred. Where only third-party guides were found, that is marked "(secondary)". Claims marked "UNVERIFIED" could not be confirmed against a primary source.

## AGENTS.md: governance, native readers, nesting rules

### Takeaway
AGENTS.md is the closest thing to a universal instruction file. It is stewarded by the Agentic AI Foundation (AAIF), a Linux Foundation directed fund formed in December 2025, and in 2026 it is read natively by almost every tool on the target list, including Claude Code (since v2.1.277). Gemini CLI reads it only through config, and Aider's default behavior is disputed. The baseline rule is "nearest file wins". Tools disagree on whether files concatenate, and only Codex documents a size cap.

### Cited Findings
- **Governance:** AGENTS.md is "stewarded by the Agentic AI Foundation under the Linux Foundation" and is "a Series of LF Projects, LLC". It came out of collaboration between OpenAI Codex, Amp, Google Jules, Cursor, and Factory. — [agents.md](https://agents.md/)
- **AAIF formation:** AAIF was announced on 2025-12-09. Its founding projects are Anthropic's MCP, Block's goose, and OpenAI's AGENTS.md. It operates as a directed fund within the Linux Foundation. Platinum members: AWS, Anthropic, Block, Bloomberg, Cloudflare, Google, Microsoft, OpenAI. — [Linux Foundation press release](https://linuxfoundation.org/press/linux-foundation-announces-the-formation-of-the-agentic-ai-foundation); [InfoQ](https://infoq.com/news/2025/12/agentic-ai-foundation/)
- **Release and scale:** OpenAI released AGENTS.md in August 2025 and claims more than 60,000 adopting projects. This is a vendor figure. — [YourStory](https://yourstory.com/ai-story/agentic-ai-foundation-openai-anthropic-linux)
- **Tools agents.md lists as supporting it:** Codex, Jules, Factory, Aider, goose, opencode, Zed, Warp, VS Code, Devin, UiPath, Junie, Amp, Cursor, RooCode, Gemini CLI, Kilo Code, Phoenix, Semgrep, GitHub Copilot coding agent, Ona, Windsurf, Augment Code. The page says the list is not exhaustive. — [agents.md](https://agents.md/)
- **Spec nesting rule:** Agents use the nearest AGENTS.md in the directory tree. The closest file to the edited file wins on conflict, and explicit chat prompts override every file. — [agents.md](https://agents.md/)

**Per-tool behavior:**
- **Claude Code:**
  - Native AGENTS.md reading requires v2.1.277 or later. It is implemented as a built-in plugin, `cc-plugin-agents-md@builtin`.
  - Default mode is `claude-md-or-agents-md`: Claude reads AGENTS.md only when there is no `CLAUDE.md`, `.claude/CLAUDE.md`, or `CLAUDE.local.md` in the working directory or above it.
  - Other modes: `claude-md-and-agents-md` (both, de-duplicated), `claude-md`, and `managed-only`. The mode is set via `/config` or the user/managed `pluginConfigs` setting.
  - Subdirectory AGENTS.md files load on demand. `@path` imports work inside AGENTS.md.
  - Not read: `AGENTS.local.md`, `AGENTS.override.md`, and anything under `.agents/`.
  - Some sessions cannot read AGENTS.md directly. On versions before v2.1.281 this included Bedrock and telemetry-disabled sessions. The docs recommend a `CLAUDE.md` containing `@AGENTS.md` there.
  - Source: [Claude Code memory docs](https://code.claude.com/docs/en/memory)
- **OpenAI Codex CLI:**
  - Global scope: checks `~/.codex` (or `$CODEX_HOME`) for `AGENTS.override.md`, then `AGENTS.md`, and uses the first non-empty one.
  - Project scope: walks from the git root down to the current directory and takes at most one file per directory, in the order `AGENTS.override.md` > `AGENTS.md` > `project_doc_fallback_filenames`.
  - Files are concatenated root to leaf, so the closest file comes last and wins.
  - Total size is capped by `project_doc_max_bytes`, default 32 KiB. Files beyond the cap are dropped.
  - Source: [Codex AGENTS.md docs](https://learn.chatgpt.com/docs/agent-configuration/agents-md)
- **Gemini CLI:**
  - Default context file is `GEMINI.md`. It does NOT read AGENTS.md by default; you opt in with `context.fileName` in settings.json, e.g. `["AGENTS.md","GEMINI.md"]`.
  - Load hierarchy: global `~/.gemini/GEMINI.md`, then workspace and parent directories, then just-in-time scans of directories that tools touch. All found files are concatenated.
  - Supports `@file.md` imports.
  - The docs page says Gemini CLI "was replaced by Antigravity CLI on June 18th, 2026" for unpaid and Google One users.
  - Source: [Gemini CLI GEMINI.md docs](https://geminicli.com/docs/cli/gemini-md/)
- **GitHub Copilot:**
  - AGENTS.md can sit anywhere in the repo. "The nearest AGENTS.md file in the directory tree will take precedence."
  - `CLAUDE.md` or `GEMINI.md` at the repo root are accepted as single-file alternatives.
  - Also reads `.github/copilot-instructions.md` (repo-wide) and `.github/instructions/*.instructions.md` (frontmatter `applyTo` glob, optional `excludeAgent`).
  - Code review and the cloud agent read agent instructions. The page gives no full IDE matrix.
  - Source: [GitHub Docs](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions)
- **Cursor:** AGENTS.md is supported in the root and in subdirectories. Nested files combine with parent files, and the more specific file wins. — [Cursor rules docs](https://cursor.com/docs/context/rules)
- **OpenCode:**
  - Project: `AGENTS.md` in the current or a parent directory, falling back to `CLAUDE.md`.
  - Global: `~/.config/opencode/AGENTS.md`, falling back to `~/.claude/CLAUDE.md`.
  - Extra files go in `instructions` in `opencode.json` (globs and URLs allowed).
  - The first match in each category wins.
  - Source: [OpenCode rules docs](https://opencode.ai/docs/rules/)
- **Amp:**
  - Reads `AGENTS.md` from the current directory and parents up to `$HOME`. Subtree files load when Amp reads a file in that subtree.
  - Falls back to `AGENT.md` or `CLAUDE.md` per directory.
  - Global files: `~/.config/amp/AGENTS.md` and `~/.config/AGENTS.md`, plus system-wide paths.
  - `@path` imports accept globs. Imported files can carry `globs:` frontmatter so they load only when a matching file is read.
  - Source: [Amp AGENTS.md docs](https://ampcode.com/docs/customize/agents-md)
- **Zed:**
  - Uses only the FIRST file found at the project root, with no merging. Order: `.rules`, `.cursorrules`, `.windsurfrules`, `.clinerules`, `.github/copilot-instructions.md`, `AGENT.md`, `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`.
  - A March 2026 commit fixed a bug where CLAUDE.md was picked over AGENTS.md, contrary to the docs.
  - Nested files are not used (secondary source).
  - Sources: [Zed commit 12bdc20](https://git.secluded.site/zed/commit/12bdc208e5298a6aaa1803b07bf6f63599cb6add); [Zed docs mirror](https://git.secluded.site/zed/blob/ep-0327/docs/src/ai/rules.md?source=1); [PromptShelf (secondary)](https://thepromptshelf.dev/blog/zed-editor-agents-md-ai-rules-setup-guide-2026). zed.dev/docs/ai/rules returned 404 at fetch time.
- **Windsurf:**
  - AGENTS.md goes through the Rules engine. A root file is always on; a subdirectory file auto-globs to its directory.
  - Workspace rules: `.windsurf/rules/*.md`, 12,000 characters per file. Global rules: `~/.codeium/windsurf/memories/global_rules.md`, 6,000 characters.
  - Source: [Windsurf memories docs](https://docs.windsurf.com/windsurf/cascade/memories.md)
  - UNVERIFIED (secondary): the docs reportedly now redirect to docs.devin.ai, and the product is rebranded "Devin Desktop" with `.devin/rules/` preferred and `.windsurf/rules/` as a fallback. — [MemoryLake blog](https://www.memorylake.ai/blogs/windsurf-forgets-project-rules)
- **Cline:**
  - Reads `AGENTS.md` (project root) and global `~/.agents/AGENTS.md`.
  - Also reads `.cursorrules` and `.windsurfrules` automatically.
  - Source: [Cline rules docs](https://docs.cline.bot/features/cline-rules)
- **Aider:** Sources conflict.
  - Mintlify's agents.md compatibility page says Aider needs a `read:` entry in `.aider.conf.yml`. — [Mintlify](https://www.mintlify.com/agentsmd/agents.md/compatibility/aider)
  - Other secondary guides claim native detection since 2025. — [PromptShelf](https://thepromptshelf.dev/blog/aider-agents-md-vs-conventions-md-2026)
  - agents.md lists Aider as supporting it. — [agents.md](https://agents.md/)

### Inferences
- A root `AGENTS.md` reaches about 11 of the 12 target tools with zero config. Gemini CLI needs a one-line setting. Aider should get `read: AGENTS.md` in `.aider.conf.yml` to be safe.
- **Nesting semantics differ.** Codex, Gemini, Cursor, and Claude concatenate along the path (later wins). Copilot is nearest-wins. Zed reads the root file only. Canonical content should therefore not depend on nested-file merge behavior; put critical rules in the root file.
- **Size:** keep the combined root-to-leaf AGENTS.md chain under 32 KiB (Codex default) and under about 200 lines per file (Claude Code guidance).
- **Override and local-file names are not portable.** Claude Code ignores `AGENTS.override.md` and `AGENTS.local.md`.
- **Pitfall:** if a bootstrapper writes a non-trivial `CLAUDE.md`, Claude Code's default mode stops reading AGENTS.md. The CLAUDE.md must import it (`@AGENTS.md`). The same applies if a user adds `CLAUDE.local.md`.

### Gaps
- The official Zed docs page was unreachable (404). The order comes from a docs mirror and a commit.
- Aider's native AGENTS.md behavior could not be confirmed from aider.chat docs.
- Roo Code status is unclear. rulesync marks the `roo` target end-of-life and points to "zoocode" ([rulesync](https://github.com/dyoshikawa/rulesync)). Whether Roo Code was discontinued, forked, or renamed is unverified.
- No first-party Windsurf/Devin confirmation of the rebrand was found.

## Vendor-specific files: can they symlink or import AGENTS.md?

### Takeaway
Most vendor files are now optional fallbacks. Where one is still needed (Gemini without config, Cursor `.mdc` scoped rules, Copilot path-specific rules, Aider), it should be a thin pointer generated from the canonical source. Import syntax (`@file`) works in Claude Code, Gemini CLI, and Amp. Elsewhere, use a symlink, a config reference, or generated copies.

### Cited Findings
- **CLAUDE.md (Claude Code):**
  - Supports `@path` imports: relative to the importing file, recursive to 4 hops, skipped inside code spans. Escape spaces in paths with backslashes.
  - Official pattern: a `CLAUDE.md` containing `@AGENTS.md` plus Claude-specific notes below it.
  - A symlink (`ln -s AGENTS.md CLAUDE.md`) also works, with two caveats. The Edit/Write tools refuse to write through symlinks. On Windows, git checks a committed symlink out as a plain text file unless `core.symlinks` is enabled, so the docs recommend the import there.
  - Path-scoped rules live in `.claude/rules/`. Symlinked rules that point outside the repo need an external-import approval.
  - `/import` copies another agent's config (AGENTS.md, MCP servers, commands, subagents, skills) into Claude Code once (v2.1.213+).
  - Source: [Claude Code memory docs](https://code.claude.com/docs/en/memory)
- **GEMINI.md (Gemini CLI):** Supports `@file.md` imports. Alternatively set `context.fileName` to include `AGENTS.md`. — [Gemini CLI docs](https://geminicli.com/docs/cli/gemini-md/)
- **Cursor rules:**
  - Project rules live in `.cursor/rules` and must be `.mdc`. The docs say plain `.md` there is "ignored by the rules system".
  - Frontmatter fields: `description`, `globs`, `alwaysApply`. They select among four modes: Always, Auto-attached by glob, Agent-requested by description, Manual @-mention.
  - Precedence: Team Rules > Project Rules > User Rules, all merged.
  - The page does not cover legacy `.cursorrules` or CLAUDE.md.
  - Source: [Cursor rules docs](https://cursor.com/docs/context/rules)
- **GitHub Copilot:**
  - `.github/copilot-instructions.md` is plain Markdown with no import syntax documented.
  - `.github/instructions/*.instructions.md` uses `applyTo` globs.
  - Copilot also accepts AGENTS.md, CLAUDE.md, or GEMINI.md, so a separate copilot-instructions file is optional for agent surfaces.
  - Source: [GitHub Docs](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions/add-repository-instructions)
- **Windsurf:**
  - `.windsurf/rules/*.md` with per-file activation modes, 12k characters per file. AGENTS.md is native.
  - Legacy `.windsurfrules` still appears in third-party precedence lists.
  - Sources: [Windsurf docs](https://docs.windsurf.com/windsurf/cascade/memories.md); [PromptShelf (secondary)](https://thepromptshelf.dev/blog/agents-md-windsurf-setup-guide-2026/)
- **Cline:**
  - Workspace rules are directories, `.clinerules/` and/or `.cline/rules/`, containing Markdown files.
  - Conditional rules use a YAML `paths:` glob array. Rules without frontmatter are always on.
  - Global rules: `~/Documents/Cline/Rules`, `~/.cline/rules`, plus `~/.agents/AGENTS.md`.
  - AGENTS.md, `.cursorrules`, and `.windsurfrules` are auto-read.
  - Source: [Cline rules docs](https://docs.cline.bot/features/cline-rules)
- **Aider CONVENTIONS.md:**
  - Not auto-loaded. Load it with `/read`, `--read`, or `read:` in `.aider.conf.yml`, which makes it read-only and prompt-cacheable.
  - The same mechanism can point at AGENTS.md.
  - Sources: [PromptLayer glossary (secondary)](https://www.promptlayer.com/glossary/aider-conventions-md); [Mintlify](https://www.mintlify.com/agentsmd/agents.md/compatibility/aider)
- **Amp:** `@path` imports with globs, plus per-file `globs:` frontmatter for conditional inclusion. — [Amp docs](https://ampcode.com/docs/customize/agents-md)
- **Codex:** No import syntax is documented. `project_doc_fallback_filenames` lets Codex read other filenames. — [Codex docs](https://learn.chatgpt.com/docs/agent-configuration/agents-md)

### Inferences
**Practical bridge table:**

| Tool | Bridge |
|---|---|
| Claude Code | Native, or `CLAUDE.md` = `@AGENTS.md` (preferred over a symlink for Windows safety) |
| Gemini CLI | `.gemini/settings.json` `context.fileName: ["AGENTS.md","GEMINI.md"]`, or `GEMINI.md` = `@AGENTS.md` |
| Cursor | Native AGENTS.md. Generate `.mdc` files only for glob-scoped rules |
| Copilot | Native AGENTS.md. Optionally generate `.github/instructions/*.instructions.md` for `applyTo` scoping |
| Windsurf | Native |
| Cline | Native |
| Zed | Native, provided no `.rules`, `.cursorrules`, `.windsurfrules`, `.clinerules`, or `.github/copilot-instructions.md` exists at the root (they would shadow AGENTS.md) |
| Aider | `.aider.conf.yml` with `read: [AGENTS.md]` |
| Codex, OpenCode, Amp | Native |

- **Zed shadowing hazard:** generating legacy `.cursorrules`, `.windsurfrules`, or `.github/copilot-instructions.md` makes Zed use that file instead of AGENTS.md. A bootstrapper should avoid emitting legacy root files, or make their content identical.
- **Path-scoped rules have no common format.** Glob frontmatter keys differ by tool:

  | Tool | Glob key |
  |---|---|
  | Cursor | `globs` |
  | Copilot | `applyTo` |
  | Cline | `paths` |
  | Amp | `globs` |
  | Claude `.claude/rules` | `paths` (per docs anchor "path-specific rules"; field name seen in docs text) |

  A generator is required to emit these.

### Gaps
- Whether Cursor still honors legacy `.cursorrules`, and whether it reads CLAUDE.md, was not confirmed from the fetched page.
- Whether Copilot IDE (VS Code) chat reads AGENTS.md was not confirmed from the fetched page. VS Code is listed on agents.md.

## MCP (Model Context Protocol): governance, clients, relevance for memory

### Takeaway
MCP is now an AAIF/Linux Foundation project and is the settled cross-vendor layer for tools, which makes it the right transport for a portable memory server. The 2026-07-28 spec revision is the largest since launch.

### Cited Findings
- Anthropic donated MCP to AAIF at AAIF's formation on 2025-12-09. — [Linux Foundation](https://linuxfoundation.org/press/linux-foundation-announces-the-formation-of-the-agentic-ai-foundation)
- Anthropic claims more than 10,000 public MCP servers. Support is claimed in ChatGPT, Cursor, Gemini, Microsoft Copilot, and VS Code (vendor claims). — [YourStory](https://yourstory.com/ai-story/agentic-ai-foundation-openai-anthropic-linux)
- **2026-07-28 spec revision:** AAIF calls it the largest revision since launch. It moves to a stateless core, adds MCP Apps and Tasks, and aligns authorization with OAuth. The 2026 roadmap (March 2026) prioritized transport evolution, agent communication, governance maturation, and enterprise readiness. — [AAIF blog](https://aaif.io/blog?hsLang=en); [mcp.directory (secondary)](https://mcp.directory/blog/mcp-foundation-linux-foundation-aaif-2026-explained); [MCP Playground roadmap summary (secondary)](https://mcpplaygroundonline.com/blog/mcp-2026-roadmap-whats-changing-for-developers)
- **Conflict:** one July 2026 blog describes stewardship transfer as happening in July 2026 rather than December 2025. Treat it as low reliability. — [agentready.it.com](https://www.agentready.it.com/blog/mcp-agentic-ai-foundation)
- **Tools on the target list with MCP support** (each tool's docs or a sync-tool target list):
  - Ruler propagates `[mcp_servers.*]` from `ruler.toml` into each agent's config. — [Ruler](https://github.com/intellectronica/ruler)
  - rulesync manages MCP config across dozens of tools. — [rulesync](https://github.com/dyoshikawa/rulesync)
  - Amp has a dedicated MCP doc page. — [Amp manual](https://ampcode.com/manual)
  - nanobot and Deep Code list MCP and Skills together. — [agentskills.io clients](https://agentskills.io/)
- A May 2026 buyer's guide describes tool calling via MCP as settled, while memory interchange is "still moving". — [cited via search summary; aiagentrank.io](https://aiagentrank.io/blog/agentic-ai-foundation-explained)

### Inferences
- **MCP config files are vendor-specific** (e.g. `.mcp.json` for Claude Code, `.cursor/mcp.json`, `.vscode/mcp.json`, `~/.codex/config.toml`, `opencode.json`, Gemini `settings.json`). A bootstrapper should keep one canonical server list and generate these, as Ruler and rulesync do.
- **For memory:** shipping the memory store as an MCP server (exposing tools like remember, recall, and search) is the only cross-vendor way to give agents read/write memory at runtime. File-based memory (Markdown in the repo) is the zero-dependency fallback that every tool can read.
- **Local models:** Ollama itself is a model server, not an agent. MCP and AGENTS.md support comes from the harness (OpenCode, Aider, goose, Cline, Zed, Continue) pointed at an Ollama endpoint. Compatibility for "local-model setups" therefore reduces to the harness chosen.

### Gaps
- No authoritative per-client MCP config-file path table was fetched. The paths above come from background knowledge and should be verified against each tool's docs.
- Whether the 2026-07-28 spec was finalized versus left at release candidate was not confirmed against modelcontextprotocol.io.
- Aider's MCP status is not confirmed. Historically it lacked native MCP.

## Agent Skills (SKILL.md) open standard

### Takeaway
Agent Skills is an Anthropic-originated open format (agentskills.io; repo github.com/agentskills/agentskills). In 2026 it is supported by more than 40 products, including Codex, Gemini CLI, Cursor, Copilot/VS Code, OpenCode, Amp, Roo Code, goose, Junie, and Kiro. `.agents/skills/` is emerging as the vendor-neutral discovery directory, but support for it is uneven.

### Cited Findings
- **Format:**
  - A skill is a folder containing `SKILL.md` (YAML frontmatter plus Markdown), with optional `scripts/`, `references/`, and `assets/`.
  - Required fields: `name` (1-64 characters, lowercase a-z, 0-9, and hyphens, no leading, trailing, or double hyphens, must match the directory name) and `description` (1-1024 characters).
  - Optional fields: `license`, `compatibility` (≤500 characters), `metadata` (string map), `allowed-tools` (experimental).
  - Progressive disclosure: about 100 tokens of metadata at startup, a body under 5,000 tokens or 500 lines, resources on demand.
  - Validator: `skills-ref validate`.
  - Source: [agentskills.io specification](https://agentskills.io/specification)
- **Governance:** "originally developed by Anthropic, released as an open standard", open to contributions via GitHub and Discord. AAIF stewardship was not stated. — [agentskills.io](https://agentskills.io/)
- **Client showcase** (with docs links): Junie, ZeroClaw, Gemini CLI, Autohand, OpenCode, OpenHands, Mux, Cursor, Amp, Letta, Firebender, goose, GitHub Copilot, VS Code, Claude Code, Claude, ChatGPT & Codex, Piebald, Factory, pi, Databricks, Agentman, TRAE, Spring AI, Roo Code, Mistral Vibe, Command Code, Ona, VT Code, Qodo, Laravel Boost, Emdash, Snowflake Cortex Code, Kiro, Workshop, Google AI Edge Gallery, nanobot, fast-agent, bub, Tabnine, Vita, Superconductor, Deep Code, Pulumi Neo, Hermes Agent, OpenClaw. — [agentskills.io](https://agentskills.io/)
- **Cursor discovery paths:**
  - Project: `.agents/skills/` and `.cursor/skills/`. User: `~/.agents/skills/` and `~/.cursor/skills/`.
  - Compatibility paths: `.claude/skills/`, `.codex/skills/`.
  - Nested folders in a repo are discovered.
  - A forum report says `.agents/skills` was not discovered on Cursor 2.5.26 (possible bug).
  - Sources: [Cursor skills docs](https://cursor.com/docs/context/skills); [Cursor forum](https://forum.cursor.com/t/agents-skills-discovery/151764)
- **Codex:** reportedly uses `.agents/skills/`, with `.codex/skills/` deprecated (secondary; first-party page 404'd). — [viblo.asia (secondary)](https://viblo.asia/p/ai-agent-config-rules-skills-XP4WEXDqL7G); [mcp.directory (secondary)](https://mcp.directory/blog/cross-agent-skills-cursor-codex-cline-antigravity-gemini-mastra-portability)
- **Claude Code:** does NOT read anything under `.agents/` for instruction files. Skills live in `.claude/skills/`. — [Claude Code memory docs](https://code.claude.com/docs/en/memory)
- **OpenCode:** reads `.claude/skills` for compatibility (disable with `OPENCODE_DISABLE_CLAUDE_CODE_SKILLS=1`). — [OpenCode rules docs](https://opencode.ai/docs/rules/)
- **Sync tools:** Ruler copies `.ruler/skills/*` to the native skills directories of about 20 agents. — [Ruler](https://github.com/intellectronica/ruler)

### Inferences
- **Directory strategy:** canonical skills in `.agents/skills/<name>/SKILL.md`, mirrored or symlinked to `.claude/skills/`. That covers Claude Code, OpenCode, and Cursor (via compat) without depending on `.agents` support.
- **Non-skills tools:** for Aider, Zed (has its own Skills feature, format unverified), and plain local-model setups, list skills by name and description in AGENTS.md so the agent can read the SKILL.md on demand.
- **Portability:** keep to the required spec fields. Vendor extras (e.g. Claude Code's extra frontmatter) should go in `metadata` or be tolerated as unknown keys.

### Gaps
- Official Codex, Gemini CLI, and Copilot skills discovery paths were not confirmed first-hand.
- Whether the Agent Skills spec has moved under AAIF governance is unknown.

## llms.txt and portable agent memory standards

### Takeaway
llms.txt is an informal, ungoverned convention for site-level docs. It is useful for documentation consumed by IDE agents but irrelevant to repo memory. No portable agent-memory standard has been adopted as of Oct 2026. Several 2026 proposals exist (W3C Community Group, MacPaw Portable Memory, an IETF draft, Engram). Vendor memory remains siloed (e.g. Claude Code auto memory in `~/.claude/projects/<project>/memory/MEMORY.md`).

### Cited Findings
- **llms.txt status:** no W3C, IETF, or schema.org recognition, no versioning, no governance body. — [Codersera (secondary)](https://codersera.com/blog/llms-txt-complete-guide-2026/amp/)
- **llms.txt adoption:** about 5.9-13% depending on sample (e.g. 8,598 of the top 100k domains as of 2026-09-13). — [Seomator](https://seomator.com/blog/llms-txt-adoption); [Digital Applied](https://www.digitalapplied.com/blog/llms-txt-in-practice-adoption-evidence-2026)
- **llms.txt consumption:** minimal. Digital Applied logged 0 verified AI-lab crawler hits out of 1,227 requests. Google says it does not use llms.txt. Some IDE agents do use it. — [Digital Applied](https://www.digitalapplied.com/blog/llms-txt-in-practice-adoption-evidence-2026); [Common Crawl](https://commoncrawl.org/blog/a-content-analysis-of-llms-txt-files-from-the-july-2026-crawl-archive)
- **Claude Code auto memory:**
  - Stored at `~/.claude/projects/<project>/memory/`, shared across worktrees. Relocatable via `autoMemoryDirectory`.
  - Layout: a `MEMORY.md` index with one line per memory, plus topic files.
  - The first 200 lines or 25 KB of MEMORY.md load every session.
  - Disable with `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1`.
  - Source: [Claude Code memory docs](https://code.claude.com/docs/en/memory)
- **Gemini CLI:** has a "Memory tool" for persistent memories, plus `/memory show` and `/memory reload`. — [Gemini CLI docs](https://geminicli.com/docs/cli/gemini-md/)
- **Windsurf Cascade memories:** auto-generated and separate from rules. Reportedly not persisted by the newer default agent, and Windsurf recommends putting durable knowledge in rules or AGENTS.md (secondary). — [MemoryLake](https://www.memorylake.ai/blogs/stop-windsurf-cascade-losing-context)
- **W3C AI Agent Memory Interoperability Community Group:** charter adopted 2026-06-19. — [W3C](https://www.w3.org/community/ai-agent-memory-interop)
- **MacPaw "Portable Memory":** July 2026 preprint proposing a `.mem` bundle. The authors admit it is not yet independently implemented or neutrally governed. — [MacPaw research](https://research.macpaw.com/publications/portable-memory)
- **Other proposals:**
  - IETF draft "ApertoMemory" (signed, encrypted memory objects). — [IETF draft mirror](https://ftp.kaist.ac.kr/ietf/draft-ferro-apertomemory-02.xml)
  - Engram spec (PLUR). Its own write-up says no adopted standard exists. — [DEV](https://dev.to/plur9/is-there-an-open-standard-for-ai-agent-memory-engrams-2ckj)
  - Universal Memory Protocol (community post). — [feedbagel](https://feedbagel.com/post/universal-memory-protocol-a-shared-format-for-agent-memory)
  - Google Cloud "Open Knowledge Format" (2026-06-12), UNVERIFIED (Medium only). — [Medium](https://medium.com/p/7e96a33898ce)

### Inferences
- A bootstrapper should not bet on any memory standard yet. The pragmatic portable memory is:
  - (a) plain Markdown files in the repo (e.g. a `memory/` or `docs/agent-memory/` directory with an index file modeled on Claude's MEMORY.md pattern: one line per entry, topic files), referenced from AGENTS.md;
  - (b) optionally, an MCP memory server for read/write at runtime.
- Write vendor auto-memory locations (Claude `~/.claude/projects/...`, Gemini, Windsurf) only via an import or sync step, if at all.
- llms.txt is worth generating only if the project publishes docs on the web.

### Gaps
- The Codex "Memories" feature exists in its docs navigation, but its storage location was not fetched.
- The state of the W3C CG deliverables is unknown.

## Rule-sync and translation tools; recommended single-source-of-truth pattern

### Takeaway
Two mature open-source generators dominate:
- **Ruler** (`@intellectronica/ruler`, `.ruler/` source, about 30 agents, MCP plus skills)
- **rulesync** (`dyoshikawa/rulesync`, `.rulesync/` source, v14.x, rules, MCP, commands, subagents, skills, hooks, permissions, import and convert)

The recommended 2026 pattern is AGENTS.md as the canonical, committed file, with thin bridges (imports or config) rather than full copies, and a generator only for the features that have no common format (scoped rules, MCP config, skills mirroring).

### Cited Findings
- **Ruler:**
  - Source: `.ruler/*.md` (recursive, concatenated with `<!-- Source -->` markers), `ruler.toml`, `.ruler/skills/`, `.ruler/agents/`, and a global `$XDG_CONFIG_HOME/ruler`.
  - Commands: `ruler apply` writes each agent's files, makes `.bak` backups, and manages a `.gitignore` block for generated files. Also `--dry-run` and `ruler revert`.
  - MCP via `[mcp_servers.<name>]` with merge or overwrite strategies, written only within the project.
  - Experimental nested mode. Labelled "Beta Research Preview"; MIT license.
  - Source: [github.com/intellectronica/ruler](https://github.com/intellectronica/ruler)
- **rulesync:**
  - Commands: `rulesync init`, `generate --targets "*" --features "*"`, `import --targets <tool>` (e.g. from CLAUDE.md, .cursorrules, copilot-instructions), and `convert --from X --to Y`.
  - Generated files remain usable without the tool.
  - Notes: the `ignore` feature is deprecated in favor of `permissions`; the `roo` and `continue` targets are end-of-life.
  - Sources: [github.com/dyoshikawa/rulesync](https://github.com/dyoshikawa/rulesync); [npm rulesync](https://www.npmjs.com/package/rulesync); [supported tools](https://rulesync.dyoshikawa.com/reference/supported-tools)
- **Name collision:** a different PyPI package, "rulesync"/agentsync, uses `.agentsync/rules.md` as its canonical file. — [PyPI rulesync](https://pypi.org/project/rulesync); [PanisHandsome/ai-rules-sync](https://github.com/PanisHandsome/ai-rules-sync)
- **Built-in importers:** Claude Code's `/import` copies another agent's config once. — [Claude Code docs](https://code.claude.com/docs/en/memory)
- **Bridge pattern:** the Claude Code docs explicitly recommend AGENTS.md as "the one file every tool shares", with a `CLAUDE.md` `@AGENTS.md` import. — [Claude Code docs](https://code.claude.com/docs/en/memory)
- **Comparison guide:** a third-party guide compares Ruler and rulesync and recommends AGENTS.md first, adding a generator only when you also sync MCP, skills, or more agents. — [developertoolkit.ai](https://developertoolkit.ai/en/ecosystem/frameworks/rules-sync/)

### Inferences

**Recommended canonical layout for the bootstrapper** (synthesis, not a vendor-endorsed spec):

```
AGENTS.md                      # canonical instructions; <~200 lines; root-level critical rules
<subdir>/AGENTS.md             # optional; never rely on merge semantics for critical rules
.agents/skills/<name>/SKILL.md # canonical skills (agentskills.io spec, required fields only)
.agents/memory/MEMORY.md       # (proposal) portable memory index + topic files, referenced from AGENTS.md
.agents/mcp.json (or toml)     # canonical MCP server list -> generated per tool
# thin bridges (generated, committed):
CLAUDE.md                      # "@AGENTS.md" + Claude-only notes (import, not symlink, for Windows)
GEMINI.md                      # "@AGENTS.md"  (or .gemini/settings.json context.fileName)
.aider.conf.yml                # read: [AGENTS.md]
.claude/skills -> mirror of .agents/skills (copy or symlink)
.cursor/rules/*.mdc, .github/instructions/*.instructions.md, .clinerules/  # only for glob-scoped rules
# avoid: .cursorrules, .windsurfrules, root .rules, .github/copilot-instructions.md (Zed shadowing)
```

- **Prefer imports over symlinks.** Claude's editing tools refuse to write through symlinks, and Windows clones break committed symlinks. Prefer generated copies over symlinks for directories when Windows matters.
- **Ruler vs rulesync:** reuse one of them as the generation engine, or emit compatible output. Ruler is closer to "AGENTS.md canonical"; rulesync is richer on features and import.
- **Things to verify before shipping:** Gemini/Antigravity CLI rename impact, the Windsurf-to-Devin path change (`.devin/rules`), Aider's native AGENTS.md behavior, and `.agents/skills` discovery in Claude Code, which per its docs does not read `.agents/`.

### Gaps
- Other sync tools (e.g. "ai-rules-sync", "agentsync", vendor CLIs like `npx skills add`) were not surveyed in depth.
- No first-hand check of each tool's MCP config path. The MCP generator targets should be taken from Ruler's or rulesync's supported-tools tables.
