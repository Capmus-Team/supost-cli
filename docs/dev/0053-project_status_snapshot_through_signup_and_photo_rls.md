# Project Status Snapshot Through Signup and Photo RLS

Date: 2026-03-02

## Summary
This increment records the current state of `supost-cli` after the signup flow work and the new `photo` public-read RLS policy migration.

## What Has Been Built So Far

### 1. Core CLI and architecture baseline
- Cobra-based command surface wired through `cmd/` with thin command handlers.
- Layered `internal/` structure for domain, service, repository, adapters, and config.
- JSON-first output rendering and command-level formatting support.

### 2. Supabase-first schema and migration workflow
- `supabase/migrations/` is used as schema source of truth.
- Progressive migration history captures taxonomy integrity, message constraints, FTS/search indexes, photo table rollout, and backfill support.
- `supabase/schema/*` artifacts are present for schema/index/scan visibility.

### 3. Home/search/post/category feature set
- Home and category-driven listing surfaces with service/repository wiring.
- Search command supports optional keyword query and Postgres FTS path.
- Post detail and reply flow with message persistence and validation coverage.

### 4. Post create + photo upload evolution
- Staged post-create flow with rules validation and persistence behavior.
- Photo support includes public `photo` table, ticker image key support, incremental backfill scripts, and S3 upload adapter/tests.
- Additional validation and image-processing regression coverage added in recent increments.

### 5. Signup flow with Supabase Auth adapter
- Signup command and service wiring added.
- Supabase auth adapter supports primary signup and documented admin fallback behavior when email send rate limits occur.
- Focused adapter tests cover fallback semantics and payload/header expectations.

### 6. Latest uncommitted schema change captured in this commit
- Added migration: `supabase/migrations/20260301013000_photo_public_read_rls.sql`.
- Enables RLS on `public.photo` and creates `photo_public_read` policy permitting `anon`/`authenticated` reads when a parent `public.post` exists.
- Removed legacy placeholder file: `migrations/.gitkeep`.

## Current Working Baseline
- CLI commands for home/search/post/post-create/post-respond/categories/signup are implemented with service-backed flows.
- In-memory and Postgres repository paths exist for key features.
- Test coverage includes service and adapter regressions across core flows.
- Development docs in `docs/dev/` now provide a continuous implementation log through this state.

## Files in This Increment
- `docs/dev/0053-project_status_snapshot_through_signup_and_photo_rls.md`
- `supabase/migrations/20260301013000_photo_public_read_rls.sql`
- `migrations/.gitkeep` (removed)
