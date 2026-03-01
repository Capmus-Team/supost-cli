-- Backfill category_id from subcategory when possible.
update public.post p
set category_id = s.category_id
from public.subcategory s
where p.category_id is null
  and p.subcategory_id = s.id
  and s.category_id is not null;

-- Ensure fallback category exists before applying non-null constraint.
do $$
begin
  if not exists (
    select 1
    from public.category
    where id = 9
  ) then
    raise exception 'fallback category id 9 not found in public.category';
  end if;
end $$;

-- Assign unresolved rows to "community" category (id = 9).
update public.post
set category_id = 9
where category_id is null;

alter table public.post
  alter column category_id set not null;
