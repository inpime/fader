package api

import (
	_ "addons/filestatic"
	_ "addons/importexport"
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
	// 	* 1. session
	// 	* 2. logger
	// ------------------------

	// ------------------------
	// 1. session
	// ------------------------

	if config.Cfg.Session.Store.Provider != "boltdb" {
		return fmt.Errorf("not supported session store %s", config.Cfg.Session.Store)
	}
	utils.EnsureDir(filepath.Dir(config.Cfg.Session.Store.BoltDBFilePath))
	db, err := bolt.Open(config.Cfg.Session.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Init Default guest user session
	file, err := store.LoadOrNewFile(config.UsersBucketName, config.GuestUserFileName)
	if err != nil {
		return err
	}

	session.DefaultGuestSession = file

	sessionAddon := &session.Extension{}
	sessionAddon.SetAppConfig(session.Config{
		Path: config.Cfg.Session.Path,

		DB:       db,
		HttpOnly: config.Cfg.Session.HttpOnly,
		Secure:   config.Cfg.Session.Secure,
		Domain:   config.Cfg.Session.Domain,
		MaxAge:   config.Cfg.Session.MaxAge,

		BucketName:  config.Cfg.Session.BucketName,
		SecretKey:   config.Cfg.Session.SecretKey,
		SessionName: config.Cfg.Session.SessionName,
	})

	logrus.WithFields(logrus.Fields{
		"_service":      "api",
		"_target":       "initaddon",
		"addon":         sessionAddon.Name(),
		"addon_version": sessionAddon.Version(),
	}).Infof("add extension")

	e.Use(sessionAddon.Middlewares()...)
	sessionAddon.RegEchoHandlers(AddSpecialHandler)
	sessionAddon.InjectTplAddons()

	// ------------------------
	// 2. logger
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
	// Registered addons (enterprise addons)
	// ------------------------

	for _, ext := range config.ListOfExtensions() {
		logrus.WithFields(logrus.Fields{
			"_service":      "api",
			"_target":       "initaddon",
			"addon":         ext.Name(),
			"addon_version": ext.Version(),
		}).Infof("add extension")
		// ext.SetLogger(logrus.StandardLogger().Out)
		e.Use(ext.Middlewares()...)
		ext.RegEchoHandlers(AddSpecialHandler)
		ext.InjectTplAddons()
	}

	// ------------------------
	// App routs
	// ------------------------

	e.Get("/", AppEntryPointHandler)
	e.Get("/*", AppEntryPointHandler)
	e.Post("/*", AppEntryPointHandler)
	e.Put("/*", AppEntryPointHandler)
	e.Delete("/*", AppEntryPointHandler)

	logrus.WithFields(logrus.Fields{
		"_service": "api",
		"_target":  "httplistener",
		"address":  config.Cfg.Address,
	}).Infof("Run API HTTP Listener")

	e.Run(standard.New(config.Cfg.Address))

	return nil
}
