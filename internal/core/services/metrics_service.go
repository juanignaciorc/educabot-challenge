package services

import (
	"context"
	"slices"

	"educabot.com/bookshop/internal/core/domain"
	"educabot.com/bookshop/internal/core/ports"
)

// metricsService implementa el puerto MetricsService
type metricsService struct {
	booksRepository ports.BooksRepository
}

// NewMetricsService crea una nueva instancia del servicio de métricas
func NewMetricsService(booksRepository ports.BooksRepository) ports.MetricsService {
	return &metricsService{
		booksRepository: booksRepository,
	}
}

// GetBooks recupera los libros usando el contexto para la operación de red
func (s *metricsService) GetBooks(ctx context.Context) []domain.Book {
	return s.booksRepository.GetBooks(ctx)
}

// GetMeanUnitsSold calcula el promedio de unidades vendidas (no requiere contexto)
func (s *metricsService) GetMeanUnitsSold(books []domain.Book) uint {
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
func (s *metricsService) GetCheapestBook(books []domain.Book) domain.Book {
	if len(books) == 0 {
		return domain.Book{}
	}

	return slices.MinFunc(books, func(a, b domain.Book) int {
		return int(a.Price - b.Price)
	})
}

// GetBooksWrittenByAuthor cuenta los libros escritos por un autor (no requiere contexto)
func (s *metricsService) GetBooksWrittenByAuthor(books []domain.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
