# Code Snippets

## Snippet 1 — N+1 Implementation

Source File: `internal/blog/repository.go`

Purpose: Demonstrasi anti-pattern. Iterasi query inside loop menghasilkan N query tambahan.

```go
func (r *Repository) GetAuthorsWithPostsNPlusOne() []AuthorWithPosts {
	authors := r.store.GetAllAuthors()
	if len(authors) == 0 {
		return []AuthorWithPosts{}
	}

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

Explanation: Di dalam loop tiap author, `GetPostsByAuthorID` dipanggil — satu query per author. Untuk 3 author, ini berarti 1 query awal + 3 query = 4 query total.

---

## Snippet 2 — Eager Loading Implementation

Source File: `internal/blog/repository.go`

Purpose: Solusi batching. Kumpulkan semua ID, jalankan satu query `IN`, grouping di memori.

```go
func (r *Repository) GetAuthorsWithPostsEager() []AuthorWithPosts {
	authors := r.store.GetAllAuthors()
	if len(authors) == 0 {
		return []AuthorWithPosts{}
	}

	authorIDs := make([]int, len(authors))
	for i, author := range authors {
		authorIDs[i] = author.ID
	}

	allPosts := r.store.GetPostsByAuthorIDs(authorIDs)

	postsByAuthor := make(map[int][]Post)
	for _, post := range allPosts {
		postsByAuthor[post.AuthorID] = append(postsByAuthor[post.AuthorID], post)
	}

	var result []AuthorWithPosts
	for _, author := range authors {
		result = append(result, AuthorWithPosts{
			Author: author,
			Posts:  postsByAuthor[author.ID],
		})
	}
	return result
}
```

Explanation: 1) Ambil semua author. 2) Ekstrak ID ke slice `[1, 2, 3]`. 3) Jalankan satu batch query `GetPostsByAuthorIDs`. 4) Kelompokkan hasil ke `map[int][]Post`. 5) Bangun struktur akhir. Total: 2 query terlepas dari jumlah author.

---

## Snippet 3 — Mock Store dengan Query Counter

Source File: `internal/blog/store.go`

Purpose: Simulasi database in-memory yang track jumlah query untuk verifikasi.

```go
type Store struct {
	mu         sync.Mutex
	authors    []Author
	posts      []Post
	queryCount int
}

func (s *Store) GetAllAuthors() []Author {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryCount++
	return s.authors
}

func (s *Store) GetPostsByAuthorID(authorID int) []Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryCount++
	var result []Post
	for _, p := range s.posts {
		if p.AuthorID == authorID {
			result = append(result, p)
		}
	}
	return result
}

func (s *Store) GetPostsByAuthorIDs(authorIDs []int) []Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryCount++
	idMap := make(map[int]bool)
	for _, id := range authorIDs {
		idMap[id] = true
	}
	var result []Post
	for _, p := range s.posts {
		if idMap[p.AuthorID] {
			result = append(result, p)
		}
	}
	return result
}

func (s *Store) GetQueryCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queryCount
}

func (s *Store) ResetQueryCount() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryCount = 0
}
```

Explanation: `queryCount` di-increment setiap kali metode data retrieval dipanggil. Ini memungkinkan pengujian kuantitatif: `TestGetAuthorsWithPostsNPlusOne` asserting 4 query, `TestGetAuthorsWithPostsEager` asserting 2 query.

---

## Snippet 4 — Demo Runner

Source File: `cmd/demo/main.go`

Purpose: CLI runner untuk verifikasi visual.

```go
func main() {
	store := blog.NewStore()
	repo := blog.NewRepository(store)

	fmt.Println("--- 1. Simulating N+1 Query Problem ---")
	store.ResetQueryCount()
	nPlusOneResults := repo.GetAuthorsWithPostsNPlusOne()
	fmt.Printf("Loaded %d authors with their posts.\n", len(nPlusOneResults))
	fmt.Printf("Total queries executed: %d (1 query for authors + %d queries for posts)\n\n",
		store.GetQueryCount(), len(nPlusOneResults))

	fmt.Println("--- 2. Simulating Eager Loading (Batching) ---")
	store.ResetQueryCount()
	eagerResults := repo.GetAuthorsWithPostsEager()
	fmt.Printf("Loaded %d authors with their posts.\n", len(eagerResults))
	fmt.Printf("Total queries executed: %d (1 query for authors + 1 batched query for posts)\n",
		store.GetQueryCount())
}
```

Example output:
```
--- 1. Simulating N+1 Query Problem ---
Loaded 3 authors with their posts.
Total queries executed: 4 (1 query for authors + 3 queries for posts)

--- 2. Simulating Eager Loading (Batching) ---
Loaded 3 authors with their posts.
Total queries executed: 2 (1 query for authors + 1 batched query for posts)
```

---

## Snippet 5 — Unit Tests

Source File: `internal/blog/repository_test.go`

Purpose: Verifikasi kuantitatif query count dan data correctness.

```go
func TestGetAuthorsWithPostsNPlusOne(t *testing.T) {
	store := NewStore()
	repo := NewRepository(store)

	store.ResetQueryCount()

	result := repo.GetAuthorsWithPostsNPlusOne()

	if len(result) != 3 {
		t.Fatalf("expected 3 authors, got %d", len(result))
	}

	queryCount := store.GetQueryCount()
	if queryCount != 4 {
		t.Errorf("expected 4 queries (N+1), got %d", queryCount)
	}
}

func TestGetAuthorsWithPostsEager(t *testing.T) {
	store := NewStore()
	repo := NewRepository(store)

	store.ResetQueryCount()

	result := repo.GetAuthorsWithPostsEager()

	if len(result) != 3 {
		t.Fatalf("expected 3 authors, got %d", len(result))
	}

	queryCount := store.GetQueryCount()
	if queryCount != 2 {
		t.Errorf("expected 2 queries (eager loading), got %d", queryCount)
	}

	// Deep equivalence check dengan N+1 result
	nPlusOneResult := repo.GetAuthorsWithPostsNPlusOne()
	if !reflect.DeepEqual(result, nPlusOneResult) {
		t.Errorf("eager result mismatch with N+1 result:\nEager: %+v\nN+1: %+v", result, nPlusOneResult)
	}
}

func TestEmptyStore(t *testing.T) {
	store := &Store{authors: nil, posts: nil}
	repo := NewRepository(store)

	n1 := repo.GetAuthorsWithPostsNPlusOne()
	eager := repo.GetAuthorsWithPostsEager()

	if len(n1) != 0 || len(eager) != 0 {
		t.Fatalf("expected empty results, got n1=%d eager=%d", len(n1), len(eager))
	}
	if !reflect.DeepEqual(n1, eager) {
		t.Errorf("empty result mismatch: n1=%#v eager=%#v", n1, eager)
	}
}
```

Explanation: Test memastikan:
- Query count sesuai ekspektasi (4 untuk N+1, 2 untuk eager).
- Hasil data identik antara kedua pendekatan — batching tidak mengubah semantik.
- Edge case empty store tetap konsisten.
