-- Migration 002: Create listings table
-- Maps to: internal/domain/listing.go → Listing struct

CREATE TABLE IF NOT EXISTS listings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price       INTEGER NOT NULL DEFAULT 0,
    category    TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_listings_user_id ON listings(user_id);
CREATE INDEX IF NOT EXISTS idx_listings_status ON listings(status);
CREATE INDEX IF NOT EXISTS idx_listings_category ON listings(category);
CREATE INDEX IF NOT EXISTS idx_listings_created_at ON listings(created_at DESC);

-- TODO: Enable RLS when deploying to Supabase
-- ALTER TABLE listings ENABLE ROW LEVEL SECURITY;
-- CREATE POLICY "Anyone can view active listings" ON listings
--     FOR SELECT USING (status = 'active');
-- CREATE POLICY "Users can manage own listings" ON listings
--     FOR ALL USING (auth.uid() = user_id);
