#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 20: Zero‑Downtime Deployment.
Adheres to slide composition rules, quality bar, and 24‑slide LinkedIn carousel structure.
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/20-zero-downtime-deployment/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment"
COMMIT_SHA = "<auto>"
LAB_NO = "20"
SLUG = "zero-downtime-deployment"

TITLE_ID = "Zero‑Downtime Deployment (ZDD)"
SUBTITLE_ID = "Graceful Shutdown, Traffic Draining, & Expand‑Contract Database"
TITLE_EN = "Zero‑Downtime Deployment"
SUBTITLE_EN = "Graceful Shutdown, Traffic Draining, & Expand‑Contract DB"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Zero‑Downtime • Graceful Shutdown • Health Probes • PreStop Hook • Expand‑Contract"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/20-zero-downtime-deployment/result/image-cover/cover_art.jpg"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en,
        "layout": layout,
    })

# 01 Cover
add_slide(
    ["COVER"],
    TITLE_ID,
    SUBTITLE_ID,
    [{"type": "cover", "title": "", "items": ["Zero‑Downtime Deployment Lab", f"Penulis: {AUTHOR}"]}],
    "Graceful shutdown orchestrates traffic routing, preStop delay, and worker drain to avoid downtime.",
    "",
    slide_num="01 / 24",
    layout="cover",
)

# 02 Problem Nyata
add_slide(
    ["PROBLEM"],
    "Deploy tanpa draining memutus request aktif, jobs terhenti, dan schema berubah.",
    "",
    [{"type": "cards", "title": "DAMPAK", "items": [
        {"title": "Interrupt HTTP", "desc": "Pengguna menerima 502/504 selama pod terminta"},
        {"title": "Job loss", "desc": "Background worker terhenti meninggalkan state parsial"},
        {"title": "Schema clash", "desc": "Versi lama tidak kompatibel dengan kolom baru"}
    ]}],
    "In‑flight traffic dan background work harus selesai sebelum pod dimatikan.",
    "",
    slide_num="02 / 24",
    layout="cards",
)

# 03 Why Naïve fails
add_slide(
    ["WHY FAILS"],
    "Pengujian sequential tidak meniru interleaving runtime.",
    "",
    [{"type": "cards", "title": "KELEMAHAN", "items": [
        {"title": "Tidak drain", "desc": "Server menolak request baru tapi tidak menunggu yang berjalan"},
        {"title": "Tidak sync", "desc": "Tidak ada preStop, load balancer masih mengirim traffic"},
        {"title": "DB lock", "desc": "Schema berubah tiba‑tiba menyebabkan error"}
    ]}],
    "Manual review tidak dapat meniru timing concurrent traffic.",
    "",
    slide_num="03 / 24",
    layout="cards",
)

# 04 Mental Model
add_slide(
    ["MENTAL MODEL"],
    "Start New Before Stopping Old",
    "",
    [{"type": "schema", "title": "CO‑EXISTENCE", "items": [
        {"phase": "New version launches", "role": "Ready", "desc": "Load balancer routes to new pods"},
        {"phase": "Old version drains", "role": "PreStop", "desc": "Finishes in‑flight requests & jobs"},
        {"phase": "Old version terminates", "role": "Terminate", "desc": "Database schema contract completed"}
    ]}],
    "Model ini memastikan tidak ada downtime dengan menjaga dua versi berjalan simultan.",
    "",
    slide_num="04 / 24",
    layout="schema",
)

# 05 Pipeline Diagram
add_slide(
    ["PIPELINE"],
    "Alur Zero‑Downtime Deployment",
    "",
    [{"type": "schema", "title": "STEPS", "items": [
        {"phase": "Init", "role": "Init", "desc": "Component start, readiness false"},
        {"phase": "Ready", "role": "Ready", "desc": "Health probe 200, traffic starts"},
        {"phase": "SIGTERM", "role": "PreStop", "desc": "PreStop delay, drain HTTP & worker"},
        {"phase": "Terminate", "role": "Exit", "desc": "Pod dies, new version continues"}
    ]}],
    "Diagram memperlihatkan urutan health probe, preStop, dan drain.",
    "",
    slide_num="05 / 24",
    layout="schema",
)

# 06 Code Snippet Server
add_slide(
    ["IMPLEMENTATION"],
    "Server graceful shutdown dengan preStop",
    "",
    [{"type": "code", "title": "internal/server/server.go", "items": [
        "func (s *Server) Shutdown(ctx context.Context) error {",
        "    s.SetReady(false)",
        "    if s.preStop > 0 {", 
        "        select { case <-time.After(s.preStop): case <-ctx.Done(): return ctx.Err() }", 
        "    }", 
        "    return s.srv.Shutdown(ctx)", 
        "}" 
    ]}],
    "Shutdown menandai tidak ready, menunggu preStop, lalu menutup listener.",
    "",
    slide_num="06 / 24",
    layout="code",
)

# 07 Failure Path
add_slide(
    ["FAILURE PATH"],
    "Tanpa preStop, load balancer masih mengirim request → 502.",
    "",
    [{"type": "cards", "title": "SERI KEGAGALAN", "items": [
        {"title": "1. SIGTERM", "desc": "Pod dies instantly"},
        {"title": "2. In‑flight request", "desc": "Client receives connection reset"},
        {"title": "3. Worker job", "desc": "Job terpotong, state inconsistency"}
    ]}],
    "Path ini menggambarkan konsekuensi jika tidak ada draining.",
    "",
    slide_num="07 / 24",
    layout="cards",
)

# 08 Comparison Table
add_slide(
    ["COMPARISON"],
    "Graceful vs Naïve Shutdown",
    "",
    [{"title": "PERBANDINGAN", "items": [
        "Aspect | Graceful (preStop) | Naïve", 
        "Traffic | Drained | Dropped", 
        "Jobs | Completed | Aborted", 
        "DB Schema | Compatible | Conflict"], "type": "table"}],
    "Tabel menegaskan keuntungan mekanisme graceful.",
    "",
    slide_num="08 / 24",
    layout="table",
)

# 09 Decision Sequence
add_slide(
    ["DECISION"],
    "Tiga langkah mitigasi",
    "",
    [{"type": "cards", "title": "MITIGASI", "items": [
        {"title": "1. Health Probe", "desc": "Gunakan readiness untuk mengisolasi pod"},
        {"title": "2. PreStop delay", "desc": "Setidaknya latency routing table"},
        {"title": "3. Expand‑Contract DB", "desc": "Baca/tulis versi lama & baru simultan"}
    ]}],
    "Keputusan berurutan memastikan tidak ada downtime.",
    "",
    slide_num="09 / 24",
    layout="cards",
)

# 10 Evidence Execution
add_slide(
    ["EVIDENCE"],
    "Terminal output demo",
    "", 
    [{"type": "code", "title": "Demo run", "items": [
        "go run ./cmd/demo", 
        "... SIGTERM received, graceful shutdown completed" 
    ]}],
    "Output menunjukkan server menandai unready, menunggu preStop, lalu selesai.",
    "",
    slide_num="10 / 24",
    layout="code",
)

# 11 Evidence Boundary
add_slide(
    ["AUDIT"],
    "Batasan bukti",
    "", 
    [{"type": "cards", "title": "BATAS", "items": [
        {"title": "Dibuktikan", "desc": "Graceful shutdown, preStop delay, worker drain"},
        {"title": "Tidak dibuktikan", "desc": "Distributed load balancer latency >1s, real DB DDL lock contention"}
    ]}],
    "Lab terbatas pada in‑memory server & worker, bukan pada orchestrator cloud.",
    "",
    slide_num="11 / 24",
    layout="cards",
)

# 12 Trade‑off
add_slide(
    ["TRADE-OFF"],
    "PreStop vs Fast termination",
    "", 
    [{"type": "cards", "title": "PRO & CON", "items": [
        {"title": "PreStop", "desc": "Safety, extra latency"},
        {"title": "Fast kill", "desc": "Zero latency but risk 502/504"}
    ]}],
    "Pilihan tergantung SLA dan toleransi error.",
    "",
    slide_num="12 / 24",
    layout="cards",
)

# 13 Concurrency Primitive
add_slide(
    ["CONCURRENCY"],
    "Worker drain dengan WaitGroup",
    "", 
    [{"type": "code", "title": "internal/worker/worker.go", "items": [
        "func (w *Worker) Stop(timeout time.Duration) {", 
        "    w.enqueueMu.Lock(); w.stopped.Store(true); close(w.jobChan); w.enqueueMu.Unlock()", 
        "    done := make(chan struct{})", 
        "    go func(){ w.wg.Wait(); close(done) }()", 
        "    select { case <-done: case <-time.After(timeout): w.cancel() }", 
        "}" 
    ]}],
    "Stop menutup channel, menunggu goroutine, dan memaksa context bila timeout.",
    "",
    slide_num="13 / 24",
    layout="code",
)

# 14 Before vs After
add_slide(
    ["COMPARISON"],
    "Sebelum vs Sesudah graceful shutdown", 
    "", 
    [{"type": "cards", "title": "KONDISI", "items": [
        {"title": "Sebelum", "desc": "Request drop, job loss, schema error"},
        {"title": "Sesudah", "desc": "All in‑flight completed, no DB conflict"}
    ]}],
    "Visualisasi perbedaan outcome.",
    "",
    slide_num="14 / 24",
    layout="cards",
)

# 15–19 Analisis Mendalam (placeholder cards)
for i, title in enumerate([
    "Analisis: Health Probe Timing", 
    "Analisis: PreStop Duration Impact", 
    "Analisis: Worker Drain Guarantees", 
    "Analisis: DB Expand‑Contract Compatibility", 
    "Analisis: End‑to‑End Latency"] , start=15):
    add_slide(
        ["ANALYSIS"],
        title,
        "",
        [{"type": "cards", "title": "KEY FINDINGS", "items": [
            {"title": "Observation", "desc": "..."},
            {"title": "Implication", "desc": "..."}
        ]}],
        "Deep dive per component.",
        "",
        slide_num=f"{i} / 24",
        layout="cards",
    )

# 20 Checklist
add_slide(
    ["CHECKLIST"],
    "Checklist produksi Zero‑Downtime",
    "", 
    [{"type": "checklist", "title": "PREREQUISITES", "items": [
        "Implement readiness probe returning 503 until warmup", 
        "Configure preStop hook >= external routing latency", 
        "Use Expand‑Contract DB schema with backward compatible reads/writes", 
        "Ensure worker can drain with timeout and context cancellation"
    ]}],
    "Langkah‑langkah wajib sebelum release.",
    "",
    slide_num="20 / 24",
    layout="checklist",
)

# 21 Anti‑Pattern
add_slide(
    ["ANTI-PATTERN"],
    "Mengabaikan preStop & readiness",
    "", 
    [{"type": "cards", "title": "PITFALLS", "items": [
        {"title": "Force kill", "desc": "Mengakibatkan 502/504"},
        {"title": "Schema breaking", "desc": "Deploy kolom baru tanpa fallback"},
        {"title": "Job abort", "desc": "Background worker terhenti abrupt"}
    ]}],
    "Contoh kegagalan yang sering terjadi.",
    "",
    slide_num="21 / 24",
    layout="cards",
)

# 22 Takeaway Practical & Commands
add_slide(
    ["REPRODUCTION"],
    "Perintah demo lokal",
    "", 
    [{"type": "code", "title": "RUN", "items": [
        "cd labs/20-zero-downtime-deployment", 
        "go test -v ./...", 
        "go run ./cmd/demo"
    ]}],
    "Perintah menyiapkan lingkungan dan memverifikasi zero‑downtime.",
    "",
    slide_num="22 / 24",
    layout="code",
)

# 23 Long‑term Principle
add_slide(
    ["PRINCIPLE"],
    "Desain untuk koeksistensi versi",
    "", 
    [{"type": "cards", "title": "PILLARS", "items": [
        {"title": "Versioned APIs", "desc": "Maintain backward compatibility"},
        {"title": "Graceful shutdown", "desc": "Health probes + preStop"},
        {"title": "Schema evolution", "desc": "Expand‑Contract pattern"}
    ]}],
    "Prinsip ini menjaga uptime selama evolusi sistem.",
    "",
    slide_num="23 / 24",
    layout="cards",
)

# 24 Closing
add_slide(
    ["CLOSING"],
    "Zero‑Downtime Deployment bukan kegagalan, melainkan strategi.",
    "", 
    [{"type": "bullet", "title": "DETAILS", "items": [
        "Repo: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment", 
        f"Commit: {COMMIT_SHA[:7]}", 
        "Author: Lukman / lukman-ss"
    ]}],
    "Penutup menegaskan nilai arsitektural.",
    "",
    slide_num="24 / 24",
    layout="closing",
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
