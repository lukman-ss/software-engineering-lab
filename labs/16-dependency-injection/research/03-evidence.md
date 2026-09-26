# Evidence Repository — DI & IoC

## Evidence 1 — Definisi DI (Spring)
Claim: DI = objek hanya mendefinisikan dependensi lewat konstruktor, argumen factory, atau properti; container menginjeksi saat membuat bean.
Source: Spring Framework 7.0.9
URL: https://docs.spring.io/spring-framework/reference/core/beans/dependencies/factory-collaborators.html
Confidence: HIGH
Corroborated By: Fowler 2004, Microsoft .NET DI docs
Notes: Tiga sumber otoritatif konvergen pada definisi ini.

## Evidence 2 — Tiga bentuk DI (Fowler)
Claim: Constructor Injection, Setter Injection, Interface Injection.
Source: Martin Fowler, 23 Jan 2004
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: Spring mendokumentasikan constructor-based dan setter-based sebagai dua varian utama.
Notes: Interface injection jarang dipakai langsung di framework modern.

## Evidence 3 — DI vs Service Locator
Claim: Service Locator membuat setiap pengguna service bergantung pada locator; DI tidak. DI lebih baik untuk komponen yang dipakai aplikasi eksternal.
Source: Martin Fowler 2004
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: PSR-11 — container sebagai Service Locator generally discouraged.
Notes: Prinsip terpenting: separation of configuration from use.

## Evidence 4 — DI meningkatkan testability
Claim: DI memudahkan penggantian implementasi asli dengan stub/mock tanpa mengubah class yang diuji.
Source: Spring docs + Fowler 2004
URL: https://docs.spring.io/spring-framework/reference/core/beans/dependencies/factory-collaborators.html
Confidence: HIGH
Corroborated By: Microsoft .NET — mock/stub tidak mungkin dengan hard-coded dependency.
Notes: Fowler mencatat Service Locator yang dirancang baik juga bisa di-stub.

## Evidence 5 — Constructor vs Setter
Claim: Prefer constructor injection; fallback ke setter bila banyak parameter, nilai primitif/string butuh nama, atau hierarki inheritance kompleks.
Source: Martin Fowler 2004
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: .NET memakai constructor sebagai pendekatan utama; Spring mendukung keduanya.
Notes: Bahasa modern dengan named args (PHP 8, C#, Kotlin) sebagian memitigasi kelemahan konstruktor.

## Evidence 6 — Service lifetimes
Claim: Singleton (satu instance seumur app), Scoped (satu instance per request/lifecycle), Transient/Prototype (instance baru per resolve).
Source: Microsoft .NET (Jan 2026), Laravel 12.x
URL: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection/overview
Confidence: HIGH
Corroborated By: Spring bean scopes (singleton, prototype, request, session).
Notes: Terminologi beda antar ekosistem; konsep sama.

## Evidence 7 — PSR-11
Claim: ContainerInterface hanya punya get dan has. get pada id tak dikenal WAJIB throw NotFoundExceptionInterface. Pengguna SHOULD NOT passing container ke objek (Service Locator discouraged).
Source: PHP-FIG PSR-11
URL: https://www.php-fig.org/psr/psr-11/
Confidence: HIGH
Notes: Standar ratified, bukan opini.

## Evidence 8 — IoC prinsip umum vs DI spesifik
Claim: IoC = framework memanggil kode pengguna (event handler, template method, EJB lifecycle). DI = satu bentuk IoC spesifik untuk perakitan service.
Source: Martin Fowler, 26 Jun 2005
URL: https://martinfowler.com/bliki/InversionOfControl.html
Confidence: HIGH
Notes: Istilah IoC Container agak misnomer; etimologi Johnson-Foote 1988, Hollywood Principle.

## Evidence 9 — Laravel container
Claim: Zero-configuration resolution via reflection; binding interface-ke-implementasi; contextual binding; tagging; singleton/scoped.
Source: Laravel 12.x docs
URL: https://laravel.com/docs/12.x/container
Confidence: HIGH
Notes: PSR-11 compliant.

## Evidence 10 — Kapan TIDAK perlu DI
Claim: Value object sederhana tanpa dependensi eksternal (DateTime, Money, Address) aman dibuat langsung. DI untuk Database, HTTP client, gateway, cache, queue, email, storage.
Source: Spesifikasi topik lab (bukan sumber primer)
Confidence: MEDIUM
Corroborated By: Tidak terverifikasi independen — Fowler hanya bicara separation of configuration from use.
Notes: Tandai sebagai klaim spesifikasi, bukan fakta terverifikasi.
