# Contradictions & Uncertainties

## 1. DI membuat test lebih mudah — setuju atau tidak?

- **Fowler (2004)**: "both [DI and Service Locator] are very amenable to stubbing. ... I suspect this observation [that DI makes testing easier] comes from projects where people don't make the effort to ensure that their service locator can be easily substituted."
- **Spring docs (2024-2026)**: "your classes become easier to test, particularly when the dependencies are on interfaces..." (menyatakan DI secara umum memudahkan testing).
- **Microsoft .NET docs (2026)**: menyatakan hard-coded dependencies membuat unit test sulit; DI memungkinkan mock/stub.

### Assessment
Kesepakatan luar biasa: DI **memungkinkan** testing dengan mock. Perbedaan adalah Fowler menekankan bahwa **Service Locator yang dirancang baik juga bisa di-mock**, jadi keuntungan testing bukan eksklusif milik DI — melainkan pola yang dipilih. Microsoft/Spring menyederhanakan dengan beralih ke DI yang direkomendasikan. Ini adalah perbedaan penekanan, bukan fakta bertolak belok.

## 2. Interface Injection

- **Fowler (2004)**: Mendeskripsikan Interface Injection sebagai salah satu dari tiga form DI utama. Menyebutkan Avalon sebagai contoh framework yang menggunakannya.
- **Spring Framework 7.0.9 docs**: Hanya mendokumentasikan Constructor-based DI dan Setter-based DI. Tidak menyebut Interface Injection sama sekali.
- **Microsoft .NET docs (2026)**: Fokus pada constructor injection.
- **Laravel 12.x docs (2024)**: Menggunakan constructor/setter + attributes (PHP 8). Interface Injection tidak disebut.

### Assessment
Interface Injection dihapuskan dari praktik modern. Fowler sendiri pada 2004 memprediksi: "Interface Injection is more invasive since you have to write a lot of interfaces..." dan framework lightweight modern tidak memilih pendekatan ini. Ini adalah evolusi, bukan kontradiksi — konsensus modern telah memilih constructor+setter injection.

## 3. "12 parameter" sebagai ambang batas constructor over-injection

- **Spesifikasi topik lab**: Menyebutkan "Kalau constructor berisi 12 parameter ... biasanya ada masalah desain."
- **Fowler (2004)**: Menyebutkan "you have a lot of constructor parameters things can look messy" dan "it's often a sign of an over-busy object that should be split" — tetapi **tidak memberikan angka spesifik**.

### Assessment
Angka "12" berasal dari spesifikasi topik, **bukan sumber primer**. Fowler memberi sinyal kualitatif (banyak, berantakan), tetapi tidak menetapkan ambang numerik. Klaim spesifikasi tidak terverifikasi secara independen. Ini harus ditandai sebagai **interpretasi/Heuristik** bukan fakta yang diverifikasi.

## 4. NestJS DI Documentation Retrieval

- **NestJS docs (Source 7)**: Saat di-fetch, halaman hanya mengembalikan redirect ke "Documentation | NestJS". Konten kode spesifik `@Injectable()`, `@Module` tidak tersedia secara penuh dari fetch.
- **Claim pada Evidence 14**: Menggambarkan pola NestJS berdasarkan pengetahuan umum, tetapi **bukan dari konten yang berhasil dibuka**.

### Assessment
NestJS dipasangkan sebagai contoh ekivalensi pola (Decorator-based DI), tetapi **bukti spesifik tidak berhasil diperoleh dari webfetch**. Confidence diturunkan ke MEDIUM. Perlu verifikasi ulang jika dibutuhkan.

## 5. Kapan harus / tidak pakai DI — nilai tambah yang tak terverifikasi

- **Spesifikasi topik**: Memberikan aturan praktis ("Gunakan DI terutama untuk: Database, HTTP Client, ...") dan nilai tambah ("Object sederhana tidak perlu di-inject").
- **Fowler (2004)**: Hanya prinsip — "separating configuration from use" — dan "leave it up to the user".

### Assessment
Spesifikasi topik berisi **pedoman heuristik praktis**, tetapi **tidak ada sumber primer yang secara eksplisit mengklasifikasikan daftar objek mana yang "harus" di-inject vs. dibuat langsung**. Fowler memberikan prinsip, bukan checklist. Klaim ini kategori perkiraan/desain, bukan fakta teknis yang terverifikasi.

---

## Kesimpulan
Tidak ada **material contradictions** — semua sumber utama saling mendukung pada inti definisi DI/IoC, tiga bentuk DI, DI vs Service Locator, dan peran Service Locator sebagai anti-pattern (PSR-11 eksplisit). Perbedaan utama adalah **luasnya penekanan** dan **evolusi praktik** (Interface Injection ditinggalkan, angka spesifik 12 param berasal dari spesifikasi topik bukan sumber primer, NestJS evidence parsial).
