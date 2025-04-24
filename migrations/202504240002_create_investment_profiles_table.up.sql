-- +migrate Up
CREATE TABLE IF NOT EXISTS investment_profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    age INT,
    max_acceptable_short_term_loss_percentage INT,
    expected_annualized_rate_of_return INT,
    time_horizon VARCHAR(255),
    years_investing INT,
    UNIQUE(user_id)
);