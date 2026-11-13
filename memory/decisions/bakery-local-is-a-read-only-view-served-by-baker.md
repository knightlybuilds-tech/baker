---
type: decision
scope: project
tags: [bakery, ui]
as_of: 2026-11-13
last_verified: 2026-11-13
source: Claude Code session 2026-11-13; internal/bakery, docs/bakery.md
confidence: high
---
# Local Bakery is a read-only view served by baker

`baker bakery` embeds the Bakery UI and serves it on 127.0.0.1 with one read-only endpoint, `/api/memory`, built from the same parsing as `baker memory list` and `check`. It never writes notes, and it never serves `hot.md`, logs, checkpoints or transcripts. That keeps the files as the source of truth and needs no backend, so teammates can try Bakery on their own repos before a hosted service exists. Review, branch and team screens stay marked "soon" until a real service with authentication backs them. See [[agents-md-is-the-single-source]] for why files stay canonical.
