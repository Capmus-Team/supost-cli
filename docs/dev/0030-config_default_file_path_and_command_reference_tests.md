# Config Default File Path and Command Reference Tests

Date: 2026-02-28

## Summary
Standardized root config-file behavior around `.supost.yaml` and added command-reference tests to lock in command tree, flags, and argument contracts.

## What Changed

### 1. Root config flag default now points to project-local `.supost.yaml`
- Updated `cmd/root.go`:
  - `--config` default changed from empty value to `.supost.yaml`
  - init logic now always sets Viper to an explicit config file path
  - removed home-directory fallback config lookup in favor of local default file behavior.

### 2. Added command-reference test coverage
- Added `cmd/command_reference_test.go` to assert:
  - top-level commands exist (`home`, `search`, `post`, `categories`, `serve`, `version`)
  - global flags exist and defaults are correct (`--verbose/-v`, `--format=json`, `--config=.supost.yaml`)
  - search flags and defaults are stable
  - post-create and post-respond flags are present
  - required post-respond flags are enforced (`--message`, `--reply-to`)
  - `post` and `post respond` argument validators require one `<post_id>`
  - `serve --port` default is `8080`.

### 3. Env example header note updated
- Updated `.env.example` top comments with an additional environment/secrets warning line.

## Why This Matters
- Makes local configuration behavior deterministic across developer machines.
- Adds regression safety for CLI command surface and defaults.
- Documents expected command/flag contracts in executable tests.

## Files in This Increment
- `.env.example`
- `cmd/root.go`
- `cmd/command_reference_test.go`
- `docs/dev/0030-config_default_file_path_and_command_reference_tests.md`
