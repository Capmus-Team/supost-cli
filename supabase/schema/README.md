# Schema Snapshot Artifacts

This folder contains a remote Supabase schema snapshot plus optimization-focused inventories.

- `full_schema.sql`
  - Full schema-only dump from remote Postgres.
  - Includes tables, views, functions, constraints, indexes, policies, and triggers.
- `indexes.csv`
  - All non-system index definitions (`pg_indexes`).
- `triggers.csv`
  - All non-system trigger definitions (`information_schema.triggers`).
- `table_sizes.csv`
  - Per-table size and row-estimate metrics for prioritizing tuning.
- `index_usage_stats.csv`
  - Per-index usage counters from `pg_stat_user_indexes`.
- `table_scan_stats.csv`
  - Per-table sequential/index scan stats from `pg_stat_user_tables`.

## Refresh Command Pattern

Use a Postgres client version matching the server major version (currently 17):

```bash
/opt/homebrew/opt/postgresql@17/bin/pg_dump "$DATABASE_URL" --schema-only --file supabase/schema/full_schema.sql
```

Then refresh the CSV inventories with the same `psql` `\copy` queries used in this session.
