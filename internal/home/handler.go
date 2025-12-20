package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type HomeHandler struct {
	router fiber.Router
}

func NewHandler(router fiber.Router) {
	h := &HomeHandler{
		router: router,
	}

	api := h.router.Group("/api")
	api.Get("/", h.home)
	api.Get("/error", h.error)
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	log.Debug().Msg("Debug: home endpoint called")
	log.Info().Msg("Info: home endpoint called")
	log.Warn().Msg("Warn: home endpoint called")
	return c.SendString("Hello!")
}

func (h *HomeHandler) error(c *fiber.Ctx) error {
	log.Error().Msg("Error: error endpoint called")
	return c.SendString("Error")
}
