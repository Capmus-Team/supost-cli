# Post Create Category Price Rules and Optional-Price Persistence

Date: 2026-02-28

## Summary
Refined post-create validation and rendering so price behavior now depends on category rules: some categories require price, some forbid it, and repositories persist price as optional rather than always set.

## What Changed

### 1. Post-create submit validation now enforces category-specific price policy
- Updated `internal/service/post_create_submit.go`:
  - switched from always-required price to category-aware checks using domain price rules.
  - if category requires price:
    - missing price is rejected
    - negative price is rejected
  - if category disallows price:
    - provided price is rejected.
  - validation errors now use a structured, multi-line format:
    - header like `N errors prohibited this post from being saved`
    - followed by field-specific messages.

### 2. Service tests expanded for price-rule scenarios
- Updated `internal/service/post_create_submit_test.go` with coverage for:
  - required-price categories (for sale) rejecting missing price
  - no-price categories (personals) rejecting provided price
  - successful personals submit without price
  - formatted validation header assertions.

### 3. Form renderer hides price field for categories where price is not allowed
- Updated `internal/adapters/post_create_output.go`:
  - form stage now conditionally renders `Price: [price]` only when category rules allow it.
- Updated `internal/adapters/post_create_output_test.go`:
  - added test verifying personals form omits price field.

### 4. Repository persistence now respects optional price
- Updated `internal/repository/inmemory_post_create.go`:
  - `HasPrice` now reflects `submission.PriceProvided`.
- Updated `internal/repository/postgres_post_create.go`:
  - inserts `NULL` for price when `PriceProvided` is false
  - inserts numeric value only when price was provided.

## Why This Matters
- Aligns create-flow behavior with category semantics instead of one-size-fits-all pricing.
- Prevents invalid data combinations (e.g., personals posts with prices).
- Keeps DB persistence consistent with whether the user actually supplied price.

## Files in This Increment
- `internal/service/post_create_submit.go`
- `internal/service/post_create_submit_test.go`
- `internal/adapters/post_create_output.go`
- `internal/adapters/post_create_output_test.go`
- `internal/repository/inmemory_post_create.go`
- `internal/repository/postgres_post_create.go`
- `docs/dev/0028-post_create_category_price_rules_and_optional_price_persistence.md`
