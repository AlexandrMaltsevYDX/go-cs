# Go HTML Templates

## Пакет html/template

```go
import "html/template"
```

Безопасные HTML шаблоны с автоэкранированием.

## Базовый синтаксис

```html
<!-- Вывод переменной -->
{{.Name}}

<!-- Вывод поля структуры -->
{{.User.Email}}

<!-- Условие -->
{{if .IsAdmin}}
    <p>Admin panel</p>
{{else}}
    <p>User panel</p>
{{end}}

<!-- Цикл -->
{{range .Items}}
    <li>{{.}}</li>
{{end}}

<!-- Цикл с индексом -->
{{range $i, $item := .Items}}
    <li>{{$i}}: {{$item}}</li>
{{end}}
```

## Создание шаблона

```go
// Из строки
tmpl, err := template.New("hello").Parse("<h1>Hello, {{.Name}}!</h1>")

// Из файла
tmpl, err := template.ParseFiles("templates/index.html")

// Несколько файлов
tmpl, err := template.ParseFiles("base.html", "header.html", "footer.html")

// По паттерну
tmpl, err := template.ParseGlob("templates/*.html")
```

## Выполнение шаблона

```go
data := struct {
    Name string
}{
    Name: "World",
}

// В io.Writer
err := tmpl.Execute(os.Stdout, data)

// В буфер
var buf bytes.Buffer
err := tmpl.Execute(&buf, data)
html := buf.String()
```

## Пример с Fiber

```go
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
```

## Шаблон из файла

`html/page.html`:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Template Test</title>
</head>
<body>
    <h1>Count: {{.Count}}</h1>
</body>
</html>
```

```go
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
```

### template.Must()

Обёртка которая паникует при ошибке — удобно для инициализации:

```go
// Без Must — нужна проверка ошибки
tmpl, err := template.ParseFiles("page.html")
if err != nil {
    panic(err)
}

// С Must — паникует автоматически
tmpl := template.Must(template.ParseFiles("page.html"))
```

### Content-Type header

Важно установить заголовок, чтобы браузер рендерил HTML:

```go
c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)  // text/html
```

Без этого браузер может показать HTML как текст.

> **См. также:** [11-fiber-template-engine.md](./11-fiber-template-engine.md) — рекомендуемый способ работы с шаблонами через Fiber Engine (эффективнее по ресурсам).

## Функции в шаблонах

```go
// Встроенные функции
{{len .Items}}              // длина
{{index .Items 0}}          // элемент по индексу
{{printf "%s" .Name}}       // форматирование

// Кастомные функции
funcs := template.FuncMap{
    "upper": strings.ToUpper,
    "add":   func(a, b int) int { return a + b },
}

tmpl := template.New("").Funcs(funcs).Parse(`
    {{upper .Name}}
    {{add 1 2}}
`)
```

## Вложенные шаблоны

```go
// Определение
{{define "header"}}
    <header>My Site</header>
{{end}}

// Использование
{{template "header"}}

// С данными
{{template "header" .}}
```

## Базовый layout

`base.html`:
```html
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
    {{template "content" .}}
</body>
</html>
```

`page.html`:
```html
{{define "content"}}
    <h1>{{.Title}}</h1>
    <p>{{.Body}}</p>
{{end}}
```

```go
tmpl, _ := template.ParseFiles("base.html", "page.html")
tmpl.ExecuteTemplate(w, "base.html", data)
```

## Композиция шаблонов (подробный пример)

Полный пример с переиспользуемым базовым layout и разными страницами.

### Структура файлов

```
html/
├── layouts/
│   └── base.html      ← базовый layout
└── compose/
    ├── home.html      ← страница "Главная"
    └── about.html     ← страница "О нас"
```

### Базовый layout

`html/layouts/base.html`:
```html
{{ define "base" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <title>{{ template "title" . }}</title>
</head>
<body>
    <header>
        <nav>
            <a href="/home">Главная</a>
            <a href="/about">О нас</a>
        </nav>
    </header>
    <main>{{ template "content" . }}</main>
    <footer>© 2025</footer>
</body>
</html>
{{ end }}
```

### Страницы

`html/compose/home.html`:
```html
{{ define "title" }}Главная страница{{ end }}

{{ define "content" }}
<h1>Главная страница</h1>
<p>Пользователь: {{ .Username }}</p>
{{ end }}
```

`html/compose/about.html`:
```html
{{ define "title" }}О нас{{ end }}

{{ define "content" }}
<h1>О нас</h1>
<p>Компания: {{ .Company }}</p>
{{ end }}
```

### Хендлеры

```go
func (h *Handler) composeHome(c *fiber.Ctx) error {
    // ParseFiles загружает base + page, page переопределяет "title" и "content"
    tmpl := template.Must(template.ParseFiles(
        "./html/layouts/base.html",
        "./html/compose/home.html",
    ))

    data := fiber.Map{"Username": "John"}

    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
        return fiber.NewError(500, err.Error())
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
    return c.Send(buf.Bytes())
}

func (h *Handler) composeAbout(c *fiber.Ctx) error {
    tmpl := template.Must(template.ParseFiles(
        "./html/layouts/base.html",
        "./html/compose/about.html",
    ))

    data := fiber.Map{"Company": "Go Corp"}

    var buf bytes.Buffer
    if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
        return fiber.NewError(500, err.Error())
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
    return c.Send(buf.Bytes())
}
```

### Как это работает

1. `ParseFiles()` загружает оба файла в один template set
2. Страница определяет блоки `title` и `content` через `{{ define }}`
3. `ExecuteTemplate(&buf, "base", data)` рендерит шаблон "base"
4. Внутри base, `{{ template "title" . }}` подставляет блок из страницы
5. Данные (`.Username`, `.Company`) доступны через точку `.`

### Преимущества

- **DRY** — header/footer/nav в одном месте
- **Гибкость** — каждая страница определяет свой title и content
- **Поддержка** — изменения в layout применяются ко всем страницам

## Экранирование

```go
// Автоматически экранируется
{{.UserInput}}  // <script> → &lt;script&gt;

// Без экранирования (опасно!)
{{.TrustedHTML | safe}}

// Или через template.HTML
data := struct {
    Content template.HTML
}{
    Content: template.HTML("<b>Bold</b>"),
}
```

## Комментарии

```html
{{/* Это комментарий */}}
```

## Пробелы

```html
{{- .Name -}}  // Убрать пробелы слева и справа
{{- .Name}}    // Убрать пробелы слева
{{.Name -}}    // Убрать пробелы справа
```

## With — изменение контекста

```html
{{with .User}}
    <p>Name: {{.Name}}</p>
    <p>Email: {{.Email}}</p>
{{end}}
```

## Сравнения

```html
{{if eq .Status "active"}}Active{{end}}
{{if ne .Count 0}}Has items{{end}}
{{if lt .Age 18}}Minor{{end}}
{{if gt .Score 90}}Excellent{{end}}
{{if le .Stock 10}}Low stock{{end}}
{{if ge .Rating 4}}Good{{end}}
```

## Логические операторы

```html
{{if and .IsAdmin .IsActive}}
    Admin is active
{{end}}

{{if or .IsAdmin .IsModerator}}
    Has privileges
{{end}}

{{if not .IsBlocked}}
    Access granted
{{end}}
```
