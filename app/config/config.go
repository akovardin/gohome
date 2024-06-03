package config

import (
	"path"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/config"
	"go.uber.org/fx"

	"gohome.4gophers.ru/gohome/app/server"
	"gohome.4gophers.ru/gohome/pkg/logger"
)

type Config struct {
	fx.Out

	Logger logger.Config
	Server server.Config
}

func New(env string, cfg string) Config {
	if env == "base" {
		panic("'base' can not be environment")
	}

	y, err := config.NewYAML(
		config.File(path.Join(cfg, "base.yml")),
		config.File(path.Join(cfg, env+".yml")),
	)
	if err != nil {
		panic(err)
	}

	c := Config{}
	err = y.Get("").Populate(&c)
	if err != nil {
		panic(err)
	}

	err = envconfig.Process(strings.ReplaceAll("lokalization", "-", ""), &c)
	if err != nil {
		panic(err)
	}

	return c
}
