# Supabase Scripts

## Incremental Photo Backfill

Use this to backfill `public.photo` from legacy columns on `public.post` without repeating work.

- Script: `supabase/scripts/backfill_photo_incremental.sql`
- Checkpoint table: `app_private.backfill_checkpoint`
- Does not modify `public.post` data.

Run one batch (default `batch_size=1000`):

```bash
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f supabase/scripts/backfill_photo_incremental.sql
```

Run with custom batch size + job name:

```bash
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -v batch_size=2000 -v job_name=photo_backfill_v1 -f supabase/scripts/backfill_photo_incremental.sql
```

Run repeatedly until output shows:

- `selected_post_count = 0`

Check progress:

```sql
select * from app_private.backfill_checkpoint order by updated_at desc;
```
