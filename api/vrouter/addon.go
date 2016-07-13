package vrouter

import (
	"github.com/inpime/fader/api/addons"
	"github.com/labstack/echo"
)

var NAME = addonName

// TODO: make public
var addonName = "routing"
var version = "0.1.0"

func init() {
	// manual init
	addons.AddAddon(&Extension{})
}

type Extension struct {
}

func (Extension) Version() string {
	return version
}

func (Extension) Destroy() {

}

func (Extension) Name() string {
	return addonName
}

func (e Extension) Setup() {

}

func (Extension) TemplateSettings() addons.Configuration {
	return &Settings{&settings{}}
}

func (*Extension) Middlewares() []echo.MiddlewareFunc {

	return []echo.MiddlewareFunc{RouterMiddleware()}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (s *Extension) InjectTplAddons() error {
	tplContext()

	return nil
}
