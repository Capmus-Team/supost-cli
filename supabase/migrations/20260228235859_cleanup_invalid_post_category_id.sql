-- Normalize orphaned category references before enforcing post.category_id FK.
-- Keeps nullable semantics by setting invalid non-null category IDs to NULL.
update public.post p
set category_id = null
where p.category_id is not null
  and not exists (
    select 1
    from public.category c
    where c.id = p.category_id
  );
