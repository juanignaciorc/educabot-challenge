package services

import (
	"context"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBooksProvider es un mock para el proveedor de libros
type MockBooksProvider struct {
	mock.Mock
}

func (m *MockBooksProvider) GetBooks(ctx context.Context) []models.Book {
	args := m.Called(ctx)
	return args.Get(0).([]models.Book)
}

func TestGetBooks(t *testing.T) {
	// Crear el mock del proveedor
	mockProvider := new(MockBooksProvider)

	// Configurar los datos de prueba
	testBooks := []models.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 2000, Price: 20},
	}

	// Configurar expectativas - el contexto es necesario para la llamada a GetBooks
	mockProvider.On("GetBooks", mock.Anything).Return(testBooks)

	// Crear el servicio con el mock
	service := NewMetricsService(mockProvider)

	// Ejecutar la función a probar
	result := service.GetBooks(context.Background())

	// Verificar el resultado
	assert.Equal(t, testBooks, result)

	// Verificar que se llamaron los métodos esperados
	mockProvider.AssertExpectations(t)
}

func TestGetMeanUnitsSold(t *testing.T) {
	// Crear el servicio con un mock (no importa para esta prueba)
	service := NewMetricsService(new(MockBooksProvider))

	// Crear datos de prueba
	testBooks := []models.Book{
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
	service := NewMetricsService(new(MockBooksProvider))

	// Crear datos de prueba
	testBooks := []models.Book{
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
	service := NewMetricsService(new(MockBooksProvider))

	// Crear datos de prueba
	testBooks := []models.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 1000, Price: 10},
		{ID: 2, Name: "Book 2", Author: "Author 1", UnitsSold: 2000, Price: 20},
		{ID: 3, Name: "Book 3", Author: "Author 2", UnitsSold: 3000, Price: 30},
	}

	// Probar directamente la función sin contexto
	result := service.GetBooksWrittenByAuthor(testBooks, "Author 1")

	assert.Equal(t, uint(2), result)
}

// Caso de prueba para GetMeanUnitsSold con slice vacío
func TestGetMeanUnitsSold_EmptySlice(t *testing.T) {
	service := NewMetricsService(new(MockBooksProvider))
	result := service.GetMeanUnitsSold([]models.Book{})
	assert.Equal(t, uint(0), result)
}

// Caso de prueba para GetCheapestBook con slice vacío
func TestGetCheapestBook_EmptySlice(t *testing.T) {
	service := NewMetricsService(new(MockBooksProvider))
	result := service.GetCheapestBook([]models.Book{})
	assert.Equal(t, models.Book{}, result)
}
