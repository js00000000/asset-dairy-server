-- +migrate Down
ALTER TABLE investment_profiles
    DROP COLUMN monthly_cash_flow,
    DROP COLUMN default_currency; 