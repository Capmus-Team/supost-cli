# Post Respond Flow With Reply Email and Message Persistence

Date: 2026-02-28

## Summary
Added a full `post respond` workflow to send response emails to post owners and persist response messages, including Mailgun `Reply-To` support, new domain models, repository persistence methods, and CLI rendering.

## What Changed

### 1. New `post respond` command
- Added `cmd/post_respond.go` under `post`:
  - usage: `post respond <post_id>`
  - required flags:
    - `--message`
    - `--reply-to`
  - optional `--dry-run` mode
  - wires repository + optional Mailgun sender
  - calls `PostRespondService.Respond(...)`
  - renders structured output via adapter.

### 2. New response/message domain models
- Added `internal/domain/post_respond.go`:
  - `PostRespondSubmission`
  - `ResponseEmailMessage`
  - `PostRespondResult`.
- Added `internal/domain/message.go`:
  - `Message` mapping for `app_private.message` records.

### 3. New post-respond service and tests
- Added `internal/service/post_respond.go`:
  - validates submission fields
  - loads target post and requires destination email + access token
  - builds response email subject/body with safety/contact lines
  - supports dry-run (no send/no save)
  - sends response email and persists message in non-dry-run mode.
- Added `internal/service/post_respond_test.go`:
  - dry-run behavior
  - send+persist happy path
  - validation error coverage.

### 4. Repository support for response message persistence
- Added `internal/repository/inmemory_post_respond.go`:
  - stores response messages in memory with incrementing message IDs.
- Added `internal/repository/postgres_post_respond.go`:
  - parameterized insert into `app_private.message`
  - returns inserted message row fields.
- Updated `internal/repository/inmemory.go`:
  - added `messages` store
  - expanded seed posts with access tokens/time-modified values and added a realistic sample response target post.
- Updated `internal/repository/inmemory_post_create.go`:
  - persisted `AccessToken` onto newly created in-memory posts.
- Updated `internal/repository/postgres.go`:
  - expanded post selects/scans to include `time_modified`, `time_modified_at`, and `access_token`.

### 5. Mailgun adapter expanded for response emails
- Updated `internal/adapters/mailgun.go`:
  - added `SendResponseEmail(...)`
  - refactored send path through shared `sendTextEmail(...)`
  - supports optional `h:Reply-To` header.
- Updated `internal/adapters/mailgun_test.go`:
  - added `Reply-To` payload assertion test.

### 6. New post-respond output adapter and test
- Added `internal/adapters/post_respond_output.go`:
  - renders send/dry-run summary and full email preview content.
- Added `internal/adapters/post_respond_output_test.go`:
  - verifies expected summary fields and subject in output.

## Why This Matters
- Completes an end-to-end response flow for posts with clear dry-run and send modes.
- Preserves clean layering: command wiring, service logic, repository persistence, adapter side effects.
- Improves operational realism by capturing reply metadata in `app_private.message` and using `Reply-To` in outgoing mail.

## Files in This Increment
- `cmd/post_respond.go`
- `internal/domain/message.go`
- `internal/domain/post_respond.go`
- `internal/service/post_respond.go`
- `internal/service/post_respond_test.go`
- `internal/repository/inmemory_post_respond.go`
- `internal/repository/postgres_post_respond.go`
- `internal/repository/inmemory.go`
- `internal/repository/inmemory_post_create.go`
- `internal/repository/postgres.go`
- `internal/adapters/mailgun.go`
- `internal/adapters/mailgun_test.go`
- `internal/adapters/post_respond_output.go`
- `internal/adapters/post_respond_output_test.go`
- `docs/dev/0027-post_respond_flow_with_reply_email_and_message_persistence.md`
