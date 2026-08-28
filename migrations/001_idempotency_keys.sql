-- Skema Tabel Idempotency Keys
-- Laboratorium 01: Idempotency

CREATE TABLE IF NOT EXISTS idempotency_keys (
    id VARCHAR(64) PRIMARY KEY,
    key VARCHAR(255) NOT NULL,
    scope VARCHAR(128) NOT NULL,
    request_hash VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL, -- 'processing', 'completed', 'failed'
    response_status INT NOT NULL DEFAULT 0,
    response_body TEXT,
    resource_id VARCHAR(128),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Unique constraint mempertimbangkan scope dan key
    -- Memungkinkan key yang sama ("checkout-123") digunakan oleh tenant/user berbeda
    CONSTRAINT idx_idempotency_scope_key UNIQUE(scope, key)
);

-- Indeks tambahan untuk pencarian cepat dan cleanup TTL
CREATE INDEX IF NOT EXISTS idx_idempotency_expires_at ON idempotency_keys(expires_at);
CREATE INDEX IF NOT EXISTS idx_idempotency_scope ON idempotency_keys(scope);

-- Komentar mengenai Mengapa Global Key Kurang Ideal:
-- 1. Collision antar User/Tenant:
--    Jika key hanya berupa string seperti "order-123", dua user berbeda yang kebetulan
--    menggunakan key yang sama (atau client yang menggunakan generator key sederhana)
--    akan saling berbenturan dan mengakibatkan salah satu request gagal atau me-return
--    data user lain.
-- 2. Multi-Endpoint Namespace:
--    Key "pay-001" untuk endpoint /payments mungkin berbeda konteks dengan "pay-001"
--    untuk endpoint /subscriptions. Scope memisahkan namespace ini secara aman.
-- 3. Security & Isolation:
--    Lingkup multi-tenant atau multi-service memerlukan isolasi agar satu tenant
--    tidak dapat menebak atau menimpa idempotency key tenant lain.
