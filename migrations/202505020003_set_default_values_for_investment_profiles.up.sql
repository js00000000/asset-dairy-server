-- +migrate Up
ALTER TABLE investment_profiles
    ALTER COLUMN monthly_cash_flow SET DEFAULT 0,
    ALTER COLUMN default_currency SET DEFAULT 'USD';

-- Update existing rows to set default values
UPDATE investment_profiles
SET monthly_cash_flow = 0,
    default_currency = 'USD'
WHERE monthly_cash_flow IS NULL
   OR default_currency IS NULL; 