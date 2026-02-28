# Post Create Submit Flow With Mailgun and Persistence

Date: 2026-02-28

## Summary
Extended `post create` from staged form rendering to a full submit workflow: validates input, creates a pending post in repository storage, generates publish-link email content, and sends via Mailgun (or dry-run preview).

## What Changed

### 1. `post create` command now supports submit-mode flags
- Updated `cmd/post_create.go`:
  - added submit flags:
    - `--name`
    - `--body`
    - `--email`
    - `--price`
    - `--dry-run`
  - detects submit mode when submit-related flags are set
  - builds `PostCreateSubmission`, calls `PostCreateService.Submit(...)`
  - configures optional Mailgun sender when not dry-run
  - adds submit output rendering path (`RenderPostCreateSubmitResult`).

### 2. Config/env expanded for publish-link email flow
- Updated `.env.example` and `configs/config.yaml.example` with:
  - `SUPOST_BASE_URL`
  - `MAILGUN_API_BASE`
  - `MAILGUN_SEND_TIMEOUT`
  - plus existing Mailgun domain/key/from values.
- Updated `internal/config/config.go`:
  - added config fields for Mailgun + base URL
  - wired Viper reads including duration parsing for send timeout.

### 3. New post-create submit domain models
- Added `internal/domain/post_create_submit.go`:
  - `PostCreateSubmission`
  - `PostCreatePersisted`
  - `PublishEmailMessage`
  - `PostCreateSubmitResult`.
- Updated `internal/domain/post.go` with additional DB-mapped fields used by submit/persist flows:
  - `time_modified`, `time_modified_at`, `access_token`.

### 4. PostCreate service now includes submit use-case
- Updated `internal/service/post_create.go` repository interface:
  - added `CreatePendingPost(...)`.
- Added `internal/service/post_create_submit.go`:
  - validates required fields and category/subcategory combination
  - validates Stanford email domains
  - enforces non-negative required price
  - generates access token
  - builds publish URL and email body/subject
  - supports dry-run (no insert/send)
  - persists pending post and sends email in non-dry-run mode.
- Added `internal/service/post_create_submit_test.go`:
  - dry-run behavior
  - persist+send path
  - invalid email rejection
  - Stanford subdomain acceptance.
- Updated `internal/service/post_create_test.go` mock to satisfy expanded repo interface.

### 5. Repository implementations for pending post insert
- Added `internal/repository/inmemory_post_create.go`:
  - appends pending post into in-memory store with generated ID and timestamps.
- Added `internal/repository/postgres_post_create.go`:
  - parameterized insert into `public.post`
  - returns inserted `id`, `access_token`, and posted time.

### 6. New Mailgun adapter + submit output adapter
- Added `internal/adapters/mailgun.go` + `mailgun_test.go`:
  - Mailgun sender with API base/domain/key/from config
  - form-encoded send with basic auth and non-2xx handling.
- Added `internal/adapters/post_create_submit_output.go` + tests:
  - renders dry-run/submit summary including publish URL, recipient, subject, and body.

## Why This Matters
- Moves post creation from static staged UI into an end-to-end publish-link workflow.
- Keeps separation of concerns clear: command wiring, service validation/orchestration, repository persistence, and adapter email/output side effects.
- Supports safe local iteration via `--dry-run` while preserving production-like behavior with Mailgun.

## Files in This Increment
- `.env.example`
- `configs/config.yaml.example`
- `internal/config/config.go`
- `cmd/post_create.go`
- `internal/domain/post.go`
- `internal/domain/post_create_submit.go`
- `internal/service/post_create.go`
- `internal/service/post_create_test.go`
- `internal/service/post_create_submit.go`
- `internal/service/post_create_submit_test.go`
- `internal/repository/inmemory_post_create.go`
- `internal/repository/postgres_post_create.go`
- `internal/adapters/mailgun.go`
- `internal/adapters/mailgun_test.go`
- `internal/adapters/post_create_submit_output.go`
- `internal/adapters/post_create_submit_output_test.go`
- `docs/dev/0026-post_create_submit_flow_with_mailgun_and_persistence.md`
