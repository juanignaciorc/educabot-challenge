package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"educabot.com/bookshop/internal/core/domain"
)

const (
	defaultBooksAPIURL = "https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books"
	timeout            = 10 * time.Second
)

// HTTPBooksRepository implementa el repositorio de libros usando HTTP
type HTTPBooksRepository struct {
	client *http.Client
	apiURL string
}

// NewHTTPBooksRepository crea una nueva instancia del repositorio de libros HTTP
func NewHTTPBooksRepository() *HTTPBooksRepository {
	return &HTTPBooksRepository{
		client: &http.Client{
			Timeout: timeout,
		},
		apiURL: defaultBooksAPIURL,
	}
}

// NewHTTPBooksRepositoryWithConfig crea una nueva instancia del repositorio de libros HTTP con configuración personalizada
// Útil para pruebas o entornos específicos
func NewHTTPBooksRepositoryWithConfig(client *http.Client, apiURL string) *HTTPBooksRepository {
	return &HTTPBooksRepository{
		client: client,
		apiURL: apiURL,
	}
}

// GetBooks obtiene libros desde la API externa
func (p *HTTPBooksRepository) GetBooks(ctx context.Context) []domain.Book {
	// Verificar si el contexto ya ha sido cancelado
	if ctx.Err() != nil {
		fmt.Printf("Context already canceled: %v\n", ctx.Err())
		return []domain.Book{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return []domain.Book{}
	}

	// Añadir un header de aceptación para especificar que esperamos JSON
	req.Header.Add("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return []domain.Book{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return []domain.Book{}
	}

	// Limitar el tamaño del cuerpo para prevenir ataques DoS
	const maxBodySize = 1 << 20 // 1 MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return []domain.Book{}
	}

	var books []domain.Book
	if err := json.Unmarshal(body, &books); err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return []domain.Book{}
	}

	// Validar los datos recibidos para asegurar integridad
	for i, book := range books {
		if book.ID == 0 || book.Name == "" {
			fmt.Printf("Warning: Book at index %d has missing required fields\n", i)
		}
	}

	return books
}