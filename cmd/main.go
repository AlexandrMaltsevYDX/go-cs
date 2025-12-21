package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/zerolog/log"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/home"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/template"
)

func main() {
	config.Init()

	// load configurations
	cfg := config.NewConfig()

	// setup zerolog
	log.Logger = cfg.Log.NewLogger()

	log.Info().Str("url", cfg.Database.URL).Msg("Database")
	log.Info().Bool("debug", cfg.Server.Debug).Msg("Debug mode")

	// setup template engine
	engine := html.New("./html", ".html")
	engine.AddFuncMap(map[string]interface{}{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"add":     func(a, b int) int { return a + b },
	})

	// create fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// apply middlewares
	app.Use(recover.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))
	home.NewHandler(app)
	template.NewHandler(app)

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
