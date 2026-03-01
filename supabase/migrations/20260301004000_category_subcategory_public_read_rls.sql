-- Public taxonomy data: allow read access, block writes via RLS.
alter table public.category enable row level security;
alter table public.subcategory enable row level security;

drop policy if exists category_public_read on public.category;
create policy category_public_read
on public.category
for select
to anon, authenticated
using (true);

drop policy if exists subcategory_public_read on public.subcategory;
create policy subcategory_public_read
on public.subcategory
for select
to anon, authenticated
using (true);
