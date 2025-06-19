package repositories

import (
	"educabot.com/bookshop/internal/core/ports"
	"educabot.com/bookshop/internal/repositories/http"
)

// NewBooksRepository creates a new HTTP-based books repository
func NewBooksRepository() ports.BooksRepository {
	// Create a new HTTP repository to fetch books from the external API
	return http.NewHTTPBooksRepository()
}
