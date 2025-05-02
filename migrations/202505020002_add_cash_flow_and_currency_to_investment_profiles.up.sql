-- +migrate Up
ALTER TABLE investment_profiles
    ADD COLUMN monthly_cash_flow DECIMAL(19,4),
    ADD COLUMN default_currency VARCHAR(10); 