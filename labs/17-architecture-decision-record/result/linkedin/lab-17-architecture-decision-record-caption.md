Lab 17 kali ini membahas pengelolaan keputusan arsitektur menggunakan Architecture Decision Record (ADR) serta validasi struktural otomatis di Go.

Masalah teknis yang sering terjadi adalah hilangnya konteks historis di balik keputusan desain penting saat anggota tim berganti. Dokumentasi yang terpisah di wiki luar cepat usang, memicu perdebatan teknis berulang atau pembatalan arsitektur secara sembarangan tanpa memahami batasan aslinya.

Lab ini membuktikan tiga mekanisme konkret:
1. Regex-Based Parsing: Mengekstrak ID, Title, Status, serta klausa silsilah SupersededBy dan Supersedes dari berkas Markdown, sekaligus memastikan bagian Context, Decision, dan Consequences tidak kosong.
2. Monotonic & Lineage Integrity: Linter mendeteksi penomoran non-monotonik, ID duplikat, tautan sepihak, target fiktif, self-supersession, serta siklus penggantian tertutup melalui DFS cycle detection.
3. Concurrent Verification: Validasi graf silsilah timbal-balik dieksekusi paralel dengan goroutine dan lolos 23 test assertions di bawah race detector tanpa peringatan data race.

Mental model inti: Decisions as versioned invariants. Simpan keputusan arsitektural langsung bersama kode sumber sebagai catatan permanen yang tidak dapat disunting retroaktif, melainkan digantikan melalui silsilah penggantian timbal-balik yang divalidasi mesin. Usulan yang ditolak dicatat sebagai status Rejected agar argumen penolakan terdokumentasi.

Satu batasan penting: Lab ini menggunakan parser ekspresi reguler dengan format heading Markdown yang ketat dan bekerja pada struktur in-memory, belum menyediakan AST CommonMark umum atau pengait otomatis git pre-commit hook.

Pelajari implementasi lengkap dan rangkaian pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record

#softwareengineering #golang #systemdesign #backend #architecture #cleancode
