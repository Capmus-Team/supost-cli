# supost

A university marketplace CLI that renders website-like terminal views on top of the live SUPost Supabase Postgres database. Browse the homepage, search listings, view posts, create new posts, and send response emails — all from your terminal.

Production stack (future): Supabase + Next.js. This CLI is the prototype: clean Go architecture today, easy port to TypeScript tomorrow.

## Quick Start

Works immediately with in-memory seed data — no database, no API keys, no Docker.

```bash
go run . version              # → v0.1.0
go run . home                 # homepage
go run . search --category 5  # search results
go run . post seed-1          # view a post
go run . serve                # preview HTTP server at localhost:8080
```

To connect to the real SUPost database:

```bash
cp .env.example .env          # fill in credentials
go run . home                 # now renders live data
```

## Commands

### Pages (read-only)

```bash
# Homepage
supost home

# Search results with filters + pagination
supost search --category 5
supost search --subcategory 14
supost search --category 5 --page 2 --per-page 20

# View a single post
supost post 130031605

# List categories (utility)
supost categories
```

### Create a Post

The `post create` command handles the full create-post wizard. Flags determine which step you're on:

```bash
# Step 1: choose category
supost post create

# Step 2: choose subcategory
supost post create --category 8

# Step 3: show form fields
supost post create --category 5 --subcategory 14

# Submit (all required fields present → validates + INSERTs + sends publish email)
supost post create \
  --category 5 \
  --subcategory 14 \
  --name "Red bike for sale" \
  --body "Pick up on campus." \
  --email "wientjes@alumni.stanford.edu" \
  --price 100

# Personals post (no price field for category 8)
supost post create \
  --category 8 \
  --subcategory 130 \
  --name "Missed connection" \
  --body "Saw you at Coupa." \
  --email "wientjes@cs.stanford.edu"

# Dry run: validate + render email, no INSERT, no send
supost post create \
  --category 5 --subcategory 14 \
  --name "Test" --body "Test" --email "test@stanford.edu" --price 50 \
  --dry-run
```

### Respond to a Post

```bash
# Send a response email to the post owner (+ saves to messages table)
supost post respond 130031783 \
  --message "Hello, I want to buy your bike" \
  --reply-to "gwientjes@gmail.com"

# Dry run: validate + render email, don't send, don't persist
supost post respond 130031783 \
  --message "Test message" \
  --reply-to "test@gmail.com" \
  --dry-run
```

### Utility

```bash
supost version                # print version
supost serve                  # preview HTTP server
supost serve --port 3000      # custom port
```

## Command Reference

```
supost
├── home                          # render homepage
├── search                        # render search results page
│     --category <id>
│     --subcategory <id>
│     --page <n>                  (default: 1)
│     --per-page <n>              (default: 20)
├── post <post_id>                # render single post page
├── post create                   # create-post wizard / submit
│     --category <id>
│     --subcategory <id>
│     --name <string>
│     --body <string>
│     --email <string>
│     --price <amount>            (required for some categories)
│     --dry-run                   (validate only, no write)
├── post respond <post_id>        # send response email
│     --message <string>          (required)
│     --reply-to <email>          (required)
│     --dry-run                   (validate only, no send)
├── categories                    # list categories + subcategories
├── serve                         # preview HTTP server
│     --port <n>                  (default: 8080)
└── version                       # print version
```

### Global Flags (available on all commands)

```
--verbose, -v       enable verbose/debug output
--format <string>   output format: json, table, text (default: json)
--config <path>     config file (default: .supost.yaml)
```

## Project Structure

```
supost-cli/
├── AGENTS.md                        # AI agent governance — read first
├── main.go                          # entrypoint (wiring only)
│
├── cmd/                             # one file per command
│   ├── root.go                      # global flags, config init
│   ├── version.go                   # supost version
│   ├── home.go                      # supost home
│   ├── search.go                    # supost search
│   ├── post.go                      # supost post <id>
│   ├── post_create.go               # supost post create
│   ├── post_respond.go              # supost post respond <id>
│   ├── categories.go                # supost categories
│   └── serve.go                     # supost serve
│
├── internal/
│   ├── config/config.go             # centralized config (Viper)
│   ├── domain/                      # types → Supabase tables
│   │   ├── post.go                  # Post struct (json + db tags)
│   │   ├── category.go              # Category, Subcategory
│   │   ├── user.go                  # User / Profile
│   │   ├── message.go               # Response messages
│   │   └── errors.go                # domain errors (HTTP-mappable)
│   ├── service/                     # business logic (the brain)
│   │   ├── posts.go                 # Create, Validate, GetByID, Search
│   │   ├── categories.go            # ListCategories, GetSubcategories
│   │   ├── email.go                 # SendPublishLink, SendResponse
│   │   └── home.go                  # HomepageData (posts + categories)
│   ├── repository/                  # data access (swappable)
│   │   ├── interfaces.go
│   │   ├── inmemory.go              # zero-dep prototype adapter
│   │   └── postgres.go              # real Supabase/Postgres adapter
│   ├── adapters/                    # external services
│   │   ├── output.go                # JSON/table/text rendering
│   │   └── mailgun.go               # email sending
│   └── util/
│
├── migrations/                      # SQL schema (Supabase-ready)
├── configs/config.yaml.example
├── testdata/seed/
└── .env.example
```

## Create-Post Validation

When submitting a post, the following rules apply:

- **Email** is required and must be Stanford-affiliated:
  - Any `*.stanford.edu` domain (e.g., `@stanford.edu`, `@cs.stanford.edu`, `@gsb.stanford.edu`)
  - Also: `@stanfordalumni.org`, `@stanfordchildrens.org`, `@stanfordhealthcare.org`, `@stanfordmed.org`, `@lpch.org`
- **Name** and **Body** are required
- **Price** is category-dependent:
  - Required for: for sale/wanted (5), housing offering (3)
  - Not available for: personals (8), housing need (4), community (9), service offered (7), campus job (1), job off-campus (2)

Validation errors follow the same format as the real site:

```
1 error prohibited this post from being saved
There were problems with the following fields:

Email must be a Stanford email (e.g., @stanford.edu, @cs.stanford.edu).
```

## Email Features (Mailgun)

### Publish-Link Confirmation

After successful post creation:
- Generates a post `access_token`
- Sends email with subject: `SUpost - Publish your post! <post name>`
- Includes publish URL: `<SUPOST_BASE_URL>/post/publish/<access_token>`

### Post Response

When sending a response:
- Sends email to the post owner's stored email
- Sets `Reply-To` header to `--reply-to` address
- Saves message to `app_private.message` table

## Environment Variables

```bash
# Database (leave empty for in-memory prototype)
DATABASE_URL=                       # read/write Postgres connection
# DATABASE_READ_URL=                # optional: separate read-only connection

# Supabase
SUPABASE_URL=
SUPABASE_ANON_KEY=
SUPABASE_SERVICE_ROLE_KEY=

# Mailgun (required for email features)
MAILGUN_DOMAIN=
MAILGUN_API_KEY=
MAILGUN_FROM_EMAIL=
MAILGUN_API_BASE=                   # https://api.mailgun.net (US) or https://api.eu.mailgun.net (EU)
MAILGUN_SEND_TIMEOUT=10s

# App
SUPOST_BASE_URL=https://n.supost.com
PORT=8080
VERBOSE=false
FORMAT=json
```

## Development

```bash
cp .env.example .env
cp configs/config.yaml.example .supost.yaml

make check    # format, vet, build, test
make build    # compile to bin/supost
make test     # tests with race detector
make serve    # preview HTTP server
make clean
```

## Migration to Production (Next.js + Supabase)

1. **Schema** → Apply `migrations/*.sql` to Supabase, uncomment RLS policies
2. **Types** → Translate `internal/domain/*.go` structs to TypeScript interfaces (`json` tags = field names)
3. **Logic** → Port `internal/service/*.go` to Next.js API routes (nearly 1:1)
4. **Data access** → Replace Go repository with Supabase JS SDK
5. **Auth** → Replace CLI email validation with Supabase Auth + RLS
6. **Seed data** → Import `testdata/seed/*.json` into Supabase

---

*Built with Cobra + clean Go architecture. See [AGENTS.md](AGENTS.md) for the full governance guide.*
