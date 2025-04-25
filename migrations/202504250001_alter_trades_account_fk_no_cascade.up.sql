-- +migrate Up
-- Step 1: Drop the existing foreign key constraint on account_id
ALTER TABLE trades DROP CONSTRAINT IF EXISTS trades_account_id_fkey;

-- Step 2: Make account_id nullable
ALTER TABLE trades ALTER COLUMN account_id DROP NOT NULL;

-- Step 3: Add a new foreign key constraint with ON DELETE SET NULL
ALTER TABLE trades ADD CONSTRAINT trades_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL;