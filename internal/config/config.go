package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	HTTP       HTTP       `envPrefix:"HTTP_"`
	PostgreSQL PostgreSQL `envPrefix:"POSTGRESQL_"`
}

type HTTP struct {
	Host string `env:"HOST"`
	Port uint16 `env:"PORT"`
}

type PostgreSQL struct {
	Host       string `env:"HOST"`
	Port       uint16 `env:"PORT"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
	DBName     string `env:"DBNAME"`
	SSLEnabled bool   `env:"SSLEnabled"`
}

func Load() (Config, error) {
	err := collectEnv()
	if err != nil {
		return Config{}, errors.Wrap(err, "collect env from files")
	}

	cfg, err := fill()
	if err != nil {
		return Config{}, errors.Wrap(err, "fill config from env")
	}
	return cfg, nil
}

func collectEnv() error {
	envFiles := []string{
		".env",
		".env.example",
	}

	for _, f := range envFiles {
		err := godotenv.Load(f)
		if err == nil {
			return nil
		}

		if os.IsNotExist(err) {
			continue
		}

		return errors.Wrapf(err, "read %q file", f)
	}

	return nil
}

func fill() (Config, error) {
	opts := env.Options{
		Prefix: "OHA_",
	}

	var cfg Config
	err := env.ParseWithOptions(&cfg, opts)
	if err != nil {
		return Config{}, errors.Wrap(err, "parse env to struct")
	}
	return cfg, nil
}
