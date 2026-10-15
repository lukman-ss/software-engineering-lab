Ketika sebuah service menginstansiasi adapter eksternal secara langsung di dalam metodenya (seperti memanggil RealGateway di dalam Processor), kode domain terkunci rapat dengan infrastruktur fisik. Pengujian unit menjadi mustahil tanpa membuka koneksi jaringan nyata, dan pengujian lokal sering kali gagal saat server pihak ketiga mengalami downtime.

Di Lab 16 seri Software Engineering Lab, kami mengeksplorasi kontras struktural antara Constructor Injection dan Service Locator di bahasa Go menggunakan manual wiring murni tanpa framework pihak ketiga.

Poin kuncinya: Inversion of Control memisahkan fase konfigurasi (perakitan graf dependensi di composition root main.go) dari fase eksekusi (logika domain). Melalui NewProcessor(g PaymentGateway), kontrak dependensi terlihat gamblang di signature kompilasi dan memungkinkan penyuntikan MockGateway seketika dalam 0.00 detik pengujian unit.

Sebaliknya, Service Locator dengan NewBadProcessor(c Container) menutupi kebutuhan dependensi di balik kontainer generik. Standar industri seperti PSR-11 menegaskan rekomendasi bahwa container tidak boleh disuntikkan ke domain logic karena memindahkan kegagalan kompilasi menjadi runtime error.

Test suite mencakup 6 test cases yang memverifikasi isolasi mock, delegasi pembayaran, propagasi error, dan penolakan nilai negatif sebelum menyentuh gateway, seluruhnya bersih dari race condition (go test -race ./...). Namun penting dicatat: pengujian mock unit membuktikan isolasi logika, bukan jaminan stabilitas koneksi jaringan atau resiliency di skala beban production.

Kode sumber lengkap, suite pengujian, dan implementasi demo tersedia di repositori:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection

#SoftwareEngineering #GoLang #DependencyInjection #CleanArchitecture #SystemDesign #SoftwareDesign