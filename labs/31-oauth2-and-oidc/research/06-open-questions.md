# Open Questions

- **RQ1**: Apakah OIDC Core 1.0 akan diperbarui agar secara eksplisit menghapus Implicit Flow seperti RFC 9700/OAuth 2.1? Belum ada draft update terbitan saat penelitian. Next: pantau OpenID Foundation / IETF OAuth WG drafts.
- **RQ2**: Bagaimana DPoP (RFC 9449) berinteraksi dengan BFF pattern? BFF sudah menjadi confidential client; DPoP lebih relevan untuk public client SPA murni. Next: bandingkan threat model DPoP vs BFF.
- **RQ3**: Klaim `azp` (authorized party) pada ID Token — kapan wajib vs opsional, dan bagaimana validasinya di multi-tenant? OIDC Core menyebut hanya saat extension digunakan; butuh contoh implementasi.
- **RQ4**: Token binding vs sender-constrained (mTLS vs DPoP) untuk refresh token — standar mana yang lebih diprioritaskan menurut RFC 9700? RFC 9700 memberikan pilihan; belum ada peringkat eksplisit antar keduanya.
- **RQ5**: Exact redirect URI matching vs wildcard — bagaimana mengimplementasikan secara aman untuk multi-env (dev/staging/prod)? RFC 8252 native apps mengizinkan localhost variasi port; web apps wajib exact match.
- **RQ6**: Apakah JWT `none` algorithm pernah disalahpakai di produksi dunia nyata? RFC 8725 menyebut serangan Mclean (RS256->HS256) dan CVE-2015-9235; butuh data publikasi insiden konkret.
- **RQ7**: Lab "menyimpan access token di localStorage" — apakah ada runtime (e.g., Node backend, mobile) di mana penyimpanan safe berbeda? Bukti di atas fokus browser; perlu perluasan untuk native/mobile (Keychain/Keystore).
- **RQ8**: OIDC Core mengizinkan ID Token dienkripsi (JWE). Bagaimana memvalidasi encrypted ID Token dalam BFF pattern? OIDC Core Sec 3.1.3.7 step 1: decrypt dengan keys negotiated at registration; butuh contoh implementasi.
- **RQ9**: Berapa `auth_time` drift yang dapat diterima? OIDC Core tidak menentukan; butuh panduan operasional.
- **RQ10**: Bagaimana PKCE `code_verifier` harus disimpan selama auth flow (single-page app vs BFF)? Browser-based-apps draft Sec 8 membahas in-memory vs persistent; perlu detail per pola arsitektur.
