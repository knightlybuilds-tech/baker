# Graphify and the Knowledge-Graph / Memory Tool Landscape for LLM Agents (as of Oct 2026)

Research date: 2026-10-18. Star counts are point-in-time snapshots from GitHub pages or aggregators fetched on this date unless noted; aggregator figures (GitTrend, dev.co, obsidianstats) may lag GitHub.

## What is graphify (author, repo, license, mechanics, agents, maturity, limitations)?

### Takeaway
"graphify" (the relevant one) is a very popular (~128k stars) Python CLI + agent skill, now under the Graphify-Labs GitHub org (originally safishamsi/graphify), that turns a folder (code, docs, PDFs, images, schemas) into a knowledge graph: code via deterministic tree-sitter AST parsing (no LLM), everything else via an LLM "semantic pass" using the host assistant's model. It emits graph.json, an interactive graph.html, GRAPH_REPORT.md, and optionally an Obsidian vault / wiki, and installs into 20+ coding agents. It is a *corpus-to-graph indexer*, not a conversational memory store.

### Cited Findings
**Identity / repo / author**
- Canonical repo now shows owner "Graphify-Labs" (organization); 127.9k stars, 12.3k forks, 400 watchers, 2,427 commits, a `v8` branch — [GitHub: Graphify-Labs/graphify](https://github.com/safishamsi/graphify) (fetched URL safishamsi/graphify resolved to Graphify-Labs page)
- Original README is served from the safishamsi/graphify path; sources disagree on owner (safishamsi vs Graphify-Labs), consistent with a transfer to an org — [jsDelivr README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md); [augmentcode summary](https://www.augmentcode.com/learn/graphify-knowledge-graph-codebase-skill)
- Repo description: "Turn any codebase, with its docs, SQL schemas, configs, and PDFs, into a queryable knowledge graph." — [GitHub](https://github.com/safishamsi/graphify)
- Releases are active; v0.9.48 notes list C++ parsing and graph-export bug fixes with contributor credits — [GitHub releases](https://github.com/safishamsi/graphify/releases)

**License**
- GitHub page lists both Apache-2.0 and MIT license files — [GitHub](https://github.com/safishamsi/graphify); GitTrend says Apache-2.0 — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify); another guide says MIT — [augmentcode](https://www.augmentcode.com/learn/graphify-knowledge-graph-codebase-skill). **Conflict: verify LICENSE file directly before depending on it.**

**Install**
- PyPI package is `graphifyy` (two y's); CLI and skill stay `graphify`. `uv tool install graphifyy` or `pipx install graphifyy`, then `graphify install` — [augmentcode](https://www.augmentcode.com/learn/graphify-knowledge-graph-codebase-skill); [GitHub](https://github.com/safishamsi/graphify)
- README warns other `graphify*` PyPI packages are not affiliated — [GitHub](https://github.com/safishamsi/graphify)
- For Claude Code, install writes `~/.claude/skills/graphify/SKILL.md`, adds a section to `~/.claude/CLAUDE.md`, and writes PreToolUse hooks to `.claude/settings.json` — [augmentcode](https://www.augmentcode.com/learn/graphify-knowledge-graph-codebase-skill)

**How it works**
- Code: tree-sitter AST parse + call-graph pass; deterministic, local, zero LLM cost; 40+ languages — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify); [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Docs (.md/.txt/.rst): concepts/relationships extracted by the LLM; PDFs: citation mining + concept extraction; images: vision model (screenshots, diagrams) — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Non-code passes use the assistant's model or a configured API key; `--code-only` runs fully offline — [GitHub](https://github.com/safishamsi/graphify)
- Leiden community detection (graspologic), "god node" detection; every edge tagged EXTRACTED / INFERRED / AMBIGUOUS — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- No vector store / no database: queries are graph traversal — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify); [horadecodar](https://horadecodar.com.br/?p=43673)
- SHA256 cache so re-runs only process changed files; `--watch` rebuilds code instantly (AST), docs/images need `--update`; `graphify hook install` adds post-commit hook — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)

**Outputs**
- `graph.html` (interactive), `GRAPH_REPORT.md` (god nodes, surprising connections, suggested questions), `graph.json` (persistent graph), `obsidian/` (Obsidian vault), `wiki/` (`--wiki`), `cache/`; optional `--svg`, `--graphml` (Gephi/yEd), `--neo4j` (cypher.txt) — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- `--mcp` starts an MCP stdio server; one guide says it exposes eight tools (structural queries, path tracing, node explanations) — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md); [claudelog](https://claudelog.com/claude-code-mcps/graphify)

**Agents supported**
- Early README: Claude Code only — [README (older)](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Current: "20+" — Claude Code, Cursor, Codex, Gemini CLI, GitHub Copilot CLI, VS Code Copilot Chat, Aider, OpenClaw, Factory Droid, Trae, Hermes, Kimi Code, Amp, Kiro, Pi, Devin CLI, Google Antigravity, CodeBuddy, OpenCode, Kilo Code, plus a generic Agent Skills target — [GitHub](https://github.com/safishamsi/graphify)

**Maturity / traction**
- On GitTrend daily trending 159 days since 2026-04-07, peaking #3 on 2026-04-14; ~9.3k stars in last 30 days — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify)
- Star counts in other snapshots ranged ~75k–82k (earlier) — [search summary of dev.co / skillrank](https://dev.co/ai/rag/graphify)

**Claims / benchmarks (vendor-reported)**
- "71.5x fewer tokens per query" vs reading raw files on a 52-file mixed corpus; 5.4x on 4 files; ~1x on a 6-file library (README concedes small corpora gain structure, not compression) — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- README reports LOCOMO QA 45.3% for graphify vs 49.7% for supermemory — [GitHub](https://github.com/safishamsi/graphify)

**Limitations**
- Inferred edges may need manual verification — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify)
- Graphs can get large on huge monorepos — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify)
- Ambiguous extensions (.inc, .h) can silently parse as wrong language unless declared in `.graphifyrc` — [GitHub](https://github.com/safishamsi/graphify)
- Per-agent gaps: Codex needs config for parallel extraction; OpenClaw/Aider sequential; Trae lacks PreToolUse hooks — [GitHub](https://github.com/safishamsi/graphify)
- Queries logged to `~/.cache/graphify-queries.log` by default (opt-out `GRAPHIFY_QUERY_LOG_DISABLE=1`) — [GitHub](https://github.com/safishamsi/graphify)
- "Runs entirely locally" claim, but docs/papers/images are sent to the LLM — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Commercial upsell: "Graphify Cloud" paid tier with 14-day trial promoted in README — [GitHub](https://github.com/safishamsi/graphify); company site [graphify.com blog](https://graphify.com/blog/how-to-give-claude-code-a-code-knowledge-graph)

### Inferences
- For an Obsidian-centric bootstrapper, graphify is attractive as a *derived-index generator*: its `obsidian/` and `wiki/` outputs and graph.json are files, not a DB, so they fit a "markdown is source of truth" design. But it is a build-time indexer of a corpus, not a read/write memory the agent updates turn by turn.
- The two-license listing plus Cloud upsell suggests open-core; check the LICENSE and whether features move to Cloud.
- The LOCOMO comparison is odd (graphify is not a conversational-memory system); treat as marketing.

### Gaps
- Could not confirm the individual author's real name / bio behind "safishamsi", or exact transfer date to Graphify-Labs.
- Could not find other notable, unrelated open-source "graphify" projects in search results (the README's warning about unaffiliated `graphify*` PyPI packages implies name collisions exist). graphify.com appears to be Graphify Labs' own commercial site, not a separate project — inferred from shared branding, not confirmed.
- No independent HN/Reddit critique found; no independent benchmark of token savings.
- Exact license (Apache-2.0 vs MIT vs dual) unresolved.

## How do alternatives compare (storage, local-first, model-agnostic, MCP, license, maturity)?

### Takeaway
The field splits into (a) markdown-native tools that keep Obsidian-compatible files as source of truth (Basic Memory, Letta Code MemFS, Obsidian MCP servers, Smart Connections, graphify outputs, Karpathy-style LLM wikis), and (b) DB-backed memory engines (Graphiti/Zep, Cognee, Mem0/OpenMemory, Supermemory, LightRAG, GraphRAG, reference MCP memory server JSONL) that are more capable at temporal/relational recall but opaque to humans. Basic Memory is the closest off-the-shelf fit for an Obsidian-based, AI-agnostic MCP memory layer.

### Cited Findings
**Basic Memory (basicmachines-co)**
- ~4.3k stars, AGPL-3.0; plain Markdown files on disk + local SQLite index (Postgres optional); Obsidian opens the same folder (wikilinks, frontmatter, graph view work) — [GitHub](https://github.com/basicmachines-co/basic-memory)
- MCP tools: write/read/edit/move/delete_note, search_notes, recent_activity, list_directory, build_context (graph navigation), schema tools; hybrid FTS + vector (FastEmbed) search, optional reranking — [GitHub](https://github.com/basicmachines-co/basic-memory)
- Clients: Claude Desktop/Code, Codex CLI, Cursor, VS Code, ChatGPT (Custom GPT actions), any MCP client; hosted cloud $15/mo for sync/mobile — [GitHub](https://github.com/basicmachines-co/basic-memory)
- Version references 0.18–0.23 (pre-1.0) — [GitHub](https://github.com/basicmachines-co/basic-memory)

**Graphiti (getzep) / Zep**
- ~31.9k stars, Apache-2.0; temporal knowledge graph with bi-temporal validity windows, fact invalidation (not deletion), episode provenance, incremental updates, hybrid semantic+BM25+traversal retrieval — [GitHub](https://github.com/getzep/graphiti)
- Requires a graph DB: Neo4j 5.26, FalkorDB (incl. embedded FalkorDB Lite), Amazon Neptune; Kuzu deprecated (upstream unmaintained) — [GitHub](https://github.com/getzep/graphiti)
- Default OpenAI; also Anthropic, Gemini, Groq, Azure, any OpenAI-compatible endpoint incl. Ollama/vLLM/LM Studio; but small/local models often produce schema-mismatched JSON causing extraction failures — [GitHub](https://github.com/getzep/graphiti)
- MCP server in `mcp_server/` (Docker + Neo4j); opt-out telemetry — [GitHub](https://github.com/getzep/graphiti)
- Zep is the managed product with users/threads/enterprise support — [GitHub](https://github.com/getzep/graphiti)

**Cognee (topoteretes)**
- ~32.5k stars, Apache-2.0, v1.6.1 (2026-09-24); remember/recall/improve/forget API; graph + vector + relational storage (can run all on one Postgres) — [GitHub](https://github.com/topoteretes/cognee)
- Runs locally with no API key (local GLiNER extraction + local embeddings), Ollama supported; `cognee-mcp` server; plugins for Claude Code, Codex, OpenClaw; Python/TS/Rust/REST SDKs — [GitHub](https://github.com/topoteretes/cognee)

**Mem0 / OpenMemory**
- Mem0: ~67.4k stars, Apache-2.0; multi-level (user/session/agent) memory, single-pass add-only extraction, semantic+BM25+entity retrieval; default LLM gpt-5-mini, embedder text-embedding-3-small; hosted platform; vendor benchmarks (LoCoMo 92.5, LongMemEval 94.4) are for managed platform with optimizations not in OSS — [GitHub](https://github.com/mem0ai/mem0)
- Memory fragments are "optimized for retrieval by software, not documents you would read" — [hjarni comparison](https://hjarni.com/blog/best-mcp-memory-server)
- OpenMemory MCP (Mem0, launched 2025-05-13): local Docker + Postgres + Qdrant; tools add_memories, search_memory, list_memories, delete_all_memories; hosted variant exists — [Mem0 blog](https://mem0.ai/blog/introducing-openmemory-mcp); [docs](https://docs.mem0.ai/openmemory)
- Name collision: CaviraOSS/OpenMemory is separate (Apache-2.0, ~4.3k stars, SQLite/Postgres, multi-sector memory, MCP) — [dev.co](https://dev.co/ai/frameworks/openmemory)

**Letta (MemGPT) / Letta Code**
- letta-ai/letta: ~25.2k stars, Apache-2.0; current development centered on `letta-ai/letta-code` (agent harness, TUI, desktop/web runtime) — [GitHub](https://github.com/letta-ai/letta)
- Letta Code MemFS: memory is a directory of markdown files in a local git repo, on by default for new agents (not available with Docker server); every change auto-committed; sleep-time reflection/defrag subagents use git worktrees — [Letta MemFS docs](https://docs.letta.com/letta-code/memfs); [Context Repositories blog (Feb 2026)](https://www.letta.com/blog/context-repositories)
- Letta announced deprecation of older features (templates, filesystem) by end of April 2026 — [Letta "our next phase"](https://www.letta.com/blog/our-next-phase)

**Microsoft GraphRAG**
- ~36.4k stars, MIT; in maintenance mode (bug fixes/CVE only, no new features/PRs); README warns indexing "can be an expensive operation" — [GitHub](https://github.com/microsoft/graphrag)

**LightRAG (HKUDS)**
- ~40.2k stars, MIT; default storages JSON KV / NanoVectorDB / NetworkX (in-memory, testing only), Postgres recommended for production; Neo4j/Memgraph/Milvus/Qdrant/Mongo options; local/global/hybrid/naive/mix modes; WebUI + REST; no MCP mentioned in README — [GitHub](https://github.com/HKUDS/LightRAG)

**Supermemory**
- ~31.3k stars, MIT; hosted + local single binary (`npx supermemory local`, data in `~/.supermemory`); BYO LLM (OpenAI/Anthropic/Gemini/Groq/compatible), local embeddings, Ollama; MCP tools memory/recall/context; connectors (Drive, Gmail, Notion, GitHub); no markdown/Obsidian mention — [GitHub](https://github.com/supermemoryai/supermemory)

**Reference MCP Memory server (modelcontextprotocol/servers)**
- MIT; stores a knowledge graph (entities, relations, observations) in a single `memory.jsonl` file; tools create/delete entities/relations/observations, read_graph, search_nodes, open_nodes — [GitHub](https://github.com/modelcontextprotocol/servers/tree/main/src/memory)

**Obsidian MCP servers / Smart Connections**
- No official Obsidian MCP server; all community-built. Two architectures: in-app (via plugin, needs Obsidian running, accesses live state) vs filesystem (reads markdown from disk, headless) — [skillselion overview](https://skillselion.com/best/obsidian-mcp-servers)
- mcp-obsidian is built on the Obsidian Local REST API plugin, fully local — [enterprisedna listing](https://enterprisedna.co/directories/mcp/mcp-obsidian); aaronsb/obsidian-mcp-plugin is an in-app option — [GitTrend](https://gittrend.io/repo/aaronsb/obsidian-mcp-plugin)
- Community MCP servers expose Smart Connections embeddings (one read-only, works without Obsidian running) — [glama](https://glama.ai/mcp/servers/ksu32q43os); [getdrio](https://www.getdrio.com/mcp/io-github-gogogadgetbytes-smart-connections/md)

**Khoj**
- Obsidian plugin ~48k downloads; open-source, chat with local or online LLMs over docs — [Obsidian community](https://community.obsidian.md/plugins/khoj)
- ~34k GitHub stars; AGPL-3.0 per aggregator; self-hostable — [obsidianstats](https://www.obsidianstats.com/plugins/khoj); [x-cmd](https://x-cmd.com/install/khoj)
- README now promotes "Pipali," a new open-source AI coworker from the same team — [search summary of Khoj README] (unverified directly)

**Other 2025–2026 entrants**
- Engram (memory MCP): markdown as source of truth + rebuildable index; vendor benchmark claims 44% fewer tool calls and 35% fewer tokens than plain markdown + Read/Grep — [mcpservers.org](https://mcpservers.org/ru/servers/veronchenko/engram-memory)
- fellowgeek/mcp-memory (Show HN): markdown+frontmatter (OKF) + SQLite FTS5 — [HN](https://hn.svelte.dev/item/49286073)
- Karpathy "LLM Wiki" gist (April 2026, no code): drop sources in a folder, LLM compiles interlinked markdown wiki, browse in Obsidian; retrieval via LLM reading its own index/summary files, not a vector DB — [Implicator (14 Apr 2026)](https://implicator.ai/karpathys-llm-wiki-got-16-million-views-the-code-isnt-the-point-2/); [agentpedia](https://agentpedia.codes/blog/karpathy-llm-wiki-idea-file); [TiddlyWiki forum](https://talk.tiddlywiki.org/t/riding-the-wave-of-andrej-karpathys-llm-wiki/15214)

**Summary table** (from findings above)

| Tool | Storage | Local-first | Model-agnostic | MCP | License | Stars (Oct 2026) |
|---|---|---|---|---|---|---|
| graphify | graph.json + HTML + optional Obsidian vault/wiki (files) | Yes (code offline; docs need LLM) | Uses host agent's model; 20+ agents | Yes (`--mcp`) | Apache-2.0/MIT (unclear) | ~128k |
| Basic Memory | Markdown + SQLite index | Yes | Yes (any MCP client) | Yes (core) | AGPL-3.0 | ~4.3k |
| Letta Code MemFS | Markdown in git repo | Yes (Local mode) | Letta harness | Letta-specific | Apache-2.0 | ~25k (letta) |
| Graphiti | Neo4j/FalkorDB/Neptune | Self-host, needs DB | Yes but needs structured-output LLM | Yes | Apache-2.0 | ~32k |
| Cognee | Graph+vector+relational DBs | Yes, keyless local mode | Yes | Yes | Apache-2.0 | ~32.5k |
| Mem0 / OpenMemory | Vector store (Qdrant) + Postgres | Self-host possible | Configurable, OpenAI default | Yes | Apache-2.0 | ~67k |
| Supermemory | Local dir / hosted | Local binary | Yes (BYO LLM) | Yes | MIT | ~31k |
| LightRAG | JSON/NetworkX dev; Postgres prod | Yes | Yes | Not in README | MIT | ~40k |
| GraphRAG | Index files (parquet; not confirmed in fetch) | Yes | — | No | MIT, maintenance mode | ~36k |
| MCP ref memory | Single JSONL | Yes | Yes | Yes | MIT | — |
| Obsidian MCP servers | Vault markdown | Yes | Yes | Yes | Various (community) | small (hundreds) |
| Khoj | Own index over docs | Self-host | Yes | Not confirmed | AGPL-3.0 | ~34k |

### Inferences
- For an AI-agnostic Obsidian bootstrapper, the natural stack is: vault markdown as source of truth + an MCP server over the vault (Basic Memory or a filesystem Obsidian MCP) + optional derived indexes (graphify for code/corpus graphs, Smart Connections embeddings, FTS). DB engines (Graphiti, Cognee, Mem0) are better reserved as optional plug-ins for conversational/temporal memory.
- Basic Memory's AGPL-3.0 matters if the bootstrapper bundles/modifies it as a network service; calling it as a separate MCP process is lower-risk (legal review advisable).
- Letta's move from DB memory blocks to git-backed markdown (Feb 2026) is a strong signal that even memory-first vendors see value in human-readable files.

### Gaps
- Did not verify Smart Connections plugin's own current version/license/storage details (it stores embeddings in a `.smart-env` folder per prior knowledge — unconfirmed this session).
- Khoj MCP support and current release date unconfirmed (sources disagree).
- GraphRAG output format (parquet) not confirmed in fetched README.
- mcp-obsidian star counts and maintenance status not confirmed.
- Zep cloud pricing/current state not researched.

## Which keep data as human-readable markdown (Obsidian-compatible) vs opaque databases?

### Takeaway
Markdown-native: Basic Memory, Letta Code MemFS, Obsidian MCP servers, Karpathy-style LLM wikis, Engram / mcp-memory style tools, and graphify's optional `obsidian/`/`wiki/` export (graph.json itself is machine JSON). Opaque/DB: Graphiti (graph DB), Cognee, Mem0/OpenMemory (Qdrant+Postgres), Supermemory, LightRAG, GraphRAG, Khoj's index; the MCP reference server uses a readable-but-not-Obsidian JSONL.

### Cited Findings
- Basic Memory: plain markdown on disk, Obsidian reads same files — [GitHub](https://github.com/basicmachines-co/basic-memory)
- Letta MemFS: directory of markdown files in local git repo — [Letta docs](https://docs.letta.com/letta-code/memfs)
- graphify outputs `obsidian/` vault and `wiki/` alongside graph.json — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Graphiti requires Neo4j/FalkorDB/Neptune — [GitHub](https://github.com/getzep/graphiti)
- Mem0 fragments not meant to be read as documents — [hjarni](https://hjarni.com/blog/best-mcp-memory-server); OpenMemory uses Postgres + Qdrant — [Mem0 blog](https://mem0.ai/blog/introducing-openmemory-mcp)
- LightRAG default JSON/NetworkX storages, Postgres for production — [GitHub](https://github.com/HKUDS/LightRAG)
- MCP reference memory: JSONL entities/relations/observations — [GitHub](https://github.com/modelcontextprotocol/servers/tree/main/src/memory)
- Engram: markdown source of truth + rebuildable index — [mcpservers.org](https://mcpservers.org/ru/servers/veronchenko/engram-memory)

### Inferences
- A "markdown source of truth + rebuildable derived index" pattern (Basic Memory, Engram, graphify-as-index) best preserves portability and AI-agnosticism: any agent can grep the vault even with no MCP server running.

### Gaps
- Whether graphify's Obsidian export round-trips (i.e., can edits in Obsidian feed back into graph.json) is not documented in the sources found; it appears one-way.

## Trade-offs: graph-based memory vs plain markdown + wikilinks + grep/embeddings

### Takeaway
Graphs win on multi-hop/relational and temporal queries over large corpora and can cut tokens (vendor-reported), but cost LLM extraction, add infra and schema lock-in, can contain inferred/hallucinated edges, and go stale. Markdown + wikilinks + grep is portable, diffable, human-auditable and "good enough" for personal-scale memory; indexes (FTS/embeddings) are added when token cost of scanning grows. Hybrid (markdown truth + derived index/graph) is the emerging consensus.

### Cited Findings
- Pro-graph: Graphiti's bi-temporal facts and invalidation enable historical queries — [GitHub](https://github.com/getzep/graphiti); graphify claims 71.5x token reduction on 52 files but ~1x on small corpora — [README](https://cdn.jsdelivr.net/gh/safishamsi/graphify@main/README.md)
- Graph cost: GraphRAG warns indexing is expensive — [GitHub](https://github.com/microsoft/graphrag); Graphiti extraction fails with small local models — [GitHub](https://github.com/getzep/graphiti); graphify inferred edges need verification — [GitTrend](https://gittrend.io/repo/Graphify-Labs/graphify)
- Schema ownership/lock-in argument favoring portable, diffable plain text; graphs defended for multi-user, relationship-heavy servers — [themenonlab](https://themenonlab.blog/blog/markdown-vs-graph-database-agent-memory-soul-py-openlobster)
- HN: index reduces O(N) scan, saves tokens at hundreds of notes; MCP works where no filesystem (claude.ai). Counter: MCP tools "tend to slow down the agent and waste tokens"; "about the same as grep in a memory/ directory"; facts go stale; hand-maintained MEMORY.md works — [HN Show HN: MCP Memory](https://hn.svelte.dev/item/49286073)
- Engram benchmark (vendor, synthetic): indexed recall used 44% fewer tool calls / 35% fewer tokens than plain markdown + grep — [mcpservers.org](https://mcpservers.org/ru/servers/veronchenko/engram-memory)
- Karpathy LLM Wiki: ~100 articles / 400k words navigated via LLM-maintained index files without vector DB (figures repeated secondhand, not primary-confirmed) — [Implicator](https://implicator.ai/karpathys-llm-wiki-got-16-million-views-the-code-isnt-the-point-2/)
- Benchmarks on conversational memory: graphify 45.3% vs supermemory 49.7% on LOCOMO (graphify README); Mem0 92.5 LoCoMo on managed platform only — [GitHub graphify](https://github.com/safishamsi/graphify); [GitHub mem0](https://github.com/mem0ai/mem0)

### Inferences
- For the bootstrapper: default to markdown + wikilinks + frontmatter (Obsidian-native graph for free), give agents grep + an MCP vault server, add embeddings/FTS for scale, and treat graph builders (graphify for code; Graphiti/Cognee for temporal conversational memory) as optional, rebuildable layers rather than the store of record.
- Obsidian's wikilink graph is itself a lightweight human-curated knowledge graph; LLM-extracted graphs add breadth but reduce trust (EXTRACTED vs INFERRED tagging helps).

### Gaps
- No independent, apples-to-apples benchmark comparing graph memory vs markdown+grep for coding agents at scale was found; all numbers are vendor-run.
