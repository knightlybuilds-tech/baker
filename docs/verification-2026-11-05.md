# Initial verification on 5 November 2026

The first reviewed-memory slice was exercised with real authenticated CLI
sessions, synthetic repositories, and assertions outside those repositories.
The working tree is based on Baker commit `51ddba0`; these are development
results, not a released compatibility guarantee.

| Check | Claude to Codex | Codex to Claude |
| --- | --- | --- |
| Source produced an unapproved crumb | Passed | Passed |
| Proposal absent from approved recall | Passed | Passed |
| Fresh teammate clone excluded personal handoff and draft | Passed | Passed |
| Target visibly invoked reviewed recall | Passed | Passed |
| Target implementation met external exact-format assertions | Passed | Passed |
| Baker handoff produced a working continuation | Passed with local handoff | Passed without hot.md or an approved crumb |
| Changed evidence withheld the approved body | Passed | Passed |
| Protected tracked files unchanged by team target | Not instrumented in first run | Passed |

The second run used Claude Code `2.1.316` and Codex CLI `0.166.0`. In its recovery
phase, the harness temporarily removed the local handoff, source draft, review
record, and approved note. Baker recovered the Codex transcript; Claude completed
the function from the generated brief. Removed fixture review files were then
restored. This simulates missing closing state, not a forced process termination
or an actual usage-limit event.

Both runs verified semicolon separation, always-quoted fields, CRLF termination,
CSV quote escaping, empty input, and preservation of case and leading zeros.
Neither run measured a baseline. Timings are retained in the local reports but
should not be interpreted as a performance comparison.

Approval in these fixtures was a **synthetic harness review**, not human or
authenticated team approval. Production review authorization remains future work.

Additional checks passed: `go test -race ./...`, `make lint`, `make build`, and
`git diff --check`. Regression coverage also verifies that rapid checkpoints
remain distinct and checkpoint fallback respects `--from`.

The evaluation runner, command contract, and rerun instructions are documented in
[Reviewed project memory](reviewed-memory.md). Local run directories retain the
JSON reports, synthetic projects, prompts, and agent output for inspection. Real
agent evaluations are intentionally excluded from ordinary CI and consume quota
only when explicitly invoked.
