package vrouter

import (
	// "api/addons"
	"api/utils"
	"github.com/labstack/echo"
	"sync"
	"time"
)

var addonName = "fader.addons.vrouter"
var version = "0.1.0"

func init() {
	// manual init
	// addons.AddAddon(&Extension{})
}

type Extension struct {
	setupOnce sync.Once
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
	e.setupOnce.Do(func() {
		ReloadAppRouts()

		// TODO: synchronization with the previous launch
		go utils.RefreshEvery(3*time.Second, ReloadAppRouts)
	})
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
