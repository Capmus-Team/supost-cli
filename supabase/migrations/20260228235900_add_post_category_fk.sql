alter table "public"."post"
  add constraint "post_category_id_fkey"
  FOREIGN KEY (category_id)
  REFERENCES public.category(id)
  ON UPDATE CASCADE
  ON DELETE SET NULL
  NOT VALID;

alter table "public"."post"
  validate constraint "post_category_id_fkey";
