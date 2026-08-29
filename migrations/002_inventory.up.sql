-- Inventory products table for race condition demonstration
CREATE TABLE IF NOT EXISTS inventory_products (
    id VARCHAR(50) PRIMARY KEY,
    name TEXT NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_inventory_products_stock ON inventory_products(stock);