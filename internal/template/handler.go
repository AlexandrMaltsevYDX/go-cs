package template

import (
	"bytes"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router fiber.Router
}

func NewHandler(router fiber.Router) {
	h := &Handler{router: router}

	t := h.router.Group("/template")
	t.Get("/", h.index)
	t.Get("/file", h.fromFile)
	t.Get("/engine", h.withEngine)
	t.Get("/ifelse", h.ifElse)
	t.Get("/range", h.rangeExample)
	t.Get("/funcs", h.funcsExample)
	t.Get("/compose", h.composeHome)
	t.Get("/compose/about", h.composeAbout)
}

func (h *Handler) index(c *fiber.Ctx) error {
	tmpl, err := template.New("test").Parse("{{.Count}}")
	data := struct{ Count int }{Count: 1}
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Template error")
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Template compile error")
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(tpl.Bytes())
}

func (h *Handler) fromFile(c *fiber.Ctx) error {
	tmpl := template.Must(template.ParseFiles("./html/page.html"))
	data := struct{ Count int }{Count: 5}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Template compile error")
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(tpl.Bytes())
}

// withEngine uses Fiber's built-in template engine
func (h *Handler) withEngine(c *fiber.Ctx) error {
	data := fiber.Map{"Count": 5}
	return c.Render("page", data)
}

// ifElse demonstrates if-else syntax in templates
func (h *Handler) ifElse(c *fiber.Ctx) error {
	data := fiber.Map{
		"IsAdmin":     true,
		"IsLoggedIn":  true,
		"Username":    "John",
		"Count":       7,
		"Status":      "active",
		"IsModerator": false,
		"IsBlocked":   false,
	}
	return c.Render("ifelse", data)
}

// User struct for range example
type User struct {
	Id   int
	Name string
}

// rangeExample demonstrates range syntax in templates
func (h *Handler) rangeExample(c *fiber.Ctx) error {
	users := []User{
		{Id: 1, Name: "Anton"},
		{Id: 2, Name: "Vasia"},
		{Id: 3, Name: "Maria"},
	}
	names := []string{"Anton", "Vasia", "Maria"}

	data := fiber.Map{
		"Names": names,
		"Users": users,
		"Empty": []string{},
	}
	return c.Render("range", data)
}

// funcsExample demonstrates functions and variables in templates
func (h *Handler) funcsExample(c *fiber.Ctx) error {
	data := fiber.Map{
		"Username":  "John",
		"Age":       25,
		"Items":     []string{"apple", "banana", "cherry"},
		"LowerText": "hello world",
	}
	return c.Render("funcs", data)
}

// composeHome demonstrates template composition with layouts
// Uses standard html/template with ParseFiles to combine base layout + page content
func (h *Handler) composeHome(c *fiber.Ctx) error {
	// ParseFiles loads both templates, page template defines "title" and "content"
	tmpl := template.Must(template.ParseFiles(
		"./html/layouts/base.html",
		"./html/compose/home.html",
	))

	data := fiber.Map{
		"Username": "John",
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(buf.Bytes())
}

// composeAbout demonstrates template composition - about page
func (h *Handler) composeAbout(c *fiber.Ctx) error {
	tmpl := template.Must(template.ParseFiles(
		"./html/layouts/base.html",
		"./html/compose/about.html",
	))

	data := fiber.Map{
		"Company": "Go Corp",
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(buf.Bytes())
}
