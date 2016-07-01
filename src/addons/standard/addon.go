package standard

import (
	"api/config"
	"fmt"
	"github.com/labstack/echo"
)

var addonName = "fader.addons.standard"

var (
	ErrNotValidData = fmt.Errorf(addonName + ": not_valid_data")
)

func init() {
	// manual init
	config.AddExtension(&StandardExtension{})
}

type StandardExtension struct {
}

func (StandardExtension) Destroy() {

}

func (StandardExtension) Name() string {
	return addonName
}

func (*StandardExtension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*StandardExtension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (s *StandardExtension) InjectTplAddons() error {
	s.initTplContext()
	s.initTplFilters()
	s.initTplTags()

	return nil
}
