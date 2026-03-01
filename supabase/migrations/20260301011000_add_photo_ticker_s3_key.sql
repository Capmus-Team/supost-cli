alter table public.photo
  add column if not exists ticker_s3_key text;
