# Zerolog — Zero Allocation JSON Logger

## Установка

```bash
go get github.com/rs/zerolog
```

## Почему zerolog?

- **Нулевые аллокации** — самый быстрый JSON логгер
- **Структурированный вывод** — JSON формат
- **Типизированные поля** — Str, Int, Bool и т.д.
- **Pretty logging** — для разработки

## Базовое использование

```go
import (
    "os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Pretty output для разработки
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    
    log.Info().Msg("Hello world")
}
```

## Уровни логирования

```go
log.Trace().Msg("trace")    // -1
log.Debug().Msg("debug")    // 0
log.Info().Msg("info")      // 1
log.Warn().Msg("warn")      // 2
log.Error().Msg("error")    // 3
log.Fatal().Msg("fatal")    // 4 — завершает программу
log.Panic().Msg("panic")    // 5 — вызывает panic
```

## SetGlobalLevel — установка уровня

```go
zerolog.SetGlobalLevel(zerolog.DebugLevel)  // показывать всё с Debug
zerolog.SetGlobalLevel(zerolog.InfoLevel)   // Info и выше
zerolog.SetGlobalLevel(zerolog.WarnLevel)   // только Warn и выше
zerolog.SetGlobalLevel(zerolog.Disabled)    // отключить логи
```

## Контекстные поля

```go
log.Info().
    Str("user", "john").
    Int("age", 25).
    Bool("active", true).
    Msg("User logged in")

// {"level":"info","user":"john","age":25,"active":true,"message":"User logged in"}
```

## Типы полей

```go
// Строки и числа
.Str("key", "value")
.Int("count", 42)
.Float64("price", 19.99)
.Bool("enabled", true)

// Ошибки
.Err(err)

// Время
.Time("created", time.Now())
.Dur("duration", elapsed)

// Массивы
.Strs("tags", []string{"go", "fiber"})
.Ints("ids", []int{1, 2, 3})
```

## Логирование ошибок

```go
err := errors.New("connection failed")

log.Error().
    Err(err).
    Str("host", "localhost").
    Msg("Database error")

// {"level":"error","error":"connection failed","host":"localhost","message":"Database error"}
```

## Создание логгера с контекстом

```go
// Logger с постоянными полями
logger := zerolog.New(os.Stderr).With().
    Timestamp().
    Str("service", "api").
    Logger()

logger.Info().Msg("Starting")
// {"level":"info","time":"...","service":"api","message":"Starting"}
```

## Sub-logger

```go
userLogger := log.With().
    Str("component", "auth").
    Logger()

userLogger.Info().Msg("Login attempt")
// {"level":"info","component":"auth","message":"Login attempt"}
```

## Pretty Console для разработки

```go
// Простой вариант
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// С настройками
output := zerolog.ConsoleWriter{
    Out:        os.Stderr,
    TimeFormat: "15:04:05",
    NoColor:    false,
}
log.Logger = log.Output(output)

// Вывод: 10:30:00 INF Hello world user=john
```

## Глобальные настройки

```go
// Формат времени
zerolog.TimeFieldFormat = time.RFC3339
zerolog.TimeFieldFormat = zerolog.TimeFormatUnix  // Unix timestamp

// Имена полей
zerolog.TimestampFieldName = "t"
zerolog.LevelFieldName = "l"
zerolog.MessageFieldName = "m"
```

## Интеграция с Fiber

```go
package main

import (
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Pretty output для разработки
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
    
    app := fiber.New()
    
    app.Get("/", func(c *fiber.Ctx) error {
        log.Info().
            Str("path", c.Path()).
            Str("method", c.Method()).
            Msg("Request")
        return c.SendString("Hello")
    })
    
    log.Info().Msg("Server starting on :3000")
    app.Listen(":3000")
}
```

## Калькуляция только при включенном уровне

```go
// Дорогая операция выполняется только если Debug включен
if e := log.Debug(); e.Enabled() {
    value := expensiveComputation()
    e.Str("result", value).Msg("Computed")
}
```

## Вывод в несколько мест

```go
multi := zerolog.MultiLevelWriter(
    zerolog.ConsoleWriter{Out: os.Stdout},
    os.Stderr,
)
log.Logger = zerolog.New(multi).With().Timestamp().Logger()
```

## Сравнение с Fiber log

| Fiber log | Zerolog |
|-----------|---------|
| `log.Info("msg")` | `log.Info().Msg("msg")` |
| `log.SetLevel(log.LevelDebug)` | `zerolog.SetGlobalLevel(zerolog.DebugLevel)` |
| Простой вывод | JSON / Pretty |
| Нет типизации полей | Str, Int, Bool и т.д. |

## Миграция с Fiber log

```go
// Было (Fiber log)
log.Info("User logged in:", username)

// Стало (zerolog)
log.Info().Str("user", username).Msg("User logged in")
```

## Fiber Zerolog Middleware

Middleware для логирования HTTP запросов через zerolog.

### Установка

```bash
go get github.com/gofiber/contrib/fiberzerolog
```

### Использование

```go
import (
    "os"
    "github.com/gofiber/contrib/fiberzerolog"
    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    // Настройка zerolog
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    
    app := fiber.New()
    
    // Подключение middleware
    app.Use(fiberzerolog.New(fiberzerolog.Config{
        Logger: &log.Logger,
    }))
    
    app.Get("/", handler)
    app.Listen(":3000")
}
```

### Вывод middleware

```
10:30:00 INF Success latency=1.234ms status=200 method=GET url=/api/
10:30:01 WRN Client error latency=0.5ms status=404 method=GET url=/notfound
10:30:02 ERR Server error latency=2ms status=500 method=POST url=/api/error
```

### Config

| Параметр | Тип | Описание | Default |
|----------|-----|----------|---------|
| Logger | `*zerolog.Logger` | Кастомный логгер | stderr с timestamp |
| Fields | `[]string` | Какие поля логировать | ip, latency, status, method, url, error |
| Messages | `[]string` | Сообщения для 5xx, 4xx, 2xx | Server error, Client error, Success |
| Levels | `[]zerolog.Level` | Уровни для 5xx, 4xx, 2xx | Error, Warn, Info |
| SkipURIs | `[]string` | URI которые не логировать | `[]` |

### Кастомные поля

```go
app.Use(fiberzerolog.New(fiberzerolog.Config{
    Logger: &log.Logger,
    Fields: []string{"ip", "latency", "status", "method", "url", "path", "ua"},
}))
```

### Доступные поля

- `ip` — IP клиента
- `latency` — время обработки
- `status` — HTTP статус
- `method` — метод (GET, POST...)
- `url` — полный URL
- `path` — путь
- `ua` — User-Agent
- `body` — тело запроса
- `resBody` — тело ответа
- `error` — ошибка
- `reqHeaders` — заголовки запроса
- `resHeaders` — заголовки ответа

## Конфигурация логгера через ENV

### Структура конфига

`config/log.go`:

```go
type LogConfig struct {
    Level  zerolog.Level
    Pretty bool
}

func NewLogConfig() *LogConfig {
    level := getString("LOG_LEVEL", "info")
    pretty := getBool("LOG_PRETTY", false)

    return &LogConfig{
        Level:  parseLogLevel(level),
        Pretty: pretty,
    }
}

func parseLogLevel(level string) zerolog.Level {
    switch level {
    case "trace":
        return zerolog.TraceLevel
    case "debug":
        return zerolog.DebugLevel
    case "info":
        return zerolog.InfoLevel
    case "warn":
        return zerolog.WarnLevel
    case "error":
        return zerolog.ErrorLevel
    default:
        return zerolog.InfoLevel
    }
}
```

### Фабрика логгера

```go
// NewLogger creates a zerolog.Logger based on config
func (c *LogConfig) NewLogger() zerolog.Logger {
    var logger zerolog.Logger

    switch c.Pretty {
    case true:
        logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
    default:
        logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
    }

    zerolog.SetGlobalLevel(c.Level)
    return logger
}
```

### Использование в main.go

```go
func main() {
    config.Init()
    cfg := config.NewConfig()

    // Одна строка — логгер настроен!
    log.Logger = cfg.Log.NewLogger()

    log.Info().Msg("Application started")
}
```

### .env

```env
LOG_LEVEL=debug     # trace, debug, info, warn, error
LOG_PRETTY=true     # true = цветной консольный вывод
                    # false = JSON (для продакшена)
```

### Вывод

**Pretty (LOG_PRETTY=true):**
```
10:30:00 INF Application started
10:30:00 DBG Debug message
```

**JSON (LOG_PRETTY=false):**
```json
{"level":"info","time":"2025-12-20T10:30:00Z","message":"Application started"}
{"level":"debug","time":"2025-12-20T10:30:00Z","message":"Debug message"}
```

