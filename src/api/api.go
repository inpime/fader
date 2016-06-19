package api

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	// "github.com/labstack/echo/middleware"
	"github.com/yosssi/boltstore/reaper"
)

func EmptyHandler(c echo.Context) error {

	return c.String(200, "empty")
}

// Run init tempaltes and start server
func Run() {

	// TODO: make easier the init

	// init templates
	initTemplates()

	if err := InitSession(); err != nil {
		panic(err)
	}

	initWidgets()

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

	e.Use(constructSessionMiddleware, MiddlewareInitSession())
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.Get("/", ExecuteWidget)
	e.Get("/*", ExecuteWidget)
	e.Post("/*", ExecuteWidget)

	// TODO: e.Get("/content/*", UserContentEntryPointHandler)

	fmt.Printf("api: listener %q\n", Cfg.Address)

	e.Run(standard.New(Cfg.Address))
}
