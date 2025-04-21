package config

import (
	"fmt"
	"go-books-api/internal/models"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	PostgresConfig struct {
		Host     string `env:"POSTGRES_HOST"`
		Port     int    `env:"POSTGRES_PORT"`
		Name     string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		SSLMode  string `env:"POSTGRES_SSL"`
	}

	HTTPConfig struct {
		Host               string        `env:"HTTP_HOST"`
		Port               string        `env:"HTTP_PORT"`
		ReadTimeout        time.Duration `env:"READ_TIMEOUT"`
		WriteTimeout       time.Duration `env:"WRITE_TIMEOUT"`
		MaxHeaderMegabytes int           `env:"MAX_HEADER_MBYTES"`
	}

	Config struct {
		Postgres       PostgresConfig
		HTTP           HTTPConfig
		MigrationsPath string `env:"MIGRATIONS_PATH"`
	}
)

func New() (*Config, error) {
	var (
		err error
		cfg Config
	)

	err = cleanenv.ReadConfig(".env", &cfg)

	if cfg == (Config{}) {
		return nil, models.ErrEmptyConfig
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
