# Home Feed: Recent Active Posts

Date: 2026-02-27

## Summary
Implemented `supost home` to render the latest active posts from `public.post`, capped at 50 items, in a terminal style that mirrors the SUPost homepage screenshot.

## What Was Implemented

### 1. New `home` command
- Added `supost home` with `--limit` (default `50`)
- Adapter selection:
  - Uses Postgres path when `DATABASE_URL` is set
  - Falls back to in-memory seed data when `DATABASE_URL` is empty
- Rendering behavior:
  - Defaults to homepage terminal view (screenshot-like) when no explicit `--format` is set
  - Supports JSON output when `--format json` is explicitly requested

File:
- `cmd/home.go`

### 2. Domain model for posts
- Added `internal/domain/post.go` with `json` + `db` tags
- Includes fields required for home rendering and post ordering:
  - `id`, `category_id`, `subcategory_id`, `email`, `name`, `status`
  - `time_posted`, `time_posted_at`
  - `price`, `has_price`, `has_image`
  - `created_at`, `updated_at`
- Added constant:
  - `PostStatusActive = 1`

File:
- `internal/domain/post.go`

### 3. Home service
- Added a dedicated home service and interface:
  - `ListRecentActive(ctx, limit int)`
  - Default limit behavior when `limit <= 0`
- Added tests validating default and explicit limit behavior

Files:
- `internal/service/home.go`
- `internal/service/home_test.go`

### 4. Repository updates
- Added read contract for home feed:
  - `ListRecentActivePosts(ctx, limit int)`
- In-memory adapter:
  - Added `posts` seed data
  - Filters to `status == 1`
  - Sorts by `time_posted DESC, id DESC`
  - Applies limit
- Postgres adapter:
  - Added `internal/repository/postgres.go`
  - Uses `psql` with `DATABASE_URL`
  - Query targets `public.post`, `status = 1`, order by newest, returns up to 50 rows
  - Computes `has_image` from image/photo columns
  - Parses JSON row output into domain model

Files:
- `internal/repository/interfaces.go`
- `internal/repository/inmemory.go`
- `internal/repository/postgres.go`

### 5. Home terminal renderer
- Added dedicated renderer for screenshot-like output:
  - Header: `recently posted`
  - Blue title
  - Gray email
  - Optional camera icon when image exists
  - Magenta relative time (`about X hours`)
  - Price formatting (`$2,000`, `Free`)

File:
- `internal/adapters/home_output.go`

### 6. Root config loading improvement
- `cmd/root.go` now loads `.env` via `gotenv`
- Binds persistent flags to Viper so global flags consistently reach config

File:
- `cmd/root.go`

## Query Behavior (Home Feed)

From `public.post`:
- `WHERE status = 1`
- `ORDER BY time_posted DESC NULLS LAST, id DESC`
- `LIMIT 50` (with command-level down-limit via `--limit`)

## Validation Performed
- `go test ./... -race` passed after implementation.
- `go run . home --format text` verified screenshot-like rendering in in-memory mode.
- Live DB path is wired and uses `DATABASE_URL`; successful execution depends on host/network DNS resolution from the local environment.

## Changed Files (Feature Scope)
- `.gitignore`
- `cmd/root.go`
- `cmd/home.go`
- `internal/domain/post.go`
- `internal/service/home.go`
- `internal/service/home_test.go`
- `internal/repository/interfaces.go`
- `internal/repository/inmemory.go`
- `internal/repository/postgres.go`
- `internal/adapters/home_output.go`
