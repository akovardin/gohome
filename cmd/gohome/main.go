package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"gohome.4gophers.ru/getapp/gohome/app/config"
	"gohome.4gophers.ru/getapp/gohome/app/handlers/home"
	"gohome.4gophers.ru/getapp/gohome/app/server"
	"gohome.4gophers.ru/getapp/gohome/pkg/logger"
)

func main() {
	app := &cli.App{
		Name:  "getapp",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
			&cli.StringFlag{
				Name:  "configs",
				Value: "./configs",
				Usage: "configs path",
			},
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start service http server",
				Action: func(ctx *cli.Context) error {
					setup(ctx).Run()

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}

func setup(c *cli.Context) *fx.App {
	env := c.String("env")
	cfg := c.String("configs")

	opts := []fx.Option{}
	opts = append(opts, home.Home)
	opts = append(opts, fx.Provide(
		func() config.Config {
			return config.New(env, cfg)
		},
		server.New,
		logger.New,
	))
	opts = append(opts, fx.Invoke(func(s *server.Server) {}))

	return fx.New(
		opts...,
	)
}
