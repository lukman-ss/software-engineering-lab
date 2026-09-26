# Contradictions

## 1. Terminologi "contract testing": provider-only vs. integration

**SOURCE A (provider-only):**
Pact Docs (Source 1, "Provider contract testing" section): menyebut alternatif istilah *"contract testing"* / *"provider contract testing"* untuk arti schema test provider‑only (misalnya verifikasi OpenAPI terhadap implementasi). Dokumen ini secara eksplisit membedakan: *"On its own, however, it does not provide any test based assurance that the consumers are calling the provider in the correct manner... and hence, it is not as effective in preventing integration bugs."*

**SOURCE B (integration):**
Pact Docs (Source 1, utama): *"Contract testing is a technique for testing an integration point by checking each application in isolation..."* + Pactflow (Source 6, definisi): *"asserts that a consumer communicates messages that match a schema, **and that a provider produces output that matches this schema**"* — fokus pada kedua belah pihak.

**ASSESSMENT:**
Bukan kontradiksi faktal, melainkan dua level makna yang sama artinya. Pact Docs sengaja membedakan:
- **provider contract testing** (schema test provider‑only) — satu sisi.
- **contract testing (integration)** — dua sisi consumer+provider.

Keclever-an: istilah "contract testing" dipakai oleh comunity secara luas untuk kedua arti. Pada laporan, dipakai arti integration contract testing (Pact code-first CDC). Provider-only schema test dikelompokkan terpisah dan jujur dibandingkan kekurangannya.

## 2. Schema (OpenAPI) testing "cukup" untuk breaking change?

**SOURCE A (Pactflow pro schema):**
Schema-based tests diyakini cukup bila schema diekstrak/Generate otomatis dari kode (DTO→schema, recording proxy) sehingi drift berkurang. Beberpa poin plus: simpler DX, faster, less duplication.

**SOURCE B (Pactflow anti schema):**
Schema testing **cons**: *"Not all schemas capture key aspects of a contract — no standard way to define HTTP verb/path/status/headers in JSON Schema";*"Coverage — it's very difficult to be sure it fully implements the spec";*"Evolution — schemas are point-in-time"*;*"Code <-> schema drift"*.

**ASSESSMENT:**
Pactflow **sendiri** (satu publikasi, satu author) mengakui keduanya — ini *intentional trade-off write-up*, bukan kontradiksi luar biasa. Kesimpulan: schema testing bukan pengganti penuh CDC untuk integration guarantee; namun dapat menjadi **layer pelengkap**. Tidak ada sumber independen pihak ketiga yang memihak.

## 3. Apakah "additive = aman"?

**SOURCE A (Lab 06):**
Additive change (tambah optional field) biasanya backward-compatible bila consumer **tolerant** ke unknown field. Namun: *"Tidak selalu aman!" bila consumer strict schema*.

**SOURCE B (Pact Docs):**
CDC test hanya meng‑cover field yang **dipakai consumer** — jadi field baru yang tidak dipakai tidak menjadi contract dan tidak akan bikin provider test gagal. Implikasinya: field *baru* tidak akan pernah dideteksi sebagai breaking oleh CDC (karena tidak ter‑test). Namun bila field baru tersebut tidak dipakai tidak berbahaya asal consumer lama toleran.

**ASSESSMENT:**
Kesepakatan: *additive* relatif **aman asal consumer toleran unknown fields**. CDC tidak akan melindungi *penambahan* secara eksplisit — protection-nya adalah *"tambahan tidak akan pernah menjadi contract consumer lama, jadi tidak akan pernah dipaksa."* Bukan kontradiksi, melengkapi. Tapi penting diketahui: jika consumer strict, penambahan optional juga bisa breaking (misal required:true bertambah). Ditandai sebagai kehandahan edge-case pada laporan.

## 4. Tanggal dokumen Pact docs

**SOURCE A:** Footer Pact Docs Introduction: *"Last updated on Aug 25, 2026 by Matt Fellows"* — tanggal di masa depan relatif "hari ini" (2026-09-26), konsisten.

**SOURCE B:** Beberapa halaman (what_is_pact, testing_scope) memberi 404 — berarti dokumen site restructuring terjadi; definisi diambil dari halaman Introduction yang hidup.

**ASSESSMENT:**
Tidak kontradiksi fakta — struktur situs berubah seiring waktu; konten definisi tetap konsisten dan terbukti dapat diakses.

## Ringkasan

Tidak ditemukan **kontradiksi fakta** yang material antar sumber primer. Diskrepansi hanya pada:
(a) istilah ganda "contract testing" (vendor) — ditegaskan secara eksplisit oleh vendor itu sendiri;
(b) batasan "additive = aman" bergantung pada ketekapan consumer — disampaikan secara jujur.

Tidak ada statistik, kutipan, atau tanggal yang diinventarisasi.
