Lab 37 kali ini membahas strategi cache invalidation dalam sistem terdistribusi, dengan fokus pada pola write, mitigasi stampede, dan teknik TTL jitter dalam implementasi Go.

Cache invalidation adalah salah satu dilema terbesar dalam arsitektur sistem terdistribusi skala tinggi. Tanpa mekanisme yang disiplin dan terukur, sistem akan menghadapi tiga masalah konkret yang fatal:
• Cache stampede: Ribuan request konkuren membanjiri database secara serentak saat kunci populer kedaluwarsa mendadak
• Stale data: Goroutine pembaca menerima data usang karena inkonsistensi jeda waktu antara penulisan database dan pembaruan cache
• Risiko kehilangan data: Pilihan throughput ekstrem pada Write-Behind mengorbankan durabilitas saat antrean memory overflow atau proses crash

Eksperimen Go standard library pada lab ini membuktikan implementasi native:
• Cache-Aside meminimalkan penggunaan memori namun rentan terhadap race condition pembacaan usang
• Write-Through menjamin konsistensi read-after-write sinkron dengan kompromi kenaikan latensi penulisan
• Write-Behind mengalirkan mutasi asinkron ke worker latar belakang dengan throughput tinggi
• Singleflight memadatkan 20 request konkuren menjadi tepat 1 query database tunggal
• XFetch menggeser titik rebuild proaktif secara probabilistik sebelum batas waktu kedaluwarsa
• Stale-While-Revalidate menyajikan data usang seketika sambil merevalidasi kunci di latar belakang
• TTL Jitter mengacak rentang waktu kedaluwarsa untuk mengeliminasi lonjakan serentak

Mental model: Cache invalidation is the hardest problem because consistency is eventual, never instant.

Batasan implementasi utama: Seluruh simulasi menggunakan memory cache dan MockDB Go murni tanpa Redis atau PostgreSQL persisten eksternal.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies

#softwareengineering #golang #systemdesign #distributed #performance #caching