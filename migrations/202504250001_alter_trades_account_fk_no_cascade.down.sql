-- +migrate Down
-- Reverse the above changes: drop the new constraint, make account_id NOT NULL, and restore ON DELETE CASCADE
ALTER TABLE trades DROP CONSTRAINT IF EXISTS trades_account_id_fkey;
ALTER TABLE trades ALTER COLUMN account_id SET NOT NULL;
ALTER TABLE trades ADD CONSTRAINT trades_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;
