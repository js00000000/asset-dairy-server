-- +migrate Down
ALTER TABLE investment_profiles
    ALTER COLUMN monthly_cash_flow DROP DEFAULT,
    ALTER COLUMN default_currency DROP DEFAULT; 