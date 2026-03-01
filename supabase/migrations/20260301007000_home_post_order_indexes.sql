create index if not exists idx_post_active_time_posted_id_desc
on public.post using btree (time_posted desc, id desc)
where (status = 1);

create index if not exists idx_post_active_category_time_posted_id_desc
on public.post using btree (category_id, time_posted desc, id desc)
where (status = 1);
