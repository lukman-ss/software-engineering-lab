# DELIVERABLES SUMMARY - Lab 16: Dependency Injection

## Quick Reference Sheet

| File | Purpose | Location |
|------|---------|----------|
| **content_def.py** | Source definition (Python) | `result/content_def.py` ✓ |
| **LinkedIn PPTX** | 24 editable carousel slides | `result/linkedin/lab-16-dependency-injection-carousel.pptx` ✓ |
| **LinkedIn PDF** | Printable 16:9 presentation | `result/linkedin/lab-16-dependency-injection-linkedin.pdf` ✓ |
| **Cover ID PNG** | Indonesian cover (1920x1080) | `result/image-cover/lab-16-dependency-injection-cover-id.png` ✓ |
| **Cover EN PNG** | English cover (1920x1080) | `result/image-cover/lab-16-dependency-injection-cover-en.png` ✓ |
| **Dev.to Article** | Full English technical post | `result/dev.to/lab-16-dependency-injection-devto.md` ✓ |
| **Medium Article** | Bahasa Indonesia (22 sections) | `result/medium/lab-16-dependency-injection-medium.html` ✓ |
| **Metadata** | SEO & Social metadata | `result/medium/lab-16-dependency-injection-metadata.md` ✓ |
| **LinkedIn Caption** | Carousel discussion caption | `result/linkedin/lab-16-dependency-injection-caption.md` ✓ |

## Verification Notes

### Build Verification
- `go test -v ./...` → 6 test cases PASS
- `go test -race ./...` → race-detector clean
- `go run ./cmd/demo` → Output: `RealGateway charging 100 USD` & `RealGateway charging 200 USD`

### Evidence Claims (Lab-Verified)
- Constructor Injection enables mock injection without infrastructure
- MockGateway replaces RealGateway completely for isolated unit tests
- Test suite validates happy path, invalid input blocking, and error propagation
- Race detector passes on all test cases
- Service Locator anti-pattern demonstrated with hidden container dependency

### Boundary Statements (Do Not Overclaim)
- DI memudahkan testing, bukan bukti performa produksi
- MockGateway menggantikan network call secara penuh (bukan berarti production behavior sama)
- Manual wiring type-safe; bukan otomatis production-ready tanpa framework
- Test membuktikan aturan struktural pada skenario uji; bukan jaminan semua kebutuhan tim

### File Structure
```text
labs/16-dependency-injection/result/
├── DELIVERABLES.md                              ← This file
├── content_def.py                               ← Source (DO NOT DELETE)
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

Seringkali kita melihat kode di mana class/struct membuat dependensinya sendiri secara langsung di dalam method. Tampak sepele, tapi ini mengunci implementasi konkret ke dalam logika bisnis. Akibatnya: pengujian unit terisolasi mustahil dilakukan tanpa infrastruktur nyata, dan pergantian vendor/provider memaksa refactoring masif di puluhan file.

Lab 16 membedah prinsip inti Dependency Injection (DI) dan Inversion of Control (IoC) lewat Go:
- Memisahkan fase *construction* dari *use* via Constructor Injection.
- Membuktikan Service Locator sebagai anti-pattern yang menyembunyikan dependensi.
- Memvalidasi mocking tanpa network overhead (test PASS 0.00s, race-clean).
- Batasan tegas: DI mempermudah pengujian, bukan jaminan otomatis untuk performa atau reliability production.

Code & Lab: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection
Snapshot: commit 8772235

#SoftwareEngineering #GoLang #DesignPatterns #DependencyInjection #CleanArchitecture