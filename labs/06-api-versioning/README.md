# Senior Software Engineer Daily #6

## Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client

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

## Breaking Change Matrix

Berikut adalah contoh-perubahan yang menjadi **breaking change**:

### Core Changes

| Perubahan | Contoh |
|-----------|--------|
| Rename field | `name` → `full_name` |
| Remove field | hapus `phone` |
| Change primitive type | `price:number` → `price:string` |
| string → object | `customer: "Budi"` → `customer: {id: 15, name: "Budi"}` |
| object → array | `customer: {...}` → `customer: [...]` |
| Change date format | `2026-07-25` → `25/07/2026` |
| Change incompatible nullability | field wajib → optional (atau sebaliknya) |

### Semantic Changes

| Perubahan | Contoh |
|-----------|--------|
| Change enum semantics | status `"ACTIVE"` berarti hal lain |
| Change HTTP status semantics | `200 OK` → `201 Created` untuk update |
| Change required request input | field `id` menjadi required |
| Change pagination contract | `page` → `page_number`, `limit` → `size` |
| Change business meaning endpoint | `/invoices` sekarang mengembalikan draft juga |

---

## Usually Compatible vs Breaking

### Biasanya Compatible

Perubahan sering **backward compatible**:

| Perubahan | Kondisi |
|-----------|---------|
| Add optional response field | consumer tolerant unknown field |
| Add endpoint | tidak memengaruhi existing endpoint |
| Add optional query parameter | client tidak wajib pakai |

### Biasanya Breaking

Perubahan yang **pasti breaking**:

| Perubahan | Alasan |
|-----------|--------|
| Rename | Consumer pakai nama lama |
| Remove | Consumer butuh data itu |
| Change type | Deserialization error |
| Change response shape | Struct mismatch |
| Change semantics | Logika bisnis salah |

### API Versioning ≠ Every Change = V2

**JANGAN** mengajarkan bahwa setiap perubahan harus membuat:
- `v2`
- `v3`
- `v4`

**Gunakan major version** saat contract incompatibility memang membutuhkan version boundary.

---

## URL vs Header Versioning

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

>**Catatan penting:** Versi di URL atau header **tidak** menentukan apakah endpoint idempotent. Idempotency ditentukan oleh semantics HTTP (GET, PUT, DELETE) dan business operation, bukan lokasi version identifier.

> **Catatan:** URL versioning dipilih untuk lab bukan karena "lebih benar secara universal", tapi karena lebih sederhana untuk demonstrasi.

---

## Consumer Inventory

Sebelum melakukan breaking change, Senior Engineer harus menanyakan:

| Pertanyaan | Contoh Jawaban |
|------------|----------------|
| Who consumes this endpoint? | Android (v1.0-v3.5), iOS, Web, Partner ERP, BI Dashboard |
| Who owns each consumer? | Mobile team (Android), Web team, Partner Integration |
| Which versions are active? | Android: 30% on v1, 60% on v2, 10% on v3 |
| Can they deploy independently? | Tidak — semua butuh backend update |
| Are external partners involved? | Ya — 3 partner ERP pakai v1 |
| How much V1 traffic remains? | ~15% di minggu pertama |
| Which mobile versions still hit V1? | Android 8.0-10.0 (ter distribusi lemot) |

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

**Yang harus diingat:**

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

## Android Release ≠ All Users Upgrade

```
Android App Store Distribution (Contoh)
         │ Total Market Share
─────────┼─────────────────────────
 v3.0     │ ■■■■■■■■■■ 45%
 v2.5     │ ■■■■■■■■■■ 35%
 v2.0     │ ■■■■■■■■■■ 15%
 v1.0     │ ■■■■■■■■■■  5%
```

**Kesimpulan:** Bahkan 1 minggu setelah rilis V2, masih ada 5% user pakai V1.

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

**Timeline hanya contoh:**

```
Day 0    → V2 released
Day 30   → V1 deprecated (no new features)
Day 60   → Warning emails sent
Day 90   → Sunset criteria check (< 5% traffic)
Day 100  → V1 removed (jika criteria terpenuhi)
```

> **INGAT:** 90 hari hanyalah **CONTOH**. Actual timeline bergantung pada:
> - **SLA** — berapa lama service harus support?
> - **Partner contract** — berapa lama perjanjian komersial?
> - **Mobile adoption** — seberapa cepat user upgrade?
> - **Business criticality** — apa konsekuensi V1 down?
> - **Security issue** — ada CVE yang harus patched cepat?
> - **Remaining traffic** — berapa persen V1?

---

## Contract Test vs Server Test

### Server Test

Memverifikasi bahwa server **bisa merespons dengan benar**:

```go
func TestServer_responds_ok(t *testing.T) {
    resp := get("/api/v2/invoices/1001")
    if resp.StatusCode != 200 {
        t.Error("server should respond")
    }
}
```

### Compatibility Test (Contract Test)

Memverifikasi bahwa **consumer lama masih dapat memahami contract**:

```go
func TestV1Contract_RemainsBackwardCompatible(t *testing.T) {
    // ... get response ...
    
    // HTTP 200 && valid JSON == OK
    // TAPI belum cukup!
    
    // Harus test:
    // - id = number
    // - customer = string (BREAKING jika object!)
    // - total = number  
    // - status = string
}
```

**Mental model utama Lab 06:**

> `HTTP 200` + `valid JSON` ≠ `backward compatible`

> `LegacyInvoicedecode → error` adalah **bukti breaking change**, bukan bug.

---

## API Documentation

Agar API versioning berhasil, dokumentasi penting banget:

| Dokumen | Tujuan |
|---------|--------|
| OpenAPI/Swagger | Schema resmi setiap version |
| Changelog | Apa yang berubah dari V1 ke V2 |
| Migration guide | Panduan untuk consumer upgrade |
| Deprecation notice | Pengumuman V1 deprecated + timeline |
| Version lifecycle | Aturan create/deprecate/sunset |
| Owner/contact | Siapa harmoni bila ada edge case |

> **Tidak perlu Swagger dependency** untuk lab ini. Cukup dokumentasi di README + komentar kode.

---

## Common Mistakes (Senior Engineer Checklist)

Berikut adalah kesalahan umum yang harus dihindari:

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

Route `/api/v1/invoices/1001` dipetakan ke `V1Handler`.

Route tak dikenal seperti `/foo` tidak akan menghasilkan invoice — `ServeMux` hanya menangani route yang didaftarkan.

Route seperti `/api/v1/invoices/9999-invalid` akan masuk handler, lalu `strconv.Atoi("")` gagal dan mengembalikan `400 Bad Request`.

---

## Final Code Structure

```
labs/06-api-versioning/
├── go.mod                   # Module go 1.22, independent
├── domain.go                # Internal domain Invoice, Customer
├── ParseLegacyInvoice()     # Helper untuk test legacy decode
├── unsafe_server.go         # Anti-pattern: breaking change
├── unsafe_test.go           # Test: legacy gagal decode
├── safe_server.go           # V1/V2 DTO + mappers
├── safe_test.go             # Contract regression tests
├── additive_server.go       # Field tambahan
├── additive_test.go         # Test: legacy tetap berhasil
└── README.md                # Panduan lengkap
```

---

## Key Takeaways

- **API adalah kontrak**, bukan teknologi. Perubahan harus dievaluasi dari perspective consumer.
- **HTTP 200 ≠ compatible**. Selalu test decoding di sisi client.
- **Versioning bukan untuk semua**. Hanya untuk breaking changes.
- **V1 selalu stabil**. Setelah di-release, jangan pernah ubah lagi.
- **Contract test wajib**. Prevent breaking change yang terlewat.
- **Consumer inventory penting**. Tanpa tahu siapa yang pakai, tidak ada cara upgrade yang aman.

---

## Usually Compatible Changes

Namun, perubahan berikut sering **backward compatible**:

| Perubahan | Kondisi |
|-----------|---------|
| Adding optional response field | consumer tolerant unknown field |
| Adding endpoint | tidak memengaruhi existing endpoint |
| Adding optional query parameter | client tidak wajib pakai |

>**Catatan tentang Penambahan Field:**

>Menambahkan optional response field biasanya backward-compatible untuk tolerant readers yang mengabaikan unknown fields. Namun **Tidak selalu aman!** Bergantung pada behavior consumer:
>- `json.Unmarshal` Go: abaikan unknown field → aman
>- JSON Schema strict: field baru bisa mismatch → tidak aman
>- Generated client (OpenAPI codegen): tergantung konfigurasi `additionalProperties`
>- Protobuf/JSON with strict parsing: field baru bisa cause error

---

## API Change Decision Rule

**Backward-compatible change** → **pertahankan version yang sama**

**Contract-breaking change** → **buat major API version baru**

### Kapan Butuh Version Baru?

| Perubahan Type | Contoh | Impact | Required Version |
|----------------|--------|--------|------------------|
| **Request-side** (client → server) | | | |
| Required field menjadi required | `field: string → required string` | Client lupa field error | ✅ Major version |
| Required field menjadi optional | `field: string → optional string` | Biasanya aman | ❌ Tidak perlu |
| Optional field menjadi required | `field: string → required string` | Client lama gagal | ✅ Major version |
| **Response-side** (server → client) | | | |
| Field dihapus | `field` hilang | Client missing data | ✅ Major version |
| Field rename | `name` → `full_name` | Consumer pakai nama lama | ✅ Major version |
| Primitive type berubah | `price: number` → `price: string` | Deserialization error | ✅ Major version |
| Object ↔ Primitive | `customer: "Budi"` → `customer: {id: 15, name: "Budi"}` | Struct mismatch | ✅ Major version |
| Object ↔ Array | `customer: {...}` → `customer: [...]` | Parse error | ✅ Major version |
| Date format berubah | `2024-01-01` → `01/01/2024` | Parsing error | ✅ Major version |
| Nullability berubah | | | |
| Non-nullable → nullable | `field: string` → `field: string | null` | **Request:** biasanya aman**<br>**Response:** breaking jika consumer tidak handle null | ✅ Major version |
| Nullable → non-nullable | `field: string | null` → `field: string` | Response: breaking jika ada null | ✅ Major version |
| **Semantic Changes** | | | |
| Enum semantics berubah | `status: "ACTIVE"` berarti hal lain | Business logic salah | ✅ Major version |
| HTTP status berubah artinya | `200 OK` → `201 Created` untuk update | Konsumen error handling | ✅ Major version |
| Pagination berubah | `page`/`limit` → `page_number`/`size` | Consumer parsing | ✅ Major version |
| Business meaning berubah | `/invoices` sekarang include draft | Logika bisnis berubah | ✅ Major version |

---

## Breaking Change ≠ Hanya JSON Structure

Perubahan yang menjadi breaking change **bukan hanya** tentang struktur JSON. Juga termasuk:

| Perubahan | Contoh | Penjelasan |
|-----------|--------|------------|
| Required → Optional (Request) | Field `authorization` menjadi optional | ✅ Aman - client tidak wajib kirim |
| Non-nullable → Nullable (Response) | `status: "PAID"` → `status: null` | JSON type berubah, consumer error bila tidak handle null |
| Optional → Required (Request) | `requestId` menjadi wajib di header | Client lama tidak kirim = 4xx error |
| Authentication requirement | Endpoint butuh API key baru | Semua client lama gagal otorisasi |
| Error response shape | `{error: "msg"}` → `{code: "ERR", message: "msg"}` | Error parser lama gagal |

---

---

## Unsafe Approach: Breaking Change Tanpa Versioning

### Implementasi

File: `unsafe_server.go`, `unsafe_test.go`

```go
// GET /api/invoices/1001
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

Tapi kurang visibilitas di logs/network traces.

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

## Exercise

Kasus:

```json
{
  "id": 1001,
  "customer": "Budi",
  "total": 500000,
  "status": "PAID"
}
```

Menjadi:

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
2. **Ya, perlu V2** (atau namespace baru) karena ini merusak contract, kecuali jika ada mekanisme header versioning yang mendukung.
3. **Android lama tetap hidup** karena endpoint V1 tetap dilayani backend dengan schema string. Routing memisah request: legacy client memanggil `/v1`, aplikasi baru memanggil `/v2`.
4. **Membuktikan compatibility** dengan **Contract Regression Test** yang secara eksplisit memverifikasi tipe data `customer` harus berupa string, tidak hanya mengandalkan raw string match.

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