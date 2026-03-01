set statement_timeout = 0;

alter table public.post
  add column if not exists fts tsvector
  generated always as (
    setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(body, '')), 'B')
  ) stored;

create index if not exists post_fts_idx
  on public.post
  using gin (fts);

reset statement_timeout;
