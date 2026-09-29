# Diagrams: Property-Based Testing (PBT)

Diagram arsitektur dan alur eksekusi Property-Based Testing yang diturunkan langsung dari implementasi lab.

---

## 1. Alur Kerja Engine Property-Based Testing vs Example-Based Testing

```text
==================================================================================
A. EXAMPLE-BASED TESTING (Unit Test Konvensional)
==================================================================================

  +-----------------------+     +-----------------------+     +------------------+
  | Hand-picked Inputs    | --> | Function Under Test   | --> | Assert Result    |
  | (e.g. $10.50, $99.99) |     | f(x)                  |     | == Expected      |
  +-----------------------+     +-----------------------+     +------------------+
                                                                       |
                                                                       v
                                                             [ PASS (False Confidence) ]


==================================================================================
B. PROPERTY-BASED TESTING & SHRINKING PIPELINE
==================================================================================

  +-----------------------+
  | Generator (Biased)    | <--- Generator Bias (0, Negatives, Boundaries)
  | e.g. testing/quick    |
  +-----------------------+
              |
              v (Random Input Candidates)
  +-----------------------+
  | Universal Invariant   |
  | Check (100-1000x)     |
  +-----------------------+
        |             |
        | PASS        | FAIL (Invariant Violated)
        v             v
  [ SUCCESS ]   +-------------------------------------------------------+
                | COUNTEREXAMPLE SHRINKING ENGINE                       |
                |                                                       |
                | 1. Binary Sectioning (Try Left/Right Half)            |
                | 2. Element Elimination (Remove index i)               |
                | 3. Magnitude Reduction (Divide by 2 / set to 0)       |
                +-------------------------------------------------------+
                                           |
                                           v
                                [ MINIMAL REPRODUCER ]
                                  (e.g. [-1] instead of 10-elem array)
```

---

## 2. Peta Invariant Kanonikal dan Contoh Terverifikasi

```text
+-----------------------------------------------------------------------------------+
|                        CANONICAL PROPERTY-BASED INVARIANTS                        |
+-----------------------------------------------------------------------------------+
|                                                                                   |
|  1. ROUNDTRIP INVARIANT                                                           |
|     +--------------+     Format()     +--------------+     Parse()     +--------+ |
|     | RobustAmount | ---------------> | String "$..."| --------------> | Amount | |
|     +--------------+                  +--------------+                 +--------+ |
|            |                                                               |      |
|            +---------------------- ASSERT EQUAL ---------------------------+      |
|                                                                                   |
|  2. IDEMPOTENCE INVARIANT                                                         |
|     +--------------+   RobustMerge()  +--------------+  RobustMerge()  +--------+ |
|     | Interval Set | ---------------> | Merged Set A | --------------> | Set B  | |
|     +--------------+                  +--------------+                 +--------+ |
|                                              |                             |      |
|                                              +----- ASSERT EQUAL (A == B) -+      |
|                                                                                   |
|  3. EQUIVALENCE / TEST ORACLE INVARIANT                                           |
|                        +------------------------------+                           |
|                        | Input Slice (Unsorted/Random)|                           |
|                        +------------------------------+                           |
|                               /                \                                  |
|               NaiveMerge(x)  /                  \  RobustMerge(x)                 |
|                             v                    v                                |
|                     +---------------+    +-----------------+                      |
|                     | Naive Output  |    | Robust Output   |                      |
|                     +---------------+    +-----------------+                      |
|                             \                    /                                |
|                              \-- ASSERT EQUAL --/                                 |
|                         (Fails on Unsorted -> 85/100 Discrepancy)                |
+-----------------------------------------------------------------------------------+
```
