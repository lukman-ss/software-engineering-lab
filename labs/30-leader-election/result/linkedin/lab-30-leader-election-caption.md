Lab 30 kali ini membahas implementasi Leader Election dengan Fencing Tokens untuk mencegah split-brain dan mutasi zombie pada sistem terdistribusi Go.

Tantangan utama pada koordinasi terdistribusi adalah ilusi waktu sinkron. Ketika leader primer mengalami stop-the-world GC pause atau partisi jaringan melebihi batas sewa Time-To-Live, koordinator menganggap node tersebut mati dan melimpahkan mandat ke kandidat baru. Saat node lama terbangun tanpa sadar masa sewanya telah hangus, ia mengeksekusi operasi penulisan basi yang menimpa mutasi baru dan merusak integritas shared storage.

Lab ini membuktikan lima prinsip rekayasa sistem:
• Linearizable Coordinator: Koordinator thread-safe menerbitkan sewa berdurasi waktu dan menaikkan nomor revisi monotonik secara atomik.
• Bounded Failure Detection: Heartbeat renewal berkala mendeteksi kegagalan leader dan mengeksekusi failover otomatis ke node siaga.
• Strict Check-and-Set: Gerbang penyimpanan bersama menolak penulisan jika nomor token yang datang bernilai kurang dari atau sama dengan nomor token terakhir yang terlihat.
• Concurrency Safety: Lima kandidat bersaing simultan memperebutkan sewa dengan hasil deterministik tepat satu pemenang tanpa race condition.
• Monotonic Clocks: Perhitungan masa kadaluwarsa sewa memanfaatkan jam monotonik Go untuk menghindari lompatan jam dinding NTP.

Mental model inti: A lease grants time-bounded authority; a fencing token converts that temporal authority into a monotonic ordering that shared storage enforces, rejecting stale zombie writes.

Batasan penting: Lab ini menguji logika koordinasi via memori internal dan mutex Go, bukan kluster etcd fisik multi-node. Parameter waktu disederhanakan untuk kecepatan eksekusi demo.

Pelajari kode sumber dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election

#softwareengineering #golang #distributedsystems #concurrency #systemdesign #backend