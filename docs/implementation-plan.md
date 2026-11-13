# Standalone Baker implementation direction

## Current scope

Build Baker as a complete standalone CLI first. It must capture, retain,
retrieve, and transfer project context without Bakery, an account, or a network
service. An individual switching between Claude and Codex is a complete use
case in its own right.

Keep the product team-first through its foundations: project knowledge is
portable, readable, and suitable for a shared Git repository. Personal task
state stays local. Bakery will add team collaboration around that foundation;
it is not a runtime dependency or a prerequisite for useful memory.

- Local capture, notes, recall, checkpoints, and handoffs belong to Baker.
- Team sharing, authenticated review, access control, and hosted coordination
  belong to future Bakery integration.
- Local note retention and retrieval do not require proposal or approval.
- Public publishing and a marketplace remain outside the present scope.
- The card and graph UI is a reference, not a backend feature specification.

## Memory responsibilities

| Layer | Purpose | Default storage |
| --- | --- | --- |
| Durable project notes | Decisions, lessons, domain knowledge and sources | Committed Markdown in memory/ |
| Current task | Goal, progress, next action, unresolved questions | Local memory/hot.md |
| Checkpoints | Mechanical recovery state when a session stops | Local .baker/checkpoints/ |
| Native transcripts | Additional evidence after interruption | Existing agent CLI stores |

Use `baker memory list`, `recall`, and `check` for ordinary local notes. Reading
the index and Markdown remains a valid fallback when Baker is unavailable.
Freshness warnings are reasons to verify, not a service approval gate.

Preserve room for stable IDs, revisions, provenance, explicit relationships,
and optional synchronization. Add fields when an actual workflow requires them;
do not impose a complete future server protocol on every local note now.

## Existing work and its status

The imported repository starts at `51ddba0`. It already supplies scaffolding,
managed agent bridges, hooks, atomic notes, checkpoints, transcript readers,
handoff and resume. Initial real-agent runs exercised transfers in both
directions and a transcript recovery case.

The prematurely introduced review lifecycle is retained under `baker
experimental crumb` for development only. It is not the default retention path
and does not define Bakery's future authorization model. Its earlier successful
evaluations do not establish that standalone note retention is fully tested.

## Build and verification order

1. Verify setup and safe repeated sync with installed Claude and Codex.
2. Verify durable note capture and retrieval across fresh agent sessions using
   ordinary memory, with no experimental review command or Bakery service.
3. Verify each agent's automatic checkpoint hook independently.
4. Verify recovery from hot.md, checkpoint only, and transcript only as separate
   scenarios, then force a real mid-task process interruption.
5. Verify scope isolation, stale notes, parallel branches/worktrees, and missing
   or incompatible integrations.
6. Repeat tests against manual handoffs and ordinary repository documentation;
   record failures, corrections, maintenance effort, and task outcomes.
7. Introduce optional Bakery integration once the standalone contract is useful
   and the team workflows justify the service boundary.

Use isolated fixture projects plus an actively developed real project. Keep
deterministic tests in CI and quota-consuming real-agent tests opt-in. Each real
failure becomes a reproducible fixture and a targeted regression check.
