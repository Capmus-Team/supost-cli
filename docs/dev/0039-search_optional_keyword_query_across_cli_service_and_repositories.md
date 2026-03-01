# Search Optional Keyword Query Across CLI, Service, and Repositories

Date: 2026-03-01

## Summary
Extended `search` to accept an optional free-text query argument and wired keyword filtering through service, in-memory repository, Postgres repository, renderer title logic, tests, and README examples.

## What Changed

### 1. Added optional query args to search command
- Updated `cmd/search.go`:
  - usage changed to `search [query]`
  - args policy set to `cobra.ArbitraryArgs`
  - joins CLI args into a trimmed query string
  - passes query into `SearchService.Search`.

### 2. Search result contract now includes query
- Updated `internal/domain/search_result.go`:
  - added `Query string` field (`json:"query"`).

### 3. Service now normalizes and forwards query text
- Updated `internal/service/search.go`:
  - `Search` signature now accepts `query string`
  - trims query via `normalizeSearchQuery`
  - forwards query to repository
  - returns normalized query in `SearchResultPage`.

### 4. In-memory search now filters by query terms in name/body
- Updated `internal/repository/inmemory_search.go`:
  - `SearchActivePosts` now accepts query
  - applies `matchesPostQuery`:
    - case-insensitive term matching
    - all query terms must appear in either post name or body.

### 5. Postgres search now applies full-text filter
- Updated `internal/repository/postgres_search.go`:
  - `SearchActivePosts` now accepts query text
  - adds optional full-text condition:
    - `to_tsvector('simple', coalesce(name,'') || ' ' || coalesce(body,'')) @@ plainto_tsquery('simple', $query)`
  - keeps filter optional when query is empty.

### 6. Search page title reflects keyword searches
- Updated `internal/adapters/search_output.go`:
  - title becomes `all posts` when query is empty
  - title becomes `search: <query>` when query is provided.

### 7. Tests and docs updated for new query behavior
- Updated:
  - `cmd/command_reference_test.go` (search accepts optional args)
  - `internal/adapters/search_output_test.go` (query title coverage)
  - `internal/repository/inmemory_search_test.go` (name/body term matching coverage)
  - `internal/service/search_test.go` (query forwarding/normalization coverage)
  - `README.md` examples and command synopsis for `search [query]`.

## Why This Matters
- Enables Craigslist-style keyword searching without requiring category filters.
- Keeps query behavior consistent across in-memory prototyping and Postgres-backed execution.
- Preserves stable defaults (`all posts`) when no query is supplied.

## Files in This Increment
- `README.md`
- `cmd/command_reference_test.go`
- `cmd/search.go`
- `internal/adapters/search_output.go`
- `internal/adapters/search_output_test.go`
- `internal/domain/search_result.go`
- `internal/repository/inmemory_search.go`
- `internal/repository/inmemory_search_test.go`
- `internal/repository/postgres_search.go`
- `internal/service/search.go`
- `internal/service/search_test.go`
- `docs/dev/0039-search_optional_keyword_query_across_cli_service_and_repositories.md`
