package home

import "go.uber.org/fx"

var Home = fx.Module("home",
	fx.Provide(
		New,
	),
)
