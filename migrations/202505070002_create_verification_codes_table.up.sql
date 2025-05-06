-- Create verification_codes table
CREATE TABLE verification_codes (
    email VARCHAR(255) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
