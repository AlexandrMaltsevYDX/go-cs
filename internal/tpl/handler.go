package tpl

import (
	tadapter "github.com/AlexandrMaltsevYDX/go-cs/pkg"
	"github.com/AlexandrMaltsevYDX/go-cs/views"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router fiber.Router
}

func NewHandler(router fiber.Router) {
	h := &Handler{router: router}
	t := h.router.Group("/tpl")
	t.Get("/", h.index)
}

func (h *Handler) index(c *fiber.Ctx) error {
	component := views.Hello("Alexandr")
	return tadapter.Render(c, component)
}
