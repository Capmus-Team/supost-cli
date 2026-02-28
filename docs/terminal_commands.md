# Terminal Commands (Without Installed Binary)

Use these commands from the repo root when you want to run the latest local code directly, without relying on `/usr/local/bin/supost`.

## Run the CLI directly

```bash
go run . version
go run . home
go run . categories
go run . post 130031900
go run . home --refresh
go run . home --format json
go run . home --cache-ttl 60s
go run . home --limit 50
```

## Build and run local repo binary

```bash
make build
./bin/supost home
```

This uses the binary built from your current working tree, not the globally installed one.

If `go build -o bin/supost .` prints no output, that means the build succeeded.

## Reinstall global binary

```bash
go build -o /usr/local/bin/supost .
supost version
```

Use this when `supost` in your shell is outdated compared to current repo code.

## Home command notes

- Default output is terminal-friendly text rendering.
- JSON output is available for web/frontend integration checks:
  - `go run . home --format json`
- `--refresh` bypasses cache.
- `--cache-ttl` controls cache duration (set `0s` to disable cache).
- Performance path:
  - cache 1: recent active posts
  - cache 2: category last-active timestamps
  - category/subcategory taxonomy is loaded from local seed data, not from runtime DB queries.
