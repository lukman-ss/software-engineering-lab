-- Lab 02: Database Index Foundation
-- Schema definition for service table

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

-- Indexes for the lab
-- Primary business query index (to be tested)
CREATE INDEX IF NOT EXISTS idx_service_branch_status_date
    ON service (branch_id, status, service_date DESC);

-- Index for status filtering only
CREATE INDEX IF NOT EXISTS idx_service_status
    ON service (status);

-- Index for service_date filtering only
CREATE INDEX IF NOT EXISTS idx_service_date
    ON service (service_date);

-- Foreign key constraints (optional, for referential integrity)
-- ALTER TABLE service ADD CONSTRAINT fk_branch FOREIGN KEY (branch_id) REFERENCES branch(id);
-- ALTER TABLE service ADD CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customer(id);
-- ALTER TABLE service ADD CONSTRAINT fk_mechanic FOREIGN KEY (mechanic_id) REFERENCES mechanic(id);