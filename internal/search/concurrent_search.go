package search

import (
	"context"
	"sync"

	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
)

func searchBlogs(tokens []string, ctx context.Context, db *database.Queries, blogsResultSet *safeSet[Blog], wg *sync.WaitGroup) {
	// searching blogs by tags
	if blogsSearchByTags, err := db.SearchBlogsByTags(ctx, tokens); err == nil && blogsSearchByTags != nil {
		for _, value := range blogsSearchByTags {
			blogsResultSet.Add(Blog{
				ID:           value.ID,
				Title:        value.Title,
				Brief:        value.Brief,
				ThumbnailUrl: value.ThumbnailUrl,
				Views:        value.Views,
			})
		}
	}

	// searching blogs by title
	for _, token := range tokens {
		if blogsSearchByTitle, err := db.SearchBlogsByTitle(ctx, database.SearchBlogsByTitleParams{
			Title:   "%" + token,
			Title_2: "%" + token + "%",
			Title_3: token + "%",
		}); err == nil && blogsSearchByTitle != nil {
			for _, value := range blogsSearchByTitle {
				blogsResultSet.Add(Blog{
					ID:           value.ID,
					Title:        value.Title,
					Brief:        value.Brief,
					ThumbnailUrl: value.ThumbnailUrl,
					Views:        value.Views,
				})
			}
		}
	}

	wg.Done()
}

func searchBooks(tokens []string, ctx context.Context, db *database.Queries, booksResultSet *safeSet[Book], wg *sync.WaitGroup) {
	// searching books by tags
	if booksSearchByTags, err := db.SearchBooksByTags(ctx, tokens); err == nil && booksSearchByTags != nil {
		for _, value := range booksSearchByTags {
			booksResultSet.Add(Book{
				ID:            value.ID,
				Name:          value.Name,
				CoverImageUrl: value.CoverImageUrl,
			})
		}
	}

	// searching books by title
	for _, token := range tokens {
		if booksSearchByTitle, err := db.SearchBooksByTitle(ctx, database.SearchBooksByTitleParams{
			Name:   "%" + token,
			Name_2: "%" + token + "%",
			Name_3: token + "%",
		}); err == nil && booksSearchByTitle != nil {
			for _, value := range booksSearchByTitle {
				booksResultSet.Add(Book{
					ID:            value.ID,
					Name:          value.Name,
					CoverImageUrl: value.CoverImageUrl,
				})
			}
		}
	}

	wg.Done()
}

func Search(tokens []string, ctx context.Context, db *database.Queries) ([]*Blog, []*Book) {
	// declaring a wait group
	var waitGroup sync.WaitGroup

	// initializing a blogs result set
	blogsResultSet := NewSafeSet[Blog]()

	// initializing a books result set
	booksResultSet := NewSafeSet[Book]()

	waitGroup.Add(2)
	go searchBlogs(tokens, ctx, db, blogsResultSet, &waitGroup)
	go searchBooks(tokens, ctx, db, booksResultSet, &waitGroup)
	waitGroup.Wait()

	// creating the results
	return blogsResultSet.Keys(), booksResultSet.Keys()
}
