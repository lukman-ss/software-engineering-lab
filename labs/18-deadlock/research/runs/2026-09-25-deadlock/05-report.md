# Research Report

## Research Question
Apa itu deadlock, bagaimana database engine menanganinya, dan bagaimana aplikasi dapat mengurangi frekuensi dan dampak dari deadlock?

## Executive Summary
Deadlock adalah kondisi sistemik saat dua atau lebih transaksi saling memegang lock yang dibutuhkan transaksi lain, membentuk siklus tunggu permanen. Database modern memonitor hal ini secara berkala dan akan secara sepihak membatalkan (abort) salah satu transaksi sebagai "deadlock victim". Solusi paling efektif berada di level aplikasi, bukan sekadar tuning infrastruktur database. Memastikan urutan akses yang seragam, durasi transaksi yang pendek, serta implementasi retry otomatis adalah standar yang harus diterapkan oleh software engineer.

## Findings

### Finding 1
Claim: Deadlock adalah siklus dependensi (circular wait) permanen.

Evidence: Konkurensi menyebabkan transaksi A menunggu transaksi B, sementara transaksi B bergantung pada lock transaksi A, menciptakan _deadly embrace_.

Sources: Source 1, Source 3

Confidence: HIGH

### Finding 2
Claim: Penanganan deadlock otomatis melalui aborsi/rollback (Deadlock Victim).

Evidence: Karena siklus tidak dapat selesai dengan sendirinya, database engine memiliki _deadlock monitor_ yang secara periodik mendeteksi siklus. Engine secara paksa membatalkan (rollback) salah satu transaksi untuk membebaskan lock, yang kemudian memunculkan exception (mis. error 1205 di SQL Server atau error deadlock di Postgres).

Sources: Source 1, Source 3

Confidence: HIGH

### Finding 3
Claim: Urutan akses (Lock Ordering) mencegah mayoritas deadlock.

Evidence: Jika seluruh bagian aplikasi memodifikasi row/tabel dalam urutan alfabetis atau hierarkis yang sama secara konsisten, siklus tidak dapat terbentuk.

Sources: Source 1, Source 3

Confidence: HIGH

### Finding 4
Claim: Durasi transaksi yang panjang (termasuk interaksi eksternal) memperbesar probabilitas deadlock.

Evidence: Semakin lama _transaction block_ bertahan, semakin lama pula lock tertahan. Praktik yang disarankan adalah _keep transactions short_.

Sources: Source 3

Confidence: HIGH

### Finding 5
Claim: Deadlock dapat ditangani dengan mekanisme Retry di aplikasi.

Evidence: Dokumentasi resmi menegaskan bahwa deadlock sering kali bukan bug logic fatal yang bisa dicegah 100%, melainkan kejadian wajar dalam sistem dengan konkurensi tinggi. Aplikasi harus menangkap exception tersebut dan melakukan proses retry otomatis.

Sources: Source 1, Source 3

Confidence: HIGH

## Areas of Agreement
- Semua sumber utama (PostgreSQL dan Microsoft) setuju bahwa aplikasi tidak boleh membiarkan transaksi berjalan lama.
- Semua sumber utama setuju bahwa lock ordering yang konsisten adalah teknik terbaik untuk pencegahan.
- Semua setuju retry logic sangat diperlukan karena pencegahan absolut tidak praktis.

## Areas of Disagreement
No material disagreements discovered.

## Limitations
Dokumentasi dari MySQL tidak dapat diakses secara langsung selama proses pengumpulan data (terkendala status 403), sehingga penelitian lebih ditopang oleh ekosistem PostgreSQL dan Microsoft SQL Server. Namun secara konsep relasional, mekanismenya sama.

## Conclusion
Deadlock bukanlah kelemahan dari mesin database, melainkan konsekuensi logis dari data konkurensi. Pengendaliannya berpusat pada perbaikan arsitektur kode aplikasi: menggunakan lock ordering yang konsisten, meminimalkan durasi query dalam satu transaksi, dan menerapkan retry pattern secara anggun (graceful).
