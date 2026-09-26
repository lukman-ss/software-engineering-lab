# Key Takeaways

1. Deadlock adalah kondisi sistemik di mana dua transaksi atau lebih saling menunggu lock yang ditahan pihak lain (*circular wait*).
2. Database engine menyelesaikan deadlock secara paksa dengan membatalkan (*abort/rollback*) salah satu transaksi sebagai *deadlock victim*.
3. Pencegahan paling efektif di tingkat aplikasi adalah menerapkan **Lock Ordering** yang konsisten di semua jalur eksekusi.
4. Pendekkan durasi transaksi secara signifikan menurunkan probabilitas terjadinya deadlock.
5. Mekanisme **Application Retry** adalah standar operasional wajib pada sistem konkuren tinggi untuk memulihkan transaksi yang menjadi korban deadlock secara transparan.
6. Lab ini adalah simulasi level aplikasi menggunakan channel Go + context timeout; tidak mengimplementasikan deteksi siklus nyata (Wait-For Graph) sebagaimana terdapat pada DBMS produksi.