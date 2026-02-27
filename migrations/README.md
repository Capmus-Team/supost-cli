# Migrations

SQL files that define the database schema. **Single source of truth.**

## Applying Locally

```bash
psql $DATABASE_URL -f migrations/001_create_profiles.sql
psql $DATABASE_URL -f migrations/002_create_listings.sql
```

## Applying to Supabase

1. Open the Supabase SQL Editor
2. Paste each migration in order
3. Uncomment the RLS policies
4. Run

## Rules

- Never modify an existing migration — create a new one
- Number sequentially: `003_add_images.sql`
- Each migration should be idempotent (`IF NOT EXISTS`)
- Include commented-out RLS policies for Supabase
