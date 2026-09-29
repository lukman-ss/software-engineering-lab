# Execution Result

## Build
Command: `go build ./...`
Result: SUCCESS (Zero build warnings/errors)

## Tests
Command: `go test -v ./...`
Result:
```text
=== RUN   TestCalculateDiscount_Strong
=== RUN   TestCalculateDiscount_Strong/VIP_customer_gets_20%_discount_and_free_shipping_regardless_of_amount
=== RUN   TestCalculateDiscount_Strong/High_spender_non-VIP_gets_20%_discount
=== RUN   TestCalculateDiscount_Strong/Premium_customer_with_amount_500_gets_10%_discount
=== RUN   TestCalculateDiscount_Strong/Standard_customer_with_amount_100_gets_5%_discount
=== RUN   TestCalculateDiscount_Strong/Coupon_with_>2_items_adds_5%_rate_boost
=== RUN   TestCalculateDiscount_Strong/Coupon_with_<=2_items_does_NOT_add_extra_rate
=== RUN   TestCalculateDiscount_Strong/Free_shipping_threshold_at_exact_200_final_amount
--- PASS: TestCalculateDiscount_Strong (0.00s)
=== RUN   TestCalculateDiscount_Weak
--- PASS: TestCalculateDiscount_Weak (0.00s)
PASS
ok  	labs/38-mutation-testing/internal/service	0.434s
=== RUN   TestEngine_GeneratesMutants
--- PASS: TestEngine_GeneratesMutants (0.00s)
=== RUN   TestEngine_MutationScoreDifference
--- PASS: TestEngine_MutationScoreDifference (0.00s)
=== RUN   TestDomainService_DirectCalculation
--- PASS: TestDomainService_DirectCalculation (0.00s)
PASS
ok  	labs/38-mutation-testing/tests	0.437s
```

## Race Detector
Command: `go test -race ./...`
Result:
```text
ok  	labs/38-mutation-testing/internal/service	1.372s
ok  	labs/38-mutation-testing/tests	1.345s
```
Status: PASS (0 data races detected)

## Demo
Command: `go run ./cmd/demo`
Result:
```text
=============================
  MUTATION TESTING LAB DEMO
=============================

--- WHAT IS MUTATION TESTING? ---
  Mutation testing proves test effectiveness by injecting small faults (mutants)
  into source code and checking if the test suite detects them.
  A 'killed' mutant = test suite found the fault.
  A 'survived' mutant = test suite MISSED a real bug pattern.

  Formula: Mutation Score = (Killed / Total) × 100%

--- TARGET: internal/service/discount.go — 15 mutants generated ---

  [RelationalOpReplace] 8 mutants
    Line 31: == → != | Replace '==' with '!='
    Line 31: >= → > | Replace '>=' with '>'
    Line 33: == → != | Replace '==' with '!='
    Line 33: >= → > | Replace '>=' with '>'
    Line 35: >= → > | Replace '>=' with '>'
    Line 40: > → >= | Replace '>' with '>='
    Line 48: == → != | Replace '==' with '!='
    Line 48: >= → > | Replace '>=' with '>'
  [BooleanOpFlip] 4 mutants
    Line 31: || → && | Replace '||' with '&&'
    Line 33: && → || | Replace '&&' with '||'
    Line 40: && → || | Replace '&&' with '||'
    Line 48: || → && | Replace '||' with '&&'
  [ArithmeticOpReplace] 2 mutants
    Line 44: * → / | Replace '*' with '/'
    Line 45: - → + | Replace '-' with '+'
  [BoundaryValueMutate] 1 mutants
    Line 40: 2 → 3 | Shift integer constant +1

--- WEAK TEST SUITE (100% Line Coverage, Weak Assertions) ---
  Total: 15, Killed: 0, Survived: 15, Score: 0.00%
  Breakdown:
    [SURVIVED] Mutant 1: Replace '||' with '&&' (|| → &&)
    [SURVIVED] Mutant 2: Replace '==' with '!=' (== → !=)
    [SURVIVED] Mutant 3: Replace '>=' with '>' (>= → >)
    [SURVIVED] Mutant 4: Replace '&&' with '||' (&& → ||)
    [SURVIVED] Mutant 5: Replace '==' with '!=' (== → !=)
    [SURVIVED] Mutant 6: Replace '>=' with '>' (>= → >)
    [SURVIVED] Mutant 7: Replace '>=' with '>' (>= → >)
    [SURVIVED] Mutant 8: Replace '&&' with '||' (&& → ||)
    [SURVIVED] Mutant 9: Replace '>' with '>=' (> → >=)
    [SURVIVED] Mutant 10: Shift integer constant +1 (2 → 3)
    [SURVIVED] Mutant 11: Replace '*' with '/' (* → /)
    [SURVIVED] Mutant 12: Replace '-' with '+' (- → +)
    [SURVIVED] Mutant 13: Replace '||' with '&&' (|| → &&)
    [SURVIVED] Mutant 14: Replace '==' with '!=' (== → !=)
    [SURVIVED] Mutant 15: Replace '>=' with '>' (>= → >)

--- STRONG TEST SUITE (Comprehensive Assertions) ---
  Total: 15, Killed: 15, Survived: 0, Score: 100.00%
  Breakdown:
    [KILLED  ] Mutant 1: Replace '||' with '&&' (|| → &&)
    [KILLED  ] Mutant 2: Replace '==' with '!=' (== → !=)
    [KILLED  ] Mutant 3: Replace '>=' with '>' (>= → >)
    [KILLED  ] Mutant 4: Replace '&&' with '||' (&& → ||)
    [KILLED  ] Mutant 5: Replace '==' with '!=' (== → !=)
    [KILLED  ] Mutant 6: Replace '>=' with '>' (>= → >)
    [KILLED  ] Mutant 7: Replace '>=' with '>' (>= → >)
    [KILLED  ] Mutant 8: Replace '&&' with '||' (&& → ||)
    [KILLED  ] Mutant 9: Replace '>' with '>=' (> → >=)
    [KILLED  ] Mutant 10: Shift integer constant +1 (2 → 3)
    [KILLED  ] Mutant 11: Replace '*' with '/' (* → /)
    [KILLED  ] Mutant 12: Replace '-' with '+' (- → +)
    [KILLED  ] Mutant 13: Replace '||' with '&&' (|| → &&)
    [KILLED  ] Mutant 14: Replace '==' with '!=' (== → !=)
    [KILLED  ] Mutant 15: Replace '>=' with '>' (>= → >)

--- CONCLUSION ---
  Weak suite kills:   0/15 mutants (0.00%)
  Strong suite kills: 15/15 mutants (100.00%)

  This demonstrates that 100% code coverage does NOT ensure
  fault detection. Only precise assertions kill mutants.
  A mutation score of 100% requires exact boundary checks,
  boolean condition verification, and value assertions.
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
