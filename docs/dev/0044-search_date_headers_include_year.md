# Search Date Headers Include Year

Date: 2026-03-01

## Summary
Updated search result date headers to include the year in rendered group labels.

## What Changed

### 1. Expanded search date header format
- Updated `internal/adapters/search_output.go`:
  - changed `formatSearchDateHeader` from:
    - `Mon, Jan 2`
  - to:
    - `Mon, Jan 2, 2006`.

## Why This Matters
- Makes grouped search dates unambiguous across year boundaries.

## Files in This Increment
- `internal/adapters/search_output.go`
- `docs/dev/0044-search_date_headers_include_year.md`
