## Memory
This project keeps its memory in `memory/` as small markdown notes. Any agent can read and write it.

- **Session start:** read `memory/hot.md` and `memory/MEMORY.md`. Open other memory notes only when relevant to the task.
- **Continuing work:** if `hot.md` has unchecked items under Next, you are picking up someone else's task (another agent, account, or provider). Say the plan back to the user in one line, then continue.
- **Checkpoint as you go:** after each finished step, before any long or risky operation, and whenever the user changes direction, update `memory/hot.md` (Goal, Done, Next, Gotchas, Open questions). Your session can end at any moment (usage limits) — never keep progress only in your head.
- **Session end:** follow the `session-close` skill: append to `memory/log/<YYYY-MM-DD>.md`, update `hot.md`, propose lessons.
- **Atomic notes:** one idea per file in `memory/decisions/`, `memory/lessons/` or `memory/knowledge/`. The title is the claim as a full sentence, the filename is its slug, under 200 words. Link related notes with `[[slug]]`. Add one line per note to `memory/MEMORY.md`.
- **Freshness:** every note has `type`, `as_of`, `last_verified`, `source` frontmatter. Treat a note not verified in 90 days as a hint, not a fact; re-verify before relying on it.
- **Promotion:** you may propose turning a lesson into a rule in this file, but only the user approves it.
- **Never** write secrets, tokens, passwords or customer data into memory.
