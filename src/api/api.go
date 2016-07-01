package api

import (
	_ "addons/filestatic"
	_ "addons/search"
	"addons/session"
	_ "addons/standard"

	"api/config"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"path/filepath"
	"store"
	"time"
	"utils"
)

// Run init tempaltes and start server
func Run() error {
	config.Init()

	var e = echo.New()

	// ------------------------
	// Special addons
	// 	* session
	// 	* logger
	// ------------------------

	// ------------------------
	// Addon: session
	// ------------------------

	if Cfg.Session.Store.Provider != "boltdb" {
		return fmt.Errorf("not supported session store %s", Cfg.Session.Store)
	}
	utils.EnsureDir(filepath.Dir(Cfg.Session.Store.BoltDBFilePath))
	db, err := bolt.Open(Cfg.Session.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Init Default guest user session
	file, err := store.LoadOrNewFile(config.UsersBucketName, config.GuestUserFileName)
	if err != nil {
		return err
	}

	session.DefaultGuestSession = file

	sessionAddon := &session.SessionExtension{}
	sessionAddon.SetAppConfig(session.Config{
		Path: Cfg.Session.Path,

		DB:       db,
		HttpOnly: Cfg.Session.HttpOnly,
		Secure:   Cfg.Session.Secure,
		Domain:   Cfg.Session.Domain,
		MaxAge:   Cfg.Session.MaxAge,

		BucketName:  Cfg.Session.BucketName,
		SecretKey:   Cfg.Session.SecretKey,
		SessionName: Cfg.Session.SessionName,
	})
	logrus.WithField("_service", "api").Infof("add extension: %q", sessionAddon.Name())
	e.Use(sessionAddon.Middlewares()...)
	sessionAddon.RegEchoHandlers(AddSpecialHandler)
	sessionAddon.InjectTplAddons()

	// ------------------------
	// Addon: logger
	// ------------------------

	if logrus.GetLevel() >= logrus.InfoLevel {
		// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// 	Format: `{"_service": "api", "time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
		// 		`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
		// 		`"latency_human":"${latency_human}","rx_bytes":${rx_bytes},` +
		// 		`"tx_bytes":${tx_bytes}}` + "\n",
		// 	Output: logrus.StandardLogger().Out,
		// }))
		e.Use(middleware.Logger())

		e.Use(middleware.Recover())
	}

	// ------------------------
	// Registered addons
	// ------------------------

	for _, ext := range config.ListOfExtensions() {
		logrus.WithField("_service", "api").Infof("add extension: %q", ext.Name())
		// ext.SetLogger(logrus.StandardLogger().Out)
		e.Use(ext.Middlewares()...)
		ext.RegEchoHandlers(AddSpecialHandler)
		ext.InjectTplAddons()
	}

	// ------------------------

	// ------------------------
	// App routs
	// ------------------------

	e.Get("/", AppEntryPointHandler)
	e.Get("/*", AppEntryPointHandler)
	e.Post("/*", AppEntryPointHandler)
	e.Put("/*", AppEntryPointHandler)
	e.Delete("/*", AppEntryPointHandler)

	logrus.WithField("_service", "api").Infof("HTTP Listener %q", Cfg.Address)

	e.Run(standard.New(Cfg.Address))

	return nil
}
