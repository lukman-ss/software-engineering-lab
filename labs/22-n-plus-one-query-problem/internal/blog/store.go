package blog

import "sync"

// Store mocks a database, tracking query counts to demonstrate N+1.
type Store struct {
	mu           sync.Mutex
	authors      []Author
	posts        []Post
	queryCount   int
}

func NewStore() *Store {
	return &Store{
		authors: []Author{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		},
		posts: []Post{
			{ID: 101, AuthorID: 1, Title: "Alice Post 1"},
			{ID: 102, AuthorID: 1, Title: "Alice Post 2"},
			{ID: 103, AuthorID: 2, Title: "Bob Post 1"},
			{ID: 104, AuthorID: 3, Title: "Charlie Post 1"},
			{ID: 105, AuthorID: 3, Title: "Charlie Post 2"},
		},
	}
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
