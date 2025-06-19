package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/internal/core/domain"
)

func TestHTTPBooksRepository_GetBooks(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verificar que se incluye el header Accept
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("Expected Accept header to be 'application/json', got '%s'", accept)
		}

		// Return a sample response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id": 1, "name": "Test Book 1", "author": "Test Author 1", "units_sold": 1000, "price": 25},
			{"id": 2, "name": "Test Book 2", "author": "Test Author 2", "units_sold": 2000, "price": 30}
		]`))
	}))
	defer server.Close()

	// Create a repository that uses the mock server
	repository := &HTTPBooksRepository{
		client: server.Client(),
		apiURL: server.URL,
	}

	// Call the GetBooks method
	books := repository.GetBooks(context.Background())

	// Verify the results
	if len(books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(books))
	}

	// Check the first book
	expectedBook1 := domain.Book{
		ID:        1,
		Name:      "Test Book 1",
		Author:    "Test Author 1",
		UnitsSold: 1000,
		Price:     25,
	}
	if books[0] != expectedBook1 {
		t.Errorf("Expected book %+v, got %+v", expectedBook1, books[0])
	}

	// Check the second book
	expectedBook2 := domain.Book{
		ID:        2,
		Name:      "Test Book 2",
		Author:    "Test Author 2",
		UnitsSold: 2000,
		Price:     30,
	}
	if books[1] != expectedBook2 {
		t.Errorf("Expected book %+v, got %+v", expectedBook2, books[1])
	}
}

func TestHTTPBooksRepository_GetBooks_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a repository that uses the mock server
	repository := &HTTPBooksRepository{
		client: server.Client(),
		apiURL: server.URL,
	}

	// Call the GetBooks method
	books := repository.GetBooks(context.Background())

	// Verify that an empty slice is returned on error
	if len(books) != 0 {
		t.Errorf("Expected 0 books on error, got %d", len(books))
	}
}

func TestHTTPBooksRepository_GetBooks_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not a valid json`))
	}))
	defer server.Close()

	// Create a repository that uses the mock server
	repository := &HTTPBooksRepository{
		client: server.Client(),
		apiURL: server.URL,
	}

	// Call the GetBooks method
	books := repository.GetBooks(context.Background())

	// Verify that an empty slice is returned on JSON parsing error
	if len(books) != 0 {
		t.Errorf("Expected 0 books on JSON parsing error, got %d", len(books))
	}
}

func TestHTTPBooksRepository_GetBooks_ContextCanceled(t *testing.T) {
	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create a repository
	repository := NewHTTPBooksRepository()

	// Call the GetBooks method with canceled context
	books := repository.GetBooks(ctx)

	// Verify that an empty slice is returned on canceled context
	if len(books) != 0 {
		t.Errorf("Expected 0 books on canceled context, got %d", len(books))
	}
}

func TestHTTPBooksRepository_GetBooks_InvalidData(t *testing.T) {
	// Create a mock server that returns data with missing fields
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id": 1, "name": "", "author": "Test Author 1", "units_sold": 1000, "price": 25},
			{"id": 0, "name": "Test Book 2", "author": "Test Author 2", "units_sold": 2000, "price": 30}
		]`))
	}))
	defer server.Close()

	// Create a repository that uses the mock server
	repository := &HTTPBooksRepository{
		client: server.Client(),
		apiURL: server.URL,
	}

	// Call the GetBooks method
	books := repository.GetBooks(context.Background())

	// Verify books are returned even with missing fields, since the repository
	// should handle this gracefully with warnings
	if len(books) != 2 {
		t.Errorf("Expected 2 books with validation warnings, got %d", len(books))
	}

	// Verify that the first book has an empty name
	if books[0].Name != "" {
		t.Errorf("Expected empty name for first book, got '%s'", books[0].Name)
	}

	// Verify that the second book has ID 0
	if books[1].ID != 0 {
		t.Errorf("Expected ID 0 for second book, got %d", books[1].ID)
	}
}