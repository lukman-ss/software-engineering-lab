# API Compatibility

## Versioning API
Evidence: Versi dipasang saat client pertama kali request. Kompatibilitas dijaga melalui interceptor atau version-change module secara transparan tanpa merusak kontrak.
Source: Stripe API Versioning
URL: https://stripe.com/blog/api-versioning
Confidence: HIGH

## Additive vs Destructive
- **Additive**: Menambah field baru, endpoint baru (non-breaking).
- **Destructive (Breaking)**: Mengubah tipe data, mengubah struktur response. Di Stripe, destructive changes diatur dalam transform layer per versi sehingga business logic selalu memakai object format terbaru.
