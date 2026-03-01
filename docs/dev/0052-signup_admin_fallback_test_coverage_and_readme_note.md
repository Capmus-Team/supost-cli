# Signup Admin Fallback Test Coverage and README Note

Date: 2026-03-01

## Summary
Documented and validated the Supabase Auth signup admin-fallback behavior when email send rate limits occur by adding focused adapter test coverage and a README usage note.

## What Changed

### 1. Added adapter test for signup admin fallback
- Updated `internal/adapters/supabase_auth_signup_test.go`:
  - added `TestSupabaseAuthSignupClient_SignUp_FallsBackToAdminCreateUserOnEmailRateLimit`
  - verifies public signup path returns 429 with `over_email_send_rate_limit`
  - verifies client retries with `/auth/v1/admin/users` when a secret key is configured
  - verifies auth headers use publishable key for public path and secret key for admin path
  - verifies admin payload fields (`email`, `password`, `email_confirm`, `user_metadata`)
  - verifies result mapping from admin response and `email_confirmation_sent=false` on fallback.

### 2. Added README note for optional secret-key behavior
- Updated `README.md`:
  - added note beneath signup usage explaining optional admin fallback behavior when `SUPABASE_SECRET_KEY` (or `SUPABASE_SERVICE_ROLE_KEY`) is configured and email rate limits are hit.

## Why This Matters
- Locks in expected behavior for signup resilience under Supabase email rate limiting.
- Makes fallback behavior explicit for operators configuring CLI environment variables.

## Files in This Increment
- `internal/adapters/supabase_auth_signup_test.go`
- `README.md`
- `docs/dev/0052-signup_admin_fallback_test_coverage_and_readme_note.md`
