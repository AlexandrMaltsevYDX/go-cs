# Middleware

## Что такое Middleware?

Функция, которая выполняется **до** или **после** обработчика запроса.

```
Запрос → [Middleware 1] → [Middleware 2] → [Handler] → Ответ
```

## app.Use()

Регистрирует middleware для всех роутов:

```go
app.Use(middleware)  // применяется ко ВСЕМ запросам
```

## Recover Middleware

Перехватывает panic и не даёт серверу упасть.

**Без recover:**
```
panic → сервер падает → все пользователи без сервиса
```

**С recover:**
```
panic → recover ловит → возвращает 500 → сервер работает дальше
```

## Установка

```bash
go get github.com/gofiber/fiber/v2/middleware/recover
```

## Реализация

`cmd/main.go`:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())  // ловит все panic

	// роуты...

	app.Listen(":3000")
}
```

## Проверка

Добавь panic в handler:

```go
func (h *HomeHandler) error(c *fiber.Ctx) error {
	panic("test panic!")  // сервер НЕ упадёт
}
```

```bash
http localhost:3000/api/error
# Вернёт 500, но сервер продолжит работать
```

## Другие полезные middleware

```go
import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

app.Use(logger.New())   // логирование запросов
app.Use(cors.New())     // CORS заголовки
app.Use(recover.New())  // recover от panic
```
