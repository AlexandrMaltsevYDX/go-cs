package example

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes устанавливает маршруты для модуля example
func (h *ExampleHandler) RegisterRoutes(router fiber.Router) {
	// Группируем маршруты под /example
	// Если router передан как /api, то итоговый путь будет /api/example
	route := router.Group("/example")

	route.Post("/", h.Create)
}
