package standard

import (
	"api/config"
	"fmt"
	"github.com/labstack/echo"
)

var addonName = "fader.addons.standard"
var version = "0.1.0"

var (
	ErrNotValidData = fmt.Errorf(addonName + ": not_valid_data")
)

func init() {
	// manual init
	config.AddExtension(&Extension{})
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

func (s *Extension) InjectTplAddons() error {
	s.initTplContext()
	s.initTplFilters()
	s.initTplTags()

	return nil
}