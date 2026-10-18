---
type: decision
scope: project
tags: [cli, distribution]
as_of: 2026-10-18
last_verified: 2026-10-18
source: design session 2026-10-18
confidence: high
---
# baker is a single Go binary

Chosen over TypeScript/Node so it installs the same way for PHP, Python and JS developers with no runtime, starts instantly (it runs in agent hooks after every turn), and cross-compiles easily. Recipes and skills are embedded with `embed.FS`.
