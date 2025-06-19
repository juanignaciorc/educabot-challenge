package ports

import (
	"context"

	"educabot.com/bookshop/internal/core/domain"
)

// BooksRepository define el puerto para acceder a los datos de libros
type BooksRepository interface {
	// GetBooks recupera todos los libros disponibles
	GetBooks(ctx context.Context) []domain.Book
}
