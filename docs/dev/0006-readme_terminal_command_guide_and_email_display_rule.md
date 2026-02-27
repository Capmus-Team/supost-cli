# README + Terminal Command Guide + Home Email Display Rule

Date: 2026-02-27

## Summary
Documented how to run the latest local code without relying on an installed global binary, linked that guide from README, and adjusted home feed display so Stanford emails are normalized to `@stanford.edu` while non-Stanford emails are hidden in the terminal view.

## What Changed

### 1. Added terminal command reference doc
- New file: `docs/terminal_commands.md`
- Covers:
  - running via `go run . ...`
  - building and running `./bin/supost`
  - reinstalling `/usr/local/bin/supost`
  - `home` command cache/format flags

### 2. README now links terminal command guide
- `README.md` now includes a link to `docs/terminal_commands.md` in Quick Start so users can avoid stale global binary confusion.

### 3. Home renderer email display normalization
- Updated `internal/adapters/home_output.go`:
  - Added `formatDisplayEmail(email string) string`
  - Emails containing `stanford.edu` are rendered as `@stanford.edu`
  - Non-Stanford emails render as empty display value in the terminal home output

### 4. Added tests for email display behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestFormatDisplayEmail_StanfordDomainCollapses`
  - `TestFormatDisplayEmail_NonStanfordUnchanged`

## Why This Matters
- Reduces local execution confusion between repo code and a stale globally installed binary.
- Keeps terminal home output aligned with expected display conventions for Stanford-domain posters.
- Adds regression coverage for the new email formatting rule.

## Files in This Increment
- `README.md`
- `docs/terminal_commands.md`
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0006-readme_terminal_command_guide_and_email_display_rule.md`
