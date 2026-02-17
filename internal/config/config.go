package config

import (
	"time"

	"github.com/ArchDevs/radix/internal/env"
)

type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	DB       DBConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type SecurityConfig struct {
	JwtSecret string
}

type DBConfig struct {
	DataSource string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         env.Int("SERVER_PORT"),
			ReadTimeout:  time.Duration(env.Int("SERVER_READ_TIMEOUT")) * time.Second,
			WriteTimeout: time.Duration(env.Int("SERVER_WRITE_TIMEOUT")) * time.Second,
		},
		DB: DBConfig{
			DataSource: env.String("DB_DSN"),
		},
	}
}
