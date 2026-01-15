package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
	. "github.com/qor5/x/v3/ui/vuetify"
	. "github.com/theplant/htmlgo"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"gohome.4gophers.ru/getapp/gohome/appv2/config"
	"gohome.4gophers.ru/getapp/gohome/appv2/database"
	"gohome.4gophers.ru/getapp/gohome/appv2/logger"
	"gohome.4gophers.ru/getapp/gohome/appv2/modules/media"
	"gohome.4gophers.ru/getapp/gohome/appv2/modules/users"
	"gohome.4gophers.ru/getapp/gohome/appv2/s3storage"
	"gohome.4gophers.ru/getapp/gohome/appv2/server"
)

func main() {
	fx.New(
		fx.Provide(
			func() prometheus.Registerer {
				// default prometheus
				return prometheus.DefaultRegisterer
			},
		),
		config.Module,
		database.Module,
		logger.Module,
		server.Module,
		s3storage.Module,

		// modules
		media.Module,
		users.Module,

		fx.Provide(
			configure,
			auth,
		),
		fx.Invoke(
			migrate,
			serve,
		),
	).Run()
}

func serve(lc fx.Lifecycle, srv *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return srv.Serve()
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func auth(users *users.Users, pb *presets.Builder) *login.Builder {
	return users.Auth(pb)
}

func configure(
	db *gorm.DB,
	media *media.Media,
	users *users.Users,
) *presets.Builder {
	b := presets.New()

	// Set up the project name, ORM and Homepage
	b.URIPrefix("/admin").
		BrandTitle("GoHome").
		DataOperator(gorm2op.DataOperator(db)).
		HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			r.Body = VContainer(
				H1("GoHome"),
				P().Text("Система управления пакетами"))
			return
		})

	media.Configure(b)
	users.Configure(b)

	b.MenuOrder(
		"advertisers",
		"campaigns",
		"bgroups",
		"banners",
		b.MenuGroup("Additions").SubItems(
			"media-library",
			"audience",
		).Icon("mdi-plus-box-multiple"),
		b.MenuGroup("Settings").SubItems(
			"users",
		).Icon("mdi-wrench-cog"),
	)

	return b
}

func migrate(users *users.Users) {
	users.Migrate()
}
