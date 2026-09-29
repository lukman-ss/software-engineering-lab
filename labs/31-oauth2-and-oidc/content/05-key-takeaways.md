# Key Takeaways

1. **OAuth 2.0 vs OIDC**: OAuth 2.0 menyediakan authorization (access token untuk akses sumber daya), sedangkan OIDC menambahkan authentication (ID Token JWT untuk identitas pengguna). Keduanya adalah mekanisme terpisah yang saling melengkapi.

2. **PKCE Wajib untuk Public Clients**: RFC 9700 mewajibkan penggunaan PKCE (S256) untuk semua client publik (browser-based, mobile). PKCE mencegah serangan interception authorization code dengan menambahkan proof-of-possession melalui code_verifier.

3. **ID Token Validation Multi-Langkah**: ID Token harus melalui validasi ketat meliputi: signature HMAC, issuer exact match, audience match, expiration check, issued-at window (±5 menit), dan nonce anti-replay. Lewat satu pun step berarti token ditolak.

4. **Refresh Token Rotation**: Setiap refresh token usage harus menghasilkan token baru. Token lama di-invalidate secara instan. Ini mengurangi window of exposure jika token dicuri.

5. **Family Revocation on Replay**: Jika refresh token yang sudah di-revoke digunakan kembali (replay), seluruh token family (termasuk access token aktif) dicabut. Client dipaksa melakukan re-authentication.

6. **Implicit Flow Deprecated**: Implicit flow (response_type=token/id_token) mengekspos token di URL fragment dan rentan XSS. RFC 9700 dan OAuth 2.1 deprecated flow ini. Gunakan Authorization Code Flow + PKCE sebagai gantinya.

7. **Token Storage Matters**: Bearer token tidak boleh disimpan di localStorage pada browser karena rentan terhadap XSS. Gunakan httpOnly Secure SameSite cookie (BFF pattern) atau in-memory storage untuk SPA.

8. **Kedua Token Berbeda Peran**: Access token digunakan untuk mengakses resource server (scope-based). ID Token digunakan oleh client untuk menampilkan identitas user (claim-based). Jangan menggunakan access token sebagai pengganti ID Token.

9. **Concurrent Safety**: Implementasi server menggunakan mutex untuk melindungi state shared. Race detector (`go test -race`) lulus tanpa race condition terdeteksi.

10. **Standar Sebagai Dasar**: Semua mekanisme dilab bersumber dari RFC 6749, RFC 7636, RFC 7519, RFC 8725, RFC 9700, dan OIDC Core 1.0. Implementasi ini mengikuti best practice terkini (BCP 2025).
