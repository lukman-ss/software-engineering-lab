Dokumentasi arsitektur sering kali menjadi artefak yang usang begitu kode berevolusi. Masalah yang lebih berbahaya: saat satu keputusan menggantikan keputusan lama, tautan penggantian sering kali lupa diperbarui di kedua berkas. Review manual pada dokumen Markdown rentan melewatkan relasi sepihak, penomoran urut yang loncat, atau status informal yang tidak baku.

Di Lab 17 dari seri Software Engineering Lab, kami mengeksplorasi bagaimana memperlakukan Architecture Decision Record (ADR) sebagai data terstruktur, bukan sekadar teks naratif bebas.

Dengan membangun parser baris berbasis regex dan linter konkuren di Go (tanpa dependensi pihak ketiga), sistem dapat memvalidasi urutan nomor monotonik, mendeteksi siklus sirkular dalam silsilah keputusan, dan memastikan tautan timbal-balik antara ADR lama dan ADR baru tetap konsisten.

Hasil uji menunjukkan seluruh skenario kerusakan integritas referensial tertolak secara deterministik sebelum kode di-merge ke branch utama. Namun penting dicatat: linter hanya memvalidasi kepatuhan format struktural. Kualitas dan kedalaman analisis dari keputusan arsitektur itu sendiri tetap sepenuhnya berada di tangan pertimbangan tim rekayasa.

Kode sumber lengkap, suite pengujian, dan implementasi demo tersedia di repositori:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record

#SoftwareEngineering #Go #SystemDesign #SoftwareArchitecture #CleanArchitecture #EngineeringManagement
