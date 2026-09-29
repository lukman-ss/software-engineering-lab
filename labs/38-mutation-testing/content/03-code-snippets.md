# Code Snippets

## Snippet 1 — AST Inspection & Mutation Plan Generation

Source File: `internal/engine/mutator.go:40-165`  
Purpose: Menelusuri Abstract Syntax Tree (AST) kode Go, mendeteksi node operator relasional, boolean, aritmetika, serta konstanta integer, lalu membentuk daftar `MutationPlan` yang dilengkapi *closure* pembalik (*undo*).

```go
ast.Inspect(parsedFile, func(n ast.Node) bool {
    if n == nil {
        return true
    }

    switch expr := n.(type) {
    case *ast.BinaryExpr:
        pos := m.Fset.Position(expr.Pos())

        // 1. Relational operators
        switch expr.Op {
        case token.GEQ:
            plans = append(plans, MutationPlan{
                Type:        RelationalOpReplace,
                Description: "Replace '>=' with '>'",
                Line:        pos.Line,
                Original:    ">=",
                Mutated:     ">",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.GTR
                    return func() { expr.Op = token.GEQ }
                },
            })
        case token.GTR:
            plans = append(plans, MutationPlan{
                Type:        RelationalOpReplace,
                Description: "Replace '>' with '>='",
                Line:        pos.Line,
                Original:    ">",
                Mutated:     ">=",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.GEQ
                    return func() { expr.Op = token.GTR }
                },
            })
        case token.EQL:
            plans = append(plans, MutationPlan{
                Type:        RelationalOpReplace,
                Description: "Replace '==' with '!='",
                Line:        pos.Line,
                Original:    "==",
                Mutated:     "!=",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.NEQ
                    return func() { expr.Op = token.EQL }
                },
            })

        // 2. Boolean operators
        case token.LOR:
            plans = append(plans, MutationPlan{
                Type:        BooleanOpFlip,
                Description: "Replace '||' with '&&'",
                Line:        pos.Line,
                Original:    "||",
                Mutated:     "&&",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.LAND
                    return func() { expr.Op = token.LOR }
                },
            })
        case token.LAND:
            plans = append(plans, MutationPlan{
                Type:        BooleanOpFlip,
                Description: "Replace '&&' with '||'",
                Line:        pos.Line,
                Original:    "&&",
                Mutated:     "||",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.LOR
                    return func() { expr.Op = token.LAND }
                },
            })

        // 3. Arithmetic operators
        case token.SUB:
            plans = append(plans, MutationPlan{
                Type:        ArithmeticOpReplace,
                Description: "Replace '-' with '+'",
                Line:        pos.Line,
                Original:    "-",
                Mutated:     "+",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.ADD
                    return func() { expr.Op = token.SUB }
                },
            })
        case token.MUL:
            plans = append(plans, MutationPlan{
                Type:        ArithmeticOpReplace,
                Description: "Replace '*' with '/'",
                Line:        pos.Line,
                Original:    "*",
                Mutated:     "/",
                Apply: func(f *ast.File) func() {
                    expr.Op = token.QUO
                    return func() { expr.Op = token.MUL }
                },
            })
        }

    case *ast.BasicLit:
        // 4. Boundary / Value Mutation for integers
        pos := m.Fset.Position(expr.Pos())
        if expr.Kind == token.INT {
            val, err := strconv.Atoi(expr.Value)
            if err == nil && val > 0 {
                orig := expr.Value
                plans = append(plans, MutationPlan{
                    Type:        BoundaryValueMutate,
                    Description: "Shift integer constant +1",
                    Line:        pos.Line,
                    Original:    orig,
                    Mutated:     strconv.Itoa(val + 1),
                    Apply: func(f *ast.File) func() {
                        expr.Value = strconv.Itoa(val + 1)
                        return func() { expr.Value = orig }
                    },
                })
            }
        }
    }

    return true
})
```

Explanation: Code ini memindai setiap node AST kode Go tanpa memodifikasi berkas asli di disk. Setiap kecocokan operator menghasilkan objek `MutationPlan` berisi fungsi `Apply` yang mengubah AST dan mengembalikan fungsi pembatalan (*undo*) untuk mengembalikan state node AST.

---

## Snippet 2 — Parallel Mutant Runner with Isolated AST Parsing

Source File: `internal/engine/runner.go:23-91`  
Purpose: Menjalankan eksekusi mutan secara paralel menggunakan `sync.WaitGroup`. Setiap *goroutine* mem-parsing ulang AST secara terisolasi untuk menghindari kontaminasi state antar-mutan.

```go
func (r *Runner) Run(sourceCode []byte, plans []MutationPlan, testFn TestFunc) Report {
    results := make([]MutantResult, len(plans))

    var wg sync.WaitGroup
    for i, plan := range plans {
        wg.Add(1)
        go func(idx int, p MutationPlan) {
            defer wg.Done()

            fset := token.NewFileSet()
            parsedFile, err := parser.ParseFile(fset, "source.go", sourceCode, 0)
            if err != nil {
                results[idx] = MutantResult{
                    Mutant: Mutant{
                        ID:          idx + 1,
                        Type:        p.Type,
                        Description: p.Description,
                        LineNumber:  p.Line,
                        Original:    p.Original,
                        Mutated:     p.Mutated,
                    },
                    Status: StatusEquivalent,
                    Output: "parse error",
                }
                return
            }

            undo := p.Apply(parsedFile)

            mutatedSrc, err := renderSource(fset, parsedFile)
            if err != nil {
                undo()
                results[idx] = MutantResult{ ... }
                return
            }

            killed := testFn(mutatedSrc)
            undo()

            status := StatusSurvived
            if killed {
                status = StatusKilled
            }

            results[idx] = MutantResult{
                Mutant: Mutant{
                    ID:          idx + 1,
                    Type:        p.Type,
                    Description: p.Description,
                    LineNumber:  p.Line,
                    Original:    p.Original,
                    Mutated:     p.Mutated,
                },
                Status: status,
            }
        }(i, plan)
    }
    wg.Wait()

    // Hitung total killed dan score...
```

Explanation: runner mengisolasi setiap mutan dalam *goroutine* independen. Karena *slice* `results` telah dialokasikan sesuai jumlah mutan (`make([]MutantResult, len(plans))`), penulisan hasil dilakukan langsung via indeks `results[idx]` sehingga bebas dari *mutex lock contention*.

---

## Snippet 3 — Domain Target: Calculate Discount Logic

Source File: `internal/service/discount.go:27-56`  
Purpose: Menyediakan logika bisnis utama untuk perhitungan tingkat diskon dan kelayakan *free shipping*. Kode ini mengandung operator relasional, boolean majemuk, aritmetika, dan konstanta batas yang menjadi target penyuntikan mutan.

```go
func CalculateDiscount(order Order) DiscountResult {
    rate := 0.0

    // Relational & boolean: VIP or high spender gets primary rate
    if order.Tier == TierVIP || order.TotalAmount >= 1000.0 {
        rate = 0.20
    } else if order.Tier == TierPremium && order.TotalAmount >= 500.0 {
        rate = 0.10
    } else if order.TotalAmount >= 100.0 {
        rate = 0.05
    }

    // Boundary & arithmetic: Coupon adds extra discount if minimum items present
    if order.HasCoupon && order.ItemCount > 2 {
        rate = rate + 0.05
    }

    discountAmount := order.TotalAmount * rate
    finalAmount := order.TotalAmount - discountAmount

    // Free shipping threshold: after discount, VIP or amount >= 200
    freeShipping := (order.Tier == TierVIP) || (finalAmount >= 200.0)

    return DiscountResult{
        OriginalTotal: order.TotalAmount,
        DiscountRate:  rate,
        DiscountTotal: discountAmount,
        FinalAmount:   finalAmount,
        FreeShipping:  freeShipping,
    }
}
```

Explanation: Logika bisnis mencakup berbagai kondisi batas. Jika salah satu kondisi ini berubah sedikit saja (misal `order.ItemCount > 2` menjadi `order.ItemCount >= 2`), perilaku sistem akan bergeser bagi pengguna dengan 2 barang.

---

## Snippet 4 — Weak Test Suite with 100% Coverage

Source File: `internal/service/discount_weak_test.go:7-51`  
Purpose: Menunjukkan skenario *test suite* yang mencapai 100.0% *statement coverage* namun hanya menggunakan asersi permusukaan (*weak assertions*).

```go
func TestCalculateDiscount_Weak(t *testing.T) {
    // Case 1: VIP customer
    res1 := CalculateDiscount(Order{
        TotalAmount: 1200.0,
        ItemCount:   3,
        Tier:        TierVIP,
        HasCoupon:   true,
    })
    if res1.FinalAmount <= 0 {
        t.Errorf("Expected positive final amount, got %v", res1.FinalAmount)
    }

    // Case 2: Premium customer with high amount
    res2 := CalculateDiscount(Order{
        TotalAmount: 600.0,
        ItemCount:   1,
        Tier:        TierPremium,
        HasCoupon:   false,
    })
    if res2.DiscountTotal < 0 {
        t.Errorf("Discount total should not be negative")
    }

    // Case 3: Standard customer with medium amount
    res3 := CalculateDiscount(Order{
        TotalAmount: 150.0,
        ItemCount:   1,
        Tier:        TierStandard,
        HasCoupon:   false,
    })
    if res3.DiscountRate < 0 {
        t.Errorf("Rate should not be negative")
    }

    // Case 4: Standard customer with low amount and coupon
    res4 := CalculateDiscount(Order{
        TotalAmount: 50.0,
        ItemCount:   4,
        Tier:        TierStandard,
        HasCoupon:   true,
    })
    if res4.OriginalTotal != 50.0 {
        t.Errorf("Original total mismatch")
    }
}
```

Explanation: Keempat kasus di atas melewati seluruh baris fungsi `CalculateDiscount`. Namun karena asersinya hanya mengecek `FinalAmount <= 0` atau `DiscountRate < 0`, mutasi logika (seperti membalik `||` menjadi `&&` atau mengubah diskon `+` menjadi `-`) tidak memicu kegagalan tes.

---

## Snippet 5 — Strong Table-Driven Assertions for Boundary & Logic Verification

Source File: `internal/service/discount_strong_test.go:6-176`  
Purpose: Menunjukkan *test suite* berkualitas tinggi dengan asersi eksak terhadap nilai persentase, total diskon, jumlah akhir, dan flag *free shipping* di berbagai kondisi batas.

```go
func TestCalculateDiscount_Strong(t *testing.T) {
    tests := []struct {
        name                 string
        order                Order
        expectedRate         float64
        expectedDiscount     float64
        expectedFinal        float64
        expectedFreeShipping bool
    }{
        {
            name: "VIP customer gets 20% discount and free shipping regardless of amount",
            order: Order{
                TotalAmount: 100.0,
                ItemCount:   1,
                Tier:        TierVIP,
                HasCoupon:   false,
            },
            expectedRate:         0.20,
            expectedDiscount:     20.0,
            expectedFinal:        80.0,
            expectedFreeShipping: true,
        },
        {
            name: "High spender non-VIP gets 20% discount",
            order: Order{
                TotalAmount: 1000.0,
                ItemCount:   1,
                Tier:        TierStandard,
                HasCoupon:   false,
            },
            expectedRate:         0.20,
            expectedDiscount:     200.0,
            expectedFinal:        800.0,
            expectedFreeShipping: true,
        },
        // ... (11 total skenario mencakup ambang batas 500.0, 100.0, ItemCount 2 vs 3, 499.99, 99.99)
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            res := CalculateDiscount(tt.order)
            if res.DiscountRate != tt.expectedRate {
                t.Errorf("DiscountRate: got %.2f, want %.2f", res.DiscountRate, tt.expectedRate)
            }
            if res.DiscountTotal != tt.expectedDiscount {
                t.Errorf("DiscountTotal: got %.2f, want %.2f", res.DiscountTotal, tt.expectedDiscount)
            }
            if res.FinalAmount != tt.expectedFinal {
                t.Errorf("FinalAmount: got %.2f, want %.2f", res.FinalAmount, tt.expectedFinal)
            }
            if res.FreeShipping != tt.expectedFreeShipping {
                t.Errorf("FreeShipping: got %v, want %v", res.FreeShipping, tt.expectedFreeShipping)
            }
        })
    }
}
```

Explanation: Tabel pengujian ini secara eksplisit menguji nilai tepat di titik batas dan memeriksa seluruh properti *struct* keluaran secara terpisah. Ketika mutan disuntikkan ke kode aplikasi, sekecil apa pun perubahannya akan memicu kegagalan pada salah satu asersi eksplisit ini.
