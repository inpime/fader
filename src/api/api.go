package api

import (
	"api/config"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/yosssi/boltstore/reaper"
)

// Run init tempaltes and start server
func Run() {
	config.Init()

	pongo2InitGlobalCustoms()
	pongo2InitAddons()

	// TODO: make easier the init

	// init templates
	// initTemplates()

	if err := InitSession(); err != nil {
		panic(err)
	}

	initWidgetVirtualRouts()

	// Init session

	constructSessionMiddleware := MiddlewareSessionWithConfig(DefaultSessionName, SessionConfig{
		Path:       Cfg.Session.Path,
		BucketName: Cfg.Session.BucketName,
		SecretKey:  Cfg.Session.SecretKey,
		DB:         sessionDb,
		HttpOnly:   Cfg.Session.HttpOnly,
		Secure:     Cfg.Session.Secure,
		Domain:     Cfg.Session.Domain,
		MaxAge:     Cfg.Session.MaxAge,
	})

	defer reaper.Quit(reaper.Run(sessionDb, reaper.Options{
		BucketName: []byte(Cfg.Session.BucketName),
	}))

	// ------------------
	// Init http server
	// ------------------

	var e = echo.New()

	// logger
	if logrus.GetLevel() >= logrus.InfoLevel {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `{"_service": "api", "time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
				`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
				`"latency_human":"${latency_human}","rx_bytes":${rx_bytes},` +
				`"tx_bytes":${tx_bytes}}` + "\n",
			Output: logrus.StandardLogger().Out,
		}))

		e.Use(middleware.Recover())
	}

	e.Use(constructSessionMiddleware, MiddlewareInitSession())

	e.Get("/", ExecuteWidget)
	e.Get("/*", ExecuteWidget)
	e.Post("/*", ExecuteWidget)

	// TODO: e.Get("/content/*", UserContentEntryPointHandler)

	logrus.Infof("Api listener %q", Cfg.Address)

	e.Run(standard.New(Cfg.Address))
}
