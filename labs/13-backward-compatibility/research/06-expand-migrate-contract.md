# The Expand -> Migrate -> Contract Pattern

## Overview

Juga dikenal sebagai **Parallel Change Pattern** (dirumuskan oleh Joshua Kerievsky dan dipopulerkan oleh Danilo Sato / Martin Fowler, 2014). Pola ini membagi *breaking incompatible change* menjadi tahapan non-breaking yang incremental:

```text
EXPAND  -->  MIGRATE  -->  CONTRACT
```

**Evidence** (Martin Fowler, Parallel Change):
"**Parallel change**, also known as **expand and contract**, is a pattern to implement backward-incompatible changes to an interface in a safe manner, by breaking the change into three distinct phases: expand, migrate, and contract."

## Phase 1: Expand

Tingkatkan sistem sehingga mendukung **kedua format lama dan baru** secara bersamaan.

### 1.1 Kiat Implementasi

1. Tambah kolom atau tabel baru di database dengan constraint yang rileks (nullable).
2. Tambah field/endpoint baru ke aplikasi.
3. Deploy perubahan. Consumer yang ada tidak menyadari perubahan dan tidak crash.

### 1.2 Evidence Kode (Martin Fowler, Grid Contoh)

```java
// BEFORE (versi lama)
class Grid {
  private Cell[][] cells;
  public void addCell(int x, int y, Cell cell) {
    cells[x][y] = cell;
  }
}

// EXPAND: support both
class Grid {
  private Cell[][] cells;                     // legacy
  private Map<Coordinate, Cell> newCells;     // baru
  public void addCell(int x, int y, Cell cell) {
    cells[x][y] = cell;
  }
  public void addCell(Coordinate coordinate, Cell cell) {
    newCells.put(coordinate, cell);
  }
}
```

**Evidence** (Fowler):
"In the *expand* phase you augment the interface to support both the old and the new versions. In our example, we introduce a new `Map<Coordinate, Cell>` data structure and the new methods that can receive `Coordinate` instances without changing the existing code."

## Phase 2: Migrate

Secara gradual transisikan semua data dan consumer dari kontrak lama ke kontrak baru.

### 2.1 Langkah Migrate

1. **Application Dual Writing**: Aplikasi menulis data ke kedua struktur lama dan baru.
2. **Backfill**: Script background atau worker menyalin semua data historis dari struktur lama ke struktur baru.
3. **Dual / Fallback Reading**: Aplikasi mulai membaca dari struktur baru, dengan fallback ke struktur lama jika kosong.
4. **Consumer Migration**: Update klien eksternal/internal, worker, dan mobile app untuk konsumsi kontrak baru.
5. **Observability Verification**: Monitor metrics hingga semua klien sepenuhnya berpindah dan permintaan read/write pada skema lama berhenti.

### 2.2 Evidence (Martin Fowler)

"During the *migrate* phase you update all clients using the old version to the new version. This can be done incrementally and, in the case of external clients, this will be the longest phase."

"This pattern is particularly useful when practicing ContinuousDelivery because it allows your code to be released in any of these three phases."

### 2.3 Penggunaan FeatureFlag selama Migrate

**Evidence** (Parallel Change):
"During the migrate phase, a [FeatureFlag] can be used to control which version of the interface is used. A feature toggle on the client side allows it to be forward-compatible with the new version of the supplier, which decouples the release of the supplier from the client."

## Phase 3: Contract

Hapus format lama setelah tidak ada yang bergantung padanya.

### 3.1 Langkah Contract

1. Berhenti menulis ke skema/field database lama.
2. Hapus code paths yang deprecated dan transformasi adapter lama.
3. Bersihkan database: hapus kolom/tabel lama atau endpoint yang deprecated.
4. Tambahkan hard constraints yang diinginkan (misal `NOT NULL` pada kolom baru) jika diperlukan.

### 3.2 Evidence Kode (Martin Fowler)

```java
// CONTRACT: remove old version
class Grid {
  private Map<Coordinate, Cell> cells;
  public void addCell(Coordinate coordinate, Cell cell) {
    cells.put(coordinate, cell);
  }
  public Cell fetchCell(Coordinate coordinate) {
    return cells.get(coordinate);
  }
  public boolean isEmpty(Coordinate coordinate) {
    return !cells.containsKey(coordinate);
  }
}
```

**Evidence** (Fowler):
"Once all usages have been migrated to the new version, you perform the *contract* phase to remove the old version and change the interface so that it only supports the new version."

### 3.3 Risiko Jika Contract Dilewati

**Evidence** (Fowler):
"The downside of using parallel change is that during the migrate phase the supplier has to support two different versions, and clients could get confused about which version is new versus old. If the contract phase is not executed you might end up in a worse state than you started, therefore you need discipline to finish the transition successfully."

## Aplikasi Pola Parallel Change

**Evidence** (Martin Fowler):
- **Refactoring**: when changing a method or function signature, especially when doing a Long Term Refactoring or when changing a PublishedInterface.
- **Database refactoring**: this is a key component to evolutionary database design. Most database refactorings follow the parallel change pattern, where the migrate phase is the transition period between the original and the new schema, until all database access code has been updated to work with the new schema.
- **Deployments**: deployment techniques such as canary releases and BlueGreenDeployment are applications of the parallel change pattern where you have both old and new versions of the code deployed side by side, and you incrementally migrate users from one version to another.
- **Remote API evolution**: parallel change can be used to evolve a remote API (e.g. a REST web service) when you can't make the change in a backwards compatible manner.

## Urutan Deployment Aman (Hipotesis Terentagn)

**Evidence validasi** (lihat 07-deployment-and-rollback.md untuk detail):

1. Add new schema
2. Deploy application compatible dengan schema lama + baru
3. Backfill existing data
4. Mulai write ke struktur baru
5. Pindahkan read ke struktur baru
6. Monitor legacy usage
7. Stop legacy writes
8. Remove legacy structure

**Status**: **VALIDATED** dengan berjalan Fowler Blue-Green Deployment ("first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that ... then deploy the new version of the application").

---

## Catatan Kualitas Bukti

| Pernyataan | Sumber | Confidence |
|------------|--------|------------|
| Expand = support keduanya tanpa merusak | Martin Fowler (2014) | HIGH |
| Migrate = fase terpanjang, incremental | Martin Fowler (2014) | HIGH |
| FeatureFlag dapat kontrol versi interface | Martin Fowler (2014) | HIGH |
| Contract = hapus lama setelah migrasi | Martin Fowler (2014) | HIGH |