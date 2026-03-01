# 0049 - post create v2 photo uploads

## Summary

- Added repeatable `--photo` flag to `supost post create` (up to 4 files).
- Extended create-submit flow to carry photo bytes through validation and submit logic.
- Added `PostCreatePhotoUploader` interface in service layer.
- Implemented S3 uploader adapter that uploads to:
  - `v2/posts/{post_id}/{random_hex}.{ext}`
  - preserves file extension when valid; otherwise defaults to `.jpg`
  - generates two variants per photo:
    - post image: max width 340px
    - ticker image: max width 220px
- Persists uploaded photo metadata to `public.photo` via repository `SavePostPhotos`.
- Added config keys:
  - `S3_PHOTO_BUCKET`
  - `S3_PHOTO_PREFIX`
  - `S3_PHOTO_REGION`
  - `S3_PHOTO_AWS_PROFILE`
- Updated `post` read/search/home query `has_image` logic to include `public.photo` existence.
- Updated post page rendering to use S3-key-based URLs when present.

## Notes

- Adapter uses AWS CLI (`aws s3 cp`) with existing AWS credential env vars/profile.
- Dry run mode (`--dry-run`) never uploads or persists photos.
