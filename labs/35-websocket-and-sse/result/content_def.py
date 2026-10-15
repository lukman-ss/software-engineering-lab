#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 35: WebSocket and Server-Sent Events (SSE).
Defines metadata, colors, and 24 slides following Lab 14 & 15 standards.
"""
import os
import sys

OUTPUT_DIR = "labs/35-websocket-and-sse/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/35-websocket-and-sse"
COMMIT_SHA = "b48fe6c30d0b3635f36b8549c5621bbf5ce7c5f8"
LAB_NO = "35"
SLUG = "websocket-and-sse"
TITLE_ID = "WebSocket vs Server-Sent Events"
SUBTITLE_ID = "Memilih Protokol Real-Time Berdasarkan Directionality dan Kompleksitas"
TITLE_EN = "WebSocket vs Server-Sent Events"
SUBTITLE_EN = "Selecting Real-Time Protocols via Directionality and Complexity"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "WebSocket • Server-Sent Events • RFC 6455 • WHATWG SSE • Go 1.22"

DEVTO_DESC = "Explore WebSocket vs SSE in Go: RFC 6455 full-duplex binary framing vs WHATWG text/event-stream push with Last-Event-ID resumption."
DEVTO_TAGS = "go, webdev, performance, architecture"
DEVTO_INTRO = "Selecting between WebSocket and Server-Sent Events dictates connection complexity, proxy resilience, and operational overhead. Lab 35 implements both protocols from scratch in pure Go to expose their architectural trade-offs."

MEDIUM_DESC = "Panduan arsitektur komparatif WebSocket (RFC 6455) vs Server-Sent Events (WHATWG SSE) menggunakan Go murni tanpa library eksternal."

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
            "Framing RFC 6455 vs text/event-stream HTTP streaming",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Directionality dictates the protocol: server-push belongs to SSE, bidirectional belongs to WebSocket.",
    source_path="labs/35-websocket-and-sse/README.md",
    slide_num="01 / 24",
    narrative="",
    title_en=TITLE_EN,
    narrative_en="Real-time protocols must be chosen based on communication directionality rather than trend."
)

# Slide 2: Dilema Protokol Real-Time
add_slide(
    categories=["MOTIVASI", "LATAR BELAKANG"],
    title="Dilema Protokol Real-Time Modern",
    subtitle="Kesalahan pemilihan protokol memicu beban operasional dan kompleksitas stateful",
    content_items=[
        {"title": "GEJALA UMUM", "items": [
            "Fitur notifikasi dan live dashboard sering otomatis dipasangi WebSocket",
            "WebSocket memaksakan upgrade TCP, heartbeat kustom, dan reconnection manual",
            "Proxy reverse dan load balancer memerlukan sticky routing yang rumit",
        ], "type": "bullet"},
        {"title": "BIAYA OPERASIONAL", "items": [
            "Stateful connection membebani memori saat client hanya membaca data pasif",
            "Kegagalan reconnection storm melumpuhkan gateway saat restart",
            "Debugging framing biner lebih sulit dibanding stream teks standar",
        ], "type": "bullet"},
        {"title": "PERTANYAAN UTAMA", "items": [
            "Kapan arsitektur benar-benar memerlukan pipa dua arah (bidirectional)?",
            "Bagaimana memanfaatkan streaming HTTP standar tanpa overhead WebSocket?",
        ], "type": "bullet"},
    ],
    mental_model="Do not pay bidirectional connection tax for unidirectional server push.",
    source_path="labs/35-websocket-and-sse/README.md",
    slide_num="02 / 24",
    narrative="Banyak tim rekayasa perangkat lunak menganggap WebSocket sebagai standar de facto untuk semua kebutuhan real-time. Padahal, memasang WebSocket untuk use case satu arah seperti notifikasi atau live dashboard memperkenalkan overhead operasional yang tidak perlu. SSE menyediakan mekanisme server push berbasis HTTP murni yang jauh lebih ramah proxy dan hemat resource.",
    title_en="The Modern Real-Time Protocol Dilemma",
    narrative_en="Engineering teams often treat WebSocket as the default choice for all real-time requirements. Adopting WebSocket for strictly unidirectional use cases introduces unnecessary operational overhead. Server-Sent Events offers standard HTTP push that is lightweight and proxy-friendly."
)

# Slide 3: Mental Model Inti: Directionality Rule
add_slide(
    categories=["MENTAL MODEL", "ATURAN KEPUTUSAN"],
    title="Mental Model: Directionality Rule",
    subtitle="Arah aliran pesan menentukan protokol, bukan sekadar latensi atau tren",
    content_items=[
        {"title": "SERVER-PUSH ONLY", "items": [
            "Arah komunikasi: Server -> Client secara sepihak",
            "Contoh: Notifikasi sistem, live stock ticker, streaming token LLM",
            "Solusi tepat: Server-Sent Events (SSE) berbasis HTTP standar",
        ], "type": "bullet"},
        {"title": "FULL-DUPLEX BIDIRECTIONAL", "items": [
            "Arah komunikasi: Server <-> Client secara independen dan simultan",
            "Contoh: Chat interaktif, collaborative canvas, multiplayer game",
            "Solusi tepat: WebSocket (RFC 6455) atas koneksi raw TCP",
        ], "type": "bullet"},
        {"title": "PEDOMAN ARSITEKTUR", "items": [
            "Jika client tidak pernah mengirim pesan di channel yang sama: pilih SSE",
            "Jika client harus push data dengan latensi sub-milidetik: pilih WebSocket",
        ], "type": "bullet"},
    ],
    mental_model="Directionality dictates the protocol: server-push belongs to SSE, bidirectional belongs to WebSocket.",
    source_path="labs/35-websocket-and-sse/content/01-content-brief.md",
    slide_num="03 / 24",
    narrative="Keputusan arsitektural dapat disederhanakan melalui satu pertanyaan fundamental: apakah client perlu mengirim pesan ke server secara independen pada saluran yang sama? Jika komunikasi hanya satu arah dari server ke client, SSE memberikan solusi paling bersih. WebSocket hanya diwajibkan ketika komunikasi dua arah berlatensi rendah mutlak diperlukan.",
    title_en="Core Mental Model: Directionality Rule",
    narrative_en="Architectural selection reduces to one question: does the client send messages upstream independently over the same channel? If communication flows strictly server-to-client, SSE is optimal. WebSocket is reserved for genuine full-duplex requirements."
)

# Slide 4: Arsitektur Komparatif Lab 35
add_slide(
    categories=["ARSITEKTUR", "TOPOLOGI"],
    title="Topologi Arsitektur Lab 35",
    subtitle="Unified HTTP Router di Go menghosting /sse, /ws, dan /publish secara murni",
    content_items=[
        {"title": "ENDPOINT DISPATCH", "items": [
            "GET /sse      -> SSE Hub (Handler streaming text/event-stream)",
            "GET /ws       -> WebSocket Upgrader (RFC 6455 Hijack + loop)",
            "POST /publish -> Trigger broadcast event ke seluruh client SSE aktif",
        ], "type": "bullet"},
        {"title": "ZERO DEPENDENCIES", "items": [
            "Implementasi 100% Go Standard Library (net/http, sync, crypto)",
            "Tanpa gorilla/websocket atau pustaka eksternal pihak ketiga",
            "Membuktikan mekanisme low-level framing dan HTTP connection hijacking",
        ], "type": "bullet"},
        {"title": "ISOLASI KOMPONEN", "items": [
            "internal/sse: Hub broadcast dengan buffer history dan Last-Event-ID",
            "internal/ws: RFC 6455 parser, unmasking XOR, frame text & binary",
        ], "type": "bullet"},
    ],
    mental_model="Unified standard library architecture demonstrates both protocols side by side without third-party bias.",
    source_path="labs/35-websocket-and-sse/internal/server/server.go",
    slide_num="04 / 24",
    narrative="Lab 35 mendemonstrasikan kedua protokol di bawah satu HTTP ServeMux di internal/server/server.go. Pendekatan zero external dependencies membongkar bagaimana Go standard library mampu menangani streaming HTTP (SSE) melalui http.Flusher sekaligus protokol non-HTTP (WebSocket) melalui interface http.Hijacker.",
    title_en="Lab 35 Comparative Architecture",
    narrative_en="Lab 35 hosts both protocols under a single Go HTTP ServeMux. The zero-dependency implementation reveals how standard library interfaces handle persistent HTTP streaming via http.Flusher and raw protocol transitions via http.Hijacker."
)

# Slide 5: Anatomi Protokol SSE (WHATWG)
add_slide(
    categories=["SPESIFIKASI", "WHATWG SSE"],
    title="Anatomi Server-Sent Events (SSE)",
    subtitle="Streaming teks persisten berbasis HTTP dengan MIME type text/event-stream",
    content_items=[
        {"title": "HEADER STANDAR HTTP", "items": [
            "Content-Type: text/event-stream (wajib untuk browser EventSource)",
            "Cache-Control: no-cache (mencegah caching di perantara)",
            "Connection: keep-alive (menjaga socket TCP tetap terbuka)",
            "X-Accel-Buffering: no (menonaktifkan buffering di Nginx proxy)",
        ], "type": "bullet"},
        {"title": "STRUKTUR FRAME EVENT", "items": [
            "id: <nomor>        -> Identifier unik untuk resumption",
            "event: <nama>      -> Tipe custom event untuk event listener",
            "retry: <milidetik> -> Durasi tunggu reconnect client browser",
            "data: <konten>     -> Payload teks (dapat multi-line dengan data: prefix)",
        ], "type": "bullet"},
        {"title": "DELIMITER BLOK", "items": [
            "Setiap baris diakhiri \\n tunggal",
            "Akhir dari satu blok event ditandai dengan \\n\\n ganda",
        ], "type": "bullet"},
    ],
    mental_model="SSE is plain HTTP with structured text frames delimited by double newlines.",
    source_path="labs/35-websocket-and-sse/internal/sse/sse.go:19",
    slide_num="05 / 24",
    narrative="Spesifikasi WHATWG mendefinisikan SSE sebagai respons HTTP berdurasi panjang dengan format text/event-stream. Setiap event terdiri dari field id, event, retry, dan data yang diakhiri oleh double newline. Karena berbasis HTTP murni, SSE langsung kompatibel dengan infrastruktur web yang ada seperti TLS termination dan monitoring.",
    title_en="Anatomy of Server-Sent Events (SSE)",
    narrative_en="The WHATWG specification defines SSE as a long-lived HTTP response using text/event-stream. Each event comprises id, event, retry, and data fields terminated by double newlines. Operating over pure HTTP ensures immediate compatibility with web infrastructure."
)

# Slide 6: SSE Event Formatting di Go
add_slide(
    categories=["BUKTI KODE", "SSE FORMATTING"],
    title="Format Event SSE: Multi-Line Framing",
    subtitle="Implementasi method Event.Format() memastikan kepatuhan spesifikasi WHATWG",
    content_items=[
        {"title": "KODE SSE FORMAT", "items": [
            "func (e Event) Format() []byte {",
            "    var sb strings.Builder",
            "    if e.ID > 0 { fmt.Fprintf(&sb, \"id: %d\\n\", e.ID) }",
            "    if e.Event != \"\" { fmt.Fprintf(&sb, \"event: %s\\n\", e.Event) }",
            "    if e.Retry > 0 { fmt.Fprintf(&sb, \"retry: %d\\n\", e.Retry) }",
            "    lines := strings.Split(e.Data, \"\\n\")",
            "    for _, l := range lines {",
            "        fmt.Fprintf(&sb, \"data: %s\\n\", l)",
            "    }",
            "    sb.WriteString(\"\\n\")",
            "    return []byte(sb.String())",
            "}",
        ], "type": "code"},
        {"title": "ATURAN MULTI-LINE", "items": [
            "Payload data yang memiliki newline harus dipecah menjadi beberapa baris data:",
            "Browser EventSource akan menggabungkan kembali baris-baris tersebut dengan newline",
            "Trailing newline kedua menutup frame dan memicu event dispatch di client",
        ], "type": "bullet"},
    ],
    mental_model="Split multi-line data payloads across multiple data prefixes before terminating with double newline.",
    source_path="labs/35-websocket-and-sse/internal/sse/sse.go:19-36",
    slide_num="06 / 24",
    narrative="Method Format() pada struct Event mematuhi aturan WHATWG secara ketat. Jika payload mengandung baris baru, setiap baris harus diawali prefiks 'data: '. Kegagalan memecah baris baru akan merusak parsing client EventSource di browser. Karakter double newline di akhir menjadi sinyal bahwa frame event telah lengkap.",
    title_en="SSE Event Formatting in Go",
    narrative_en="The Format() method enforces WHATWG rules. When payloads contain newlines, each line must carry a 'data: ' prefix. Client parsers assemble these lines sequentially. The final double newline marks the completion of the event frame."
)

# Slide 7: SSE Hub & Last-Event-ID Resumption
add_slide(
    categories=["LOGIKA SISTEM", "RECOVERY"],
    title="SSE Hub: Resumption via Last-Event-ID",
    subtitle="State management memutar ulang event tertinggal saat client melakukan reconnect",
    content_items=[
        {"title": "MEKANISME RESUMPTION", "items": [
            "Client terputus -> browser otomatis reconnect mengirim header Last-Event-ID",
            "Hub membaca Last-Event-ID dan memfilter history: evt.ID > lastID",
            "Event yang terlewat dikirimkan seketika sebelum mendengarkan stream live",
        ], "type": "bullet"},
        {"title": "KODE SUBSCRIBE & REPLAY", "items": [
            "func (h *Hub) Subscribe(lastID int) (chan Event, []Event, func()) {",
            "    h.mu.Lock(); defer h.mu.Unlock()",
            "    var replay []Event",
            "    for _, evt := range h.history {",
            "        if evt.ID > lastID { replay = append(replay, evt) }",
            "    }",
            "    ch := make(chan Event, 64)",
            "    h.clients[ch] = struct{}{}",
            "    return ch, replay, func() { ... }",
            "}",
        ], "type": "code"},
        {"title": "JAMINAN KONSISTENSI", "items": [
            "Client tidak kehilangan update penting akibat network flap sementara",
            "Proses terjadi transparan tanpa perlu logika retry kustom di aplikasi",
        ], "type": "bullet"},
    ],
    mental_model="Last-Event-ID is a built-in recovery offset: replay what was missed, then stream live.",
    source_path="labs/35-websocket-and-sse/internal/sse/sse.go:74-96",
    slide_num="07 / 24",
    narrative="Fitur paling unggul dari SSE adalah kemampuan auto-reconnect bawaan yang dipadukan dengan header Last-Event-ID. Ketika koneksi putus, browser secara otomatis mengirimkan ID terakhir yang diterimanya. Hub di backend menyaring riwayat in-memory dan langsung memutar ulang event yang terlewat sebelum menyambungkan client ke antrean live.",
    title_en="SSE Hub: Resumption via Last-Event-ID",
    narrative_en="SSE's defining advantage is native browser reconnection coupled with the Last-Event-ID header. Upon disconnection, the browser retransmits its latest ID. The hub replays missed events from buffer memory before piping real-time events."
)

# Slide 8: HTTP Streaming & Keep-Alive Loop
add_slide(
    categories=["BUKTI KODE", "CONCURRENCY"],
    title="Loop Streaming HTTP & Heartbeat",
    subtitle="Flush berkala dan komentar keep-alive mencegah timeout reverse proxy",
    content_items=[
        {"title": "KODE SERVEHTTP SSE", "items": [
            "flusher, ok := w.(http.Flusher)",
            "for {",
            "    select {",
            "    case <-ctx.Done(): return",
            "    case evt, ok := <-ch:",
            "        if !ok { return }",
            "        w.Write(evt.Format())",
            "        flusher.Flush()",
            "    case <-ticker.C:",
            "        w.Write([]byte(\": keep-alive\\n\\n\"))",
            "        flusher.Flush()",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "PERAN FLUSHER", "items": [
            "http.Flusher mendorong buffer HTTP kernel seketika tanpa menunggu penuh",
            "Mencegah penumpukan paket di TCP socket buffer",
        ], "type": "bullet"},
        {"title": "KOMENTAR KEEP-ALIVE", "items": [
            "Prefix titik dua (':') diabaikan oleh parser EventSource browser",
            "Mengirim data setiap 15 detik untuk menjaga state NAT/Firewall proxy tetap open",
        ], "type": "bullet"},
    ],
    mental_model="Flushing pushes frames instantly; comment frames defeat proxy idle timeouts.",
    source_path="labs/35-websocket-and-sse/internal/sse/sse.go:98-151",
    slide_num="08 / 24",
    narrative="Koneksi SSE memerlukan http.Flusher agar frame teks terkirim ke client seketika tanpa tertahan di buffer Go runtime. Loop select memantau tiga hal: sinyal pemutusan dari context request, event baru dari broadcast channel, dan ticker 15 detik yang mengirim komentar ': keep-alive' untuk mencegah proxy memutus koneksi akibat idle timeout.",
    title_en="HTTP Streaming Loop & Keep-Alive",
    narrative_en="SSE streaming requires http.Flusher to push frames through Go runtime buffers immediately. The select loop handles client termination contexts, live channel broadcasts, and a 15-second ticker sending ': keep-alive' comments to prevent intermediate proxy timeouts."
)

# Slide 9: Anatomi Protokol WebSocket (RFC 6455)
add_slide(
    categories=["SPESIFIKASI", "RFC 6455"],
    title="Anatomi Protokol WebSocket (RFC 6455)",
    subtitle="Protokol non-HTTP full-duplex melalui upgrade TCP connection tunggal",
    content_items=[
        {"title": "HANDSHAKE HTTP/1.1", "items": [
            "Request: GET /ws dengan Upgrade: websocket & Connection: Upgrade",
            "Sec-WebSocket-Key: Nilai random 16-byte yang di-encode base64",
            "Respons: HTTP/1.1 101 Switching Protocols",
            "Sec-WebSocket-Accept: base64(SHA-1(Key + WS_GUID))",
        ], "type": "bullet"},
        {"title": "KONEKSI NON-HTTP", "items": [
            "Setelah status 101, socket dilepas dari HTTP engine (Hijacked)",
            "Komunikasi beralih ke framing biner raw TCP berarah ganda",
            "Server dan client dapat mengirim frame kapan saja secara simultan",
        ], "type": "bullet"},
        {"title": "TIPE PAYLOAD", "items": [
            "Mendukung teks UTF-8 murni (Opcode 0x1)",
            "Mendukung raw binary byte array (Opcode 0x2)",
            "Control frames: Close (0x8), Ping (0x9), Pong (0xA)",
        ], "type": "bullet"},
    ],
    mental_model="WebSocket begins as HTTP handshake then metamorphoses into raw TCP binary framing.",
    source_path="labs/35-websocket-and-sse/internal/ws/ws.go:31-76",
    slide_num="09 / 24",
    narrative="WebSocket dirancang oleh IETF RFC 6455 untuk menghadirkan kanal dua arah berlatensi rendah. Koneksi diawali dengan request HTTP biasa yang meminta Upgrade. Server merespons dengan 101 Switching Protocols dan menghitung Sec-WebSocket-Accept menggunakan GUID standar. Setelah itu, socket sepenuhnya beralih ke parsing frame biner.",
    title_en="Anatomy of WebSocket Protocol (RFC 6455)",
    narrative_en="WebSocket was standardized by IETF RFC 6455 to provide low-latency full-duplex communication over single TCP connections. Following an HTTP 101 Switching Protocols handshake, the underlying socket shifts entirely to raw binary framing."
)

# Slide 10: WebSocket Upgrade & Hijacking di Go
add_slide(
    categories=["BUKTI KODE", "HIJACKING"],
    title="Handshake Upgrade & Connection Hijacking",
    subtitle="Mengambil alih raw socket TCP dari HTTP Server menggunakan http.Hijacker",
    content_items=[
        {"title": "KODE UPGRADE GO", "items": [
            "h := sha1.New()",
            "h.Write([]byte(key + \"258EAFA5-E914-47DA-95CA-C5AB0DC85B11\"))",
            "accept := base64.StdEncoding.EncodeToString(h.Sum(nil))",
            "",
            "hj, ok := w.(http.Hijacker)",
            "netConn, bufrw, err := hj.Hijack()",
            "resp := \"HTTP/1.1 101 Switching Protocols\\r\\n\" +",
            "    \"Upgrade: websocket\\r\\nConnection: Upgrade\\r\\n\" +",
            "    \"Sec-WebSocket-Accept: \" + accept + \"\\r\\n\\r\\n\"",
            "bufrw.WriteString(resp); bufrw.Flush()",
        ], "type": "code"},
        {"title": "IMPLIKASI HIJACK", "items": [
            "net/http tidak lagi mengelola timeout atau penutupan socket",
            "Aplikasi bertanggung jawab penuh atas buffer I/O dan life cycle koneksi",
        ], "type": "bullet"},
    ],
    mental_model="Hijacking frees the raw TCP socket from net/http lifecycle constraints.",
    source_path="labs/35-websocket-and-sse/internal/ws/ws.go:43-76",
    slide_num="10 / 24",
    narrative="Fungsi Upgrade() mendemonstrasikan kalkulasi kriptografis Sec-WebSocket-Accept sesuai RFC 6455. Dengan memanggil w.(http.Hijacker).Hijack(), Go melepaskan kendali socket TCP dari net/http runtime. Server menulis respons 101 Switching Protocols secara langsung ke buffer TCP dan membungkus socket ke dalam struct Conn kustom.",
    title_en="WebSocket Upgrade & Connection Hijacking",
    narrative_en="The Upgrade() function computes the Sec-WebSocket-Accept challenge. By invoking http.Hijacker, Go transfers socket ownership from net/http. The server writes 101 Switching Protocols headers directly to TCP buffers and wraps the connection into a custom struct."
)

# Slide 11: Struktur Frame Biner RFC 6455
add_slide(
    categories=["SPESIFIKASI", "WIRE FRAMING"],
    title="Struktur Frame Biner RFC 6455",
    subtitle="Encoding header hemat byte dengan mask flag, opcode, dan multi-tier length",
    content_items=[
        {"title": "LAYOUT FRAME RFC 6455", "items": [
            "Byte 0: FIN (1 bit) + RSV (3 bits) + Opcode (4 bits)",
            "Byte 1: MASK (1 bit) + Payload Length (7 bits)",
            "Extended Length (opsional): 2 bytes (jika len=126) atau 8 bytes (jika len=127)",
            "Masking Key: 4 bytes (wajib ada untuk frame dari client ke server)",
            "Payload Data: N bytes (dimasking dengan XOR jika MASK=1)",
        ], "type": "diagram"},
        {"title": "ENCODING PANJANG PAYLOAD", "items": [
            "<= 125 bytes     -> Cukup 7 bit pada Byte 1",
            "126 - 65,535     -> 7 bit bernilai 126 + 2 byte uint16 BigEndian",
            "> 65,535 bytes   -> 7 bit bernilai 127 + 8 byte uint64 BigEndian",
        ], "type": "bullet"},
        {"title": "ATURAN MASKING", "items": [
            "Client -> Server: Wajib dimasking (mencegah cache poisoning perantara)",
            "Server -> Client: Dilarang dimasking (RFC 6455 Section 5.1)",
        ], "type": "bullet"},
    ],
    mental_model="Bit-level efficiency: 2-byte minimum header accommodates small frames to gigabytes.",
    source_path="labs/35-websocket-and-sse/internal/ws/ws.go:84-174",
    slide_num="11 / 24",
    narrative="Struktur frame RFC 6455 dirancang sangat efisien. Header terkecil hanya membutuhkan 2 byte. Panjang payload dikodekan secara adaptif: 7 bit untuk pesan kecil, 16 bit untuk ukuran menengah, dan 64 bit untuk data raksasa. Mekanisme masking wajib diterapkan oleh client untuk mencegah eksploitasi cache poisoning pada proxy perantara.",
    title_en="RFC 6455 Binary Frame Structure",
    narrative_en="RFC 6455 framing minimizes wire overhead with a 2-byte base header. Payload lengths encode dynamically across 7-bit, 16-bit, and 64-bit boundaries. Client-to-server masking protects intermediaries against cache poisoning vulnerabilities."
)

# Slide 12: ReadFrame & XOR Unmasking
add_slide(
    categories=["BUKTI KODE", "PARSING WIRE"],
    title="Parsing Frame & Operasi XOR Masking",
    subtitle="Membaca byte header dan memulihkan payload asli client secara bitwise",
    content_items=[
        {"title": "KODE READFRAME", "items": [
            "var head [2]byte",
            "io.ReadFull(c.rw, head[:])",
            "opcode = head[0] & 0x0F",
            "masked := (head[1] & 0x80) != 0",
            "length := uint64(head[1] & 0x7F)",
            "// ... handle extended length ...",
            "if masked {",
            "    io.ReadFull(c.rw, mask[:])",
            "    io.ReadFull(c.rw, data)",
            "    for i := range data {",
            "        data[i] ^= mask[i%4] // Bitwise XOR unmasking",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "SIFAT SIMETRIS XOR", "items": [
            "Operasi (A ^ B) ^ B = A",
            "Proses unmasking di server identik dengan proses masking di client",
            "Biaya CPU rendah namun membutuhkan iterasi byte-by-byte",
        ], "type": "bullet"},
    ],
    mental_model="XOR unmasking restores wire bytes in-place: data[i] ^= mask[i % 4].",
    source_path="labs/35-websocket-and-sse/internal/ws/ws.go:84-127",
    slide_num="12 / 24",
    narrative="Method ReadFrame() mengekstrak header biner dari stream TCP. Jika bit MASK aktif, server membaca 4 byte masking key. Setiap byte payload kemudian di-XOR dengan mask[i % 4] untuk membuka enkripsi transport. Sifat simetris XOR memungkinkan decoding langsung di memori secara cepat tanpa alokasi buffer sekunder.",
    title_en="Frame Parsing & XOR Unmasking",
    narrative_en="The ReadFrame() method decodes wire bytes from the stream. When the mask bit is set, the server reads the 4-byte key and performs in-place XOR restoration. This guarantees protocol fidelity with minimal memory allocation."
)

# Slide 13: WriteFrame: Framing Text & Binary
add_slide(
    categories=["BUKTI KODE", "SERIALISASI"],
    title="WriteFrame: Konstruksi Frame Outbound",
    subtitle="Serialisasi data text dan biner dengan pengaturan bit FIN dan masking",
    content_items=[
        {"title": "KODE WRITEFRAME", "items": [
            "func (c *Conn) WriteFrame(opcode byte, payload []byte, mask bool) error {",
            "    c.mu.Lock(); defer c.mu.Unlock()",
            "    var buf []byte",
            "    b0 := byte(0x80) | (opcode & 0x0F) // FIN=1",
            "    buf = append(buf, b0)",
            "    // ... encode length (7, 16, or 64 bit) ...",
            "    if mask {",
            "        // append maskKey and maskedPayload",
            "    } else {",
            "        buf = append(buf, payload...)",
            "    }",
            "    _, err := c.rw.Write(buf)",
            "    return err",
            "}",
        ], "type": "code"},
        {"title": "THREAD-SAFETY SOCKET", "items": [
            "Mutex memastikan frame tidak terfragmentasi saat ditulis bersamaan",
            "Bit FIN=1 menandai frame lengkap (tidak terpecah dalam continuation)",
        ], "type": "bullet"},
    ],
    mental_model="Mutex synchronization guarantees atomic frame delivery over raw TCP sockets.",
    source_path="labs/35-websocket-and-sse/internal/ws/ws.go:129-174",
    slide_num="13 / 24",
    narrative="Method WriteFrame() menyusun paket biner untuk dikirim ke socket. Mutex c.mu menjamin bahwa penulisan frame bersifat atomic sehingga dua goroutine tidak mencampur aduk byte frame di TCP stream. Bit FIN diset ke 1 untuk menandakan pesan mandiri tanpa fragmentasi lanjutan.",
    title_en="WriteFrame: Outbound Frame Construction",
    narrative_en="The WriteFrame() method serializes outbound binary frames. Mutex synchronization guarantees that frame writes remain atomic across concurrent goroutines, avoiding frame byte interleaving on TCP streams."
)

# Slide 14: Echo Server & Payload Reversal
add_slide(
    categories=["LOGIKA SISTEM", "VERIFIKASI INTEGRITAS"],
    title="Server Loop: Echo Teks & Reversal Biner",
    subtitle="Menangani OpClose, OpPing, OpText, dan OpBinary pada endpoint /ws",
    content_items=[
        {"title": "KODE HANDLER LOOP /WS", "items": [
            "switch op {",
            "case ws.OpClose:",
            "    conn.WriteFrame(ws.OpClose, nil, false); return",
            "case ws.OpPing:",
            "    conn.WriteFrame(ws.OpPong, data, false)",
            "case ws.OpText:",
            "    conn.WriteFrame(ws.OpText, append([]byte(\"echo: \"), data...), false)",
            "case ws.OpBinary:",
            "    rev := make([]byte, len(data))",
            "    for i := range data { rev[i] = data[len(data)-1-i] }",
            "    conn.WriteFrame(ws.OpBinary, rev, false)",
            "}",
        ], "type": "code"},
        {"title": "PEMBUKTIAN BINER", "items": [
            "Client mengirim 0xDEADBEEF -> Server membalikkan menjadi 0xEFBEADDE",
            "Membuktikan integritas byte murni tanpa distorsi encoding karakter",
        ], "type": "bullet"},
    ],
    mental_model="Byte reversal conclusively verifies true binary transport fidelity.",
    source_path="labs/35-websocket-and-sse/internal/server/server.go:34-56",
    slide_num="14 / 24",
    narrative="Loop pemrosesan WebSocket di server menangani berbagai tipe frame kontrol dan data. Untuk teks, server menambahkan prefiks 'echo: '. Untuk biner, payload dibalik urutan byte-nya sebelum dikirim kembali. Reversal ini memberikan bukti deterministik bahwa protokol mampu mentransfer biner murni tanpa kompromi.",
    title_en="Server Loop: Text Echo & Binary Reversal",
    narrative_en="The WebSocket event loop routes incoming frames by opcode. Text frames receive an echo prefix, while binary payloads undergo byte-order reversal. This reversal proves transparent raw binary transfer without encoding loss."
)

# Slide 15: Hasil Eksekusi Demo & Verifikasi
add_slide(
    categories=["HASIL PENGUJIAN", "DEMO RUN"],
    title="Eksekusi Demo: Bukti Terminal Langsung",
    subtitle="Output deterministik dari cmd/demo memvalidasi resumption dan duplexing",
    content_items=[
        {"title": "OUTPUT DEMO SSE", "items": [
            "> 1. SSE: Unidirectional Server Push with Last-Event-ID",
            "Client connected with Last-Event-ID: 1. Received resumed events:",
            "  [SSE Stream] id: 2",
            "  [SSE Stream] event: news",
            "  [SSE Stream] retry: 2000",
            "  [SSE Stream] data: Update: Faster routing in net/http",
        ], "type": "code"},
        {"title": "OUTPUT DEMO WEBSOCKET", "items": [
            "> 2. WebSocket: Full-Duplex Text and Binary Protocol",
            "Client sending WS text frame: \"Hello RFC 6455\"",
            "Server responded with Opcode 0x1 (Text): \"echo: Hello RFC 6455\"",
            "Client sending WS binary frame: DEADBEEF",
            "Server responded with Opcode 0x2 (Binary): EFBEADDE",
            "Binary frame payload verified successfully.",
        ], "type": "code"},
    ],
    mental_model="Deterministic verification: SSE resumes exactly from offset; WebSocket mirrors bytes exactly.",
    source_path="labs/35-websocket-and-sse/cmd/demo/main.go",
    slide_num="15 / 24",
    narrative="Eksekusi cmd/demo membuktikan klaim teknis secara nyata. Pada bagian SSE, pengiriman Last-Event-ID: 1 berhasil mengabaikan event 1 dan memutar ulang event 2. Pada bagian WebSocket, frame teks di-echo dengan benar dan frame biner DEADBEEF berhasil dibalik menjadi EFBEADDE tanpa kecacatan data.",
    title_en="Demo Execution & Verification Results",
    narrative_en="Running cmd/demo validates theoretical claims. SSE resumption honors Last-Event-ID: 1 by skipping initial events and replaying the subsequent update. WebSocket mirrors text and reverses binary DEADBEEF into EFBEADDE flawlessly."
)

# Slide 16: Verifikasi Unit Test & Concurrency
add_slide(
    categories=["VERIFIKASI TEST", "SUITE TEST"],
    title="Hasil Pengujian Test Suite & Concurrency",
    subtitle="4 test case diverifikasi lolos dengan Race Detector aktif (Go 1.22)",
    content_items=[
        {"title": "HASIL EKSEKUSI TEST", "items": [
            "=== RUN   TestSSE_Formatting            -> PASS (0.00s)",
            "=== RUN   TestSSE_LastEventID_Resumption -> PASS (0.00s)",
            "=== RUN   TestWebSocket_TextAndBinary    -> PASS (0.00s)",
            "=== RUN   TestSSE_Concurrency           -> PASS (0.05s)",
            "PASS: go test -race ./... (1.559s - Race Clean)",
        ], "type": "code"},
        {"title": "SKENARIO CONCURRENCY", "items": [
            "10 client subscriber terhubung secara paralel ke SSE Hub",
            "5 event broadcast dipicu secara simultan di goroutine terpisah",
            "Mutex dan unbuffered select mencegah deadlock maupun blocking antar client",
        ], "type": "bullet"},
        {"title": "DETERMINISTIK FORMATTING", "items": [
            "Test memvalidasi ordering field id, event, retry, data, dan delimiter ganda",
        ], "type": "bullet"},
    ],
    mental_model="Race-clean test pass proves thread-safe synchronization across concurrent streams.",
    source_path="labs/35-websocket-and-sse/tests/protocol_test.go",
    slide_num="16 / 24",
    narrative="Seluruh rangkaian pengujian di tests/protocol_test.go lulus dengan Go Race Detector aktif. Test concurrency membuktikan bahwa non-blocking broadcast channel pada Hub mampu melayani 10 client paralel saat menerima 5 event acak tanpa terjadi data race atau goroutine starvation.",
    title_en="Test Suite Verification & Concurrency",
    narrative_en="All test suites pass under the Go race detector. The concurrency test confirms that the SSE Hub non-blocking select pattern safely serves 10 parallel subscribers across 5 concurrent broadcasts without data races or lockups."
)

# Slide 17: Evaluasi Matriks 8 Dimensi
add_slide(
    categories=["TABEL EVALUASI", "PERBANDINGAN"],
    title="Matriks Evaluasi: WebSocket vs SSE",
    subtitle="Perbandingan mendalam 8 dimensi arsitektur dan operasional",
    content_items=[
        {"title": "TABEL PERBANDINGAN", "items": [
            "Dimensi           | WebSocket (RFC 6455)      | Server-Sent Events (SSE)",
            "Arah Komunikasi   | Full-Duplex (2 Arah)      | Unidirectional (Server->Client)",
            "Transport / Pipa  | Raw TCP (setelah Upgrade) | Standard HTTP Streaming",
            "Tipe Payload      | UTF-8 Teks & Raw Biner    | UTF-8 Teks Murni Saja",
            "Auto-Reconnect    | Manual / Custom Client    | Bawaan Browser EventSource",
            "Resumption State  | Tidak ada di spec level   | Bawaan via Last-Event-ID",
            "Proxy / LB Compat | Butuh Sticky / WS Bypass  | Standar HTTP (No buffering)",
            "HTTP/2 Multiplex  | Butuh RFC 8441 CONNECT    | Native HTTP/2 Multiplexing",
        ], "type": "table"},
        {"title": "RINGKASAN TRADE-OFF", "items": [
            "WebSocket unggul dalam fleksibilitas payload dan interaktivitas upstream",
            "SSE unggul mutlak dalam kesederhanaan operasional dan keandalan jaringan",
        ], "type": "bullet"},
    ],
    mental_model="WebSocket optimizes for protocol capability; SSE optimizes for operational simplicity.",
    source_path="labs/35-websocket-and-sse/content/02-master-draft.md",
    slide_num="17 / 24",
    narrative="Matriks perbandingan memperjelas trade-off arsitektural. WebSocket menawarkan kemampuan full-duplex dan dukungan biner dengan kompensasi hilangnya auto-reconnect bawaan dan meningkatnya kerumitan proxy. SSE memberikan ketahanan jaringan otomatis dan dukungan native HTTP/2 stream multiplexing.",
    title_en="Evaluation Matrix: WebSocket vs SSE",
    narrative_en="The evaluation matrix highlights the architectural trade-offs. WebSocket offers full-duplex streaming and binary payloads at the cost of operational proxy complexity. SSE delivers native network resilience and automatic HTTP/2 multiplexing."
)

# Slide 18: Interoperabilitas Proxy & Load Balancer
add_slide(
    categories=["INFRASTRUKTUR", "NETWORKING"],
    title="Tantangan Proxy, Firewall & Gateway",
    subtitle="Bagaimana middleware jaringan memperlakukan koneksi persisten",
    content_items=[
        {"title": "ISU PROXY PADA WEBSOCKET", "items": [
            "Load balancer modern (seperti AWS ALB) butuh sticky sessions jika node memegang state in-memory",
            "Beberapa firewall perusahaan memblokir request non-HTTP setelah Upgrade",
            "Idle timeout TCP memutus socket jika tidak ada traffic Ping/Pong aktif",
        ], "type": "bullet"},
        {"title": "ISU PROXY PADA SSE", "items": [
            "Buffering di perantara: Nginx menahan chunk response jika proxy_buffering aktif",
            "Solusi: Tambahkan header X-Accel-Buffering: no pada respons server",
            "Batas koneksi HTTP/1.1 browser (maks 6 koneksi per domain) sebelum HTTP/2",
        ], "type": "bullet"},
        {"title": "REKOMENDASI PRODUKSI", "items": [
            "Deploy SSE di atas HTTP/2 untuk menghilangkan batas koneksi domain",
            "Terapkan distributed message broker (Redis Pub/Sub) untuk scaling multi-node",
        ], "type": "bullet"},
    ],
    mental_model="Reverse proxies buffer HTTP unless instructed otherwise; firewalls distrust protocol upgrades.",
    source_path="labs/35-websocket-and-sse/content/02-master-draft.md",
    slide_num="18 / 24",
    narrative="Koneksi jangka panjang rentan terhadap timeout dan pemotongan oleh proxy perantara. WebSocket kerap terhambat oleh firewall korporat yang mencurigai protokol non-HTTP. SSE menghadapi risiko proxy buffering di Nginx, yang dapat diselesaikan dengan menyetel X-Accel-Buffering: no atau menonaktifkan proxy_buffering.",
    title_en="Proxy, Firewall & Gateway Interoperability",
    narrative_en="Long-lived connections are vulnerable to gateway buffering and idle terminations. WebSocket often encounters strict enterprise firewalls suspicious of protocol switching. SSE requires disabling proxy buffering via X-Accel-Buffering: no headers."
)

# Slide 19: Dampak HTTP/2 & HTTP/3
add_slide(
    categories=["STANDAR PROTOKOL", "HTTP/2"],
    title="Dampak HTTP/2 & Multiplexing Aliran",
    subtitle="Mengapa evolusi HTTP memperkuat posisi Server-Sent Events",
    content_items=[
        {"title": "SSE NATIVE DI HTTP/2", "items": [
            "SSE berjalan di atas stream HTTP/2 standar tanpa ekstensi",
            "Multiplexing memungkinkan puluhan stream SSE berbagi 1 koneksi TCP fisik",
            "Mengeliminasi batas browser 6 koneksi per origin pada HTTP/1.1",
        ], "type": "bullet"},
        {"title": "WEBSOCKET DI HTTP/2 (RFC 8441)", "items": [
            "HTTP/2 melarang header Connection: Upgrade (RFC 9113)",
            "WebSocket di HTTP/2 memerlukan Extended CONNECT method (:protocol=websocket)",
            "Dukungan server dan proxy gateway masih belum universal di industri",
        ], "type": "bullet"},
        {"title": "KEUNGGULAN OPERASIONAL", "items": [
            "SSE secara natural mendapatkan kompresi header HPACK dan flow control HTTP/2",
        ], "type": "bullet"},
    ],
    mental_model="HTTP/2 natively supercharges SSE while demanding complex extensions for WebSocket.",
    source_path="labs/35-websocket-and-sse/research/05-report.md",
    slide_num="19 / 24",
    narrative="Revolusi HTTP/2 secara dramatis meningkatkan daya tarik SSE. Melalui stream multiplexing, puluhan saluran SSE dapat berbagi satu koneksi TCP tanpa batasan 6 koneksi per domain. Sebaliknya, WebSocket memerlukan implementasi RFC 8441 yang rumit karena HTTP/2 melarang mekanisme connection upgrade tradisional.",
    title_en="Impact of HTTP/2 & Stream Multiplexing",
    narrative_en="HTTP/2 amplifies SSE's viability by multiplexing many event streams across a single TCP connection. Conversely, WebSocket over HTTP/2 mandates RFC 8441 Extended CONNECT mechanics because traditional connection upgrades are prohibited."
)

# Slide 20: Studi Kasus: Notifikasi & Streaming LLM
add_slide(
    categories=["STUDI KASUS", "INDUSTRI"],
    title="Studi Kasus Nyata: Notifikasi & AI Tokens",
    subtitle="Membedah pilihan arsitektur OpenAI, Claude, dan Sistem Notifikasi Skala Besar",
    content_items=[
        {"title": "STREAMING TOKEN LLM (OPENAI / ANTHROPIC)", "items": [
            "Karakteristik: Server menghasilkan token kata demi kata secara real-time",
            "Arah: 100% Server -> Client setelah prompt awal dikirim",
            "Pilihan arsitektur: Server-Sent Events (SSE) adalah standar industri universal",
        ], "type": "bullet"},
        {"title": "KOLABORASI DOKUMEN & CHAT (FIGMA / SLACK)", "items": [
            "Karakteristik: Kursor user bergerak cepat, pesan teks dikirim dua arah",
            "Arah: Bidirectional interaktif dengan latensi minimal",
            "Pilihan arsitektur: WebSocket dengan enkripsi frame biner",
        ], "type": "bullet"},
        {"title": "PEMBELAJARAN ARSITEKTURAL", "items": [
            "Penggunaan SSE pada API LLM membuktikan ketahanan dan kesederhanaan protokol",
        ], "type": "bullet"},
    ],
    mental_model="Match real-world traffic patterns: LLM tokens stream via SSE; multi-user editing flows via WebSocket.",
    source_path="labs/35-websocket-and-sse/content/02-master-draft.md",
    slide_num="20 / 24",
    narrative="Tren modern membuktikan penerapan aturan directionality. Layanan AI generative seperti OpenAI ChatGPT memilih SSE untuk men-stream token kata demi kata ke browser karena request prompt sudah selesai di awal. Sebaliknya, aplikasi seperti Figma atau multiplayer gaming mengandalkan WebSocket karena user terus-menerus mengirimkan data interaktif.",
    title_en="Case Study: Notifications & LLM Streaming",
    narrative_en="Production adoption reflects directionality rules. Generative AI providers stream tokens via SSE because traffic is strictly server-to-client following the initial prompt. Real-time collaborative tools like Figma depend on WebSocket for low-latency bidirectional interaction."
)

# Slide 21: Batasan Implementasi Lab 35
add_slide(
    categories=["BATAS LAB", "DISCLAIMER"],
    title="Batasan Model & Implementasi Lab 35",
    subtitle="Penyederhanaan pedagogis untuk menonjolkan prinsip mekanika protokol",
    content_items=[
        {"title": "BATAS 1: WEBSOCKET FRAMING", "items": [
            "Masking key client menggunakan nilai konstan (bukan crypto-random)",
            "Belum mendukung fragmented continuation frames (OpContinuation 0x0)",
            "Payload > 65,535 bytes memiliki parser BigEndian namun belum diuji ekstensif",
        ], "type": "bullet"},
        {"title": "BATAS 2: SSE HUB BUFFER", "items": [
            "History event disimpan dalam slice in-memory tak terbatas (unbounded buffer)",
            "Produksi memerlukan ring buffer dengan batas retensi atau Redis Streams",
        ], "type": "bullet"},
        {"title": "BATAS 3: SCALE & DISTRIBUTED STATE", "items": [
            "Uji coba berjalan pada single node HTTP server tanpa distributed message bus",
        ], "type": "bullet"},
    ],
    mental_model="A lab isolates protocol mechanics; production demands ring buffers, random masks, and brokers.",
    source_path="labs/35-websocket-and-sse/README.md",
    slide_num="21 / 24",
    narrative="Lab ini dirancang sebagai sarana edukasi teknis murni. Masking key diset statis untuk memudahkan debugging inspeksi byte, dan history SSE disimpan dalam memori tanpa penggusuran. Dalam lingkungan produksi riil, implementasi memerlukan ring buffer berkapasitas tetap, masking acak kriptografis, dan distributed broker.",
    title_en="Lab 35 Model Boundaries & Constraints",
    narrative_en="This lab focuses on educational clarity. Static masking keys simplify byte-level inspection, and the SSE hub retains history without buffer eviction. Production deployments require bounded ring buffers, cryptographically random masks, and message brokers."
)

# Slide 22: Panduan Checklist Pemilihan Protokol
add_slide(
    categories=["PANDUAN", "CHECKLIST KEPUTUSAN"],
    title="Checklist Keputusan Arsitektur Real-Time",
    subtitle="8 langkah verifikasi sebelum menentukan protokol transport aplikasi",
    content_items=[
        {"title": "EVALUASI KEBUTUHAN", "items": [
            "[ ] Apakah client perlu mengirim pesan upstream di socket yang sama? (Ya->WS, Tidak->SSE)",
            "[ ] Apakah data mengandung payload biner murni (audio/protobuf)? (Ya->WS)",
            "[ ] Apakah sistem berjalan di balik corporate proxy / strict firewall? (Pilih SSE)",
            "[ ] Apakah pemulihan event otomatis (reconnect) diperlukan secara gratis? (Pilih SSE)",
        ], "type": "bullet"},
        {"title": "VERIFIKASI PRODUKSI", "items": [
            "[ ] Pastikan Nginx menonaktifkan proxy_buffering untuk SSE endpoint",
            "[ ] Pasang heartbeat / ping-pong berkala untuk mencegah idle connection drops",
            "[ ] Atur limit OS file descriptor (ulimit -n) untuk menangani ribuan koneksi",
            "[ ] Siapkan Redis Pub/Sub untuk replikasi broadcast lintas cluster node",
        ], "type": "bullet"},
    ],
    mental_model="Verify directionality, payload type, and proxy topology before committing to a protocol.",
    source_path="labs/35-websocket-and-sse/content/02-master-draft.md",
    slide_num="22 / 24",
    narrative="Gunakan checklist ini dalam tahap perancangan sistem real-time. Jika komunikasi upstream dapat dipenuhi melalui HTTP POST standar dan client hanya mendengarkan event, SSE adalah pilihan terbaik. Pastikan konfigurasi reverse proxy telah disesuaikan dan kapasitas file descriptor sistem operasi memadai.",
    title_en="Protocol Decision Checklist",
    narrative_en="Utilize this checklist during architectural design phases. If upstream requests can be satisfied via standard HTTP POST while clients listen passively, SSE is the superior path. Verify reverse proxy settings and OS file descriptor limits."
)

# Slide 23: Ringkasan Prinsip Arsitektur
add_slide(
    categories=["PRINSIP INTI", "KESIMPULAN TEKNIS"],
    title="Tiga Prinsip Inti Komunikasi Real-Time",
    subtitle="Ringkasan aksioma untuk membangun sistem streaming yang efisien",
    content_items=[
        {"title": "1. DIRECTIONALITY IS KING", "items": [
            "Jangan bayar overhead koneksi dua arah jika trafik hanya mengalir satu arah",
            "SSE adalah peluru perak untuk live notifications, feed, dan token streaming",
        ], "type": "bullet"},
        {"title": "2. LEVERAGE PLATFORM PRIMITIVES", "items": [
            "Manfaatkan EventSource dan Last-Event-ID bawaan browser untuk auto-healing",
            "Gunakan HTTP/2 multiplexing untuk menghemat handshake TCP koneksi",
        ], "type": "bullet"},
        {"title": "3. PROTOCOL RESPECTS BOUNDARIES", "items": [
            "WebSocket adalah pilihan mutlak ketika interaktivitas biner dua arah dibutuhkan",
            "Keduanya saling melengkapi, bukan saling meniadakan",
        ], "type": "bullet"},
    ],
    mental_model="Right protocol for the right direction: simplicity scales better than unnecessary duplexing.",
    source_path="labs/35-websocket-and-sse/content/05-key-takeaways.md",
    slide_num="23 / 24",
    narrative="Sistem real-time yang elegan memilih protokol berdasarkan karakteristik data transport. Hindari godaan over-engineering dengan memasang WebSocket di setiap use case. Dengan memilih SSE untuk pola satu arah dan WebSocket untuk pola interaktif dua arah, sistem tetap efisien, mudah di-scale, dan tahan banting.",
    title_en="Three Core Real-Time Principles",
    narrative_en="Resilient systems choose transport protocols matching traffic patterns. Avoid over-engineering by deploying WebSocket indiscriminately. Reserving SSE for unidirectional streams and WebSocket for true full-duplex interactions yields scalable architectures."
)

# Slide 24: Closing (Dark Navy)
add_slide(
    categories=["RINGKASAN", "REPOSITORI"],
    title="Eksperimen kecil, eliminasi latensi tersembunyi.",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | 4/4 Tests PASS (Go 1.22 Race Clean) | {AUTHOR}",
    content_items=[
        {"title": "KEY RESULTS", "items": [
            "internal/sse • internal/ws",
            "internal/server • cmd/demo",
            "tests/protocol_test.go • README.md",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Semua test PASS dengan Go 1.22 Race Detector aktif",
            "Implementasi murni zero external dependencies",
        ], "type": "bullet"},
    ],
    mental_model="Directionality dictates the protocol: server-push belongs to SSE, bidirectional belongs to WebSocket.",
    source_path="labs/35-websocket-and-sse/README.md",
    slide_num="24 / 24",
    narrative="Penelitian dan pengujian Lab 35 membuktikan bahwa pemahaman mendalam atas protokol transport menghemat biaya komputasi dan kompleksitas arsitektur.",
    title_en="Small experiment, zero hidden latency.",
    narrative_en="Lab 35 demonstrates that deep protocol comprehension eliminates architectural complexity and operational overhead."
)

if __name__ == "__main__":
    print(f"Loaded content_def for Lab {LAB_NO}: {len(SLIDES)} slides defined.")
