package services

import (
	"context"
	"slices"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/providers"
)

type MetricsService interface {
	// Mantener contexto en GetBooks ya que podría necesitar cancelación HTTP
	GetBooks(ctx context.Context) []models.Book
	// Eliminar contexto de funciones de cálculo puro
	GetMeanUnitsSold(books []models.Book) uint
	GetCheapestBook(books []models.Book) models.Book
	GetBooksWrittenByAuthor(books []models.Book, author string) uint
}

type metricsService struct {
	booksProvider providers.BooksProvider
}

func NewMetricsService(booksProvider providers.BooksProvider) MetricsService {
	return &metricsService{
		booksProvider: booksProvider,
	}
}

// GetBooks recupera los libros usando el contexto para la operación de red
func (s *metricsService) GetBooks(ctx context.Context) []models.Book {
	return s.booksProvider.GetBooks(ctx)
}

// GetMeanUnitsSold calcula el promedio de unidades vendidas (no requiere contexto)
func (s *metricsService) GetMeanUnitsSold(books []models.Book) uint {
	if len(books) == 0 {
		return 0
	}

	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

// GetCheapestBook encuentra el libro más barato (no requiere contexto)
func (s *metricsService) GetCheapestBook(books []models.Book) models.Book {
	if len(books) == 0 {
		return models.Book{}
	}

	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

// GetBooksWrittenByAuthor cuenta los libros escritos por un autor (no requiere contexto)
func (s *metricsService) GetBooksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
