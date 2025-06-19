package handlers

import (
	"net/http"

	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewGetMetrics(metricsService services.MetricsService) GetMetrics {
	return GetMetrics{metricsService}
}

type GetMetrics struct {
	metricsService services.MetricsService
}

func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query GetMetricsRequest
		if err := ctx.ShouldBindQuery(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		// Usar el contexto de la petición solo para la operación que lo necesita (obtener libros)
		requestCtx := ctx.Request.Context()
		books := h.metricsService.GetBooks(requestCtx)

		// Verificar si se obtuvieron libros correctamente
		if len(books) == 0 {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Could not retrieve books data"})
			return
		}

		// Las operaciones de cálculo puro no necesitan contexto
		meanUnitsSold := h.metricsService.GetMeanUnitsSold(books)
		cheapestBook := h.metricsService.GetCheapestBook(books).Name
		booksWrittenByAuthor := h.metricsService.GetBooksWrittenByAuthor(books, query.Author)

		ctx.JSON(http.StatusOK, gin.H{
			"mean_units_sold":         meanUnitsSold,
			"cheapest_book":           cheapestBook,
			"books_written_by_author": booksWrittenByAuthor,
		})
	}
}
