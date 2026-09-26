# Case Studies in Backward Compatibility

## Case Study A — Customer Phone (1:1 ke 1:N)

### Initial State
```sql
CREATE TABLE customers (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(20)  -- single phone
);
```

### Requirement
Customer dapat memiliki beberapa nomor telepon.

### Target Schema
```sql
CREATE TABLE customer_phones (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL REFERENCES customers(id),
  phone_number VARCHAR(20) NOT NULL,
  is_primary BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index untuk lookup cepat
CREATE INDEX idx_customer_phones_customer_id ON customer_phones(customer_id);
```

### Strategy: Expand → Migrate → Contract

#### Phase 1: Expand
```sql
-- 1. Buat tabel baru (non-breaking)
-- 2. JANGAN sentuh customers.phone
-- 3. Deploy aplikasi yang compatible dengan keduanya
```

**Application Code (Expand)**:
```go
// Dual write: tulis ke phone lama DAN ke tabel baru
func (s *Store) UpdatePhone(ctx context.Context, customerID int64, phone string) error {
    tx := s.db.Begin()
    
    // Legacy write (untuk V1 consumers)
    tx.Exec("UPDATE customers SET phone = $1 WHERE id = $2", phone, customerID)
    
    // New write
    tx.Exec(`
        INSERT INTO customer_phones (customer_id, phone_number, is_primary)
        VALUES ($1, $2, TRUE)
        ON CONFLICT (customer_id, phone_number) DO UPDATE SET is_primary = TRUE
    `, customerID, phone)
    
    return tx.Commit()
}
```

#### Phase 2: Migrate

**Backfill Script** (idempotent, resumable):
```go
func BackfillPhones(ctx context.Context, db *sql.DB, batchSize int) error {
    lastID := int64(0)
    
    for {
        rows, err := db.QueryContext(ctx, `
            SELECT id, phone FROM customers 
            WHERE phone IS NOT NULL 
            AND id > $1 
            ORDER BY id 
            LIMIT $2
        `, lastID, batchSize)
        if err != nil { return err }
        
        count := 0
        for rows.Next() {
            var id int64
            var phone string
            rows.Scan(&id, &phone)
            
            // Idempotent insert
            _, err = db.ExecContext(ctx, `
                INSERT INTO customer_phones (customer_id, phone_number, is_primary)
                VALUES ($1, $2, TRUE)
                ON CONFLICT (customer_id, phone_number) DO NOTHING
            `, id, phone)
            if err != nil { return err }
            
            lastID = id
            count++
        }
        rows.Close()
        
        if count == 0 { break }
        time.Sleep(100 * time.Millisecond) // throttling
    }
    return nil
}
```

**Fallback Read**:
```go
func (s *Store) GetPhones(ctx context.Context, customerID int64) ([]string, error) {
    // Coba tabel baru
    rows, _ := s.db.QueryContext(ctx, `
        SELECT phone_number FROM customer_phones 
        WHERE customer_id = $1 ORDER BY is_primary DESC
    `, customerID)
    
    var phones []string
    for rows.Next() {
        var p string
        rows.Scan(&p)
        phones = append(phones, p)
    }
    rows.Close()
    
    // Fallback ke legacy jika tabel baru kosong
    if len(phones) == 0 {
        var legacyPhone sql.NullString
        s.db.QueryRowContext(ctx, `SELECT phone FROM customers WHERE id = $1`, customerID).
            Scan(&legacyPhone)
        if legacyPhone.Valid {
            phones = []string{legacyPhone.String}
        }
    }
    return phones, nil
}
```

#### Phase 3: Contract

```sql
-- 1. Stop write ke customers.phone (deploy flag-controlled)
-- 2. Verifikasi metrics legacy_read = 0, legacy_write = 0
-- 3. Drop column
ALTER TABLE customers DROP COLUMN phone;
```

### API Compatibility

**Legacy Response** (V1):
```json
{
  "id": 10,
  "name": "Budi",
  "phone": "+628111"  // first phone
}
```

**Modern Response** (V2):
```json
{
  "id": 10,
  "name": "Budi",
  "phones": ["+628111", "+628222"],
  "primary_phone": "+628111"
}
```

**Transformation Layer** (Stripe-style):
```go
func transformCustomer(c *Customer, apiVersion string) map[string]any {
    if apiVersion == "v1" {
        return map[string]any{
            "id": c.ID,
            "name": c.Name,
            "phone": c.PrimaryPhone(),  // flatten ke single field
        }
    }
    return map[string]any{
        "id":           c.ID,
        "name":         c.Name,
        "phones":       c.Phones,
        "primary_phone": c.PrimaryPhone(),
    }
}
```

---

## Case Study B — CMMS Invoice Mechanics (1:1 ke N:M)

### Initial State
```sql
CREATE TABLE invoices (
  id BIGSERIAL PRIMARY KEY,
  mechanic_id BIGINT REFERENCES mechanics(id),
  amount NUMERIC(12,2),
  status VARCHAR(20)
);
```

### Requirement
Satu invoice dapat memiliki beberapa mechanic (N:M relationship).

### Target Schema
```sql
CREATE TABLE invoice_mechanics (
  invoice_id BIGINT NOT NULL REFERENCES invoices(id),
  mechanic_id BIGINT NOT NULL REFERENCES mechanics(id),
  role VARCHAR(50) DEFAULT 'primary',
  PRIMARY KEY (invoice_id, mechanic_id)
);
```

### Strategy

#### Expand
```sql
CREATE TABLE invoice_mechanics ( ... );
-- JANGAN hapus invoices.mechanic_id
```

#### Migrate (Code)
```go
func (s *Store) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error) {
    tx := s.db.Begin()
    
    // Legacy: single mechanic_id (first mechanic)
    firstMechanic := req.Mechanics[0]
    tx.Exec(`
        INSERT INTO invoices (mechanic_id, amount, status)
        VALUES ($1, $2, $3) RETURNING id
    `, firstMechanic.ID, req.Amount, req.Status)
    
    // New: semua mechanics
    for _, m := range req.Mechanics {
        tx.Exec(`
            INSERT INTO invoice_mechanics (invoice_id, mechanic_id, role)
            VALUES ($1, $2, $3)
        `, invoiceID, m.ID, m.Role)
    }
    
    return tx.Commit()
}
```

#### Backfill
```go
func BackfillInvoiceMechanics(ctx context.Context, db *sql.DB) error {
    _, err := db.ExecContext(ctx, `
        INSERT INTO invoice_mechanics (invoice_id, mechanic_id, role)
        SELECT id, mechanic_id, 'primary' 
        FROM invoices 
        WHERE mechanic_id IS NOT NULL
        ON CONFLICT (invoice_id, mechanic_id) DO NOTHING
    `)
    return err
}
```

#### API Compatibility Layer

**Legacy Endpoint**: `GET /api/v1/invoices/{id}`
```json
{
  "id": 100,
  "mechanic_id": 42,     // first mechanic
  "amount": 500.00
}
```

**Modern Endpoint**: `GET /api/v2/invoices/{id}`
```json
{
  "id": 100,
  "mechanics": [
    {"id": 42, "name": "John", "role": "primary"},
    {"id": 43, "name": "Jane", "role": "secondary"}
  ],
  "amount": 500.00
}
```

#### Contract
1. Monitor akses `GET /api/v1/invoices` dan `invoices.mechanic_id` read
2. Jika 0 traffic legacy → drop `invoices.mechanic_id`

---

## Case Study C — Multi Currency (Schema Splitting)

### Initial State
```sql
CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  price NUMERIC(12,2) NOT NULL  -- assumed USD
);
```

### Requirement
Product memiliki harga dalam multiple currencies.

### Target Schema
```sql
CREATE TABLE product_prices (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT NOT NULL REFERENCES products(id),
  currency CHAR(3) NOT NULL,  -- ISO 4217
  amount NUMERIC(12,2) NOT NULL,
  UNIQUE (product_id, currency)
);
```

### Strategy

#### Expand
```sql
CREATE TABLE product_prices ( ... );
-- products.price tetap ada untuk legacy
```

#### Migrate (Data)
```go
func BackfillPrices(ctx context.Context, db *sql.DB) error {
    // Asumsikan existing prices dalam USD
    _, err := db.ExecContext(ctx, `
        INSERT INTO product_prices (product_id, currency, amount)
        SELECT id, 'USD', price 
        FROM products 
        WHERE price IS NOT NULL
        ON CONFLICT (product_id, currency) DO UPDATE SET amount = EXCLUDED.amount
    `)
    return err
}
```

#### Migrate (API Compatibility)

```go
func GetProduct(ctx context.Context, id int64, currency string) (*ProductResponse, error) {
    // Coba product_prices
    var price ProductPrice
    err := db.QueryRowContext(ctx, `
        SELECT currency, amount FROM product_prices 
        WHERE product_id = $1 AND currency = $2
    `, id, currency).Scan(&price.Currency, &price.Amount)
    
    if err == sql.ErrNoRows && currency == "USD" {
        // Fallback ke legacy price column
        var legacyPrice sql.NullFloat64
        db.QueryRowContext(ctx, `SELECT price FROM products WHERE id = $1`, id).
            Scan(&legacyPrice)
        if legacyPrice.Valid {
            price = ProductPrice{Currency: "USD", Amount: legacyPrice.Float64}
        }
    }
    // ... return response
}
```

**Legacy Response**:
```json
{
  "id": 1,
  "name": "Widget",
  "price": 99.99  // USD
}
```

**Modern Response**:
```json
{
  "id": 1,
  "name": "Widget",
  "prices": [
    {"currency": "USD", "amount": 99.99},
    {"currency": "EUR", "amount": 89.99},
    {"currency": "IDR", "amount": 1500000.00}
  ]
}
```

#### Contract

**Perhatian Khusus**: `products.price` mungkin tidak bisa dihapus jika:
- Mobile app lama tidak bisa diupdate (user tidak update app)
- Partner API eksternal tidak bisa bermigrasi

**Solusi**:
- Pertahankan `products.price` sebagai **computed/view column** atau sync via trigger
- API layer selalu flatten ke `.price` untuk V1 consumers (melalui version transformer)
- Hanya hapus physical column jika **semua** internal services bermigrasi

---

## Ringkasan Cross-Case Patterns

| Aspek | Case A (Phone) | Case B (Invoice Mechanics) | Case C (Multi Currency) |
|-------|----------------|---------------------------|------------------------|
| **Expansion** | Tabel baru `customer_phones` | Tabel baru `invoice_mechanics` | Tabel baru `product_prices` |
| **Legacy Field** | `customers.phone` | `invoices.mechanic_id` | `products.price` |
| **Dual Write** | Ya (legacy + new) | Ya (legacy + new) | Ya (legacy + new) |
| **Backfill** | Synchronous batch | Single INSERT | Single INSERT |
| **Fallback Read** | Perlu (sampai backfill done) | Tidak (backfill atomic) | Perlu (currency tidak USD) |
| **API Transform** | Array → single field | Array → single field | Multi → single field |
| **Contract Trigger** | Metrics = 0 | Metrics = 0 | **Tidak pernah** (mobile constraint) |

---

## Catatan Kualitas Bukti

| Pernyataan | Sumber | Confidence | Keterangan |
|------------|--------|------------|------------|
| Pola Expand-Migrate-Contract berlaku | Martin Fowler | HIGH | Framework utama |
| Dual write pattern | Martin Fowler Parallel Change | HIGH | Described in migrate phase |
| Fallback read pattern | Fowler (implisit) | MEDIUM | Not explicit code example |
| Legacy API tidak bisa dihapus (mobile) | GitHub API policy | HIGH | GitHub 24-month support window |