-- Ensure categories with NULL-subcategory posts have at least one subcategory.
insert into public.subcategory (category_id, name)
select c.id, 'general'
from public.category c
where exists (
    select 1
    from public.post p
    where p.category_id = c.id
      and p.subcategory_id is null
)
  and not exists (
    select 1
    from public.subcategory s
    where s.category_id = c.id
);

-- Backfill NULL subcategory_id from a deterministic per-category fallback.
with fallback_subcategory as (
    select category_id, min(id) as subcategory_id
    from public.subcategory
    group by category_id
)
update public.post p
set subcategory_id = f.subcategory_id
from fallback_subcategory f
where p.subcategory_id is null
  and p.category_id = f.category_id;

do $$
begin
  if exists (
    select 1
    from public.post
    where subcategory_id is null
  ) then
    raise exception 'cannot set post.subcategory_id NOT NULL: unresolved NULL values remain';
  end if;
end $$;

alter table public.post
  alter column subcategory_id set not null;
