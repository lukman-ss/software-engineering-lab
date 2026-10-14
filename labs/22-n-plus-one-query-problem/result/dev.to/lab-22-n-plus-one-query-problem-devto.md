---
title: "N+1 Query Problem: Batching Database Relations to Eliminate Hidden Latency"
published: false
description: "Explore N+1 Query Problem in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# N+1 Query Problem: Batching Database Relations to Eliminate Hidden Latency

A production system must remain resilient under load. Lab 22 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> N+1 is invisible until it hurts. Batch related queries, don't loop them.

## Satu Query Awal, Puluhan Query Mengikuti

Perilaku cepat di local development menipu engineer saat masuk ke beban production. Setiap query individu terlihat sub-milidetik, sehingga lolos dari slow query logs. Degradasi akumuler: request yang seharusnya 20ms membengkak menjadi beberapa detik. Store query counter mencatat 1 query author + 3 query post = 4 query total pada 3 data. Eager loading memangkas eksekusi menjadi 1 query author + 1 query post = 2 query total. Lab memodelkan in-memory Store; pada RDBMS nyata overhead socket TCP memperparah bottleneck.

## Arsitektur Tiga Layer Lab 22

Komponen | Peran | File Sumber
|---|---|---|
Store | Mock database thread-safe & query counter | internal/blog/store.go
Repository | Logika query Naive loop vs Eager batching | internal/blog/repository.go
Demo | Eksekusi CLI & verifikasi statistik I/O | cmd/demo/main.go

## Entitas Domain: Author dan Post

```go
type Author struct {
    ID   int
    Name string
}

type Post struct {
    ID       int
    AuthorID int
    Title    string
}

type AuthorWithPosts struct {
    Author Author
    Posts  []Post
}
```

## Loop-Driven Data Retrieval

```go
func (r *Repository) GetAuthorsWithPostsNPlusOne() []AuthorWithPosts {
    authors := r.store.GetAllAuthors()
    if len(authors) == 0 { return []AuthorWithPosts{} }

    var result []AuthorWithPosts
    for _, author := range authors {
        posts := r.store.GetPostsByAuthorID(author.ID)
        result = append(result, AuthorWithPosts{
            Author: author,
            Posts:  posts,
        })
    }
    return result
}
```

## Mengapa Slow Query Log Gagal Mendeteksi

Slow query log hanya menandai eksekusi dengan durasi di atas ambang batas (mis. >1000ms). Query N+1 masing-masing berjalan sub-milidetik (<1ms) karena memakai primary/foreign key index. Database engine menganggap sistem bekerja normal tanpa anomali query berat. Development: 3 author memicu 4 query, selesai dalam tempo <5ms. Production: 500 author memicu 501 query, menghabiskan waktu 2.4 detik. Beban berpindah ke latency transfer jaringan dan serialisasi paket di aplikasi.

## Biaya Tersembunyi: Network Round-Trip Time

Komponen Biaya | 1 Query Batch | 100 Query N+1
|---|---|---|
Network RTT | 1 x 2ms = 2ms | 100 x 2ms = 200ms
TCP Wire Overhead | 1 frame exchange | 100 frame exchanges
DB Engine Parsing | 1x AST parse | 100x AST parse
App Context Switch | Minimal | Sangat tinggi (100x I/O)

## Bahaya Connection Pool Exhaustion

Satu HTTP worker menahan koneksi database selama seluruh siklus iterasi loop N. Koneksi tidak dapat dipakai oleh worker lain hingga iterasi terakhir selesai. Kenaikan concurrency membuat kapasitas connection pool langsung habis (exhausted). Request baru mengalami blocking queue di pool dan berakhir dengan gateway timeout (504). Satu endpoint lambat melumpuhkan seluruh service yang berbagi connection pool database.

## Solusi Eager Loading: Pola Batch Data Retrieval

TAHAP 1: Ambil Parent ── GetAllAuthors() → [ID: 1, 2, 3] (Query #1) TAHAP 2: Ekstraksi Kunci ── authorIDs = [1, 2, 3] di memori TAHAP 3: Ambil Child Sekaligus ── GetPostsByAuthorIDs([1, 2, 3]) (Query #2) TAHAP 4: Asosiasi Memori ── Kelompokkan posts ke author via map[int][]Post

## Query Batching di Sisi Store

```go
func (s *Store) GetPostsByAuthorIDs(authorIDs []int) []Post {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.queryCount++

    idMap := make(map[int]bool)
    for _, id := range authorIDs { idMap[id] = true }

    var result []Post
    for _, p := range s.posts {
        if idMap[p.AuthorID] {
            result = append(result, p)
        }
    }
    return result
}
```

## In-Memory Map Grouping pada Layer Aplikasi

```go
// 1. Kelompokkan posts berdasarkan author ID
postsByAuthor := make(map[int][]Post)
for _, post := range allPosts {
    postsByAuthor[post.AuthorID] = append(
        postsByAuthor[post.AuthorID], post)
}

// 2. Susun agregat akhir dengan pencarian O(1)
var result []AuthorWithPosts
for _, author := range authors {
    result = append(result, AuthorWithPosts{
        Author: author,
        Posts:  postsByAuthor[author.ID],
    })
}
```

## Komparasi Kuantitatif Kompleksitas Sistem

Parameter | Naive (Loop) | Eager (Batch)
|---|---|---|
Query Count | 1 + N (Linear) | 1 + 1 (Konstan 2)
Kompleksitas I/O | O(N) relatif record | O(1) relatif record
Data Transfer | Header repetitif | Single packet payload
Memory Allocation | Fragmented allocations | Single pre-sized buffer
Network RTT | N + 1 round-trips | 2 round-trips

## SQL Translation: IN Clause vs Single JOIN

```go
-- Strategi 1: Batching (Dua Query Terpisah)
SELECT id, name FROM authors;
SELECT id, author_id, title FROM posts
WHERE author_id IN (1, 2, 3);

-- Strategi 2: Single JOIN (Satu Query Tunggal)
SELECT a.id, a.name, p.id, p.title
FROM authors a
LEFT JOIN posts p ON a.id = p.author_id;
```

## Batasan Parameter IN Clause pada RDBMS

PostgreSQL & SQLite memiliki limit bind parameter (mis. 65,535 atau 999 parameter). Melebihi batas ini menyebabkan error query execution langsung dari driver database. Statement SQL dengan ribuan literal ID membebani query planner dan parser CPU. Waktu penyusunan AST di engine database dapat melebihi waktu eksekusi indeks. Pecah daftar ID besar menjadi chunk terkendali (misal 200–500 ID per batch query).

## Bahaya Cartesian Product Explosion

Tingkat Relasi | Entitas Record | Total Baris Hasil JOIN
|---|---|---|
1 Author | 1 author | 1 baris
10 Posts | 10 post per author | 1 x 10 = 10 baris
5 Tags | 5 tag per post | 1 x 10 x 5 = 50 baris
10 Comments | 10 komentar per post | 1 x 10 x 5 x 10 = 500 baris!

## Mapping-Level Eager vs Query-Level Eager

FetchType.EAGER pada schema entity menarik seluruh relasi setiap kali entitas dibaca. Operasi sederhana seperti validasi email author terpaksa memuat seluruh riwayat postingan. Menimbulkan memory bloat parah dan membebani garbage collector aplikasi. Relasi dibiarkan lazy secara default pada level definisi domain model. Eager loading hanya diaktifkan secara eksplisit pada query yang benar-benar membutuhkannya.

## Solusi Ekosistem: GraphQL DataLoader Pattern

Resolver #1 (Author 1) ──┐ Resolver #2 (Author 2) ──┼──> [DataLoader Queue] ──> Single Batch Query Resolver #3 (Author 3) ──┘        (Tick Event)        WHERE author_id IN (1, 2, 3)

## Column Selection & Aggregation Optimization

```go
-- 1. Ambil ringkasan agregasi alih-alih seluruh baris post
SELECT author_id, COUNT(*) AS post_count
FROM posts WHERE author_id IN (1, 2, 3)
GROUP BY author_id;

-- 2. Proyeksikan kolom esensial saja
SELECT id, author_id, title
FROM posts WHERE author_id IN (1, 2, 3);
```

## Demo CLI Execution Output

```sh
$ go run ./cmd/demo
--- 1. Simulating N+1 Query Problem ---
Loaded 3 authors with their posts.
Total queries executed: 4 (1 for authors + 3 queries for posts)

--- 2. Simulating Eager Loading (Batching) ---
Loaded 3 authors with their posts.
Total queries executed: 2 (1 for authors + 1 batched query for posts)
```

## Pengujian Otomatis Query Count Invariant

```go
func TestGetAuthorsWithPostsNPlusOne(t *testing.T) {
    store := NewStore()
    repo := NewRepository(store)
    store.ResetQueryCount()

    result := repo.GetAuthorsWithPostsNPlusOne()
    if len(result) != 3 { t.Fatalf("expected 3, got %d", len(result)) }

    queryCount := store.GetQueryCount()
    if queryCount != 4 {
        t.Errorf("expected 4 queries (N+1), got %d", queryCount)
    }
}
```

## Validasi Kesetaraan Semantik (Deep Equivalence)

```go
// Memastikan hasil batching 100% identik dengan naive
nPlusOneResult := repo.GetAuthorsWithPostsNPlusOne()
if !reflect.DeepEqual(result, nPlusOneResult) {
    t.Errorf("eager result mismatch with N+1 result:\n"+
        "Eager: %+v\nN+1: %+v", result, nPlusOneResult)
}
```

## Observability & Deteksi N+1 di Production

APM trace menampilkan 'waterfall graph' berisi puluhan span query identik berurutan. Waktu tunggu jaringan antar span memperjelas latensi akumulatif yang dialami client. Pasang query logger di database driver untuk mendeteksi repetisi query template per HTTP request. Beri threshold alarm jika satu request mengeksekusi >20 query serupa dalam durasi <50ms. Gunakan linter kode untuk menandai eksekusi database calls di dalam perulangan loop for.

## Reproduksi Mandiri & Panduan Verifikasi

```sh
cd labs/22-n-plus-one-query-problem
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Source

[Software Engineering Lab 22 - N+1 Query Problem](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/22-n-plus-one-query-problem), snapshot commit `6649acb9`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
