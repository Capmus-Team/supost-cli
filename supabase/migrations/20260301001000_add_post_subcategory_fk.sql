-- Normalize orphaned subcategory references before enforcing FK.
update public.post p
set subcategory_id = null
where p.subcategory_id is not null
  and not exists (
    select 1
    from public.subcategory s
    where s.id = p.subcategory_id
  );

alter table "public"."post"
  add constraint "post_subcategory_id_fkey"
  FOREIGN KEY (subcategory_id)
  REFERENCES public.subcategory(id)
  ON UPDATE CASCADE
  ON DELETE SET NULL
  NOT VALID;

alter table "public"."post"
  validate constraint "post_subcategory_id_fkey";
