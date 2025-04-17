package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string        `env:"LOG_LEVEL"`
	Addr            string        `env:"ADDR"`
	WriteTimeout    time.Duration `env:"WRITETIMEOUT"`
	DatabaseURL     string        `envconfig:"DATABASE_URL" required:"true"`
	JWTSecret       string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"720h"`
	Port            string        `envconfig:"PORT" default:"8080"`
	MockEmail       string        `envconfig:"MOCK_EMAIL"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
