-- Schema for Lab 05 - Race Condition

CREATE TABLE IF NOT EXISTS inventory_products (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS service_bookings (
    id SERIAL PRIMARY KEY,
    branch_id VARCHAR(50) NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    service_date DATE NOT NULL,
    slot_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Unique constraint protects the booking race invariant
    CONSTRAINT uq_service_bookings_slot UNIQUE (branch_id, service_date, slot_time)
);

CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    invoice_no INT NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_invoice_no UNIQUE (invoice_no)
);
