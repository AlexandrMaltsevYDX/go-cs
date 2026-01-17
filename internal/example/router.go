package example

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes устанавливает маршруты для модуля example
func RegisterRoutes(router fiber.Router, ctrl *ExampleController) {
	// Группируем маршруты под /example
	// Если router передан как /api, то итоговый путь будет /api/example
	route := router.Group("/example")

	route.Post("/", ctrl.Create)
}
