# Category Vocabulary and Architecture Patterns for AI Agent Memory Bootstrappers (as of Oct 2026)

## Terms in use, their nuance, and momentum

### Takeaway
There is no single settled category name. "Context engineering" is the umbrella discipline term (established by Anthropic, Manus, OpenAI cookbooks in 2025-26). "Agent memory layer" is the established product-category term for API/DB-style memory services (Mem0, Zep, Letta). "LLM Wiki" (Karpathy, April 2026) and "second brain for AI" / "AI second brain" are the momentum terms for file-based, markdown, human-readable knowledge bases maintained by agents. "Memory bank" is a Cline-specific convention name. AGENTS.md / CLAUDE.md / SKILL.md are the de facto file standards a bootstrapper should target.

### Cited Findings
- **Context engineering**: Anthropic (Sept 29, 2025) defines it as "the set of strategies for curating and maintaining the optimal set of tokens" during inference, presented as the successor to prompt engineering. — [Anthropic: Effective context engineering for AI agents](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)
- Manus published "Context Engineering for AI Agents: Lessons from Building Manus" (Yichao 'Peak' Ji, July 18, 2025), a widely cited practitioner post using the same term. — [Manus blog](https://manus.im/blog/Context-Engineering-for-AI-Agents-Lessons-from-Building-Manus)
- OpenAI also uses the term in its cookbook "Context Engineering for Personalization - State Management with Long-Term Memory Notes". — [OpenAI Cookbook](https://developers.openai.com/cookbook/examples/agents_sdk/context_personalization)
- Shopify CEO Tobi Lütke is credited (secondary source only, unverified) with popularizing "context engineering". — [Rework profile](https://resources.rework.com/de/libraries/leadership-legends/tobi-lutke-leadership)
- **Agent memory layer**: 2026 comparison guides use "memory layer for AI agents" as the category label, listing Mem0, Letta, Zep, Graphiti, Cognee; Mem0 describes itself as "a memory layer for AI applications". — [Powerdrill: 10 Best Free Memory layer for AI agents in 2026](https://powerdrill.ai/blog/best-free-memory-layer-for-ai-agents); [Developers Digest: Mem0 vs Zep vs Letta vs Cloudflare](https://www.developersdigest.tech/blog/best-ai-agent-memory-providers-2026)
- These products are "not interchangeable": one extracts facts (Mem0), one builds a temporal knowledge graph (Zep), one lets the agent edit its own memory (Letta), Cloudflare Durable Objects are substrate. — [Developers Digest](https://www.developersdigest.tech/blog/best-ai-agent-memory-providers-2026)
- **LLM Wiki (Karpathy)**: Gist created April 4, 2026; an "idea file" describing an LLM that "incrementally builds and maintains a persistent wiki" of interlinked markdown, as an alternative to per-query RAG. Vocabulary in the gist: "compounding artifact", "Obsidian is the IDE", "the LLM is the programmer", "the wiki is the codebase", nod to Vannevar Bush's Memex. — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- The gist spawned implementations, including an MIT-licensed `karpathy-llm-wiki` agent skill distributed via skills marketplaces. — [skills.sh listing](https://www.skills.sh/learnprompt/andrej-karpathy-skills/karpathy-llm-wiki); critics call it a relabeled caching layer and note no RBAC/ACID/concurrency for enterprise. — [Towards AI](https://pub.towardsai.net/andrej-karpathy-killed-rag-or-did-he-the-llm-wiki-pattern-7824d876e790)
- Press framed the LLM Wiki as Karpathy "build[ing] a second brain" — the two terms are now conflated in coverage. — [Medium/NeuralNotions](https://medium.com/neuralnotions/andrej-karpathy-stopped-using-ai-to-write-code-hes-using-it-to-build-a-second-brain-instead-cddceadc5df5); [Codersera](https://codersera.com/blog/karpathy-llm-knowledge-base-second-brain/amp/)
- **Second brain for AI / AI-native PKM**: An active 2026 scene pairs Obsidian vaults with Claude Code/Codex as an "AI second brain"; Tiago Forte reportedly (via a secondary post) wrote in March 2026 that "Personal Context Management is taking over Personal Knowledge Management." "AI-native PKM" is used loosely with no authoritative definition. — [Pasquale Pillitteri](https://pasqualepillitteri.it/en/news/962/obsidian-claude-code-second-brain-persistent-memory); [HackerNoon](https://hackernoon.com/obsidian-holds-the-notes-claude-does-the-work-second-brain); [Obsidian Second Brain project](https://gittrend.io/repo/eugeniughelbur/obsidian-second-brain)
- **Memory bank**: Cline's documented convention — a `memory-bank/` folder with projectbrief.md, productContext.md, systemPatterns.md, techContext.md, activeContext.md, progress.md, initialized via a `.clinerules` file. — [Cline docs: Memory Bank](https://docs.cline.bot/prompting/cline-memory-bank); [Cline blog](https://cline.bot/blog/memory-bank-how-to-make-cline-an-ai-agent-that-never-forgets)
- Critique: committed memory-bank files cause git merge conflicts; git-ignored ones are trapped on one machine (third-party view). — [fast.io comparison](https://www.fast.io/resources/cline-claude-code-memory-comparison.md)
- **AGENTS.md** became a founding project of the Linux Foundation's Agentic AI Foundation (Dec 9, 2025) alongside MCP and goose; OpenAI claims 60,000+ projects adopted it (Campus Technology reports 40,000 — figures conflict). — [OpenAI: Agentic AI Foundation](https://openai.com/blog/agentic-ai-foundation); [InfoQ](https://infoq.com/news/2025/12/agentic-ai-foundation/); [Campus Technology](https://campustechnology.com/articles/2025/12/15/newly-launched-agentic-ai-foundation-brings-together-tech-giants-for-open-source-ai-development.aspx?p=1)
- **Agent Skills (SKILL.md)**: Anthropic published Agent Skills as an open standard (agentskills.io, Dec 18, 2025); adoption reported at 26+ to 40+ platforms including OpenAI Codex (counts vary). — [Unite.AI](https://www.unite.ai/anthropic-opens-agent-skills-standard-continuing-its-pattern-of-building-industry-infrastructure); [Atlan explainer](https://atlan.com/know/ai-agent/ai-agent-skills/what-are-agent-skills/)

### Inferences
- Best positioning vocabulary for a tool that *bootstraps folders + persistent memory*: lead with a concrete pattern people already search for ("LLM wiki", "AI second brain", "memory bank") and describe the mechanism with the discipline term ("context engineering", "file-based agent memory"). "Agent memory layer" signals an API/vector-DB product, which a scaffolding tool is not — using it invites comparison with Mem0/Zep.
- "Workspace scaffolding" and "agentic knowledge base" are descriptive but did not surface as momentum terms in searches.
- Being "AI-agnostic" is best signaled by targeting AGENTS.md (vendor-neutral, LF-governed) plus SKILL.md and MCP, with CLAUDE.md as a symlink/alias.

### Gaps
- No Google Trends / HN-count data was retrieved; momentum judgments rest on volume and recency of coverage, not measured search trends.
- "Personal context protocol" returned no meaningful usage — appears not to be an established term (good for naming freedom, poor for discoverability).
- Could not confirm Lütke/Karpathy primary posts on coining "context engineering" (June 2025 X posts are the commonly cited origin, but not verified here).

## Memory taxonomy for agents

### Takeaway
The dominant vocabulary is the cognitive-science split (working / episodic / semantic / procedural), popularized by LangChain/LangMem (derived from the CoALA paper). Letta uses an OS metaphor (core memory blocks in-context vs. archival/recall memory out-of-context). Anthropic is deliberately simple: a file directory (`/memories`) the model CRUDs, plus context editing and compaction. The Dec 2025 academic survey reframes as forms/functions/dynamics.

### Cited Findings
- **LangMem**: semantic (facts; stored as a profile or collection), episodic (past experiences, few-shot examples, conversation summaries), procedural (system behavior, stored as prompt rules). — [LangMem conceptual guide](https://langchain-ai.github.io/langmem/concepts/conceptual_guide/)
- LangMem distinguishes **collections** (growing records that must be reconciled—delete/invalidate/consolidate) vs **profiles** (single schema'd document updated in place, easy for humans to edit). — [LangMem](https://langchain-ai.github.io/langmem/concepts/conceptual_guide/)
- LangMem distinguishes **hot-path** memory formation (during conversation; immediate but adds latency) vs **background** formation (after conversation/inactivity; no latency, higher recall, delayed). Memories are organized in hierarchical **namespaces** (org/user/app). — [LangMem](https://langchain-ai.github.io/langmem/concepts/conceptual_guide/)
- **Letta (MemGPT lineage)**: OS-inspired design; context window as "main memory", external storage as "external memory", paging in/out; core memory is "stored in-context composed of memory blocks". — [Developers Digest](https://www.developersdigest.tech/blog/best-ai-agent-memory-providers-2026); [Letta: Sleep-time compute](https://www.letta.com/blog/sleep-time-compute)
- **Anthropic memory tool** (`memory_20250818`): client-side file directory under `/memories` with commands view/create/str_replace/insert/delete/rename; Claude checks memory before starting a task; storage is entirely in the developer's app. — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- Anthropic pairs memory with **context editing** (client clears specific tool results) and **compaction** (server summarizes conversation near the limit): "compaction keeps the active context small... and memory preserves the information that must survive summarization." — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- Reported: memory tool + context editing improved agentic search performance 39% (context editing alone 29%) and cut token use 84% in a 100-turn web search eval (Anthropic figures via press). — [The Decoder](https://the-decoder.com/anthropic-claims-context-engineering-beats-prompt-engineering-when-managing-ai-agents/)
- **OpenAI**: Cookbook pattern uses a local state object (structured profile + notes); session notes distilled during a run, consolidated into global notes at end, injected at next run start with precedence rules; no embeddings. — [OpenAI Cookbook](https://developers.openai.com/cookbook/examples/agents_sdk/context_personalization); [aihola write-up](https://aihola.com/article/openai-agents-sdk-memory-personalization)
- **Academic survey** "Memory in the Age of AI Agents: A Survey, Forms, Functions and Dynamics" (arXiv 2512.13564, Dec 15, 2025): forms = token-level, parametric, latent; functions = factual, experiential, working; dynamics = formation, evolution, retrieval. Explicitly separates agent memory from RAG and context engineering. — [arXiv](https://arxiv.org/pdf/2512.13564); [HF Papers](https://huggingface.co/papers/2512.13564)
- Alternative taxonomy: "Rethinking Memory in AI" (arXiv 2505.00675) — parametric vs contextual, six atomic operations. — [arXiv](https://arxiv.org/html/2505.00675v1)

### Inferences
- For a file-based bootstrapper, a practical mapping: working memory = session scratch/todo files; episodic = append-only log.md / session journals; semantic = wiki/entity/concept pages and profile files; procedural = AGENTS.md/CLAUDE.md rules + SKILL.md skills. This lets the project speak both the CoALA/LangMem vocabulary and the file conventions.
- "Profile vs collection" maps neatly to "single curated file (e.g., user.md) vs folder of notes".

### Gaps
- Did not fetch Letta's own memory docs (archival/recall terminology) directly; Letta descriptions are via its sleep-time post and third-party guides.
- Did not verify the CoALA paper link directly.

## Best practices: file-based memory, progressive disclosure, compaction, consolidation, staleness

### Takeaway
Experts converge on: (1) keep the context small and high-signal (attention budget / context rot); (2) store memory as files on disk the agent reads just-in-time, with a small always-loaded index/instruction file and deeper files loaded on demand (progressive disclosure); (3) append-only logs plus curated, periodically consolidated summaries; (4) background "sleep-time" or "lint" passes to consolidate, dedupe and flag stale/contradictory claims; (5) explicit initializer + end-of-session update rituals.

### Cited Findings
- **Attention budget / context rot**: recall degrades as context fills; aim for "the smallest set of high-signal tokens". — [Anthropic](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)
- **Just-in-time retrieval**: keep lightweight identifiers (paths, queries, links) and load data at runtime; Claude Code uses a hybrid — CLAUDE.md loaded up front, glob/grep on demand. Folder hierarchy, naming conventions, and timestamps act as relevance signals. — [Anthropic](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)
- **Structured note-taking**: agent writes notes (e.g., NOTES.md, to-do lists) outside the context window and reloads later; pick compaction for long back-and-forth, note-taking for milestone work, sub-agents (returning ~1,000-2,000 token summaries) for parallel research. — [Anthropic](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)
- **Compaction**: preserve architectural decisions and unresolved bugs, drop redundant tool output; tune for recall first, then precision; tool-result clearing is the lightest form. — [Anthropic](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents)
- **Progressive disclosure (Skills)**: only name+description (~30-50 tokens per skill) load at startup; full SKILL.md on trigger; reference files only when needed. — [Atlan](https://atlan.com/know/ai-agent/ai-agent-skills/what-are-agent-skills/); [SwirlAI](https://www.newsletter.swirlai.com/p/agent-skills-progressive-disclosure?open=false)
- **Manus**: file system as "the ultimate context" — compression must be restorable (drop a page but keep its URL/path); KV-cache hit rate is "the single most important metric" so keep prompt prefix stable and context append-only (a timestamp at the top of the system prompt breaks caching); recitation via a rewritten todo.md keeps goals in recent attention; keep errors in context; add variation to avoid few-shot ruts; mask tools rather than remove them. — [Manus](https://manus.im/blog/Context-Engineering-for-AI-Agents-Lessons-from-Building-Manus)
- **Karpathy LLM Wiki**: three layers — immutable raw sources, LLM-maintained wiki, schema file (CLAUDE.md/AGENTS.md, "the key configuration file", co-evolving with the user). Operations: ingest (one source may touch 10-15 pages), query (file useful answers back as pages), lint (contradictions, stale claims, orphan pages, missing cross-refs). `index.md` (catalog with one-line summaries, read first) and `log.md` (append-only, consistent prefixes for grep). Index approach works to ~100 sources / few hundred pages without embeddings; beyond that add a local search tool (e.g., qmd). On contradictions: note the conflict and update synthesis rather than silently overwrite. — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- **Letta sleep-time compute** (Apr 21, 2025): a separate sleep-time agent owns memory-editing tools and refines the primary agent's memory blocks in the background ("clean, concise, and detailed memories"); can use a different model and run on a schedule; authors claim a Pareto improvement (self-reported). — [Letta](https://www.letta.com/blog/sleep-time-compute)
- **Anthropic multisession pattern**: initializer session creates progress log, feature checklist, init script reference; each later session reads them first; end-of-session update; mark features complete only after end-to-end verification. — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool); see also [Effective harnesses for long-running agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents)
- Anthropic's auto-injected memory protocol tells the model to "ASSUME INTERRUPTION" and record progress; suggested reinforcement: keep memory "up-to-date, coherent and organized", rename/delete irrelevant files, "Do not create new files unless necessary." — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- **Staleness / expiration**: "Periodically delete memory files that haven't been accessed in a long time." — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- OpenAI cookbook: distill session notes → consolidate into global notes at session end → inject with precedence rules. — [OpenAI Cookbook](https://developers.openai.com/cookbook/examples/agents_sdk/context_personalization)

### Inferences
- MVP design implied by consensus: a generated folder with (a) AGENTS.md (+ CLAUDE.md alias) as the always-loaded schema/router, kept short; (b) `index.md` catalog; (c) `log.md` append-only episodic log; (d) topic/entity pages (semantic); (e) `skills/` with SKILL.md (procedural, progressive disclosure); (f) `raw/` immutable sources; (g) a "lint/consolidate" command or skill that runs between sessions (sleep-time analogue) to dedupe, flag contradictions, mark stale entries with dates.
- Put dates/"last verified" metadata on memory entries so lint can detect staleness; keep the always-loaded file stable to preserve prompt-cache hits (Manus lesson).
- Plain markdown in git gives human readability, diffs and rollback — a natural defense against memory corruption, though it creates merge-conflict friction (Cline critique).

### Gaps
- No quantitative evidence found comparing file-based vs vector-DB memory for coding/knowledge agents; benchmark claims in the space (e.g., MemPalace LongMemEval) are disputed.

## Risks: prompt injection via memory, secrets in memory, memory bloat

### Takeaway
Memory poisoning is now a named top-10 risk (OWASP Agentic ASI06, "Memory and Context Poisoning"): it is persistent prompt injection that survives sessions. Anthropic's docs list concrete mitigations: path-traversal protection, stripping sensitive data before writes, file-size caps, and expiration. Bloat is handled with caps, consolidation, and pruning.

### Cited Findings
- OWASP Top 10 for Agentic Applications (2026) includes Memory & Context Poisoning as ASI06, covering tampering with conversation history, vector stores, scratchpads, RAG indexes (via vendor summaries; verify wording against OWASP primary doc). — [Kiteworks](https://www.kiteworks.com/cybersecurity-risk-management/owasp-agent-memory-poisoning-guard/); [Christian Schneider](https://christian-schneider.net/blog/persistent-memory-poisoning-in-ai-agents/)
- Poisoned entries persist across sessions/restarts and can trigger days or weeks later; research attack families include MINJA, AgentPoison. Vendor-reported attack success rates (e.g., 84.3% avg from Agent Security Bench) — treat as vendor-cited. — [tianpan.co](https://tianpan.co/blog/2026-04-10-agent-memory-poisoning-persistent-compromise); [Vectorize](https://vectorize.io/articles/ai-memory-poisoning)
- Anthropic: "Claude usually refuses to write sensitive information to memory files. For stronger guarantees, add validation that strips sensitive data before your handler writes the file." — [Claude docs: Memory tool](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- Anthropic: validate every path to block traversal like `/memories/../../secrets.env`, including URL-encoded `%2e%2e%2f`; use `Path.resolve()`/`relative_to()`. — [Claude docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- Anthropic: cap file sizes and characters returned by `view`; page with `view_range`. — [Claude docs](https://platform.claude.com/docs/en/agents-and-tools/tool-use/memory-tool)
- Karpathy LLM Wiki keeps raw sources immutable and separate from LLM-written wiki — provenance separation. — [Karpathy gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- MemPalace security note: impostor domains distributing malware — supply-chain risk for popular memory tools. — [Vectorize: What is MemPalace](https://vectorize.io/articles/what-is-mempalace)

### Inferences
- Bootstrapper mitigations to ship by default: `.gitignore` + secret-scanning pre-commit (e.g., gitleaks) on the memory folder; a "never store secrets" rule in AGENTS.md; treat content in `raw/` and web-ingested material as untrusted data (provenance tags on wiki pages citing sources); human-reviewable diffs via git; size caps and lint warnings for oversized files; quarantine/inbox folder for unreviewed ingests.

### Gaps
- Did not access the OWASP primary document directly.
- No empirical data found on memory bloat thresholds for file-based systems beyond Karpathy's ~100 sources / few hundred pages heuristic and Anthropic's 16,000-character view truncation.

## Positioning: naming conventions and category leaders

### Takeaway
Category leaders by sub-segment: Mem0 ("memory layer"), Letta (stateful agents/MemGPT), Zep/Graphiti (temporal knowledge graph), Cline Memory Bank (convention), Karpathy LLM Wiki (pattern), Obsidian+Claude Code "second brain" (community). A critical finding: **MemPalace**, a viral open-source "memory palace" AI memory project (April 2026), occupies the name space closest to "mindpalace".

### Cited Findings
- MemPalace: MIT-licensed, local-first AI memory by Milla Jovovich and Ben Sigman; verbatim storage + semantic search; hierarchy of "wings/rooms/drawers" (method of loci metaphor); Claude Code hooks; MCP integration; reported ~19.5k GitHub stars in first week (April 2026) and ~57k later (aggregator figures); benchmark claims (100% LongMemEval R@5) contested. — [Better Stack guide](https://betterstack.com/community/guides/ai/mempalace-ai-memory/); [Let's Data Science](https://letsdatascience.com/news/milla-jovovich-launches-open-source-ai-memory-system-ac41e53c); [Vectorize](https://vectorize.io/articles/what-is-mempalace); [noqta](https://noqta.tn/en/blog/mempalace-open-source-ai-memory-method-loci-2026)
- Mem0 (formerly EmbedChain) has open-source self-hostable + managed cloud; Letta from UC Berkeley MemGPT research; Zep targets enterprise compliance. — [Powerdrill](https://powerdrill.ai/blog/best-free-memory-layer-for-ai-agents); [Developers Digest](https://www.developersdigest.tech/blog/best-ai-agent-memory-providers-2026)
- Common naming patterns observed: "mem-" prefix (Mem0, MemGPT, MemPalace, Memori), "brain" (Obsidian Second Brain), cognitive metaphors (Letta "sleep-time", Karpathy "Memex"). — [Edenai overview](https://www.edenai.co/post/ai-agent-memory-mempalace-mem0-and-persistent-context); [Memori](https://memorilabs.ai/blog/how-to-build-an-openai-agent-with-persistent-memory)

### Inferences
- "mindpalace" will be read as adjacent to/derivative of MemPalace (same method-of-loci metaphor, same domain, viral in 2026). Positioning should differentiate explicitly: MemPalace = verbatim conversation storage + vector search runtime; a bootstrapper = scaffolding/organization of human-readable markdown projects + memory conventions, agent-agnostic (AGENTS.md/SKILL.md/MCP). Consider a name check on GitHub/PyPI/npm and trademark before committing.
- A tagline formula that uses searched terms: "Scaffold an LLM wiki / AI second brain for any agent — context engineering in plain markdown."
- Clear category leadership for "bootstrapper/scaffolder of agent memory folders" was not found — this appears to be an open niche between conventions (Cline Memory Bank, LLM Wiki skill) and runtimes (Mem0/Letta/MemPalace).

### Gaps
- No direct competitor search was done for "scaffold"/"init" CLIs for agent memory (e.g., `npx create-...` style); a dedicated competitive scan is recommended.
- GitHub star counts are from aggregators, not verified on GitHub.
