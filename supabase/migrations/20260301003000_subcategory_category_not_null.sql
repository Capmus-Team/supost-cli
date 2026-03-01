alter table public.subcategory
  alter column category_id set not null;

alter table public.subcategory
  drop constraint if exists subcategory_category_id_fkey;

alter table public.subcategory
  add constraint subcategory_category_id_fkey
  foreign key (category_id)
  references public.category(id)
  on update cascade
  on delete restrict
  not valid;

alter table public.subcategory
  validate constraint subcategory_category_id_fkey;
