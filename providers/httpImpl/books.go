package httpImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/providers"
)

const (
	defaultBooksAPIURL = "https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books"
	timeout            = 10 * time.Second
)

// HTTPBooksProvider implements the BooksProvider interface by fetching books from an HTTP API
type HTTPBooksProvider struct {
	client *http.Client
	apiURL string
}

// NewHTTPBooksProvider creates a new HTTPBooksProvider with a configured HTTP client
func NewHTTPBooksProvider() providers.BooksProvider {
	return &HTTPBooksProvider{
		client: &http.Client{
			Timeout: timeout,
		},
		apiURL: defaultBooksAPIURL,
	}
}

// NewHTTPBooksProviderWithConfig creates a new HTTPBooksProvider with custom configuration
// Útil para pruebas o entornos específicos
func NewHTTPBooksProviderWithConfig(client *http.Client, apiURL string) providers.BooksProvider {
	return &HTTPBooksProvider{
		client: client,
		apiURL: apiURL,
	}
}

// GetBooks fetches books from the external API
func (p *HTTPBooksProvider) GetBooks(ctx context.Context) []models.Book {
	// Verificar si el contexto ya ha sido cancelado
	if ctx.Err() != nil {
		fmt.Printf("Context already canceled: %v\n", ctx.Err())
		return []models.Book{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return []models.Book{}
	}

	// Añadir un header de aceptación para especificar que esperamos JSON
	req.Header.Add("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return []models.Book{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return []models.Book{}
	}

	// Limitar el tamaño del cuerpo para prevenir ataques DoS
	const maxBodySize = 1 << 20 // 1 MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return []models.Book{}
	}

	var books []models.Book
	if err := json.Unmarshal(body, &books); err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return []models.Book{}
	}

	// Validar los datos recibidos para asegurar integridad
	for i, book := range books {
		if book.ID == 0 || book.Name == "" {
			fmt.Printf("Warning: Book at index %d has missing required fields\n", i)
		}
	}

	return books
}
