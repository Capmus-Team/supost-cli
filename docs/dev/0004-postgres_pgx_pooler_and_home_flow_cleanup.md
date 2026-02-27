# Postgres Adapter Upgrade + Home Flow Cleanup

Date: 2026-02-27

## Summary
This increment replaces the `psql` subprocess strategy with a native Go PostgreSQL adapter using `database/sql` + `pgx`, improves Supabase pooler compatibility via DSN normalization, and simplifies `supost home` execution flow by centralizing output rendering and tightening cache usage.

## What Changed

### 1. Postgres repository migrated to native DB driver
- `internal/repository/postgres.go` now uses `database/sql` with `github.com/jackc/pgx/v5/stdlib`.
- Added a real connection pool (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`).
- `Close()` now closes the DB pool (previously a no-op when shelling out to `psql`).
- `ListRecentActivePosts` now executes a parameterized SQL query with `$1`, `$2` and scans directly into `domain.Post`.

### 2. Supabase pooler-safe connection string normalization
- Added `ensurePoolerSafeConnectionString(databaseURL string)` in `internal/repository/postgres.go`.
- Ensures these settings are present for both URL and key/value DSN styles:
  - `default_query_exec_mode=simple_protocol`
  - `statement_cache_capacity=0`
  - `description_cache_capacity=0`
- Goal: avoid prepared statement/cache behavior that can conflict with transaction poolers.

### 3. Home command flow simplification
- `cmd/home.go` now attempts cache reads only when `DATABASE_URL` is configured.
- Extracted format dispatch into `renderHomeOutput(...)` to remove duplicated branching.
- Fetch path is now straightforward:
  - try cache
  - fetch fresh from service when needed
  - save cache (DB mode + TTL > 0)
  - render via one helper

### 4. Tests updated for new behavior
- `internal/repository/postgres_test.go` now tests:
  - limit clamping behavior (`clampRecentLimit`)
  - DSN normalization for URL and key/value styles (`ensurePoolerSafeConnectionString`)
- Removed legacy parser-specific test coverage tied to the old delimiter-based `psql` parsing path.

### 5. Module dependencies refreshed
- `go.mod` and `go.sum` now include `pgx/v5`-related dependencies used by the new adapter implementation.

## Why This Matters
- Removes runtime dependency on `psql` CLI and shell parsing.
- Improves correctness and maintainability through typed row scanning.
- Keeps SQL injection-safe behavior with parameterized queries.
- Improves compatibility with Supabase poolers in production-like configurations.
- Makes `supost home` command behavior easier to read and reason about.

## Files in This Increment
- `cmd/home.go`
- `internal/repository/postgres.go`
- `internal/repository/postgres_test.go`
- `go.mod`
- `go.sum`
- `docs/dev/0004-postgres_pgx_pooler_and_home_flow_cleanup.md`
