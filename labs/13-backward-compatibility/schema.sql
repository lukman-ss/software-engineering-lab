-- V1: Baseline (1:1 relationship)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL
);

-- V2: Expand (1:N relationship added without touching old column)
CREATE TABLE user_phones (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_user_phones_user_id ON user_phones(user_id);

-- V3: Contract (executed ONLY after legacy consumer and dual-write are fully phased out)
-- ALTER TABLE users DROP COLUMN phone;
