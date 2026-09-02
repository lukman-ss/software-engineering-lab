# Senior Software Engineer Daily #6

## Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client

## Mental Model Utama

API bukan sekadar endpoint yang menghasilkan HTTP 200. API adalah **kontrak** antara backend dan konsumen. Perubahan backend yang berhasil compile dan mengembalikan HTTP 200 tidak otomatis berarti perubahan tersebut **backward compatible**.

> **HTTP 200 ≠ backward compatible**

Tim backend sering menganggap "backend OK" berarti "semua client masih dapat". Padahal, perubahan pada representasi data (JSON/XML) sering merusak client yang sudah ada.

---

## Problem: Customer String → Customer Object

```json
// API Lama
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}

// Breaking Change
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

## Breaking Changes

Berikut contoh perubahan yang menjadi **breaking change**:

### Core Changes

| Perubahan | Contoh |
|-----------|--------|
| Rename field | `name` → `full_name` |
| Remove field | hapus `phone` |
| Primitive type berubah | `price:number` → `price:string` |
| string ↔ object | `customer: "Budi"` ↔ `customer: {id: 15, name: "Budi"}` |
| object ↔ array | `customer: {...}` ↔ `customer: [...]` |
| Tanggal format berubah | `2026-07-25` ↔ `25/07/2026` |
| Nullability berubah | Request: optional → required (breaking)<br>Response: non-nullable → nullable (bisa breaking) |

### Semantic Changes

| Perubahan | Contoh |
|-----------|--------|
| Enum semantics berubah | `status: "ACTIVE"` berarti hal lain |
| HTTP status semantics berubah | `200 OK` ↔ `201 Created` untuk update |
| Required request input berubah | field `id` menjadi required |
| Pagination contract berubah | `page` → `page_number`, `limit` → `size` |
| Business meaning endpoint berubah | `/invoices` sekarang mengembalikan draft juga |
| Authentication requirement berubah | Endpoint butuh API key baru |
| Error response shape berubah | `{error: "msg"}` ↔ `{code: "ERR", message: "msg"}` |

---

## Usually Compatible vs Breaking

### Backward Compatible

Perubahan yang **biasanya compatible**:

| Perubahan | Kondisi |
|-----------|---------|
| Tambah field optional di response | consumer toleran unknown field |
| Tambah endpoint | tidak memengaruhi existing endpoint |
| Tambah query parameter optional | client tidak wajib pakai |

>**Catatan tentang Penambahan Field:**

>Menambah optional response field biasanya backward-compatible untuk tolerant readers. Namun **Tidak selalu aman!** Bergantung pada behavior consumer:
>- `json.Unmarshal` Go: abaikan unknown field → aman
>- JSON Schema strict: field baru bisa mismatch → tidak aman

### Breaking

Perubahan yang **pasti breaking**:

| Perubahan | Alasan |
|-----------|--------|
| Rename | Consumer pakai nama lama |
| Remove | Consumer butuh data itu |
| Type change | Deserialization error |
| Shape change | Struct mismatch |
| Semantics change | Logika bisnis berubah |

---

## API Change Decision Rule

**Backward-compatible change** → **pertahankan version yang sama**  
**Contract-breaking change** → **buat major API version baru**

---

## Versioning Strategies

### URL Versioning (Pendekatan Lab Ini)

```http
GET /api/v1/invoices/1001
GET /api/v2/invoices/1001
```

Dipilih untuk lab karena:
- **Simple** — tidak perlu middleware tambahan
- **Visible** — langsung terlihat di logs
- **Debuggable** — easy to test with curl
- **Easy to document** — clear URL pattern
- **Easy to route** — `http.ServeMux` handles ini

>**Catatan penting:** Versi di URL **tidak** menentukan apakah endpoint idempotent. Idempotency ditentukan oleh semantics HTTP (GET, PUT, DELETE) dan business operation, bukan lokasi version identifier.

### Header Versioning (Alternatif)

```http
GET /api/invoices/1001
Accept: application/vnd.company.v2+json
```

atau:

```http
GET /api/invoices/1001
API-Version: 2
```

Keuntungan:
- URL resource tetap stabil (memungkinkan caching & routing yang konsisten)
- Memungkinkan multiple versions dalam satu endpoint

Kerugian:
- Kurang visibilitas di logs/network traces
- Membutuhkan middleware parsing header

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
| Support period? | contoh: 90 hari |
| Monitor old-version traffic? | Ya |
| Compatibility test? | Automated di CI |
| Deprecation communication? | Release notes |
| Sunset criteria? | < 5% V1 traffic |

> **API versioning tanpa mengetahui consumer lama tetap berbahaya.** Jangan buat V2 hanya karena "mau bangun".

---

## Android Migration Strategy

```
Time:      0d           30d          60d          90d          120d
           │             │            │            │            │
           ▼             ▼            ▼            ▼            ▼
Day 0: Release V2
       Monitor      Send         Deprecate      Sunset
       adoption     warning      V1 endpoint
```

> **Android release terbaru ≠ semua user sudah upgrade.** Proses upgrade memakan waktu.

Strategi:
1. Keep V1 stable. (Jangan pernah ubah V1)
2. Deploy V2.
3. New Android uses V2.
4. Old Android remains on V1.
5. Monitor adoption.
6. Monitor V1 traffic.
7. Announce deprecation.
8. Migrate remaining consumers.
9. Sunset V1 after criteria are met.

---

## Deprecation Lifecycle

```
V1 active
    ↓
V2 released
    ↓
V1 deprecated
    ↓
Migration monitoring
    ↓
V1 traffic monitoring
    ↓
Sunset criteria reached
    ↓
V1 removed
```

>**INGAT:** 90 hari hanyalah **CONTOH**. Actual timeline bergantung pada:
>- **SLA** — berapa lama service harus support?
>- **Partner contract** — berapa lama perjanjian komersial?
>- **Mobile adoption** — seberapa cepat user upgrade?
>- **Business criticality** — apa konsekuensi V1 down?
>- **Security issue** — ada CVE yang harus patched cepat?
>- **Remaining traffic** — berapa persen V1?

---

## Contract Tests

### V1 Contract Protection

Memverifikasi **consumer lama tetap dapat memahami contract**:

```go
func TestV1Contract_RemainsBackwardCompatible(t *testing.T)
```

Checking secara eksplisit:
- `id` = number, id == 1001
- `customer` = string, customer == "Budi"
- `total` = number, total == 500000
- `status` = string, status == "PAID"

> `HTTP 200` + `valid JSON` ≠ `backward compatible`  
> `LegacyInvoice.decode → error` adalah **bukti breaking change**, bukan bug.

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

**Prinsip:** Domain model ≠ public API contract. DTO terpisah untuk setiap version.

---

## Additive Changes

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID",
  "currency": "IDR"
}
```

Legacy Go client tanpa field `currency` tetap bisa decode karena `json.Unmarshal` mengabaikan field tak dikenal.

---

## API Versioning ≠ Database Versioning

**JANGAN** membuat tabel:

```sql
customers_v2
invoices_v2
products_v2
```

Hanya karena API response berubah! API versioning adalah **representasi public contract**, bukan data model.

---

## Common Mistakes (Senior Engineer Checklist)

| Mistake | Dampak | Solusi |
|---------|--------|--------|
| `name` langsung ganti `full_name` | Client crash | Gunakan `name` + `full_name` paralel |
| Hapus field tanpa notice | Consumer error | Deprecation cycle dulu |
| `price:number` jadi `price:string` | Type error | Version baru |
| Reuse DTO V1/V2 | Keduanya saling terikat | DTO terpisah |
| Tidak ada contract test | Breaking lompat | Tambahkan test |
| Tidak ada consumer inventory | Upgrade tak terduga | Survei dulu |
| Hapus V1 terlalu cepat | Production outage | Monitor traffic |
| Versioning semua perubahan kecil | VERSI EXPLOSION | Hanya major changes |
| `HTTP 200 == compatible` | Assumption salah | Test dari perspective client |
| Tabel `customers_v2` | Database bloat | API contract versioning |

---

## HTTP Routing Notes

Untuk lab ini, routing menggunakan `http.NewServeMux()`:

```go
mux.HandleFunc("/api/v1/invoices/", V1Handler)
mux.HandleFunc("/api/v2/invoices/", V2Handler)
```

- Route `/api/v1/invoices/1001` dipetakan ke `V1Handler`
- Route tak dikenal seperti `/foo` akan mengembalikan 404
- Route dengan ID tidak valid (`/api/v1/invoices/abc`) mengembalikan 400

---

## Exercise

```json
// API Lama
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}

// Breaking Change
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

**Pertanyaan:**
1. Breaking change?
2. Perlu V2?
3. Bagaimana Android lama tetap hidup?
4. Bagaimana membuktikan compatibility?

**Expected Reasoning:**
1. **Ya, breaking change** karena tipe data `customer` berubah dari `string` menjadi `object`.
2. **Ya, perlu V2** karena ini merusak contract.
3. **Android lama tetap hidup** karena endpoint V1 tetap dilayani backend dengan schema string. Routing memisah request.
4. **Membuktikan compatibility** dengan Contract Regression Test.

---

## Key Takeaways

- **API adalah kontrak**, bukan teknologi. Perubahan harus dievaluasi dari perspective consumer.
- **HTTP 200 ≠ compatible**. Selalu test decoding di sisi client.
- **Versioning bukan untuk semua**. Hanya untuk breaking changes.
- **V1 selalu stabil**. Setelah di-release, jangan pernah ubah lagi.
- **Contract test wajib**. Prevent breaking change yang terlewat.
- **Consumer inventory penting**. Tanpa tahu siapa yang pakai, tidak ada cara upgrade yang aman.

---

## Cara Menjalankan Test

```bash
# Masuk ke lab
cd labs/06-api-versioning

# Jalankan semua test
go test -v ./...

# Dengan race detector
go test -race -v ./...
```

---

## File Structure

```
labs/06-api-versioning/
├── go.mod                  # Module go 1.22
├── domain.go               # Model Invoice, Customer
├── ParseLegacyInvoice()    # Helper untuk test legacy decode
├── unsafe_server.go        # Anti-pattern: breaking change
├── unsafe_test.go          # Test: legacy client fails
├── safe_server.go          # V1/V2 DTO + mappers
├── safe_test.go            # Contract regression tests
├── additive_server.go      # Field tambahan currency
├── additive_test.go        # Test: legacy tetap berhasil
└── README.md               # Panduan lengkap
```

---

## Navigasi

- **Previous**: [Lab 05 — Race Condition](../05-race-condition/)
- **Next**: [Lab 07 — Outbox Pattern](../07-outbox-pattern/)