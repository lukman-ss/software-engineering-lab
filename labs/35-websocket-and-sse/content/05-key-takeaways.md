# Key Takeaways

1. **Arah komunikasi menentukan protokol.** SSE untuk server→client push. WebSocket untuk bidirectional.

2. **SSE adalah HTTP standard.** Tidak perlu protocol upgrade, tidak perlu RFC 8441 untuk HTTP/2. Browser `EventSource` menangani reconnection otomatis.

3. **WebSocket butuh handshake khusus.** HTTP/1.1 `Upgrade: websocket` + `101 Switching Protocols`. Over HTTP/2 perlu RFC 8441 Extended CONNECT.

4. **`Last-Event-ID` memungkinkan resumption.** Client reconnect dengan header ini, server meng-replay event yang terlewat. Ini built-in di browser.

5. **WebSocket mendukung binary; SSE hanya UTF-8 text.** Untuk data biner, WebSocket satu-satunya pilihan.

6. **HTTP/1.1 membatasi SSE ke ~6 koneksi per origin.** HTTP/2 menghilangkan bottleneck ini dengan stream multiplexing.

7. **Proxy harus dikonfigurasi untuk real-time.** `X-Accel-Buffering: no` (nginx), `proxy_read_timeout` yang cukup tinggi.

8. **Scaling 100k connections bukan masalah protokol.** Terbatas oleh OS file descriptor limits dan per-connection memory, bukan pilihan WebSocket vs SSE.

9. **Lab ini adalah minimal implementation.** Bukan full RFC 6455 spec. Static mask key, unbounded history buffer, fragmented frames belum dihandle.

10. **Gunakan SSE untuk: notification, live dashboard, LLM token streaming.** Gunakan WebSocket untuk: chat, multiplayer games, collaborative editing.
