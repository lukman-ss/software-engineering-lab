# DELIVERABLES SUMMARY - Lab 16: Dependency Injection & IoC

## Quick Reference Sheet

| File | Purpose | Location |
|------|---------|----------|
| **content_def.py** | Source definition (Python) | `result/content_def.py` ✓ |
| **LinkedIn PPTX** | 24 editable carousel slides | `result/linkedin/lab-16-carousel.pptx` ✓ |
| **LinkedIn PDF** | Printable 16:9 presentation | `result/linkedin/lab-16-linkedin.pdf` ✓ |
| **Cover ID PNG** | Indonesian cover (1920x1080) | `result/image-cover/lab-16-cover-id.png` ✓ |
| **Cover EN PNG** | English cover (1920x1080) | `result/image-cover/lab-16-cover-en.png` ✓ |
| **Dev.to Article** | Full English technical post | `result/dev.to/lab-16-devto.md` ✓ |
| **Medium Article** | Bahasa Indonesia (22 sections) | `result/medium/lab-16-medium.html` ✓ |
| **Metadata** | SEO & Social metadata | `result/medium/lab-16-metadata.md` ✓ |
| **LinkedIn Caption** | Carousel discussion caption | `result/linkedin/lab-16-caption.md` ✓ |

## Verification Notes

### Build Verification
- `go test -v ./...` → 6 test cases PASS
- `go test -race ./...` → race-detector clean
- `go run ./cmd/demo` → Both DI and Service Locator execute successfully

### Evidence Claims (Lab-Verified)
- Constructor injection exposes dependency contracts at compile-time
- Service Locator hides dependencies behind container interface
- MockGateway replaces RealGateway without network calls
- `go test -race ./...` passes without warnings
- Validations reject invalid input before gateway delegation

### Boundary Statements (Do Not Overclaim)
- Lab tests only unit logic isolation; does not prove network resiliency
- No measurement of production throughput or connection pool behavior
- No framework DI integration (Wire, Dig, Fx)
- Manual wiring demonstrated; not proof of enterprise-scale adoption

### File Structure
```text
labs/16-dependency-injection/result/
├── DELIVERABLES.md              ← This file
├── content_def.py               ← Source (DO NOT DELETE)
├── dev.to/
│   └── lab-16-dependency-injection-devto.md
├── medium/
│   ├── lab-16-dependency-injection-medium.html
│   └── lab-16-dependency-injection-metadata.md
├── image-cover/
│   ├── lab-16-dependency-injection-cover-id.png
│   └── lab-16-dependency-injection-cover-en.png
└── linkedin/
    ├── lab-16-dependency-injection-carousel.pptx
    ├── lab-16-dependency-injection-linkedin.pdf
    └── lab-16-dependency-injection-caption.md
```

---

## Ready-to-Copy LinkedIn Caption

Ketika sebuah service menginstansiasi adapter eksternal secara langsung di dalam metodenya (seperti memanggil RealGateway di dalam Processor), kode domain terkunci rapat dengan infrastruktur fisik. Pengujian unit menjadi mustahil tanpa membuka koneksi jaringan nyata, dan pengujian lokal sering kali gagal saat server pihak ketiga mengalami downtime.

Di Lab 16 seri Software Engineering Lab, kami mengeksplorasi kontras struktural antara Constructor Injection dan Service Locator di bahasa Go menggunakan manual wiring murni tanpa framework pihak ketiga.

Poin kuncinya: Inversion of Control memisahkan fase konfigurasi (perakitan graf dependensi di composition root main.go) dari fase eksekusi (logika domain). Melalui NewProcessor(g PaymentGateway), kontrak dependensi terlihat gamblang di signature kompilasi dan memungkinkan penyuntikan MockGateway seketika dalam 0.00 detik pengujian unit.

Sebaliknya, Service Locator dengan NewBadProcessor(c Container) menutupi kebutuhan dependensi di balik kontainer generik. Standar industri seperti PSR-11 menegaskan rekomendasi bahwa container tidak boleh disuntikkan ke domain logic karena memindahkan kegagalan kompilasi menjadi runtime error.

Test suite mencakup 6 test cases yang memverifikasi isolasi mock, delegasi pembayaran, propagasi error, dan penolakan nilai negatif sebelum menyentuh gateway, seluruhnya bersih dari race condition (go test -race ./...). Namun penting dicatat: pengujian mock unit membuktikan isolasi logika, bukan jaminan stabilitas koneksi jaringan atau resiliency di skala beban production.

Kode sumber lengkap, suite pengujian, dan implementasi demo tersedia di repositori:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection

#SoftwareEngineering #GoLang #DependencyInjection #CleanArchitecture #SystemDesign #SoftwareDesign