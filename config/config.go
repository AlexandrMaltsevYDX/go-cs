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
