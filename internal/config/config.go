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
	RateLimit    int
	RateBurst    int
}

type SecurityConfig struct {
	JwtSecret string
	JwtTTL    time.Duration
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
			RateLimit:    env.Int("RATE_LIMIT_PER_SECOND"),
			RateBurst:    env.Int("RATE_BURST"),
		},
		Security: SecurityConfig{
			JwtSecret: env.String("JWT_SECRET"),
			JwtTTL:    time.Duration(env.Int("JWT_TTL_HOURS")) * time.Hour,
		},
		DB: DBConfig{
			DataSource: env.String("DB_DSN"),
		},
	}
}
