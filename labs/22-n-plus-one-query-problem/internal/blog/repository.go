package blog

type Repository struct {
	store *Store
}

func NewRepository(store *Store) *Repository {
	return &Repository{store: store}
}

// GetAuthorsWithPostsNPlusOne demonstrates the N+1 problem.
// It executes 1 query for authors, and N queries for posts.
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

// GetAuthorsWithPostsEager demonstrates the solution using eager loading (batching).
// It executes 1 query for authors, and 1 query for posts.
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
