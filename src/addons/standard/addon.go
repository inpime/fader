package standard

import (
	"api/addons"
	"fmt"
	"github.com/labstack/echo"
	"utils"
)

var addonName = "basic"
var version = "0.1.0"

var (
	ErrNotValidData = fmt.Errorf(addonName + ": not_valid_data")
)

func init() {
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

func (*Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (Extension) Setup() {
}

func (s Extension) TemplateSettings() addons.Configuration {
	return Settings{&settings{
		Config: utils.Map(),
	}}
}

func (s *Extension) InjectTplAddons() error {
	s.initTplContext()
	s.initTplFilters()
	s.initTplTags()

	return nil
}
