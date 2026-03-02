-- Public photo data: allow read access via visible parent posts, block writes via RLS.
alter table public.photo enable row level security;

drop policy if exists photo_public_read on public.photo;
create policy photo_public_read
on public.photo
for select
to anon, authenticated
using (
  exists (
    select 1
    from public.post p
    where p.id = photo.post_id
  )
);
