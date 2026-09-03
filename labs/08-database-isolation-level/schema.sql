-- Schema for Lab 08: Database Isolation Level
-- Wallet accounts for transfer experiments and isolation level demonstrations

DROP TABLE IF EXISTS isolation_invoices CASCADE;
DROP TABLE IF EXISTS isolation_transfer_audit CASCADE;
DROP TABLE IF EXISTS isolation_accounts CASCADE;

CREATE TABLE isolation_accounts (
    id SERIAL PRIMARY KEY,
    owner VARCHAR(100) NOT NULL,
    balance BIGINT NOT NULL CHECK (balance >= 0)
);

CREATE TABLE isolation_transfer_audit (
    id SERIAL PRIMARY KEY,
    from_account_id INT NOT NULL,
    to_account_id INT NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE isolation_invoices (
    id SERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed initial test accounts
-- Alice: 1,000,000
-- Bob:   1,000,000
-- Charlie: 1,000,000
INSERT INTO isolation_accounts (id, owner, balance) VALUES
    (1, 'Alice', 1000000),
    (2, 'Bob', 1000000),
    (3, 'Charlie', 1000000)
ON CONFLICT (id) DO UPDATE SET balance = EXCLUDED.balance;
