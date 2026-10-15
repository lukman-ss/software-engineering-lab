#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 36: CORS and CSRF Defense.
Defines metadata, colors, and 24 slides following Lab 14 & 15 standards.
"""
import os
import sys

OUTPUT_DIR = "labs/36-cors-and-csrf/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/36-cors-and-csrf"
COMMIT_SHA = "7ced81edaaa9ca043063d77e19a19b9438264382"
LAB_NO = "36"
SLUG = "cors-and-csrf"
TITLE_ID = "CORS & CSRF: Mengamankan Backend"
SUBTITLE_ID = "Membedah Realitas CORS dan Pertahanan CSRF Multi-Layer di Sisi Server"
TITLE_EN = "CORS & CSRF: Securing Backend from Cross-Origin Attacks"
SUBTITLE_EN = "Demystifying CORS and Implementing Server-Side Multi-Layer Anti-CSRF Defense"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "CORS • CSRF • Web Security • HMAC-SHA256 • Fetch Metadata • Go 1.22"

DEVTO_DESC = "Demystify CORS vs CSRF in Go: prove why CORS is not a backend firewall, and implement multi-layer anti-CSRF defense using signed tokens and Sec-Fetch-Site."
DEVTO_TAGS = "go, security, webdev, architecture"
DEVTO_INTRO = "Assuming CORS protects your backend from malicious write mutations is one of the most widespread security misconceptions. Lab 36 builds a spec-compliant CORS middleware and an HMAC-SHA256 signed CSRF protection engine in pure Go to expose why CORS fails against CSRF and how defense-in-depth mitigates cross-origin exploits."

MEDIUM_DESC = "Panduan teknis mendalam arsitektur keamanan web: mengapa CORS bukan pelindung mutasi backend, dan implementasi pertahanan anti-CSRF multi-layer di Go."

# Color constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", title_en="", narrative_en=""):
    """Add a slide to the presentation definition."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "title_en": title_en if title_en else title,
        "narrative_en": narrative_en if narrative_en else narrative,
    })

# Slide 1: Cover (Dark Navy)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor dengan scanlines dan kurva phosphor",
            "Simulasi SOP, CORS Preflight, dan Eksploitasi Form CSRF",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="SOP restricts cross-origin read; CSRF exploits cross-origin write. CORS relaxes read, it never defends write.",
    source_path="labs/36-cors-and-csrf/README.md",
    slide_num="01 / 24",
    narrative="",
    title_en=TITLE_EN,
    narrative_en="CORS does not stop cross-origin write requests; true protection requires server-side token validation."
)

# Slide 2: Miskonsepsi Fatal Keamanan Web
add_slide(
    categories=["MOTIVASI", "LATAR BELAKANG"],
    title="Miskonsepsi Fatal: CORS Bukan Firewall",
    subtitle="Banyak tim backend mengira CORS melindungi endpoint dari mutasi berbahaya",
    content_items=[
        {"title": "ASUMSI KELIRU", "items": [
            "Menganggap CORS adalah tembok api backend penolak request asing",
            "Meyakini penolakan header CORS otomatis menggagalkan serangan CSRF",
            "Merasa aman hanya karena origin penyerang tidak ada di whitelist CORS",
        ], "type": "bullet"},
        {"title": "FAKTA LAPANGAN", "items": [
            "Browser tetap mengeksekusi simple request POST langsung ke backend",
            "Database memproses mutasi saldo sebelum browser mengevaluasi header respons",
            "Penolakan CORS hanya memblokir pembacaan respons oleh JavaScript attacker",
        ], "type": "bullet"},
        {"title": "DAMPAK KERENTANAN", "items": [
            "Transaksi finansial ilegal dieksekusi tanpa otorisasi sadar dari korban",
            "Perubahan email dan password akun terjadi tanpa hambatan jaringan",
        ], "type": "bullet"},
    ],
    mental_model="CORS governs read permissions in JavaScript, not execution permissions on backend servers.",
    source_path="labs/36-cors-and-csrf/content/01-content-brief.md",
    slide_num="02 / 24",
    narrative="Miskonsepsi bahwa CORS bertindak sebagai firewall backend masih sangat umum ditemukan. Padahal, CORS dirancang oleh browser murni untuk mengontrol izin pembacaan respons oleh JavaScript, bukan membatasi server dalam memproses request. Tanpa pertahanan CSRF di backend, penyerang dapat memicu mutasi state di database hanya melalui tag form HTML biasa.",
    title_en="Fatal Misconception: CORS Is Not a Firewall",
    narrative_en="A common fallacy assumes CORS acts as a backend firewall. In reality, CORS only governs whether client JavaScript may inspect response bodies. Backend state mutation completes before browser CORS checks occur."
)

# Slide 3: Mental Model Inti: SOP, CORS, dan CSRF
add_slide(
    categories=["MENTAL MODEL", "BATAS KEAMANAN"],
    title="Mental Model: SOP vs CORS vs CSRF",
    subtitle="Memisahkan pembatasan read, relaksasi browser, dan eksploitasi write",
    content_items=[
        {"title": "SAME-ORIGIN POLICY (SOP)", "items": [
            "Membatasi JavaScript membaca data respons dari origin berbeda",
            "Secara historis mengizinkan pengiriman request write (misal form submit)",
        ], "type": "bullet"},
        {"title": "CORS (RELAXATION MECHANISM)", "items": [
            "Mekanisme negosiasi header HTTP untuk merelaksasi pembatasan SOP",
            "Mengizinkan frontend membaca data respons dari API domain lain",
            "Bukan mekanisme keamanan sisi server!",
        ], "type": "bullet"},
        {"title": "CSRF (WRITE EXPLOITATION)", "items": [
            "Mengeksploitasi perilaku browser yang otomatis melampirkan cookie sesi",
            "Penyerang memicu mutasi write tanpa perlu membaca data respons",
        ], "type": "bullet"},
    ],
    mental_model="SOP restricts reading; CORS relaxes reading; CSRF exploits writing. Defend the write at the server.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="03 / 24",
    narrative="Arsitektur keamanan web berdiri di atas tiga pilar ini. SOP melarang pembacaan respons lintas origin. CORS hadir sebagai mekanisme izin agar JavaScript dapat membaca data respons tersebut. Sementara itu, CSRF mengeksploitasi izin pengiriman data form yang secara historis diizinkan oleh SOP dengan memanfaatkan cookie sesi yang terlampir otomatis.",
    title_en="Core Mental Model: SOP vs CORS vs CSRF",
    narrative_en="Web security rests on three pillars: SOP blocks cross-origin reads; CORS selectively relaxes read restrictions; CSRF exploits unrestricted cross-origin writes via ambient credentials."
)

# Slide 4: Arsitektur Perbankan Lab 36
add_slide(
    categories=["ARSITEKTUR", "TOPOLOGI LAB"],
    title="Topologi Arsitektur Simulasi Bank",
    subtitle="Membandingkan endpoint rentan vs endpoint terproteksi dalam ekosistem Go",
    content_items=[
        {"title": "ENDPOINT DISPATCH", "items": [
            "GET  /api/balance             -> Cek saldo akun via cookie sesi",
            "GET  /api/csrf-token          -> Endpoint penerbit token HMAC terikat sesi",
            "POST /api/transfer/vulnerable -> Transfer rentan (hanya dibungkus CORS)",
            "POST /api/transfer/protected  -> Transfer aman (CORS + CSRF Token)",
            "POST /api/transfer/fetch-metadata -> Transfer aman via Sec-Fetch-Site",
            "POST /api/transfer/custom-header  -> Transfer aman via X-Requested-With",
        ], "type": "bullet"},
        {"title": "ZERO DEPENDENCIES", "items": [
            "Murni Go Standard Library (net/http, crypto/hmac, crypto/subtle, sync)",
            "Membuktikan mekanisme otentikasi sesi dan penolakan payload secara presisi",
        ], "type": "bullet"},
        {"title": "KOMPONEN MODULAR", "items": [
            "internal/cors: Evaluasi origin, header preflight OPTIONS, allow-credentials",
            "internal/csrf: HMAC token generator, token validator, metadata middleware",
        ], "type": "bullet"},
    ],
    mental_model="Modular middleware isolates transport negotiation from state-changing business handlers.",
    source_path="labs/36-cors-and-csrf/internal/bank/app.go",
    slide_num="04 / 24",
    narrative="Lab 36 mengimplementasikan bank service lengkap di internal/bank/app.go. Terdapat endpoint transfer yang sengaja dibiarkan rentan hanya dengan middleware CORS, serta varian endpoint terlindungi yang dibungkus oleh Anti-CSRF Token, Fetch Metadata, dan Custom Header middleware untuk membandingkan perilaku eksekusinya secara langsung.",
    title_en="Lab 36 Banking Architecture Topology",
    narrative_en="Lab 36 hosts a banking service comparing vulnerable and protected endpoints. The vulnerable endpoint relies solely on CORS, whereas protected endpoints compose HMAC token verification, Fetch Metadata, and Custom Header middleware."
)

# Slide 5: Simple Request vs Preflight OPTIONS
add_slide(
    categories=["SPESIFIKASI", "WHATWG FETCH"],
    title="Simple Request vs Preflight OPTIONS",
    subtitle="Kapan browser meminta izin sebelum mengirim dan kapan langsung mengeksekusi",
    content_items=[
        {"title": "SIMPLE REQUEST (NO PREFLIGHT)", "items": [
            "Method: GET, HEAD, atau POST",
            "Content-Type: application/x-www-form-urlencoded, multipart, text/plain",
            "Tanpa custom headers; browser LANGSUNG mengirim request ke backend",
            "Target utama serangan formulir HTML CSRF klasik",
        ], "type": "bullet"},
        {"title": "PREFLIGHT REQUEST (OPTIONS)", "items": [
            "Method: PUT, DELETE, PATCH atau Content-Type: application/json",
            "Memakai custom header seperti X-CSRF-Token atau Authorization",
            "Browser mengirim OPTIONS terlebih dahulu untuk meminta izin CORS",
            "Jika preflight gagal, actual request TIDAK AKAN PERNAH dikirim",
        ], "type": "bullet"},
        {"title": "IMPLIKASI KEAMANAN", "items": [
            "Penyerang CSRF memanfaatkan celah simple request agar server langsung mengeksekusi",
        ], "type": "bullet"},
    ],
    mental_model="Simple requests fire immediately without asking; preflight requests ask permission first.",
    source_path="labs/36-cors-and-csrf/internal/cors/middleware.go:44-92",
    slide_num="05 / 24",
    narrative="Standar WHATWG Fetch membedakan request menjadi simple request dan non-simple request. Simple request menggunakan method standar dan encoding form tanpa custom header, sehingga dikirimkan langsung oleh browser tanpa konfirmasi awal. Sebaliknya, non-simple request memicu preflight OPTIONS. Penyerang CSRF memanfaatkan sifat simple request untuk mengeksekusi mutasi server seketika.",
    title_en="Simple Request vs Preflight OPTIONS",
    narrative_en="The WHATWG Fetch standard bifurcates traffic into simple and preflighted requests. Simple requests fire immediately without permission checks, making them the primary vehicle for CSRF exploits. Non-simple requests mandate OPTIONS preflight before sending actual payloads."
)

# Slide 6: Anatomi Serangan CSRF
add_slide(
    categories=["ANALISIS SERANGAN", "ANATOMI CSRF"],
    title="Anatomi Serangan Cross-Site Request Forgery",
    subtitle="Bagaimana situs penyerang mengeksploitasi cookie korban yang tersimpan di browser",
    content_items=[
        {"title": "LANGKAH EKSPLOITASI", "items": [
            "1. Korban login ke bank target dan memegang cookie sesi aktif di browser",
            "2. Korban terpancing membuka tab situs penyerang (evil-attacker.com)",
            "3. Halaman jahat menjalankan form submit tersembunyi via script JavaScript",
            "4. Browser otomatis menyertakan cookie sesi korban ke request POST bank",
            "5. Server backend memvalidasi cookie sesi dan mengeksekusi transfer dana",
        ], "type": "diagram"},
        {"title": "KENAPA CORS GAGAL MENCEGAH?", "items": [
            "Server mengeksekusi mutasi data sebelum mengirimkan HTTP response",
            "Browser memblokir pembacaan respons oleh penyerang karena origin mismatch",
            "Penyerang tidak butuh membaca respons; saldo korban sudah terlanjur lenyap",
        ], "type": "bullet"},
        {"title": "KERUGIAN", "items": [
            "Pemindahan aset finansial atau pembajakan email akun korban secara diam-diam",
        ], "type": "bullet"},
    ],
    mental_model="CSRF is blind execution: the attacker does not need to read the reply to steal the prize.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="06 / 24",
    narrative="Serangan CSRF mengeksploitasi kepercayaan server terhadap browser korban. Ketika form submit otomatis dipicu dari situs penyerang, browser melampirkan cookie sesi milik korban. Server bank memproses transfer tersebut hingga tuntas. Walaupun browser kemudian menolak memberikan isi response kepada penyerang karena aturan CORS, mutasi saldo sudah selesai dicatat di database.",
    title_en="Anatomy of Cross-Site Request Forgery",
    narrative_en="CSRF exploits the server's ambient trust in the user's browser session. When a hidden form submits cross-site, the browser attaches authenticated session cookies automatically. The backend executes the transfer before CORS can block response reading."
)

# Slide 7: Bukti Empiris: CORS Mengizinkan Eksekusi Backend
add_slide(
    categories=["BUKTI KODE", "INTEGRATION TEST"],
    title="Bukti Empiris: CORS Tidak Mencegah Eksekusi",
    subtitle="Test membuktikan saldo korban berkurang meskipun origin penyerang ditolak CORS",
    content_items=[
        {"title": "KODE TEST INTEGRASI", "items": [
            "req := httptest.NewRequest(POST, \"/api/transfer/vulnerable\", body)",
            "req.Header.Set(\"Origin\", \"https://evil-attacker.com\")",
            "req.AddCookie(&http.Cookie{Name: \"session_id\", Value: \"session-victim-secret\"})",
            "",
            "app.ServeHTTP(rec, req)",
            "// Response: 200 OK!",
            "// ACAO Header: Kosong (ditolak untuk evil-attacker.com)",
            "// Saldo Korban: Turun dari $1000 menjadi $700!",
        ], "type": "code"},
        {"title": "FAKTA TERVERIFIKASI", "items": [
            "Status HTTP: 200 OK (handler bisnis transfer tetap tereksekusi)",
            "Header Access-Control-Allow-Origin: Omitted (tidak diberikan ke attacker)",
            "Saldo akun korban: Berkurang $300 ke rekening penyerang",
        ], "type": "bullet"},
    ],
    mental_model="Backend execution occurs before header evaluation. No CORS header does not mean no mutation.",
    source_path="labs/36-cors-and-csrf/tests/integration_test.go:54-87",
    slide_num="07 / 24",
    narrative="Pengujian TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution membuktikan realitas ini secara definitif. Request POST dari origin evil-attacker.com berhasil memicu pemotongan saldo korban sebesar $300 dan menambah saldo penyerang. Meskipun middleware CORS tidak menyertakan header Access-Control-Allow-Origin, mutasi transaksi telah tuntas.",
    title_en="Empirical Proof: CORS Allows Backend Execution",
    narrative_en="Integration tests confirm that cross-origin simple POST requests mutate backend state regardless of CORS configuration. The vulnerable transfer endpoint returns 200 OK and deducts $300 from the victim despite omitting Access-Control-Allow-Origin headers."
)

# Slide 8: Implementasi Middleware CORS Standar W3C
add_slide(
    categories=["BUKTI KODE", "CORS SPEC"],
    title="Implementasi Middleware CORS Standar W3C",
    subtitle="Menangani evaluasi origin, preflight OPTIONS, dan aturan kredensial",
    content_items=[
        {"title": "KODE MIDDLEWARE CORS", "items": [
            "if !m.isOriginAllowed(origin) {",
            "    if r.Method == http.MethodOptions {",
            "        w.WriteHeader(http.StatusForbidden); return",
            "    }",
            "    next.ServeHTTP(w, r) // Non-preflight tetap diteruskan!",
            "    return",
            "}",
            "if m.config.AllowCredentials {",
            "    w.Header().Set(\"Access-Control-Allow-Origin\", origin)",
            "    w.Header().Set(\"Access-Control-Allow-Credentials\", \"true\")",
            "}",
        ], "type": "code"},
        {"title": "ATURAN SPESIFIKASI", "items": [
            "Jika AllowCredentials: true, wildcard '*' dilarang keras oleh browser",
            "Origin yang tidak ada di whitelist ditolak 403 HANYA pada preflight OPTIONS",
            "Header Vary: Origin wajib disetel untuk mencegah cache poisoning di CDN",
        ], "type": "bullet"},
    ],
    mental_model="Standard CORS middleware allows unwhitelisted simple requests through to the next handler.",
    source_path="labs/36-cors-and-csrf/internal/cors/middleware.go:44-78",
    slide_num="08 / 24",
    narrative="Di dalam internal/cors/middleware.go, kepatuhan terhadap spesifikasi W3C Fetch terlihat jelas. Ketika origin tidak terdaftar, request preflight OPTIONS langsung ditolak dengan HTTP 403. Namun, untuk request non-preflight (seperti simple POST), middleware harus meneruskan request ke handler berikutnya tanpa menyertakan header ACAO.",
    title_en="W3C Spec-Compliant CORS Middleware",
    narrative_en="The CORS middleware adheres strictly to W3C Fetch rules. Disallowed origins receiving preflight OPTIONS are rejected with 403 Forbidden. However, non-preflight simple requests pass through to subsequent handlers without CORS headers."
)

# Slide 9: Bahaya Access-Control-Allow-Origin: * dan Kredensial
add_slide(
    categories=["KEAMANAN", "CORS CAVEAT"],
    title="Aturan Kredensial & Wildcard Origin",
    subtitle="Larangan spesifikasi: mengapa wildcard '*' ditolak saat membawa kredensial",
    content_items=[
        {"title": "LARANGAN SPESIFIKASI FETCH", "items": [
            "Kombinasi ACAO: * dan ACAC: true dianggap ILEGAL oleh browser modern",
            "Browser akan menolak memberikan respons ke JavaScript client",
            "Mencegah kebocoran data terotentikasi secara tidak sengaja ke seluruh web",
        ], "type": "bullet"},
        {"title": "IMPLEMENTASI VALIDASI GO", "items": [
            "TestCORS_Credentials_With_Wildcard_DisallowedInSpec memvalidasi aturan ini",
            "Jika AllowCredentials aktif, server wajib merefleksikan origin spesifik pemanggil",
            "Server tidak boleh memancarkan '*' jika request menyertakan cookie otentikasi",
        ], "type": "bullet"},
        {"title": "KESALAHAN UMUM", "items": [
            "Merefleksikan Origin request secara membabi buta tanpa whitelist sama berbahayanya dengan wildcard",
        ], "type": "bullet"},
    ],
    mental_model="Never pair wildcard origins with ambient credentials; explicit whitelisting is non-negotiable.",
    source_path="labs/36-cors-and-csrf/internal/cors/middleware_test.go:42-63",
    slide_num="09 / 24",
    narrative="Spesifikasi web secara eksplisit melarang penggunaan Access-Control-Allow-Origin: * jika Access-Control-Allow-Credentials bernilai true. Jika server melanggar aturan ini, browser otomatis menolak respons. Backend wajib merefleksikan origin secara selektif hanya jika origin tersebut telah terverifikasi dalam whitelist aman.",
    title_en="Credential Rules & Wildcard Origins",
    narrative_en="Web standards forbid pairing Access-Control-Allow-Origin: * with Access-Control-Allow-Credentials: true. If credentials are included, the server must explicitly whitelist and echo the validated origin, accompanied by a Vary: Origin header."
)

# Slide 10: Arsitektur Pertahanan CSRF Multi-Layer
add_slide(
    categories=["ARSITEKTUR", "DEFENSE IN DEPTH"],
    title="Arsitektur Pertahanan CSRF Multi-Layer",
    subtitle="Menggabungkan 4 lapisan pertahanan untuk menjamin keamanan mutasi state",
    content_items=[
        {"title": "LAPISAN 1: SIGNED CSRF TOKEN", "items": [
            "HMAC-SHA256 Token terikat Session ID, TTL, dan nonce (Double-Submit)",
        ], "type": "bullet"},
        {"title": "LAPISAN 2: FETCH METADATA", "items": [
            "Evaluasi header Sec-Fetch-Site (tolak jika cross-site pada method mutasi)",
        ], "type": "bullet"},
        {"title": "LAPISAN 3: CUSTOM API HEADER", "items": [
            "Wajibkan header kustom X-Requested-With (kebal dari submit form HTML)",
        ], "type": "bullet"},
        {"title": "LAPISAN 4: SAMESITE COOKIE", "items": [
            "Konfigurasi cookie SameSite=Lax atau SameSite=Strict pada layer auth",
        ], "type": "bullet"},
    ],
    mental_model="Defense-in-depth: combine cryptographic tokens, browser metadata, and cookie boundaries.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="10 / 24",
    narrative="Pertahanan CSRF yang tangguh tidak boleh bergantung pada satu mekanisme tunggal. Pendekatan defense-in-depth menggabungkan Signed Double-Submit Token berbasis HMAC, verifikasi Fetch Metadata dari browser, penegakan header custom untuk API AJAX, serta konfigurasi cookie SameSite yang ketat pada sisi klien.",
    title_en="Multi-Layer Anti-CSRF Defense Architecture",
    narrative_en="Robust anti-CSRF defense avoids single points of failure. The defense-in-depth model composes HMAC-signed double-submit tokens, Fetch Metadata verification, API custom header enforcement, and strict SameSite cookie policies."
)

# Slide 11: Signed Double-Submit CSRF Token
add_slide(
    categories=["KRIPTOGRAFI", "TOKEN PATTERN"],
    title="Signed Double-Submit Token (HMAC-SHA256)",
    subtitle="Membuat token stateless tanpa beban memori session store di server",
    content_items=[
        {"title": "STRUKTUR PAYLOAD TOKEN", "items": [
            "Payload: sessionID : timestamp : nonce",
            "Signature: HMAC_SHA256(secret_key, payload)",
            "Token Wire: base64(payload : base64(signature))",
        ], "type": "diagram"},
        {"title": "KEUNGGULAN POLA INI", "items": [
            "Stateless: Server tidak perlu menyimpan token di Redis atau database",
            "Cryptographically Bound: Terikat pada ID sesi pengguna aktif",
            "Anti-Replay & Expiring: Dilengkapi batas waktu kedaluwarsa (TTL)",
        ], "type": "bullet"},
        {"title": "KEBAL SUBDOMAIN TAKEOVER", "items": [
            "Naive Double-Submit cookie dapat disuntik attacker; HMAC signature menggagalkannya",
        ], "type": "bullet"},
    ],
    mental_model="HMAC-signed double-submit tokens provide stateless security immune to cookie injection.",
    source_path="labs/36-cors-and-csrf/internal/csrf/token.go:43-59",
    slide_num="11 / 24",
    narrative="Pola Signed Double-Submit Token menggabungkan efisiensi stateless dengan integritas kriptografis. Payload token berisi Session ID, timestamp pembuatan, dan random nonce 16-byte yang ditandatangani dengan HMAC-SHA256. Penyerang yang menguasai subdomain tidak dapat memalsukan token karena tidak memiliki secret key server.",
    title_en="Signed Double-Submit CSRF Token (HMAC)",
    narrative_en="The Signed Double-Submit pattern delivers stateless CSRF defense without server memory overhead. Tokens combine Session ID, timestamp, and random nonce signed via HMAC-SHA256, protecting against cookie injection and subdomain compromise."
)

# Slide 12: Generator & Validator Token di Go
add_slide(
    categories=["BUKTI KODE", "HMAC VALIDATOR"],
    title="Implementasi Generate & Validate Token",
    subtitle="Validasi tanda tangan, Session ID, masa aktif, dan constant-time compare",
    content_items=[
        {"title": "KODE VALIDASI TOKEN", "items": [
            "parts := strings.Split(string(decoded), \":\")",
            "sessionID, tsStr, nonceStr, sigB64 := parts[0], parts[1], parts[2], parts[3]",
            "if sessionID != expectedSessionID { return ErrInvalidToken }",
            "",
            "ts, err := strconv.ParseInt(tsStr, 10, 64)",
            "if err != nil || time.Since(time.Unix(ts, 0)) > tm.ttl {",
            "    return ErrExpiredToken",
            "}",
            "expectedSig := hmac.New(sha256.New, tm.secret).Sum(payload)",
            "if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {",
            "    return ErrInvalidToken",
            "}",
        ], "type": "code"},
        {"title": "ANTI-TIMING ATTACK", "items": [
            "subtle.ConstantTimeCompare mengeliminasi kebocoran informasi via waktu eksekusi CPU",
        ], "type": "bullet"},
    ],
    mental_model="Constant-time verification prevents side-channel timing attacks on cryptographic signatures.",
    source_path="labs/36-cors-and-csrf/internal/csrf/token.go:62-110",
    slide_num="12 / 24",
    narrative="Method ValidateToken() di internal/csrf/token.go menjalankan validasi berlapis. Token diurai menjadi empat komponen. Validator memeriksa kecocokan Session ID, memastikan timestamp belum kedaluwarsa, dan menghitung ulang HMAC signature. subtle.ConstantTimeCompare digunakan untuk mencegah penyerang membongkar signature via timing attack.",
    title_en="Token Generation & Validation in Go",
    narrative_en="The ValidateToken method validates session binding, timestamp expiration, and cryptographic authenticity. Leveraging crypto/subtle.ConstantTimeCompare eliminates timing side-channel leakage during signature verification."
)

# Slide 13: Middleware Validasi Token CSRF
add_slide(
    categories=["BUKTI KODE", "MIDDLEWARE CSRF"],
    title="Middleware Proteksi CSRF di HTTP Handler",
    subtitle="Mengabaikan method aman dan memvalidasi token dari header maupun form body",
    content_items=[
        {"title": "KODE CSRF HANDLER", "items": [
            "if r.Method == GET || r.Method == HEAD || r.Method == OPTIONS {",
            "    next.ServeHTTP(w, r); return",
            "}",
            "cookie, err := r.Cookie(m.cookieName)",
            "if err != nil || cookie.Value == \"\" {",
            "    http.Error(w, \"Unauthorized: missing cookie\", 401); return",
            "}",
            "token := r.Header.Get(\"X-CSRF-Token\")",
            "if token == \"\" { token = r.PostFormValue(\"csrf_token\") }",
            "",
            "if err := m.tokenManager.ValidateToken(token, cookie.Value); err != nil {",
            "    http.Error(w, \"Forbidden: CSRF validation failed\", 403); return",
            "}",
        ], "type": "code"},
        {"title": "FLEXIBILITAS TRANSPORT", "items": [
            "Mendukung AJAX SPA via header X-CSRF-Token",
            "Mendukung form HTML tradisional via field csrf_token",
        ], "type": "bullet"},
    ],
    mental_model="Safeguard state mutation: exempt safe methods, extract token from header/form, verify against session.",
    source_path="labs/36-cors-and-csrf/internal/csrf/middleware.go:24-53",
    slide_num="13 / 24",
    narrative="Middleware CSRF bertindak sebagai gerbang pengaman mutasi data. Request dengan method aman (GET, HEAD, OPTIONS) langsung dilewatkan. Untuk method mutasi (POST, PUT, DELETE), middleware mengekstrak cookie sesi dan mengambil token CSRF dari header atau form field. Jika token tidak cocok atau kedaluwarsa, server menolak dengan HTTP 403 Forbidden.",
    title_en="Anti-CSRF Middleware in HTTP Handler",
    narrative_en="The anti-CSRF middleware wraps sensitive HTTP routes. It passes safe idempotent methods transparently while interrogating state-mutating requests. It extracts tokens from either custom headers or form payloads and rejects mismatches with 403 Forbidden."
)

# Slide 14: Fetch Metadata: Sec-Fetch-Site
add_slide(
    categories=["BROWSER MODERN", "FETCH METADATA"],
    title="Fetch Metadata: Memanfaatkan Sec-Fetch-Site",
    subtitle="Header antipemalsuan browser untuk mengidentifikasi asal usul inisiasi request",
    content_items=[
        {"title": "NILAI SEC-FETCH-SITE", "items": [
            "same-origin : Request diinisiasi dari origin yang sama persis (Aman)",
            "same-site   : Request diinisiasi dari subdomain dalam registable domain sama",
            "cross-site  : Request diinisiasi dari situs luar sepenuhnya (Waspada!)",
            "none        : Navigasi langsung oleh user (misal ketik URL di address bar)",
        ], "type": "bullet"},
        {"title": "KODE MIDDLEWARE FETCH METADATA", "items": [
            "fetchSite := r.Header.Get(\"Sec-Fetch-Site\")",
            "if fetchSite == \"cross-site\" {",
            "    if r.Method != GET && r.Method != HEAD && r.Method != OPTIONS {",
            "        http.Error(w, \"Forbidden: Cross-site rejected\", 403)",
            "        return",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "KEUNGGULAN UTAMA", "items": [
            "Header disetel langsung oleh kernel browser; tidak dapat diubah oleh JavaScript frontend",
        ], "type": "bullet"},
    ],
    mental_model="Browser-controlled metadata tells the unvarnished truth about request origin.",
    source_path="labs/36-cors-and-csrf/internal/csrf/middleware.go:55-70",
    slide_num="14 / 24",
    narrative="Fetch Metadata merupakan standar W3C modern di mana browser melampirkan header Sec-Fetch-Site pada setiap request keluar. Header ini bernilai cross-site jika dipicu dari domain luar. Karena header ini dikelola langsung oleh mesin browser dan dilindungi dari manipulasi JavaScript, backend dapat menolak request mutasi cross-site tanpa perlu token sama sekali.",
    title_en="Fetch Metadata: Leveraging Sec-Fetch-Site",
    narrative_en="Fetch Metadata provides browser-asserted contextual headers that JavaScript cannot tamper with. A Sec-Fetch-Site value of cross-site on state-changing methods allows backend middleware to drop malicious requests instantly."
)

# Slide 15: Custom Header Enforcement (X-Requested-With)
add_slide(
    categories=["DEFENSE TRICK", "CUSTOM HEADER"],
    title="Custom Header Enforcement untuk API",
    subtitle="Memanfaatkan keterbatasan form HTML standar untuk mematikan eksploitasi CSRF",
    content_items=[
        {"title": "MEKANIKA TEKNIS", "items": [
            "Formulir HTML standar (<form>) TIDAK BISA menyetel custom HTTP headers",
            "Hanya JavaScript terprogram (fetch / axios) yang dapat menambahkan custom header",
            "Menambahkan custom header otomatis memaksa browser mengirim preflight OPTIONS",
        ], "type": "bullet"},
        {"title": "KODE MIDDLEWARE CUSTOM HEADER", "items": [
            "func RequireCustomHeaderMiddleware(headerKey, expectedValue string) {",
            "    val := r.Header.Get(headerKey)",
            "    if val == \"\" || !strings.EqualFold(val, expectedValue) {",
            "        http.Error(w, \"Forbidden: Missing custom header\", 403); return",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "BATASAN STRATEGI", "items": [
            "Sangat efektif untuk REST API murni, namun tidak dapat dipakai pada aplikasi berbasis server-side rendered form",
        ], "type": "bullet"},
    ],
    mental_model="HTML forms cannot emit custom headers. Require a custom header, defeat form CSRF by definition.",
    source_path="labs/36-cors-and-csrf/internal/csrf/middleware.go:72-86",
    slide_num="15 / 24",
    narrative="Salah satu pertahanan paling sederhana namun ampuh untuk API backend adalah mewajibkan header kustom seperti X-Requested-With: XMLHttpRequest. Karena tag form HTML standar di browser tidak memiliki kemampuan menambahkan custom header HTTP, seluruh upaya eksploitasi form CSRF klasik langsung gagal sebelum mencapai handler bisnis.",
    title_en="API Custom Header Enforcement",
    narrative_en="Requiring custom headers like X-Requested-With completely neutralizes classic HTML form CSRF. Standard HTML form tags cannot attach custom HTTP headers, immediately breaking exploit attempts against API endpoints."
)

# Slide 16: Peran Atribut Cookie SameSite
add_slide(
    categories=["COOKIE SECURITY", "SAMESITE ATTRIBUTE"],
    title="Peran Atribut Cookie SameSite: Lax vs Strict",
    subtitle="Membendung pengiriman kredensial otomatis pada batas navigasi browser",
    content_items=[
        {"title": "SAMESITE=STRICT", "items": [
            "Cookie tidak pernah dikirim pada request lintas situs dalam kondisi apapun",
            "Bahkan saat user mengklik link dari email atau search engine, cookie ditahan",
            "Sangat aman, namun mengorbankan kenyamanan login navigasi eksternal",
        ], "type": "bullet"},
        {"title": "SAMESITE=LAX (DEFAULT BROWSER MODERN)", "items": [
            "Cookie ditahan pada cross-site POST (memblokir serangan form CSRF)",
            "Cookie tetap dikirim pada top-level GET navigation (misal klik link)",
            "Wajib: Endpoint GET harus idempotent dan tidak boleh memutasi state!",
        ], "type": "bullet"},
        {"title": "SAMESITE=NONE", "items": [
            "Cookie selalu dikirim lintas situs; wajib menyertakan flag Secure (HTTPS)",
        ], "type": "bullet"},
    ],
    mental_model="SameSite=Lax blocks cross-origin POST cookies; keep GET endpoints strictly read-only.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="16 / 24",
    narrative="Atribut SameSite pada cookie menentukan kapan browser diizinkan melampirkan cookie sesi. SameSite=Lax yang menjadi default browser modern menahan cookie pada request POST lintas situs. Namun, karena cookie tetap dikirim pada navigasi top-level GET, sistem tidak boleh melakukan mutasi saldo atau perubahan data melalui method GET.",
    title_en="Cookie SameSite Attributes: Lax vs Strict",
    narrative_en="SameSite cookie attributes govern ambient credential transmission. SameSite=Lax withholds cookies on cross-origin POST requests while transmitting them on top-level GET links. Consequently, GET endpoints must remain strictly read-only and idempotent."
)

# Slide 17: Eksekusi Demo & Bukti Terminal
add_slide(
    categories=["HASIL PENGUJIAN", "DEMO RUN"],
    title="Eksekusi Demo: Bukti Terminal Langsung",
    subtitle="Simulasi nyata serangan dan pertahanan pada cmd/demo/main.go",
    content_items=[
        {"title": "OUTPUT RUN CMD/DEMO", "items": [
            "[1] Initial State: Victim: $1000 | Attacker: $50",
            "[2] DEMO: Attack on Vulnerable Endpoint (Origin: https://evil.com)",
            "    HTTP Response Code: 200 | CORS Header: \"\"",
            "    [RESULT] Transfer EXECUTED! Victim: $600 | Attacker: $450",
            "[3] DEMO: Attack on Protected Endpoint (No CSRF Token)",
            "    HTTP Response Code: 403",
            "    [RESULT] Transfer BLOCKED! Victim: $600 | Attacker: $450",
            "[4] DEMO: Legitimate Flow with Valid CSRF Token",
            "    Obtained Signed Token: c2Vzc2lvbi12aWN0aW0tc2Vjc...",
            "    [RESULT] Transfer Succeeded! Victim: $500 | Attacker: $550",
        ], "type": "code"},
    ],
    mental_model="Terminal output confirms: vulnerable endpoints fall to simple POST; protected endpoints hold firm.",
    source_path="labs/36-cors-and-csrf/cmd/demo/main.go",
    slide_num="17 / 24",
    narrative="Program demo di cmd/demo/main.go memperlihatkan realitas serangan secara nyata. Pada endpoint rentan, serangan transfer $400 dari evil.com sukses dieksekusi menghasilkan saldo korban $600. Namun, ketika serangan diarahkan ke endpoint terlindungi tanpa token CSRF, server mengembalikan status 403 Forbidden dan saldo korban tetap aman.",
    title_en="Demo Execution: Direct Terminal Evidence",
    narrative_en="The demonstration in cmd/demo/main.go provides conclusive terminal evidence. A cross-origin exploit against the vulnerable route executes successfully, while the protected route rejects unauthenticated mutations with HTTP 403 Forbidden."
)

# Slide 18: Verifikasi Test Suite & Race Detector
add_slide(
    categories=["VERIFIKASI TEST", "TEST SUITE"],
    title="Verifikasi Test Suite & Bebas Race Condition",
    subtitle="13 test case diverifikasi 100% lulus dengan Go Race Detector aktif",
    content_items=[
        {"title": "HASIL EKSEKUSI PENGUJIAN", "items": [
            "PASS: TestBankTransfer_VulnerableEndpoint       (0.00s)",
            "PASS: TestCORS_Preflight_Success               (0.00s)",
            "PASS: TestCORS_Credentials_With_Wildcard       (0.00s)",
            "PASS: TestTokenManager_GenerateAndValidate     (0.00s)",
            "PASS: TestTokenManager_ExpiredToken            (0.02s)",
            "PASS: TestIntegration_CORS_Does_Not_Prevent    (0.00s)",
            "PASS: TestIntegration_CSRF_Token_Prevents      (0.00s)",
            "PASS: TestIntegration_SecFetchSite_Protection  (0.00s)",
            "PASS: TestIntegration_Concurrency_RaceCondition (0.00s)",
            "PASS: go test -race ./... (Race-Clean Verified)",
        ], "type": "code"},
        {"title": "UJI KONKURENSI", "items": [
            "20 worker simultan melakukan pembuatan dan validasi token tanpa tabrakan memori",
        ], "type": "bullet"},
    ],
    mental_model="Deterministic test assertions prove security invariants and concurrency safety.",
    source_path="labs/36-cors-and-csrf/tests/integration_test.go",
    slide_num="18 / 24",
    narrative="Seluruh 13 pengujian unit dan integrasi lulus dengan Go Race Detector aktif. Pengujian mencakup verifikasi penolakan origin CORS, validasi tanda tangan HMAC, penolakan token kedaluwarsa, blokade cross-site Sec-Fetch-Site, pencegahan token reuse lintas sesi, hingga stabilitas konkurensi di bawah beban paralel.",
    title_en="Test Suite Verification & Race Safety",
    narrative_en="All 13 unit and integration test cases pass cleanly under the Go race detector. Tests validate origin whitelisting, cryptographic signatures, expiration timeouts, Sec-Fetch-Site drops, cross-session reuse barriers, and thread concurrency."
)

# Slide 19: Matriks Komparasi 5 Dimensi Pertahanan
add_slide(
    categories=["TABEL EVALUASI", "PERBANDINGAN"],
    title="Matriks Evaluasi: Mekanisme Pertahanan",
    subtitle="Perbandingan komparatif 5 pendekatan keamanan web lintas origin",
    content_items=[
        {"title": "TABEL PERBANDINGAN MEKANISME", "items": [
            "Mekanisme Keamanan | Target Proteksi | Lokasi Penegakan | Tangkal Form CSRF? | Ketergantungan Browser",
            "CORS Middleware    | Relaksasi Read  | Browser & API    | TIDAK              | Penuh (Fetch Spec)",
            "Signed CSRF Token  | Mutasi Write    | Backend Server   | YA                 | Nol (Bekerja di semua)",
            "Fetch Metadata     | Mutasi Write    | Gateway / App    | YA                 | Butuh browser modern",
            "Custom Header API  | Mutasi Write    | Gateway / App    | YA                 | Nol (API Client)",
            "Cookie SameSite    | Ambient Cookie  | Browser Engine   | YA (Lax/Strict)    | Butuh browser modern",
        ], "type": "table"},
        {"title": "KESIMPULAN EVALUASI", "items": [
            "Hanya Signed CSRF Token yang memberikan pertahanan 100% mandiri di sisi backend",
            "Lapisan lain memperkuat pertahanan di layer browser dan network gateway",
        ], "type": "bullet"},
    ],
    mental_model="Only cryptographic tokens provide absolute backend independence; other layers reinforce the perimeter.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="19 / 24",
    narrative="Tabel komparasi memperjelas posisi masing-masing mekanisme. CORS sama sekali tidak menangkal form CSRF karena fokusnya adalah relaksasi pembacaan respons. Signed CSRF Token adalah satu-satunya mekanisme yang 100% independen dari peramban karena divalidasi langsung oleh logika kriptografis backend.",
    title_en="Comparative Matrix: 5 Security Dimensions",
    narrative_en="The comparative matrix clarifies defense boundaries. CORS is ineffective against form CSRF as it targets read access. Signed CSRF tokens remain the sole mechanism providing 100% server-side autonomy independent of browser capabilities."
)

# Slide 20: Ancaman Tersembunyi: Ketika XSS Melumpuhkan Anti-CSRF
add_slide(
    categories=["ANALISIS ANCAMAN", "XSS INTERACTION"],
    title="Ancaman Tersembunyi: XSS Melumpuhkan CSRF",
    subtitle="Mengapa pertahanan CSRF terbaik sekalipun runtuh jika origin memiliki celah XSS",
    content_items=[
        {"title": "ALUR BYPASS VIA XSS", "items": [
            "1. Penyerang menemukan celah Stored / Reflected XSS pada domain target",
            "2. Script jahat penyerang dieksekusi di dalam context origin yang sah",
            "3. Karena berada di origin yang sama, SOP mengizinkan script membaca DOM",
            "4. Script melakukan fetch ke /api/csrf-token dan membaca token sah",
            "5. Script mengeksekusi POST transfer lengkap dengan token dan cookie korban",
        ], "type": "diagram"},
        {"title": "IMPLIKASI KEAMANAN", "items": [
            "XSS adalah kerentanan tingkat origin yang mengalahkan seluruh pertahanan CSRF",
            "Pertahanan CSRF wajib didampingi Content Security Policy (CSP) dan sanitasi input",
        ], "type": "bullet"},
    ],
    mental_model="XSS invalidates all anti-CSRF defenses by executing within the trusted origin context.",
    source_path="labs/36-cors-and-csrf/content/01-content-brief.md",
    slide_num="20 / 24",
    narrative="Penting untuk diingat bahwa seluruh pertahanan anti-CSRF (baik token, SameSite, maupun Fetch Metadata) runtuh jika aplikasi memiliki kerentanan Cross-Site Scripting (XSS). Jika penyerang dapat menyuntikkan script ke origin target, script tersebut memiliki hak istimewa SOP untuk membaca token CSRF secara sah dan mengeksekusi transaksi atas nama korban.",
    title_en="Hidden Threat: XSS Defeats Anti-CSRF",
    narrative_en="All anti-CSRF mechanisms collapse if an application contains Cross-Site Scripting vulnerabilities. An injected script executes within the trusted origin context, allowing it to read CSRF tokens legally under SOP and authorize malicious transactions."
)

# Slide 21: Batasan Implementasi Lab 36
add_slide(
    categories=["BATAS LAB", "DISCLAIMER"],
    title="Batasan Model & Implementasi Lab 36",
    subtitle="Penyederhanaan pedagogis untuk menonjolkan prinsip mekanika protokol",
    content_items=[
        {"title": "BATAS 1: SECRET KEY IN-MEMORY", "items": [
            "Kunci rahasia HMAC menggunakan byte slice statis dalam kode demo",
            "Produksi wajib menggunakan AWS Secrets Manager, Vault, dan key rotation",
        ], "type": "bullet"},
        {"title": "BATAS 2: DELIMITER STRING PARSER", "items": [
            "Token memisahkan field menggunakan karakter titik dua (':')",
            "Jika sessionID sistem memuat karakter ':', parsing token akan terganggu",
            "Produksi memerlukan format terstruktur seperti JSON, Protobuf, atau hex encoding",
        ], "type": "bullet"},
        {"title": "BATAS 3: LINGKUP BROWSER", "items": [
            "SOP, CORS, dan CSRF murni berada pada lingkup runtime peramban web",
            "HTTP client non-browser (cURL, script backend) tidak terikat aturan ini",
        ], "type": "bullet"},
    ],
    mental_model="A lab isolates cryptographic mechanics; production demands secret management and structured payloads.",
    source_path="labs/36-cors-and-csrf/README.md",
    slide_num="21 / 24",
    narrative="Lab ini menyederhanakan beberapa aspek infrastruktur untuk menjaga fokus edukasi. Kunci rahasia HMAC di-hardcode dalam memori, dan token diurai menggunakan string delimiter sederhana. Dalam sistem produksi riil, kunci HMAC wajib dikelola melalui key management service dan token sebaiknya menggunakan format terstruktur seperti JSON Web Signature (JWS) atau Protobuf.",
    title_en="Lab 36 Model Boundaries & Constraints",
    narrative_en="This lab simplifies production infrastructure for clarity. HMAC secret keys reside in memory and tokens parse via colon delimiters. Production deployments demand managed key vaults, automatic rotation, and structured token payloads."
)

# Slide 22: Panduan Checklist Audit Keamanan Backend
add_slide(
    categories=["PANDUAN", "CHECKLIST KEAMANAN"],
    title="Checklist Audit Keamanan CORS & CSRF",
    subtitle="9 langkah verifikasi wajib sebelum merilis backend ke lingkungan produksi",
    content_items=[
        {"title": "VERIFIKASI CORS & HTTP", "items": [
            "[ ] Method GET, HEAD, OPTIONS dipastikan bebas dari mutasi data database",
            "[ ] CORS whitelist hanya mencantumkan domain tepercaya (tanpa refleksi liar)",
            "[ ] AllowCredentials: true TIDAK PERNAH dipadukan dengan origin '*'",
        ], "type": "bullet"},
        {"title": "VERIFIKASI ANTI-CSRF", "items": [
            "[ ] Seluruh rute mutasi state (POST/PUT/DELETE) dilindungi Anti-CSRF Token",
            "[ ] Token CSRF ditandatangani HMAC-SHA256, terikat Session ID, dan memiliki TTL",
            "[ ] Validasi signature menggunakan subtle.ConstantTimeCompare",
            "[ ] Cookie otentikasi dipasangi flag HttpOnly, Secure, dan SameSite=Lax",
            "[ ] API gateway memverifikasi header Sec-Fetch-Site dan Custom Header",
            "[ ] Aplikasi terlindungi dari celah XSS via Content Security Policy (CSP)",
        ], "type": "bullet"},
    ],
    mental_model="Verify read-only GETs, strict origin whitelists, signed tokens, and CSP before shipping.",
    source_path="labs/36-cors-and-csrf/content/02-master-draft.md",
    slide_num="22 / 24",
    narrative="Gunakan checklist ini saat melakukan security audit pada sistem web Anda. Pastikan endpoint GET bebas dari efek samping mutasi data, jangan pernah menggabungkan wildcard CORS dengan kredensial, wajibkan validasi token HMAC yang terikat ke sesi pada setiap mutasi POST, dan perkuat perimeter dengan header Sec-Fetch-Site serta cookie SameSite.",
    title_en="Security Audit Checklist: CORS & CSRF",
    narrative_en="Apply this checklist during backend security audits. Guarantee that GET routes remain strictly read-only, prohibit wildcard origins with credentials, mandate session-bound HMAC tokens for mutations, and configure SameSite cookies."
)

# Slide 23: Tiga Prinsip Inti Keamanan Cross-Origin
add_slide(
    categories=["PRINSIP INTI", "KESIMPULAN"],
    title="Tiga Prinsip Inti Keamanan Cross-Origin",
    subtitle="Aksioma fundamental untuk melindungi backend dari eksploitasi browser",
    content_items=[
        {"title": "1. NEVER TRUST CLIENT BOUNDARIES FOR SERVER INTEGRITY", "items": [
            "CORS adalah aturan kepatuhan browser, bukan pelindung integritas database",
            "Validasi otorisasi dan mutasi data wajib ditegakkan sepenuhnya di sisi server",
        ], "type": "bullet"},
        {"title": "2. BIND TOKENS CRYPTOGRAPHICALLY TO SESSIONS", "items": [
            "Anti-CSRF Token harus terikat secara kriptografis ke identitas sesi aktif",
            "Mencegah eksploitasi cross-session token reuse dan pemalsuan token",
        ], "type": "bullet"},
        {"title": "3. DEFENSE IN DEPTH WINS", "items": [
            "Kombinasikan Signed Token + SameSite Cookie + Fetch Metadata",
            "Kekuatan sistem berlapis menjamin keamanan saat salah satu mekanisme gagal",
        ], "type": "bullet"},
    ],
    mental_model="CORS governs browser reads; cryptographic server tokens defend backend writes. Build multi-layered security.",
    source_path="labs/36-cors-and-csrf/content/05-key-takeaways.md",
    slide_num="23 / 24",
    narrative="Menjaga keamanan backend menuntut pemahaman arsitektural yang jelas. Jangan pernah mempercayai CORS untuk menjaga integritas database. Selalu ikat token CSRF ke ID sesi pengguna secara kriptografis, dan terapkan arsitektur pertahanan berlapis untuk menangkal berbagai vektor serangan lintas origin secara efektif.",
    title_en="Three Core Cross-Origin Security Principles",
    narrative_en="Securing backend systems demands conceptual clarity. Never rely on CORS to protect database integrity. Cryptographically bind anti-CSRF tokens to active sessions, and compose multi-layer defenses to neutralize cross-origin attack vectors."
)

# Slide 24: Closing (Dark Navy)
add_slide(
    categories=["RINGKASAN", "REPOSITORI"],
    title="Eksperimen kecil, interpretasi yang teliti.",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | 13/13 Tests PASS (Go 1.22 Race Clean) | {AUTHOR}",
    content_items=[
        {"title": "KEY RESULTS", "items": [
            "internal/cors • internal/csrf",
            "internal/bank • cmd/demo",
            "tests/integration_test.go • README.md",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Semua test PASS dengan Go 1.22 Race Detector aktif",
            "Terbukti: CORS tidak mencegah mutasi data sisi server",
        ], "type": "bullet"},
    ],
    mental_model="SOP restricts cross-origin read; CSRF exploits cross-origin write. CORS relaxes read, it never defends write.",
    source_path="labs/36-cors-and-csrf/README.md",
    slide_num="24 / 24",
    narrative="Penelitian dan pengujian Lab 36 membuktikan bahwa pemahaman tepat atas batas-batas protokol browser melindungi backend dari eksploitasi mutasi berbahaya.",
    title_en="Small experiment, careful interpretation.",
    narrative_en="Lab 36 demonstrates that precise understanding of browser security boundaries protects backend databases from unauthorized write mutations."
)

if __name__ == "__main__":
    print(f"Loaded content_def for Lab {LAB_NO}: {len(SLIDES)} slides defined.")
