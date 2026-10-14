Lab 29 kali ini membahas implementasi Saga Pattern untuk mengelola transaksi terdistribusi tanpa Two-Phase Commit pada arsitektur microservices berbasis Go.

Tantangan utama pada pola Database-per-Service adalah hilangnya jaminan transaksi ACID lintas layanan. Memaksakan protokol Two-Phase Commit memicu penguncian baris data secara sinkron lintas jaringan, meningkatkan latensi drastis, serta memunculkan titik kegagalan tunggal yang melumpuhkan ketersediaan sistem. Ketika proses pemesanan gagal di tengah jalan, database lokal yang telah melakukan commit tidak dapat dibatalkan melalui mekanisme rollback storage engine standar.

Lab ini membuktikan lima prinsip arsitektur terdistribusi:
• Orchestration Engine: Koordinator sentral mengeksekusi rangkaian transaksi lokal secara sekuensial dan memantau status setiap tahapan bisnis.
• LIFO Rollback: Kegagalan reservasi stok memicu kompensasi mundur dari pemotongan dana kartu kredit hingga pembatalan pesanan.
• Event-Driven Choreography: Komunikasi asinkron via EventBus memungkinkan layanan bereaksi mandiri terhadap domain event.
• Idempotency Protection: Pengecekan payment ID mencegah duplikasi pendebetan saldo saat terjadi retry request akibat network glitch.
• Semantic Locking: Status pesanan berstatus pending mengunci modifikasi konkuren sehingga mencegah terjadinya anomali data dirty reads.

Mental model inti: A saga is a sequence of local transactions where each step commits immediately. If failure occurs, compensating transactions execute in reverse order.

Satu batasan penting: Lab ini menguji logika koordinasi menggunakan in-memory state dan pemanggilan fungsi internal Go, bukan database relasional fisik atau message broker eksternal seperti Apache Kafka.

Pelajari implementasi dan seluruh test suite di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern

#softwareengineering #golang #microservices #distributedsystems #systemdesign #backend
