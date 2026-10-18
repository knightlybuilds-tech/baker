## Stack
Laravel{{with .Facts.laravel_version}} {{.}}{{end}}{{with .Facts.php_version}} on PHP {{.}}{{end}}. Tests: {{.Facts.test_runner}}.

## Commands
- Dev server (app, queue, logs, Vite): `composer run dev`
- Tests: `php artisan test` (single file: `php artisan test tests/Feature/ExampleTest.php`)
- Format: `./vendor/bin/pint` (run before committing)
- Migrations: `php artisan migrate` · fresh with seed: `php artisan migrate:fresh --seed`
- Generate code with artisan (`php artisan make:model Post -mfc`) instead of writing boilerplate by hand.

## Conventions
- Follow Laravel defaults before inventing structure. Check `php artisan route:list` and existing code first.
- Validation in Form Requests, not controllers. Authorization in Policies.
- Every migration has a working `down()`. Never edit a migration that has run in production; add a new one.
- Database access through Eloquent; avoid raw queries unless measured.
- Every bug fix or feature comes with a {{.Facts.test_runner}} test.
{{- if not .Facts.boost}}

## Tooling tip
Laravel Boost adds an MCP server that lets agents read routes, schema and logs:
`composer require laravel/boost --dev && php artisan boost:install`, then run `baker sync` to register it for every agent.
{{- end}}
