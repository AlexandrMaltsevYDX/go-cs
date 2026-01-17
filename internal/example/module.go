package example

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Register инициализирует и регистрирует модуль example
func Register(router fiber.Router, db *pgxpool.Pool) {
	// 1. Создаем репозиторий
	repo := NewExampleRepository(db)

	// 2. Создаем контроллер, внедряя репозиторий
	ctrl := NewController(repo)

	// 3. Регистрируем маршруты, связывая роутер и контроллер
	RegisterRoutes(router, ctrl)
}
