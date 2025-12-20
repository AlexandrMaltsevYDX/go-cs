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
