#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 28: Timeouts & Deadlines.
Defines metadata, visual palette, and exact 24 slides.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/28-timeouts-and-deadlines/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/28-timeouts-and-deadlines"
COMMIT_SHA = "bdb285d35a71c9a309e83a0e4497b9dbd842ceb8"
LAB_NO = "28"
SLUG = "timeouts-and-deadlines"
TITLE_ID = "Timeouts & Deadlines"
SUBTITLE_ID = "Mencegah Satu Dependency Lambat Menghancurkan Seluruh Ekosistem Distributed System"
TITLE_EN = "Timeouts & Deadlines"
SUBTITLE_EN = "Preventing a Single Slow Dependency from Crashing an Entire Distributed System"
AUTHOR = "Lukman / lukman-ss"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })

# Slide 1: Cover (Dark Navy, Retro Terminal)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines bercahaya",
            "Context propagation, exponential backoff full jitter, circuit breaker, idempotency store",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Slow vs Dead Dependency
add_slide(
    categories=["MASALAH", "LITTLE'S LAW"],
    title="Slow Dependency Lebih Berbahaya daripada Dead Dependency",
    subtitle="Fast failure melepaskan resource seketika; latency tinggi menahan thread pool",
    content_items=[
        {"title": "FAST FAILURE (DEAD)", "items": [
            "Mengembalikan ECONNREFUSED atau reset paket secara instan",
            "Worker thread langsung dilepas ke pool dalam orde mikrodetik",
            "Sistem caller tetap responsif untuk request lain",
        ], "type": "bullet"},
        {"title": "SLOW DEPENDENCY", "items": [
            "Latency melonjak drastis, misal dari 300ms menjadi 60 detik",
            "Request in-flight tertahan di thread pool, memory, dan socket",
            "Endpoint independen (/profile, /health) kehabisan worker thread",
        ], "type": "bullet"},
        {"title": "FORMULA LITTLE'S LAW", "items": [
            "L = λ × W (Konkurensi = Arrival Rate × Latency)",
            "Jika W naik 200× (0.3s ke 60s), kapasitas L meledak 200× lipat",
            "Thread pool 100 worker hanya mampu melayani 1.6 QPS pada W=60s",
        ], "type": "bullet"},
    ],
    mental_model="Fast failure frees resources immediately; slow dependencies hoard threads until total starvation.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="2 / 24"
)

# Slide 3: Mental Model of Timeout
add_slide(
    categories=["KONSEP", "MENTAL MODEL"],
    title="Timeout Adalah Resource Guardrail, Bukan Tebakan Aman",
    subtitle="Empat prinsip fundamental dalam perancangan time budget",
    content_items=[
        {"title": "TIMEOUT ≠ FAILURE", "items": [
            "Timeout berarti caller tidak menerima response tepat waktu",
            "Status backend adalah indeterminate: mungkin sukses, mungkin gagal",
            "Jangan pernah mengasumsikan transaksi dibatalkan hanya karena timeout",
        ], "type": "bullet"},
        {"title": "BUDGET HIERARKIS", "items": [
            "Deadline harus diwariskan dari hulu ke hilir distributed trace",
            "Child timeout tidak boleh melebihi sisa waktu parent deadline",
            "Downstream harus membatalkan eksekusi jika parent telah timeout",
        ], "type": "bullet"},
        {"title": "PROTEKSI TOTAL", "items": [
            "Retry tanpa jitter menciptakan thundering herd retry storm",
            "HTTP timeout tidak melindungi DB connection pool dari query liar",
            "Gunakan statement_timeout pada database storage layer",
        ], "type": "bullet"},
    ],
    mental_model="Timeout means caller gave up waiting, not that the backend operation was aborted.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="3 / 24"
)

# Slide 4: Lab Architecture
add_slide(
    categories=["ARSITEKTUR", "DESAIN"],
    title="Komponen Pertahanan Terdistribusi Lab 28",
    subtitle="Empat subsistem mandiri yang bekerja bersinergi menahan kegagalan",
    content_items=[
        {"title": "1. CONTEXT DEADLINE", "items": [
            "Package internal/deadline",
            "Pewarisan context.WithTimeout dengan channel non-blocking",
            "Pencegahan goroutine leak saat caller membatalkan eksekusi",
        ], "type": "bullet"},
        {"title": "2. RETRY WITH FULL JITTER", "items": [
            "Package internal/retry",
            "Exponential backoff 2^(attempt-1) dikombinasikan seragam acak",
            "Menghilangkan thundering herd antar worker",
        ], "type": "bullet"},
        {"title": "3. CIRCUIT BREAKER & DEDUP", "items": [
            "Package internal/circuit: State CLOSED, OPEN, HALF_OPEN",
            "Package internal/idempotency: Store in-memory dengan lazy eviction",
            "Melindungi downstream yang kolaps dan mencegah double charge",
        ], "type": "bullet"},
    ],
    mental_model="A resilient pipeline combines deadlines, jittered retries, circuit breaking, and idempotency.",
    source_path="labs/28-timeouts-and-deadlines/engineering/01-design.md",
    slide_num="4 / 24"
)

# Slide 5: Context Deadline Propagation
add_slide(
    categories=["KODE", "DEADLINE"],
    title="Implementasi Context Deadline di Go",
    subtitle="Memastikan child budget tunduk pada parent context tanpa memory leak",
    content_items=[
        {"title": "CHILD TIMEOUT BUDGET", "items": [
            "childCtx, cancel := context.WithTimeout(ctx, budget)",
            "defer cancel() selalu dipanggil untuk membersihkan timer",
            "Child deadline = min(parent.Deadline(), now + budget)",
        ], "type": "bullet"},
        {"title": "BUFFERED CHANNEL (CAP=1)", "items": [
            "done := make(chan error, 1)",
            "Kapasitas 1 menjamin goroutine pekerja tidak pernah block",
            "Mencegah goroutine leak saat select memilih childCtx.Done()",
        ], "type": "bullet"},
        {"title": "BUKTI PENGUJIAN", "items": [
            "TestExecuteWithBudget_Timeout: task 100ms dipotong budget 20ms",
            "TestExecuteWithBudget_ParentTimeoutInherited: parent 20ms batalkan budget 500ms",
            "Mengembalikan context.DeadlineExceeded secara deterministik",
        ], "type": "bullet"},
    ],
    mental_model="Always buffer the completion channel to capacity 1 to prevent goroutine leaks on timeout.",
    source_path="labs/28-timeouts-and-deadlines/internal/deadline/deadline.go",
    slide_num="5 / 24"
)

# Slide 6: Deadline Inheritance
add_slide(
    categories=["PROPAGASI", "HIERARKI"],
    title="Pewarisan Deadline dalam Hierarki Distributed Trace",
    subtitle="Waktu tidak bertambah di downstream; deadline selalu menyusut",
    content_items=[
        {"title": "SKENARIO API GATEWAY", "items": [
            "API Gateway menetapkan timeout 500ms untuk request user",
            "Service A butuh 150ms pemrosesan lokal sebelum memanggil Service B",
            "Budget aktual yang diteruskan ke Service B tersisa 350ms",
        ], "type": "bullet"},
        {"title": "DOWNSTREAM OVER-ESTIMATION", "items": [
            "Jika Service B membuat timeout 2 detik sendiri tanpa parent ctx",
            "Saat detik ke-1 user sudah menerima 504 Gateway Timeout",
            "Service B tetap membakar CPU untuk request yang tidak lagi dinanti",
        ], "type": "bullet"},
        {"title": "MEKANISME GOLANG", "items": [
            "context.WithTimeout secara otomatis mewarisi parent deadline",
            "Jika parent deadline lebih pendek dari budget baru, parent menang",
            "Sub-operasi berhenti bersamaan saat parent deadline terpicu",
        ], "type": "bullet"},
    ],
    mental_model="Time budgets strictly shrink as calls descend through microservice trees.",
    source_path="labs/28-timeouts-and-deadlines/internal/deadline/deadline_test.go",
    slide_num="6 / 24"
)

# Slide 7: Thundering Herd & Retries
add_slide(
    categories=["MASALAH", "RETRY STORM"],
    title="Bahaya Retry Tanpa Jitter: Thundering Herd",
    subtitle="Retry deterministik menyinkronkan ribuan client menjadi gelombang perusak",
    content_items=[
        {"title": "SINKRONISASI CLIENT", "items": [
            "1.000 client gagal bersamaan saat database spike sejenak",
            "Semua client menjadwalkan retry tepat pada attempt × 100ms",
            "Downstream dihantam gelombang request periodik yang sinkron",
        ], "type": "bullet"},
        {"title": "RETRY STORM COLLAPSE", "items": [
            "Server downstream yang baru pulih langsung tumbang kembali",
            "Cascading failure merambat naik ke API Gateway",
            "Load amplification: 1 request asli menjadi 3-5 request bersamaan",
        ], "type": "bullet"},
        {"title": "SOLUSI AWS ARCHITECTURE", "items": [
            "Bukan sekadar exponential backoff berlipat ganda",
            "Wajib menambahkan randomisasi penuh (Full Jitter)",
            "Meredistribusikan beban secara merata di sepanjang sumbu waktu",
        ], "type": "bullet"},
    ],
    mental_model="Deterministic retries synchronize clients into destructive thundering herds.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="7 / 24"
)

# Slide 8: Full Jitter Algorithm
add_slide(
    categories=["ALGORITMA", "JITTER"],
    title="Matematika Exponential Backoff dengan Full Jitter",
    subtitle="sleep = random(0, min(maxBackoff, baseBackoff × 2^(attempt-1)))",
    content_items=[
        {"title": "FORMULA KALKULASI", "items": [
            "multiplier := 1 << uint(attempt - 1)",
            "temp := min(MaxBackoff, BaseBackoff * multiplier)",
            "sleep := rand.Float64() * temp (Distribusi Seragam)",
        ], "type": "bullet"},
        {"title": "KEUNGGULAN FULL JITTER", "items": [
            "AWS Research membuktikan Full Jitter meminimalkan total waktu kerja",
            "Mengurangi server contention hingga titik terendah dibanding Equal Jitter",
            "Setiap client memiliki waktu tunggu yang independen secara statistik",
        ], "type": "bullet"},
        {"title": "BATASAN IMPLEMENTASI", "items": [
            "attempt <= 0 ditangani aman dengan membatasi shift bit",
            "Backoff dibatasi MaxBackoff untuk mencegah overflow",
            "Default zero config: 3 attempts, base 10ms, max 100ms",
        ], "type": "bullet"},
    ],
    mental_model="Full jitter samples uniformly from zero to the exponential ceiling, maximizing dispersion.",
    source_path="labs/28-timeouts-and-deadlines/internal/retry/retry.go",
    slide_num="8 / 24"
)

# Slide 9: Retrier Implementation
add_slide(
    categories=["KODE", "RETRY"],
    title="Loop Eksekusi Retrier dengan Respek Context",
    subtitle="Retry loop yang berhenti seketika saat context timeout tercapai",
    content_items=[
        {"title": "STRUKTUR EKSEKUSI", "items": [
            "Iterasi 1 hingga cfg.MaxAttempts secara berurutan",
            "Mengeksekusi operasi fn(ctx) di tiap iterasi",
            "Jika err == nil, fungsi langsung return tanpa sleep",
        ], "type": "bullet"},
        {"title": "RESPEK TERHADAP CONTEXT", "items": [
            "Sleep dilakukan menggunakan select { case <-time.After(): ... }",
            "case <-ctx.Done(): return ctx.Err() membatalkan jeda tidur",
            "Tidak membuang waktu menunggu jika caller sudah membatalkan",
        ], "type": "bullet"},
        {"title": "ERROR HANDLING", "items": [
            "Mengembalikan ErrMaxRetriesExceeded jika semua percobaan gagal",
            "Membungkus lastErr untuk debugging riwayat kegagalan",
            "Test membuktikan retry berhenti saat parent context hangus",
        ], "type": "bullet"},
    ],
    mental_model="A retrier must abort its sleep immediately when the governing context is cancelled.",
    source_path="labs/28-timeouts-and-deadlines/internal/retry/retry.go",
    slide_num="9 / 24"
)

# Slide 10: Circuit Breaker State Machine
add_slide(
    categories=["CIRCUIT", "STATE MACHINE"],
    title="Circuit Breaker: Menghindari Hammering Downstream",
    subtitle="Tiga status proteksi: CLOSED, OPEN, dan HALF-OPEN",
    content_items=[
        {"title": "CLOSED STATE (NORMAL)", "items": [
            "Semua request diteruskan langsung ke downstream service",
            "Menghitung kegagalan berturut-turut (consecutive failures)",
            "Jika failure count >= FailureThreshold (misal 2), transisi ke OPEN",
        ], "type": "bullet"},
        {"title": "OPEN STATE (ISOLASI)", "items": [
            "Semua request langsung ditolak dengan ErrCircuitOpen",
            "Tidak ada satupun request yang dikirim ke downstream",
            "Memberikan ruang napas bagi downstream yang kolaps",
        ], "type": "bullet"},
        {"title": "HALF-OPEN STATE (UJI COBA)", "items": [
            "Aktif otomatis setelah CooldownWindow (misal 50ms) berlalu",
            "Mengizinkan request uji coba masuk secara bertahap",
            "Sukses berturut-turut memulihkan ke CLOSED; 1 gagal kembali ke OPEN",
        ], "type": "bullet"},
    ],
    mental_model="Circuit breaker fails fast to give dying downstream dependencies room to recover.",
    source_path="labs/28-timeouts-and-deadlines/internal/circuit/circuit.go",
    slide_num="10 / 24"
)

# Slide 11: Circuit Breaker Concurrency
add_slide(
    categories=["KODE", "THREAD-SAFETY"],
    title="Thread-Safe State Transitions pada Breaker",
    subtitle="Sinkronisasi state menggunakan sync.RWMutex untuk throughput tinggi",
    content_items=[
        {"title": "METODE ALLOW() & READ LOCK", "items": [
            "b.mu.RLock() membaca state saat ini secara konkuren",
            "Mengecek apakah durasi OPEN telah melampaui CooldownWindow",
            "Promosi ke b.mu.Lock() hanya saat transisi status ke HALF_OPEN",
        ], "type": "bullet"},
        {"title": "RECORD SUCCESS & FAILURE", "items": [
            "b.mu.Lock() eksklusif memperbarui counter dan timestamp",
            "Di HALF_OPEN: 1 kegagalan langsung memicu OPEN kembali",
            "Di HALF_OPEN: sukses mencapai threshold memicu pemulihan ke CLOSED",
        ], "type": "bullet"},
        {"title": "VERIFIKASI INTEGRITAS", "items": [
            "TestCircuitBreaker_StateTransitions memvalidasi siklus penuh",
            "TestCircuitBreaker_HalfOpenFailureTripsOpen memvalidasi sensitivitas",
            "Lolos pengujian race detector tanpa race condition",
        ], "type": "bullet"},
    ],
    mental_model="Fast reads under RLock; atomic state mutations under exclusive Lock.",
    source_path="labs/28-timeouts-and-deadlines/internal/circuit/circuit.go",
    slide_num="11 / 24"
)

# Slide 12: Ambiguity of Timeout
add_slide(
    categories=["MASALAH", "IDEMPOTENCY"],
    title="Ambiguitas Status Saat Operasi Mengalami Timeout",
    subtitle="Apakah request benar-benar gagal atau berhasil tapi respon hilang?",
    content_items=[
        {"title": "SKENARIO JARINGAN ASINKRON", "items": [
            "Client mengirim request POST /charges untuk potong saldo $100",
            "Database bank sukses mendebit saldo nasabah di backend",
            "Kabel fiber putus atau reverse proxy timeout sebelum ACK terkirim",
        ], "type": "bullet"},
        {"title": "DILEMA CLIENT CALLER", "items": [
            "Client hanya menerima pesan: context deadline exceeded",
            "Client TIDAK BISA membedakan apakah transaksi sukses atau gagal",
            "Jika caller langsung melakukan retry mentah, saldo terpotong $200",
        ], "type": "bullet"},
        {"title": "PERLINDUNGAN IDEMPOTENSI", "items": [
            "Setiap operasi yang mengubah state wajib menyertakan Idempotency-Key",
            "Server menyimpan hash request dan response hasil pemrosesan",
            "Retry dengan key yang sama tidak boleh mengeksekusi ulang transaksi",
        ], "type": "bullet"},
    ],
    mental_model="Timeout is indeterminate: the database write may have committed even if the network dropped.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Idempotency Key Store
add_slide(
    categories=["KODE", "DEDUP"],
    title="In-Memory Idempotency Store & Lazy Eviction",
    subtitle="Penyimpanan response dengan TTL dan pembersihan pasif tanpa timer leak",
    content_items=[
        {"title": "RECORD STRUKTUR & TTL", "items": [
            "type Record struct { Response string; CreatedAt time.Time }",
            "Set(key, response) menyimpan payload dengan timestamp saat ini",
            "Store menjaga masa retensi data sesuai configured TTL (misal 50ms)",
        ], "type": "bullet"},
        {"title": "LAZY EVICTION PATTERN", "items": [
            "Tidak menggunakan background goroutine periodik pemakan resource",
            "Get(key) memeriksa selisih time.Since(rec.CreatedAt) > s.ttl",
            "Record kadaluarsa dihapus seketika saat diakses pembaca",
        ], "type": "bullet"},
        {"title": "WRITE LOCK PADA METHOD GET", "items": [
            "Method Get() menggunakan s.mu.Lock(), bukan RLock()",
            "Penting karena lazy eviction melakukan mutasi delete(s.records, key)",
            "Menjamin konsistensi memori dan lolos race detector di bawah 50 goroutine",
        ], "type": "bullet"},
    ],
    mental_model="Lazy eviction cleans expired keys during lookups, eliminating background cleaner goroutines.",
    source_path="labs/28-timeouts-and-deadlines/internal/idempotency/idempotency.go",
    slide_num="13 / 24"
)

# Slide 14: Integrating Retry + Circuit Breaker
add_slide(
    categories=["INTEGRASI", "CIRCUIT+RETRY"],
    title="Integrasi: Retry Loop Bersanding dengan Circuit Breaker",
    subtitle="Breaker memotong rantai retry seketika saat downstream terbukti down",
    content_items=[
        {"title": "EKSEKUSI BERSAMA", "items": [
            "Retrier dikonfigurasi 4 attempt dengan backoff 5ms - 20ms",
            "Circuit Breaker memiliki failure threshold 2 dan cooldown 50ms",
            "Downstream disimulasikan selalu mengalami kegagalan 500 error",
        ], "type": "bullet"},
        {"title": "PEMOTONGAN CEPAT (FAST ABORT)", "items": [
            "Attempt #1 gagal -> Breaker failure count = 1",
            "Attempt #2 gagal -> Breaker failure count = 2 -> State beralih ke OPEN",
            "Attempt #3 & #4 langsung ditolak lokal oleh Allow() tanpa call network",
        ], "type": "bullet"},
        {"title": "PROTEKSI TOTAL", "items": [
            "Mencegah retry membebani backend yang sedang berjuang pulih",
            "Waktu eksekusi caller terpangkas drastis tanpa sleep yang sia-sia",
            "TestIntegration_RetryWithCircuitBreaker membuktikan pemotongan ini",
        ], "type": "bullet"},
    ],
    mental_model="Circuit breakers short-circuit futile retries before they hit an already failing downstream.",
    source_path="labs/28-timeouts-and-deadlines/tests/integration_test.go",
    slide_num="14 / 24"
)

# Slide 15: Integrating Idempotent Retry
add_slide(
    categories=["INTEGRASI", "IDEMPOTENT RETRY"],
    title="Integrasi: Safe Retries Melalui Idempotency Store",
    subtitle="Memastikan operasi non-idempoten aman dieksekusi berulang kali",
    content_items=[
        {"title": "SKENARIO SIMULASI", "items": [
            "Operasi bisnis: actualExecutions counter mencatat eksekusi nyata",
            "Request pertama berhasil menyimpan hasil ke Idempotency Store",
            "Tetapi caller mengalami transient network disconnect sebelum terima data",
        ], "type": "bullet"},
        {"title": "EKSEKUSI RETRY LOOP", "items": [
            "Retrier menjalankan 3 attempt pengulangan dengan request yang sama",
            "Store.Get(idempotencyKey) menemukan cache record transaksi sebelumnya",
            "Mengembalikan response tersimpan dengan label DEDUPLICATED",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI", "items": [
            "actualExecutions tetap bernilai 1 meskipun attempt berlangsung berkali-kali",
            "Uang nasabah hanya terdebit sekali; saldo aman dari double charging",
            "TestIntegration_IdempotentRetry mengonfirmasi proteksi absolut ini",
        ], "type": "bullet"},
    ],
    mental_model="Idempotency keys make retries safe by returning cached results for duplicate requests.",
    source_path="labs/28-timeouts-and-deadlines/tests/integration_test.go",
    slide_num="15 / 24"
)

# Slide 16: Verification & Test Suite
add_slide(
    categories=["VERIFIKASI", "TESTING"],
    title="Verifikasi Komprehensif: 13 Test Case Lolos 100%",
    subtitle="Pengujian menyeluruh mencakup unit, konkurensi, dan integrasi lintas modul",
    content_items=[
        {"title": "DISTRIBUSI TEST SUITE", "items": [
            "internal/circuit: 3 tests (StateTransitions, HalfOpenFailure, ZeroConfig)",
            "internal/deadline: 3 tests (Success, Timeout, ParentInherited)",
            "internal/retry: 5 tests (FirstTry, RetryUntilSuccess, ExceedMax, ContextCanceled, Bounds)",
            "internal/idempotency: 3 tests (GetSet, ConcurrentAccess, LazyEviction)",
        ], "type": "bullet"},
        {"title": "INTEGRATION SUITE", "items": [
            "TestIntegration_RetryWithCircuitBreaker -> PASS",
            "TestIntegration_IdempotentRetry -> PASS",
            "Race detector (go test -race ./...) berjalan bersih tanpa data race",
        ], "type": "bullet"},
        {"title": "DURASI & STABILITAS", "items": [
            "Total 13 automated test cases selesai dalam ~1.6 detik",
            "Semua race window diuji dengan goroutine agresif",
            "Zero flake, deterministik, dan reproducible",
        ], "type": "bullet"},
    ],
    mental_model="13/13 tests pass cleanly under the Go race detector, verifying concurrent correctness.",
    source_path="labs/28-timeouts-and-deadlines/tests/integration_test.go",
    slide_num="16 / 24"
)

# Slide 17: Interactive Demo Program
add_slide(
    categories=["DEMO", "EKSEKUSI"],
    title="Program Demo Interaktif: cmd/demo/main.go",
    subtitle="Demonstrasi langsung 4 skenario ketahanan dalam terminal",
    content_items=[
        {"title": "SKENARIO 1 & 2", "items": [
            "Demo 1: Context Deadline membatalkan task yang melampaui budget",
            "Output: Deadline propagation result: context deadline exceeded",
            "Demo 2: Retrier mencoba 3 kali dengan Full Jitter hingga sukses",
            "Output: Attempt #1, #2, #3 dieksekusi dengan jeda acak seragam",
        ], "type": "bullet"},
        {"title": "SKENARIO 3", "items": [
            "Demo 3: Circuit Breaker dari CLOSED -> 2 gagal -> OPEN",
            "Call ditolak instan saat OPEN; cooldown -> HALF_OPEN -> sukses -> CLOSED",
            "Output membuktikan pemulihan otomatis saat downstream sehat kembali",
        ], "type": "bullet"},
        {"title": "SKENARIO 4", "items": [
            "Demo 4: Idempotent Request Retry Protection",
            "Attempt 1: Charged $100 successfully",
            "Attempt 2 (Retry): Charged $100 successfully (DEDUPLICATED)",
        ], "type": "bullet"},
    ],
    mental_model="The runnable demo proves deadline cancellation, jittered retry, trip/recovery, and dedup in action.",
    source_path="labs/28-timeouts-and-deadlines/cmd/demo/main.go",
    slide_num="17 / 24"
)

# Slide 18: Common Architectural Anti-Patterns
add_slide(
    categories=["ANTI-PATTERN", "KESALAHAN"],
    title="Empat Anti-Pattern Umum Terkait Timeout",
    subtitle="Kebiasaan buruk di industri yang memicu cascading failure di produksi",
    content_items=[
        {"title": "1. TIMEOUT ANGKA BESAR", "items": [
            "Menyetel timeout 30s-60s 'supaya aman dan tidak gampang error'",
            "Menyebabkan antrian thread meledak saat dependency melambat sejenak",
            "Harus berbasis latency P99 riil + toleransi sempit (misal 2.5× P99)",
        ], "type": "bullet"},
        {"title": "2. RETRY TANPA JITTER", "items": [
            "Retry terjadwal kaku pada interval eksak (1s, 2s, 4s)",
            "Menciptakan fenomena thundering herd yang melumpuhkan recovery backend",
            "Wajib acak penuh dengan Full Jitter",
        ], "type": "bullet"},
        {"title": "3. BLIND RETRY & ONLY HTTP", "items": [
            "Mengulang request mutasi tanpa idempotency key -> data corrupt",
            "Hanya memasang timeout di reverse proxy/HTTP, lupa membatasi SQL query",
            "PostgreSQL connection pool tetap terkuras habis oleh query gantung",
        ], "type": "bullet"},
    ],
    mental_model="Arbitrary high timeouts and jitterless retries turn minor glitches into catastrophic outages.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="18 / 24"
)

# Slide 19: Database Timeouts: The Missing Layer
add_slide(
    categories=["DATABASE", "STORAGE LAYER"],
    title="PostgreSQL Timeouts: Perlindungan Connection Pool",
    subtitle="Mencegah query lambat atau lock gantung membajak seluruh koneksi database",
    content_items=[
        {"title": "STATEMENT TIMEOUT", "items": [
            "SET statement_timeout = '5s'",
            "Membatalkan query secara paksa jika durasi komputasi melebihi batas",
            "Mencegah table scan tak terindeks menghabiskan CPU database",
        ], "type": "bullet"},
        {"title": "LOCK TIMEOUT", "items": [
            "SET lock_timeout = '2s'",
            "Membatalkan transaksi jika menunggu exclusive lock terlalu lama",
            "Mencegah antrian DDL/UPDATE memblokir ribuan pembaca (connection pileup)",
        ], "type": "bullet"},
        {"title": "IDLE IN TRANSACTION", "items": [
            "SET idle_in_transaction_session_timeout = '10s'",
            "Memutus koneksi yang membuka BEGIN tapi tidak kunjung COMMIT",
            "Mencegah hold lock tak sengaja dan bloat vacuum pada tabel",
        ], "type": "bullet"},
    ],
    mental_model="Application timeouts are useless if slow database queries hold open connection pools indefinitely.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="19 / 24"
)

# Slide 20: gRPC & Cross-Service Propagation
add_slide(
    categories=["PROTOKOL", "GRPC & HTTP"],
    title="Propagasi Batas Waktu Lintas Protokol Jaringan",
    subtitle="Meneruskan deadline melewati batas proses dan host fisik",
    content_items=[
        {"title": "GRPC HEADERS (STANDAR)", "items": [
            "Header resmi: grpc-timeout (misal 500m untuk 500 milidetik)",
            "Downstream gRPC client otomatis mengurangi durasi transmisi",
            "Menerapkan context deadline langsung ke server handler",
        ], "type": "bullet"},
        {"title": "HTTP/REST HEADERS", "items": [
            "Belum ada standar tunggal ISO; heterogen antar vendor",
            "Pola umum: Request-Timeout, X-Request-Timeout, atau Baggage",
            "W3C Baggage / OpenTelemetry Context Injection semakin populer",
        ], "type": "bullet"},
        {"title": "CLOCK SKEW DEFENSE", "items": [
            "Jangan pernah mengirim absolut wall-clock time (misal 14:00:05.123)",
            "Perbedaan NTP antar host (clock skew) dapat menyebabkan false timeout",
            "Selalu transmisikan sisa durasi relatif (remaining duration)",
        ], "type": "bullet"},
    ],
    mental_model="Transmit remaining duration, never absolute timestamps, to eliminate clock skew hazards.",
    source_path="labs/28-timeouts-and-deadlines/research/05-report.md",
    slide_num="20 / 24"
)

# Slide 21: Lab Boundaries & Simulator Scope
add_slide(
    categories=["BATAS", "AUDIT"],
    title="Batasan Eksplisit Implementasi Lab 28",
    subtitle="Pemisahan transparan antara simulator lab dan kebutuhan produksi nyata",
    content_items=[
        {"title": "IN-MEMORY IDEMPOTENCY", "items": [
            "Lab menyimpan key di memori dengan sync.Mutex",
            "Data hilang saat server crash atau restart",
            "Produksi butuh Redis cluster atau database transaksional",
        ], "type": "bullet"},
        {"title": "PROSES TUNGGAL (SINGLE PROCESS)", "items": [
            "Deadline diuji dalam context Go proses tunggal",
            "Tidak mencakup transmisi serialisasi wire protocol gRPC/HTTP",
            "Tidak memodelkan dynamic adaptive timeout berbasis P99 live",
        ], "type": "bullet"},
        {"title": "DATABASE STORAGE LAYER", "items": [
            "Lab membuktikan konsep melalui arsitektur software",
            "Konfigurasi PostgreSQL statement_timeout dibahas konseptual",
            "Tidak menjalankan live cluster Docker PostgreSQL",
        ], "type": "bullet"},
    ],
    mental_model="Educational implementations highlight core principles; production requires distributed state stores.",
    source_path="labs/28-timeouts-and-deadlines/engineering/02-implementation-notes.md",
    slide_num="21 / 24"
)

# Slide 22: Before vs After Comparison
add_slide(
    categories=["KOMPARASI", "BEFORE VS AFTER"],
    title="Perbandingan: Tanpa Perlindungan vs Sistem Resilien",
    subtitle="Dampak langsung penerapan pola Timeouts & Deadlines",
    content_items=[
        {"title": "KONDISI TANPA PERLINDUNGAN", "items": [
            "Default HTTP Client Go: No Timeout (menunggu selamanya)",
            "Dependency lambat 60s: Thread pool terkuras, QPS anjlok ke 1.6",
            "Retry instan: Terjadi thundering herd retry storm",
            "Network timeout pada mutasi: Rentan double charge transaksi",
        ], "type": "bullet"},
        {"title": "DENGAN TIME BUDGET & RESILIENCE", "items": [
            "Time Budget: Eksekusi dibatasi tegas; resource segera dilepas",
            "Hierarchical Deadline: Downstream berhenti jika parent dibatalkan",
            "Full Jitter Retries: Beban tersebar merata secara statistik",
            "Circuit Breaker & Idempotency: Downstream dilindungi, data aman",
        ], "type": "bullet"},
        {"title": "HASIL OBSERVABILITAS", "items": [
            "Sistem tahan terhadap kegagalan parsial tanpa downtime total",
            "Resource usage (thread, socket, memory) tetap terkendali stabil",
            "Zero orphan goroutine / memory leak di background",
        ], "type": "bullet"},
    ],
    mental_model="Resilience patterns transform uncontrolled systemic failure into graceful degraded operation.",
    source_path="labs/28-timeouts-and-deadlines/content/02-master-draft.md",
    slide_num="22 / 24"
)

# Slide 23: Core Principles & Checklist
add_slide(
    categories=["PRINSIP", "BEST PRACTICE"],
    title="Prinsip Inti Perancangan Timeouts di Produksi",
    subtitle="Tiga pilar utama untuk membangun distributed system berdaya tahan tinggi",
    content_items=[
        {"title": "1. SET EXPLICIT BOUNDS EVERYWHERE", "items": [
            "Jangan pernah menggunakan default timeout nol pada client manapun",
            "Pasang batas waktu pada HTTP client, gRPC, database, dan messaging queue",
            "Tetapkan timeout berdasarkan data observabilitas P99 latency aktual",
        ], "type": "bullet"},
        {"title": "2. ALWAYS JITTER RETRIES", "items": [
            "Gunakan Full Jitter untuk mendispersi waktu tunggu antar client",
            "Kombinasikan dengan Circuit Breaker untuk memotong retry sia-sia",
            "Batasi jumlah maksimal retry secara konservatif (2-3 kali)",
        ], "type": "bullet"},
        {"title": "3. SAFEGUARD MUTATIONS WITH KEYS", "items": [
            "Anggap timeout sebagai status indeterminate",
            "Kirim Idempotency-Key pada semua operasi POST / mutasi state",
            "Simpan status eksekusi dengan TTL untuk menghentikan double execution",
        ], "type": "bullet"},
    ],
    mental_model="Bound every call, jitter every retry, and key every mutation for absolute resilience.",
    source_path="labs/28-timeouts-and-deadlines/content/05-key-takeaways.md",
    slide_num="23 / 24"
)

# Slide 24: Closing (Dark Navy)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 28: Timeouts & Deadlines",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 13/13 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "KOMPONEN TERVERIFIKASI", "items": [
            "Context Deadline: non-blocking channel, inheritance terjaga",
            "Full Jitter Retrier: distribusi acak seragam, anti thundering herd",
            "Circuit Breaker: 3-state transisi aman konkuren (CLOSED/OPEN/HALF_OPEN)",
            "Idempotency Store: thread-safe dedup dengan lazy TTL eviction",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI", "items": [
            "13 unit & integration test PASS (100%)",
            "Race detector clean (go test -race ./...)",
            "Demo interaktif berhasil membuktikan seluruh skenario",
        ], "type": "bullet"},
        {"title": "SUMBER KODE & REPOSITORI", "items": [
            f"Repository: {REPO_URL}",
            f"Commit Hash: {COMMIT_SHA}",
            "Uji mandiri: go test -v -race ./... && go run ./cmd/demo",
            f"Author: {AUTHOR}",
        ], "type": "bullet"},
    ],
    mental_model="Distributed resilience requires time bounds at every layer: deadline, retry jitter, breaker, and dedup.",
    source_path="labs/28-timeouts-and-deadlines/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
