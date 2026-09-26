# Claim Audit

## Claim 1

Claim: Deadlock terjadi hanya jika empat kondisi Coffman terpenuhi secara bersamaan (mutual exclusion, hold and wait, no preemption, circular wait).
Location: 05-report.md (Finding 1)
Evidence Provided: Wikipedia & PostgreSQL explicit locking docs.
Source: https://en.wikipedia.org/wiki/Deadlock_(computer_science)#Conditions
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-established theoretical foundation.

## Claim 2

Claim: PostgreSQL otomatis mendeteksi deadlock dan membatalkan satu transaksi (victim), pemilihan transaksi victim tidak dapat diprediksi.
Location: 05-report.md (Finding 2)
Evidence Provided: PostgreSQL Documentation 13.3.4.
Source: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: PostgreSQL docs confirm this behavior directly.

## Claim 3

Claim: PostgreSQL mengembalikan SQLSTATE 40P01 untuk deadlock dan 40001 untuk serialization failure; aplikasi harus meretry transaksi.
Location: 05-report.md (Finding 3)
Evidence Provided: Appendix A Error Codes & Transaction Isolation.
Source: https://www.postgresql.org/docs/current/errcodes-appendix.html
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Error codes and recommended retry behavior confirmed.

## Claim 4

Claim: deadlock_timeout adalah interval sebelum deteksi deadlock dijalankan, sedangkan lock_timeout adalah durasi maksimum tunggu lock.
Location: 05-report.md (Finding 4)
Evidence Provided: Runtime config documentation.
Source: https://www.postgresql.org/docs/current/runtime-config-locks.html
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Important distinction accurately detailed.

## Claim 5

Claim: Konsistensi urutan perolehan lock (lock ordering) di seluruh transaksi mencegah terbentuknya circular wait.
Location: 05-report.md (Finding 5)
Evidence Provided: PostgreSQL Documentation 13.3.4.
Source: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-supported standard prevention technique.

## Claim 6

Claim: Operasi eksternal (API WhatsApp, generate PDF) dalam transaksi memperpanjang waktu penguncian dan memperbesar peluang deadlock.
Location: 05-report.md (Finding 6)
Evidence Provided: Lab case study and generic best practices.
Source: Topic spec / Application patterns.
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: LOW
Notes: Correct application engineering principle derived from duration risk.

## Claim 7

Claim: Retry dengan backoff dan jitter adalah praktik standar untuk penanganan deadlock transaksi idempoten di Go.
Location: 05-report.md (Finding 8)
Evidence Provided: Go standard library pkg/time docs.
Source: https://pkg.go.dev/time
Source Actually Supports Claim: PARTIAL
Classification: EXAMPLE
Severity: MEDIUM
Notes: Standard pattern, though standard library does not specifically mandate retry patterns for deadlocks.

## Claim 8

Claim: Monitoring deadlock production mengandalkan log_lock_waits, pg_locks, dan pg_stat_activity.
Location: 05-report.md (Finding 9)
Evidence Provided: PostgreSQL monitoring and logging docs.
Source: https://www.postgresql.org/docs/current/runtime-config-logging.html
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standard PostgreSQL observability practices confirmed.
