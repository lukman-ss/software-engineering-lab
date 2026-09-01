# Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client

## Mental Model Utama

API bukan sekadar endpoint yang menghasilkan HTTP 200. API adalah **kontrak** antara backend dan konsumen. Perubahan backend yang berhasil compile dan mengembalikan HTTP 200 tidak otomatis berarti perubahan tersebut **backward compatible**.

> **HTTP 200 ≠ backward compatible**

Ini adalah yang paling sering saya temui sebagai akar masalah di perusahaan: tim backend menganggap "backend OK" sudah berarti "semua client masih kann". Padahal, perubahan pada representasi data (JSON/XML) sering merusak client yang sudah ada.

## Studi Kasus: Invoice Bengkel CMMS

### Alur Bisnis

Kita punya sistem billing untuk bengkel perawatan mesin (CMMS - Computerized Maintenance Management System). Setiap service mengenerate invoice.

**Contract V1 (versi lama):**

```http
GET /api/v1/invoices/1001
```

Response:

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}
```

**Legacy Android Client** sudah sudah terdeploy ke 5.000 perangkat. Model-nya:

```go
type LegacyInvoice struct {
    ID       int    `json:"id"`
    Customer string `json:"customer"`  // HARUS string
    Total    int64  `json:"total"`
    Status   string `json:"status"`
}
```

### Produk Kemudian Minta Fitur Baru

Tim produk menambah field `customer` menjadi object:

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

**Anda tahu apa yang terjadi?** `customer` berubah dari `string` menjadi `object`. Ini adalah **BREAKING CHANGE**.

Legacy client mencoba `json.Unmarshal` ke `string` → **FAIL**. Tapi server tetap HTTP 200.

---

## Tujuan Laboratorium

1. Reproduce failure (membuktikan breaking change)
2. Pahami root cause (kontrak berubah, tidak ada versioning)
3. Implementasikan solusi (API versioning dengan V1/V2 contract terpisah)
4. Verifikasi dengan test otomatis

---

## Tiga Skenario

### 1. Unsafe: Breaking Change Tanpa Versioning

Backend secara radung mengubah struktur `customer` dari string ke object.

**File:** `unsafe_server.go`, `unsafe_test.go`
**Command:** `go test -v -run TestBreakingChange_LegacyClientFails`

Hasil: Test PASS karena `json.Unmarshal` gagal. Ini membuktikan: **HTTP 200 ≠ backward compatible**.

### 2. Safe: API Versioning dengan V1/V2

Backend menyediakan dua contract berdampingan:
- `GET /api/v1/invoices/1001` → `customer: "Budi"` (string)
- `GET /api/v2/invoices/1001` → `customer: {id: 15, name: "Budi", phone: "08123"}` (object)

**File:** `safe_server.go`, `safe_test.go`
**Command:** `go test -v -run TestSafeVersioning`

Hasil: Kedua client (legacy + new) berhasil. Kontrak dipisah secara eksplisit.

**Arsitektur:**

```
Domain Model (Invoice, Customer)
       │
       ├── mapper → InvoiceV1Response
       │
       └── mapper → InvoiceV2Response
```

### 3. Additive: Penambahan Field Tanpa Breaking Change

Backend menambah field `currency: "IDR"` tanpa mengubah tipe existing.

**File:** `additive_server.go`, `additive_test.go`
**Command:** `go test -v -run TestAdditiveField_LegacyClientStillWorks`

Hasil: Legacy client tetap berhasil decode karena `json.Unmarshal` Go secara default mengabaikan unknown field.

> **Catatan penting:** Penambahan field aman HANYA jika consumer tolerant unknown field. Jika consumer menggunakan strict schema (misal: JSON Schema validation), penambahan field tetap bisa menjadi breaking change.

---

## Cara Menjalankan Test

```bash
# Masuk ke lab
cd labs/06-api-versioning

# Jalankan semua test
go test -v ./...

# Test masing-masing skenario
go test -v -run TestBreakingChange_LegacyClientFails
go test -v -run TestSafeVersioning
go test -v -run TestAdditiveField_LegacyClientStillWorks

# Dengan race detector
go test -race -v ./...
```

---

## Kesimpulan

| Skenario | Perubahan | Breaking Change? | Solusi |
|----------|-----------|------------------|--------|
| Unsafe | `customer: "Budi"` → `customer: {...}` | ✅ YES | API Versioning |
| Safe | V1 = string, V2 = object | ❌ NO | Parallel Contract |
| Additive | Tambah `currency` field | ❌ JIKA tolerant | Unknown field skip |

**Prinsip:** Versioning bukan tentang membuat database baru (`invoice_v2`, `customer_v2`). Versioning adalah tentang **representasi public contract**.

---

## File

```
labs/06-api-versioning/
├── go.mod                    # Module go 1.22
├── domain.go                 # Model Invoice, Customer, LegacyInvoice
├── unsafe_server.go          # Anti-pattern: breaking change
├── unsafe_test.go            # Test: legacy client fails
├── safe_server.go            # V1/V2 contract dengan mapper
├── safe_test.go              # Test: kedua contract berhasil
├── additive_server.go        # Penambahan field currency
├── additive_test.go          # Test: legacy tetap berhasil
└── README.md                 # This file
```