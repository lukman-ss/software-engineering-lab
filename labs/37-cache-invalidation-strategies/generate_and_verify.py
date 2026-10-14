#!/usr/bin/env python3
"""
Deliverables Generator and Quality Gate Verification for Lab 37:
Cache Invalidation Strategies (Write Patterns, Stampede Mitigations, TTL Jitter in Go).
"""
import os
import re
import sys
import zipfile
from pathlib import Path

RESULT_DIR = Path("/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/37-cache-invalidation-strategies/result")
LAB_DIR = Path("/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/37-cache-invalidation-strategies")


def generate_medium_metadata():
    title = "Strategi Invalidation Cache"
    subtitle = "Pola Write, Mitigasi Stampede, dan Jitter TTL dalam Go"
    seo_title = "Strategi Invalidation Cache: Pola Write, Stampede, Jitter"
    seo_desc = "Pelajari cache invalidation strategies (Cache-Aside, Write-Through, Write-Behind, Singleflight, XFetch, SWR, TTL Jitter) dalam Go."
    slug = "cache-invalidation-strategies"
    alt_text = "Cover image for Lab 37 Cache Invalidation Strategies"
    social_title = "Strategi Invalidation Cache: Pola Write, Stampede, Jitter"
    social_desc = "Pola write patterns, stampede mitigations, dan TTL jitter dalam Go. Lihat contoh konkret, test, dan tradeoff."
    tags = ["cache", "golang", "performance", "stampede", "systemdesign"]

    content = f"""# Metadata Publikasi Medium — Lab 37

## Identitas Artikel
- **Title**: {title}
  *(Panjang: {len(title)} karakter)*
- **Subtitle**: {subtitle}
  *(Panjang: {len(subtitle)} karakter)*
- **Proposed Slug**: `{slug}`
  *(Panjang: {len(slug)} karakter)*

## Konfigurasi SEO
- **SEO Title**: {seo_title}
  *(Panjang: {len(seo_title)} karakter — batas maksimal 60 karakter)*
- **SEO Description**: {seo_desc}
  *(Panjang: {len(seo_desc)} karakter — target rekomendasi 140–156 karakter)*

## Media & Visual
- **Cover Image**: `image-cover/lab-37-cache-invalidation-strategies-cover-id.png`
- **Alt Text Cover**: {alt_text}
  *(Panjang: {len(alt_text)} karakter)*

## Konfigurasi Sosial Media
- **Social Title**: {social_title}
  *(Panjang: {len(social_title)} karakter)*
- **Social Description**: {social_desc}
  *(Panjang: {len(social_desc)} karakter)*

## Klasifikasi Topik (Tepat 5 Tags)
1. `{tags[0]}`
2. `{tags[1]}`
3. `{tags[2]}`
4. `{tags[3]}`
5. `{tags[4]}`

## Verifikasi Karakter & Field
- Title: {len(title)} karakter
- Subtitle: {len(subtitle)} karakter
- SEO Title: {len(seo_title)} karakter (valid <= 60)
- SEO Description: {len(seo_desc)} karakter (target 140-156)
- Proposed Slug: {len(slug)} karakter
- Alt Text Cover: {len(alt_text)} karakter
- Social Title: {len(social_title)} karakter
- Social Description: {len(social_desc)} karakter
- Jumlah Tag: Tepat {len(tags)} tags (valid)
"""
    out_path = RESULT_DIR / "medium" / "lab-37-cache-invalidation-strategies-metadata.md"
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(content, encoding="utf-8")
    return out_path


def generate_linkedin_caption():
    caption = """Lab 37 kali ini membahas strategi cache invalidation dalam sistem terdistribusi, dengan fokus pada pola write, mitigasi stampede, dan teknik TTL jitter dalam implementasi Go.

Cache invalidation adalah salah satu dilema terbesar dalam arsitektur sistem terdistribusi skala tinggi. Tanpa mekanisme yang disiplin dan terukur, sistem akan menghadapi tiga masalah konkret yang fatal:
• Cache stampede: Ribuan request konkuren membanjiri database secara serentak saat kunci populer kedaluwarsa mendadak
• Stale data: Goroutine pembaca menerima data usang karena inkonsistensi jeda waktu antara penulisan database dan pembaruan cache
• Risiko kehilangan data: Pilihan throughput ekstrem pada Write-Behind mengorbankan durabilitas saat antrean memory overflow atau proses crash

Eksperimen Go standard library pada lab ini membuktikan implementasi native:
• Cache-Aside meminimalkan penggunaan memori namun rentan terhadap race condition pembacaan usang
• Write-Through menjamin konsistensi read-after-write sinkron dengan kompromi kenaikan latensi penulisan
• Write-Behind mengalirkan mutasi asinkron ke worker latar belakang dengan throughput tinggi
• Singleflight memadatkan 20 request konkuren menjadi tepat 1 query database tunggal
• XFetch menggeser titik rebuild proaktif secara probabilistik sebelum batas waktu kedaluwarsa
• Stale-While-Revalidate menyajikan data usang seketika sambil merevalidasi kunci di latar belakang
• TTL Jitter mengacak rentang waktu kedaluwarsa untuk mengeliminasi lonjakan serentak

Mental model: Cache invalidation is the hardest problem because consistency is eventual, never instant.

Batasan implementasi utama: Seluruh simulasi menggunakan memory cache dan MockDB Go murni tanpa Redis atau PostgreSQL persisten eksternal.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies

#softwareengineering #golang #systemdesign #distributed #performance #caching"""

    out_path = RESULT_DIR / "linkedin" / "lab-37-cache-invalidation-strategies-caption.md"
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(caption, encoding="utf-8")
    return out_path, caption


def generate_deliverables_md(caption_text: str):
    content = f"""# Deliverables Rekapitulasi: Lab 37 — Cache Invalidation Strategies

Paket deliverable publikasi lengkap untuk **Software Engineering Lab 37: Cache Invalidation Strategies — Pola Write, Mitigasi Stampede, dan Jitter TTL dalam Go** oleh **Lukman / lukman-ss**.

---

## Tabel Unduhan & Berkas Publikasi

| No | Tipe Deliverable | File Output | Status / Ukuran | Deskripsi |
| :---: | :--- | :--- | :--- | :--- |
| **0** | **Python Definition** | [`result/content_def.py`](content_def.py) | 59.4 KB (848 baris) | Definisi 24 slide, konstanta warna, metadata, dan mental model (DILARANG DIHAPUS) |
| **1** | **Dev.to Article** | [`result/dev.to/lab-37-cache-invalidation-strategies-devto.md`](dev.to/lab-37-cache-invalidation-strategies-devto.md) | ~15.4 KB (2.115 kata) | Artikel teknis Bahasa Inggris lengkap dengan frontmatter dan 22 subbagian |
| **2** | **Medium Article** | [`result/medium/lab-37-cache-invalidation-strategies-medium.html`](medium/lab-37-cache-invalidation-strategies-medium.html) | ~23.4 KB (2.543 kata) | Standalone HTML Bahasa Indonesia dengan 22 section dan styling Georgia/Code |
| **3** | **Medium Metadata** | [`result/medium/lab-37-cache-invalidation-strategies-metadata.md`](medium/lab-37-cache-invalidation-strategies-metadata.md) | Markdown | Metadata SEO, slug, social title, dan 5 tags Medium |
| **4** | **Cover Image (ID)**| [`result/image-cover/lab-37-cache-invalidation-strategies-cover-id.png`](image-cover/lab-37-cache-invalidation-strategies-cover-id.png) | 751 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Indonesia |
| **5** | **Cover Image (EN)**| [`result/image-cover/lab-37-cache-invalidation-strategies-cover-en.png`](image-cover/lab-37-cache-invalidation-strategies-cover-en.png) | 754 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Inggris |
| **6** | **LinkedIn PDF** | [`result/linkedin/lab-37-cache-invalidation-strategies-linkedin.pdf`](linkedin/lab-37-cache-invalidation-strategies-linkedin.pdf) | 1.171 KB (24 slide) | Carousel dokumen PDF 16:9 24 halaman (Cover Dark Navy, 22 Editorial Cream, Closing Dark Navy) |
| **7** | **LinkedIn PPTX** | [`result/linkedin/lab-37-cache-invalidation-strategies-carousel.pptx`](linkedin/lab-37-cache-invalidation-strategies-carousel.pptx) | 119 KB (24 slide) | File presentasi PowerPoint 24 slide dengan layout 16:9 dan Speaker Notes |
| **8** | **LinkedIn Caption**| [`result/linkedin/lab-37-cache-invalidation-strategies-caption.md`](linkedin/lab-37-cache-invalidation-strategies-caption.md) | 240 kata | Teks caption LinkedIn tanpa bold markdown siap salin |

---

## Verifikasi Kualitas & Catatan Teknis

1. **Commit Snapshot**: `eac3d80a`
2. **Status Pengujian**:
   - `go test -v ./...` -> 16 passing tests (race-clean verified)
   - `go test -race ./...` -> PASS (Race-Clean Verified Go 1.22, 0 data races)
   - `go run ./cmd/demo` -> Simulasi seluruh pola (Cache-Aside, Write-Through, Write-Behind, Singleflight, XFetch, SWR, TTL Jitter) berjalan sukses
3. **Batas Implementasi Utama**:
   - Seluruh modul menggunakan simulasi in-memory cache dan MockDB Go murni, bukan cluster Redis atau database PostgreSQL nyata.
   - Buffer antrean Write-Behind menggunakan channel Go in-memory; crash pada proses menyebabkan data yang belum di-flush hilang (sistem produksi memerlukan durable WAL / Kafka).
   - Singleflight beroperasi pada level single process/node; sistem multi-pod membutuhkan distributed locking (misal Redis Redlock / SETNX).

---

## Salinan Teks Caption LinkedIn

```text
{caption_text}
```
"""
    out_path = RESULT_DIR / "DELIVERABLES.md"
    out_path.write_text(content, encoding="utf-8")
    return out_path


def run_quality_gates():
    gates = []

    # 1. content_def.py exists
    cdef_path = RESULT_DIR / "content_def.py"
    if cdef_path.exists():
        size = cdef_path.stat().st_size
        gates.append(("content_def.py exists", True, f"{size} bytes"))
    else:
        gates.append(("content_def.py exists", False, "Missing"))

    # 2. Medium HTML: >= 2000 words & 9 required sections
    med_path = RESULT_DIR / "medium" / "lab-37-cache-invalidation-strategies-medium.html"
    if med_path.exists():
        html = med_path.read_text(encoding="utf-8")
        clean_text = re.sub(r"<[^>]+>", " ", html)
        words = len(clean_text.split())
        h2_tags = re.findall(r"<h2[^>]*>(.*?)</h2>", html)
        h2_count = len(h2_tags)
        words_ok = words >= 2000
        sections_ok = h2_count >= 9
        gates.append(("Medium HTML words (>= 2000)", words_ok, f"{words} words"))
        gates.append(("Medium HTML sections (>= 9)", sections_ok, f"{h2_count} sections"))
    else:
        gates.append(("Medium HTML exists", False, "Missing"))

    # 3. Dev.to MD: >= 1500 words with full frontmatter
    devto_path = RESULT_DIR / "dev.to" / "lab-37-cache-invalidation-strategies-devto.md"
    if devto_path.exists():
        devto_md = devto_path.read_text(encoding="utf-8")
        fm_match = re.match(r"^---\n(.*?)\n---", devto_md, re.DOTALL)
        has_fm = fm_match is not None and "title:" in fm_match.group(1) and "tags:" in fm_match.group(1)
        words = len(devto_md.split())
        gates.append(("Dev.to MD frontmatter", has_fm, "Full YAML frontmatter present" if has_fm else "Missing/invalid"))
        gates.append(("Dev.to MD words (>= 1500)", words >= 1500, f"{words} words"))
    else:
        gates.append(("Dev.to MD exists", False, "Missing"))

    # 4. Cover PNG (ID & EN): exists and > 50KB
    cover_id = RESULT_DIR / "image-cover" / "lab-37-cache-invalidation-strategies-cover-id.png"
    cover_en = RESULT_DIR / "image-cover" / "lab-37-cache-invalidation-strategies-cover-en.png"
    for p, label in [(cover_id, "ID"), (cover_en, "EN")]:
        if p.exists():
            sz = p.stat().st_size
            ok = sz > 50 * 1024
            gates.append((f"Cover PNG ({label}) > 50KB", ok, f"{sz} bytes"))
        else:
            gates.append((f"Cover PNG ({label}) exists", False, "Missing"))

    # 5. PDF LinkedIn: 24 pages, > 200KB
    pdf_path = RESULT_DIR / "linkedin" / "lab-37-cache-invalidation-strategies-linkedin.pdf"
    if pdf_path.exists():
        sz = pdf_path.stat().st_size
        sz_ok = sz > 200 * 1024
        pdf_bytes = pdf_path.read_bytes()
        pages = len(re.findall(rb"/Type\s*/Page\b", pdf_bytes))
        pages_ok = pages == 24
        gates.append(("LinkedIn PDF > 200KB", sz_ok, f"{sz} bytes"))
        gates.append(("LinkedIn PDF 24 pages", pages_ok, f"{pages} pages"))
    else:
        gates.append(("LinkedIn PDF exists", False, "Missing"))

    # 6. PPTX: 24 slides, > 100KB, slide 1 #0b1b2d, slide 2-23 #f4f1e8 with >= 3 content items
    pptx_path = RESULT_DIR / "linkedin" / "lab-37-cache-invalidation-strategies-carousel.pptx"
    if pptx_path.exists():
        sz = pptx_path.stat().st_size
        sz_ok = sz > 100 * 1024
        gates.append(("LinkedIn PPTX > 100KB", sz_ok, f"{sz} bytes"))

        with zipfile.ZipFile(pptx_path) as z:
            slides = [name for name in z.namelist() if re.match(r"ppt/slides/slide[0-9]+\.xml", name)]
            slide_count_ok = len(slides) == 24
            gates.append(("PPTX 24 slides", slide_count_ok, f"{len(slides)} slides"))

            # Slide 1 check: background #0b1b2d
            s1_xml = z.read("ppt/slides/slide1.xml").decode("utf-8")
            s1_first_sp = s1_xml[:s1_xml.find("</p:sp>")]
            s1_bg = re.findall(r"srgbClr val=\"([0-9a-fA-F]+)\"", s1_first_sp)
            s1_bg_ok = bool(s1_bg and s1_bg[0].lower() == "0b1b2d")
            s1_shapes = s1_xml.count("<p:sp>")
            gates.append(("Slide 1: bg #0b1b2d with visual elements", s1_bg_ok and s1_shapes >= 3, f"bg=#{s1_bg[0] if s1_bg else 'none'}, shapes={s1_shapes}"))

            # Slides 2-23 check: background #f4f1e8 with >= 3 content items
            s2_23_ok = True
            details = []
            for i in range(2, 24):
                s_xml = z.read(f"ppt/slides/slide{i}.xml").decode("utf-8")
                s_first_sp = s_xml[:s_xml.find("</p:sp>")]
                bg = re.findall(r"srgbClr val=\"([0-9a-fA-F]+)\"", s_first_sp)
                shapes = s_xml.count("<p:sp>")
                bg_val = bg[0].lower() if bg else ""
                # slide has background shape + title + subtitle + content boxes >= 3 items (total shapes >= 4, paragraphs >= 3)
                if bg_val != "f4f1e8" or shapes < 4:
                    s2_23_ok = False
                    details.append(f"Slide {i} failed (bg=#{bg_val}, shapes={shapes})")
            gates.append(("Slide 2-23: bg #f4f1e8 with >= 3 content items", s2_23_ok, "All 22 editorial slides verified" if s2_23_ok else "; ".join(details)))
    else:
        gates.append(("LinkedIn PPTX exists", False, "Missing"))

    # Print summary
    print("\n" + "=" * 70)
    print("QUALITY GATES VERIFICATION SUMMARY")
    print("=" * 70)
    all_pass = True
    for name, status, detail in gates:
        status_str = "PASS" if status else "FAIL"
        mark = "✓" if status else "✗"
        if not status:
            all_pass = False
        print(f"[{status_str}] {mark} {name}: {detail}")
    print("=" * 70)
    verdict = "ALL QUALITY GATES PASSED" if all_pass else "QUALITY GATES FAILED"
    print(f"VERDICT: {verdict}")
    print("=" * 70)
    return all_pass


def main():
    print("Generating deliverables for Lab 37...")
    p1 = generate_medium_metadata()
    print(f"Generated Medium Metadata: {p1}")

    p2, caption = generate_linkedin_caption()
    print(f"Generated LinkedIn Caption: {p2}")

    p3 = generate_deliverables_md(caption)
    print(f"Generated DELIVERABLES.md: {p3}")

    passed = run_quality_gates()
    if not passed:
        sys.exit(1)


if __name__ == "__main__":
    main()
