# Remove Legacy Listings and Enforce Project Structure Contracts

Date: 2026-02-28

## Summary
Removed the old `listings` feature surface (domain/service/cmd/seed/repository interface remnants), updated preview server output to posts, and added/expanded command tests to enforce command and filesystem structure contracts reflected in README.

## What Changed

### 1. Removed legacy listings feature artifacts
- Deleted:
  - `cmd/listings.go`
  - `internal/domain/listing.go`
  - `internal/service/listings.go`
  - `internal/service/listings_test.go`
  - `testdata/seed/listings.json`
- Updated `internal/repository/interfaces.go`:
  - removed `ListingStore` interface.
- Updated `internal/repository/inmemory.go`:
  - removed listing map and listing CRUD/list methods
  - removed legacy listing seed loader
  - in-memory repository now focuses on posts/categories/messages.
- Updated `internal/domain/errors.go`:
  - removed listing-specific validation errors no longer used.

### 2. Updated serve preview endpoint semantics
- Updated `cmd/serve.go`:
  - switched from listing service to home service
  - replaced `GET /api/listings` with `GET /api/posts`
  - log output now documents `/api/posts` endpoint.

### 3. Root config flag description simplified
- Updated `cmd/root.go`:
  - `--config` help text now uses neutral `config file path` wording while keeping `.supost.yaml` default behavior.

### 4. README command and structure documentation refreshed
- Updated `README.md`:
  - `post` example now uses numeric post ID (`go run . post 130031605`)
  - project tree section updated to current command/domain/service/repository/adapter file layout
  - reflects removal of legacy listing paths and addition of newer page/flow files.

### 5. Command/project-structure contract tests expanded
- Updated `cmd/command_reference_test.go` with new checks:
  - assert legacy `listings` command is absent
  - assert README-listed files/directories exist
  - assert legacy listing files are removed
  - added helper assertions for repo-root path/file/dir checks.

## Why This Matters
- Eliminates obsolete listing-era code paths and seed data to reduce maintenance confusion.
- Aligns runtime preview behavior with current post-centric architecture.
- Adds stronger regression protection for CLI surface and documented repository structure.

## Files in This Increment
- `README.md`
- `cmd/root.go`
- `cmd/serve.go`
- `cmd/command_reference_test.go`
- `internal/domain/errors.go`
- `internal/repository/interfaces.go`
- `internal/repository/inmemory.go`
- `docs/dev/0031-remove_legacy_listings_and_enforce_project_structure_contracts.md`
- deleted: `cmd/listings.go`
- deleted: `internal/domain/listing.go`
- deleted: `internal/service/listings.go`
- deleted: `internal/service/listings_test.go`
- deleted: `testdata/seed/listings.json`
