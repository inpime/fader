package session

import (
	// "github.com/flosch/pongo2"
	// "api/config"
	"github.com/labstack/echo"
	"github.com/yosssi/boltstore/reaper"
	"net/http"
	"time"
)

var addonName = "fader.addons.session"
var version = "0.1.0"

func init() {
	// manual init
	// config.AddExtension(&Extension{})
}

type Extension struct {
	config Config

	reaperQuitC chan<- struct{}
	reaperDoneC <-chan struct{}
}

func (Extension) Version() string {
	return version
}

func (Extension) Name() string {
	return addonName
}

func (s Extension) Destroy() {
	reaper.Quit(s.reaperQuitC, s.reaperDoneC)
}

func (s *Extension) SetAppConfig(config Config) {
	s.config = config

	// cleaner expiring sessions
	s.reaperQuitC, s.reaperDoneC = reaper.Run(s.config.DB, reaper.Options{
		BucketName:    []byte(config.BucketName),
		BatchSize:     100,             // TODO: move in the system settings
		CheckInterval: time.Minute * 1, // TODO: move in the system settings
	})
}

func (s *Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		SessionStoreMiddleware(s.config.SessionName, s.config),
		InitializerUserSessionMiddleware(),
	}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(addonName+".logout_handler", func(ctx echo.Context) error {

		if _session := GetSession(ctx); _session != nil {
			_session.Logout()
		}

		return ctx.NoContent(http.StatusOK)
	})
}

func (*Extension) InjectTplAddons() error {
	return nil
}
