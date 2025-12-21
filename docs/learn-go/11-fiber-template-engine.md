# Fiber Template Engine

Встроенный движок шаблонов для Fiber — эффективнее ручного способа.

## Установка

```bash
go get github.com/gofiber/template/html/v2
```

## Настройка

### main.go

```go
import "github.com/gofiber/template/html/v2"

func main() {
    // Создаём engine: папка + расширение
    engine := html.New("./html", ".html")
    
    // Передаём в Fiber
    app := fiber.New(fiber.Config{
        Views: engine,
    })
    
    app.Listen(":3000")
}
```

### Структура файлов

```
project/
├── cmd/main.go
├── html/
│   ├── page.html
│   ├── index.html
│   └── layouts/
│       └── main.html
```

## Использование

### Хендлер

```go
func (h *Handler) home(c *fiber.Ctx) error {
    data := fiber.Map{"Count": 5}
    return c.Render("page", data)  // → ./html/page.html
}
```

**Одна строка** вместо 10!

### Шаблон

`html/page.html`:
```html
<!DOCTYPE html>
<html>
<body>
    <h1>Count: {{.Count}}</h1>
</body>
</html>
```

## Сравнение: ручной vs engine

| Ручной способ | Template Engine |
|---------------|-----------------|
| `ParseFiles()` на каждый запрос | Парсинг **один раз** при старте |
| `bytes.Buffer` каждый раз | Engine управляет буферами |
| Ручной `Content-Type` | Автоматически |
| ~10 строк | **1 строка** |

## Почему engine эффективнее?

```
РУЧНОЙ СПОСОБ (каждый запрос):
───────────────────────────────────────
Запрос 1: Parse → Alloc → Execute → GC
Запрос 2: Parse → Alloc → Execute → GC
Запрос 3: Parse → Alloc → Execute → GC
          ↑ парсим файл заново!

ENGINE:
───────────────────────────────────────
Старт:    Parse all → cache ✓
Запрос 1: Execute (from cache)
Запрос 2: Execute (from cache)
Запрос 3: Execute (from cache)
          ↑ шаблоны в памяти!
```

### Экономия ресурсов:

1. **CPU** — нет повторного парсинга файлов
2. **Memory** — переиспользование буферов
3. **I/O** — файлы читаются один раз
4. **GC** — меньше мусора

## Reload в development

```go
engine := html.New("./html", ".html")
engine.Reload(true)  // перечитывать при изменении
```

В production — `Reload(false)` (по умолчанию).

## Layouts

### Базовый layout

`html/layouts/main.html`:
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <nav>Menu</nav>
    {{embed}}
    <footer>Footer</footer>
</body>
</html>
```

### Страница

`html/page.html`:
```html
<h1>{{.Title}}</h1>
<p>Content here</p>
```

### Рендеринг с layout

```go
return c.Render("page", fiber.Map{
    "Title": "Home",
}, "layouts/main")
```

Результат:
```html
<!DOCTYPE html>
<html>
<head><title>Home</title></head>
<body>
    <nav>Menu</nav>
    <h1>Home</h1>
    <p>Content here</p>
    <footer>Footer</footer>
</body>
</html>
```

## Partials (частичные шаблоны)

`html/partials/header.html`:
```html
<header>
    <h1>{{.SiteName}}</h1>
</header>
```

Использование:
```html
{{template "partials/header" .}}
```

## Кастомные функции

```go
engine := html.New("./html", ".html")

engine.AddFunc("upper", strings.ToUpper)
engine.AddFunc("formatDate", func(t time.Time) string {
    return t.Format("02.01.2006")
})
```

В шаблоне:
```html
{{upper .Name}}
{{formatDate .CreatedAt}}
```

## Передача данных

### fiber.Map

```go
c.Render("page", fiber.Map{
    "Title":   "Home",
    "Count":   42,
    "Items":   []string{"a", "b", "c"},
    "IsAdmin": true,
})
```

### Структура

```go
type PageData struct {
    Title string
    Count int
}

c.Render("page", PageData{
    Title: "Home",
    Count: 42,
})
```

## Обработка ошибок

```go
func (h *Handler) page(c *fiber.Ctx) error {
    err := c.Render("page", data)
    if err != nil {
        log.Error().Err(err).Msg("Template render error")
        return fiber.NewError(500, "Render failed")
    }
    return nil
}
```

## Другие движки

Fiber поддерживает разные template engines:

```bash
go get github.com/gofiber/template/pug/v2      # Pug
go get github.com/gofiber/template/mustache/v2 # Mustache
go get github.com/gofiber/template/handlebars/v2 # Handlebars
go get github.com/gofiber/template/jet/v2      # Jet
```

Все работают одинаково — меняется только синтаксис шаблонов.
