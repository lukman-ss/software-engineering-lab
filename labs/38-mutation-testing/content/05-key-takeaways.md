# Key Takeaways

1. **Cakupan Kode Bukan Tolok Ukur Kualitas Pengujian**: *Code coverage* (line/statement/branch) hanya mengukur baris kode apa saja yang dilewati saat tes berjalan, bukan apakah pengujian tersebut sanggup mendeteksi kegagalan logika.
2. **Asersi Lemah Menciptakan Ilusi Keamanan Semu**: Pengujian tanpa asersi ketat dapat dengan mudah mencapai cakupan 100% namun menghasilkan *mutation score* 0%, membiarkan bug logika lolos ke tahap *production*.
3. **Penyuntikan Mutasi Sistematis Meniru Kesalahan Nyata**: *Mutation testing* secara otomatis menghasilkan mutan sintaksis berbasis Abstract Syntax Tree (AST)—meliputi operator relasional (`>` vs `>=`), logika boolean (`&&` vs `||`), operasi aritmetika (`*` vs `/`), dan pergeseran nilai batas (+1).
4. **Model RIP Menjelaskan Kegagalan Deteksi**: Agar mutan terbunuh, tes harus mencapai kode (*Reach*), mutasi harus mengubah *state* (*Infect*), dan perubahan tersebut harus merambat hingga ke asersi pengujian (*Propagate*). Pengujian lemah umumnya gagal pada tahap propagasi.
5. **Dua Hipotesis Dasar Pengujian**: *Mutation testing* bertumpu pada *Competent Programmer Hypothesis* (program mendekati benar, kesalahan biasanya kecil) dan *Coupling Effect* (kesalahan sederhana saling bertingkat membentuk kegagalan kompleks).
6. **Tantangan Equivalent Mutants**: Mutan yang secara semantik tidak mengubah keluaran program merupakan persoalan yang *mathematically undecidable*; perkakas praktis mengatasinya melalui filter heuristik atau analisis terarah.
7. **Strategi Adopsi Produksi**: Jangan memaksakan skor 100% di seluruh modul sekaligus. Gunakan *incremental/differential mutation testing* pada jalur CI/CD untuk fokus pada modul domain inti dan perubahan kode terkini.
