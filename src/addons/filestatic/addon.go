package filestatic

import (
	"api/addons"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

const (
	NAME    = "filestatic"
	VERSION = "0.1.0"
)

var (
	// route or special handler name
	ByNameRouteName = NAME + ".byname"
	ByIDRouteName   = NAME + ".byid"
)

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
	return Settings{&settings{}}
}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (e *Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(ByNameRouteName, FileContentByNameHandler)
	fnReg(ByIDRouteName, FileContentByNameHandler)
}

func (*Extension) InjectTplAddons() error {
	pongo2.RegisterFilter("fc", filterUrlFileByName)
	pongo2.RegisterFilter("filecontenturl", filterUrlFileByName)
	pongo2.RegisterFilter("urlfile", filterUrlFileByName) // OLD
	return nil
}
