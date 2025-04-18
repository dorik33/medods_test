package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string        `env:"LOG_LEVEL"`
	Addr            string        `env:"ADDR"`
	WriteTimeout    time.Duration `env:"WRITETIMEOUT"`
	DatabaseURL     string        `env:"DATABASE_URL" required:"true"`
	JWTSecret       string        `env:"JWT_SECRET" required:"true"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" default:"720h"`
	MockEmail       string        `env:"MOCK_EMAIL"`
	SMTPConfig      SMTPConfig
}

type SMTPConfig struct {
	EmailFrom     string `env:"EMAIL_FROM"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
	EmailTo       string `env:"EMAIL_TO"`
	SMTPServer    string `env:"EMAIL_SMTP_SERVER"`
	SMTPPort      string `env:"EMAIL_SMTP_PORT"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
