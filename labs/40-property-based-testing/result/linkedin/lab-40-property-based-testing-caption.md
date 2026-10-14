Lab 40 kali ini membahas Property-Based Testing (PBT), kelemahan mendasar unit test berbasis contoh, serta pembuktian kebenaran sistem via universal invariant dan counterexample shrinking di Go.

Masalah teknis yang sering terjadi di produksi adalah unit test konvensional memberikan rasa aman palsu. Pengembang hanya menguji skenario yang mereka antisipasi, sehingga kode berstatus lulus seratus persen meski menyimpan cacat batas domain. Pada contoh mata uang berbasis float64, format naif lulus seluruh contoh statis, namun gagal total saat menerima angka pecahan acak akibat keterbatasan aproksimasi biner IEEE 754. Begitu pula pada penggabungan interval, algoritma naif gagal pada 85 dari 100 irisan acak karena berasumsi data selalu terurut.

Lab ini membuktikan empat pilar utama:
• Roundtrip Invariant: Representasi integer cents int64 lulus 1.000 iterasi acak tanpa kehilangan presisi pecahan sedikit pun.
• Test Oracle & Idempotence: RobustMerge membuktikan relasi f(f(x)) == f(x) dan hasil non-overlapping pada 1.000 kombinasi irisan acak.
• Automated Shrinking: Algoritma pemotongan biner dan reduksi nilai memangkas array acak 10 elemen menjadi satu counterexample minimal [-1] hanya dalam 21 langkah evaluasi.
• Zero Dependency: Seluruh pengujian properti dibangun murni menggunakan pustaka standar Go testing/quick.

Mental model inti: Example-based tests check developer assumptions; property-based tests explore domain boundaries.

Pelajari implementasi lengkap dan hasil pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/40-property-based-testing

#softwareengineering #golang #testing #qualityassurance #systemdesign
