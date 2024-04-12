package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"path"
	"time"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Addr            string        `yaml:"addr" env-default:"localhost:8080"`
		Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"30s"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"3s"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		URL string `env-required:"true"   yaml:"url" env:"PG_URL"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	godotenv.Load()

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("errs reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("errs updating env: %w", err)
	}
	return cfg, nil
}
