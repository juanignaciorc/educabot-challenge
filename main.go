package main

import (
	"fmt"
	"log"

	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/repositories/mockImpls"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Inicializar el provider - Para una aplicación real, usa HTTPBooksProvider en lugar de MockBooksProvider
	booksProvider := mockImpls.NewMockBooksProvider()

	// Inicializar el servicio - Aquí el contexto se propagará correctamente
	metricsService := services.NewMetricsService(booksProvider)

	// Inicializar el handler con el servicio
	metricsHandler := handlers.NewGetMetrics(metricsService)
	router.GET("/", metricsHandler.Handle())

	fmt.Println("Starting server on :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
