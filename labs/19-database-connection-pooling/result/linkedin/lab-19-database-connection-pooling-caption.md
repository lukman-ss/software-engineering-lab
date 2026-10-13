Lab 19 kali ini membahas analisis mendalam database connection pooling, perumusan kapasitas berbasis hardware, serta mitigasi kegagalan fatal akibat kebocoran koneksi.

Masalah teknis yang sering terjadi adalah lonjakan error FATAL: sorry, too many clients already yang keliru direspons dengan memperbesar connection pool atau menaikkan max_connections server, sehingga latensi memburuk drastis dan server mengalami kehabisan slot.

Lab ini membuktikan empat hal konkret:
• Overhead Koneksi: Mode unpooled membutuhkan waktu berkali-kali lipat lebih lama akibat TCP handshake berulang, sedangkan pooled memakai ulang koneksi aktif seketika.
• Server Exhaustion: Mengizinkan pool aplikasi melebihi kapasitas max_connections server secara deterministik memicu penolakan koneksi ErrServerOverloaded.
• Connection Starvation: Menahan koneksi database saat menjalankan I/O eksternal lambat menyebabkan starvation dan timeout context deadline exceeded pada request lain yang sehat.
• Anti-Deadlock Sizing: Mengakuisisi koneksi kedua secara simultan pada pool berukuran 1 membuktikan risiko pool-locking deadlock yang memerlukan rumus kapasitas minimum.

Mental model inti: Connection capacity is dictated by database hardware constraints, not by client concurrency. Keep pools lean, reuse connections, and never hold sockets during external I/O.

Satu batasan penting: Lab ini menggunakan simulasi in-memory berbasis driver Go database/sql untuk menguji batas konkurensi dan latensi handshake, bukan proses database PostgreSQL nyata, serta belum memodelkan OOM process backend atau PgBouncer session mode.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/19-database-connection-pooling

#softwareengineering #golang #database #backend #systemdesign #performance
