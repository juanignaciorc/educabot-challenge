package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricsService es un mock del servicio de métricas para pruebas
type MockMetricsService struct {
	mock.Mock
}

// Implementación del mock con el diseño correcto de contexto
func (m *MockMetricsService) GetBooks(ctx context.Context) []domain.Book {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Book)
}

func (m *MockMetricsService) GetMeanUnitsSold(books []domain.Book) uint {
	args := m.Called(books)
	return args.Get(0).(uint)
}

func (m *MockMetricsService) GetCheapestBook(books []domain.Book) domain.Book {
	args := m.Called(books)
	return args.Get(0).(domain.Book)
}

func (m *MockMetricsService) GetBooksWrittenByAuthor(books []domain.Book, author string) uint {
	args := m.Called(books, author)
	return args.Get(0).(uint)
}

func TestGetMetrics_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un mock del servicio
	mockService := new(MockMetricsService)

	// Configurar libros de prueba
	testBooks := []domain.Book{
		{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
		{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
		{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
	}

	// Configurar expectativas del mock
	mockService.On("GetBooks", mock.Anything).Return(testBooks)
	mockService.On("GetMeanUnitsSold", testBooks).Return(uint(11000))
	mockService.On("GetCheapestBook", testBooks).Return(domain.Book{Name: "The Go Programming Language"})
	mockService.On("GetBooksWrittenByAuthor", testBooks, "Alan Donovan").Return(uint(1))

	// Crear el handler con el mock del servicio
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Alan+Donovan", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 11000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, "The Go Programming Language", resBody["cheapest_book"])
	assert.Equal(t, 1, int(resBody["books_written_by_author"].(float64)))

	// Verificar que se llamaron los métodos esperados
	mockService.AssertExpectations(t)
}

func TestGetMetrics_EmptyBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un mock del servicio
	mockService := new(MockMetricsService)

	// Configurar el mock para que devuelva una lista vacía de libros
	mockService.On("GetBooks", mock.Anything).Return([]domain.Book{})

	// Crear el handler con el mock del servicio
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusServiceUnavailable, res.Code)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Contains(t, resBody, "error")

	// Verificar que se llamaron los métodos esperados
	mockService.AssertExpectations(t)
}
