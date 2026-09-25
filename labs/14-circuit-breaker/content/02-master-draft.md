# Pola Circuit Breaker: Mencegah Kegagalan Beruntun pada Sistem Terdistribusi

## Masalah (Problem)

Dalam arsitektur sistem terdistribusi dan layanan mikro (*microservices*), setiap komponen bergantung pada pemanggilan jaringan jarak jauh (*remote calls*). Panggilan jaringan melekat dengan ketidakpastian. Layanan *downstream* (tujuan) bisa menjadi lambat secara mendadak atau mati (*hang*) sama sekali akibat lonjakan lalu lintas data, kegagalan *database*, atau pemadaman sebagian infrastruktur jaringan.

Tanpa adanya mekanisme perlindungan yang ketat, *caller threads* (utas dari layanan yang memanggil) akan terblokir karena harus menunggu proses *timeout*. Proses penahanan *thread* ini secara bersamaan menahan soket jaringan dan alokasi memori lokal. Seiring berjalannya waktu dan bertambahnya jumlah panggilan yang masuk, seluruh sumber daya sistem pada layanan pemanggil akan terkuras hingga habis (*exhausted*).

## Mengapa Ini Penting (Why This Matters)

Ketidakmampuan sistem untuk menangani kelambatan *downstream* akan berujung pada fenomena yang dikenal sebagai **Cascade Failure** (Kegagalan Beruntun). 
Dalam infrastruktur kritikal, dampak dari satu layanan kecil yang mengalami degradasi dapat menyebar bagaikan efek domino. Jika layanan A (pemanggil) lumpuh karena menunggu layanan B (*downstream*), maka layanan di atas A (seperti API Gateway) juga akan kehabisan koneksi. Pada titik puncaknya, *health-check endpoints* akan berhenti merespons, memicu pemadaman total (*total outage*) pada platform yang tadinya sehat hanya karena satu kegagalan yang tidak terisolasi. Membatasi *blast radius* (radius dampak) dari sebuah insiden adalah kewajiban mutlak dalam rekayasa keandalan sistem (*site reliability engineering*).

## Mental Model

Pola *Circuit Breaker* (Pemutus Sirkuit) terinspirasi dari sistem kelistrikan di dunia nyata. Pemutus sirkuit dipasang untuk memonitor beban kelistrikan dan seketika "terputus" (*trips open*) ketika beban melebih kapasitas, guna mencegah rumah terbakar. 

Dalam perangkat lunak, sebuah *Circuit Breaker* berwujud sebagai objek proksi (*state machine*) yang diletakkan di antara *caller* dan layanan *downstream*. Proksi ini memonitor jumlah kegagalan berturut-turut. Saat jumlah kegagalan melampaui ambang batas yang dikonfigurasi, proksi akan bergeser statusnya dan memutuskan koneksi sementara. Setiap panggilan baru yang masuk tidak akan diteruskan melalui jaringan, melainkan langsung ditolak atau digagalkan secepat mungkin (*fail-fast*) dalam hitungan nanodetik. Circuit Breaker memberikan ruang bernapas yang aman bagi *downstream* untuk pulih dan secara bertahap melakukan pengujian acak (*canary probe*) sebelum mengembalikan lalu lintas jaringan secara penuh.

## Skenario Kegagalan Beruntun (Failure Scenario)

Sebagai ilustrasi verifikasi dari eksperimen lab ini, perhatikan alur di bawah ini saat *Circuit Breaker* **tidak** diimplementasikan:
1. Layanan *Payment* mengalami *down* atau berhenti memproses permintaan.
2. *Worker threads* di Layanan *Checkout* tertahan (*blocked*) pada proses respons HTTP yang tidak kunjung datang.
3. Seluruh *thread pools* dan *connection pools* di Layanan *Checkout* terisi penuh oleh permintaan yang menggantung.
4. Karena Layanan *Checkout* macet total, antrean HTTP di *API Gateway* yang meneruskan lalu lintas eksternal ikut penuh.
5. Skenario fatal ini membuktikan bahwa latensi akibat permintaan macet secara komulatif berbanding lurus dengan persamaan `N * timeout`. 

## Konsep Inti dan Status Sirkuit (Core Concept)

Mesin status (*state machine*) pada implementasi *Circuit Breaker* lab ini dipecah ke dalam tiga kondisi (status) definitif:

### 1. Status CLOSED (Tertutup)
Ini merupakan kondisi operasional normal sistem. Sirkuit tertutup rapat sehingga semua panggilan jaringan (*requests*) diizinkan untuk dikirim secara langsung ke layanan *downstream*. 
- Jika panggilan sukses: Penghitung kegagalan (*failure count*) akan dikembalikan menjadi nol.
- Jika panggilan gagal (berupa error jaringan/server): Penghitung kegagalan ditambah satu (di-*increment*).
- Ketika jumlah kegagalan mencapai batas `FailureThreshold`, status sirkuit seketika berubah menjadi **OPEN**.

### 2. Status OPEN (Terbuka)
Status darurat di mana sirkuit diputus. Pada fase ini:
- Sama sekali **tidak ada** lalu lintas pemanggilan jaringan (*network calls*) yang dikirimkan.
- Semua permintaan masuk akan seketika digagalkan dengan membuang instansiasi objek `ErrCircuitOpen`. Operasi penolakan ini diselesaikan hanya dalam pecahan hitungan mikrodetik tanpa beban jaringan I/O.
- Layanan *downstream* benar-benar diistirahatkan agar proses pemulihannya (*recovery*) tidak dibanjiri oleh re-transmisi otomatis dari pemanggil.
- Timer jeda istirahat atau *cooldown* (`OpenTimeout`) mulai berjalan mundur.

### 3. Status HALF-OPEN (Setengah Terbuka)
Status uji coba paska-kegagalan yang terpicu secara otomatis begitu masa *cooldown* `OpenTimeout` habis berlalu.
- Circuit breaker akan melewatkan sejumlah panggilan spesifik dengan limit kecil (diatur oleh parameter `HalfOpenMaxCalls`) sebagai *canary probe* (pemeriksaan kondisi).
- Jika *probe* berhasil dan mengembalikan respons yang diharapkan, sistem meyakini bahwa *downstream* telah pulih. Status akan langsung dialihkan kembali ke **CLOSED** dan operasional kembali normal secara penuh.
- Jika *probe* kembali gagal, *Circuit Breaker* menyimpulkan gangguan masih persisten. Status segera ditarik lagi ke **OPEN**, dan siklus *cooldown* kembali diulang dari awal.

## Arsitektur dan Implementasi Internal

Di dalam repositori lab `14-circuit-breaker`, implementasi arsitektur divalidasi ke dalam tiga komponen sentral:
```text
Layanan Checkout (Service)
       │
       ▼
Proksi Circuit Breaker (State Machine)
  ├── [CLOSED]   ──► Memanggil Payment Client ──► Fake Payment Server HTTP
  ├── [OPEN]     ──► Penolakan Fail Fast (ErrCircuitOpen)
  └── [HALF-OPEN]──► Pengujian Ulang/Probe 
```

Struktur *Circuit Breaker* merangkum implementasi rekayasa konkurensi bawaan dari bahasa Go. Alih-alih merancang struktur bebas-kunci (*lock-free structures*) yang rumit dengan penyangga cincin (*ring buffer*) bertenaga atomik, implementasi lab ini disederhanakan menggunakan pelindung standar `sync.Mutex`. Pilihan ini memberikan keseimbangan terbaik antara kebenaran status (ketiadaan *data race*) dan kode yang mudah dirawat.

Perpindahan masa berlaku *cooldown* diselesaikan secara lekas (tanpa perlunya perulangan *goroutine* asinkron di belakang layar) melalui metodologi *lazy evaluation*. Pada fungsi perantara `checkStateTransitionLocked()`, sistem memeriksa selisih waktu secara matematis (`now - lastStateChange`) ketika objek metode tereksekusi. Pendekatan kalkulasi *on-demand* ini mencegah *memory leak* dari proses pekerja tak terbatas. 

## Timeout vs Retry vs Circuit Breaker

Ketiga mekanisme ini sering disalahartikan padahal fungsinya sangat berbeda di tingkat infrastruktur ketahanan (*resiliency*).
- **Timeout**: Sekadar batas waktu tenggang maksimal bagi satu permintaan individu untuk menyerah guna mencegah *hang* abadi. Tidak melindungi sistem dari banjir sirkulasi kegagalan.
- **Retry**: Upaya agresif memicu ulang transmisi secara otomatis yang mengasumsikan kegagalan bersifat cacat minor (*transient*). Saat sebuah server benar-benar kelebihan kapasitas (*overload*), tindakan *retry* membabi buta dari ratusan *client* akan memicu kondisi ekstrem yang disebut *retry storm*.
- **Circuit Breaker**: Mekanisme defensif menyeluruh. Menghentikan total laju request pada ambang keandalan tertentu. Memberikan perlindungan ganda (mencegah *caller* bunuh diri membuang koneksi memori, sambil memulihkan *downstream*). Pada praktiknya, Circuit Breaker bekerja sama dengan Timeout, dan sering meredam (silence) *retry storm*.

## Pertimbangan Produksi (Production Considerations)

Hasil reviu arsitektural (audit) menyepakati berbagai catatan penting sebelum menerapkan kode simulasi ini pada skala beban berat di tahap produksi:

1. **Konfigurasi Ambang Batas Waktu**: 
Waktu putus (timeout HTTP 100ms) dan jeda pemulihan (cooldown 300ms) pada kode simulasi ini hanyalah *ilustrasi semata* demi pengujian eksekusi unit yang responsif dan berjalan cepat (kurang dari sedetik). Di skenario dunia nyata, nilai *timeout* harus selalu ditelusuri dari kesepakatan *Service Level Agreement* (SLA) dan diukur berdasarkan metrik penyebaran latensi P99 dari layanan tersebut (contoh: *cooldown* 30 sampai 60 detik atau *backoff* terukur).
2. **Algoritma Penghitungan Error**: 
Lab membuktikan konsep proteksi sirkuit dengan mengandalkan sistem "penghitungan kegagalan berturut-turut" yang sederhana (*consecutive failures count*). Aplikasi skala *enterprise* biasanya mengadopsi penghitungan dengan berbasis persentase rasio kegagalan atau kerangka pemetaan waktu bergerak (*sliding time-window / leaky bucket*) agar lebih tahan terhadap gangguan sesaat.
3. **Isolasi Node Mutex**:
Variabel status mesin (`state`) diamankan dengan Mutex sehingga berjalan stabil khusus pada memori tunggal (*single instance*). Untuk menjaga presisi pemutusan sirkuit pada topologi dengan puluhan node mikro, pertimbangkan penggunaan penyimpanan sinkronisasi terdistribusi seperti Redis, terlepas sistem memori lokal sering kali sudah cukup memberikan perlindungan batas alokasi *thread*.

## Bukti Pengujian Terverifikasi (What the Tests Prove)

Kode lab (`circuit_breaker_test.go` dan `integration_test.go`) melewati fase pengujian teknik komprehensif, mengonfirmasi behavior valid:
1. Validasi Transisi: Skema `CLOSED -> OPEN -> HALF_OPEN -> CLOSED` berhasil melewati seluruh kriteria pengujian unit (*asserting* status dengan tepat). 
2. Proteksi Pemanggilan Turunan: Selama pengujian integrasi di status `OPEN`, terbukti *secara empiris* bahwa metode downstream *tidak* dieksekusi secara fisikal ke jaringan (`downstream calls = 0`).
3. Perlindungan *Fail-Fast*: Waktu respons HTTP untuk status kegagalan `OPEN` pada eksekusi pelaporan demo (`cmd/demo`) turun sangat drastis, bergeser dari tingkat waktu tunggu ratusan milidetik (karena latensi layanan lambat) menjadi *fail-fast* rata-rata dalam ukuran hitungan < 1µs (mikrodetik). 
4. Aman secara Konkurensi: Seluruh eksekusi multi-goroutine lulus *Go Race Detector* (`go test -race`). Tidak ada pembacaan status tumpang tindih meskipun diserbu secara paralel oleh ratusan *worker*.

## Fallback dan Pencegahan (Recovery / Rollback)

Salah satu langkah mitigasi tingkat lanjut paska terjadinya gagal sirkuit adalah melakukan *Fallback* (Mekanisme Penanganan Darurat). 
*Fallback* yang valid termasuk mengekstrak hasil penyimpanan sementara (data *cache*) pada sistem pembacaan (*read cache*), mengembalikan nilai pengaturan standar (pengganti statis), atau mengarahkan pesan tertunda (pesan asinkronus ke struktur *queue*).

**Peringatan Penting (Verified Audit Constraint)**: *Jangan pernah menggunakan Fallback tersembunyi (silent fallback) yang memaksa kesuksesan semu pada permintaan mutasi kritikal berstatus esensial.* Contoh mutasi esensial mencakup deduksi nilai e-money, pembukuan ledger debit akun saldo bank, atau pengajuan pesanan inventaris fisik akhir. Pada skenario tersebut, laporkan secara tegas respons error terbuka dari *Circuit Breaker* kembali ke sistem hulu atau antar-muka pengguna.

## Studi Kasus Infrastruktur Asli (Case Study)

### Decoupling Non-Kritikal (Kasus CMMS)
Skenario CMMS (sistem faktur order):
Alur: `Buat Invoice` -> `Render PDF` -> `Kirim Pesan WhatsApp`
Gangguan asinkron tidak boleh membahayakan basis transaksi utama. Jika aplikasi API *WhatsApp* mati, tidak masuk akal membatalkan pembuatan transaksi utama *Invoice*. Solusi disetujui: Modul pencatat pesanan dieksekusi secara sikron ke Basis Data Master. Selanjutnya, pekerja di balik layar (*background worker*) meluncurkan panggilan *WhatsApp* di dalam perlindungan *Circuit Breaker*. Apabila pemutus terpicu *OPEN*, retensi pengiriman diamankan secara otomatis dengan menahannya sementara (*leveling load*) dalam *message queue* seperti RabbitMQ.

### Skenario Batasan (Kasus PPOB)
Alur transaksi produk digital (Pulsa/Tagihan):
Alur Utama: `Pembeli` -> `Server Operator Utama` -> `Integrasi Payment Gateway` -> `Provider (Pihak 3) Pulsa`
Layanan sinkron kritis seperti *Payment Gateway* harus berjalan tepat waktu. Jika bergantung dengan server transaksi pembayaran secara mutlak dan server mati total, proksi berstatus *OPEN* harus memaksa penggagalan sekuens secepat mungkin (*fail-fast*) tanpa penanganan asinkronis atau pengurangan saldo. Sebaliknya, proses *Provider Pulsa* hulu bersifat bisa dicoba ulang secara asinkron (*retryable async dependency*) yang berpadanan kuat dengan *Idempotency Key* yang terkunci aman. 

## Kesalahan Umum Praktisi (Common Mistakes)

1. **Threshold Terlalu Rendah**: Merancang pemutus dengan jumlah 1-2 metrik gagal sehingga berujung pemutusan prematur karena blip sementara dari *network jitter*. Traffic valid sering direkayasa tertolak. 
2. **Kelebihan Probe Banjir (Probe Flood)**: Mengonfigurasi `HalfOpenMaxCalls` terlalu masif. Ini menyebabkan ketika layanan hilir yang baru sembuh masih tertatih perlahan (*slow ramp up*), proksi justru menembak ulang ratusan *canary* probe secara bersamaan sehingga layanan hilir terkapar hancur kembali. 
3. **Mengabaikan Pemisahan Tipe Error**: Menyamaratakan validasi masukan error. Menghitung kode HTTP seri *4xx* (*Client Errors/Bad Request*) secara bodoh sebagai komponen kerusakan server internal, ketimbang membatasi deteksi sirkuit pada kesalahan murni koneksi, *timeouts*, atau status galat server infrastruktur *5xx*.

## Sumber Utama (Sources)
- **Martin Fowler**: *Circuit Breaker* (martinfowler.com)
- **Microsoft Learn**: *Circuit Breaker Pattern - Azure Architecture* 
- **AWS Builders Library**: *Timeouts, retries, and backoff with jitter*
