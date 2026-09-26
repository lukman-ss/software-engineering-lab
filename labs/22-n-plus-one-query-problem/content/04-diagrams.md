# Diagrams

## N+1 Query Flow
```text
App -> DB: GetAllAuthors()
DB -> App: 3 Authors
App -> DB: GetPostsByAuthorID(1)
App -> DB: GetPostsByAuthorID(2)
App -> DB: GetPostsByAuthorID(3)
```

## Eager Loading (Batched) Flow
```text
App -> DB: GetAllAuthors()
DB -> App: 3 Authors
App -> DB: GetPostsByAuthorIDs([1, 2, 3])
DB -> App: All related posts
App: Group posts in memory
```
