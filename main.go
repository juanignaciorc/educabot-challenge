package main

import (
	"fmt"
	"log"

	"educabot.com/bookshop/internal/adapters/handlers"
	"educabot.com/bookshop/internal/core/services"
	"educabot.com/bookshop/internal/repositories/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Inicializar el repositorio - Usando el repositorio HTTP para obtener datos reales
	booksRepository := http.NewHTTPBooksRepository()

	// Inicializar el servicio - Aquí el contexto se propagará correctamente
	metricsService := services.NewMetricsService(booksRepository)

	// Inicializar el handler con el servicio
	metricsHandler := handlers.NewGetMetrics(metricsService)
	router.GET("/", metricsHandler.Handle())

	fmt.Println("Starting server on :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
