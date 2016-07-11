package api

import (
	_ "addons/filestatic"
	_ "addons/importexport"
	_ "addons/payment"
	_ "addons/search"
	"addons/session"
	_ "addons/standard"
	"api/addons"
	"api/vrouter"

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

	apiutils "api/utils"
)

// InitAppSettings
func InitAppSettings() {
	config.Reload()
	vrouter.ReloadAppRouts()

	go apiutils.RefreshEvery(3*time.Second, func() {
		config.Reload()
		vrouter.ReloadAppRouts()
	})
}

// Run init tempaltes and start server
func Run() error {
	config.Init()

	var e = echo.New()

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 400 << 10, // 400 KB
	}))

	if logrus.GetLevel() >= logrus.InfoLevel {

		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `{"_service": "api", "time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
				`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
				`"latency_human":"${latency_human}","rx_bytes":${rx_bytes},` +
				`"tx_bytes":${tx_bytes}}` + "\n",
			Output: logrus.StandardLogger().Out,
		}))
	}

	// ------------------------
	// Special addons
	//  * 1. routing
	// 	* 2. session
	// ------------------------

	// ------------------------
	// 1. routing
	// ------------------------

	vrouterAddon := &vrouter.Extension{}
	config.Cfgx.AddConfig(
		vrouterAddon.Name(), // routing
		vrouterAddon.TemplateSettings())
	e.Use(vrouterAddon.Middlewares()...)
	vrouterAddon.Setup()
	vrouterAddon.RegEchoHandlers(AddSpecialHandler)
	vrouterAddon.InjectTplAddons()

	// ------------------------
	// 2. session
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
	config.Cfgx.AddConfig(
		sessionAddon.Name(), // session
		sessionAddon.TemplateSettings())
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
	// Registered addons (enterprise addons)
	// ------------------------

	for _, ext := range addons.ListOfAddons() {
		if vrouter.NAME == ext.Name() || session.NAME == ext.Name() {
			// игнорируем вручную установленные расширения
			continue
		}

		logrus.WithFields(logrus.Fields{
			"_service":      "api",
			"_target":       "initaddon",
			"addon":         ext.Name(),
			"addon_version": ext.Version(),
		}).Infof("add extension")

		config.Cfgx.AddConfig(
			ext.Name(), // component name
			ext.TemplateSettings())

		e.Use(ext.Middlewares()...)
		if logrus.GetLevel() >= logrus.InfoLevel {
			e.Use(debugMiddleware(ext.Name()))
		}
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

	// Init app settings
	InitAppSettings()

	e.Run(standard.New(config.Cfg.Address))

	return nil
}

func debugMiddleware(servicename string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			uri := ctx.Request().URI()

			logrus.WithFields(logrus.Fields{
				"_service": servicename,
				"uri":      uri,
			}).Debug("trace")

			return h(ctx)
		}
	}
}
