# Home Performance Docs + Compact Age Edge-Case Test

Date: 2026-02-27

## Summary
Updated user-facing documentation for home performance/caching behavior and added adapter test coverage for zero-time compact-age formatting.

## What Changed

### 1. README performance notes for home command
- Updated `README.md` with a `Home Performance Notes` section covering:
  - silent success behavior of `go build -o bin/supost .`
  - recommendation to use `./bin/supost home` for repeated runs
  - current home cache layers (posts + category latest timestamps)
  - taxonomy source (local seed data instead of runtime category-table queries).

### 2. Terminal commands doc expanded
- Updated `docs/terminal_commands.md`:
  - clarified that no output from `go build -o bin/supost .` means success
  - added explicit home performance path notes:
    - cache for recent active posts
    - cache for category last-active timestamps
    - taxonomy loaded from local seed data.

### 3. Adapter test for compact age zero-time behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestFormatCompactAge_ZeroTime`
- Verifies empty/zero timestamp is rendered as `no active posts`.

## Why This Matters
- Makes expected build/runtime behavior clearer for contributors using local binaries.
- Documents current home rendering data flow and caching strategy.
- Adds edge-case regression coverage for category age formatting.

## Files in This Increment
- `README.md`
- `docs/terminal_commands.md`
- `internal/adapters/home_output_test.go`
- `docs/dev/0014-home_performance_docs_and_compact_age_test.md`
