#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 31: OAuth 2.0 & OpenID Connect (OIDC).
Generates 24-slide LinkedIn carousel covering Authorization Code Grant, PKCE S256 protection,
OIDC ID Token multi-step cryptographic verification, Refresh Token Rotation, and Token Family Revocation.
"""
import os

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/31-oauth2-and-oidc/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/31-oauth2-and-oidc"
COMMIT_SHA = "dec7514f81d0f9a0c26c95cd892d8c2cc51e0f5f"
LAB_NO = "31"
SLUG = "oauth2-and-oidc"
TITLE_ID = "OAuth 2.0 & OpenID Connect (OIDC)"
SUBTITLE_ID = "Autentikasi vs Otorisasi, Proteksi PKCE, dan Deteksi Replay Rotasi Token"
TITLE_EN = "OAuth 2.0 & OpenID Connect (OIDC)"
SUBTITLE_EN = "Authentication vs Authorization, PKCE Protection, and Token Rotation Replay Detection"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/31-oauth2-and-oidc/result/image-cover/cover_art.jpg"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en,
        "layout": layout
    })

# =====================================================================
# 24 SLIDE DEFINITIONS
# =====================================================================

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Visual Hero: Keamanan Protokol Identitas & Otorisasi Web Terdistribusi",
            "Mekanisme Inti: PKCE S256, ID Token Cryptographic Verification, Refresh Token Rotation",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="OAuth grants delegated authority; OIDC asserts verified identity; PKCE binds the exchange; family rotation neutralizes replay.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="OAuth 2.0 dan OpenID Connect merupakan fondasi keamanan identitas dan delegasi wewenang pada arsitektur web modern. Kesalahan paling fatal dalam sistem terdistribusi adalah mencampuradukkan kartu akses dengan kartu tanda pengenal, mengekspos token via flow usang, atau membiarkan refresh token statis dapat direplay selamanya. Lab ini membuktikan secara konkret implementasi murni Go untuk Authorization Code Flow yang dilindungi PKCE S256, verifikasi kriptografis ID Token bertingkat, serta rotasi refresh token otomatis dengan penarikan token family seketika saat terdeteksi penyusupan.",
    narrative_en="OAuth 2.0 and OpenID Connect form the bedrock of distributed identity and delegated authorization in modern architectures. The most critical failure in production systems is confusing access tokens with identity credentials, exposing tokens through deprecated flows, or permitting static refresh tokens to be replayed indefinitely. This lab demonstrates a pure Go implementation of the Authorization Code Flow fortified with PKCE S256, multi-step cryptographic ID Token verification, and single-use refresh token rotation backed by instantaneous family revocation.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "SECURITY RISKS"],
    title="Masalah Nyata: Rapuhnya Identitas & Sesi Web",
    subtitle="Ketika protokol otorisasi disalahgunakan sebagai autentikasi, celah keamanan fatal terbuka lebar.",
    content_items=[
        {"title": "PSEUDO-AUTHENTICATION", "items": [
            "Access token diperlakukan keliru sebagai bukti identitas pengguna oleh client aplikasi.",
            "Access token tidak memuat informasi audiens pemanggil, membuka celah confused deputy attack.",
            "Client tidak mampu memvalidasi apakah pengguna benar-benar hadir saat sesi diinisiasi."
        ], "type": "card"},
        {"title": "CODE INTERCEPTION HAZARD", "items": [
            "Pada public client (SPA/Mobile), authorization code dikirim via browser redirect terbuka.",
            "Aplikasi berbahaya di OS yang sama atau malware browser dapat mencuri authorization code.",
            "Tanpa cryptographic binding, penyerang dapat menukar code curian menjadi token akses sah."
        ], "type": "card"},
        {"title": "STATIC REFRESH REPLAY", "items": [
            "Refresh token berumur panjang disimpan statis tanpa mekanisme rotasi single-use.",
            "Kebocoran token melalui XSS atau network logging memberi akses persisten bagi penyerang.",
            "Authorization server tidak memiliki cara mendeteksi apakah token telah disusupi pihak ketiga."
        ], "type": "card"}
    ],
    narrative="Dalam arsitektur terdistribusi nyata, kebocoran sesi sering kali bukan diakibatkan oleh kelemahan algoritma enkripsi, melainkan kekeliruan desain alur otorisasi. Banyak tim rekayasa perangkat lunak menganggap access token OAuth 2.0 dapat digunakan langsung untuk mengenali identitas user, padahal token tersebut dirancang khusus untuk resource server, bukan client. Selain itu, authorization code yang melintasi browser redirect pada aplikasi mobile atau SPA sangat rentan disadap jika tidak diproteksi oleh bukti kepemilikan kriptografis seperti PKCE.",
    narrative_en="In production distributed architectures, session breaches rarely stem from broken cryptographic algorithms, but rather from flawed authorization flows. Many engineering teams mistakenly use OAuth 2.0 access tokens to identify users, ignoring that access tokens are scoped for resource servers rather than clients. Furthermore, authorization codes traversing browser redirects on mobile or single-page applications remain highly vulnerable to interception unless anchored by cryptographic proof-of-possession like PKCE.",
    mental_model="Access tokens prove resource permission; only cryptographically signed ID tokens prove caller identity.",
    source_path="labs/31-oauth2-and-oidc/content/01-content-brief.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "PENDEKATAN NAIF"],
    title="Mengapa Pendekatan Naif & Implicit Flow Gagal?",
    subtitle="Melewatkan validasi claims JWT dan mengandalkan implicit token fragment adalah bencana keamanan.",
    content_items=[
        {"title": "KODE VALIDASI TOKEN NAIF", "items": [
            "func HandleAuthNaive(rawJWT string) (*User, error) {",
            "    // BAHAYA FATAL: Decode payload tanpa verifikasi signature HMAC/RSA!",
            "    parts := strings.Split(rawJWT, \".\")",
            "    claimsBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])",
            "    var claims NaiveClaims",
            "    json.Unmarshal(claimsBytes, &claims)",
            "    // Tidak memeriksa iss, aud, exp, maupun nonce anti-replay!",
            "    return &User{ID: claims.Subject}, nil",
            "}"
        ], "type": "code"},
        {"title": "KEGAGALAN FUNDAMENTAL PROTOKOL", "items": [
            "Signature Bypassing: Penyerang bebas memalsukan payload JSON dan menyamar sebagai admin siapapun.",
            "Implicit Flow Flaw: Token terekspos di URL fragment riwayat browser dan HTTP Referer headers.",
            "Audience Confusion: Token sah dari Resource Server A diterima mentah-mentah oleh Resource Server B.",
            "Replay Attack: Token kedaluwarsa atau nonce usang dapat digunakan berulang kali tanpa batas."
        ], "type": "bullet"}
    ],
    narrative="Pendekatan naif kerap melakukan decode base64 pada payload JWT tanpa memverifikasi tanda tangan kriptografis dan claims wajibnya. Kode semacam ini memungkinkan penyerang membuat token palsu berisi identitas sembarang dan langsung diterima oleh sistem. Selain itu, penggunaan Implicit Flow yang dahulu populer kini telah resmi didepresiasi oleh RFC 9700 karena mengekspos token secara telanjang pada URL fragment dan history browser, menjadikannya sasaran empuk serangan Cross-Site Scripting (XSS).",
    narrative_en="The naive approach frequently decodes base64 JWT payloads directly without validating the cryptographic signature or mandatory claims. Such code allows an attacker to forge arbitrary user claims and instantly gain elevated privileges. Moreover, the legacy Implicit Flow has been formally deprecated by RFC 9700 because it leaks access tokens directly in URL hash fragments and browser histories, exposing credentials to catastrophic cross-site scripting (XSS) extraction.",
    mental_model="Unverified tokens are merely untrusted JSON; cryptographic signature and claim bounds establish integrity.",
    source_path="labs/31-oauth2-and-oidc/pkg/oidc/oidc.go:64-116",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SCHEMA"],
    title="Mental Model: Dual-Token Architecture & Family Lineage",
    subtitle="Memisahkan izin akses (Access Token) dari kartu identitas (ID Token) dan silsilah refresh token.",
    content_items=[
        {"title": "STRUKTUR DATA OIDC & SERVER METADATA", "items": [
            "type IDTokenClaims struct {",
            "    Issuer     string `json:\"iss\"`  // OIDC Provider sah",
            "    Subject    string `json:\"sub\"`  // Identitas unik user",
            "    Audience   string `json:\"aud\"`  // Client ID target",
            "    Expiration int64  `json:\"exp\"`  // Waktu kedaluwarsa",
            "    IssuedAt   int64  `json:\"iat\"`  // Waktu terbit (±5m)",
            "    Nonce      string `json:\"nonce\"`// Anti-replay challenge",
            "}",
            "type RefreshTokenMeta struct {",
            "    FamilyID  string    // Silsilah sesi tunggal",
            "    Revoked   bool      // Status single-use consumption",
            "    ExpiresAt time.Time // Kedaluwarsa absolut (30 hari)",
            "}"
        ], "type": "code"},
        {"title": "PRINSIP ARSITEKTUR KEMBAR", "items": [
            "OAuth 2.0 Gate: Mengeluarkan access_token opaque beraudience Resource Server untuk otorisasi scope.",
            "OIDC Identity Gate: Mengeluarkan id_token berformat JWT yang ditandatangani untuk client display identitas.",
            "Cryptographic Binding: Nonce mengikat token ke browser session; PKCE mengikat authorization code ke client runtime.",
            "Lineage Invariant: Satu FamilyID merepresentasikan satu rantai login aktif; penggunaan token revoked mematikan seluruh keluarga."
        ], "type": "bullet"}
    ],
    narrative="Mental model sistem ini bertumpu pada pemisahan tanggung jawab secara tegas antara otorisasi dan autentikasi. ID Token claims mendefinisikan identitas terverifikasi pengguna dengan enam perimeter verifikasi wajib: penerbit (iss), subjek (sub), penerima sah (aud), masa berlaku (exp), jendela waktu terbit (iat), dan tiket anti-replay (nonce). Di sisi backend sesi, struktur RefreshTokenMeta melacak lineage silsilah token melalui FamilyID dan penanda status Revoked untuk menjamin aturan konsumsi sekali pakai (single-use semantics).",
    narrative_en="The architectural mental model rests on a rigorous separation of concerns between authorization and authentication. ID Token claims establish verified caller identity bounded by six mandatory checkpoints: issuer (iss), subject (sub), audience (aud), expiration (exp), issuance window (iat), and anti-replay nonce. Concurrently, RefreshTokenMeta tracks session lineage via a unique FamilyID alongside a Revoked state flag to enforce strict single-use semantics across refresh token rotations.",
    mental_model="Access tokens unlock APIs; ID tokens prove identity; Family IDs trace stolen token lineage.",
    source_path="pkg/oidc/oidc.go:29-38 • pkg/server/server.go:39-46",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline End-to-End: PKCE, Token Minting, & Rotation",
    subtitle="Alur sekuensial lengkap: Auth Request → Code Issuance → Exchange → ID Verify → Rotated Refresh.",
    content_items=[
        {"title": "SEKUENSI PIPELINE IDENTITAS & OTORISASI", "items": [
            "1. Client.BuildAuthorizationRequest() → Hitung S256(verifier) = challenge & generate nonce acak.",
            "2. Client → AS: GET /authorize?response_type=code&code_challenge=...&nonce=...&scope=openid profile",
            "3. AS.Authorize() → Verifikasi client & redirect URI, simpan challenge & nonce, terbitkan AuthCode (TTL 5m).",
            "4. Client → AS: POST /token grant_type=authorization_code&code=...&code_verifier=...",
            "5. AS.ExchangeCode() → Validasi SHA256(code_verifier) == challenge. Tandai code used = true.",
            "6. AS Menerbitkan: AccessToken (opaque), IDToken (JWT HS256 bertanda tangan), & RefreshToken (FamilyID).",
            "7. Client.ParseAndVerifyIDToken() → Verifikasi signature HMAC, iss, aud, exp, dan nonce match.",
            "8. Client.RefreshTokens() → AS rotasi: tandai RT lama Revoked=true, terbitkan pasangan baru dalam FamilyID sama."
        ], "type": "diagram"}
    ],
    narrative="Pipeline otorisasi dan identitas ini memperlihatkan interaksi end-to-end antara Client, Authorization Server, dan Resource Server. Pada tahap 1 hingga 3, client menyiapkan proteksi PKCE S256 dan nonce sebelum meminta authorization code yang berumur pendek (5 menit). Pada tahap 4 hingga 6, client menukar authorization code dengan menyertakan bukti kepemilikan code_verifier asli, memungkinkan server menerbitkan access token, ID token terenkripsi, dan refresh token berkategori family. Tahap akhir memverifikasi integritas ID token serta merotasi refresh token secara aman.",
    narrative_en="This end-to-end identity and authorization pipeline maps the sequential interactions between Client, Authorization Server, and Resource Server. During stages 1 through 3, the client establishes PKCE S256 parameters and a session nonce to obtain a short-lived authorization code. In stages 4 through 6, the client exchanges the code by supplying the unhashed code_verifier proof, prompting the server to issue access tokens, signed ID tokens, and a family-tracked refresh token. Subsequent stages verify claims and rotate tokens cleanly.",
    mental_model="Challenge in front-channel; verifier in back-channel; token rotation in ongoing sessions.",
    source_path="pkg/client/client.go:36-86 • pkg/server/server.go:92-217",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="Implementasi Inti: Refresh Token Rotation & Family Revocation",
    subtitle="Pendeteksian replay otomatis pada server.go: penggunaan token basi membatalkan seluruh keluarga.",
    content_items=[
        {"title": "REFRESH TOKEN REPLAY DEFENSE", "items": [
            "if s.revokedFams[meta.FamilyID] {",
            "    return nil, ErrTokenReplayDetected // Keluarga sudah terkompromi",
            "}",
            "if meta.Revoked {",
            "    // REPLAY TERDETEKSI: Penyerang menggunakan token yang sudah dirotasi!",
            "    s.revokedFams[meta.FamilyID] = true // Cabut seluruh sesi aktif",
            "    return nil, ErrTokenReplayDetected",
            "}",
            "meta.Revoked = true // Tandai token saat ini sebagai terpakai",
            "// Terbitkan refresh token baru dengan FamilyID yang sama..."
        ], "type": "code"}
    ],
    narrative="Potongan kode ini merupakan benteng pertahanan utama terhadap serangan replay refresh token pada pkg/server/server.go. Setiap kali klien meminta pembaruan token, server memeriksa apakah token tersebut sudah pernah digunakan sebelumnya melalui flag meta.Revoked. Jika token yang sudah berstatus revoked dipresentasikan kembali, sistem menyimpulkan bahwa token tersebut telah dicuri dan sedang direplay oleh pihak jahat. Server seketika memasukkan meta.FamilyID ke dalam daftar hitam revokedFams, membatalkan seluruh sesi login terkait.",
    narrative_en="This core snippet forms the primary defense against refresh token replay attacks inside pkg/server/server.go. Whenever a client presents a refresh token, the server verifies whether the token was previously consumed via the meta.Revoked flag. If an already-revoked token is submitted, the server infers that the credential was stolen and is currently being replayed by an adversary. It immediately adds the meta.FamilyID to the revokedFams blacklist, terminating all associated active sessions.",
    mental_model="Replay of a consumed token is proof of theft; revoke the entire family immediately.",
    source_path="labs/31-oauth2-and-oidc/pkg/server/server.go:228-250",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "BUKTI"],
    title="Failure Path: Penolakan Intersepsi & Replay Serangan",
    subtitle="Log terminal membuktikan mitigasi penyerang dengan verifier salah dan replay token curian.",
    content_items=[
        {"title": "TERMINAL ATTACK MITIGATION LOG", "items": [
            "[Step 6] Test Interception Attack (Attempt code exchange with wrong verifier)...",
            "Interception Attack Blocked! Error: pkce verification failed: code_verifier does not match code_challenge",
            "[Step 7] Test Refresh Token Rotation...",
            "Rotated New Access Token : at_8c8facd2a7dce0b35b74243d16172970a30c4e94f05853bd",
            "Rotated New Refresh Token: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4d4e6503170fee590",
            "[Step 8] Test Refresh Token Replay Detection (Attacker reuses stolen Refresh Token)...",
            "Refresh Token Replay Detected! Error: refresh token reuse detected: entire token family revoked",
            "Re-testing active refresh token after family revocation...",
            "Active session revoked due to family revocation! Error: refresh token reuse detected: entire token family revoked"
        ], "type": "terminal"},
        {"title": "ANALISIS KEGAGALAN SERANGAN", "items": [
            "Langkah 6 membuktikan pencurian auth code sia-sia tanpa code_verifier asli dari client.",
            "Langkah 8 mendeteksi penyerang yang mengeksekusi replay menggunakan token rt_c62478... yang telah usang.",
            "Mitigasi Reaktif: Klien sah yang memegang token rt_a0aafd... juga ditolak, memaksa user login ulang secara aman."
        ], "type": "bullet"}
    ],
    narrative="Failure path konkret ini diambil langsung dari hasil eksekusi program demo lab cmd/demo/main.go. Pada Langkah 6, penyerang yang berhasil menyadap authorization code mencoba menukarnya menggunakan verifier palsu, namun server menolak transaksi secara mutlak dengan error ketidakcocokan PKCE. Pada Langkah 8, penyerang mencoba menggunakan refresh token usang yang telah dirotasi pada Langkah 7. Server mendeteksi upaya replay ini, menandai keluarga token sebagai terkompromi, dan membatalkan seluruh sesi termasuk token sah milik korban.",
    narrative_en="This concrete failure path is extracted directly from the demo run in cmd/demo/main.go. In Step 6, an attacker who intercepted the authorization code attempts an exchange using a bogus code_verifier, but the server rejects the request due to PKCE challenge mismatch. In Step 8, the attacker replays the consumed refresh token from Step 7. The server instantly detects the replay, marks the entire token family as compromised, and invalidates all session tokens, forcing a clean re-authentication.",
    mental_model="Cryptographic binding blocks intercepted codes; family revocation neutralizes compromised refresh tokens.",
    source_path="cmd/demo/main.go:70-109 • engineering/03-execution-result.md:98-109",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme vs Gejala: Evaluasi Solusi Autentikasi & Otorisasi",
    subtitle="Mengapa solusi parsial tetap menyisakan celah keamanan kritis pada sistem web modern.",
    content_items=[
        {"title": "TABEL EVALUASI MEKANISME KEAMANAN", "items": [
            "Mekanisme | Peran Primer | Celah Kegagalan / Vektor Serangan | Ketahanan Produksi",
            "OAuth 2.0 Murni | Otorisasi Akses API | Token ditukar sebagai identitas; no identity claims | Rentan (Confused Deputy)",
            "Implicit Flow (Legacy) | Otorisasi SPA Langsung | Token bocor via URL history & referer header | Tidak Aman (RFC 9700 Deprecated)",
            "OIDC ID Token | Autentikasi Pengguna | Skipping aud/iss/nonce verification | Aman jika divalidasi ketat",
            "PKCE (RFC 7636) | Code Exchange Binding | Menggunakan metode 'plain' alih-alih S256 | Aman Mutlak (S256 Mandatory)",
            "Refresh Rotation | Sesi Berkelanjutan | Static token replay tanpa pelacakan family | Aman Mutlak (dengan Family Revocation)"
        ], "type": "table"}
    ],
    narrative="Tabel ini mengevaluasi berbagai mekanisme otorisasi dan autentikasi dalam menghadapi berbagai vektor ancaman dunia nyata. Menggunakan OAuth 2.0 murni tanpa OIDC untuk login pengguna menimbulkan kerentanan fatal karena access token tidak dirancang untuk memverifikasi audiens pemanggil. Implicit flow kini resmi dilarang oleh standar IETF RFC 9700 karena tingginya risiko kebocoran kredensial di URL browser. Hanya kombinasi terpadu antara OIDC, PKCE S256, dan Refresh Token Rotation dengan Family Revocation yang memberikan perlindungan komprehensif.",
    narrative_en="This table evaluates multiple authorization and authentication mechanisms against modern attack vectors. Relying solely on raw OAuth 2.0 for user sign-in creates severe vulnerabilities because access tokens lack caller audience constraints. The legacy Implicit Flow is formally banned by IETF RFC 9700 due to severe token leakage risks via browser history. Only an integrated combination of OIDC, PKCE S256, and Refresh Token Rotation with Token Family Revocation offers comprehensive defense.",
    mental_model="Partial protocols leave gaping security holes; defense-in-depth requires layered standard mechanisms.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md:9-36",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "ARSITEKTUR"],
    title="3 Keputusan Desain Kritis Arsitektur Identitas",
    subtitle="Urutan keputusan arsitektur untuk membangun otorisasi web yang kebal sadap dan replay.",
    content_items=[
        {"title": "1. MANDATORY PKCE S256", "items": [
            "Wajibkan parameter code_challenge dan code_challenge_method=S256 pada endpoint /authorize.",
            "Tolak metode 'plain' untuk semua client publik demi integritas cryptographic hash.",
            "Client wajib mengenerate code_verifier dengan entropi tinggi (32 byte acak / 43-128 karakter)."
        ], "type": "card"},
        {"title": "2. MULTI-STEP ID TOKEN VERIFICATION", "items": [
            "Client wajib memverifikasi cryptographic signature (HMAC/RSA) sebelum memproses claims.",
            "Periksa kecocokan exact match issuer (iss), client ID audience (aud), dan expiration (exp).",
            "Cocokkan parameter nonce dengan nonce acak yang disimpan di sesi klien awal."
        ], "type": "card"},
        {"title": "3. TOKEN ROTATION & FAMILY REVOCATION", "items": [
            "Perlakukan setiap refresh token sebagai single-use credential; terbitkan token baru setiap refresh.",
            "Tautkan token dalam satu silsilah FamilyID untuk mendeteksi penggunaan ulang token kedaluwarsa.",
            "Jika replay terdeteksi, hapus seketika seluruh sesi dalam family tersebut dan paksa re-login."
        ], "type": "card"}
    ],
    narrative="Tiga keputusan arsitektural ini memastikan implementasi OAuth 2.0 dan OIDC memenuhi standar ketat IETF RFC 9700. Pertama, PKCE metode S256 diwajibkan tanpa kompromi untuk seluruh jenis client guna memastikan authorization code tidak dapat dieksploitasi meski berhasil disadap di jaringan publik. Kedua, verifikasi ID Token tidak boleh hanya membaca JSON, melainkan wajib mengevaluasi tanda tangan digital dan parameter nonce anti-replay. Ketiga, rotasi refresh token wajib didukung oleh pelacakan FamilyID untuk mematikan sesi penyerang seketika.",
    narrative_en="These three architectural decisions ensure OAuth 2.0 and OIDC implementations satisfy rigorous IETF RFC 9700 standards. First, PKCE S256 is enforced unconditionally across all clients to ensure intercepted authorization codes cannot be redeemed. Second, ID Token validation must cryptographically verify digital signatures alongside strict issuer, audience, and nonce constraints. Third, refresh token rotation must maintain FamilyID lineage tracking to trigger immediate session revocation upon replay detection.",
    mental_model="Hash the verifier with S256; verify all token claims strictly; treat refresh tokens as single-use consumables.",
    source_path="labs/31-oauth2-and-oidc/content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Siklus Lengkap OIDC, PKCE, & Replay Defense",
    subtitle="Output langsung dari go run ./cmd/demo memperlihatkan keberhasilan flow normal dan mitigasi serangan.",
    content_items=[
        {"title": "DEMO EXECUTION TRACE", "items": [
            "=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===",
            "[Step 1] Initiate Authorization Request with PKCE & OIDC scope...",
            " Generated PKCE code_challenge (S256): dbNVgXttw5HU018kaBwk3JusvwqMVanN4c4ghBsH88w",
            " Generated OIDC nonce: 2152fa0c9254b322d11b27dba96b99de",
            "[Step 2] Authorization Server issues Authorization Code...",
            " Issued Code: 5fabbdcdb6bbd582c3811481395a5161 (expires: 15:31:47)",
            "[Step 3] Client exchanges Code + code_verifier for Tokens...",
            " Access Token : at_4e8dd945eaf2aed814aaf9717a7ee6bab364db0914917bb2",
            " Refresh Token: rt_c624781adb52c6effa11d85bbe24fccf7f73eb20feeaf928",
            " ID Token     : eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...",
            "[Step 4] Validate OIDC ID Token claims... Subject: user_alice_99 | Nonce Match: SUCCESS",
            "[Step 5] Access Resource Server using Access Token... Resource verified. Subject: user_alice_99",
            "[Step 6] Test Interception Attack... Interception Attack Blocked! Error: pkce verification failed",
            "[Step 7] Test Refresh Token Rotation... Rotated New RT: rt_a0aafd72efbf3ffb2f3cfafeacf71cf4...",
            "[Step 8] Test Refresh Token Replay Detection... Replay Detected! Error: entire token family revoked",
            "=== Demo Completed Successfully ==="
        ], "type": "terminal"}
    ],
    narrative="Bukti eksekusi terminal membuktikan bahwa seluruh siklus autentikasi dan otorisasi berjalan mulus tanpa celah keamanan. Langkah 1 hingga 5 menunjukkan pembuatan tantangan PKCE S256, penerbitan authorization code, penukaran token yang sah, verifikasi klaim ID Token, hingga akses ke resource server dengan scope yang sesuai. Langkah 6 hingga 8 memperlihatkan simulasi serangan nyata di mana upaya pertukaran kode oleh penyerang digagalkan oleh PKCE, dan upaya penggunaan ulang refresh token curian langsung memicu pemblokiran seluruh keluarga sesi.",
    narrative_en="The terminal execution trace verifies that the full authentication and authorization lifecycle functions flawlessly. Steps 1 through 5 demonstrate PKCE S256 challenge generation, code issuance, legitimate token exchange, ID Token claim verification, and authorized resource access. Steps 6 through 8 execute active attack simulations where code interception is neutralized by PKCE and stolen refresh token replays immediately trigger complete token family revocation.",
    mental_model="Terminal logs verify the dual guarantee: legitimate traffic succeeds seamlessly; unauthorized replays fail loudly.",
    source_path="cmd/demo/main.go:1-112 • engineering/03-execution-result.md:74-111",
    slide_num="10 / 24",
    layout="terminal-output"
)

# Slide 11: Analisis batasan bukti
add_slide(
    categories=["BATASAN", "ANALISIS"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Batasan ruang lingkup laboratorium yang wajib dipahami sebelum implementasi produksi.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Penyimpanan token dan code bersifat in-memory; tidak menguji persistensi distributed database (Redis/PostgreSQL).",
            "Menggunakan penandatanganan simetris HMAC-SHA256 (HS256); tidak menguji asimetris RS256/ES256 via JWKS endpoint.",
            "Validasi string equality pada PKCE tidak menggunakan constant-time compare (crypto/subtle.ConstantTimeCompare).",
            "Tidak mencakup form login UI/HTML, multi-factor authentication (MFA), atau user consent dialog.",
            "Tidak menguji flow yang telah didepresiasi seperti Implicit Grant atau Resource Owner Password Credentials (ROPC).",
            "Tidak membuktikan mekanisme token-binding tingkat transport seperti DPoP (RFC 9449) atau mTLS (RFC 8705)."
        ], "type": "bullet"}
    ],
    narrative="Klaim teknis dalam laboratorium ini dibatasi secara tegas oleh parameter lingkungan pengujian yang dirancang mandiri (self-contained). Authorization server menyimpan state di dalam memori proses tunggal tanpa integrasi database terdistribusi seperti Redis kluster. Penandatanganan ID Token menggunakan skema simetris HMAC-SHA256 untuk kesederhanaan komputasi lokal, bukan infrastruktur kunci publik asimetris dengan endpoint JWKS dinamis. Pengembang tidak boleh mengabaikan kebutuhan constant-time comparison dan mekanisme sender-constrained token (DPoP) saat beralih ke skala produksi.",
    narrative_en="The technical claims in this laboratory are strictly bounded by its self-contained test design. The authorization server manages state in-process without distributed database backends like Redis clusters. ID Token signing relies on symmetric HMAC-SHA256 rather than asymmetric public-key cryptography backed by dynamic JWKS endpoints. Production deployments must account for constant-time comparisons, database persistence, and advanced sender-constrained mechanisms like DPoP (RFC 9449).",
    mental_model="Acknowledge educational simplifications to ensure production architectures address persistent storage and public JWKS.",
    source_path="labs/31-oauth2-and-oidc/content/01-content-brief.md:37-41 • engineering/02-implementation-notes.md:28-48",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Keamanan Sesi vs Kompleksitas Server",
    subtitle="Menilai biaya penyimpanan state dan rotasi token dibandingkan manfaat eliminasi pembajakan sesi.",
    content_items=[
        {"title": "MANFAAT & KEUNGGULAN SISTEM", "items": [
            "Imunitas Serangan Intersepsi: PKCE S256 menghilangkan ancaman pencurian authorization code di redirect URL.",
            "Deteksi Dini Pembajakan Sesi: Replay token yang dicuri terdeteksi seketika melalui status konsumsi token tunggal.",
            "Decoupled Identity & Scope: ID Token membawa konteks pengguna, access token mengisolasi batasan wewenang API."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI & BIAYA IMPLEMENTASI", "items": [
            "Penyimpanan State Terpusat: Authorization Server wajib melacak histori FamilyID dan status Revoked di storage.",
            "Latensi Rotasi Jaringan: Setiap pembaruan access token mewajibkan round-trip penukaran refresh token ke server.",
            "Potensi False Revocation: Kegagalan koneksi saat rotasi dapat menyebabkan klien sah tidak sengaja memicu deteksi replay."
        ], "type": "bullet"}
    ],
    narrative="Penerapan arsitektur OAuth 2.0 dan OIDC modern menuntut kompromi yang cermat antara tingkat keamanan dan kompleksitas operasional. Menghapus token statis dan memberlakukan Refresh Token Rotation memberikan perlindungan maksimal dari pembajakan sesi jangka panjang, namun mewajibkan server mengelola state silsilah FamilyID secara konsisten di penyimpanan terdistribusi. Selain itu, ada risiko kegagalan jaringan sementara (race condition jaringan) di mana klien yang mencoba refresh ulang dapat disalahartikan sebagai penyerang yang melakukan replay.",
    narrative_en="Adopting modern OAuth 2.0 and OIDC architectures requires a calculated trade-off between security posture and operational complexity. Eliminating static tokens in favor of Refresh Token Rotation guarantees defense against persistent session hijacking, yet forces the server to maintain stateful FamilyID lineage metadata across distributed storage. Furthermore, transient network partitions during rotation exchanges can cause legitimate retries to be mistakenly flagged as replay attacks.",
    mental_model="Single-use tokens eliminate persistent exposure but require robust state tracking and retry handling.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md:262-273",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "VALIDATION"],
    title="Validasi Konkurensi: 10 Goroutine Berebut Refresh Token",
    subtitle="TestOAuth2_ConcurrentRefreshReplay membuktikan: minimal 9 dari 10 request ditolak secara aman.",
    content_items=[
        {"title": "KODE STRESS TEST KONKURENSI", "items": [
            "func TestOAuth2_ConcurrentRefreshReplay(t *testing.T) {",
            "    stolenRT := resp.RefreshToken",
            "    var wg sync.WaitGroup; errCount := 0; var mu sync.Mutex",
            "    for i := 0; i < 10; i++ {",
            "        wg.Add(1)",
            "        go func() {",
            "            defer wg.Done()",
            "            _, err := as.Refresh(stolenRT, \"client_conc_ref\")",
            "            if err != nil { mu.Lock(); errCount++; mu.Unlock() }",
            "        }()",
            "    }",
            "    wg.Wait()",
            "    // Assert: errCount >= 9 akibat rotasi dan family revocation",
            "}"
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI KONKURENSI", "items": [
            "Mutex Protection: Server mengeksekusi method Refresh() di bawah s.mu.Lock() secara linier.",
            "Tepat 1 Menang / Sisanya Ditolak: Maksimal 1 goroutine berhasil merotasi; 9 goroutine lainnya ditolak replay detection.",
            "Race Detector: PASS tanpa deteksi data race condition pada 20 worker simultan (TestOAuth2_ConcurrencyAndRace)."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi pada TestOAuth2_ConcurrentRefreshReplay membuktikan ketangguhan sistem saat satu refresh token yang sama ditembak secara simultan oleh 10 goroutine. Karena method Refresh di Authorization Server dilindungi oleh sync.Mutex eksklusif, evaluasi status token berlangsung secara linearizable. Tepat satu permintaan berhasil merotasi token, sementara sembilan permintaan lainnya seketika memicu error ErrTokenReplayDetected dan mencabut seluruh family. Pengujian go test -race membuktikan ketiadaan data race pada 20 worker konkruen.",
    narrative_en="Concurrency testing in TestOAuth2_ConcurrentRefreshReplay validates system resilience when a single stolen refresh token is replayed simultaneously by 10 goroutines. Because the Authorization Server protects the Refresh method via an exclusive sync.Mutex, state evaluations remain strictly linearizable. At most one request succeeds in rotating the credential, while the remaining nine immediately trigger ErrTokenReplayDetected, invalidating the entire family with zero data races flagged under -race.",
    mental_model="Mutex serialization ensures atomicity: the first refresh consumes the token, subsequent concurrent callers trigger replay defense.",
    source_path="tests/oauth_test.go:365-455",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "PERUBAHAN"],
    title="Before vs After: Evolusi Keamanan Sesi Web",
    subtitle="Transformasi dari arsitektur warisan yang rentan kebocoran menuju standar modern RFC 9700.",
    content_items=[
        {"title": "BEFORE (Legacy Implicit & Static Tokens)", "items": [
            "Alur: Implicit Flow mengekspos access token langsung di URL fragment browser.",
            "Identitas: Klien membedah access token tanpa verifikasi tanda tangan digital.",
            "Refresh Sesi: Refresh token statis tanpa batas pakai disimpan di localStorage.",
            "Dampak Kebocoran: Penyerang yang mencuri token menguasai akun selamanya tanpa terdeteksi."
        ], "type": "card"},
        {"title": "AFTER (Auth Code + PKCE & Token Rotation)", "items": [
            "Alur: Authorization Code Flow wajib dilindungi oleh cryptographic challenge PKCE S256.",
            "Identitas: OIDC ID Token ditandatangani dan divalidasi 7 langkah (iss, aud, exp, nonce).",
            "Refresh Sesi: Refresh token sekali pakai dirotasi setiap request dengan pelacakan FamilyID.",
            "Dampak Kebocoran: Replay token curian memicu pembatalan instan seluruh sesi keluarga token."
        ], "type": "card"}
    ],
    narrative="Perbandingan before-after ini mengilustrasikan lompatan besar standar keamanan web selama satu dekade terakhir. Pada arsitektur lama, pengembang mengandalkan Implicit Flow dan menyimpan refresh token statis di browser localStorage, menciptakan kerentanan masif terhadap pencurian sesi melalui skrip jahat. Pada arsitektur modern berbasis RFC 9700, authorization code diikat secara kriptografis oleh PKCE, ID Token diverifikasi secara matematis, dan refresh token dirotasi pada setiap pemanggilan sehingga penyusupan dapat dideteksi dan dilumpuhkan secara otomatis.",
    narrative_en="This before-and-after comparison highlights the architectural paradigm shift in web session security over the past decade. Legacy systems relied on the Implicit Flow and stored static refresh tokens inside browser localStorage, inviting devastating session hijacking via cross-site scripting. Under modern RFC 9700 architectures, authorization codes are cryptographically bound by PKCE, ID Tokens undergo multi-step signature verification, and refresh tokens rotate continuously with automated replay detection.",
    mental_model="Before: static tokens leak invisibly. After: dynamic rotations turn stolen tokens into automated tripwires.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md:414-434",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 PKCE S256 Challenge & Bounds
add_slide(
    categories=["BUKTI TEST", "PKCE & RFC 7636"],
    title="Test 1: Validasi Kriptografis Tantangan PKCE S256",
    subtitle="TestPKCE_S256_Valid dan TestPKCE_VerifierLength_Bounds membuktikan kepatuhan spesifikasi RFC 7636.",
    content_items=[
        {"title": "KODE PENGUJIAN BOUNDS & METHOD", "items": [
            "func TestPKCE_S256_Valid(t *testing.T) {",
            "    pair, err := pkce.GeneratePKCEPair(\"S256\")",
            "    // Assert: len(pair.CodeVerifier) >= 43 && len(pair.CodeVerifier) <= 128",
            "    err = pkce.Verify(pair.CodeVerifier, pair.CodeChallenge, \"S256\")",
            "    // Assert: err == nil (SHA-256 base64url matching)",
            "}",
            "func TestPKCE_VerifierLength_Bounds(t *testing.T) {",
            "    // Verifier < 43 karakter -> Tolak ErrInvalidVerifierLength",
            "    // Verifier > 128 karakter -> Tolak ErrInvalidVerifierLength",
            "}"
        ], "type": "code"},
        {"title": "VERIFIKASI INVARIAN RFC 7636", "items": [
            "Entropi Kriptografis: 32 byte acak menghasilkan 43 karakter base64url tanpa padding, memenuhi batas 256-bit entropy.",
            "Batas Panjang Ketat: Kode memvalidasi panjang verifier berada dalam rentang minimum 43 dan maksimum 128 karakter.",
            "Penolakan Metode Ilegal: Penggunaan method di luar 'S256' dan 'plain' seketika ditolak dengan ErrInvalidMethod."
        ], "type": "bullet"}
    ],
    narrative="Unit test pertama memverifikasi pemenuhan standar RFC 7636 pada modul pkg/pkce/pkce.go. Fungsi GeneratePKCEPair terbukti menghasilkan code_verifier berkekuatan 256-bit entropy dengan panjang tepat 43 karakter base64url. Pengujian memvalidasi bahwa perhitungan BASE64URL(SHA256(verifier)) cocok secara sempurna dengan challenge yang dikirimkan. Test verifier bounds memastikan bahwa string yang lebih pendek dari 43 karakter atau lebih panjang dari 128 karakter langsung ditolak demi mencegah serangan degradasi entropi.",
    narrative_en="The first unit test suite validates RFC 7636 compliance within pkg/pkce/pkce.go. GeneratePKCEPair demonstrably generates 256-bit cryptographic entropy yielding a 43-character base64url code_verifier. The test verifies that BASE64URL(SHA256(verifier)) matches the issued code_challenge. Furthermore, bounds testing ensures any verifier shorter than 43 characters or exceeding 128 characters is systematically rejected, preventing entropy degradation vulnerabilities.",
    mental_model="High entropy verifiers bound between 43 and 128 chars guarantee unforgeable SHA-256 challenges.",
    source_path="tests/oauth_test.go:17-47 • tests/oauth_test.go:335-346",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 OIDC ID Token Verification
add_slide(
    categories=["BUKTI TEST", "OIDC ID TOKEN"],
    title="Test 2: Verifikasi Kriptografis & Claims ID Token",
    subtitle="Memvalidasi tanda tangan HMAC, modifikasi payload, token expired, dan pencocokan nonce anti-replay.",
    content_items=[
        {"title": "PENGUJIAN KETAHANAN JWT OIDC", "items": [
            "// TestOIDC_IDToken_TamperedSignature: Modifikasi signature",
            "tampered := parts[0] + \".\" + parts[1] + \".invalidSig==\"",
            "_, err = oidc.ParseAndVerifyIDToken(tampered, secret, iss, aud, \"\", now)",
            "// Assert: err == ErrSignatureInvalid",
            "// TestOIDC_IDToken_MismatchClaims: Evaluasi claims",
            "// 1. Wrong issuer -> Tolak ErrIssuerMismatch",
            "// 2. Wrong audience -> Tolak ErrAudienceMismatch",
            "// 3. Wrong nonce -> Tolak ErrNonceMismatch",
            "// 4. Expiration lampau -> Tolak ErrTokenExpired"
        ], "type": "code"},
        {"title": "CHECKPOINT KEAMANAN MULTI-STEP", "items": [
            "Constant-Time Signature Check: HMAC-SHA256 diverifikasi via hmac.Equal untuk menangkal timing attack.",
            "Issuer & Audience Pinning: Token ditolak jika audiens bukan Client ID pemanggil atau diterbitkan pihak asing.",
            "Anti-Replay Nonce: Klaim nonce mengikat ID Token secara unik ke sesi authorization request browser pemohon."
        ], "type": "bullet"}
    ],
    narrative="Pengujian modul OIDC membuktikan bahwa ID Token tidak dapat dimanipulasi oleh pihak ketiga. Pada TestOIDC_IDToken_TamperedSignature, setiap pengubahan satu karakter pun pada tanda tangan JWT langsung digagalkan oleh pemeriksaan hmac.Equal. Sementara itu, rangkaian pengujian klaim memverifikasi bahwa token yang memiliki penerbit (iss) tidak cocok, ditujukan untuk klien lain (aud mismatch), memiliki masa berlaku yang telah lewat (exp), atau memiliki nonce yang berbeda akan ditolak secara mutlak.",
    narrative_en="The OIDC test suite proves that ID Tokens are tamper-proof against unauthorized manipulation. In TestOIDC_IDToken_TamperedSignature, altering a single character of the JWT signature triggers instant rejection via hmac.Equal constant-time verification. Simultaneously, claim tests verify that any token presenting an unrecognized issuer, mismatched audience, expired timestamp, or deviated nonce is intercepted and rejected without exception.",
    mental_model="Every claim is a cryptographic barrier; skipping one invalidates the entire identity assertion.",
    source_path="tests/oauth_test.go:48-143 • pkg/oidc/oidc.go:64-116",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Full Flow & Interception Defense
add_slide(
    categories=["BUKTI TEST", "CODE INTERCEPTION"],
    title="Test 3: Pencegahan Serangan Intersepsi Authorization Code",
    subtitle="TestOAuth2_FullFlowAndPKCEInterception membuktikan pencurian auth code sia-sia tanpa verifier.",
    content_items=[
        {"title": "SIMULASI PENYADAPAN KODE PADA TEST", "items": [
            "// 1. Client sah menginisiasi flow dan meminta authorization code",
            "ch, _ := cli.BuildAuthorizationRequest(\"openid profile\")",
            "ac, _ := as.Authorize(\"client_1\", \"https://app.com/cb\", \"openid profile\", \"user_sub_1\", ch, \"S256\", cli.Nonce)",
            "// 2. Attacker menyadap ac.Code, mencoba exchange dengan verifier tebakan",
            "_, err = as.ExchangeCode(ac.Code, \"client_1\", \"https://app.com/cb\", \"attacker_wrong_verifier_12345678901234567890123\")",
            "// Assert: err != nil (PKCE verification failed: challenge mismatch)",
            "// 3. Client sah menukar code dengan verifier asli -> SUKSES",
            "resp, err := cli.Exchange(ac.Code) // Menerima Access, Refresh, & ID Token",
            "// 4. Percobaan penggunaan ulang ac.Code yang sama -> DITOLAK (Used = true)"
        ], "type": "code"},
        {"title": "DUAL-LAYER CODE SECURITY", "items": [
            "Proof-of-Possession: Penyerang yang menyadap URL redirect tidak dapat menukar code tanpa code_verifier lokal klien.",
            "One-Time Use Enforcement: Flag ac.Used = true mencegah serangan replay authorization code oleh pihak manapun."
        ], "type": "bullet"}
    ],
    narrative="Test ketiga mendemonstrasikan bagaimana proteksi PKCE menggagalkan skenario penyerangan nyata pada client publik. Saat penyerang berhasil menyadap authorization code yang melintasi sistem operasi dan mencoba menukarnya ke server menggunakan verifier palsu, server mendeteksi ketidakcocokan SHA-256 dan membatalkan pertukaran. Klien sah kemudian berhasil menukar kode tersebut dengan verifier aslinya. Ketika kode yang sama dicoba untuk digunakan kembali untuk kedua kalinya, server menolaknya karena kode otorisasi bersifat strictly single-use.",
    narrative_en="The third test orchestrates an end-to-end interception attack simulation on public clients. When an adversary intercepts the authorization code and attempts redemption using a forged code_verifier, the server detects the cryptographic hash mismatch and aborts the exchange. The legitimate client then redeems the code using its authentic verifier. When anyone attempts to replay that redeemed code a second time, the server rejects it because authorization codes are strictly single-use.",
    mental_model="Authorization codes are worthless without local verifier secrets; one-time use prevents post-exchange replay.",
    source_path="tests/oauth_test.go:145-182",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Studi Kasus Industri: RFC 9700 & Deprekasi Implicit Flow
add_slide(
    categories=["STUDI KASUS", "INDUSTRI"],
    title="Studi Kasus Industri: Konsensus Keamanan RFC 9700",
    subtitle="Mengapa IETF secara resmi mendepresiasi Implicit Flow dan mewajibkan PKCE untuk semua aplikasi.",
    content_items=[
        {"title": "DEPREKASI IMPLICIT FLOW (OAUTH 2.1 / RFC 9700)", "items": [
            "Sejarah Masa Lalu: Implicit Flow diciptakan saat browser belum mendukung CORS secara luas.",
            "Vektor Kebocoran: Access token dikembalikan pada URL hash fragment (#access_token=...), bocor ke Access Log, Browser History, dan referer headers.",
            "Vonis IETF 2025: Dihapus sepenuhnya dari spesifikasi; semua SPA wajib beralih ke Auth Code Flow + PKCE."
        ], "type": "bullet"},
        {"title": "BEST CURRENT PRACTICE MODERN (BCP)", "items": [
            "Wajib PKCE untuk Semua: Tidak hanya aplikasi mobile/SPA, server-rendered confidential client juga dianjurkan menggunakan PKCE.",
            "Strict Redirect URI Matching: Larangan penggunaan wildcard URI atau redirect berbasis skema kustom tanpa app-linking.",
            "Sender-Constrained Tokens: Tren industri beralih ke DPoP (RFC 9449) untuk mengikat token akses ke kunci privat browser."
        ], "type": "bullet"}
    ],
    narrative="Riset standar IETF menunjukkan pergeseran paradigma fundamental dalam keamanan web terdistribusi. Implicit Flow yang dahulu menjadi standar bagi Single Page Application resmi dihapus dalam OAuth 2.1 dan RFC 9700 karena token akses rentan bocor melalui riwayat peramban dan skrip analitik pihak ketiga. Standar industri terkini mewajibkan Authorization Code Flow yang dilindungi PKCE S256 untuk seluruh aplikasi. Di masa depan, proteksi ini semakin diperkuat dengan sender-constrained tokens seperti DPoP untuk menghentikan pencurian token di lapisan jaringan.",
    narrative_en="IETF standardization reflects a decisive evolution in distributed web security. The legacy Implicit Flow originally designed for Single Page Applications has been excised from OAuth 2.1 and RFC 9700 because access tokens routinely leaked via browser histories and third-party analytics scripts. Modern best current practices mandate Authorization Code Flow backed by PKCE S256 across all client types, reinforced by emerging sender-constrained standards like DPoP to bind tokens to browser private keys.",
    mental_model="Implicit flows leak credentials in URLs; PKCE and sender-constraining bind tokens to execution contexts.",
    source_path="research/02-sources.md • research/05-report.md",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: Analisis Teknis: Refresh Token Rotation & Family Replay
add_slide(
    categories=["ANALISIS TEKNIS", "SESSION DEFENSE"],
    title="Analisis Teknis: Anatomi Token Family Revocation",
    subtitle="Bagaimana pelacakan silsilah FamilyID mengubah refresh token basi menjadi jebakan alarm penyusup.",
    content_items=[
        {"title": "MEKANISME SILSILAH FAMILYID", "items": [
            "Penerbitan Perdana: Saat authorization code ditukar, server membuat FamilyID baru (fam_...) untuk sesi tersebut.",
            "Rotasi Berkelanjutan: Setiap kali /refresh dipanggil, token aktif ditandai Revoked=true, dan token baru diterbitkan di FamilyID sama.",
            "Integritas Silsilah: Hanya tepat satu token refresh yang berstatus Revoked=false dalam satu FamilyID pada satu waktu."
        ], "type": "bullet"},
        {"title": "TRIGGER REPLAY & REAKSI SISTEM", "items": [
            "Deteksi Pelanggaran: Jika token dengan status Revoked=true masuk kembali, server mendeteksi replay pencurian token.",
            "Pembersihan Total: Server menetapkan revokedFams[FamilyID] = true, membatalkan semua token dalam silsilah tersebut.",
            "Proteksi Pengguna Sah: Klien sah yang mencoba merefresh berikutnya akan ditolak, memutus akses penyerang dan memaksa login baru."
        ], "type": "bullet"}
    ],
    narrative="Mekanisme Token Family Revocation per RFC 9700 Section 4.14 mengubah sifat pasif refresh token menjadi jebakan aktif bagi penyusup. Ketika token yang dicuri oleh penyerang digunakan setelah klien sah melakukan rotasi, server mengenali bahwa token tersebut sudah berstatus terpakai (Revoked=true). Alih-alih hanya menolak satu permintaan tersebut, server secara cerdas mengidentifikasi bahwa seluruh keluarga sesi telah terkompromi. Server langsung mencabut FamilyID secara menyeluruh, memastikan penyerang tidak dapat mempertahankan aksesnya.",
    narrative_en="The Token Family Revocation mechanism specified in RFC 9700 Section 4.14 converts passive refresh tokens into active intrusion alarms. When an adversary attempts to replay a stolen refresh token after the legitimate client has rotated it, the server recognizes the token was already consumed. Rather than merely rejecting the individual request, the server concludes that the session lineage has been breached. It revokes the entire FamilyID instantly, terminating unauthorized access and alerting the system.",
    mental_model="Consumed refresh tokens become tripwires: any subsequent replay instantly incinerates the entire session tree.",
    source_path="pkg/server/server.go:219-286 • tests/oauth_test.go:184-217",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "PRODUKSI"],
    title="Checklist Implementasi Produksi OAuth 2.0 & OIDC",
    subtitle="9 poin verifikasi wajib sebelum menggelar sistem otorisasi dan identitas ke lingkungan produksi.",
    content_items=[
        {"title": "CHECKLIST KEAMANAN PROTOKOL PRODUKSI", "items": [
            "[ ] Parameter code_challenge dan code_challenge_method=S256 wajib ada pada setiap authorization request.",
            "[ ] Code verifier di-generate dengan 32 byte cryptographic random entropy (43-128 base64url characters).",
            "[ ] Authorization code berumur pendek (maks 5 menit) dan strictly single-use (ditandai used setelah exchange).",
            "[ ] ID Token diverifikasi secara kriptografis: signature, iss, aud, exp, iat window, dan nonce match.",
            "[ ] Refresh token dirotasi setiap kali digunakan; token lama seketika ditandai revoked di database.",
            "[ ] Token family tracking aktif: penggunaan ulang token yang sudah revoked membatalkan seluruh silsilah sesi.",
            "[ ] Bearer access token TIDAK PERNAH disimpan di browser localStorage; gunakan httpOnly SameSite cookies (BFF).",
            "[ ] Implicit Grant dan Resource Owner Password Credentials (ROPC) dinonaktifkan sepenuhnya di server.",
            "[ ] Endpoint token dilindungi dari race condition melalui transaksi atomik atau distributed lock."
        ], "type": "bullet"}
    ],
    narrative="Daftar periksa ini merangkum sembilan prinsip krusial yang wajib dipenuhi sebelum menggelar infrastruktur OAuth 2.0 dan OIDC ke sistem produksi perbankan atau enterprise. Seluruh alur harus mewajibkan PKCE S256 tanpa toleransi metode plain, dan masa berlaku authorization code dibatasi maksimal lima menit. Verifikasi ID Token di sisi client wajib mencakup seluruh aspek kriptografi dan claims. Di lapisan penyimpanan browser, penyimpanan token di localStorage harus dihindari sepenuhnya demi memitigasi risiko XSS.",
    narrative_en="This checklist codifies nine critical requirements before deploying OAuth 2.0 and OIDC infrastructure to production enterprise environments. All authorization requests must mandate PKCE S256, and authorization codes must expire within five minutes under strict single-use rules. Client-side ID Token verification must evaluate signatures and all mandatory claims without shortcuts. Storing bearer tokens in localStorage must be eliminated in favor of hardened httpOnly SameSite cookie architectures.",
    mental_model="Enforce S256, verify every token claim, rotate single-use credentials, and quarantine tokens from localStorage.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md:334-345",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Kesalahan Desain yang Kerap Terjadi",
    subtitle="Kekeliruan arsitektur yang membuka celah pembajakan akun dan cara membuktikannya di lab.",
    content_items=[
        {"title": "ANTI-PATTERN UMUM DI INDUSTRI", "items": [
            "Access Token Pseudo-Auth: Memperlakukan access token sebagai identitas tanpa ID Token OIDC.",
            "Skipping Token Claims: Mengabaikan validasi aud atau nonce pada JWT ID Token.",
            "Token di LocalStorage: Menyimpan bearer token di browser Web Storage yang mudah diakses XSS.",
            "Refresh Token Statis: Tidak menerapkan rotasi, membiarkan token curian hidup berminggu-minggu."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI GAGAL DI LAB", "items": [
            "1. Hapus pengecekan `ac.Used` di `server.go` baris 141: auth code dapat diexchange berulang kali!",
            "2. Ubah `meta.Revoked` di `server.go` menjadi no-op: penyerang bebas mereplay refresh token tanpa terdeteksi.",
            "3. Hapus validasi signature di `oidc.go`: test `TestOIDC_IDToken_TamperedSignature` langsung gagal membuktikan celah pemalsuan identitas."
        ], "type": "card"}
    ],
    narrative="Mengabaikan prinsip dasar keamanan identitas web dapat berakibat fatal bagi perlindungan data pengguna. Anti-pattern yang paling sering dijumpai adalah menggunakan access token untuk autentikasi dan menyimpan token berumur panjang di localStorage. Hal ini memberikan jalan bebas bagi serangan XSS untuk membajak akun secara permanen. Pengembang dapat mereproduksi kerentanan ini secara langsung di lab dengan melonggarkan validasi single-use pada server.go dan menyaksikan test suite mendeteksi kegagalan tersebut.",
    narrative_en="Ignoring fundamental identity security principles exposes enterprise user data to catastrophic compromises. The most pervasive anti-patterns include substituting access tokens for authentication and persisting long-lived credentials in browser localStorage, allowing XSS scripts to harvest sessions. Developers can reproduce these failure states locally by relaxing single-use validations in server.go and watching the test suite instantly flag security violations.",
    mental_model="Never substitute access tokens for identity; never store bearer tokens in localStorage; test failure paths proactively.",
    source_path="labs/31-oauth2-and-oidc/content/02-master-draft.md:274-282",
    slide_num="21 / 24",
    layout="anti-pattern-cards"
)

# Slide 22: Takeaway praktis & command block dominan
add_slide(
    categories=["PRAKTIS", "COMMAND"],
    title="Verifikasi Mandiri: Uji & Jalankan Lab",
    subtitle="Reproduksi seluruh 17 unit test, pengujian race detector, dan simulasi demo interaktif di terminal Anda.",
    content_items=[
        {"title": "PERINTAH VERIFIKASI RESMI", "items": [
            "cd labs/31-oauth2-and-oidc",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "RINGKASAN METRIK LAB", "items": [
            "Test Suite: 17/17 test cases lulus 100% tanpa kegagalan (0.317s runtime).",
            "Race Detector: Bersih dari thread contention dan data races pada 20 worker (1.366s pass).",
            "Zero Dependencies: Murni Go Standard Library (crypto/hmac, crypto/sha256, encoding/base64, sync) tanpa vendor eksternal."
        ], "type": "bullet"}
    ],
    narrative="Lab ini dirancang agar dapat diverifikasi secara mandiri dalam hitungan detik oleh siapapun tanpa memerlukan instalasi database eksternal. Seluruh logika kriptografi JWT HS256 dan PKCE dibangun murni menggunakan pustaka standar Go. Perintah go test -v mengeksekusi 17 pengujian komprehensif mulai dari bounds verifier hingga rotasi refresh token. Perintah go test -race memvalidasi integritas state di bawah konkurensi goroutine tinggi, dan go run demo menyajikan visualisasi serangan nyata.",
    narrative_en="This laboratory is engineered for immediate, frictionless peer verification without external database infrastructure. All JWT HS256 signing and PKCE cryptographic logic is implemented using Go standard library primitives. Running go test -v executes 17 comprehensive test suites covering verifier bounds and token rotation. Executing go test -race proves thread safety under multi-worker contention, while go run demo demonstrates real attack mitigations.",
    mental_model="Zero-dependency standard library code enables deterministic, transparent verification of security guarantees.",
    source_path="README.md:1-50 • engineering/03-execution-result.md",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "ARSITEKTUR"],
    title="Prinsip Arsitektur: Desain Identitas Terdistribusi",
    subtitle="Tiga hukum fundamental dalam merancang sistem otorisasi dan identitas yang tahan uji zaman.",
    content_items=[
        {"title": "1. OTORISASI BUKAN AUTENTIKASI", "items": [
            "OAuth 2.0 memberikan izin akses ke resource server via access token dengan scope terbatas.",
            "OIDC memberikan bukti identitas pengguna kepada client via ID Token JWT bertanda tangan."
        ], "type": "card"},
        {"title": "2. BIND KREDENSIAL KE RUNTIME PEMINTA", "items": [
            "Authorization code tidak boleh dipercaya tanpa cryptographic proof-of-possession (PKCE S256).",
            "Sesi login browser diikat secara unik ke ID Token melalui parameter nonce acak."
        ], "type": "card"},
        {"title": "3. ANGGAP TOKEN AKAN BOCOR: SINGLE-USE & ROTATE", "items": [
            "Jangan pernah mempercayai refresh token berumur panjang yang bersifat statis.",
            "Gunakan rotasi token dengan pelacakan family sehingga replay token basi bertindak sebagai alarm pembajakan."
        ], "type": "card"}
    ],
    narrative="Tiga prinsip arsitektural ini melampaui bahasa pemrograman dan framework apapun dalam rekayasa perangkat lunak modern. Pertama, tegakkan batas tegas antara otorisasi dan autentikasi agar sistem tidak terjebak dalam confused deputy attacks. Kedua, jangan pernah mempercayai data yang melintasi browser tanpa pembuktian kepemilikan kriptografis seperti PKCE dan nonce. Ketiga, terapkan paradigma Zero Trust pada sesi pengguna: anggap kredensial dapat bocor sewaktu-waktu, dan gunakan rotasi token keluarga untuk melumpuhkan penyusup seketika.",
    narrative_en="These three architectural principles transcend specific frameworks and languages in modern software engineering. First, enforce an absolute separation between authorization and authentication to eliminate confused deputy vulnerabilities. Second, never trust credentials traversing public browsers without cryptographic proof-of-possession like PKCE and nonces. Third, adopt a Zero Trust mindset for user sessions: assume credentials leak, and leverage family-tracked token rotation as an automated defense mechanism.",
    mental_model="OAuth delegates permission; OIDC proves identity; PKCE binds exchanges; family rotation incinerates compromised sessions.",
    source_path="labs/31-oauth2-and-oidc/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Identitas Terverifikasi, Otorisasi Terlindungi, Sesi Aman.",
    subtitle="",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR LAB", "items": [
            "pkg/pkce/pkce.go — RFC 7636 PKCE S256 verifier & challenge generator and validator",
            "pkg/oidc/oidc.go — OpenID Connect Core 1.0 ID Token JWT signer & 7-step claims verifier",
            "pkg/server/server.go — In-memory AS with PKCE validation, token minting, & family rotation",
            "pkg/client/client.go — Client flow orchestrator for PKCE, token exchange, and refresh",
            "tests/oauth_test.go — 17 comprehensive unit & concurrent race test suites",
            "cmd/demo/main.go — Complete interactive lifecycle and attack mitigation demonstration"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & REPOSITORI", "items": [
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "Unit Tests: 17/17 PASS (100% clean, 0.317s)",
            "Race Detector: PASS (zero data race, 1.366s)",
            f"Repositori: github.com/lukman-ss/software-engineering-lab/tree/main/labs/{LAB_NO}-{SLUG}"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium OAuth 2.0 & OpenID Connect ini membuktikan secara konkret bahwa keamanan identitas dan otorisasi modern dapat dibangun secara elegan tanpa kompromi performa. Melalui kombinasi Authorization Code Flow yang diproteksi PKCE S256, verifikasi kriptografis ID Token bertingkat, serta rotasi refresh token dengan deteksi replay keluarga, arsitektur ini kebal terhadap intersepsi kode dan pembajakan sesi. Seluruh kode sumber, dokumentasi arsitektur, dan instruksi verifikasi tersedia lengkap di repositori.",
    narrative_en="This OAuth 2.0 & OpenID Connect laboratory provides definitive proof that modern identity and authorization security can be achieved cleanly without sacrificing performance. By pairing Authorization Code Flow fortified with PKCE S256, multi-step cryptographic ID Token verification, and family-tracked refresh token rotation, the architecture remains fully resilient against code interception and session hijacking. Complete source code and specifications are published in the repository.",
    mental_model="OAuth for permissions; OIDC for identity; PKCE for proof; Family Rotation for session safety.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
