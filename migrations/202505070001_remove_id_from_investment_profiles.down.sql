BEGIN;

-- Recreate the id column with UUID type
ALTER TABLE investment_profiles ADD COLUMN id UUID PRIMARY KEY;

-- Remove the primary key constraint on user_id
ALTER TABLE investment_profiles DROP CONSTRAINT IF EXISTS investment_profiles_pkey;

COMMIT;
