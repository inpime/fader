package payment

import (
	"github.com/inpime/fader/api/addons"
	"github.com/labstack/echo"
)

const (
	NAME    = "payment"
	VERSION = "0.1.0"
)

var ()

func init() {
	addons.AddAddon(&Extension{})
}

type Extension struct {
}

func (Extension) Version() string {
	return VERSION
}

func (Extension) Name() string {
	return NAME
}

func (s Extension) Destroy() {

}

func (Extension) Setup() {
}

func (s Extension) TemplateSettings() addons.Configuration {
	return &Settings{&settings{}}
}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (e *Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (*Extension) InjectTplAddons() error {
	initTplContext()
	return nil
}
