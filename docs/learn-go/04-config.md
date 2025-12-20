# Конфигурационный сервис

## Структура файлов

```
config/
├── config.go      # корневой конфиг, собирает все
├── env.go         # загрузка .env и хелперы
├── database.go    # конфиг базы данных
├── server.go      # конфиг сервера
└── log.go         # конфиг логирования
```

```bash
mkdir -p config
touch .env
```

## Установка godotenv

```bash
go get github.com/joho/godotenv
```

## config/env.go — загрузка .env и хелперы

```go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file")
		return
	}
	log.Println(".env file loaded")
}

// getString возвращает строковое значение или дефолтное
func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getInt возвращает целое число или дефолтное значение
func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

// getBool возвращает boolean или дефолтное значение
func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}
```

## config/database.go — конфиг базы данных

```go
package config

type DatabaseConfig struct {
	URL string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URL: getString("DATABASE_URL", "postgres://localhost:5432/myapp"),
	}
}
```

## config/server.go — конфиг сервера

```go
package config

type ServerConfig struct {
	Port  int
	Debug bool
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:  getInt("PORT", 3000),
		Debug: getBool("DEBUG", false),
	}
}
```

## config/log.go — конфиг логирования

```go
package config

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

## config/config.go — корневой конфиг

Собирает все конфиги в один:

```go
package config

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Log      *LogConfig
}

func NewConfig() *Config {
	return &Config{
		Server:   NewServerConfig(),
		Database: NewDatabaseConfig(),
		Log:      NewLogConfig(),
	}
}
```

## Использование в main.go

```go
package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/AlexandrMaltsevYDX/go-cs/internal/home"
)

func main() {
	config.Init()

	cfg := config.NewConfig()

	log.SetLevel(cfg.Log.Level)

	log.Info("Database URL:", cfg.Database.URL)
	log.Info("Debug mode:", cfg.Server.Debug)

	app := fiber.New()
	app.Use(recover.New())

	home.NewHandler(app)

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

## .env

```env
PORT=3000
DEBUG=true
DATABASE_URL=postgres://postgres:secret@localhost:5432/myapp
LOG_LEVEL=debug
```

## Граф зависимостей

```
main()
  │
  ├── config.Init()           // загружаем .env
  │
  └── cfg := config.NewConfig()
        │
        ├── cfg.Server   → ServerConfig{Port, Debug}
        ├── cfg.Database → DatabaseConfig{URL}
        └── cfg.Log      → LogConfig{Level}
```

## .gitignore

```gitignore
.env
```
