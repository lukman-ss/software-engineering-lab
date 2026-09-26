# Open Questions

## Unanswered / Weak Evidence

1. **Case-insensitive enum breaking?**
   Topik spesifikasi menyebutkan perubahan `IN_PROGRESS → in_progress` (casing). Apakah ini breaking secara universal?
   - Jika klien semua dilakukan dengan `status.toLowerCase() === "in_progress"`, maka change ini **bisa** tidak breaking.
   - Namun **tidak ada sumber primer** yang menyatakan "case change pada enum dapat dianggap additive."
   - **Status:** LOW. Perlu verifikasi: carilah pedoman API besar (Stripe, GitHub, Google APIs) — apakah mereka memperlakukan enum case change sebagai breaking?

2. **Berapa ambang toleransi consumer?**
   - "Consumer toleran unknown field" → aman. Tapi apakah semua framework JSON (Jackson, Gson, serde) benar-benar default mengabaikan unknown field? Apakah konfigurasi strict mode ada?
   - **Status:** LOW-MEDIUM. Perlu verifikasi per-library (Go json.Unmarshal, Java Jackson, JS JSON.parse default).

3. **Provider-driven CDC verification — blocking deploy secara eksplisit?**
   - Martin Fowler menyebut contract test *bisa* tidak langsung break build. Pactflow blog menyatakan bisa.  Apakah konfigurasi "block deploy" adalah default?
   - **Status:** MEDIUM. Perlu eksplorasi pada `pactbroker can-i-deploy` atau GitHub Actions pact plugin — apakah exit non-zero secara otomatis?

4. **Schema-based contract testing sebagai pengganti penuh?**
   - Pada repositori ini (Go microservice), apakah ada rencana/standar untuk **schema-based** (mis. OAS) daripada **code-based** (Pact)?
   - **Status:** LOW. Perlu keputusan tim/proyek.

5. **Provider states complexity di skala Go?**
   - Pada Spring Cloud Contract/JAX, provider states dipetakan ke `@State` annotation. Apakah ekosistem Pact untuk Go menyediakan helper setara?
   - **Status:** LOW-MEDIUM. Butuh survey dependensi Go (`github.com/pact-foundation/pact-go`).

6. **Event contract testing — format apa yang diekspektasi?**
   - Untuk Kafka/webhooks, kontrak pesan meliputi: payload schema, message key, headers, topic name, partition. Apakah Pact (message queue plugin) meng-coverall semua ini, atau hanya payload?
   - **Status:** MEDIUM. Perlu verifikasi dokumentasi Pact message plugin.

7. **Consumer inventory sebenarnya ada di repo ini?**
   Lab tidak menyertakan daftar konsumen API. Apakah Mobile App satu-satunya konsumen, atau ada lainnya (web dashboard, partner)?
   - **Status:** LOW. Perlu stakeholder interview.

8. **Kebijanan versioning (URL vs header) belum final**
   - Topic spesifikasi tidak menyatakan kebijakan. Lab 06 pakai URL versioning. Konsistensi?
   - **Status:** LOW. Perlu keputusan tim arsitek.

## Possible Next Research Directions

- **Audit konsumen API** — kumpulkan inventor konsumen GET `/api/work-orders/{id}`.
- **Benchmarking tooling Go** — bandingkan `pact-go` (code-based) vs `Spectral` (schema linting) vs `OpenAPI-diff` (regression).
- **Integrasi ke CI** — proof-of-concept `pact-broker can-i-deploy` di pipeline backend.
- **Event-first contract testing** — eksplorasi [AsyncAPI](https://www.asyncapi.com/) + Pact message plugin untuk Kafka/webhooks.
- **Policy enforcement** — tambahkan ke `opencode`/`AGENTS` workflow rule: setiap PR API harus lampirkan contract test CDC.
