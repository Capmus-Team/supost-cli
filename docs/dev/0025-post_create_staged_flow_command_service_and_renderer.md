# Post Create Staged Flow: Command, Service, and Renderer

Date: 2026-02-28

## Summary
Added a staged post-creation flow under `post create`, including a dedicated domain page model, service orchestration for stage selection, and adapter rendering for category/subcategory/form steps.

## What Changed

### 1. New `post create` command
- Added `cmd/post_create.go` as a `post` subcommand.
- Supports flags:
  - `--category`
  - `--subcategory`
- Loads repository adapter from config, builds page view model via service, and renders through post-create adapter.

### 2. New staged post-create domain contract
- Added `internal/domain/post_create_page.go`:
  - stage constants:
    - `choose_category`
    - `choose_subcategory`
    - `form`
  - `PostCreatePage` view model with selected IDs/names and category/subcategory lists.

### 3. New post-create service and tests
- Added `internal/service/post_create.go`:
  - `PostCreateRepository` interface for taxonomy reads
  - `BuildPage(...)` stage resolution logic:
    - no category: choose category stage
    - category only: choose subcategory stage
    - category + subcategory: form stage
    - infers category from subcategory when category ID is omitted
  - not-found handling for invalid category/subcategory selections.
- Added `internal/service/post_create_test.go`:
  - covers all stages, category inference, and invalid-category error behavior.

### 4. New post-create page renderer and tests
- Added `internal/adapters/post_create_output.go`:
  - renders shared header/footer with taxonomy breadcrumb support
  - renders campus band + staged content:
    - category menu with curated order/labels
    - subcategory list view
    - form skeleton (title/price/description/email/photos/preview)
  - reuses existing housing-policy notice in form stage.
- Added `internal/adapters/post_create_output_test.go`:
  - verifies output for category stage, subcategory stage, and form stage.

## Why This Matters
- Establishes a clear multi-step post creation UX in CLI form that maps well to future web flow.
- Keeps flow logic in service layer and rendering concerns in adapter layer.
- Reuses shared header/footer/taxonomy conventions for consistency across page-like commands.

## Files in This Increment
- `cmd/post_create.go`
- `internal/domain/post_create_page.go`
- `internal/service/post_create.go`
- `internal/service/post_create_test.go`
- `internal/adapters/post_create_output.go`
- `internal/adapters/post_create_output_test.go`
- `docs/dev/0025-post_create_staged_flow_command_service_and_renderer.md`
