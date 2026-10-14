Lab 25 kali ini membahas strategi pertahanan beban bertingkat di Go: algoritma Token Bucket dan Leaky Bucket, bounded queue backpressure, penegakan HTTP 429 berbasis RFC 6585, dan exponential backoff dengan AWS Full Jitter.

Masalah teknis yang sering terjadi adalah lonjakan request yang melebihi kapasitas pemrosesan downstream worker. Tanpa mekanisme kendali beban, antrean memori tumbuh tanpa batas (unbounded queue), memicu degradasi latensi berbasis Little's Law, kegagalan memory exhaustion (OOM), dan badai retry tersinkronisasi (thundering herd) yang melipatgandakan beban saat sistem mencoba pulih.

Lab ini membuktikan empat pilar ketahanan sistem terdistribusi:
• Kendali Batas Token Bucket: Mengakomodasi lonjakan burst sesaat hingga kapasitas batas sambil mempertahankan laju rata-rata melalui penambahan token proporsional waktu.
• Isolasi Tenant RFC 6598: Middleware mengekstrak identitas X-API-Key untuk mencegah tabrakan kuota pada arsitektur Carrier-Grade NAT.
• Bounded Queue Backpressure: Buffer channel berukuran tetap menolak request baru secara non-blocking dengan ErrQueueFull saat antrean penuh.
• AWS Full Jitter Retry: Algoritma jitter mendistribusikan interval pengulangan klien secara seragam pada rentang eksponensial, melenyapkan sinkronisasi thundering herd.

Mental model inti: Rate limiting protects system boundaries from excessive ingress; bounded queue backpressure protects internal execution workers; full jitter retry prevents clients from synchronizing into destructive thundering herds.

Satu batasan penting: Lab ini menggunakan penyimpanan in-memory Go yang terikat pada satu host, sehingga deployment multi-node di lingkungan produksi memerlukan koordinasi terdistribusi seperti Redis dan message broker partisi Kafka.

Pelajari implementasi lengkap dan verifikasi pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/25-rate-limiting-and-backpressure

#softwareengineering #golang #systemdesign #distributedsystems #architecture #backend
