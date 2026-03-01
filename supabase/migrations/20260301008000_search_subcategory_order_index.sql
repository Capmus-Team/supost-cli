create index if not exists idx_post_active_subcategory_time_posted_id_desc
on public.post using btree (subcategory_id, time_posted desc, id desc)
where (status = 1);
