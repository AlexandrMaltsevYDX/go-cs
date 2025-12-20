package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/home"
)

func main() {
	config.Init()

	// load configurations
	cfg := config.NewConfig()

	// setup logger
	log.SetLevel(cfg.Log.Level)

	log.Info("Database URL:", cfg.Database.URL)
	log.Info("Debug mode:", cfg.Server.Debug)

	// create fiber app
	app := fiber.New()

	// apply middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	home.NewHandler(app)

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
