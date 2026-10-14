Lab 31 kali ini membahas pemisahan mendasar antara otorisasi OAuth 2.0 dan autentikasi OpenID Connect (OIDC) serta mekanisme pertahanan token modern di Go.

Kesalahan umum arsitektur modern adalah memperlakukan access token sebagai kartu identitas pengguna. Access token hanya memberikan izin akses delegasi terhadap sumber daya, bukan membuktikan siapa pemilik sesi. Mengabaikan batas ini membuka celah pembajakan akun, intersepsi kode otorisasi pada browser, dan pencurian sesi melalui refresh token statis.

Lab ini menguji empat mekanisme pertahanan:
• Pertahanan Intersepsi Kode: Proof Key for Code Exchange (PKCE) S256 memastikan hanya client pemegang verifier asli yang dapat menukar kode otorisasi.
• Verifikasi Identitas Kriptografis: ID Token berformat JWT divalidasi tanda tangannya dengan HMAC-SHA256 serta diperiksa kecocokan issuer, audience, kedaluwarsa, dan nonce anti-replay.
• Rotasi Refresh Token Sekali Pakai: Setiap pembaruan sesi menerbitkan pasangan token baru dan menghanguskan token lama.
• Deteksi Replay Satu Keluarga: Penggunaan ulang token yang sudah dicabut memicu pencabutan seluruh token aktif dalam keluarga yang sama.
• Concurrency Safety: Seluruh mutasi token dilindungi sync.RWMutex dan lolos pengujian race detector.

Mental model inti: OAuth 2.0 grants delegated access authority while OIDC asserts cryptographically verified user identity. PKCE binds authorization codes to the initiator, and refresh token rotation with family revocation neutralizes token replay.

Batasan implementasi: Pengujian berjalan in-memory tanpa persistensi database dan menggunakan symmetric signing HMAC-SHA256 untuk kesederhanaan demonstrasi.

Pelajari kode sumber lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc

#softwareengineering #golang #security #oauth2 #oidc #systemdesign
