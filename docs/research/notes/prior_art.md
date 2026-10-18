# Prior Art: Folder/Vault/Project Bootstrappers and Memory for AI Agents (as of Oct 2026)

Star counts were pulled from the public GitHub REST API (api.github.com/repos/...) on 2026-10-18. They change daily, so treat them as a snapshot. "Last push" means the repo's `pushed_at` date on that day.

## Q1: Obsidian + Claude Code / Codex / Gemini vault templates and "second brain" starter kits

### Takeaway
This area is crowded and moving fast. The biggest Obsidian-native projects are claude-obsidian (~15.6k stars), the Claudian plugin (~15.7k), and obsidian-second-brain (~4.8k). Most of the newer projects implement Karpathy's April 2026 "LLM Wiki" pattern. Almost all started out Claude Code-first. They have since added multi-agent support through the shared Agent Skills format: one `.agents/skills/` or `.skills/` folder, symlinked into each agent's folder, with AGENTS.md as the fallback. So hard lock-in has weakened, but Claude Code is still the main target everywhere.

### Cited Findings

**Snapshot of key repos (stars as of 2026-10-18):**

| Repo | Stars | Last push | Approach | Agent coverage |
|---|---|---|---|---|
| [AgriciDaniel/claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian) | 15,606 | 2026-10-17 | LLM Wiki; Python CLI `init`/`adopt` with plan-hash approval; claim ledgers | Claude Code plugin first; `setup-multi-agent.sh` for Codex, OpenCode, Gemini, ZCode |
| [YishenTu/claudian](https://github.com/YishenTu/claudian) | 15,714 | 2026-10-18 | Obsidian plugin that embeds an agent in the vault (not a template) | Claude Code, Codex (plus others per listings) |
| [eugeniughelbur/obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain) | 4,802 | 2026-10-17 | Vault as memory; LLM Wiki extension; "Crew" shared logbook | Claude Code, Codex, Gemini CLI, OpenCode, Antigravity, Hermes, Pi, Grok; Cursor through Crew |
| [Ar9av/obsidian-wiki](https://github.com/Ar9av/obsidian-wiki) | 3,584 | 2026-10-13 | LLM Wiki framework; `pip install obsidian-wiki` and `obsidian-wiki setup --vault` | Claude Code, Cursor, Codex, Windsurf, Gemini CLI, "a dozen more" |
| [heyitsnoah/claudesidian](https://github.com/heyitsnoah/claudesidian) | 2,599 | 2026-04-11 | PARA vault template plus `/init-bootstrap` wizard | Claude Code first; skills in `.agents/skills/` symlinked for OpenCode, Codex, Cursor, Pi |
| [Astro-Han/karpathy-llm-wiki](https://github.com/Astro-Han/karpathy-llm-wiki) | 2,512 | 2026-07-23 | LLM Wiki as an installable Agent Skill | Claude Code, Cursor, Codex |
| [ballred/obsidian-claude-pkm](https://github.com/ballred/obsidian-claude-pkm) | 1,882 | 2026-02-18 | "Complete starter kit" for Obsidian + Claude Code PKM | Claude Code (by its description) |
| [huytieu/COG-second-brain](https://github.com/huytieu/COG-second-brain) | 1,294 | 2026-10-02 | Cognition + Obsidian + Git; numbered folders; 33–35 skills, verifier agents | Claude Code (full), Antigravity, Cursor, Kiro, Gemini CLI; AGENTS.md fallback for Codex |
| [kepano/obsidian-skills](https://github.com/kepano/obsidian-skills) | 49,702 | 2026-09-15 | Agent skills that teach agents Obsidian formats (Markdown, Bases, JSON Canvas) and the Obsidian CLI. By Obsidian's CEO. | Agent Skills format (agent-agnostic) |
| [MarkusPfundstein/mcp-obsidian](https://github.com/MarkusPfundstein/mcp-obsidian) | 4,480 | 2026-08-31 | MCP server using the Obsidian Local REST API plugin | Any MCP client |
| [ScrapingArt/Karpathy-LLM-Wiki-Stack](https://github.com/ScrapingArt/Karpathy-LLM-Wiki-Stack) | 95 | 2026-10-14 | Reference build for Obsidian + Claude Code | Claude Code |
| [ignromanov/llm-obsidian-wiki](https://github.com/ignromanov/llm-obsidian-wiki) | 17 | 2026-05-30 | Claude Code plugin; typed pages (concept, decision, comparison, open question) | Claude Code |

**Karpathy LLM Wiki pattern (the dominant design in 2026):**
- Gist created April 4, 2026, with 5,000+ stars and 5,000+ forks. It has three layers: (1) immutable raw sources, (2) an LLM-maintained markdown wiki, and (3) a schema file (CLAUDE.md or AGENTS.md). It has three operations: ingest, query, lint. — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- The special files are `index.md` (a catalog with one-line summaries that the LLM reads first, which avoids embedding RAG at moderate scale) and `log.md` (append-only, with parseable prefixes). — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- Tools are optional: Obsidian as the "IDE," Web Clipper, qmd (local markdown search CLI/MCP), Marp, Dataview, git. The gist names Codex, Claude Code, and OpenCode/Pi, so it is agent-agnostic by design. — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- Karpathy's framing: "Obsidian is the IDE, the LLM is the programmer, the wiki is the codebase." — [DenserAI summary](https://denser.ai/blog/llm-wiki-karpathy-knowledge-base/)
- A community ranking of implementations names claude-obsidian the top pick for Obsidian integration. — [srisatyalokesh](https://srisatyalokesh.is-a.dev/learn-ai/best-llm-wiki-community-implementations/)
- green-dalii/obsidian-llm-wiki is an Apache-2.0 Obsidian community plugin. It puts ingest, query, and lint behind buttons and runs locally. (The repo URL now returns a 301 redirect, so it may have been renamed or moved.) — [mythos.one](https://mythos.one/me/brianswichkow/95e1d3)

**claude-obsidian details:**
- The product checkout is kept separate from the user vault. The vault layout is `inbox/`, `.raw/` (immutable, content-addressed), `wiki/`, `.claude-obsidian.json`, and `.vault-meta/`. `init` prints a plan with an `approved_plan_sha256`, and you re-run it with `--apply`. `adopt` handles existing vaults without destroying anything. — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian)
- `wiki-mode` routes new notes using Generic, LYT, PARA, or Zettelkasten conventions "without reorganizing existing notes." — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian)
- Source and claim ledgers track authority, freshness, contradiction, and confidence. High-risk claims need two independent sources. Writes are transactional with `transaction recover`. Retrieval is BM25, and embedding reranking needs explicit consent. Native Windows writes need WSL. — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian)

**claudesidian details:**
- Numbered PARA folders (`00_Inbox`, `01_Projects`, `02_Areas`, `03_Resources`, `04_Archive`) plus metadata, templates, and scripts. The `/init-bootstrap` wizard installs dependencies, detaches from the template's git origin, imports an existing vault, interviews the user, writes a personalized CLAUDE.md, and runs git init. `/upgrade` merges template updates and keeps customizations. It has no dedicated memory system: "the vault itself" is the memory. — [claudesidian](https://github.com/heyitsnoah/claudesidian)

**obsidian-second-brain details:**
- "OKM" rule: every stored fact must be timeless, dated, or a pointer to a live source, which is an explicit anti-staleness rule. The vault layout is `raw/`, `wiki/` (entities, concepts, projects, daily notes, decisions), `boards/`, `templates/`, `index.md`, and `log.md`. — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain)
- "Crew" is a plain-markdown shared logbook: each agent writes one line per job, and other agents read it through hooks. A lesson becomes a rule only after a human runs `keep`. — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain)
- Cross-project: you set `OBSIDIAN_VAULT_PATH` in each project's `.claude/settings.json`. The docs say there is "no isolation within one vault yet." — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain)

**Ar9av/obsidian-wiki details:**
- Setup registers SessionStart (inject memory) and SessionEnd (capture knowledge) hooks. One `.skills/` folder is symlinked into each agent. `index.md`, `log.md`, and `hot.md` go through a locking atomic writer so parallel agents don't collide. A manifest records ingested sources so re-ingests only process changes. `/wiki-history-ingest` mines past sessions from Claude, Codex, Copilot, Hermes, OpenClaw, and Pi. `user_id` scoping is supported, and the vault can be served over HTTP/MCP via Docker. — [obsidian-wiki](https://github.com/Ar9av/obsidian-wiki)

**COG details:**
- Folders: `00-inbox`, `01-daily`, `02-personal`, `03-professional`, `04-projects`, `05-knowledge` (including a people CRM), `06-templates`. Install with `npx skills add huytieu/COG-second-brain` or clone and say "run onboarding" (~2 minutes). A `memory-hygiene` skill re-checks claims and stamps them with `last_verified` and a confidence level. The design is adapted from Garry Tan's "gbrain" ("Compiled Truth + Timeline"). Read-only verifier agents check worker output ("the worker never grades its own homework"). — [COG-second-brain](https://github.com/huytieu/COG-second-brain)

**Non-Obsidian-template memory tools often used alongside vaults:**
- [basicmachines-co/basic-memory](https://github.com/basicmachines-co/basic-memory) (4,261 stars): an MCP server over markdown. Each note is an entity with typed observations (`[category]`) and wikilink relations, which form a knowledge graph. It has a SQLite index, works with Obsidian as-is, and manages multiple projects (`project add/list`, per-project local or cloud). AGPL-3.0, with a paid cloud tier at $15/mo. — [basic-memory](https://github.com/basicmachines-co/basic-memory)
- [thedotmack/claude-mem](https://github.com/thedotmack/claude-mem) (101,471 stars): "Captures everything your agent does during sessions, compresses it with AI, and inject[s]" it into later sessions. It is now pitched as "for Every Agent." — [GitHub API description](https://github.com/thedotmack/claude-mem)
- Infrastructure-grade memory: [mem0](https://github.com/mem0ai/mem0) (67,372), [cognee](https://github.com/topoteretes/cognee) (32,453), [graphiti](https://github.com/getzep/graphiti) (31,935), [letta](https://github.com/letta-ai/letta) (25,175), [mcp-memory-service](https://github.com/doobidoo/mcp-memory-service) (2,021). These are databases or services, not folder bootstrappers.

### Inferences
- The field has converged on a shared stack: plain markdown + Obsidian as the viewer + an AGENTS.md/CLAUDE.md schema + Agent Skills symlinked per agent + hooks for session start and end + an index/log file pair. A new tool would not stand out by doing any one of these.
- The most mature design ideas to borrow are claude-obsidian's plan-hash approval and transactional writes, the "dated or pointer" fact rule from obsidian-second-brain and COG, Ar9av's atomic writer and source manifest, and claudesidian's interview-driven bootstrap with an `/upgrade` that preserves customizations.
- Every Obsidian project is either a vault template (one vault, personal knowledge) or a skill pack. None of them handles both the personal vault and the many code repos a developer works in under one model. Each treats "project" as either a PARA folder or a separate vault path.

### Gaps
- No primary star history was gathered. The "917 stars in 30 days" figure for obsidian-second-brain comes from [gittrend](https://gittrend.io/repo/eugeniughelbur/obsidian-second-brain) and was not verified.
- The ballred/obsidian-claude-pkm and Claudian READMEs were not fetched in detail.
- No independent quality review of these repos was found. Most write-ups are promotional.

## Q2: Memory-bank and spec-driven patterns (Cline, Cursor, Roo, context engineering, BMAD, spec-kit, Agent OS, Task Master)

### Takeaway
Two groups bootstrap project structure for agents. (a) "Memory bank" conventions are a fixed set of markdown files, generated when you ask the agent to initialize. (b) Spec-driven development (SDD) toolkits install folders, commands, and templates per project. The SDD group is far larger by stars (spec-kit ~142k, OpenSpec ~72k, GSD ~64k, BMAD ~54k). All of them are scoped to one code repo. None of them handles personal knowledge or memory across projects.

### Cited Findings
- **Cline Memory Bank:** a `memory-bank/` folder with `projectbrief.md`, `productContext.md`, `activeContext.md`, `systemPatterns.md`, `techContext.md`, and `progress.md`. You bootstrap it by pasting the custom instructions into `.clinerules/memory-bank.md` and saying "initialize memory bank." Updates happen on new patterns, after significant changes, on "update memory bank" (which re-reads every file), or when context fills up. The docs say it "works with any AI that can read docs." — [Cline docs](https://docs.cline.bot/prompting/cline-memory-bank)
- [vanzan01/cursor-memory-bank](https://github.com/vanzan01/cursor-memory-bank) (3,060 stars, last push 2026-05): a "modular, documentation-driven framework using Cursor custom modes (VAN, PLAN, CREATIVE, IMPLEMENT)." It is locked to Cursor. — [GitHub](https://github.com/vanzan01/cursor-memory-bank)
- [GreatScottyMac/roo-code-memory-bank](https://github.com/GreatScottyMac/roo-code-memory-bank) (1,667 stars, last push 2025-05-15, effectively dormant). — [GitHub](https://github.com/GreatScottyMac/roo-code-memory-bank)
- [centminmod/my-claude-code-setup](https://github.com/centminmod/my-claude-code-setup) (2,671 stars): "Shared starter template configuration and CLAUDE.md memory bank system for Claude Code." — [GitHub](https://github.com/centminmod/my-claude-code-setup)
- [coleam00/context-engineering-intro](https://github.com/coleam00/context-engineering-intro) (13,880 stars): "Context engineering is the new vibe coding… Claude Code is the best for this so that's what [it uses]." This is explicit Claude orientation. — [GitHub](https://github.com/coleam00/context-engineering-intro)
- [github/spec-kit](https://github.com/github/spec-kit) (141,709 stars): "Toolkit to help you get started with SDD or any other process." — [GitHub](https://github.com/github/spec-kit)
- [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) (72,283 stars): "Spec-driven development (SDD) for AI coding assistants." — [GitHub](https://github.com/Fission-AI/OpenSpec)
- [gsd-build/get-shit-done](https://github.com/gsd-build/get-shit-done) (64,277 stars): "meta-prompting, context engineering and spec-driven development system for Claude Code." — [GitHub](https://github.com/gsd-build/get-shit-done)
- [bmad-code-org/BMAD-METHOD](https://github.com/bmad-code-org/BMAD-METHOD) (54,271 stars, very active): "Breakthrough Method for Agile AI Driven Development." — [GitHub](https://github.com/bmad-code-org/BMAD-METHOD)
- [buildermethods/agent-os](https://github.com/buildermethods/agent-os) (5,515 stars): "a system for injecting your codebase standards and writing better specs." It "works alongside Claude Code, Cursor, Antigravity, and other AI tools." The repo contains `profiles/default/global`, `commands/agent-os`, and `config.yml`. Install docs are at buildermethods.com/agent-os. — [Agent OS](https://github.com/buildermethods/agent-os)
- [eyaltoledano/claude-task-master](https://github.com/eyaltoledano/claude-task-master) (28,240 stars, last push 2026-04-28): an "AI-powered task-management system you can drop into Cursor, Lovable, Windsurf, Roo, and others." — [GitHub](https://github.com/eyaltoledano/claude-task-master)
- [obra/superpowers](https://github.com/obra/superpowers) (301,328 stars): "An agentic skills framework & software development methodology." It is the biggest repo in the skills/methodology space. — [GitHub](https://github.com/obra/superpowers)
- [oraios/serena](https://github.com/oraios/serena) (30,352 stars): an MCP coding toolkit, "the IDE for your agent." — [GitHub](https://github.com/oraios/serena)

### Inferences
- Memory banks were a 2025 pattern. The Roo bank has gone dormant and Cursor's slowed after May 2026. Agent-native features (auto-compaction, native memory files, hooks) and SDD toolkits have absorbed the idea.
- SDD toolkits are the closest thing to `create-*` scaffolders for agent workspaces. They write folders and templates into a repo and register slash commands for several agents. But their structure is about specs and tasks, not knowledge or memory, and it lives and dies with one repo.
- Fixed memory-bank files (`activeContext.md`, `progress.md`) are a known staleness risk: they are only as fresh as the last explicit "update memory bank."

### Gaps
- Install details were not re-checked this session: spec-kit's `specify init` agent list, BMAD's installer and the agents it supports, and Agent OS's project folder layout (the README points to external docs).
- No figures were found on how many users actually keep memory banks updated.

## Q3: Organizing frameworks (PARA, Johnny Decimal, Zettelkasten, ACE) adapted for agents

### Takeaway
Numbered PARA is the default (claudesidian, many CLAUDE.md templates). Zettelkasten and LYT appear as alternative routing modes. Johnny Decimal shows up only as numeric folder prefixes. No search turned up a real JD-for-agents template. A newer, agent-native layout from the LLM Wiki pattern (`raw/` → `wiki/` + `index.md`/`log.md`) is replacing PARA as the default in 2026 projects. Some practitioners argue for flat vaults that use only links and tags.

### Cited Findings
- claudesidian uses numbered PARA: `00_Inbox` through `04_Archive`. — [claudesidian](https://github.com/heyitsnoah/claudesidian)
- claude-obsidian `wiki-mode` supports Generic, LYT, PARA, and Zettelkasten routing for new notes, without reorganizing existing ones. — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian)
- COG uses numbered domain folders (inbox, daily, personal, professional, projects, knowledge, templates). It cites PARA as an inspiration but does not implement it. — [COG](https://github.com/huytieu/COG-second-brain)
- One public vault uses a "modified PARA" with two non-PARA folders for the agent loop (including an inbox that gets drained). AGENTS.md is the contract, and CLAUDE.md is symlinked to it. — [alexanderweichart.de/AGENTS](https://alexanderweichart.de/AGENTS)
- A Chinese-language template combines PARA and Zettelkasten with numbered top-level folders. — [tomevault obsidian-ai-vault-template](https://tomevault.io/tome/iamzhihuix/obsidian-ai-vault-template) (the GitHub repo iamzhihuix/obsidian-ai-vault-template returned 404 on the API).
- The para-obsidian Claude Code plugin offers PARA vault management with frontmatter validation, template versioning, and git auto-commit. — [vibeindex](https://www.vibeindex.ai/plugins/nathanvale/side-quest-marketplace/para-obsidian)
- One practitioner argues for a flat vault with no folders, organized only by tags and wikilinks, as shared memory for Claude Code. — [sphelps substack](https://sphelps.substack.com/p/a-shared-memory-for-claude-code)
- Basic Memory gets its graph from note semantics (typed observations plus relations), not from folders. — [basic-memory](https://github.com/basicmachines-co/basic-memory)

### Inferences
- How the folders are organized matters less to agents than having an index file, consistent frontmatter, and typed links. An opinionated but swappable "organization profile" (PARA, JD, Zettelkasten, LLM-Wiki, flat) applied at bootstrap would generalize what claude-obsidian does only for routing new notes.
- Johnny Decimal's stable IDs (e.g. `21.04`) could make good agent-addressable keys. No one has used them that way, which is a small open niche.

### Gaps
- No agent template built on Johnny Decimal was found.
- "ACE" could mean Agentic Context Engineering (the 2025 paper on evolving context "playbooks") or the Obsidian "ACE" (Atlas/Calendar/Efforts) framework from Linking Your Thinking. No sources were collected for either in an agent bootstrap context, so this needs follow-up.

## Q4: CLI scaffolders (create-* / cookiecutter style) for agent workspaces

### Takeaway
No dominant `create-agent-workspace` scaffolder exists. What exists falls into four groups: (a) per-vault init scripts inside templates (claudesidian `/init-bootstrap`, claude-obsidian `init`/`adopt`, `obsidian-wiki setup`, `npx skills add` for COG), (b) SDD per-repo installers, (c) rule/config sync CLIs that fan one source of truth out to each agent's file (ruler, rulesync, agentsync, agentsmd), and (d) generic templating (cookiecutter). Nothing combines a structure scaffold, a memory layout, and multi-agent wiring for both a personal vault and many code repos.

### Cited Findings
- [intellectronica/ruler](https://github.com/intellectronica/ruler) (2,947 stars): "apply the same rules to all coding agents." — [GitHub](https://github.com/intellectronica/ruler)
- [dyoshikawa/rulesync](https://github.com/dyoshikawa/rulesync) (1,531 stars): "A Utility CLI for AI Coding Agents." — [GitHub](https://github.com/dyoshikawa/rulesync)
- agentsmd (Go) makes AGENTS.md the source of truth and symlinks CLAUDE.md and GEMINI.md to it. `--check` exits with an error if the state is wrong. — [joargp/agentmd](https://github.com/joargp/agentmd)
- agentsync (Rust crate) keeps the source of truth in `.agents/` and symlinks it into each tool's location. — [lib.rs agentsync](https://lib.rs/crates/agentsync)
- sync-agents-settings (`npx`) syncs MCP configs and instructions from Claude Code to Gemini CLI, Codex, OpenCode, Kiro, and Cursor (7 stars). — [GitHub](https://github.com/Leoyang183/sync-agents-settings)
- The AGENTS.md open format repo has 24,977 stars. — [agentsmd/agents.md](https://github.com/agentsmd/agents.md)
- Claude Code does not read AGENTS.md natively, so an import or symlink workaround is needed. Gemini CLI can @-import. — [inventivehq comparison](https://inventivehq.com/blog/claude-md-vs-agents-md-vs-gemini-md)
- cookiecutter (25,148 stars) is the generic project-template CLI. No agent-workspace cookiecutter of note turned up. — [cookiecutter](https://github.com/cookiecutter/cookiecutter)
- The template-embedded bootstrappers are claudesidian's `/init-bootstrap` + `/upgrade`, claude-obsidian's `python3 scripts/claude-obsidian.py init <path>` (plan → hash → apply) + `adopt`, `obsidian-wiki setup --vault ~/brain` (registers hooks), and COG's `npx skills add`. — [claudesidian](https://github.com/heyitsnoah/claudesidian), [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian), [obsidian-wiki](https://github.com/Ar9av/obsidian-wiki), [COG](https://github.com/huytieu/COG-second-brain)

### Inferences
- The scaffolding pieces exist but are scattered: rule sync (ruler/rulesync), vault seeds (claudesidian), safe apply (claude-obsidian), hooks (obsidian-wiki). A single `npx create-mindpalace` / `uvx` style CLI that does interview → plan → approve → apply for both a home vault and per-repo memory stubs would fill a real hole.
- The "template plus upgrade that keeps customizations" problem (claudesidian `/upgrade`, COG per-file approval) shows that keeping a scaffold in sync after the first run is a known pain point. A scaffolder should handle versioned migrations from day one.

### Gaps
- npm and PyPI were not searched directly for `create-*` agent-workspace packages, so a small one may exist that wasn't found.

## Q5: Common complaints and gaps (lock-in, bloat, staleness, cross-project memory, privacy)

### Takeaway
The documented pain points are: stale facts surfacing as confidently as true ones, CLAUDE.md bloat (static files load everything every session), context rot, manual upkeep, Claude-first design, weak isolation and sharing across projects, and privacy leaks through optional cloud embeddings and research APIs. Leading projects are starting to tackle staleness (claim ledgers, `last_verified`, the dated-fact rule) and multi-agent coordination (Crew, atomic writers). Cross-project memory and portable privacy controls are still mostly unsolved. Direct Reddit and HN complaint threads were hard to find, and most evidence comes from project docs and vendor blogs.

### Cited Findings
- **Staleness:** "a stale fact surfaces with the same fluent confidence as a true one." — [Mem0 library (vendor)](https://mem0.ai/library/agent-memory/memory-staleness). A memory-tool builder on Reddit argues that retrieving a stale fact is worse than no memory. — [Reddit snapshot](https://reddit.sentinel-team.org/posts/1u3qael/snapshots/2026-06-13T02%3A56%3A06.64152Z)
- **Bloat:** CLAUDE.md doesn't maintain itself, and because it is static, everything in it loads every time regardless of relevance. Flat markdown also handles relational knowledge poorly. — [AgentsRoom (vendor blog)](https://agentsroom.dev/blog/ai-agent-memory-solutions)
- **Context rot:** output quality drops as context grows, well before the window is full. — [Unblocked](https://getunblocked.com/blog/what-is-context-rot/index.md)
- **Cross-project isolation:** obsidian-second-brain says there is "no isolation within one vault yet." Multi-repo users point each repo at a separate vault through per-project `OBSIDIAN_VAULT_PATH`. — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain)
- **Privacy:** cloud embedding endpoints send note text off-machine, so obsidian-second-brain recommends excluding private folders or using Ollama. Its optional background agent runs with `--dangerously-skip-permissions` (sandboxed to file tools, off by default). — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain). claude-obsidian requires explicit consent before model-based reranking. — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian). Basic Memory's cloud tier is $15/mo. — [basic-memory](https://github.com/basicmachines-co/basic-memory)
- **Unreviewed agent writes:** one Obsidian plugin warns that CLAUDE.md is handed to an agent with read/write access to the vault, so changes should be reviewed through git or file history. — [claude-vault-assistant](https://community.obsidian.md/plugins/claude-vault-assistant)
- **Platform limits:** Agent SDK-based Obsidian plugins are desktop-only (no Node on mobile). — [Roasbeef/obsidian-claude-code](https://github.com/Roasbeef/obsidian-claude-code). Native Windows writes in claude-obsidian need WSL. — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian)
- **Lock-in signals:** many flagship repos are Claude-branded and Claude-first (claude-obsidian, claudesidian, claude-task-master, context-engineering-intro "Claude Code is the best," GSD "for Claude Code"). Cursor memory bank depends on Cursor custom modes. — GitHub descriptions: [context-engineering-intro](https://github.com/coleam00/context-engineering-intro), [GSD](https://github.com/gsd-build/get-shit-done), [cursor-memory-bank](https://github.com/vanzan01/cursor-memory-bank)
- **Mitigations already in the wild:** the OKM rule (timeless, dated, or pointer) — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain); `memory-hygiene` with `last_verified` — [COG](https://github.com/huytieu/COG-second-brain); claim ledgers and `lint --as-of` — [claude-obsidian](https://github.com/AgriciDaniel/claude-obsidian); a source manifest for delta ingestion and an atomic writer — [obsidian-wiki](https://github.com/Ar9av/obsidian-wiki); human-gated `keep` to promote lessons into rules — [obsidian-second-brain](https://github.com/eugeniughelbur/obsidian-second-brain).

### Inferences: gaps a new AI-agnostic bootstrapper could fill
1. **One model for a personal vault and many code repos.** Vault templates cover personal knowledge, and SDD and memory banks cover single repos. No one ties them together. A home vault (global memory, people, decisions) plus lightweight per-repo memory stubs that link back to it, with explicit scopes (global, project, session), is missing.
2. **Agent-neutral by construction, not by adapter.** Most projects bolt on multi-agent support after being built for Claude. A tool whose source of truth is AGENTS.md + Agent Skills + MCP, generating the CLAUDE.md, GEMINI.md, and Cursor rules files (as ruler and rulesync do), would avoid lock-in from the start.
3. **Bootstrapping as a first-class, safe, repeatable CLI.** Interview → plan → hash-approved apply → `adopt` existing folders → versioned `upgrade`/migrate. Today this exists only inside individual templates.
4. **Swappable organization profiles** (PARA, Johnny Decimal IDs, Zettelkasten, LYT/ACE, LLM-Wiki, flat), applied at bootstrap and switchable later. No one offers Johnny Decimal for agents.
5. **Freshness built in.** Required `as_of`/`last_verified`/source frontmatter, lint for stale and orphaned pages, a tiered "hot" set vs. an index to fight bloat, and a human gate before memory becomes a rule.
6. **Privacy defaults.** Local-only by default, per-folder "never send" markers respected by every agent adapter, and an audit log of what was sent to which model.
7. **A typed knowledge graph without a database.** Basic Memory-style observations and relations inside plain Obsidian markdown, readable by any agent without a running server, with MCP optional.

### Gaps
- No primary Reddit (r/ObsidianMD, r/ClaudeAI, r/LocalLLaMA) or HN threads with direct lock-in or bloat complaints were found. Searches surfaced mostly promotional blogs, so community sentiment is under-sourced.
- No numbers were found on how widely these memory systems are used or kept up (e.g. how many users keep running "update memory bank").
- Claude Code's own built-in memory features (auto-memory or a memory tool, if any as of Oct 2026) and those of Codex and Gemini CLI were not researched here. They affect how big the gap really is.
