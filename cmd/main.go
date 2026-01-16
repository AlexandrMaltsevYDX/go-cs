package main

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/example"
	"github.com/AlexandrMaltsevYDX/go-cs/pkg/database"
)

func main() {
	config.Init()

	// load configurations
	cfg := config.NewConfig()

	// setup zerolog
	log.Logger = cfg.Log.NewLogger()

	log.Info().Str("url", cfg.Database.URL).Msg("Database")
	log.Info().Bool("debug", cfg.Server.Debug).Msg("Debug mode")

	// setup database
	db := database.CreateDbPool(cfg.Database)
	defer db.Close()

	// create fiber app
	app := fiber.New(fiber.Config{
		// Views: engine,
	})

	// apply middlewares
	app.Use(recover.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	// serve static files
	app.Static("/", "./static")

	// api group
	api := app.Group("/api")

	// example module
	exampleRepo := example.NewExampleRepository(db)
	exampleHandler := example.NewHandler(exampleRepo)
	exampleHandler.RegisterRoutes(api)

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
