package main

import (
	"fmt"
	"github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog"
)

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
