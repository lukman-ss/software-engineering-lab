Lab 18 kali ini membahas mekanisme terbentuknya deadlock pada sistem konkuren, penanganan oleh database engine, serta strategi pencegahan di tingkat aplikasi.

Masalah teknis yang sering terjadi adalah dua transaksi atau goroutine mengakses dan mengunci sumber daya bersama secara bersamaan dengan urutan berlawanan, menciptakan siklus circular wait yang membekukan pemrosesan transaksi.

Lab ini membuktikan empat hal konkret:
• Deadlock Occurrence: Transfer naif dua arah yang dieksekusi simultan terbukti memicu kegagalan deadlock victim saat context timeout terlampaui.
• Lock Ordering: Pengurutan akuisisi lock berdasarkan identitas akun memutus terbentuknya siklus secara deterministik dan menjaga konsistensi saldo.
• Application Retry: Mekanisme retry pada tingkat aplikasi berhasil memulihkan transaksi yang sempat dibatalkan oleh database engine.
• Transaction Duration: Pengujian durasi transaksi membuktikan bahwa menahan lock lebih lama secara nyata melipatgandakan peluang terjadinya deadlock.

Mental model inti: Deadlock is a circular dependency. Database engine membatalkan salah satu transaksi sebagai victim. Lock ordering mencegah siklus terbentuk, sementara application retry memulihkan transaksi korban.

Satu batasan penting: Lab ini adalah simulasi tingkat aplikasi berbasis channel Go dan context timeout, bukan implementasi Wait-For Graph (WFG) internal database engine, serta menggunakan jeda retry tetap 2ms untuk kebutuhan pengujian.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock

#softwareengineering #golang #concurrency #database #systemdesign #backend
