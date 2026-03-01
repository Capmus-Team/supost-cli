# Subcategory Category Not-Null and Restrictive FK

Date: 2026-03-01

## Summary
Added a Supabase migration to enforce non-null `subcategory.category_id` and tighten foreign-key behavior to prevent category deletion when dependent subcategories exist.

## What Changed

### 1. Enforced non-null subcategory category link
- Added `supabase/migrations/20260301003000_subcategory_category_not_null.sql`:
  - sets `public.subcategory.category_id` to `NOT NULL`.

### 2. Rebuilt subcategory→category foreign key with stricter delete behavior
- Same migration:
  - drops existing `subcategory_category_id_fkey` if present
  - recreates FK referencing `public.category(id)`
  - uses `ON UPDATE CASCADE`
  - uses `ON DELETE RESTRICT` (instead of nullable/delete-relaxing behavior)
  - adds as `NOT VALID` then validates.

## Why This Matters
- Guarantees every subcategory belongs to a category.
- Protects taxonomy integrity by blocking category deletion while subcategories still reference it.

## Files in This Increment
- `supabase/migrations/20260301003000_subcategory_category_not_null.sql`
- `docs/dev/0036-subcategory_category_not_null_and_restrict_fk.md`
