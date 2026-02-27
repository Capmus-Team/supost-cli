# Home Feed: Performance Optimizations

Date: 2026-02-27

## Summary
Improved `supost home` speed while preserving:
- plain-text terminal rendering by default
- JSON output support for future web/TypeScript consumption (`--format json`)

## What Changed

### 1. Faster Postgres row parsing
- Kept `psql` execution path, but removed per-row JSON encode/decode overhead.
- Switched to delimiter-based row output (`-F` and `-R`) and lightweight parsing.
- Reduced parse cost for the home feed path.

File:
- `internal/repository/postgres.go`

### 2. Smaller query payload
- Home query now selects only fields needed for home rendering and JSON home output:
  - `id`, `email`, `name`, `status`, `time_posted`, `time_posted_at`, `price`, `has_price`, `has_image`
- Removed non-essential columns from this endpoint path to reduce transfer + parsing work.

File:
- `internal/repository/postgres.go`

### 3. Short-lived local cache for repeated runs
- Added local cache for DB-backed home feed calls.
- Default TTL: `30s`.
- New flags:
  - `--cache-ttl` to tune/disable cache (`0s` disables)
  - `--refresh` to bypass cache and force fresh fetch
- Goal: make repeated `supost home` runs near-instant during active usage.

Files:
- `cmd/home.go`
- `internal/adapters/home_cache.go`

### 4. Parser/limit tests
- Added repository tests for:
  - limit clamping
  - row parser correctness

File:
- `internal/repository/postgres_test.go`

## Command Behavior After Optimization

### Default terminal output
```bash
supost home
```

### JSON output (for frontend-aligned contract checks)
```bash
supost home --format json
```

### Force fresh fetch (skip cache)
```bash
supost home --refresh
```

### Custom cache duration
```bash
supost home --cache-ttl 60s
```

## Validation
- `go test ./... -race` passes with the updated code.

## Files in This Increment
- `cmd/home.go`
- `internal/repository/postgres.go`
- `internal/adapters/home_cache.go`
- `internal/repository/postgres_test.go`
- `docs/dev/0003-home_performance_optimizations.md`
