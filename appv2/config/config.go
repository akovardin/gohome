package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"

	"gohome.4gophers.ru/getapp/gohome/appv2/database"
	"gohome.4gophers.ru/getapp/gohome/appv2/s3storage"
	"gohome.4gophers.ru/getapp/gohome/appv2/server"
)

type Config struct {
	fx.Out

	Database  database.Config  `yaml:"database"`
	Server    server.Config    `yaml:"server"`
	S3Storage s3storage.Config `yaml:"s3storage"`
}

func New(file string) (Config, error) {
	provider, err := config.NewYAML(
		config.Expand(os.LookupEnv),
		config.File("configs/"+file),
		config.Permissive(),
	)

	if err != nil {
		return Config{}, err
	}

	cfg := Config{}

	err = provider.Get("").Populate(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
