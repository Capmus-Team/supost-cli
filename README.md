
# supost

**A university marketplace CLI prototype.**  
It renders beautiful **website-like terminal views** directly on top of the live `supost.com` Supabase Postgres database.  
You can browse the homepage, search results, individual posts, and even submit new posts — all from your terminal.  

Production stack (future): **Supabase Postgres + Next.js/TypeScript web frontend**.  
This CLI is the perfect prototype: clean Go architecture today, easy port to Next.js tomorrow.

## Quick Start

No database, no API keys, no Docker required. Works immediately with in-memory seed data.

```bash
# 1. Build & run (after scaffolding)
go run ./cmd/supost version          # → v0.1.0

# 2. Try the website pages (Cobra structure)
go run ./cmd/supost website home                    # Full homepage shell
go run ./cmd/supost website search-results --subcategory 14
go run ./cmd/supost website post 130031605
```

For create-post flow and emails, copy the env template:

```bash
cp .env.example .env
```

Then run:
```bash
go run ./cmd/supost website create --help
```

## Essential Commands (Cobra Style)

All website simulation lives under the `website` parent command.

```bash
# Homepage
supost website home

# Search results page (filters + pagination)
supost website search-results --subcategory 14
supost website search-results --category 5 --page 2
supost website search-results --subcategory 59 --page 1 --limit 20

# Single post page
supost website post 130031605

# Reply to a post (sends Mailgun + saves message)
supost website post 130031783 \
  --message "Hello, I want to buy your bike" \
  --email-reply-to "gwientjes@gmail.com"

# Create-post flow
supost website create                              # Step 1: category chooser
supost website create --category 8                 # Step 2: subcategory
supost website create --category 5 --subcategory 14 # Step 3: full form

# Submit a new post (INSERT + publish email)
supost website create \
  --category 5 \
  --subcategory 14 \
  --name "Red bike for sale" \
  --body "Pick up on campus." \
  --email "wientjes@alumni.stanford.edu" \
  --price 100 \
  --submit

# Personals post (no price)
supost website create \
  --category 8 \
  --subcategory 130 \
  --name "Missed connection" \
  --body "Saw you at Coupa." \
  --email "wientjes@cs.stanford.edu" \
  --submit
```

## Create-Post Validation

Exact same rules as the real site:
- Email required + must be Stanford-affiliated
- Name & Body required
- Price required only for certain categories
- Errors shown in the exact friendly format you know.

## Email Features (Mailgun)

- Successful submit → generates `access_token` + sends “SUpost - Publish your post!” email
- Reply → sends response email + saves to `app_private.message` table

Required env vars (in `.env`):
- `MAILGUN_DOMAIN`, `MAILGUN_API_KEY`, `MAILGUN_FROM_EMAIL`
- Optional: `MAILGUN_API_BASE`, `MAILGUN_SEND_TIMEOUT`, `SUPOST_BASE_URL`

## Development

```bash
cp .env.example .env
cp configs/config.yaml.example .supost.yaml

make check   # format, lint, build, test
make build   # compile to bin/supost
make test
make serve   # optional HTTP preview
make clean
```

## Project Structure

See **[AGENTS.md](AGENTS.md)** for the full architecture guide (Cobra commands in `internal/cli/commands/`, pure logic in `internal/core/`, adapters for Supabase/Mailgun).

## Connecting Real Supabase

Defaults to fast in-memory storage.  
To go live:
1. Add your Supabase credentials to `.env`
2. Adapter auto-switches
3. All commands work against real data

## Migration to Production (Next.js)

1. Apply same migrations to Supabase
2. Port `internal/core/domain/` structs → TypeScript
3. Port `internal/core/service/` → Next.js API routes
4. Keep validation & email logic forever

## Full Command Reference & Docs

- Detailed workflows → `docs/cli-command-reference.md`
- Schema & progress → `docs/`
- Latest schema → `docs/reference/supost_2-22-26_schema.sql`

---

**Last updated**: February 2026  
Built with clean Cobra + standard Go layout. Ready for you and any AI to extend forever.

Happy prototyping! Run `supost website home` and watch the website appear in your terminal.
