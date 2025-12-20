package main

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/home"
)

func main() {
	config.Init()

	// load configurations
	cfg := config.NewConfig()

	// setup zerolog
	log.Logger = cfg.Log.NewLogger()

	log.Info().Str("url", cfg.Database.URL).Msg("Database")
	log.Info().Bool("debug", cfg.Server.Debug).Msg("Debug mode")

	// create fiber app
	app := fiber.New()

	// apply middlewares
	app.Use(recover.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))
	home.NewHandler(app)

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
