# 05 - Code Audit: Bloom Filters Research

## Scope Override
Audit ini difokuskan khusus pada fase riset (`research/`). Audit implementasi kode/source code tidak dieksekusi pada tahapan ini sesuai dengan pipeline override instruksi.

## Evaluation of Pseudocode & Formulas in Research Report
1. **Index Calculation Formula**:
   $$g_i(x) = (h_1(x) + i \cdot h_2(x)) \pmod m$$
   Formulasi ini tepat dan sesuai dengan konvensi Kirsch-Mitzenmacher.
2. **Optimal Hash Count Formula**:
   $$k = \frac{m}{n} \ln 2$$
   Formulasi tepat dan konsisten dengan turunan kalkulus $\frac{d}{dk} (1 - e^{-kn/m})^k = 0$.
3. **Optimal Bit Size Formula**:
   $$m = - \frac{n \ln p}{(\ln 2)^2}$$
   Formulasi tepat dan menghasilkan nilai teoretis minimum untuk bit array.
