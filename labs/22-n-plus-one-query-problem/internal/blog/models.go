package blog

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
