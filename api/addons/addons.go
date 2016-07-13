package addons

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

var Addons = make(addons)

type Addon interface {
	Name() string
	Version() string

	// SetLogger(io.Writer)
	Middlewares() []echo.MiddlewareFunc
	RegEchoHandlers(func(string, func(echo.Context) error))

	InjectTplAddons() error

	TemplateSettings() Configuration

	Setup()
	Destroy()
}

type addons map[string]Addon

// AddAddon
func AddAddon(a Addon) {
	if _, exists := Addons[a.Name()]; exists {
		logrus.WithField("_service", "addons_registrator").
			Warningf("replace addons %q settings", a.Name())
	}

	Addons[a.Name()] = a
}

func GetAddon(name string) Addon {

	return Addons[name]
}

// ListOfAddons
func ListOfAddons() addons {
	// TODO: Sort by priority

	return Addons
}
