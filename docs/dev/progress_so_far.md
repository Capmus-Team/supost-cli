# Development Progress So Far

Last updated: 2026-02-27

## Purpose
This document captures what has already been implemented in the `supost-cli` prototype so new contributors can quickly understand the current baseline.

## Milestones Completed

### 1. Initial prototype scaffold
Commit: `c0c3164`

Delivered foundational CLI architecture and project structure:
- Cobra command tree (`root`, `version`, `listings`, `serve`)
- Core packages under `internal/` (`config`, `domain`, `service`, `repository`, `adapters`, `util`)
- In-memory listing repository for zero-dependency prototyping
- Listing service and unit tests
- SQL migrations for `profiles` and `listings`
- Seed data and Makefile development workflow

### 2. URL preview support
Commit: `7076f00`

Added functionality for URL preview behavior in listing data/flow (as reflected in current domain/config evolution).

### 3. Mailgun environment support
Commit: `f521158`

Added Mailgun-related environment variables to support email integration groundwork.

### 4. `.env.example` refresh
Commit: `56741a3`

Updated the environment template to keep local setup aligned with currently supported config values.

## Current Working Capabilities

As of this update, the repository provides:
- `go run . version` for CLI version output
- `go run . listings` returning active listings (JSON by default)
- `go run . serve` exposing preview HTTP endpoints:
  - `GET /api/listings`
  - `GET /api/health`
- In-memory data path with no required external services
- Baseline tests for listing service behavior and validation logic

## Architecture Baseline In Place

- Domain-first structure in `internal/domain` intended as future API/schema contract
- Service layer separated from CLI transport concerns
- Repository abstraction with in-memory implementation as default adapter
- Output rendering centralized in `internal/adapters/output.go`
- Migration files as schema source of truth (`migrations/*.sql`)

## Next Recommended Work

1. Add/complete PostgreSQL repository adapter and switch adapter selection via config.
2. Expand command behavior for currently-declared flags (for example `--category` filtering on `listings`).
3. Wire seed loading from `testdata/seed/*.json` in the in-memory adapter.
4. Increase service test coverage and include more regression-focused cases.
5. Replace `log.Printf` in `serve` path with structured logging to match project standards.
