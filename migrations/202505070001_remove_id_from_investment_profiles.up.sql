BEGIN;

-- Remove the existing primary key constraint on id column
ALTER TABLE investment_profiles DROP CONSTRAINT IF EXISTS investment_profiles_pkey;

-- Drop the id column
ALTER TABLE investment_profiles DROP COLUMN IF EXISTS id;

-- Set user_id as the primary key
ALTER TABLE investment_profiles ADD PRIMARY KEY (user_id);

COMMIT;
