-- +migrate Up
CREATE TABLE trades (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('buy', 'sell')),
    asset_type VARCHAR(10) NOT NULL CHECK (asset_type IN ('stock', 'crypto')),
    ticker VARCHAR(20) NOT NULL,
    trade_date DATE NOT NULL,
    quantity NUMERIC NOT NULL,
    price NUMERIC NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);