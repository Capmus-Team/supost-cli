create index if not exists idx_post_active_fts_idx
on public.post using gin (fts)
where (status = 1);
