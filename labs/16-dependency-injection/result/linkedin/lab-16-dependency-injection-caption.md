Seringkali kita melihat kode di mana class/struct membuat dependensinya sendiri secara langsung di dalam method. Tampak sepele, tapi ini mengunci implementasi konkret ke dalam logika bisnis. Akibatnya: pengujian unit terisolasi mustahil dilakukan tanpa infrastruktur nyata, dan pergantian vendor/provider memaksa refactoring masif di puluhan file.

Lab 16 membedah prinsip inti Dependency Injection (DI) dan Inversion of Control (IoC) lewat Go:
- Memisahkan fase *construction* dari *use* via Constructor Injection.
- Membuktikan Service Locator sebagai anti-pattern yang menyembunyikan dependensi.
- Memvalidasi mocking tanpa network overhead (test PASS 0.00s, race-clean).
- Batasan tegas: DI mempermudah pengujian, bukan jaminan otomatis untuk performa atau reliability production.

Code & Lab: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection
Snapshot: commit 8772235

#SoftwareEngineering #GoLang #DesignPatterns #DependencyInjection #CleanArchitecture