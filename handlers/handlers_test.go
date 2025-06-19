package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricsService es un mock del servicio de métricas para pruebas
type MockMetricsService struct {
	mock.Mock
}

// Implementación del mock con el nuevo diseño correcto de contexto
func (m *MockMetricsService) GetBooks(ctx context.Context) []models.Book {
	args := m.Called(ctx)
	return args.Get(0).([]models.Book)
}

func (m *MockMetricsService) GetMeanUnitsSold(books []models.Book) uint {
	args := m.Called(books)
	return args.Get(0).(uint)
}

func (m *MockMetricsService) GetCheapestBook(books []models.Book) models.Book {
	args := m.Called(books)
	return args.Get(0).(models.Book)
}

func (m *MockMetricsService) GetBooksWrittenByAuthor(books []models.Book, author string) uint {
	args := m.Called(books, author)
	return args.Get(0).(uint)
}

func TestGetMetrics_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un mock del servicio
	mockService := new(MockMetricsService)

	// Configurar libros de prueba
	testBooks := []models.Book{
		{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
		{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
		{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
	}

	// Configurar expectativas del mock
	mockService.On("GetBooks", mock.Anything).Return(testBooks)
	mockService.On("GetMeanUnitsSold", testBooks).Return(uint(11000))
	mockService.On("GetCheapestBook", testBooks).Return(models.Book{Name: "The Go Programming Language"})
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
	mockService.On("GetBooks", mock.Anything).Return([]models.Book{})

	// Crear el handler con el mock del servicio
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/", handler.Handle())

	req := httptest.NewRequest(http.MethodGet, "/?author=Alan+Donovan", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	// Verificar que la respuesta indica un error del servicio
	assert.Equal(t, http.StatusServiceUnavailable, res.Code)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Contains(t, resBody, "error")
	assert.Equal(t, "Could not retrieve books data", resBody["error"])

	// Verificar que se llamaron los métodos esperados
	mockService.AssertExpectations(t)
}

func TestGetMetrics_InvalidQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un struct de prueba con un tipo de parámetro de consulta inválido
	type InvalidQueryParams struct {
		Author int `form:"author"` // Intentamos pasar un int donde se espera un string
	}

	// Mock del context.ShouldBindQuery para simular un error de binding
	mockEngine := gin.New()
	mockEngine.GET("/", func(c *gin.Context) {
		var q InvalidQueryParams
		err := c.ShouldBindQuery(&q)
		assert.Error(t, err, "Debería haber un error de binding")
		c.Status(http.StatusBadRequest)
	})

	req := httptest.NewRequest(http.MethodGet, "/?author=invalid", nil)
	res := httptest.NewRecorder()
	mockEngine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	// Ahora probamos el handler real con un parámetro inválido
	mockService := new(MockMetricsService)
	handler := NewGetMetrics(mockService)

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		// Manipulamos los parámetros para forzar un error de binding
		c.Request.URL.RawQuery = "author[]=%invalid"
		handler.Handle()(c)
	})

	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	res = httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Contains(t, resBody, "error")
}
