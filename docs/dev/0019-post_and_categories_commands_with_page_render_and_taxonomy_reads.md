# Post + Categories Commands, Taxonomy Reads, and Post Page Renderer

Date: 2026-02-28

## Summary
Added new `post` and `categories` commands, expanded service/repository/domain support for single-post and taxonomy reads, and introduced a dedicated terminal post page renderer with tests.

## What Changed

### 1. New CLI commands
- Added `cmd/categories.go`:
  - command: `categories`
  - loads config, selects repo adapter (in-memory or Postgres), calls `CategoryService`, renders output via shared adapter rendering.
- Added `cmd/post.go`:
  - command: `post <post_id>`
  - validates numeric post ID input
  - loads one post via `PostService`
  - maps not-found to a user-facing error message
  - renders with `RenderPostPage` for default/text/table formats.

### 2. New domain taxonomy model + expanded post fields
- Added `internal/domain/category.go`:
  - `Category`
  - `Subcategory`
  - `CategoryWithSubcategories`
  - all with `json` and `db` tags.
- Updated `internal/domain/post.go`:
  - added image/file source fields needed by post-page rendering:
    - `photo1_file_name` .. `photo4_file_name`
    - `image_source1` .. `image_source4`.

### 3. New services and tests
- Added `internal/service/categories.go` + `internal/service/categories_test.go`:
  - orchestrates category list + subcategory grouping
  - stable sorting for categories and subcategories.
- Added `internal/service/post.go` + `internal/service/post_test.go`:
  - simple single-post retrieval use-case (`GetByID`).

### 4. Repository contracts and implementations expanded
- Updated `internal/repository/interfaces.go`:
  - added `PostStore` (`GetPostByID`)
  - added `CategoryStore` (`ListCategories`, `ListSubcategories`).
- Updated `internal/repository/inmemory.go`:
  - added in-memory `GetPostByID`
  - added category/subcategory storage + list methods
  - added category/subcategory seed loading from `testdata/seed/*.json` with fallback defaults.
- Updated `internal/repository/postgres.go`:
  - added parameterized `GetPostByID` query
  - added `ListCategories` and `ListSubcategories` queries.

### 5. New post page adapter and tests
- Added `internal/adapters/post_output.go` + `internal/adapters/post_output_test.go`:
  - shared page header/footer integration
  - post header block (title/email/date/price)
  - 2x2 photo URL grid with quadrant mapping
  - body rendering
  - right-side "Message Poster" panel
  - commercial-service warning text.

### 6. Command docs updated
- Updated `docs/terminal_commands.md` with:
  - `go run . categories`
  - `go run . post 130031900`.

## Why This Matters
- Moves the CLI closer to full page-level parity with web-style post detail views.
- Keeps layering aligned (`cmd` wiring, `service` orchestration, `repository` data access, `adapters` rendering).
- Preserves swapability between in-memory and Postgres adapters while adding new read capabilities.

## Files in This Increment
- `cmd/categories.go`
- `cmd/post.go`
- `docs/terminal_commands.md`
- `internal/adapters/post_output.go`
- `internal/adapters/post_output_test.go`
- `internal/domain/category.go`
- `internal/domain/post.go`
- `internal/repository/inmemory.go`
- `internal/repository/interfaces.go`
- `internal/repository/postgres.go`
- `internal/service/categories.go`
- `internal/service/categories_test.go`
- `internal/service/post.go`
- `internal/service/post_test.go`
- `docs/dev/0019-post_and_categories_commands_with_page_render_and_taxonomy_reads.md`
