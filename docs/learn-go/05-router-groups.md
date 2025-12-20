# Группы роутов (Router Groups)

## Зачем нужны группы?

1. **Общий префикс** — `/api/users`, `/api/posts` вместо повторения `/api`
2. **Общие middleware** — авторизация для группы роутов
3. **Организация кода** — логическое разделение эндпоинтов

## Базовый пример

```go
api := app.Group("/api")
api.Get("/", handler)        // GET /api/
api.Get("/users", getUsers)  // GET /api/users
```

## Реализация в handler

`internal/home/handler.go`:

```go
package home

import "github.com/gofiber/fiber/v2"

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
	return c.SendString("Hello")
}

func (h *HomeHandler) error(c *fiber.Ctx) error {
	return c.SendString("Error")
}
```

## Результат

```
GET /api/       →  "Hello"
GET /api/error  →  "Error"
```

## Проверка

```bash
http localhost:3000/api/
http localhost:3000/api/error
```

## Вложенные группы

```go
api := app.Group("/api")

v1 := api.Group("/v1")
v1.Get("/users", getUsersV1)  // GET /api/v1/users

v2 := api.Group("/v2")
v2.Get("/users", getUsersV2)  // GET /api/v2/users
```

## Middleware для группы

```go
api := app.Group("/api", authMiddleware)  // middleware для всей группы
api.Get("/users", getUsers)   // защищён authMiddleware
api.Get("/posts", getPosts)   // защищён authMiddleware

// Публичные роуты без middleware
app.Get("/health", healthCheck)
```
