# Senior Software Engineer Daily #6

## Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client

## Mental Model Utama

API bukan sekadar endpoint yang menghasilkan HTTP 200. API adalah **kontrak** antara backend dan konsumen. Perubahan backend yang berhasil compile dan mengembalikan HTTP 200 tidak otomatis berarti perubahan tersebut **backward compatible**.

> **HTTP 200 ≠ backward compatible**

Tim backend sering menganggap "backend OK" berarti "semua client masih dapat". Padahal, perubahan pada representasi data (JSON/XML) sering merusak client yang sudah ada.

---

## Problem: Customer String → Customer Object

API lama:

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}
```

Breaking change:

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

## Breaking Changes

Berikut contoh perubahan yang menjadi **breaking change**:

### Core Changes

| Perubahan | Contoh |
|-----------|--------|
| Rename field | `name` → `full_name` |
| Remove field | hapus `phone` (breaking karena contract berubah) |
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
| Required request header berubah | Endpoint tidak butuh `X-Tenant-ID`, kemudian header menjadi wajib. Client lama tanpa header akan gagal. |

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

### Contract-Breaking Changes

Perubahan berikut mengubah published contract dan dapat memutus consumer yang bergantung pada contract lama.

| Perubahan | Reasoning |
|-----------|-----------|
| Rename | nama field lama tidak lagi tersedia dalam published contract |
| Remove | published contract berubah; consumer yang bergantung pada field tersebut dapat gagal |
| Type change | schema lama tidak lagi cocok dan deserialization dapat gagal |
| Shape change | representasi lama tidak lagi sesuai contract |
| Semantic change | payload dapat tetap valid secara syntax tetapi behavior consumer dapat menjadi salah |

---

## API Change Decision Rule

**Backward-compatible change** → **pertahankan version yang sama**  
**Contract-breaking change** → **gunakan major API version baru atau compatibility/migration layer**

> **Catatan:** Compatibility layer dapat digunakan jika contract lama masih perlu dipertahankan sementara consumer bermigrasi ke versi baru. Layer ini bertindak sebagai translator antara V1 dan V2 contract.

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
- **Cache key natural berbeda** — URL berbeda = cache key berbeda otomatis

>**Catatan penting:** Versi di URL **tidak** menentukan apakah endpoint idempotent. Idempotency ditentukan oleh semantics HTTP (GET, PUT, DELETE) dan business operation, bukan lokasi version identifier.

>**Trade-off:** URL versioning lebih eksplisit dan mudah di-observe, namun membuat URL lebih panjang. Header versioning memberi URL resource yang bersih, namun routing/observability lebih tersembunyi dan cache/CDN harus aware terhadap version header. Pilihan tergantung prioritas observability vs. clean resource naming.

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

Karakteristik:
- URL resource tetap stabil
- Version selection berada di representation metadata
- Routing/observability lebih tersembunyi namun lebih bersih
- Memungkinkan multiple versions dalam satu endpoint

Tantangan:
- Kurang visibilitas di logs/network traces
- Membutuhkan middleware parsing header
- Caching/CDN harus aware terhadap version header

>**Penting untuk Caching:** Jika representasi berbeda berdasarkan version header (Accept atau API-Version), infrastructure/cache harus memasukkan versi ke cache key atau menggunakan `Vary` header yang sesuai. Contoh:
>```http
>Vary: Accept
>```
>Jika tidak, client dapat menerima representation dari API version yang salah.

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
| Deprecation communication? | Notifikasi eksplisit ke consumer (bukan hanya release notes) |
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
1. V1 tidak boleh menerima breaking contract changes selama masih didukung. Backward-compatible fixes dan additive changes masih dapat dilakukan sesuai compatibility policy.
2. Deploy V2.
3. New Android uses V2.
4. Old Android remains on V1.
5. Monitor adoption.
6. Monitor V1 traffic.
7. Announce deprecation.
8. Consumer communication.
9. Migration monitoring.
10. V1 traffic monitoring.
11. Sunset V1 after criteria are met.

---

## Deprecation Lifecycle

```
V1 active
    ↓
V2 released
    ↓
V1 deprecated
    ↓
Consumer communication
    ↓
Migration monitoring
    ↓
V1 traffic monitoring
    ↓
Sunset criteria reached
    ↓
V1 removed
```

> **Consumer communication harus eksplisit** bukan hanya tersirat lewat release notes atau warning. Komunikasi ini mencakup:
> - Notifikasi resmi ke semua pengguna endpoint V1
> - Dokumentasi dengan timeline deprecasi
> - Support channels untuk membantu migrasi

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
| Hapus field dari published contract tanpa migration/deprecation | consumer yang bergantung pada field dapat gagal | Deprecation cycle dulu |
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
mux := http.NewServeMux()

mux.HandleFunc("/api/v1/invoices", V1Handler)
mux.HandleFunc("/api/v1/invoices/", V1Handler)

mux.HandleFunc("/api/v2/invoices", V2Handler)
mux.HandleFunc("/api/v2/invoices/", V2Handler)
```

Exact collection path dan subtree path didaftarkan eksplisit agar request tanpa trailing slash masuk langsung ke handler dan menghasilkan API error 400, bukan implicit redirect dari http.ServeMux.

Dokumentasi exact routing & error behavior:

- `GET /api/v1/invoices/1001`
  → `200 OK`
- `GET /api/v1/invoices/`
  → `400 Bad Request` (missing invoice ID)
- `GET /api/v1/invoices`
  → `400 Bad Request` (invalid/missing invoice ID)
- `GET /api/v1/invoices/abc`
  → `400 Bad Request` (ID harus numeric)
- `GET /api/v1/invoices/0`
  → `400 Bad Request` (ID harus positive)
- `GET /api/v1/invoices/-1`
  → `400 Bad Request` (ID harus positive)
- `GET /api/v1/invoices/9999`
  → `404 Not Found` (ID valid tetapi invoice tidak ditemukan)
- `GET /api/v1/invoices/1001/extra`
  → `400 Bad Request` (extra path segment)
- `POST /api/v1/invoices/1001`
  → `405 Method Not Allowed` (`Allow: GET`, `Content-Type: application/json`)

V2 menggunakan routing/error semantics yang sama, hanya wire response contract yang berbeda.

Unknown route:
- `GET /api/v3/invoices/1001`
  → `404 Not Found` dari router (default `ServeMux` 404 mengembalikan `text/plain`, bukan JSON API response).

---

## Exercise

API lama:

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}
```

Breaking change:

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
- **Versioning bukan untuk semua**. Major API version biasanya diperlukan ketika perubahan tidak dapat dipertahankan secara backward-compatible.
- **V1 tidak boleh menerima breaking contract changes selama masih didukung**. Backward-compatible fixes dan additive changes masih dapat dilakukan sesuai compatibility policy.
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
├── domain.go               # Model Invoice, Customer, LegacyInvoice, InvoiceRepository
├── unsafe_server.go        # Anti-pattern: breaking change tanpa versioning
├── unsafe_test.go          # Test: legacy client fails, unsafe route behavior
├── safe_server.go          # V1/V2 DTO + mappers + parseInvoiceID helper
├── safe_test.go            # Contract regression tests, V1/V2 semantic assertions
├── additive_server.go      # Additive change: field currency tanpa breaking
├── additive_test.go        # Test: legacy client tetap berhasil
├── routing_test.go         # HTTP routing, 400/404/405, method validation
└── README.md               # Panduan lengkap
```

---

## Navigasi

- **Previous**: [Lab 05 — Race Condition](../05-race-condition/)
- **Next**: [Lab 07 — Observability](../07-observability/)