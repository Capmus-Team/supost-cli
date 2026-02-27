# Home Output Wrapping + Renderer Tests

Date: 2026-02-27

## Summary
Improved `supost home` terminal rendering so each row wraps cleanly within a fixed width while preserving ANSI color styling, and added adapter-level tests to validate wrapping behavior.

## What Changed

### 1. Row width control for home output
- Added `homeRowWidth = 54` in `internal/adapters/home_output.go`.
- Header rendering now uses this shared width constant for consistent visual alignment.

### 2. Styled word pipeline for wrapped rendering
- Added `styledWord` type to carry text + color metadata.
- Reworked row rendering to build a token list in this order:
  - title (blue)
  - email (gray)
  - optional camera icon (uncolored)
  - relative time (magenta)
- Added helper functions:
  - `splitStyledWords`
  - `wrapStyledWords`
  - `renderStyledLine`
- Long lines now wrap into multiple terminal lines without dropping color formatting.

### 3. Added tests for wrapping behavior
- New file: `internal/adapters/home_output_test.go`
- Tests included:
  - `TestWrapStyledWords_RespectsWidth`
  - `TestWrapStyledWords_SplitsLongWord`
- Included minimal ANSI-stripping helper for assertions on visible width.

## Why This Matters
- Home output remains readable for long titles/emails.
- Formatting behavior is now deterministic and test-backed.
- Adapter output logic is safer to evolve without visual regressions.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0005-home_output_wrapping_and_tests.md`
