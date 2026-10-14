Lab 36 kali ini membahas miskonsepsi fatal dalam keamanan web modern bahwa konfigurasi CORS dapat melindungi backend dari mutasi data berbahaya. Pada kenyataannya, CORS hanya mengatur izin pembacaan respons oleh JavaScript di browser, sedangkan serangan CSRF mengeksploitasi pengiriman request write lintas origin yang diizinkan secara historis oleh Same-Origin Policy.

Banyak engineer mengira jika origin penyerang tidak terdaftar di whitelist CORS, server otomatis menolak request tersebut. Padahal, browser tetap mengirimkan simple POST langsung ke server beserta cookie otentikasi korban, sehingga mutasi database selesai sebelum browser memeriksa respons.

Eksperimen Go standard library pada lab ini membuktikan:
• CORS tidak mencegah eksekusi backend: simple POST dari evil-attacker.com memotong saldo rekening korban dari $1000 menjadi $700 meskipun header Access-Control-Allow-Origin tidak diberikan.
• Signed Double-Submit Token HMAC-SHA256 yang terikat ke Session ID dan timestamp TTL berhasil memblokir serangan CSRF tanpa token dengan HTTP 403 Forbidden.
• Fetch Metadata (Sec-Fetch-Site) menolak mutasi cross-site di layer gateway secara mandiri.
• Custom Header X-Requested-With menetralisasi eksploitasi form HTML klasik secara total.
• Konfigurasi cookie SameSite=Lax menahan pengiriman cookie pada request mutasi lintas situs.

Mental model: SOP restricts cross-origin read, CSRF exploits cross-origin write. CORS relaxes read, it never defends write.

Batasan implementasi: Lab ini menggunakan HMAC secret in-memory dan string delimiter untuk parsing token; sistem produksi memerlukan secret vault dan serialisasi data terstruktur.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/36-cors-and-csrf

#softwareengineering #golang #cybersecurity #webdevelopment #systemdesign #backend