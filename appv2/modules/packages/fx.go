package packages

import "go.uber.org/fx"

var Module = fx.Module(
	"packages",
	fx.Provide(
		New,
	),
)
