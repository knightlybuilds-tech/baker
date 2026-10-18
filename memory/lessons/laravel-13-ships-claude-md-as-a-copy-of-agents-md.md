---
type: lesson
scope: project
tags: [laravel, recipes]
as_of: 2026-10-18
last_verified: 2026-10-18
source: design session 2026-10-18
confidence: high
---
# Laravel 13 ships CLAUDE.md as a copy of AGENTS.md

`composer create-project laravel/laravel` (Laravel 13) creates both files with identical Laravel Boost guidelines. Once baker's bridge imports AGENTS.md, the copy would load twice, so `planBridges` drops bridge-file text that is already contained in AGENTS.md.
