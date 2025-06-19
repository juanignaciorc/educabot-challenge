//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"educabot.com/bookshop/internal/adapters/handlers"
	"educabot.com/bookshop/internal/adapters/repositories/http"
	"educabot.com/bookshop/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHandlerIntegration prueba el handler con datos reales
// Para ejecutar esta prueba: go test -tags=integration ./internal/adapters/handlers
func TestHandlerIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un cliente HTTP con timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Crear el repositorio real de libros
	booksRepository := http.NewHTTPBooksRepositoryWithConfig(
		client,
		"https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books",
	)

	// Crear el servicio con el repositorio real
	metricsService := services.NewMetricsService(booksRepository)

	// Crear el handler
	handler := handlers.NewGetMetrics(metricsService)

	// Configurar el router
	r := gin.New()
	r.GET("/metrics", handler.Handle())

	// Realizar una solicitud de prueba
	req := httptest.NewRequest(http.MethodGet, "/metrics?author=Test+Author", nil)

	// Usar un contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	// Verificar que la respuesta es exitosa
	assert.Equal(t, http.StatusOK, res.Code)

	// Verificar que el cuerpo tiene el formato esperado
	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err, "El cuerpo debería ser JSON válido")

	// Verificar que los campos esperados están presentes
	assert.Contains(t, resBody, "mean_units_sold")
	assert.Contains(t, resBody, "cheapest_book")
	assert.Contains(t, resBody, "books_written_by_author")
}

// TestHandlerTimeout prueba que el handler maneja correctamente los timeouts
func TestHandlerTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear un cliente HTTP con un timeout muy corto
	client := &http.Client{
		Timeout: 1 * time.Millisecond, // Timeout demasiado corto para completar
	}

	// Crear el repositorio real de libros con un timeout que causará fallo
	booksRepository := http.NewHTTPBooksRepositoryWithConfig(
		client,
		"https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books",
	)

	// Crear el servicio con el repositorio
	metricsService := services.NewMetricsService(booksRepository)

	// Crear el handler
	handler := handlers.NewGetMetrics(metricsService)

	// Configurar el router
	r := gin.New()
	r.GET("/metrics", handler.Handle())

	// Realizar una solicitud de prueba
	req := httptest.NewRequest(http.MethodGet, "/metrics?author=Test+Author", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	// Verificar que la respuesta indica un error
	assert.Equal(t, http.StatusServiceUnavailable, res.Code)

	// Verificar que el mensaje de error es el esperado
	var resBody map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.Contains(t, resBody, "error")
}
