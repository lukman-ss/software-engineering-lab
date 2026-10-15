#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 21: Transactional Outbox Pattern.
Generates 24-slide LinkedIn carousel covering dual-write problem, atomic outbox writes,
asynchronous relay polling, idempotent consumers, concurrency safety, and operational trade-offs.
"""
import os

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/21-outbox-pattern/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/21-outbox-pattern"
COMMIT_SHA = "b2d32740bf482b67be48a216c69e5f8e6b83fe8b"
LAB_NO = "21"
SLUG = "outbox-pattern"
TITLE_ID = "Transactional Outbox Pattern"
SUBTITLE_ID = "Atomic Database Writes Meet Eventual Message Delivery"
TITLE_EN = "Transactional Outbox Pattern"
SUBTITLE_EN = "Atomic Database Writes Meet Eventual Message Delivery"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/21-outbox-pattern/result/image-cover/cover_art.jpg"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en,
        "layout": layout
    })

# =====================================================================
# 24 SLIDE DEFINITIONS
# =====================================================================

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "Diagram transactional outbox: DB commit + relay polling",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Atomic outbox writes + eventual relay delivery guarantee no silent event loss.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Transactional Outbox Pattern adalah solusi klasik untuk dual-write problem dalam arsitektur microservice. Lab ini membuktikan bahwa menulis event ke tabel outbox dalam transaksi database yang sama dengan entitas bisnis menghilangkan risiko kehilangan pesan ketika broker mengalami kegagalan. Pola ini memisahkan penulisan atomik dari pengiriman asinkron melalui relay polling, memastikan konsistensi eventual tanpa bergantung pada transaksi terdistribusi dua-fase.",
    narrative_en="The Transactional Outbox Pattern is a classic solution to the dual-write problem in microservice architectures. This lab proves that writing events to an outbox table within the same database transaction as the business entity eliminates the risk of message loss when the broker fails. The pattern decouples atomic writes from asynchronous delivery via a polling relay, ensuring eventual consistency without relying on two-phase distributed transactions.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "DUAL-WRITE PROBLEM"],
    title="Dual-Write Problem: Ketidakonsistenan Senyap",
    subtitle="Database commit sukses, tetapi broker publish gagal — order ada, event hilang.",
    content_items=[
        {"title": "SILENT EVENT LOSS", "items": [
            "Database commit berhasil menyimpan order.",
            "Broker publish gagal karena jaringan down atau timeout.",
            "Order persisten, tetapi event OrderCreated tidak pernah dikirim ke downstream."
        ], "type": "card"},
        {"title": "COUPLED OPERATIONS", "items": [
            "Business logic langsung memanggil broker.publish().",
            "Hard to test, hard to retry, hard to recover.",
            "Kegagalan broker menggagalkan seluruh flow pembuatan order."
        ], "type": "card"},
        {"title": "MANUAL RECONCILIATION", "items": [
            "Perlu audit log untuk mendeteksi order tanpa event.",
            "Script repair manual untuk replay event yang hilang.",
            "Biaya operasional tinggi dan rentan human error."
        ], "type": "card"}
    ],
    narrative="Dual-write problem muncul ketika aplikasi mencoba menulis ke database dan message broker secara terpisah dalam satu flow bisnis. Transaksi database hanya mengontrol operasi di dalam database itu sendiri; ia tidak dapat melakukan ROLLBACK pada sistem eksternal seperti Kafka atau RabbitMQ. Akibatnya, jika commit database sukses tetapi publish ke broker gagal, sistem masuk ke keadaan inkonsisten: order tersimpan tetapi event tidak pernah dikirim. Kondisi ini sulit dideteksi secara real-time dan memerlukan mekanisme rekonsiliasi manual yang kompleks.",
    narrative_en="The dual-write problem occurs when an application attempts to write to a database and message broker separately within one business flow. Database transactions only control operations within the database itself; they cannot ROLLBACK external systems like Kafka or RabbitMQ. As a result, if database commit succeeds but broker publish fails, the system enters an inconsistent state: the order is persisted but the event is never sent. This condition is hard to detect in real-time and requires complex manual reconciliation mechanisms.",
    mental_model="Separate writes to independent systems cannot be atomic; decouple them with an outbox.",
    source_path="labs/21-outbox-pattern/content/01-content-brief.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "DUA-WRITE NAIEF"],
    title="Mengapa Pendekatan Naif Gagal?",
    subtitle="Kode dual-write terlihat sederhana, tetapi tidak memiliki jaminan atomisitas lintas sistem.",
    content_items=[
        {"title": "NAIVE DUAL-WRITE CODE", "items": [
            "func CreateOrderDualWriteNaive(broker, orderID, amount) error {",
            "    tx := db.BeginTx()",
            "    tx.SaveOrder(order)",
            "    tx.Commit()          // DB berhasil",
            "    broker.Publish(msg)  // Broker gagal → inkonsisten",
            "}"
        ], "type": "code"},
        {"title": "KEGAGALAN YG TERJADI", "items": [
            "Order persisten di DB, tetapi event tidak sampai ke broker.",
            "Tanpa outbox, tidak ada jejak pesan yang tertinggal untuk di-retry.",
            "Reliable delivery must be decoupled dari business transaction."
        ], "type": "bullet"}
    ],
    narrative="Implementasi naive dual-write tampaknya intuitif: simpan order ke database, lalu publish event ke broker. Namun, karena kedua operasi berada di luar boundary transaksi yang sama, kegagalan publish tidak dapat di-rollback dari sisi database. Order tetap tersimpan sementara event hilang selamanya. Pendekatan ini juga memperkenalkan coupling yang kuat antara service pemesanan dan broker, sehingga setiap kegagalan broker akan menggagalkan seluruh flow bisnis.",
    narrative_en="The naive dual-write implementation seems intuitive: save the order to the database, then publish the event to the broker. However, because both operations are outside the same transactional boundary, a publish failure cannot be rolled back from the database side. The order remains persisted while the event is lost forever. This approach also introduces strong coupling between the ordering service and the broker, so any broker failure will abort the entire business flow.",
    mental_model="Naive dual-write lacks cross-system atomicity; the outbox pattern restores it by staging events.",
    source_path="internal/outbox/service.go:55-90",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SCHEMA"],
    title="Mental Model: Transactional Outbox",
    subtitle="Tulis entitas bisnis + event ke tabel outbox dalam satu transaksi database.",
    content_items=[
        {"title": "OUTBOX MESSAGE SCHEMA", "items": [
            "type OutboxMessage struct {",
            "    ID        string",
            "    EventType string",
            "    Payload   string",
            "    Status    MessageStatus  // PENDING, PROCESSED",
            "    CreatedAt time.Time",
            "}"
        ], "type": "code"},
        {"title": "PRINSIP INTI", "items": [
            "Atomic write: order + outbox message dalam satu Tx.",
            "Decoupled relay: polling worker membaca PENDING records.",
            "At-least-once delivery: retry hingga publish sukses.",
            "Idempotent consumer: deduplication via event ID."
        ], "type": "bullet"}
    ],
    narrative="Mental model inti dari pola Transactional Outbox adalah memindahkan penulisan event ke dalam boundary transaksi yang sama dengan penulisan entitas bisnis. Alih-alih melakukan dua write terpisah, aplikasi menulis order dan outbox message secara bersamaan dalam satu transaksi database. Setelah commit berhasil, kedua rekaman tersebut persisten. Proses asynchronous yang disebut message relay kemudian secara periodik polling record PENDING dan mengirimkannya ke broker. Jika relay gagal, record tetap PENDING dan akan dicoba lagi pada siklus berikutnya.",
    narrative_en="The core mental model of the Transactional Outbox pattern is to move the event write into the same transactional boundary as the business entity write. Instead of performing two separate writes, the application writes the order and outbox message simultaneously within one database transaction. After successful commit, both records are persisted. An asynchronous process called the message relay then periodically polls PENDING records and delivers them to the broker. If the relay fails, the record remains PENDING and will be retried in the next cycle.",
    mental_model="Commit business entity and its event together; let a relay eventually deliver.",
    source_path="internal/outbox/model.go • internal/outbox/service.go:17-53",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Alur Eksekusi: Dari Transaksi Hingga Consumer",
    subtitle="Pipeline端到端: Service → DB Tx → Relay → Broker → Consumer.",
    content_items=[
        {"title": "EKSEKUSI PIPELINE", "items": [
            "Client → OrderService.CreateOrderWithOutbox",
            "Service → DB.BeginTx() → SaveOrder + SaveOutbox (staged)",
            "Service → DB.Commit() → flushed atomically",
            "Relay (ticker) → GetPendingOutbox() → [PENDING]",
            "Relay → Broker.Publish(msg) → success?",
            "Relay → MarkOutboxProcessed(msg.ID) → PROCESSED",
            "Consumer.Handle(msg) → idempotent check → accepted"
        ], "type": "diagram"}
    ],
    narrative="Pipeline eksekusi Transactional Outbox terdiri dari lima tahap utama. Pertama, service menerima request pembuatan order dan memulai transaksi database. Kedua, service menyimpan order dan outbox message secara staged di dalam transaksi. Ketiga, service melakukan commit yang secara atomik flush kedua staged record ke database. Keempat, relay worker yang berjalan periodik polling record PENDING dan mencoba publish ke broker. Kelima, consumer menerima pesan dan melakukan idempotent check sebelum memproses lebih lanjut.",
    narrative_en="The Transactional Outbox execution pipeline consists of five main stages. First, the service receives an order creation request and starts a database transaction. Second, the service stages both the order and outbox message within the transaction. Third, the service commits, which atomically flushes both staged records to the database. Fourth, a relay worker running periodically polls PENDING records and attempts to publish them to the broker. Fifth, the consumer receives the message and performs an idempotent check before further processing.",
    mental_model="Atomic commit ensures both entity and event survive together; relay ensures eventual delivery.",
    source_path="internal/outbox/db.go:112-129 • internal/outbox/relay.go:50-67",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="Implementasi Atomic Outbox Write",
    subtitle="Service.go: CreateOrderWithOutbox — staged save + atomic commit.",
    content_items=[
        {"title": "ATOMIC OUTBOX WRITE", "items": [
            "func (s *OrderService) CreateOrderWithOutbox(orderID string, customerID string, amount float64) error {",
            "    tx := s.db.BeginTx()",
            "    order := Order{ID: orderID, CustomerID: customerID, Amount: amount, Status: OrderStatusCreated}",
            "    payload, _ := json.Marshal(order)",
            "    msg := OutboxMessage{ID: fmt.Sprintf(\"evt-%s\", orderID), EventType: \"OrderCreated\", Payload: string(payload), Status: MessageStatusPending, CreatedAt: time.Now()}",
            "    tx.SaveOrder(order)",
            "    tx.SaveOutbox(msg)",
            "    return tx.Commit()",
            "}"
        ], "type": "code"}
    ],
    narrative="Implementasi atomic outbox write berada di method CreateOrderWithOutbox pada service.go. Method ini memulai transaksi baru, membuat entitas order dan outbox message, lalu menyimpan keduanya secara staged menggunakan SaveOrder dan SaveOutbox. Jika Marshalling JSON gagal atau salah satu save operation error, method akan rollback dan mengembalikan error. Hanya setelah kedua staged record tersimpan dengan sukses, method melakukan commit yang akan secara atomik flush semua perubahan ke database.",
    narrative_en="The atomic outbox write implementation resides in the CreateOrderWithOutbox method in service.go. This method starts a new transaction, creates the order entity and outbox message, then stages both using SaveOrder and SaveOutbox. If JSON marshalling fails or any save operation errors, the method rolls back and returns an error. Only after both staged records are successfully saved does the method commit, which atomically flushes all changes to the database.",
    mental_model="Stage both writes, commit once, guarantee both succeed or both fail.",
    source_path="internal/outbox/service.go:17-53",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "BUKTI"],
    title="Failure Path: Dual-Write Naive Menghasilkan Dangling Order",
    subtitle="Order persisten, event hilang — bukti dari TestDualWriteProblem_Failure.",
    content_items=[
        {"title": "DEMO OUTPUT", "items": [
            "=== Lab 21: Transactional Outbox Pattern Demo ===",
            "[Scenario 1: The Dual-Write Problem]",
            "Direct write failed: failed to publish to broker after DB commit: broker unavailable",
            "State Inconsistency: Order in DB = true, Broker Message Count = 0"
        ], "type": "terminal"},
        {"title": "ANALISIS", "items": [
            "TestDualWriteProblem_Failure secara eksplisit mengassert ketidakonsistenan ini.",
            "Order ada di DB, tetapi broker published count = 0.",
            "Inconsistency terjadi karena broker publish dilakukan setelah commit."
        ], "type": "bullet"}
    ],
    narrative="Failure path ini membuktikan bagaimana dual-write naive menciptakan dangling reference: order yang persisten di database tanpa event yang sesuai di broker. Dalam demo, broker disimulasikan gagal melalui SetFailNext(true). Setelah commit database sukses, method CreateOrderDualWriteNaive mencoba publish ke broker yang gagal, menyebabkan error dikembalikan. Namun order sudah tersimpan permanen di database. Test Case TestDualWriteProblem_Failure mengassert kondisi ini: order ditemukan di DB, tetapi broker tidak menerima pesan apa pun.",
    narrative_en="This failure path proves how naive dual-write creates a dangling reference: an order persisted in the database without a corresponding event in the broker. In the demo, the broker is simulated to fail via SetFailNext(true). After successful database commit, the CreateOrderDualWriteNaive method attempts to publish to the failing broker, causing an error to be returned. However, the order is already permanently stored in the database. The TestDualWriteProblem_Failure test case asserts this condition: the order is found in the DB, but the broker received no messages.",
    mental_model="When broker publish follows DB commit, failure leaves dangling orders without events.",
    source_path="tests/outbox_test.go:118-140 • cmd/demo/main.go:20-28",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme vs Gejala: Dual-Write vs Outbox",
    subtitle="Tabel perbandingan perilaku kedua pendekatan saat broker failure.",
    content_items=[
        {"title": "PERBANDINGAN", "items": [
            "Aspek | Dual-Write Naive | Transactional Outbox",
            "Write Boundary | Terpisah (DB commit dulu, lalu broker) | Sama (order + outbox dalam satu Tx)",
            "Atomicity | Tidak atomik lintas sistem | Atomik untuk write pertama",
            "Broker Failure Impact | Order persisten, event hilang | Order + event persisten, relay retry",
            "Idempotency Requirement | Tidak diperlukan (karena konsisten) | Diperlukan (at-least-once delivery)",
            "Operational Complexity | Rendah, tetapi risiko tinggi | Sedang, dengan relay dan cleanup"
        ], "type": "table"}
    ],
    narrative="Tabel perbandingan ini menyoroti perbedaan fundamental antara pendekatan dual-write naive dan transactional outbox. Pada dual-write naive, boundary write terpisah antara database commit dan broker publish, sehingga tidak ada atomisitas lintas sistem. Ketika broker gagal setelah commit, order tetap persisten sedangkan event hilang. Sebaliknya, pada transactional outbox, kedua write dilakukan dalam satu transaksi database, menjamin atomisitas untuk tahap pertama. Broker failure hanya menunda delivery, bukan menghapus event, karena relay akan retry hingga sukses.",
    narrative_en="This comparison table highlights the fundamental differences between naive dual-write and transactional outbox approaches. In naive dual-write, the write boundary is separated between database commit and broker publish, so there is no cross-system atomicity. When the broker fails after commit, the order persists while the event is lost. Conversely, in transactional outbox, both writes are performed within one database transaction, guaranteeing atomicity for the first stage. Broker failure only delays delivery, not deletes the event, because the relay will retry until success.",
    mental_model="Outbox shifts the failure window from permanent loss to eventual delivery with retry.",
    source_path="labs/21-outbox-pattern/content/02-master-draft.md",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "PATTERN"],
    title="3 Keputusan Desain Kritis",
    subtitle="Urutan keputusan arsitektural yang harus diambil saat mengimplementasikan outbox.",
    content_items=[
        {"title": "KEPUTUSAN 1: ATOMIC BOUNDARY", "items": [
            "Tentukan boundary transaksi yang mencakup write business entity + write outbox event.",
            "Gunakan staging pattern: simpan sementara di memory, baru flush saat commit."
        ], "type": "card"},
        {"title": "KEPUTUSAN 2: RELIABLE RELAY", "items": [
            "Pilih strategi polling (Polling Publisher) atau log-tailing (CDC).",
            "Relay harus retry hingga publish sukses dan hanya tandai PROCESSED setelah sukses."
        ], "type": "card"},
        {"title": "KEPUTUSAN 3: IDEMPOTENT CONSUMER", "items": [
            "Consumer harus memiliki deduplication mechanism berbasis event ID.",
            "Handle(msg) mengembalikan true untuk pertama, false untuk duplikat."
        ], "type": "card"}
    ],
    narrative="Tiga keputusan desain kritis ini membentuk fondasi implementasi transactional outbox yang robust. Keputusan pertama menetapkan boundary transaksi yang atomik untuk write bisnis dan outbox. Keputusan kedua memilih mekanisme relay yang reliable dengan retry logic. Keputusan ketiga memastikan consumer dapat menangani duplicate delivery dengan aman melalui idempotency check. Ketiga keputusan ini saling bergantung dan harus diimplementasikan bersama untuk mencapai eventual consistency yang diinginkan.",
    narrative_en="These three critical design decisions form the foundation of a robust transactional outbox implementation. The first decision establishes an atomic transaction boundary for business and outbox writes. The second decision selects a reliable relay mechanism with retry logic. The third decision ensures the consumer can safely handle duplicate delivery through idempotency checks. These three decisions are interdependent and must be implemented together to achieve the desired eventual consistency.",
    mental_model="Atomic boundary → reliable relay → idempotent consumer = eventual consistency.",
    source_path="labs/21-outbox-pattern/content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Terminal Output Demo",
    subtitle="Output langsung dari go run ./cmd/demo yang menunjukkan dual-write failure dan outbox success.",
    content_items=[
        {"title": "DEMO TERMINAL", "items": [
            "=== Lab 21: Transactional Outbox Pattern Demo ===",
            "",
            "[Scenario 1: The Dual-Write Problem]",
            "Direct write failed: failed to publish to broker after DB commit: broker unavailable",
            "State Inconsistency: Order in DB = true, Broker Message Count = 0",
            "",
            "[Scenario 2: Transactional Outbox Solution]",
            "Broker received messages: 1",
            " - Event ID: evt-order-outbox-success, Type: OrderCreated, Payload: {...}",
            "",
            "[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]",
            "Consumer processing duplicate delivery: accepted=false"
        ], "type": "terminal"}
    ],
    narrative="Terminal output di atas adalah bukti eksekusi langsung dari demo command yang menjalankan tiga skenario. Skenario 1 menunjukkan dual-write failure: order ada di DB tetapi broker message count = 0. Skenario 2 menunjukkan outbox success: satu pesan berhasil dipublikasikan ke broker setelah relay polling. Skenario 3 menunjukkan idempotent consumer: duplicate delivery ditolak dengan accepted=false, memastikan received count tetap 1. Output ini memverifikasi bahwa outbox pattern berhasil mengatasi dual-write problem.",
    narrative_en="The terminal output above is direct execution proof from the demo command running three scenarios. Scenario 1 demonstrates dual-write failure: the order exists in the DB but the broker message count equals zero. Scenario 2 demonstrates outbox success: one message was successfully published to the broker after relay polling. Scenario 3 demonstrates the idempotent consumer: duplicate delivery is rejected with accepted=false, ensuring the received count remains one. This output verifies that the outbox pattern successfully overcomes the dual-write problem.",
    mental_model="Observable terminal output validates both failure and success paths.",
    source_path="cmd/demo/main.go • tests/outbox_test.go",
    slide_num="10 / 24",
    layout="terminal-output"
)

# Slide 11: Analisis batasan bukti
add_slide(
    categories=["BATASAN", "ANALISIS"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Batasan lab ini yang penting diketahui sebelum menerapkan ke production.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Lab menggunakan in-memory DB, bukan PostgreSQL/MySQL real.",
            "Hanya Polling Publisher yang diimplementasikan; CDC/log-tailing tidak ada.",
            "Tidak ada dead-letter queue untuk poison payload non-retryable.",
            "Skema OutboxMessage disederhanakan: tidak ada kolom aggregateid/aggregatetype.",
            "Nilai SLA contoh (oldest-event age ~2s normal, 47+ menit abnormal) bersifat ilustratif.",
            "Tidak ada benchmark throughput atau partisi skala tinggi."
        ], "type": "bullet"}
    ],
    narrative="Laboratorium ini memiliki beberapa batasan yang perlu diingat. Pertama, database yang digunakan adalah in-memory mock, bukan RDBMS produksi seperti PostgreSQL atau MySQL. Kedua, hanya strategi Polling Publisher yang diimplementasikan; pendekatan Transaction Log Tailing menggunakan CDC tools seperti Debezium tidak ditunjukkan. Ketiga, tidak ada dead-letter queue untuk menangani payload yang gagal berulang kali. Keempat, skema OutboxMessage disederhanakan tanpa kolom agregat yang biasanya ada di desain produksi. Kelima, nilai SLA yang disebutkan bersifat ilustratif dan harus disesuaikan dengan SLO masing-masing sistem.",
    narrative_en="This laboratory has several limitations that need to be noted. First, the database used is an in-memory mock, not a production RDBMS like PostgreSQL or MySQL. Second, only the Polling Publisher strategy is implemented; the Transaction Log Tailing approach using CDC tools like Debezium is not shown. Third, there is no dead-letter queue to handle payloads that fail repeatedly. Fourth, the OutboxMessage schema is simplified without the aggregate columns typically present in production designs. Fifth, the SLA values mentioned are illustrative and must be adjusted to each system's SLOs.",
    mental_model="Understand lab boundaries to map patterns correctly to production constraints.",
    source_path="labs/21-outbox-pattern/content/01-content-brief.md",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Outbox Pattern",
    subtitle="Konsistensi eventual vs kompleksitas operasional dan latency tambahan.",
    content_items=[
        {"title": "KEUNTUNGAN", "items": [
            "Menghilangkan silent event loss dari dual-write naive.",
            "Menjamin at-least-once delivery melalui retry mechanism.",
            "Decoupling service pemesanan dari broker implementation.",
            "Idempotent consumer mencegah double processing."
        ], "type": "bullet"},
        {"title": "KEMALAHAN", "items": [
            "Latency tambahan: event tidak langsung sampai ke broker, perlu polling interval.",
            "Operational overhead: perlu monitor oldest-event age dan pending count.",
            "Complexity tambahan: relay process, cleanup job, idempotency tracking.",
            "Storage growth: outbox table bertambah hingga di-purge processed records."
        ], "type": "bullet"}
    ],
    narrative="Implementasi Transactional Outbox membawa trade-off yang perlu dipertimbangkan. Di sisi keuntungan, pola ini menghilangkan silent event loss yang sering terjadi pada dual-write naive dan menjamin at-least-once delivery melalui mekanisme retry. Pola ini juga mendekopling service pemesanan dari implementasi broker tertentu. Namun, ada kerugian yang signifikan: latency tambahan karena event tidak langsung sampai ke broker tetapi perlu menunggu polling interval. Selain itu, ada operational overhead berupa monitoring oldest-event age dan pending count, serta kompleksitas tambahan untuk relay process, cleanup job, dan idempotency tracking.",
    narrative_en="Implementing the Transactional Outbox pattern brings trade-offs that need to be considered. On the benefit side, the pattern eliminates silent event loss that often occurs in naive dual-write and guarantees at-least-once delivery through retry mechanisms. It also decouples the ordering service from specific broker implementations. However, there are significant drawbacks: additional latency because events do not reach the broker immediately but must wait for polling intervals. Additionally, there is operational overhead in monitoring oldest-event age and pending counts, as well as added complexity for relay processes, cleanup jobs, and idempotency tracking.",
    mental_model="Trade consistency for operability: accept eventual delivery to prevent permanent loss.",
    source_path="labs/21-outbox-pattern/content/02-master-draft.md",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "VALIDATION"],
    title="Concurrency Safety: 10 Workers × 10 Orders = 100 Messages",
    subtitle="Race detector bersih, consumer idempoten mencegah double processing.",
    content_items=[
        {"title": "TEST CONCURRENCY", "items": [
            "func TestTransactionalOutbox_ConcurrentWrites(t *testing.T) {",
            "    workers := 10",
            "    ordersPerWorker := 10",
            "    totalExpected := workers * ordersPerWorker  // 100",
            "    // Goroutines menulis order + outbox secara konkuren",
            "    // Relay polling dan publish",
            "    // Assert: len(broker.GetPublished()) == 100",
            "    // Assert: len(db.GetPendingOutbox()) == 0",
            "}"
        ], "type": "code"},
        {"title": "HASIL", "items": [
            "Race detector: PASS",
            "Total published: 100",
            "Pending after processing: 0",
            "Consumer concurrent test: 20 unique processed, duplicates rejected."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi membuktikan bahwa implementasi outbox pattern aman terhadap akses multi-goroutine. Dalam TestTransactionalOutbox_ConcurrentWrites, 10 worker goroutines masing-masing menghasilkan 10 order secara simultan, total 100 order. Race detector Go tidak melaporkan data race apapun. Setelah relay selesai memproses semua pending message, broker menerima tepat 100 pesan dan tidak ada yang tersisa di pending state. Test lain, TestTransactionalOutbox_ConcurrentConsumers, memverifikasi bahwa 10 worker yang mengirimkan pesan dengan ID sama hanya menghasilkan 20 unique processed message, membuktikan idempotency kerja dengan baik di bawah konkurensi.",
    narrative_en="Concurrency testing proves that the outbox pattern implementation is safe against multi-goroutine access. In TestTransactionalOutbox_ConcurrentWrites, 10 worker goroutines each generate 10 orders simultaneously, totaling 100 orders. Go's race detector reports no data races. After the relay finishes processing all pending messages, the broker receives exactly 100 messages and none remain in pending state. Another test, TestTransactionalOutbox_ConcurrentConsumers, verifies that 10 workers sending messages with identical IDs result in only 20 unique processed messages, proving idempotency works correctly under concurrency.",
    mental_model="Mutex-protected state and deterministic ID tracking ensure safe concurrent processing.",
    source_path="tests/outbox_test.go:142-194 • tests/outbox_test.go:266-293",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "PERUBAHAN"],
    title="Before vs After: Transformasi Arsitektur",
    subtitle="Perbandingan visual keadaan sebelum dan sesudah menerapkan outbox pattern.",
    content_items=[
        {"title": "BEFORE (Dual-Write Naive)", "items": [
            "Struktur: Service → DB.Commit → Broker.Publish",
            "Keadaan: Order persisten, event mungkin hilang",
            "Kegagalan: Broker down = inconsistent state",
            "Testing: Sulit mengassert consistency lintas sistem",
            "Ops: Perlu script rekonsiliasi manual"
        ], "type": "card"},
        {"title": "AFTER (Transactional Outbox)", "items": [
            "Struktur: Service → DB.Tx(Order+Outbox) → Relay → Broker",
            "Keadaan: Order + event persisten atomik, delivery eventual",
            "Kegagalan: Broker down = retry, tidak ada yang hilang",
            "Testing: Unit test atomis + integration test end-to-end",
            "Ops: Monitor relay health dan oldest-event age"
        ], "type": "card"}
    ],
    narrative="Transformasi dari dual-write naive ke transactional outbox mengubah arsitektur secara fundamental. Sebelum menerapkan pola ini, service menulis ke database dan broker secara berurutan tanpa jaminan atomisitas. Jika broker gagal, order tetap tersimpan tetapi event hilang, menciptakan inconsistent state yang sulit dideteksi. Setelah menerapkan pola ini, kedua write dilakukan dalam satu transaksi database, sehingga baik order maupun event persisten bersama. Broker failure hanya menunda delivery, bukan menghapus event, karena relay akan terus retry hingga sukses.",
    narrative_en="The transformation from naive dual-write to transactional outbox fundamentally changes the architecture. Before implementing this pattern, the service writes to the database and broker sequentially without atomicity guarantees. If the broker fails, the order persists but the event is lost, creating an inconsistent state that is hard to detect. After implementing the pattern, both writes are performed within one database transaction, so both order and event persist together. Broker failure only delays delivery, not deletes the event, because the relay will keep retrying until success.",
    mental_model="Before: coupled and fragile. After: decoupled and resilient.",
    source_path="internal/outbox/service.go • internal/outbox/relay.go",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 Happy Path
add_slide(
    categories=["BUKTI TEST", "HAPPY PATH"],
    title="Test 1: Happy Path - Atomic Write + Relay Dispatch",
    subtitle="TestTransactionalOutbox_HappyPath membuktikan全流程 sukses.",
    content_items=[
        {"title": "TEST CODE", "items": [
            "func TestTransactionalOutbox_HappyPath(t *testing.T) {",
            "    db := outbox.NewDB()",
            "    broker := outbox.NewMockBroker()",
            "    relay := outbox.NewRelay(db, broker, 10*time.Millisecond)",
            "    service := outbox.NewOrderService(db)",
            "    consumer := outbox.NewConsumer()",
            "    relay.Start()",
            "    err := service.CreateOrderWithOutbox(\"o-1\", \"c-1\", 100.50)",
            "    assert.NoError(t, err)",
            "    time.Sleep(50 * time.Millisecond)",
            "    published := broker.GetPublished()",
            "    assert.Len(t, published, 1)",
            "    assert.Equal(t, \"OrderCreated\", published[0].EventType)",
            "    msg, _ := db.GetOutbox(\"evt-o-1\")",
            "    assert.Equal(t, outbox.MessageStatusProcessed, msg.Status)",
            "}"
        ], "type": "code"}
    ],
    narrative="Test ini membuktikan happy path lengkap dari transactional outbox. Pertama, service membuat order dengan ID o-1 secara atomik bersama outbox message. Setelah transaksi commit, relay yang berjalan dengan interval 10ms akan polling dan menemukan message PENDING. Relay kemudian publish ke broker dan menandai message sebagai PROCESSED. Consumer menerima pesan dan memprosesnya. Test ini mengassert bahwa broker menerima tepat satu pesan dengan EventType OrderCreated, dan status outbox message berubah menjadi PROCESSED. Semua assertion lolos, membuktikan alur端到端 berfungsi dengan benar.",
    narrative_en="This test proves the complete happy path of the transactional outbox. First, the service creates an order with ID o-1 atomically along with the outbox message. After transaction commit, the relay running at 10ms intervals will poll and find the PENDING message. The relay then publishes to the broker and marks the message as PROCESSED. The consumer receives the message and processes it. This test asserts that the broker receives exactly one message with EventType OrderCreated, and the outbox message status changes to PROCESSED. All assertions pass, proving the end-to-end flow works correctly.",
    mental_model="Happy path: atomic write → relay poll → broker publish → consumer accept.",
    source_path="tests/outbox_test.go:12-60",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 Rollback
add_slide(
    categories=["BUKTI TEST", "ROLLBACK"],
    title="Test 2: Rollback - Nothing Persists",
    subtitle="TestTransactionalOutbox_Rollback membuktikan staged mutation dibuang sempurna.",
    content_items=[
        {"title": "TEST CODE", "items": [
            "func TestTransactionalOutbox_Rollback(t *testing.T) {",
            "    db := outbox.NewDB()",
            "    tx := db.BeginTx()",
            "    tx.SaveOrder(outbox.Order{ID: \"o-2\", Amount: 50})",
            "    tx.SaveOutbox(outbox.OutboxMessage{ID: \"evt-o-2\", Status: outbox.MessageStatusPending})",
            "    tx.Rollback()",
            "    _, ok := db.GetOrder(\"o-2\")",
            "    assert.False(t, ok, \"order should not exist\")",
            "    _, ok = db.GetOutbox(\"evt-o-2\")",
            "    assert.False(t, ok, \"outbox should not exist\")",
            "    time.Sleep(30 * time.Millisecond)",
            "    assert.Len(t, broker.GetPublished(), 0, \"no messages to broker\")",
            "}"
        ], "type": "code"}
    ],
    narrative="Test rollback membuktikan bahwa staged mutation dibuang sepenuhnya ketika transaksi di-rollback. Method SaveOrder dan SaveOutbox hanya menaruh data di map stagedOrders dan stagedOutbox, belum menulis ke database utama. Ketika Rollback dipanggil, kedua map staged tersebut dibuang tanpa pernah menyentuh database. Test mengassert bahwa order dan outbox message tidak dapat diambil dari database setelah rollback. Selain itu, setelah menunggu relay poll, broker tidak menerima pesan apa pun, membuktikan bahwa rollback mencegah event dispatch sama sekali.",
    narrative_en="The rollback test proves that staged mutations are completely discarded when a transaction is rolled back. The SaveOrder and SaveOutbox methods only place data in the stagedOrders and stagedOutbox maps, without writing to the main database. When Rollback is called, both staged maps are discarded without ever touching the database. The test asserts that the order and outbox message cannot be retrieved from the database after rollback. Additionally, after waiting for relay polling, the broker receives no messages, proving that rollback prevents event dispatch entirely.",
    mental_model="Rollback discards staged writes; no persistence, no dispatch.",
    source_path="tests/outbox_test.go:62-91",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Idempotency
add_slide(
    categories=["BUKTI TEST", "IDEMPOTENSI"],
    title="Test 3: Idempotency - Duplicate Delivery Rejected",
    subtitle="TestTransactionalOutbox_Idempotency_DuplicateDelivery membuktikan consumer safe.",
    content_items=[
        {"title": "TEST CODE", "items": [
            "func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {",
            "    consumer := outbox.NewConsumer()",
            "    msg := outbox.OutboxMessage{ID: \"evt-duplicate\", EventType: \"OrderCreated\", Payload: \"data\"}",
            "    p1 := consumer.Handle(msg)",
            "    assert.True(t, p1, \"first delivery should be accepted\")",
            "    p2 := consumer.Handle(msg)",
            "    assert.False(t, p2, \"duplicate delivery should be rejected\")",
            "    assert.Equal(t, 1, consumer.GetReceivedCount(), \"count stays 1\")",
            "}"
        ], "type": "code"}
    ],
    narrative="Test idempotensi membuktikan bahwa consumer dapat dengan aman menangani duplicate delivery. Pertama, Handle dipanggil dengan pesan unik dan mengembalikan true (accepted). Kedua, Handle dipanggil lagi dengan pesan yang sama dan mengembalikan false (rejected). Test mengassert bahwa received count tetap 1, membuktikan bahwa duplicate tidak menambah state konsumen. Skenario ini mensimulasikan kasus di mana relay crash setelah publish tetapi sebelum menandai message sebagai PROCESSED, menyebabkan pesan yang sama dipublish ulang oleh relay pada siklus berikutnya.",
    narrative_en="The idempotency test proves that the consumer can safely handle duplicate delivery. First, Handle is called with a unique message and returns true (accepted). Second, Handle is called again with the same message and returns false (rejected). The test asserts that the received count remains 1, proving that duplicates do not add to the consumer's state. This scenario simulates a case where the relay crashes after publish but before marking the message as PROCESSED, causing the same message to be republished by the relay in the next cycle.",
    mental_model="Idempotent consumers deduplicate by event ID; second delivery is ignored.",
    source_path="tests/outbox_test.go:93-116",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Analisis mendalam - Test 4 Relay Retry
add_slide(
    categories=["BUKTI TEST", "RELAY RETRY"],
    title="Test 4: Relay Retry - Broker Failure Tolerated",
    subtitle="TestTransactionalOutbox_RelayRetryAfterBrokerFailure membuktikan retry logic.",
    content_items=[
        {"title": "TEST CODE", "items": [
            "func TestTransactionalOutbox_RelayRetryAfterBrokerFailure(t *testing.T) {",
            "    db := outbox.NewDB()",
            "    broker := outbox.NewMockBroker()",
            "    relay := outbox.NewRelay(db, broker, 10*time.Millisecond)",
            "    service := outbox.NewOrderService(db)",
            "    broker.SetFailNext(true)  // first publish fails",
            "    service.CreateOrderWithOutbox(\"o-retry\", \"c-retry\", 55.0)",
            "    relay.Start()",
            "    time.Sleep(25 * time.Millisecond)  // first poll fails",
            "    msg, _ := db.GetOutbox(\"evt-o-retry\")",
            "    assert.Equal(t, outbox.MessageStatusPending, msg.Status)",
            "    time.Sleep(50 * time.Millisecond)  // second poll succeeds",
            "    published := broker.GetPublished()",
            "    assert.Len(t, published, 1)",
            "    msg, _ = db.GetOutbox(\"evt-o-retry\")",
            "    assert.Equal(t, outbox.MessageStatusProcessed, msg.Status)",
            "}"
        ], "type": "code"}
    ],
    narrative="Test relay retry membuktikan mekanisme ketahanan terhadap kegagalan broker. Pertama, broker diatur untuk gagal pada publish berikutnya melalui SetFailNext(true). Service membuat order dengan outbox, kemudian relay dimulai. Pada poll pertama, publish gagal karena broker dalam mode fail, sehingga message tetap berstatus PENDING. Pada poll berikutnya (setelah interval ticker), broker kembali normal, publish sukses, dan message ditandai PROCESSED. Test ini mengassert bahwa setelah retry, broker menerima pesan dan status outbox berubah menjadi PROCESSED.",
    narrative_en="The relay retry test proves the failure tolerance mechanism against broker failures. First, the broker is set to fail on the next publish via SetFailNext(true). The service creates an order with outbox, then the relay starts. On the first poll, the publish fails because the broker is in fail mode, so the message remains PENDING. On the next poll (after the ticker interval), the broker returns to normal, the publish succeeds, and the message is marked PROCESSED. This test asserts that after retry, the broker receives the message and the outbox status changes to PROCESSED.",
    mental_model="Relay retry converts transient broker failures into eventual delivery.",
    source_path="tests/outbox_test.go:219-264",
    slide_num="18 / 24",
    layout="test-analysis"
)

# Slide 19: Analisis mendalam - Test 5 Purge Processed
add_slide(
    categories=["BUKTI TEST", "CLEANUP"],
    title="Test 5: Purge Processed - Memory Leak Prevention",
    subtitle="TestTransactionalOutbox_PurgeProcessed membuktikan cleanup mechanism bekerja.",
    content_items=[
        {"title": "TEST CODE", "items": [
            "func TestTransactionalOutbox_PurgeProcessed(t *testing.T) {",
            "    db := outbox.NewDB()",
            "    tx := db.BeginTx()",
            "    tx.SaveOutbox(outbox.OutboxMessage{ID: \"m1\", Status: outbox.MessageStatusPending})",
            "    tx.SaveOutbox(outbox.OutboxMessage{ID: \"m2\", Status: outbox.MessageStatusProcessed})",
            "    tx.Commit()",
            "    purged := db.PurgeProcessedOutbox()",
            "    assert.Equal(t, 1, purged)",
            "    _, ok := db.GetOutbox(\"m1\")",
            "    assert.True(t, ok, \"pending should remain\")",
            "    _, ok = db.GetOutbox(\"m2\")",
            "    assert.False(t, ok, \"processed should be purged\")",
            "}"
        ], "type": "code"}
    ],
    narrative="Test purge processed membuktikan mekanisme pembersihan tabel outbox untuk mencegah memory leak. Database awalnya berisi dua message: m1 dengan status PENDING dan m2 dengan status PROCESSED. Method PurgeProcessedOutbox akan menghapus semua record dengan status PROCESSED dan mengembalikan jumlah record yang dihapus. Test mengassert bahwa purged count = 1, message PENDING masih ada, dan message PROCESSED sudah dihapus. Tanpa cleanup mechanism seperti ini, tabel outbox akan tumbuh tak terbatas karena record PROCESSED tidak pernah dihapus.",
    narrative_en="The purge processed test proves the outbox table cleanup mechanism to prevent memory leaks. The database initially contains two messages: m1 with PENDING status and m2 with PROCESSED status. The PurgeProcessedOutbox method will delete all records with PROCESSED status and return the count of deleted records. The test asserts that purged count equals 1, the PENDING message still exists, and the PROCESSED message has been deleted. Without such a cleanup mechanism, the outbox table would grow indefinitely because PROCESSED records are never deleted.",
    mental_model="Purge processed records to prevent unbounded outbox table growth.",
    source_path="tests/outbox_test.go:196-217",
    slide_num="19 / 24",
    layout="test-analysis"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Checklist Implementasi Outbox Pattern",
    subtitle="Langkah-langkah wajib saat menerapkan pola ini ke production.",
    content_items=[
        {"title": "IMPLEMENTATION CHECKLIST", "items": [
            "□ Buat tabel outbox_events dengan kolom: id, aggregate_type, aggregate_id, type, payload, status, created_at",
            "□ Pastikan business transaction menyertakan INSERT ke outbox_events",
            "□ Implementasikan polling relay dengan interval yang sesuai (misal 1-5 detik)",
            "□ Relay hanya menandai PROCESSED setelah broker ACK",
            "□ Consumer mengimplementasikan idempotency check berdasarkan event ID",
            "□ Tambahkan scheduled job untuk purge processed records",
            "□ Monitor oldest-event age dan pending count untuk deteksi degradasi",
            "□ Siapkan dead-letter queue untuk poison payload (opsional)"
        ], "type": "bullet"}
    ],
    narrative="Checklist implementasi ini mencakup langkah-langkah esensial untuk menerapkan transactional outbox pattern di production. Pertama, buat tabel outbox_events dengan skema yang memadai termasuk kolom agregat untuk partitioning Kafka. Kedua, pastikan setiap business transaction menyertakan INSERT ke outbox_events. Ketiga, implementasikan polling relay dengan interval yang seimbang antara latensi dan beban database. Keempat, relay hanya menandai message sebagai PROCESSED setelah broker memberikan ACK. Kelima, consumer harus memiliki idempotency check berbasis event ID. Keenam, tambahkan scheduled job untuk purge processed records guna mencegah table bloat.",
    narrative_en="This implementation checklist covers essential steps for applying the transactional outbox pattern in production. First, create an outbox_events table with an adequate schema including aggregate columns for Kafka partitioning. Second, ensure every business transaction includes an INSERT to outbox_events. Third, implement a polling relay with an interval balanced between latency and database load. Fourth, the relay only marks messages as PROCESSED after the broker provides ACK. Fifth, the consumer must have an idempotency check based on event ID. Sixth, add a scheduled job to purge processed records to prevent table bloat.",
    mental_model="Implementation requires atomic writes, reliable relay, idempotent consumers, and cleanup.",
    source_path="labs/21-outbox-pattern/content/05-key-takeaways.md",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Apa yang Harus Dihindari",
    subtitle="Kesalahan umum yang membuat outbox pattern menjadi tidak efektif.",
    content_items=[
        {"title": "ANTI-PATTERN 1", "items": [
            "Menulis ke broker sebelum commit database.",
            "Menggunakan ID event yang tidak unik (misal timestamp saja).",
            "Tidak memiliki mekanisme retry pada relay."
        ], "type": "card"},
        {"title": "ANTI-PATTERN 2", "items": [
            "Consumer tanpa idempotency check.",
            "Tidak ada purge untuk processed records.",
            "Relay polling interval terlalu lama (>1 menit) sehingga latency tinggi."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI GAGAL", "items": [
            "Jalankan TestDualWriteProblem_Failure untuk melihat inconsistency.",
            "Matikan broker selama relay berjalan, perhatikan pending accumulation.",
            "Hapus idempotency check dari consumer, maka duplicate akan diproses ganda."
        ], "type": "card"}
    ],
    narrative="Anti-pattern dalam implementasi outbox pattern dapat merusak seluruh jaminan konsistensi yang dimaksudkan. Menulis ke broker sebelum commit database mengembalikan ke keadaan dual-write naive. Menggunakan ID event yang tidak unik menyebabkan deduplication gagal. Consumer tanpa idempotency check akan memproses duplicate message secara ganda. Tidak ada purge untuk processed records akan menyebabkan tabel outbox tumbuh tak terbatas. Relly dengan polling interval terlalu lama akan meningkatkan latency delivery secara signifikan.",
    narrative_en="Anti-patterns in outbox pattern implementation can undermine all the consistency guarantees intended. Writing to the broker before database commit reverts to naive dual-write. Using non-unique event IDs causes deduplication to fail. Consumers without idempotency checks will process duplicate messages twice. No purge for processed records will cause the outbox table to grow indefinitely. A relay with too-long polling intervals will significantly increase delivery latency.",
    mental_model="Avoid coupling, non-unique IDs, and missing idempotency; they undo the pattern's benefits.",
    source_path="internal/outbox/service.go:55-90 • tests/outbox_test.go:118-140",
    slide_num="21 / 24",
    layout="anti-pattern-cards"
)

# Slide 22: Takeaway praktis & command block
add_slide(
    categories=["PRAKTIS", "COMMAND"],
    title="Takeaway Praktis & Command Block",
    subtitle="Cara mereproduksi lab dan menjalankan verifikasi.",
    content_items=[
        {"title": "COMMAND BLOCK", "items": [
            "cd labs/21-outbox-pattern",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "PRAKTIK TERBAIK", "items": [
            "Gunakan event ID yang stabil (UUID berbasis aggregate ID).",
            "Monitor oldest-event age sebagai sinyal degradasi relay/broker.",
            "Pertahankan pending count di bawah threshold tertentu.",
            "Implementasikan dead-letter queue untuk payload yang gagal berulang."
        ], "type": "bullet"}
    ],
    narrative="Takeaway praktis dari lab ini adalah pemahaman tentang cara mereproduksi dan memverifikasi implementasi outbox pattern. Command go test -v ./... menjalankan seluruh suite test dengan verbose output. Command go test -race ./... menjalankan race detector untuk memastikan thread safety. Command go run ./cmd/demo menjalankan demo interaktif yang menampilkan tiga skenario. Praktik terbaik meliputi penggunaan event ID yang stabil berdasarkan aggregate ID, monitoring oldest-event age sebagai sinyal degradasi, mempertahankan pending count di bawah threshold, dan mengimplementasikan dead-letter queue untuk payload yang gagal berulang.",
    narrative_en="The practical takeaway from this lab is understanding how to reproduce and verify the outbox pattern implementation. The go test -v ./... command runs the entire test suite with verbose output. The go test -race ./... command runs the race detector to ensure thread safety. The go run ./cmd/demo command runs the interactive demo displaying three scenarios. Best practices include using stable event IDs based on aggregate IDs, monitoring oldest-event age as a degradation signal, maintaining pending counts below thresholds, and implementing dead-letter queues for payloads that fail repeatedly.",
    mental_model="Reproduce, measure, monitor — then trust eventual consistency.",
    source_path="README.md • cmd/demo/main.go",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "JARAK PANJANG"],
    title="Prinsip Arsitektural Jangka Panjang",
    subtitle="Filosofi desain yang melampaui pola outbox spesifik.",
    content_items=[
        {"title": "PRINSIP 1", "items": [
            "Eventual consistency lebih baik daripada strong consistency yang mahal.",
            "Decouple write boundary dari delivery boundary.",
            "At-least-once delivery lebih aman daripada at-most-once."
        ], "type": "bullet"},
        {"title": "PRINSIP 2", "items": [
            "Idempotency adalah keharusan, bukan optional.",
            "Observable metrics (oldest-event age, pending count) wajib diimplementasikan.",
            "Cleanup mechanism harus berjalan scheduled, bukan manual."
        ], "type": "bullet"}
    ],
    narrative="Prinsip arsitektural jangka panjang ini melampaui implementasi spesifik transactional outbox. Pertama, eventual consistency seringkali lebih praktis daripada strong consistency yang memerlukan koordinasi terdistribusi mahal. Kedua, memisahkan boundary write dari boundary delivery memungkinkan masing-masing bagian dioptimalkan secara independen. Ketiga, at-least-once delivery lebih aman daripada at-most-once karena mencegah kehilangan data. Keempat, idempotency bukan fitur optional melainkan keharusan karena retry adalah hal yang wajar. Kelima, observable metrics seperti oldest-event age dan pending count harus diimplementasikan untuk deteksi dini degradasi.",
    narrative_en="These long-term architectural principles extend beyond the specific transactional outbox implementation. First, eventual consistency is often more practical than strong consistency that requires expensive distributed coordination. Second, separating the write boundary from the delivery boundary allows each part to be optimized independently. Third, at-least-once delivery is safer than at-most-once because it prevents data loss. Fourth, idempotency is not an optional feature but a necessity because retries are normal. Fifth, observable metrics like oldest-event age and pending counts must be implemented for early degradation detection.",
    mental_model="Design for failure, measure relentlessly, accept eventual consistency.",
    source_path="labs/21-outbox-pattern/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Konsistensi tercapai, event tak pernah hilang.",
    subtitle="",
    content_items=[
        {"title": "SUMBER LAB", "items": [
            "internal/outbox/db.go",
            "internal/outbox/service.go",
            "internal/outbox/relay.go",
            "internal/outbox/consumer.go",
            "internal/outbox/broker.go",
            "tests/outbox_test.go",
            "cmd/demo/main.go"
        ], "type": "bullet"},
        {"title": "VERIFIKASI", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing: 8/8 test cases",
            "Race detector: clean",
            "Demo: all 3 scenarios pass"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium Transactional Outbox Pattern ini membuktikan bahwa dual-write problem dapat diatasi dengan menulis event ke tabel outbox dalam transaksi database yang sama dengan entitas bisnis. Relay polling yang asynchronous memastikan event akhirnya sampai ke broker, sementara consumer idempoten mencegah pemrosesan ganda. Lab ini menggunakan Go standard library, tidak bergantung pada framework eksternal, dan dapat direproduksi dengan mudah melalui perintah go test dan go run demo. Seluruh kode sumber dan test tersedia di repository GitHub yang tercantum.",
    narrative_en="This Transactional Outbox Pattern laboratory proves that the dual-write problem can be solved by writing events to an outbox table within the same database transaction as the business entity. The asynchronous polling relay ensures events eventually reach the broker, while the idempotent consumer prevents double processing. This lab uses Go standard library, has no external framework dependencies, and can be easily reproduced through go test and go run demo commands. All source code and tests are available in the GitHub repository listed.",
    mental_model="Atomic outbox writes + eventual relay delivery = no silent event loss.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
