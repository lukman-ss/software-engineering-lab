# WebSocket vs Server-Sent Events: Memilih Protokol Real-Time yang Tepat

## Problem

Ketika membangun sistem real-time — entah itu notifikasi, live dashboard, atau streaming respons LLM — developer menghadapi pilihan arsitektural fundamental: gunakan WebSocket atau Server-Sent Events (SSE)? Keduanya memungkinkan server mendorong data ke client secara real-time, tetapi berbeda drastis dalam model komunikasi, kompleksitas protokol, dan perilaku operasional.

Kesalahan pemilihan protokol menimbulkan biaya operasional yang tidak perlu: WebSocket membawa kompleksitas handshake, reconnection manual, dan konfigurasi proxy yang lebih rumit. SSE, sebaliknya, terbatas pada satu arah (server→client) dan tidak mendukung payload biner.

## Why This Matters

Keputusan ini bukan sekadar akademis. Untuk aplikasi notifikasi atau live feed, SSE menawarkan kompleksitas yang lebih rendah dan interoperabilitas proxy yang lebih baik. Untuk chat atau game multiplayer, WebSocket adalah satu-satunya pilihan karena membutuhkan komunikasi bidirectional. Menggunakan WebSocket ketika SSE cukup berarti menambah beban operasional tanpa benefit fungsional. Sebaliknya, menggunakan SSE ketika bidirectional dibutuhkan berarti tidak ada solusi sama sekali.

## Mental Model

Inti keputusan ini dapat direduksi menjadi satu pertanyaan: **apakah client perlu mengirim pesan ke server secara independen?**

- Jika **tidak** — server hanya mendorok data ke client (notifikasi, update dashboard, streaming token) → **SSE**
- Jika **ya** — client juga harus mengirim pesan ke server (chat, kolaboratif, game) → **WebSocket**

Model mental ini mengeliminasi banyak pertimbangan teknis yang rumit. Arah komunikasi menentukan segalanya: model protokol, perilaku browser, konfigurasi proxy, dan pola scaling.

## Core Concept

**WebSocket** (RFC 6455) adalah protokol aplikasi non-HTTP yang menyediakan saluran full-duplex melalui satu koneksi TCP. Handshake dimulai dengan permintaan HTTP/1.1 `Upgrade: websocket` yang menghasilkan respons `101 Switching Protocols`. Setelah handshake selesai, komunikasi berlangsung di luar siklus request-response HTTP — kedua pihak dapat mengirim frame secara independen. WebSocket mendukung payload UTF-8 text (opcode 0x1) dan binary data (opcode 0x2).

**Server-Sent Events (SSE)** adalah pola HTTP persistent dengan format `text/event-stream` yang didefinisikan oleh WHATWG HTML Living Standard. SSE secara ketat unidirectional: server mengirim event ke client, client tidak bisa mengirim pesan melalui koneksi SSE yang sama. Browser mengelola koneksi melalui API `EventSource` yang menyertakan auto-reconnect dan `Last-Event-ID` — sebuah mekanisme yang secara native memungkinkan client meminta event yang terlewat saat reconnect.

## Failure Scenario

Bayangkan sebuah sistem notifikasi yang menggunakan WebSocket untuk mengirim update ke 10,000 client. Setiap proxy load balancer umumnya membutuhkan konfigurasi WebSocket sticky sessions — koneksi harus tetap di node yang sama karena state koneksi ada di memori lokal — kecuali ada broker pesan terpusat. Ketika sebuah node restart, semua koneksi WebSocket terputus dan client harus reconnect secara manual (tidak ada auto-reconnect di level protokol). Proses reconnect memicu storm notifikasi ulang ke semua client secara bersamaan.

Skenario yang sama dengan SSE: proxy tidak memerlukan sticky session karena SSE adalah HTTP standard. Load balancer bisa mendistribusikan koneksi ke node manapun. Ketika node restart, client `EventSource` secara otomatis reconnect dan mengirim `Last-Event-ID` header untuk mengambil event yang terlewat. Proses recovery terjadi tanpa intervensi developer. Untuk replay lintas node tetap perlu shared history/broker seperti Redis Pub/Sub — in-memory history bersifat per-node.

## How It Works

### SSE: Alur Event

Ketika client melakukan GET ke endpoint `/sse` dengan header `Last-Event-ID`, server menghitung event mana yang terlewat dan mengirimkannya sebagai replay. Kemudian, server memasuki loop streaming di mana event broadcast dikirim ke channel per-client dan keep-alive comment (`:`) dikirim setiap 15 detik untuk mencegah timeout proxy.

Setiap event di-format sesuai spesifikasi WHATWG: field `id`, `event`, `data` (dengan support multi-line), dan `retry` dipisahkan oleh newline, diakhiri dengan double-newline sebagai delimiter blok.

### WebSocket: Alur Frame

Ketika client mengirim koneksi ke `/ws`, server memvalidasi header `Upgrade: websocket`, membaca `Sec-WebSocket-Key`, dan menghitung `Sec-WebSocket-Accept` sebagai `base64(SHA1(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))`. Setelah respons `101 Switching Protocols`, koneksi di-hijack dari `net/http` untuk raw TCP access.

Frame-frame selanjutnya menggunakan struktur RFC 6455: dua byte header berisi opcode (4 bit), FIN flag (1 bit), dan payload length (7 bit atau extended 16/64 bit). Jika payload dimasking (dari client ke server), XOR unmasking dilakukan dengan `data[i] ^= mask[i % 4]`.

## Architecture

Server menggunakan `http.ServeMux` tunggal yang menghosting tiga endpoint: `/sse` (Handler SSE Hub), `/ws` (WebSocket upgrade + handler loop), dan `/publish` (API untuk memicu broadcast SSE). SSE Hub memelihara in-memory history buffer dan map channel per-client. WebSocket handler menjalankan loop read-frame/process/write per koneksi.

Arsitektur ini menggunakan **pure Go standard library** — tidak ada dependensi eksternal. Handshake RFC 6455, frame masking/unmasking, dan framing dilakukan langsung melalui connection hijacking `net/http`. Ini menjamin zero-dependency reproducibility dan memberikan insight langsung ke mekanisme byte-level.

## Implementation

### SSE Hub (`internal/sse/sse.go`)

`Hub` mengelola state dengan `sync.RWMutex`. `Broadcast` menambahkan event ke history dengan monotonic ID, lalu mengirim ke semua channel client melalui non-blocking `select`/`default` — jika channel client penuh (capacity 64), event di-drop tanpa memblokir broadcast ke client lain. `Subscribe` menghitung replay slice (event dengan `ID > lastID`), membuat channel buffered, dan mendaftarkannya ke clients map. `ServeHTTP` mengatur header SSE yang diperlukan (`Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`), membaca `Last-Event-ID` header, dan menjalankan loop select untuk live events, replay, dan heartbeat.

### WebSocket (`internal/ws/ws.go`)

`Upgrade` memvalidasi handshake HTTP, menghitung accept key, dan hijack koneksi. `ReadFrame` membaca header 2-byte, mengekstrak opcode dan length (dengan support extended 16/64-bit), membaca mask key 4-byte jika ada, dan menerapkan XOR unmasking. `WriteFrame` membangun frame dengan FIN=1, length encoding yang benar untuk semua rentang, dan optional masking. `Dial` melakukan handshake client-side atas raw TCP connection.

### Server (`internal/server/server.go`)

`NewServer` membuat SSE Hub dan ServeMux. Handler `/ws` menjalankan loop yang memproses `OpClose` (kirim close frame, return), `OpPing` (kirim pong), `OpText` (echo dengan prefix "echo: "), dan `OpBinary` (reversed payload). Handler `/publish` memicu broadcast ke semua SSE subscriber.

## Code Walkthrough

**SSE Event Formatting** — `Event.Format()` membangun string sesuai WHATWG spec. Multi-line data dipisah per baris, masing-masing mendapat prefix `data: `. ID hanya dicetak jika `> 0`, retry jika `> 0`. Double-newline dihasilkan dari trailing newline di baris `data:` terakhir plus `sb.WriteString("\n")`.

**WebSocket Frame Reading** — Dua byte pertama berisi opcode (4 bit) dan payload length (7 bit). Jika length == 126, 2 byte tambahan dibaca untuk 16-bit extended length. Jika length == 127, 8 byte untuk 64-bit. XOR unmasking `data[i] ^= mask[i%4]` sesuai RFC 6455 Section 5.3.

**WebSocket Binary Reversal** — Ketika server menerima binary frame, payload dibalik sebelum dikirim kembali. Ini membuktikan byte-level integrity over the wire, yang diverifikasi dalam test dengan membandingkan `[]byte{0x01, 0x02, 0x03, 0x04}` yang dikirim terhadap `[]byte{0x04, 0x03, 0x02, 0x01}` yang diterima.

**Last-Event-ID Resumption** — Client subscribe dengan `Last-Event-ID: 1`. Server menghitung semua event dengan `ID > 1` dari history dan mengirimkannya sebagai replay. Event ID 1 tidak muncul di stream. Ini dibuktikan oleh `TestSSE_LastEventID_Resumption`.

## What the Tests Prove

Empat test dalam `tests/protocol_test.go` memverifikasi perilaku kunci:

- **`TestSSE_Formatting`** — Membandingkan string output `Event.Format()` secara eksplisit terhadap WHATWG wire format: `"id: 42\nevent: update\nretry: 1000\ndata: hello\ndata: world\n\n"`. Ini membuktikan multi-line data framing dan field ordering.

- **`TestSSE_LastEventID_Resumption`** — Broadcast 3 event, client connect dengan `Last-Event-ID: 1`, verifikasi event ID 1 tidak ada di stream, event ID 2 dan 3 ada. Ini membuktikan mekanisme replay bekerja.

- **`TestWebSocket_TextAndBinary`** — Mengirim text frame dan memverifikasi echo prefix, mengirim binary frame `[]byte{0x01, 0x02, 0x03, 0x04}` dan memverifikasi reversed payload `[]byte{0x04, 0x03, 0x02, 0x01}`. Ini membuktikan bidirectional framing bekerja untuk kedua tipe payload.

- **`TestSSE_Concurrency`** — 10 concurrent client subscribe, 5 concurrent broadcasts, race detector bersih. Ini membuktikan thread-safety pada SSE Hub.

Semua test lolos dengan dan tanpa race detector (`go test -race ./...`). Demo (`go run ./cmd/demo`) juga menghasilkan output yang sesuai harapan.

## Recovery / Rollback

SSE memiliki mekanisme recovery bawaan di level browser: `EventSource` API secara otomatis reconnect dengan jeda sesuai field `retry` (implementation-defined, biasanya beberapa detik) dan mengirim `Last-Event-ID` header. Server menghitung event yang terlewat dari history dan mengirimkannya sebagai replay. Ini adalah perilaku spec-defined, bukan implementasi kustom.

WebSocket tidak memiliki mekanisme recovery otomatis di level protokol. RFC 6455 Section 7.2.3 menggambarkan recovery dari abnormal closure sebagai implement concern. Developer harus membangun logika reconnect sendiri — pustaka seperti `ReconnectingWebSocket` ada untuk mengisi gap ini, tetapi bukan bagian dari spesifikasi.

## Production Considerations

**Proxy dan Load Balancer:** SSE membutuhkan `proxy_buffering` dimatikan (nginx: `proxy_buffering off`) dan `proxy_read_timeout` yang cukup tinggi. `X-Accel-Buffering: no` header yang digunakan dalam implementasi ini secara spesifik menonaktifkan buffering nginx. WebSocket lebih rentan terhadap proxy issues dan memerlukan sticky session untuk stateful routing — AWS ELB mendokumentasikan ini sebagai persyaratan khusus WebSocket.

**HTTP/2:** SSE berjalan native di HTTP/2 melalui stream multiplexing — tidak memerlukan ekstensi khusus. WebSocket over HTTP/2 memerlukan RFC 8441 Extended CONNECT method dengan `:protocol = websocket` pseudo-header karena HTTP/2 melarang connection-wide headers seperti `Upgrade` dan `Connection`.

**Browser Limits:** SSE over HTTP/1.1 dibatasi hingga ~6 koneksi per origin (kebijakan browser de facto, bukan web standard). HTTP/2 menghilangkan bottleneck ini dengan negotiated stream limits (default 100). WebSocket tidak dihitung dalam batas 6 koneksi HTTP/1.1 tersebut (tetapi browser tetap membatasi jumlah koneksi WS secara terpisah).

**Scaling 100k Connections:** Klaim tentang 100,000 koneksi bersamaan adalah prinsip arsitektural yang dibahas dalam penelitian, bukan di-benchmark di lab ini. Batasan utama adalah OS file descriptor limits (`ulimit -n`) dan per-connection memory overhead dalam runtime. Untuk broadcast lintas node, diperlukan message broker seperti Redis Pub/Sub.

## Common Mistakes

1. **Menggunakan WebSocket untuk kebutuhan SSE.** Ini menambah kompleksitas handshake, reconnect manual, dan konfigurasi proxy tanpa benefit fungsional jika client hanya perlu menerima data.

2. **Menganggap SSE tidak bisa scale.** SSE sebenarnya lebih mudah di-scale secara horizontal karena merupakan HTTP standard — tidak memerlukan sticky session, proxy tidak perlu WebSocket-specific configuration.

3. **Mengasumsikan auto-reconnection ada di WebSocket.** RFC 6455 tidak mendefinisikan mekanisme reconnect. Developer harus implement sendiri atau gunakan pustaka pihak ketiga.

4. **Mengabaikan proxy buffering.** Baik WebSocket maupun SSE memerlukan tuning proxy untuk long-lived connections, tetapi SSE lebih rentan karena proxy mungkin men-buffer response HTTP yang panjang.

5. **Menganggap lab ini adalah implementasi production-ready.** Implementasi ini adalah minimal demonstration — static mask key untuk WS, unbounded history untuk SSE, dan tidak ada fragmented frame support.

## Case Study

**Demo Lab: SSE Resumption dari Last-Event-ID**

Server melakukan broadcast 2 event berita. Client connect dengan `Last-Event-ID: 1`. Server menghitung bahwa event ID 2 harus di-replay (karena `evt.ID > 1`), sementara event ID 1 dilewati. Client menerima:

```
[SSE Stream] id: 2
[SSE Stream] event: news
[SSE Stream] retry: 2000
[SSE Stream] data: Update: Faster routing in net/http
```

Ini membuktikan bahwa `Last-Event-ID` mechanism bekerja sesuai WHATWG HTML spec — event yang terjadi sebelum reconnect tidak diulangi, event yang terjadi sesudahnya di-replay.

**Demo Lab: WebSocket Binary Integrity**

Client mengirim binary payload `DEADBEEF` melalui WebSocket frame. Server membalik payload byte-by-byte dan mengirim kembali `EFBEADDE`. Client memverifikasi kecocokan byte. Ini membuktikan bahwa WebSocket framing mempertahankan integritas payload pada level byte — tidak ada transformasi atau encoding yang mengubah data biner.

## Checklist

- [ ] Tentukan apakah komunikasi hanya server→client (SSE) atau bidirectional (WebSocket)
- [ ] Jika SSE: atur `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`
- [ ] Jika WebSocket: validasi header `Upgrade: websocket`, hitung `Sec-WebSocket-Accept`, hijack koneksi
- [ ] Implementasikan `Last-Event-ID` handling untuk SSE recovery
- [ ] Konfigurasikan proxy: nonaktifkan buffering, set timeout cukup tinggi
- [ ] Jika over HTTP/2: gunakan RFC 8441 untuk WebSocket (bukan upgrade header)
- [ ] Implementasikan reconnect logic jika WebSocket (tidak ada di level protokol)
- [ ] Atur file descriptor limits dan memory monitoring untuk scale

## Key Takeaways

- **Arah komunikasi adalah faktor penentu utama.** Server-push-only = SSE. Bidirectional = WebSocket.
- **SSE memiliki auto-reconnect bawaan browser.** WebSocket tidak — developer harus build sendiri.
- **SSE berjalan native di HTTP/2.** WebSocket memerlukan RFC 8441 untuk HTTP/2.
- **SSE lebih sederhana untuk proxy.** Tidak perlu sticky session, tidak perlu WebSocket-specific proxy config.
- **WebSocket mendukung binary, SSE hanya UTF-8.** Untuk data biner, WebSocket wajib.
- **Scaling 100k connections adalah masalah OS/runtime, bukan protokol.** File descriptor dan memory adalah bottleneck utama.
- **Implementasi lab ini adalah minimal.** Bukan production-ready — static mask key, unbounded buffer, dan missing edge cases adalah trade-off pedagogis yang disengaja.

## Sources

- RFC 6455 – The WebSocket Protocol (December 2011) — https://datatracker.ietf.org/doc/html/rfc6455
- RFC 8441 – Bootstrapping WebSockets with HTTP/2 (September 2018) — https://datatracker.ietf.org/doc/html/rfc8441
- RFC 9113 – HTTP/2 (June 2022, obsoletes RFC 7540) — https://datatracker.ietf.org/doc/html/rfc9113
- WHATWG HTML Living Standard Section 9.2 — Server-Sent Events — https://html.spec.whatwg.org/multipage/server-sent-events.html
- MDN: WebSocket — https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- MDN: EventSource — https://developer.mozilla.org/en-US/docs/Web/API/EventSource
- Linux epoll(7) — https://man7.org/linux/man-pages/man7/epoll.7.html
- NGINX ngx_http_proxy_module — https://nginx.org/en/docs/http/ngx_http_proxy_module.html
- AWS ELB WebSocket Support — https://docs.aws.amazon.com/elasticloadbalancing/latest/application/websockets-support.html

---

**Lab:** labs/35-websocket-and-sse
**Status:** APPROVED (Research & Engineering)
**Implementation:** Pure Go standard library, zero external dependencies, Go 1.22
**Tests:** 4 tests, all PASS with race detector
