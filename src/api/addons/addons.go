package addons

import (
	"github.com/labstack/echo"
)

type Addons map[string]Addon

type Addon interface {
	Name() string
	Version() string

	// SetLogger(io.Writer)
	Middlewares() []echo.MiddlewareFunc
	RegEchoHandlers(func(string, func(echo.Context) error))

	InjectTplAddons() error

	Setup()
	Destroy()
}

var addons = make(Addons)

func AddAddon(a Addon) {
	addons[a.Name()] = a
}

func ListOfAddons() Addons {
	return addons
}
