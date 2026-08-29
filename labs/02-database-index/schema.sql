-- Lab 02: Database Index Foundation
-- Schema definition for service table
-- Baseline has table, primary key, unique constraints, but NO query-supporting secondary indexes.

CREATE TABLE service (
    id           BIGSERIAL PRIMARY KEY,
    branch_id    INTEGER NOT NULL,
    customer_id  INTEGER NOT NULL,
    mechanic_id  INTEGER NOT NULL,
    status       VARCHAR(20) NOT NULL,
    service_date DATE NOT NULL,
    invoice_no   VARCHAR(50) UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
