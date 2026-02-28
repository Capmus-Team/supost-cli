# Post Create and Post Respond Validation Regression Tests

Date: 2026-02-28

## Summary
Expanded service-layer validation tests for `post create submit` and `post respond` to cover additional edge-case ordering and input-validation behaviors.

## What Changed

### 1. Added post-create submit validation edge-case tests
- Updated `internal/service/post_create_submit_test.go` with:
  - `TestPostCreateService_Submit_PriceForbiddenBeforeSubcategoryMismatch`
    - verifies category-level price-forbidden rule is reported when personals category gets a provided price, even if subcategory ID belongs to another category.
  - `TestPostCreateService_Submit_SubcategoryMismatchWithoutPrice`
    - verifies that once price rule does not block validation, mismatched subcategory/category combinations return the expected not-found mismatch message.

### 2. Added post-respond reply-to validation tests
- Updated `internal/service/post_respond_test.go` with:
  - `TestPostRespondService_ValidationReplyToRequired`
    - ensures empty `reply_to` is rejected.
  - `TestPostRespondService_ValidationReplyToEmailFormat`
    - ensures malformed reply-to addresses are rejected with the expected message.

## Why This Matters
- Improves confidence in validation ordering for post-create submissions.
- Strengthens response-flow validation guarantees around sender contact requirements.
- Prevents regressions in user-facing error behavior for key CLI flows.

## Files in This Increment
- `internal/service/post_create_submit_test.go`
- `internal/service/post_respond_test.go`
- `docs/dev/0029-post_create_and_post_respond_validation_regression_tests.md`
