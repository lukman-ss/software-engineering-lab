# API Backward Compatibility

## 1. Modifikasi Payload Secara Aman

### 1.1 Additive Fields (Penambahan Non-Breaking)

**Evidence** (Google AIP-180):
"New required fields **must not** be added to existing request messages or resources ... Any field previously populated by the server **must** continue to be populated, even if it introduces redundancy."

**Evidence** (GitHub API):
Additive changes meliputi "Adding an operation, Adding an optional parameter, Adding a response field, Adding enum values."

**Keterbatasan (AIP-180 Caveats)**:
Menurut AIP-180, penambahan field bersifat non-breaking **hanya jika**:
- Field bersifat **optional** (bukan required)
- Default serialization tidak bergeser (misalnya JSON field order atau default value changes)
- Penambahan enum tidak memaksa client switch/case crash (strict enum handling)
- Tidak ada naming conflict di generated code stubs (protobuf, OpenAPI generators)

**Praktik**:
```json
{
  "customer_name": "Budi",
  "customer": { "id": 10, "name": "Budi" }
}
```

### 1.2 Internal Transformation Pipelines (Stripe Model)

**Evidence** (Stripe):
"Every possible response from the Stripe API is codified by a class that we call an API resource. API resources define their possible fields using a DSL ... When we need to make a backwards-incompatible change, we encapsulate it in a version change module which defines documentation about the change, a transformation, and the set of API resource types that are eligible to be modified."

**Evidence** (Stripe, lanjut):
"When generating a response, the API initially formats data by describing an API resource at the current version, then determines a target API version from one of: A Stripe-Version header if one was supplied. The user's pinned version, which is set on their very first request to Stripe. It then walks back through time and applies each version change module that finds along the way until that target version is reached."

### 1.3 URL Versioning

**Evidence** (GitHub):
When a new REST API version is released, the previous API version will be supported for at least 24 more months.

**Trade-off**:
- **Pro**: Isolasi kode lama dan baru; rollback mudah
- **Kontra**: Duplikasi controller; technical debt; cost definisi yang tinggi

## 2. Field Removal, Rename, Type Changes

### 2.1 Removing atau Renaming Components

**Evidence** (Google AIP-180):
"Existing components (interfaces, methods, messages, fields, enums, or enum values) **must not** be removed from existing APIs in the same major version. Removing a component is a backwards incompatible change."
"Renaming a component is semantically equivalent to 'remove and add'."

### 2.2 Field Type Changes

**Evidence** (Google AIP-180):
"Existing fields and messages **must not** have their type changed, even if the new type is wire-compatible, because type changes alter generated code in a breaking way."

### 2.3 Semantic Changes

**Evidence** (Google AIP-180):
"APIs **must not** change visible behavior or semantics in ways that are likely to break reasonable user code, as such changes will be seen as breaking by those users."

## 3. API Lifecycle dan Deprecation

### 3.1 The Deprecation Window

**Evidence** (GitHub API):
| Fase | Deskripsi |
|------|-----------|
| **ACTIVE** | Current standard |
| **DEPRECATED** | Endpoint masih berfungsi, tapi klien diminta migrasi |
| **SUNSET** | API diterbitkan formal dan dimatikan (410 Gone atau 404) |
| **REMOVED** | Kode dihapus dari repository |

**Evidence** (GitHub, lanjut):
"While a version is within its support window but approaching closing down, GitHub includes [Deprecation and Sunset] headers in API responses to help you prepare for migration ... After the support window ends: Requests that specify a closing down API version receive a 410 Gone response."

### 3.2 Metadata Deprecation (RFC 8594)

**Evidence** (GitHub):
"[API] version interfaces ... Deprecation — The date when the API version will be closing down, formatted as an HTTP date per RFC 7231. Sunset — The date when the API version will be completely removed (retired), after which requests will return a 410 Gone response. Follows RFC 8594."

### 3.3 Breaking vs Non-Breaking Changes Daftar Lengkap

**Evidence** (GitHub API):
| Tipe | Contoh | Kebutuhan |
|------|--------|-----------|
| Breaking | Removing an entire operation, Removing/renaming parameter, Removing/renaming response field, Adding required parameter, Making previously optional parameter required, Changing type, Removing enum values, Changing auth requirements | Membutuhkan versi baru |
| Non-breaking (Additive) | Adding operation, Adding optional parameter, Adding optional header, Adding response field, Adding response header, Adding enum values | Tersedia di semua versi yang didukung |

## 4. Deteksi Consumer Aktif (Observability)

**Evidence** (GitHub):
"Never guess if a field is still used. Ensure visibility before deletion: Log metrics when clients access deprecated fields ... Analyze load balancer / API gateway logs for old endpoint access."

**Evidence** (Martin Fowler):
"Any system using feature flags should expose some way for an operator to discover the current state of the toggle configuration. In an HTTP-oriented SOA system this is often accomplished via some sort of metadata API endpoint."

---

## Catatan Kualitas Bukti

| Pernyataan | Sumber | Confidence | Keterangan |
|------------|--------|------------|------------|
| Additive fields = non-breaking | Google AIP-180, GitHub | HIGH | Didokumentasikan secara eksplisit |
| Renaming = remove + add | Google AIP-180 | HIGH | Dinyatakan secara eksplisit |
| Semantic changes = breaking | Google AIP-180 | HIGH | Alur keputusan dijelaskan |
| Sunset headers (Deprecation, Sunset) | GitHub API, RFC 8594 | HIGH | Implementasi nyata GitHub |