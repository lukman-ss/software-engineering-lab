const fs = require('fs');
let content = fs.readFileSync('labs/03-database-transaction/README.md', 'utf8');

// 1. Kapan Butuh Consistency Langsung?
// The prompt asks to find the part implicitly/explicitly stating Payment+ledger+OPL marking -> harus satu transaksi.
// Since it's not present exactly, I'll add a new section under Business Invariant to clarify this.
const replacement1 = `## 3. Business Invariant: Aturan yang Harus Selalu Benar

**Definisi:**

> **Business invariant** adalah aturan bisnis yang harus selalu benar meskipun terjadi failure.

### Kapan Membutuhkan Strong Consistency?

Strong/local atomic consistency diperlukan ketika beberapa perubahan merupakan satu business invariant, berada pada transactional datastore yang sama, dan partial state tidak boleh terlihat.

**Contoh:**
Jika \`payment\` + \`journal\` + \`OPL status\` berada pada database transactional yang sama dan merupakan satu business invariant, maka gunakan **satu local database transaction**.

Tetapi jika:
- \`payment-service DB\` + \`accounting-service DB\`
- atau \`local DB\` + \`external inventory service\`

Maka satu local DB transaction tidak dapat memberikan atomicity lintas boundary. 

> **Penting**: Jenis data tidak menentukan transaction strategy. Business invariant + transaction boundary yang menentukan.`;

content = content.replace(/## 3\. Business Invariant: Aturan yang Harus Selalu Benar\n\n\*\*Definisi:\*\*\n\n> \*\*Business invariant\*\* adalah aturan bisnis yang harus selalu benar meskipun terjadi failure\.\n\n### Contoh Business Invariant\n\n\*\*Pembayaran:\*\*\n\nJika payment tercatat PAID, maka payment record dan financial journal yang wajib menjadi bagian dari transaksi tidak boleh berada pada state setengah jadi\.\n\n```\nBusiness invariant:\n- Payment PAID → Journal harus balance\n- Payment FAILED → Tidak ada journal\n- State setengah jadi → TIDAK BOLEH\n```/g, replacement1);


// 2. Retry
// Checking for retry mention. It seems there's no exact match of "retry menjadi mekanisme wajib", but I will replace the relevant row in the Failure Matrix or add a new section about Retry.
// Wait, the prompt says "Cari wording seperti: retry menjadi mekanisme wajib". I didn't see it in the read, but I'll add a section right before or after Eventual Consistency.
// Let's add a section for Retry explicitly.

const replacement2 = `Retry merupakan mekanisme umum untuk menangani transient failure pada distributed workflow, tetapi tidak semua failure boleh di-retry.

**Retryable:**
- timeout
- temporary network failure
- HTTP 429
- HTTP 502/503/504
- temporary broker unavailable

**Umumnya tidak retryable (tanpa perubahan input/state):**
- validation error
- malformed payload
- unauthorized / forbidden
- business rule rejection
- resource not found yang memang permanent

> Retry policy harus membedakan transient failure dan permanent failure.`;

// Let's just append this to the Glossary or Outbox/Consumer section if "retry menjadi mekanisme wajib" is not found.
// Wait, is "retry menjadi mekanisme wajib" in the Outbox section? No, outbox section doesn't have it.
// Let's search the file for "wajib" or "retry".

