package memory

import (
	"context"

	"educabot.com/bookshop/internal/core/domain"
)

// MemoryBooksRepository implementa el repositorio de libros en memoria
type MemoryBooksRepository struct{}

// NewMemoryBooksRepository crea una nueva instancia del repositorio de libros en memoria
func NewMemoryBooksRepository() *MemoryBooksRepository {
	return &MemoryBooksRepository{}
}

// GetBooks implementa la interfaz BooksRepository
// Nota: el contexto se ignora con _ ya que los datos son estáticos
// En una implementación real, el contexto se usaría para cancelación de operaciones
func (m *MemoryBooksRepository) GetBooks(_ context.Context) []domain.Book {
	// Este repositorio en memoria no usa el contexto porque devuelve datos estáticos
	// En una implementación real, el contexto sería útil para cancelar
	// operaciones HTTP o de base de datos
	return []domain.Book{
		{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
		{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
		{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
	}
}