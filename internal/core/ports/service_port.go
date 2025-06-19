package ports

import (
	"context"

	"educabot.com/bookshop/internal/core/domain"
)

// MetricsService define el puerto para los servicios de métricas
type MetricsService interface {
	// GetBooks recupera todos los libros disponibles
	GetBooks(ctx context.Context) []domain.Book
	// GetMeanUnitsSold calcula el promedio de unidades vendidas
	GetMeanUnitsSold(books []domain.Book) uint
	// GetCheapestBook encuentra el libro más barato
	GetCheapestBook(books []domain.Book) domain.Book
	// GetBooksWrittenByAuthor cuenta los libros escritos por un autor
	GetBooksWrittenByAuthor(books []domain.Book, author string) uint
}
