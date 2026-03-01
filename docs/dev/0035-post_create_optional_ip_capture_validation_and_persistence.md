# Post Create Optional IP Capture, Validation, and Persistence

Date: 2026-03-01

## Summary
Extended `post create` submit flow with an optional `--ip` field, validated as IPv4/IPv6, and propagated it through domain models, service validation, repository persistence, tests, and README command docs.

## What Changed

### 1. Added optional `--ip` flag to post create command
- Updated `cmd/post_create.go`:
  - reads `--ip` from flags
  - trims and forwards it in `domain.PostCreateSubmission`
  - added flag help: `poster IP address (optional)`.

### 2. Extended domain types with IP field
- Updated `internal/domain/post_create_submit.go`:
  - added `IP string` with `json:"ip"` and `db:"ip"` tags to `PostCreateSubmission`.
- Updated `internal/domain/post.go`:
  - added `IP string` with `json:"ip"` and `db:"ip"` tags to `Post`.

### 3. Added service-level IP validation
- Updated `internal/service/post_create_submit.go`:
  - normalizes `IP` input via `strings.TrimSpace`
  - validates optional IP using `net/netip.ParseAddr`
  - returns validation problem: `IP must be a valid IPv4 or IPv6 address.` when invalid.

### 4. Persisted IP in in-memory and Postgres repositories
- Updated `internal/repository/inmemory_post_create.go`:
  - copies `submission.IP` onto persisted in-memory post record.
- Updated `internal/repository/postgres_post_create.go`:
  - insert now includes `ip` column/value
  - stores `NULL` when IP is empty and provided IP string when present.

### 5. Updated tests and command reference expectations
- Updated `internal/service/post_create_submit_test.go`:
  - verifies valid IP is forwarded to repository
  - adds regression test rejecting invalid IP.
- Updated `cmd/command_reference_test.go`:
  - includes `ip` in expected `post create` flags
  - updates expected directory contract list to keep `supabase/migrations`, `testdata/seed`, and `docs`.

### 6. Updated README examples/command tree
- Updated `README.md`:
  - `post create` examples now include `--ip` (IPv4 + IPv6 example)
  - command tree documents optional `--ip <address>` flag.

## Why This Matters
- Captures optional poster network metadata for moderation/audit workflows.
- Ensures user input quality with explicit IP format validation.
- Keeps CLI docs and command/test contracts aligned with the new submit payload shape.

## Files in This Increment
- `README.md`
- `cmd/command_reference_test.go`
- `cmd/post_create.go`
- `internal/domain/post.go`
- `internal/domain/post_create_submit.go`
- `internal/repository/inmemory_post_create.go`
- `internal/repository/postgres_post_create.go`
- `internal/service/post_create_submit.go`
- `internal/service/post_create_submit_test.go`
- `docs/dev/0035-post_create_optional_ip_capture_validation_and_persistence.md`
