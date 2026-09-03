# Lab 10 — Project Estimation

**Estimasi Proyek Berbasis Breakdown, Risk, dan Uncertainty**

> Estimasi software bukan menebak lama coding. Estimasi adalah proses memecah scope, mengukur effort, dependency, risk, uncertainty, dan menghasilkan range yang dapat dijelaskan.

---

## Problem

Junior engineer sering menganggap estimasi sebagai "berapa hari ini fitur selesai". Pendekatan ini:

1. **Single-point estimate** - hanya satu angka (misal: "5 hari")
2. **Tanpa risk analysis** - tidak mempertimbangkan ketidakpastian
3. **Tanpa breakdown** - langsung menebak seluruh fitur
4. **Sangat sensitive** terhadap requirement change

Hasil: Estimasi tidak pernah akurat, deadline sering terlewat, tim frustrasi.

---

## Cara Berpikir Junior

```go
// Mental model: 10 halaman * 1 hari = 10 hari
func EstimateByPageCount(pageCount int) int {
    return pageCount * 1
}
```

```
10 halaman dokumentasi
× 1 hari per halaman
= 10 hari

Tidak ada kalkulasi:
- Risk
- Uncertainty
- Multiple engineer
- Availability
- Contingency
```

---

## Cara Berpikir Senior

```
Skill Breakdown
↓
Risk Assessment  
↓
Range Estimation (Min/MostLikely/Max)
↓
Spike for Unknowns
↓
PERT Calculation
↓
Calendar Duration (Effort ≠ Duration)
↓
Contingency Buffer
↓
Assumptions Documentation
↓
Confidence Level
```

---

## Task Breakdown

Representasikan feature menjadi task-task kecil yang dapat diukur.

### Contoh: Aplikasi Booking Servis

| No | Feature | Activities |
|----|---------|------------|
| 1 | Login | Backend/API, Validation, Testing |
| 2 | Booking Online | Database, Backend/API, Validation, Testing |
| 3 | Pilih Cabang | Frontend, Backend/API |
| 4 | Pilih Mekanik | Frontend, Backend/API |
| 5 | Payment Gateway | Backend/API, External API, Testing, Spike |
| 6 | WhatsApp Notification | Backend/API, External API, Spike |
| 7 | Dashboard Admin | Frontend, Backend/API, Testing |
| 8 | Laporan Excel | Backend/API, Database |

Setiap task dapat dijabarkan lebih lanjut:

```
Requirement Analysis
Database Schema
Backend/API
Frontend
Validation
Testing
Code Review
Deployment
UAT
```

---

## Risk & Uncertainty

| Risk Level | Kriteria | Contoh |
|------------|----------|--------|
| **Low** | Well-known, proven tech | Tambah field, CRUD biasa |
| **Medium** | Complexity menengah, butanya ada integrasi | API endpoint baru, query kompleks |
| **High** | Dependencies kunci, high cost of failure | Database migration besar, cache strategy |
| **Unknown** | Belum pernah dibuat, teknologi baru | Payment gateway baru, API vendor |

> **Unknown berarti ada sesuatu yang belum cukup dipahami untuk dipercaya sebagai implementation estimate final. Spike adalah mekanisme eksplisit untuk mengurangi uncertainty.**

---

## Unknown != Zero

Unknown tidak boleh dianggap 0 effort. Gunakan **Spike**.

```go
task := Task{
    Name: "Payment Gateway Vendor API",
    Estimate: EstimateRange{Min: 0, MostLikely: 0, Max: 0},
    Risk: RiskUnknown,
    SpikeDays: 1.5, // Wajib untuk eksplorasi
    Assumptions: []string{"Vendor sandbox available"},
}
```

**Spike** = waktu yang diperlukan untuk mengungkap ketidakpastian.

---

## Spike

Spike adalah eksplorasi teknis untuk mengatasi unknown.

### Alur:

```
Payment Gateway Vendor API
↓
Unknown Risk (0 effort)
↓
Spike: 1-2 days
↓
Baru dapat estimasi implementasi:
- Min: 3 days
- MostLikely: 5 days  
- Max: 8 days (jika ada masalah integrasi)
```

Spike harus tercatat secara eksplisit di estimation model.

---

## Estimate Range

Gunakan 3-point estimation untuk mewakili ketidakpastian.

```go
type EstimateRange struct {
    Min        float64  // Best case
    MostLikely float64  // Mode
    Max        float64  // Worst case
}
```

### PERT Formula

Expected = (Min + 4 × MostLikely + Max) / 6

```go
func (r EstimateRange) Expected() float64 {
    return (r.Min + 4.0*r.MostLikely + r.Max) / 6.0
}
```

> Ini adalah **model estimasi**, bukan jaminan tanggal selesai.

---

## Effort vs Calendar Duration

Bedakan effort dengan durasi kalender.

```
Effort: 18–24 engineer-days
↓
Spike: 2.5 days
↓
Base Effort: 20.5–26.5 engineer-days
↓
Contingency: 15%
↓
Final Effort: 23–30 engineer-days

1 engineer
70% productivity

Calendar Duration: 33–43 working days (~7–9 weeks)
```

```go
type DurationRange struct {
    MinDays      float64
    ExpectedDays float64
    MaxDays      float64
}

// Calendar Duration = Effort / (Engineers × Availability)
calendarDays = finalEffort / (engineerCount * availability)
```

**Jangan menyamakan effort dengan tanggal selesai!**

---

## Contingency

Buffer berdasarkan risiko project.

### Default Buffer:

| Risk Level | Contingency |
|------------|-------------|
| High/Unkown| 25%         |
| Medium     | 15%         |
| Low        | 10%         |

**Perbedaan Auto vs Manual Contingency:**

- `AutoContingency = true`: derive otomatis dari risk level
- `AutoContingency = false` + `ContingencyRate = 0`: NO buffer (intentional)

---

## Assumptions

Estimasi tanpa assumptions akan punya confidence rendah.

```go
project := Project{
    Assumptions: []string{
        "UI design final",
        "Vendor sandbox API accessible",
        "1 engineer dedicated with 70% allocation",
        "No major requirement changes",
    },
}
```

---

## Confidence Level

| Confidence | Conditions |
|------------|------------|
| **High**   | Semua task low/medium risk, assumptions documented, no unknown dependencies |
| **Medium** | Ada task medium risk, contingency dihitung, beberapa assumptions |
| **Low**    | Ada unknown risk, tidak ada spike, assumptions tidak lengkap |

---

## Failure Scenarios

| No | Scenario | Expected Behavior |
|----|----------|-------------------|
| 1 | Invalid range: min > mostLikely | Error |
| 2 | Invalid range: mostLikely > max | Error |
| 3 | Negative estimate | Error |
| 4 | Negative spike days | Error |
| 5 | Availability = 0 | Error |
| 6 | Availability > 1 | Error |
| 7 | Negative contingency | Error |
| 8 | Empty project | Error |
| 9 | Invalid risk level (e.g., "Critical") | Error |
| 10 | Unknown risk without SpikeDays > 0 | Error |
| 11 | EngineerCount = 0 | Error |
| 12 | EngineerCount < 0 | Error |
| 13 | Unknown risk with SpikeDays = 0 but implementation estimate > 0 | **Allowed** - spike reduces uncertainty |

---

## Cara Menjalankan Lab

```bash
cd labs/10-project-estimation
go test -v ./...
```

---

## Expected Result

```
Project: Aplikasi Booking Servis
Implementation Effort: 17.0–49.0 engineer-days
Spike Effort: 2.5 days
Base Effort: 19.5–51.5 engineer-days
Contingency: 8.2 days
Final Effort: 24.4–64.4 engineer-days
Calendar Duration: 35–92 working days (7.0–18.4 weeks)
Risk: High, Confidence: Low
```

Semua test PASS (27 tests):

```
=== RUN   TestNaiveEstimator
--- PASS: TestNaiveEstimator (0.00s)
... (semua test passing)
PASS
```

---

## Senior Engineer Takeaways

1. **Estimasi bukan prediction, tapi range-based analysis** - Selalu berikan batas bawah, yang paling mungkin, dan batas atas
2. **Unknown != 0** - Jika tidak yakin, buat spike dulu
3. **Effort ≠ Calendar Duration** - Account availability (70% productivity is realistic)
4. **Risk profile menentu contingency** - High/Unknown = 25%, Medium = 15%, Low = 10%
5. **AutoContingency vs Manual** - 0 = intentional no-buffer, not magic value
6. **Range ordering matter** - Min ≤ Expected ≤ Max (PERT formula)
7. **Validation is explicit** - Invalid inputs return errors, not silent failures
8. **EngineerCount must be ≥ 1** - Invalid planning input = explicit error
9. **Unknown with SpikeDays = 0 + non-zero estimate = VALID** - Spike helps, but implementation estimate can exist
10. **Deterministic testing** - Core estimation logic tidak perlu network atau DB

---

## Navigasi

- **Previous**: [Lab 09 — Code Review](../09-code-review/)
- **Next**: [Lab 11 — Pessimistic Locking](../11-pessimistic-locking/)