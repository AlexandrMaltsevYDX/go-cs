package config

type DatabaseConfig struct {
	URL string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URL: getString("DATABASE_URL", "postgres://localhost:5432/myapp"),
	}
}
