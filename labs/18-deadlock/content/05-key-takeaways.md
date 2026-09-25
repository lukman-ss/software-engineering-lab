# Key Takeaways

1. **Deadlock Adalah Konsekuensi Konkurensi**: Terjadinya deadlock dalam transaksi basis data bukanlah bukti kerusakan engine, melainkan akibat benturan urutan akses data yang terjadi secara simultan (*circular wait*).
2. **Aborsi sebagai Resolusi Otomatis**: Basis data tidak membiarkan deadlock bertahan selamanya. Sistem akan memilih satu transaksi sebagai korban (*deadlock victim*), membatalkan aksinya (*rollback*), dan melepas kuncinya agar sistem dapat berjalan kembali.
3. **Pencegahan Utama via Lock Ordering**: Siklus saling menunggu dapat dicegah secara total dengan mengurutkan ID entitas yang dimutasi. Jika semua goroutine/koneksi mengakses baris dengan urutan yang identik (misalnya A sebelum B), maka deadlock tidak akan pernah terjadi.
4. **Pendekkan Durasi Transaksi**: Semakin lama sebuah transaksi aktif (menahan lock), semakin besar jendela probabilitas terjadinya tabrakan lock. Selalu usahakan blok transaksi hanya berisi proses tulis ke database.
5. **Dilarang I/O Eksternal Dalam Transaksi**: Jangan pernah memanggil layanan API pihak ketiga, upload file, atau menunggu respons eksternal di saat transaksi basis data sedang memegang lock.
6. **Retry Adalah Standar Wajib Aplikasi**: Karena mencegah deadlock 100% pada sistem terdistribusi sulit dan membatasi throughput, developer harus menangani kode error deadlock (seperti error 1205/40P01) dengan mengimplementasikan iterasi perulangan otomatis (*backoff & retry*).
7. **Deteksi Membutuhkan Waktu**: Basis data tidak seketika membatalkan transaksi yang menunggu. Ada variabel penundaan (seperti `deadlock_timeout`) yang menahan resolusi agar engine tidak menghabiskan CPU untuk mencari siklus WFG (*Wait-For Graph*) secara berlebihan, sehingga latensi aplikasi terdampak sebelum korban dibatalkan.
