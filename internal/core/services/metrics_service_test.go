package services

import (
	"context"
	"testing"

	"educabot.com/bookshop/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBooksRepository es un mock para el repositorio de libros
type MockBooksRepository struct {
	mock.Mock
}

func (m *MockBooksRepository) GetBooks(ctx context.Context) []domain.Book {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Book)
}

func TestGetBooks(t *testing.T) {
	// Crear el mock del repositorio
	mockRepo := new(MockBooksRepository)

	// Configurar los datos de prueba
	testBooks := []domain.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 2000, Price: 20},
	}

	// Configurar expectativas - el contexto es necesario para la llamada a GetBooks
	mockRepo.On("GetBooks", mock.Anything).Return(testBooks)

	// Crear el servicio con el mock
	service := NewMetricsService(mockRepo)

	// Ejecutar la función a probar
	result := service.GetBooks(context.Background())

	// Verificar el resultado
	assert.Equal(t, testBooks, result)

	// Verificar que se llamaron los métodos esperados
	mockRepo.AssertExpectations(t)
}

func TestGetMeanUnitsSold(t *testing.T) {
	// Crear el servicio con un mock (no importa para esta prueba)
	service := NewMetricsService(new(MockBooksRepository))

	// Crear datos de prueba
	testBooks := []domain.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 2000, Price: 20},
	}

	// Probar directamente la función sin contexto
	result := service.GetMeanUnitsSold(testBooks)

	// Verificar el resultado
	assert.Equal(t, uint(1500), result)
}

func TestGetCheapestBook(t *testing.T) {
	// Crear el servicio con un mock (no importa para esta prueba)
	service := NewMetricsService(new(MockBooksRepository))

	// Crear datos de prueba
	testBooks := []domain.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 2000, Price: 20},
	}

	// Probar directamente la función sin contexto
	result := service.GetCheapestBook(testBooks)

	assert.Equal(t, "Book 1", result.Name)
	assert.Equal(t, uint(10), result.Price)
}

func TestGetBooksWrittenByAuthor(t *testing.T) {
	// Crear el servicio con un mock (no importa para esta prueba)
	service := NewMetricsService(new(MockBooksRepository))

	// Crear datos de prueba
	testBooks := []domain.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 1", UnitsSold: 2000, Price: 20},
		{ID: 3, Name: "Book 3", Author: "Author 2", UnitsSold: 3000, Price: 30},
	}

	// Probar directamente la función sin contexto
	result := service.GetBooksWrittenByAuthor(testBooks, "Author 1")

	assert.Equal(t, uint(2), result)
}
