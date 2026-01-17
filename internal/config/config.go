package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" env-default:"8080"`
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-default:"user"`
	DBPassword string `env:"DB_PASSWORD" env-default:""`
	DBName     string `env:"DB_NAME" env-default:""`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}
	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
