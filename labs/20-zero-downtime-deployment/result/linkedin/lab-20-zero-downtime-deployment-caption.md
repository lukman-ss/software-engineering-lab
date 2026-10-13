Lab 20 kali ini membahas arsitektur zero-downtime deployment, koordinasi probe, penundaan terminasi jaringan, serta evolusi skema database paralel di Go.

Masalah teknis yang sering terjadi adalah pemutusan sepihak instance lama saat deployment memicu lonjakan error HTTP 502 Bad Gateway pada request in-flight, membatalkan job background di tengah eksekusi, dan menyebabkan crash akibat skema database yang belum kompatibel.

Lab ini membuktikan empat pilar arsitektur konkret:
• Decoupled Probes: Liveness probe tetap melaporkan status proses sehat sementara readiness probe menandai unready untuk melepas instance dari rotasi load balancer.
• Asynchronous Detachment: Hook preStop berbasis select-timer menunda penutupan listener HTTP agar tabel routing eksternal sempat diperbarui sebelum soket ditutup.
• Connection Draining: Server menyelesaikan seluruh request aktif in-flight secara utuh sebelum proses dihentikan, diverifikasi tanpa kebocoran transaksi.
• Expand and Contract: Lapisan penyimpanan database mendukung pembacaan fallback transparan antara skema warisan dan skema modern tanpa downtime migrasi.

Mental model inti: Start new before stopping old. Coexistence requires version N and version N-1 to run simultaneously, route traffic via decoupled probes, and evolve database schemas in backward-compatible phases.

Satu batasan penting: Lab ini menggunakan simulasi in-memory Go untuk menguji alur draining HTTP, worker queue, dan skema database, bukan cluster Kubernetes fisik atau database PostgreSQL nyata, serta mengasumsikan drain timeout worker membatalkan konteks jika pekerjaan melampaui batas waktu.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment

#softwareengineering #golang #devops #kubernetes #systemdesign #backend
