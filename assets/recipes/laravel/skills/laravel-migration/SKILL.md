---
name: laravel-migration
description: Create or change database tables in this Laravel app safely. Use for any schema change — new table, new column, index, rename.
---

# Laravel migration

1. Check the current schema first: read existing migrations in `database/migrations/` for the table.
2. Create a new migration with artisan — never edit one that has already run outside your machine:
   `php artisan make:migration add_status_to_orders_table`
3. Write both `up()` and a real `down()` that reverses it exactly.
4. Large tables: add columns as nullable or with a default; build indexes in a separate migration.
5. Update the model (`$fillable`/`$casts`), factory and any Form Request rules that touch the column.
6. Run `php artisan migrate`, then `php artisan migrate:rollback` and `php artisan migrate` again to prove `down()` works.
7. Run the test suite.
