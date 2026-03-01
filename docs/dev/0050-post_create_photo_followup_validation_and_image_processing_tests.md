# Post Create Photo Follow-up Validation and Image Processing Tests

Date: 2026-03-01

## Summary
Added follow-up hardening for post-create photo flow: stricter CLI photo-file validation tests, S3 uploader image decode/resize/ticker upload behavior with coverage, dry-run uploader requirements in service tests, and in-memory photo persistence tests.

## What Changed

### 1. Expanded post-create CLI photo input validation tests
- Updated `cmd/post_create_test.go`:
  - added rejection test for blank photo paths
  - added rejection test for empty photo files.

### 2. Enhanced S3 photo uploader behavior
- Updated `internal/adapters/s3_photo_uploader.go`:
  - decodes source images via stdlib image codecs
  - resizes post and ticker variants to max widths (340 and 220)
  - encodes output bytes by detected format (jpeg/png/gif)
  - uploads both main and `ticker_` objects and returns both keys.

### 3. Added S3 uploader unit tests
- Updated `internal/adapters/s3_photo_uploader_test.go`:
  - validates bucket requirement for uploader construction
  - verifies resize behavior (downscale + no-upscale)
  - verifies format selection and JPEG fallback behavior.

### 4. Added post-create service uploader requirement tests
- Updated `internal/service/post_create_submit_test.go`:
  - rejects non-dry-run submissions with photos when uploader is missing
  - confirms dry-run still allows photos without uploader
  - preserves `photo_count` reporting in dry-run path.

### 5. Added in-memory repository photo persistence tests
- Added `internal/repository/inmemory_post_create_test.go`:
  - verifies `SavePostPhotos` stores rows and marks post `HasImage=true`
  - validates blank `S3Key` rejection.

## Why This Matters
- Prevents invalid photo inputs from reaching upload/persistence layers.
- Ensures uploaded photo assets include both post and ticker variants with predictable sizing.
- Strengthens regression safety for dry-run vs real-submit behavior.

## Files in This Increment
- `cmd/post_create_test.go`
- `internal/adapters/s3_photo_uploader.go`
- `internal/adapters/s3_photo_uploader_test.go`
- `internal/service/post_create_submit_test.go`
- `internal/repository/inmemory_post_create_test.go`
- `docs/dev/0050-post_create_photo_followup_validation_and_image_processing_tests.md`
