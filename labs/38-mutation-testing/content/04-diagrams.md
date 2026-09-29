# Diagrams

## 1. Perbedaan Mendasar: Code Coverage vs Mutation Testing

```text
┌─────────────────────────────────────────────────────────────┐
│                 Kode Sumber (discount.go)                   │
│                                                             │
│   if order.Tier == TierVIP || order.TotalAmount >= 1000.0   │
└──────────────────────────────┬──────────────────────────────┘
                               │
            ┌──────────────────┴──────────────────┐
            ▼                                     ▼
┌───────────────────────────────┐   ┌───────────────────────────────┐
│     Traditional Coverage      │   │       Mutation Testing        │
├───────────────────────────────┤   ├───────────────────────────────┤
│ Pertanyaan:                   │   │ Pertanyaan:                   │
│ "Apakah baris ini dieksekusi  │   │ "Jika baris ini dirusak,      │
│  oleh pengujian?"             │   │  apakah pengujian gagal?"     │
│                               │   │                               │
│ Evaluasi:                     │   │ Evaluasi:                     │
│ Test suite berjalan ->        │   │ Mutasi '>=' menjadi '>'       │
│ Statement hit: YES (100%)     │   │ Weak Test: Lolos (SURVIVED)   │
│                               │   │ Strong Test: Gagal (KILLED)   │
│                               │   │                               │
│ Metrik: Statement Hit Ratio   │   │ Metrik: Mutation Score        │
└───────────────────────────────┘   └───────────────────────────────┘
```

---

## 2. Alur Eksekusi Mesin Mutasi AST (Lab Architecture)

```text
       ┌────────────────────────┐
       │   discount.go (Source) │
       └───────────┬────────────┘
                   │
                   ▼
       ┌────────────────────────┐
       │     go/parser AST      │
       └───────────┬────────────┘
                   │
                   ▼
       ┌────────────────────────┐
       │       ASTMutator       │
       │ (Inspect BinaryExpr /  │
       │    BasicLit Nodes)     │
       └───────────┬────────────┘
                   │ Menghasilkan 15 Mutation Plans
                   ▼
┌──────────────────────────────────────────────────┐
│                   Runner.Run                     │
│                                                  │
│  Goroutine 1: Mutant 1 (|| -> &&)                │
│  Goroutine 2: Mutant 2 (== -> !=)                │
│  Goroutine 3: Mutant 3 (>= -> >)                 │
│  ...                                             │
│  Goroutine 15: Mutant 15 (>= -> >)               │
│                                                  │
│  Setiap goroutine:                               │
│  1. Parse FileSet baru (terisolasi)              │
│  2. Apply plan & render mutated source           │
│  3. Eksekusi TestFunc(mutatedSrc)                │
│  4. Tentukan Status: KILLED atau SURVIVED        │
└──────────────────────────┬───────────────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │      Report Engine      │
              ├─────────────────────────┤
              │ Total Mutants: 15       │
              │ Killed:        k        │
              │ Survived:      s        │
              │ Score: (k/15) * 100%    │
              └─────────────────────────┘
```

---

## 3. Model RIP (Reach, Infect, Propagate) pada Evaluasi Mutan

```text
                        ┌──────────────────┐
                        │   Tes Berjalan   │
                        └────────┬─────────┘
                                 │
                                 ▼
                     ┌───────────────────────┐
                     │ 1. REACH              │  Apakah tes mengeksekusi
                     │    Baris Termutasi?   │  baris kode yang dimutasi?
                     └───────────┬───────────┘
                                 │ Ya
                                 ▼
                     ┌───────────────────────┐
                     │ 2. INFECT             │  Apakah mutasi mengubah
                     │    State Internal?    │  nilai variabel program?
                     └───────────┬───────────┘
                                 │ Ya
                                 ▼
                     ┌───────────────────────┐
                     │ 3. PROPAGATE          │  Apakah nilai salah tersebut
                     │    Sampai ke Asersi?  │  diperiksa oleh asersi tes?
                     └───────────┬───────────┘
                                 │
                 ┌───────────────┴───────────────┐
                 │ Ya                            │ Tidak (Asersi Lemah)
                 ▼                               ▼
      ┌─────────────────────┐         ┌─────────────────────┐
      │     Test GAGAL      │         │     Test LOLOS      │
      │   (Mutant KILLED)   │         │  (Mutant SURVIVED)  │
      └─────────────────────┘         └─────────────────────┘
```

---

## 4. Perbandingan Hasil Uji Nyata (Lab 38 Execution Result)

```text
Coverage & Mutation Score Comparison

Weak Test Suite:
Line Coverage  [████████████████████] 100.0% (4/4 Statement Blocks)
Mutation Score [                    ]   0.0% (0/15 Mutants Killed)

Strong Test Suite:
Line Coverage  [████████████████████] 100.0% (4/4 Statement Blocks)
Mutation Score [████████████████████] 100.0% (15/15 Mutants Killed)
```
