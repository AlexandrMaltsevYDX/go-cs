# Выделение контроллера

## Структура

```bash
mkdir -p internal/home
```

## Простой контроллер

`internal/home/handler.go`:

```go
package home

import "github.com/gofiber/fiber/v2"

func Handler(c *fiber.Ctx) error {
    return c.SendString("Hello, World!\n")
}
```

## Main

`cmd/main.go`:

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/AlexandrMaltsevYDX/go-cs/internal/home"
)

func main() {
    app := fiber.New()

    app.Get("/", home.Handler)

    app.Listen(":3000")
}
```

## Расширенный контроллер (с регистрацией роутов)

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
	h.router.Get("/", h.home)
}

func (h *HomeHandler) home(c *fiber.Ctx) error {
	return c.SendString("Hello, World!\n")
}
```

### Порядок объявления в файле

В Go порядок объявления **не влияет на работу**, но есть общепринятое соглашение:

```go
// 1. Структура — "что это"
type HomeHandler struct { ... }

// 2. Конструктор New... — "как создать"
func NewHandler(router fiber.Router) { ... }

// 3. Публичные методы — "что умеет" (API)
func (h *HomeHandler) Home(...) { ... }

// 4. Приватные методы — детали реализации
func (h *HomeHandler) helper(...) { ... }
```

Читать код сверху вниз удобнее когда сначала общее, потом детали.

### Зачем такая структура?

1. **Инкапсуляция** — все роуты и зависимости контроллера в одном месте
2. **Расширяемость** — легко добавить новые зависимости (DB, логгер, конфиг):
   ```go
   type HomeHandler struct {
       router fiber.Router
       db     *sql.DB
       logger *log.Logger
   }
   ```
3. **Тестируемость** — можно подменить зависимости моками
4. **Самодостаточность** — контроллер сам регистрирует свои роуты

### Зачем указатель `&HomeHandler{}`?

```go
h := &HomeHandler{...}  // h — это *HomeHandler (указатель)
```

**Без указателя:**
```go
h := HomeHandler{...}   // копия структуры
h.router.Get("/", h.home)  // метод получит копию h
```

**С указателем:**
```go
h := &HomeHandler{...}  // указатель на структуру
h.router.Get("/", h.home)  // метод получит тот же объект
```

### Как работает память?

```
Stack (быстрая)          Heap (медленная)
┌─────────────┐         ┌─────────────────┐
│ h (8 байт)  │ ──────► │ HomeHandler {   │
│ (указатель) │         │   router: ...   │
└─────────────┘         │ }               │
                        └─────────────────┘
```

1. `&HomeHandler{}` — создаёт структуру в **heap** и возвращает указатель
2. `h` — хранит только адрес (8 байт на 64-bit системе)
3. При вызове `h.home()` — передаётся указатель, а не копия всей структуры

### Жизненный цикл объекта и GC

```go
func NewHandler(router fiber.Router) {
    h := &HomeHandler{router: router}  // 1. Создаём объект
    h.router.Get("/", h.home)          // 2. Регистрируем метод
}                                       // 3. Функция завершилась, h выходит из scope
```

**Вопрос:** Почему объект не удаляется после выхода из `NewHandler()`?

**Ответ:** Потому что Fiber сохранил ссылку на `h.home`.

```
После вызова NewHandler():

┌─────────────────────────────────────────────────────────────┐
│ Fiber App (живёт всё время работы сервера)                  │
│                                                             │
│   routes: [                                                 │
│     {                                                       │
│       path: "/",                                            │
│       handler: h.home ─────────┐                            │
│     }                          │                            │
│   ]                            │                            │
└────────────────────────────────│────────────────────────────┘
                                 │
                                 ▼
                    ┌─────────────────────┐
                    │ HomeHandler {       │  ◄── Объект в heap
                    │   router: ...       │
                    │   home() method     │
                    │ }                   │
                    └─────────────────────┘
```

**Как работает GC (Garbage Collector):**

1. GC сканирует все "корни" (глобальные переменные, стек, регистры)
2. Fiber app — это корень (живёт в `main()`)
3. Fiber держит ссылку на `h.home`
4. `h.home` — это метод, привязанный к `*HomeHandler`
5. Значит `HomeHandler` **достижим** → GC его **не удалит**

```go
// h.home — это на самом деле:
(*HomeHandler).home(h, c *fiber.Ctx)
//             ▲
//             └── указатель на структуру сохраняется!
```

**Когда объект удалится:**
- Когда Fiber app завершится (`app.Shutdown()`)
- Или когда роут будет удалён
- GC увидит, что на `HomeHandler` никто не ссылается → удалит

### Это НЕ синглтон

```go
// Каждый вызов создаёт НОВЫЙ объект:
home.NewHandler(app)  // объект по адресу 0x1000
home.NewHandler(app)  // объект по адресу 0x2000 (другой!)

// Синглтон выглядел бы так:
var globalHandler *HomeHandler

func GetHandler() *HomeHandler {
    if globalHandler == nil {
        globalHandler = &HomeHandler{}
    }
    return globalHandler  // всегда 0x1000
}
```

| | Синглтон | Наш паттерн |
|--|----------|-------------|
| Экземпляров | Всегда 1 | Сколько угодно |
| Глобальная переменная | Да | Нет |
| Контроль создания | Централизованный | Нет |
| Паттерн | Singleton | Constructor / DI |

### Dependency Injection (DI)

**DI** — это паттерн, когда зависимости передаются объекту снаружи, а не создаются внутри.

**Без DI (плохо):**
```go
type HomeHandler struct {}

func NewHandler() {
    h := &HomeHandler{}
    db := sql.Open("postgres", "...")  // создаём зависимость ВНУТРИ
    // ...
}
```

Проблемы:
- Нельзя подменить БД для тестов
- Handler знает как создавать подключение
- Жёсткая связь с конкретной реализацией

**С DI (хорошо):**
```go
type HomeHandler struct {
    router fiber.Router
    db     *sql.DB       // зависимость передаётся снаружи
}

func NewHandler(router fiber.Router, db *sql.DB) {
    h := &HomeHandler{
        router: router,
        db:     db,       // получаем готовую зависимость
    }
    h.router.Get("/", h.home)
}
```

Использование:
```go
func main() {
    db := sql.Open("postgres", "...")  // создаём зависимость в main
    app := fiber.New()
    
    home.NewHandler(app, db)           // передаём (инжектим) зависимость
    
    app.Listen(":3000")
}
```

**Преимущества DI:**

| Аспект | Описание |
|--------|----------|
| Тестируемость | Можно передать mock вместо реальной БД |
| Гибкость | Легко заменить реализацию |
| Читаемость | Видно все зависимости в конструкторе |
| Single Responsibility | Handler не отвечает за создание зависимостей |

**Пример теста с DI:**
```go
func TestHomeHandler(t *testing.T) {
    // Создаём mock БД
    mockDB := &MockDB{}
    
    app := fiber.New()
    home.NewHandler(app, mockDB)  // передаём mock
    
    // Тестируем...
}
```

**Граф зависимостей:**
```
main()
  │
  ├── создаёт db
  ├── создаёт app
  │
  └── home.NewHandler(app, db)
        │
        └── HomeHandler использует db
              │
              └── user.NewHandler(app, db)
                    │
                    └── UserHandler использует тот же db
```

Все зависимости создаются в одном месте (`main`) и передаются вниз — это **Constructor Injection**.

### Когда использовать указатель?

| Ситуация | Указатель | Значение |
|----------|-----------|----------|
| Структура > 64 байт | ✅ | ❌ |
| Нужно изменять поля | ✅ | ❌ |
| Много методов | ✅ | ❌ |
| Маленькая структура (read-only) | ❌ | ✅ |

В нашем случае `HomeHandler` хранит зависимости и будет расширяться — поэтому указатель.

### Один роут = одна структура = один указатель

**При запуске сервера создаётся ровно одна структура:**

```
Heap (память):
┌─────────────────────────┐
│ HomeHandler (1 штука)   │  ◄── создана ОДИН раз при старте
│   router: 0x...         │
└─────────────────────────┘
         ▲
         │ указатель (8 байт)
         │
┌────────┴────────────────────────────────────────┐
│ Fiber routes table                              │
│   GET "/" → h.home (хранит указатель)           │
└─────────────────────────────────────────────────┘
```

**При каждом HTTP запросе — НЕ создаётся новый объект:**

```go
// Fiber просто вызывает метод по указателю:
h.home(ctx)  // h — тот же указатель, что и при регистрации
```

```
Запрос 1: GET /  →  h.home()  →  HomeHandler (0x1000)
Запрос 2: GET /  →  h.home()  →  HomeHandler (0x1000)  ← тот же!
Запрос 3: GET /  →  h.home()  →  HomeHandler (0x1000)  ← тот же!
...
1 000 000 запросов → всё ещё один и тот же HomeHandler
```

**Почему это эффективно:**

| Подход | Память | Действие |
|--------|--------|----------|
| Без указателя | Копировать структуру на каждый вызов | Дорого |
| С указателем | Передать 8 байт адреса | Дёшево |

**Вывод:** Нам не нужно копировать структуру, потому что она реально одна. Достаточно ссылки.

**Поэтому в Go почти всегда:**
- Хендлеры → указатели
- Сервисы → указатели
- Репозитории → указатели

Создаём один раз при старте → используем по ссылке всё время работы сервера.

## Main (расширенный)

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/home"
)

func main() {
	app := fiber.New()

	home.NewHandler(app)

	app.Listen(":3000")
}
```
