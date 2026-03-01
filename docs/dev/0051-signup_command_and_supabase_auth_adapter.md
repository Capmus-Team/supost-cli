# Signup Command and Supabase Auth Adapter

Date: 2026-03-01

## Summary
Added a new `signup` top-level command that validates user signup inputs in the service layer and calls Supabase Auth via a dedicated adapter. Also updated config and command-reference coverage to include signup flags and Supabase key fields.

## What Changed

### 1. Added signup command wiring
- Added `cmd/signup.go`:
  - introduces `supost signup`
  - requires `--display-name`, `--email`, `--phone`, and `--password`
  - loads config and picks `supabase_publishable_key` (fallback to `supabase_anon_key`)
  - constructs signup adapter/service and renders JSON/table output via existing renderer.

### 2. Added signup domain contract
- Added `internal/domain/user_signup.go`:
  - `UserSignupSubmission` for command/service input
  - `UserSignupResult` for output payload
  - uses `json` + `db` tags aligned with project conventions.

### 3. Added signup service validation/orchestration
- Added `internal/service/user_signup.go`:
  - defines a consumed `UserSignupProvider` interface
  - trims/normalizes inputs (including lowercased email)
  - validates display name/email/phone/password and returns combined validation errors
  - delegates actual account creation to provider.

### 4. Added Supabase Auth adapter
- Added `internal/adapters/supabase_auth_signup.go`:
  - posts to `/auth/v1/signup` with apikey + bearer headers
  - sends display name/phone as `data` metadata
  - falls back to `/auth/v1/admin/users` using `supabase_secret_key` when Auth email-send rate limit is hit
  - handles both nested `user` and top-level auth response shapes
  - returns structured signup result and detailed non-2xx/no-user errors.

### 5. Added tests for new command and signup flow
- Added `internal/service/user_signup_test.go`:
  - success path asserts normalization and provider invocation
  - validation path asserts bad submissions are rejected.
- Added `internal/adapters/supabase_auth_signup_test.go`:
  - verifies request method/path/headers/payload
  - covers non-2xx API failure handling
  - covers top-level user response variant
  - covers no-user response with detailed message.
- Updated `cmd/command_reference_test.go`:
  - asserts `signup` command exists
  - asserts required signup flags
  - includes new signup files in structure contract checks.

### 6. Updated docs and configuration examples
- Updated `README.md`:
  - added `supost signup` usage examples and command tree entries
  - added signup-related files to project map
  - added `SUPABASE_PUBLISHABLE_KEY` in env example block.
- Updated `configs/config.yaml.example`:
  - added `supabase_publishable_key`
  - added `supabase_secret_key`.
- Updated `internal/config/config.go`:
  - added config fields + load bindings for publishable/secret Supabase keys.

## Why This Matters
- Provides an end-to-end CLI signup workflow backed by Supabase Auth.
- Keeps validation and orchestration in `internal/service` while isolating HTTP side effects in adapters.
- Maintains command/structure contract coverage so the new surface stays discoverable and stable.

## Files in This Increment
- `cmd/signup.go`
- `cmd/command_reference_test.go`
- `internal/domain/user_signup.go`
- `internal/service/user_signup.go`
- `internal/service/user_signup_test.go`
- `internal/adapters/supabase_auth_signup.go`
- `internal/adapters/supabase_auth_signup_test.go`
- `internal/config/config.go`
- `configs/config.yaml.example`
- `README.md`
- `docs/dev/0051-signup_command_and_supabase_auth_adapter.md`
