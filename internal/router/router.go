package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AlexandrMaltsevYDX/go-cs/internal/example"
)

// SetupRoutes - centralized router setup
func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// --- Static Files ---
	// Moved from main.go
	app.Static("/", "./static")

	// --- API Groups ---
	api := app.Group("/api")

	// --- Group Middlewares ---
	// Here you can add auth middlewares, e.g.:
	// api.Use(middleware.Auth())

	// --- Modules Registration ---
	// Pass the specific group (api, v1, public, etc.) to the module
	example.Register(api, db)

	// Future modules:
	// users.Register(api, db)
	// auth.Register(api, db)
}
