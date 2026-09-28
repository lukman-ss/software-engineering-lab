# Research Report: Saga Pattern — Mengelola Transaksi Terdistribusi Tanpa 2PC

## Research Question

Bagaimana menjaga konsistensi data di beberapa service/database berbeda ketika 2PC (Two-Phase Commit) terlalu lambat atau tidak memungkinkan?

## Executive Summary

Saga Pattern adalah solusi untuk transaksi terdistribusi di microservices tanpa menggunakan 2PC. Pattern ini memecah transaksi terdistribusi menjadi series local transactions, dengan compensating transactions sebagai mekanisme rollback. Dua pendekatan utama: Choreography (event-driven, loose coupling) dan Orchestration (centralized coordinator). Kelemahan utama: tidak ada built-in isolation, sehingga memerlukan countermeasures untuk mencegah data anomalies.

## Findings

### Finding 1: 2PC Tidak Feasible untuk Microservices Database-Per-Service

**Claim:** Traditional ACID transactions dan 2PC tidak applicable untuk multiple independently managed data stores

**Evidence:** Microsoft Azure dan Chris Richardson konsisten menyatakan bahwa karena database per microservice di-scale independently dan isolated, 2PC tidak bisa digunakan, dan saga menjadi solusinya

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Chris Richardson: https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

---

### Finding 2: Saga = Sequence of Local Transactions dengan Atomicity di Level Saga

**Claim:** Saga menjamin atomicity di level saga, bukan di level individual transaction

**Evidence:** "The Saga pattern manages transactions by breaking them into a sequence of local transactions. Each local transaction: 1. Completes its work atomically within a single service. 2. Updates the service's database. 3. Initiates the next transaction via an event or message."

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Chris Richardson: https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

---

### Finding 3: Dua Pendekatan Implementasi dengan Trade-off Berbeda

**Claim:** Choreography dan orchestration memiliki kelebihan dan kelemahan masing-masing

**Evidence:**
| Aspect | Choreography | Orchestration |
|--------|-------------|---------------|
| Complexity | Simple, loose coupling | More complex, centralized |
| Scalability | Better (no SPOF) | Worse (SPOF risk) |
| Debugging | Difficult (spaghetti events) | Easier (central flow) |
| Cyclic Dependencies | Risk of cyclic deps | No cyclic deps |
| Integration Testing | Difficult | Easier |

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Chris Richardson: https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

---

### Finding 4: Compensating Transactions Bukan Simple Rollback

**Claim:** Compensating transactions adalah business-level corrective actions, bukan automatic rollback

**Evidence:** "lack of automatic rollback - a developer must design compensating transactions that explicitly undo changes made earlier in a saga rather than relying on the automatic rollback feature of ACID transactions"

**Sources:**
- Chris Richardson: https://microservices.io/patterns/data/saga.html
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

### Finding 5: Tiga Jenis Transaksi dalam Saga

**Claim:** Saga terdiri dari Compensable, Pivot, dan Retryable transactions

**Evidence:**
- Compensable: bisa di-undo dengan opposite effect
- Pivot: point of no return
- Retryable: idempotent, dijamin selesai

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

### Finding 6: Tidak Ada Built-in Isolation

**Claim:** Sagas tidak memiliki isolation, yang bisa menyebabkan data anomalies

**Evidence:** "Because each service manages its own data, called participant data, there's no built-in isolation across services."

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Chris Richardson: https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

---

### Finding 7: Countermeasures untuk Anomali Data

**Claim:** Enam countermeasures tersedia: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency

**Evidence:** Microsoft Azure Architecture Center merekomendasikan keenam countermeasures ini untuk mencegah data anomalies

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

### Finding 8: Idempotency Requirement

**Claim:** Transaction dalam saga harus idempotent

**Evidence:** "The system must handle transient failures effectively and ensure idempotence, when repeating the same operation doesn't alter the outcome"

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- Chris Richardson: https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

---

### Finding 9: Limitations - Compensating Transactions Mungkin Gagal

**Claim:** Compensating transactions tidak selalu berhasil

**Evidence:** "Limitations of compensating transactions: Compensating transactions might not always succeed, which can leave the system in an inconsistent state"

**Sources:**
- Microsoft Azure: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Areas of Agreement

1. **Saga definition:** Both sources agree on core definition
2. **Two approaches:** Both agree on choreography vs orchestration
3. **Lack of isolation:** Both identify this as key limitation
4. **Compensating transactions:** Both explain need for business-level corrective actions
5. **Idempotency:** Both emphasize importance of idempotency

## Areas of Disagreement

No material disagreements found.

Minor complementary differences:
- Microsoft source provides more detail on countermeasures
- Chris Richardson source emphasizes atomic update + publish requirement

## Limitations

1. Original academic paper (Garcia-Molina & Salem, 1987) may be behind paywall
2. Some URLs resulted in 404 (infoq, kylewbanks, eventuate)
3. Focus on implementation patterns vs theoretical foundations
4. Limited coverage of saga in specific technology stacks

## Conclusion

Saga Pattern efektif menggantikan 2PC untuk microservices dengan database-per-service architecture. Pattern ini memberikan atomicity di level saga melalui sequence of local transactions dan compensating transactions, meskipun dengan trade-off: tidak ada built-in isolation. Dua pendekatan utama (choreography dan orchestration) memiliki trade-off berbeda yang harus dipilih berdasarkan complexity workflow. Compensating transactions harus didesain secara manual dan idempotent, dengan pemahaman bahwa ada risiko compensating transactions gagal.
