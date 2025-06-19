//go:build integration

package services_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"educabot.com/bookshop/internal/adapters/repositories/http"
	"educabot.com/bookshop/internal/core/services"
	"github.com/stretchr/testify/assert"
)

// TestMetricsServiceIntegration realiza una prueba de integración
// con el servicio HTTP real de libros
// Para ejecutar esta prueba: go test -tags=integration ./internal/core/services
func TestMetricsServiceIntegration(t *testing.T) {
	// Crear un cliente HTTP con timeout bajo para la prueba
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Usar la API real para esta prueba de integración
	booksRepository := http.NewHTTPBooksRepositoryWithConfig(
		client,
		"https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books",
	)

	// Crear el servicio con el repositorio real
	service := services.NewMetricsService(booksRepository)

	// Crear un contexto con timeout para la prueba
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener los libros
	books := service.GetBooks(ctx)

	// Verificar que se obtuvieron libros
	assert.NotEmpty(t, books, "Se deberían obtener libros de la API real")

	// Verificar las métricas
	meanUnitsSold := service.GetMeanUnitsSold(books)
	assert.NotZero(t, meanUnitsSold, "El promedio de unidades vendidas no debería ser cero")

	cheapestBook := service.GetCheapestBook(books)
	assert.NotEmpty(t, cheapestBook.Name, "El libro más barato debería tener un nombre")

	// Buscar algún autor que exista en los datos
	if len(books) > 0 {
		author := books[0].Author
		booksWrittenByAuthor := service.GetBooksWrittenByAuthor(books, author)
		assert.NotZero(t, booksWrittenByAuthor, "Debería haber al menos un libro del autor")
	}
}
