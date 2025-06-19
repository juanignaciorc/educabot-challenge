package mockImpls

import (
	"context"

	"educabot.com/bookshop/models"
)

//comentario

type MockBooksProvider struct{}

func NewMockBooksProvider() *MockBooksProvider {
	return &MockBooksProvider{}
}

// GetBooks implementa la interfaz BooksProvider
// Nota: el contexto se ignora con _ ya que los datos son estáticos
// En una implementación real, el contexto se usaría para cancelación de operaciones
func (m *MockBooksProvider) GetBooks(_ context.Context) []models.Book {
	// Este mock no usa el contexto porque devuelve datos estáticos
	// En una implementación real, el contexto sería útil para cancelar
	// operaciones HTTP o de base de datos
	return []models.Book{
		{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
		{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
		{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
	}
}
