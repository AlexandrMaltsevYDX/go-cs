# Логирование

## Встроенный пакет log

Go имеет встроенный пакет `log` — простой и без зависимостей.

```go
import "log"
```

## Уровни логирования

```go
log.Print("Info message")     // обычное сообщение
log.Println("With newline")   // с переносом строки
log.Printf("User %s", name)   // форматированный вывод

log.Fatal("Fatal error")      // выводит и завершает программу (os.Exit(1))
log.Panic("Panic!")           // выводит и вызывает panic
```

## Пример

```go
func (h *HomeHandler) error(c *fiber.Ctx) error {
	log.Print("Info")
	log.Println("Debug")
	log.Print("Warn")
	log.Print("Error")
	log.Panic("Panic")  // программа упадёт (но recover поймает)
	
	return c.SendString("Error")
}
```

## Вывод

```
2025/01/15 10:30:00 Info
2025/01/15 10:30:00 Debug
2025/01/15 10:30:00 Warn
2025/01/15 10:30:00 Error
2025/01/15 10:30:00 Panic
```

## Настройка формата

```go
// Добавить файл и строку
log.SetFlags(log.LstdFlags | log.Lshortfile)
// Вывод: 2025/01/15 10:30:00 main.go:15: message

// Только время
log.SetFlags(log.Ltime)
// Вывод: 10:30:00 message

// Микросекунды
log.SetFlags(log.Ltime | log.Lmicroseconds)
// Вывод: 10:30:00.123456 message
```

## Флаги

```go
log.Ldate         // дата: 2025/01/15
log.Ltime         // время: 10:30:00
log.Lmicroseconds // микросекунды: 10:30:00.123456
log.Llongfile     // полный путь: /home/user/app/main.go:15
log.Lshortfile    // короткий путь: main.go:15
log.LUTC          // UTC время
log.LstdFlags     // Ldate | Ltime (по умолчанию)
```

## Запись в файл

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
	log.Fatal(err)
}
log.SetOutput(file)
```

## Ограничения встроенного log

- Нет уровней (Info, Debug, Warn, Error) из коробки
- Нет структурированного вывода (JSON)
- Нет ротации файлов

## Fiber Log — встроенный логгер с уровнями

Fiber имеет свой пакет `log` с поддержкой уровней:

```go
import "github.com/gofiber/fiber/v2/log"
```

### Уровни логирования

```go
log.Trace("trace message")
log.Debug("debug message")
log.Info("info message")
log.Warn("warn message")
log.Error("error message")
log.Fatal("fatal message")  // завершает программу
log.Panic("panic message")  // вызывает panic
```

### SetLevel — установка уровня

```go
import "github.com/gofiber/fiber/v2/log"

log.SetLevel(log.LevelDebug)  // показывать всё начиная с Debug
log.SetLevel(log.LevelInfo)   // показывать Info, Warn, Error (по умолчанию)
log.SetLevel(log.LevelWarn)   // только Warn и Error
log.SetLevel(log.LevelError)  // только Error
```

### Настройка через конфиг

`config/env.go`:

```go
import "github.com/gofiber/fiber/v2/log"

type LogConfig struct {
    Level log.Level
}

func NewLogConfig() *LogConfig {
    level := getString("LOG_LEVEL", "info")
    
    var logLevel log.Level
    switch level {
    case "trace":
        logLevel = log.LevelTrace
    case "debug":
        logLevel = log.LevelDebug
    case "info":
        logLevel = log.LevelInfo
    case "warn":
        logLevel = log.LevelWarn
    case "error":
        logLevel = log.LevelError
    default:
        logLevel = log.LevelInfo
    }
    
    return &LogConfig{Level: logLevel}
}
```

### Использование в main.go

```go
import "github.com/gofiber/fiber/v2/log"

func main() {
    config.Init()
    
    logConfig := config.NewLogConfig()
    log.SetLevel(logConfig.Level)
    
    app := fiber.New()
    app.Use(recover.New())
    
    log.Info("Server starting")
    app.Listen(":3000")
}
```

### .env

```env
LOG_LEVEL=debug
```
