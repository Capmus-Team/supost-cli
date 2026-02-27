# supost

A university marketplace CLI prototype. Production stack: Supabase + Next.js.

## Quick Start

```bash
go run . version              # → v0.1.0
go run . listings             # → JSON list of seed data (no DB needed)
go run . serve                # → preview server at http://localhost:8080
curl localhost:8080/api/listings  # → JSON response
```

No database, no API keys, no Docker required. Works immediately.

## Development

```bash
cp .env.example .env                # optional: configure database
cp configs/config.yaml.example .supost.yaml

make check    # format, vet, build, test
make build    # compile binary to bin/
make test     # run tests with race detector
make serve    # start preview HTTP server
make migrate  # show migration commands
make clean    # remove build artifacts
```

## Project Structure

See [AGENTS.md](AGENTS.md) for the full architecture guide.

## Connecting a Database

The app uses in-memory seed data by default. To connect Postgres/Supabase:

1. Set `DATABASE_URL` in `.env`
2. Run migrations: `make migrate`
3. Swap the repository adapter in `cmd/` files (see AGENTS.md §6.5)

## Migration to Production

1. Apply `migrations/*.sql` to Supabase (uncomment RLS policies)
2. Translate `internal/domain/` structs → TypeScript interfaces
3. Port `internal/service/` logic → Next.js API routes
4. Import `testdata/seed/` into Supabase for test data
