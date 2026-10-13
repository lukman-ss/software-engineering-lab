# Publication Deliverables — Software Engineering Lab 20

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-20-zero-downtime-deployment-carousel.pptx](linkedin/lab-20-zero-downtime-deployment-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes (via python-pptx) |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-20-zero-downtime-deployment-linkedin.pdf](linkedin/lab-20-zero-downtime-deployment-linkedin.pdf) | File PDF 24 slide rasio 16:9 (Cover CRT + 22 slide krem editorial + Penutup navy) |
| 3 | Dev.to Article | [dev.to/lab-20-zero-downtime-deployment-devto.md](dev.to/lab-20-zero-downtime-deployment-devto.md) | Artikel Markdown Bahasa Inggris lengkap dengan frontmatter |
| 4 | Medium Article | [medium/lab-20-zero-downtime-deployment-medium.html](medium/lab-20-zero-downtime-deployment-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.340 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-20-zero-downtime-deployment-cover-id.png](image-cover/lab-20-zero-downtime-deployment-cover-id.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-20-zero-downtime-deployment-cover-en.png](image-cover/lab-20-zero-downtime-deployment-cover-en.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-20-zero-downtime-deployment-metadata.md](medium/lab-20-zero-downtime-deployment-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-20-zero-downtime-deployment-caption.md](linkedin/lab-20-zero-downtime-deployment-caption.md) | Caption teks 227 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `04890623a7f333e780a38551686fa7dba2284471`
- **Status Test**: 16/16 test cases PASS (`go test -race -v ./...` durasi 2.17s, 0 data races).
- **Batas Implementasi**: Lab ini memodelkan simulasi in-memory Go untuk menguji siklus hidup HTTP draining, jeda propagasi routing jaringan (`preStop`), terminasi antrean worker kooperatif, serta logika fallback database Expand/Contract secara deterministik. Lab tidak mensimulasikan kluster Kubernetes fisik, load balancer cloud eksternal, atau proses database PostgreSQL nyata dengan lock contention antrean kueri. Sizing timeout produksi wajib menyelaraskan hierarki: `preStop delay < HTTP connection draining timeout < orchestrator terminationGracePeriodSeconds`.

## Caption LinkedIn (Siap Salin)

```text
Lab 20 kali ini membahas arsitektur zero-downtime deployment, koordinasi probe, penundaan terminasi jaringan, serta evolusi skema database paralel di Go.

Masalah teknis yang sering terjadi adalah pemutusan sepihak instance lama saat deployment memicu lonjakan error HTTP 502 Bad Gateway pada request in-flight, membatalkan job background di tengah eksekusi, dan menyebabkan crash akibat skema database yang belum kompatibel.

Lab ini membuktikan empat pilar arsitektur konkret:
• Decoupled Probes: Liveness probe tetap melaporkan status proses sehat sementara readiness probe menandai unready untuk melepas instance dari rotasi load balancer.
• Asynchronous Detachment: Hook preStop berbasis select-timer menunda penutupan listener HTTP agar tabel routing eksternal sempat diperbarui sebelum soket ditutup.
• Connection Draining: Server menyelesaikan seluruh request aktif in-flight secara utuh sebelum proses dihentikan, diverifikasi tanpa kebocoran transaksi.
• Expand and Contract: Lapisan penyimpanan database mendukung pembacaan fallback transparan antara skema warisan dan skema modern tanpa downtime migrasi.

Mental model inti: Start new before stopping old. Coexistence requires version N and version N-1 to run simultaneously, route traffic via decoupled probes, and evolve database schemas in backward-compatible phases.

Satu batasan penting: Lab ini menggunakan simulasi in-memory Go untuk menguji alur draining HTTP, worker queue, dan skema database, bukan cluster Kubernetes fisik atau database PostgreSQL nyata, serta mengasumsikan drain timeout worker membatalkan konteks jika pekerjaan melampaui batas waktu.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment

#softwareengineering #golang #devops #kubernetes #systemdesign #backend
```
