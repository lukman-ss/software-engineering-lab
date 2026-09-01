# Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client

## Mental Model Utama

API bukan sekadar endpoint yang menghasilkan HTTP 200. API adalah **kontrak** antara backend dan konsumen. Perubahan backend yang berhasil compile dan mengembalikan HTTP 200 tidak otomatis berarti perubahan tersebut **backward compatible**.

> **HTTP 200 ≠ backward compatible**

Ini adalah yang paling sering saya temui sebagai akar masalah di perusahaan: tim backend menganggap "backend OK" sudah berarti "semua client masih kann". Padahal, perubahan pada representasi data (JSON/XML) sering merusak client yang sudah ada.

---

## Problem: Customer String → Customer Object

### API Lama

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}
```

### Developer Melakukan Refactor

```json
{
  "id": 1001,
  "customer": {
    "id": 15,
    "name": "Budi",
    "phone": "08123"
  },
  "total": 500000,
  "status": "PAID"
}
```

**Perubahan:** `customer` berubah dari `string` menjadi `object`. Client yang sudah ada gagal parse.

---

## Breaking Change

Berikut adalah contoh-perubahan yang menjadi breaking change:

| Perubahan | Contoh |
|-----------|--------|
| Rename field | `name` → `full_name` |
| Remove field | hapus `phone` |
| Change type | `"500000"` (string) → `500000` (number) |
| Change object shape | `customer.name` → `customer.fullName` |
| Change date format | `"2024-01-01"` → `"01/01/2024"` |
| Change nullability | field wajib → optional (atau sebaliknya) |
| Change enum semantics | status `"ACTIVE"` berarti hal lain |
| Change pagination | `page` → `page_number` |
| Change endpoint semantics | `/invoices` sekarang mengembalikan draft juga |

---

## Usually Compatible Changes

Namun, perubahan berikut sering **backward compatible**:

| Perubahan | Kondisi |
|-----------|---------|
| Adding optional response field | consumer tolerant unknown field |
| Adding endpoint | tidak memengaruhi existing endpoint |
| Adding optional query parameter | client tidak wajib pakai |

**Penting:** Penambahan field **SELALU aman**? Tidak. Bergantung pada behavior consumer:
- `json.Unmarshal` Go: abaikan unknown field → aman
- JSON Schema strict: field baru bisa mismatch → tidak aman

---

## Unsafe Approach: Breaking Change Tanpa Versioning

### Implementasi

File: `unsafe_server.go`, `unsafe_test.go`

```go
// POST /invoices/1001
// Response: customer berubah jadi object
```

### Masalah

```
HTTP 200
├── Valid JSON ✓
├── Backend OK ✓
└── Legacy client FAIL ✗ (json: cannot unmarshal object into string)
```

### Output Test

```
✅ Breaking change confirmed: legacy client fails to decode:
   json: cannot unmarshal object into Go struct field LegacyInvoice.customer of type string
```

---

## Safe Approach: API Versioning

### Dual Contracts

```
GET /api/v1/invoices/1001 → customer: string
GET /api/v2/invoices/1001 → customer: object
```

### Arsitektur

```
Domain Model (Invoice, Customer)
       │
       ├── mapper → InvoiceV1Response (customer: string)
       │
       └── mapper → InvoiceV2Response (customer: object)
```

### File

- `domain.go`: Domain model + `ParseLegacyInvoice()` helper
- `safe_server.go`: V1/V2 DTOs dan handler
- `safe_test.go`: Contract regression test

### Contract Tests

**V1 Contract Protection:**
```go
func TestV1Contract_RemainsBackwardCompatible(t *testing.T)
```

Menggunakan `json.RawMessage` untuk semantic assertion (bukan string comparison).

**V2 Contract Verification:**
```go
func TestV2Contract_UsesNestedCustomer(t *testing.T)
```

Memastikan `customer` adalah object dengan field yang benar.

---

## API Versioning ≠ Database Versioning

JANGAN membuat tabel:

```sql
customers_v2
invoices_v2
products_v2
```

Hanya karena API response berubah! API versioning adalah **representasi public contract**, bukan data model.

---

## URL Versioning vs Header Versioning

### URL Versioning (Pendekatan Lab Ini)

```http
GET /api/v1/invoices/1001
GET /api/v2/invoices/1001
```

Pragmatis, mudah debug, cocok untuk lab.

### Header Versioning (Alternatif)

```http
GET /api/invoices/1001
Accept: application/vnd.company.v1+json
Accept: application/vnd.company.v2+json
```

Lebih idempoten (URL tetap), tapi kurang visibilitas.

---

## Consumer Inventory

Sebelum melakukan breaking change, tanyakan:

| Pertanyaan | Tujuan |
|------------|--------|
| Siapa yang konsumsi endpoint ini? | Android, iOS, Web, Partner, BI |
| Contract lama apa? | `customer: string` |
| Additive atau breaking? | Type change = breaking |
| Field/type/shape berubah? | Ya, customer Type |
| Semua client bisa deploy bersamaan? | Tidak, ada 5000 Android lama |
| Perlu major version? | Ya, V1 → V2 |
| Migration strategy? | Keep V1 + rilis V2 |
| Support period? | 90 hari |
| Monitor old-version traffic? | Ya |
| Compatibility test? | Automated di CI |
| Deprecation communication? | Release notes |
| Sunset criteria? | < 5% V1 traffic |

---

## Android Migration Strategy

```
Time: 0d           30d          60d          90d          120d
      │             │            │            │            │
      ▼             ▼            ▼            ▼            ▼
Release V2   Monitor     Send      Deprecate     Sunset
             adoption    warning   V1 endpoint
```

**Kunci:**
- `android-release` ≠ semua user upgrade
- Keep V1 aktif sampai adopsi V2 cukup tinggi
- Monitoring traffic V1 vs V2

---

## Senior Checklist: API Change Review

Sebelum deploy perubahan API:

1. ✅ Siapa consumer?
2. ✅ Contract lama apa?
3. ✅ Additive atau breaking?
4. ✅ Field/type/shape berubah?
5. ✅ Semua client deploy bersamaan?
6. ✅ Perlu major version?
7. ✅ Migration strategy?
8. ✅ Support period?
9. ✅ Monitor old-version traffic?
10. ✅ Compatibility test?
11. ✅ Deprecation communication?
12. ✅ Sunset criteria?

---

## Closing Mindset

**Junior:** "Endpoint baru bekerja, sudah."

**Senior:** "Apakah perubahan contract ini tetap aman untuk consumer yang sudah berjalan di production?"

---

## Cara Menjalankan Test

```bash
# Masuk ke lab
cd labs/06-api-versioning

# Jalankan semua test
go test -v ./...

# Test masing-masing skenario
go test -v -run TestBreakingChange_LegacyClientFails
go test -v -run TestV1Contract
go test -v -run TestV2Contract
go test -v -run TestAdditiveField_LegacyClientStillWorks

# Dengan race detector
go test -race -v ./...
```

---

## File

```
labs/06-api-versioning/
├── go.mod                  # Module go 1.22
├── domain.go               # Model Invoice, Customer, LegacyInvoice, ParseLegacyInvoice
├── unsafe_server.go        # Anti-pattern: breaking change
├── unsafe_test.go          # Test: legacy client fails
├── safe_server.go          # V1/V2 contract dengan mapper
├── safe_test.go            # Test: contract regression
├── additive_server.go      # Penambahan field currency
├── additive_test.go        # Test: legacy tetap berhasil
└── README.md               # This file
```

---

## Navigasi

- **Previous**: [Lab 05 — Race Condition](../05-race-condition/)
- **Next**: [Lab 07 — Outbox Pattern](../07-outbox-pattern/)