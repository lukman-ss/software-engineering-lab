Lab 35 kali ini membahas perbandingan arsitektur antara WebSocket (RFC 6455) dan Server-Sent Events (WHATWG SSE) untuk menentukan protokol transport real-time yang tepat berdasarkan directionality komunikasi.

Banyak developer mengasumsikan WebSocket adalah solusi default untuk seluruh kebutuhan real-time. Pilihan ini sering memicu kompleksitas yang tidak perlu ketika use case sebenarnya murni satu arah dari server ke client, seperti notifikasi atau live dashboard. WebSocket memaksakan koneksi TCP persisten non-HTTP, jabat tangan Upgrade, serta penanganan rekoneksi manual di sisi client. Sebaliknya, SSE menyediakan streaming berbasis HTTP standar dengan dukungan rekoneksi otomatis bawaan peramban.

Eksperimen Go standard library pada lab ini membuktikan:
• SSE Resumption: Header Last-Event-ID secara presisi melewati event lama dan memutar ulang pembaruan data yang terlewat setelah koneksi terputus.
• WebSocket Duplexing: Framing biner RFC 6455 dengan XOR unmasking membalikkan payload biner secara deterministik membuktikan integritas byte murni.
• Thread-Safe Concurrency: SSE Hub melayani sepuluh subscriber simultan dengan broadcast paralel tanpa data race.

Mental model: Directionality dictates the protocol. Server-push belongs to SSE, bidirectional belongs to WebSocket.

Batasan implementasi: Lab ini menggunakan masking key statis untuk WebSocket dan buffer memori tak terbatas pada SSE Hub sebagai penyederhanaan pedagogis.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/35-websocket-and-sse

#softwareengineering #golang #webdevelopment #networking #systemdesign