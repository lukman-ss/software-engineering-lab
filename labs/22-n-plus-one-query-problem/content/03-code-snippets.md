## Snippet 1 — N+1 Implementation

Source File: `internal/blog/repository.go`
Purpose: Demonstrasi anti-pattern (iterasi query).

```go
func (r *Repository) GetAuthorsWithPostsNPlusOne() []AuthorWithPosts {
	authors := r.store.GetAllAuthors()
	
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
Explanation: Di dalam *loop*, aplikasi memanggil `GetPostsByAuthorID`, menghasilkan N query.

## Snippet 2 — Eager Loading Implementation

Source File: `internal/blog/repository.go`
Purpose: Solusi batching dan penggabungan di memori.

```go
func (r *Repository) GetAuthorsWithPostsEager() []AuthorWithPosts {
	authors := r.store.GetAllAuthors()
	if len(authors) == 0 {
		return nil
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
Explanation: Aplikasi melakukan ekstraksi *Author IDs*. Satu query *batch* dijalankan via `GetPostsByAuthorIDs`. Data digabungkan (*grouping*) ke dalam *map* lokal sebelum dirangkai, memastikan total hanya 2 query yang dikirim.
