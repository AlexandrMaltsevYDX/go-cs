# Обработка ошибок в Go и Fiber

## В Go нет try/catch

В Go **нет исключений** как в других языках. Вместо этого — **явная обработка ошибок**:

```go
// ❌ Так в Go НЕ делают (нет try/catch)
try {
    result = divide(10, 0)
} catch (error) {
    // handle error
}

// ✅ Так делают в Go
result, err := divide(10, 0)
if err != nil {
    // handle error
}
```

## Пакет errors

Стандартная библиотека для создания ошибок:

```go
import "errors"

func divide(lhs, rhs int) (int, error) {
    if rhs == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return lhs / rhs, nil
}
```

## Идиома `value, err` и `value, ok`

В Go есть два похожих паттерна:

### `value, err` — для ошибок

```go
result, err := divide(10, 0)
if err != nil {
    log.Println("Error:", err)
    return
}
fmt.Println(result)
```

### `value, ok` — для проверки существования

```go
// Получение из map
value, ok := myMap["key"]
if !ok {
    // ключ не найден
}

// Type assertion
str, ok := value.(string)
if !ok {
    // value не является string
}

// Чтение из channel
value, ok := <-ch
if !ok {
    // канал закрыт
}
```

### Разница

| Паттерн | Когда использовать |
|---------|-------------------|
| `value, err` | Операция может завершиться ошибкой |
| `value, ok` | Проверка существования/типа (bool) |

## Паттерн "ошибка — последнее возвращаемое значение"

```go
// Стандартный паттерн Go
func doSomething() (Result, error)

// Примеры из стандартной библиотеки
file, err := os.Open("file.txt")
data, err := json.Marshal(obj)
resp, err := http.Get(url)
```

## Почему так?

1. **Явность** — видно где может быть ошибка
2. **Контроль** — нельзя случайно проигнорировать ошибку
3. **Простота** — нет скрытого control flow
4. **Производительность** — нет overhead от exceptions

## errors.Is() и errors.As()

### errors.Is() — проверка типа ошибки

```go
import "errors"

var ErrNotFound = errors.New("not found")

func findUser(id int) error {
    return ErrNotFound
}

err := findUser(1)
if errors.Is(err, ErrNotFound) {
    // это ошибка "not found"
}
```

### errors.As() — получение ошибки определённого типа

```go
type DivError struct {
    a, b int
}

func (d *DivError) Error() string {
    return fmt.Sprintf("Cannot divide: %d / %d", d.a, d.b)
}

err := div(10, 0)

var divErr *DivError
if errors.As(err, &divErr) {
    // теперь можно получить поля ошибки
    fmt.Println(divErr.a, divErr.b)  // 10, 0
}
```

### Разница

| Функция | Для чего |
|---------|----------|
| `errors.Is(err, target)` | Проверить: это **та самая** ошибка? |
| `errors.As(err, &target)` | Получить ошибку как **конкретный тип** |

## Интерфейс error

```go
type error interface {
    Error() string
}
```

Любой тип с методом `Error() string` реализует интерфейс `error`.

## Кастомные ошибки через структуру

Рекомендуется реализовывать ошибку как **метод с указателем** (receiver function):

```go
type DivError struct {
    a, b int
}

func (d *DivError) Error() string {
    return fmt.Sprintf("Cannot divide by zero: %d / %d", d.a, d.b)
}
```

### Почему указатель `*DivError`?

1. **Уникальность** — каждая ошибка имеет свой адрес в памяти
2. **Сравнение** — можно безопасно сравнивать ошибки через `errors.Is()`
3. **Интерфейс error** — Go ожидает что Error() может быть на указателе

```go
// С указателем — разные ошибки
err1 := &DivError{10, 0}  // адрес 0x1000
err2 := &DivError{10, 0}  // адрес 0x2000
err1 == err2  // false (разные адреса)

// Без указателя — могут быть равны по значению
err1 := DivError{10, 0}
err2 := DivError{10, 0}
err1 == err2  // true (одинаковые значения)
```

### Использование

```go
type DivError struct {
    a, b int
}

func (d *DivError) Error() string {
    return fmt.Sprintf("Cannot divide by zero: %d / %d", d.a, d.b)
}

func div(a, b int) (int, error) {
    if b == 0 {
        return 0, &DivError{a, b}
    } else {
        return a / b, nil
    }
}

answer1, err := div(9, 0)
if err != nil {
    // "Cannot divide by zero: 9 / 0"
    fmt.Println(err)
    return
}
fmt.Println("The answer is:", answer1)
```

### Где вызывается метод Error()?

Вызов происходит **неявно**:

```go
fmt.Println(err)         // Go сам вызовет err.Error()
fmt.Println(err.Error()) // явный вызов — тот же результат
```

**Как это работает:**

1. `div(9, 0)` возвращает `&DivError{9, 0}` — указатель на структуру
2. `err` имеет тип `error` (интерфейс)
3. `fmt.Println(err)` видит что `err` реализует интерфейс `error`
4. Go **автоматически вызывает `err.Error()`** для получения строки

```
fmt.Println(err)
       │
       ▼
err реализует error?
       │ да
       ▼
вызов err.Error()
       │
       ▼
"Cannot divide by zero: 9 / 0"
```

---

## Встроенные статус-коды Fiber

Fiber предоставляет константы для всех HTTP статусов:

```go
// 2xx Success
fiber.StatusOK                    // 200
fiber.StatusCreated               // 201
fiber.StatusNoContent             // 204

// 4xx Client Errors
fiber.StatusBadRequest            // 400
fiber.StatusUnauthorized          // 401
fiber.StatusPaymentRequired       // 402
fiber.StatusForbidden             // 403
fiber.StatusNotFound              // 404
fiber.StatusMethodNotAllowed      // 405
fiber.StatusRequestTimeout        // 408
fiber.StatusConflict              // 409
fiber.StatusUnprocessableEntity   // 422
fiber.StatusTooManyRequests       // 429

// 5xx Server Errors
fiber.StatusInternalServerError   // 500
fiber.StatusBadGateway            // 502
fiber.StatusServiceUnavailable    // 503
```

Использование:
```go
return c.Status(fiber.StatusCreated).JSON(user)
return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
```

## Встроенные ошибки

Fiber предоставляет готовые ошибки:

```go
fiber.ErrBadRequest          // 400
fiber.ErrUnauthorized        // 401
fiber.ErrForbidden           // 403
fiber.ErrNotFound            // 404
fiber.ErrInternalServerError // 500
```

## Использование

```go
func (h *HomeHandler) home(c *fiber.Ctx) error {
	return fiber.ErrBadRequest
}
```

Ответ:
```json
{
  "code": 400,
  "message": "Bad Request"
}
```

## fiber.NewError()

Создание ошибки с кастомным сообщением:

```go
func (h *HomeHandler) home(c *fiber.Ctx) error {
	return fiber.NewError(400, "Limit params is undefined")
}
```

Ответ:
```json
{
  "code": 400,
  "message": "Limit params is undefined"
}
```

## Кастомные ошибки

### Простой вариант

```go
func (h *HomeHandler) getUser(c *fiber.Ctx) error {
	user, err := h.service.FindUser(id)
	if err != nil {
		return fiber.NewError(404, "User not found")
	}
	return c.JSON(user)
}
```

### Структура ошибки

```go
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *AppError) Error() string {
	return e.Message
}
```

### Использование

```go
func (h *HomeHandler) createUser(c *fiber.Ctx) error {
	if err := validate(input); err != nil {
		return c.Status(400).JSON(NewAppError(400, "Validation failed", err.Error()))
	}
	// ...
}
```

## Глобальный обработчик ошибок

```go
app := fiber.New(fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		// Проверяем тип ошибки Fiber
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return c.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	},
})
```

## Проверка

```bash
http localhost:3000/api/
```

Ответ:
```
HTTP/1.1 400 Bad Request

{
    "code": 400,
    "message": "Limit params is undefined"
}
```
