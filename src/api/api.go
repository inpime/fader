// Copyright (c) Fader, OOO. All Rights Reserved.
// See LICENSE for license information.

package api

import (
	"addons"
	"api/router"
	"api/templates"
	"interfaces"
	"log"
	"os"
	"store/boltdb"
	"time"

	"github.com/boltdb/bolt"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	_       interfaces.Router = (*router.Router)(nil)
	logger  *log.Logger
	vrouter *router.Router

	settings *Settings
	config   *Config

	fileLoaderForRouting interfaces.FileLoader

	fileManager   interfaces.FileManager
	bucketManager interfaces.BucketManager
	inMemStore    *InMemoryStore

	db *bolt.DB

	TESTING bool

	ConfigBucketName   = "settings"
	MainConfigFileName = "main.toml"
)

// Setup
func Setup(e *echo.Echo, _settings *Settings) error {
	settings = SettingsOrDefault(_settings)

	// Logger -----------------------------------------------------------------

	logger = log.New(os.Stderr, "[api]", 1)
	logger.Printf("init... % v\n", settings)

	// Database ---------------------------------------------------------------
	var err error
	db, err = bolt.Open(settings.DatabasePath, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})

	if err != nil {
		logger.Println("[ERR] setup database ", err)
		return err
	}

	// Sys components -------------------------------------------------------------

	bucketManager = boltdb.NewBucketManager(db)
	fileManager = boltdb.NewFileManager(db)
	inMemStore = NewInMemoryStore(db)

	templates.DefaultTemplatesLoader = interfaces.NewTemplatesStore(fileManager)
	templates.SetupSettings()

	// for routings: file controller
	// only used lua script and meta informations
	fileLoaderForRouting = NewFileProvider(
		fileManager,
		interfaces.FileWithoutRawData,
	)

	// App compoenents ------------------------------------------------

	logger.Println("init... manager routes")
	vrouter = router.NewRouter()

	// Setup init settings
	logger.Println("init... app check")
	if !TESTING {
		if err := InitFirstRunIfNeed(); err != nil {
			logger.Fatalln("[FAIL] installation of first run:", err)
		}
	} else {
		logger.Println("\t skiped check first run")
	}

	// Setup app config
	logger.Println("init... app config")
	config = newConfig()
	logger.Println("\t  run auto update settings every 10 seconds")
	appConfigUpdateFn()
	go RefreshEvery(time.Second*10, appConfigUpdateFn)

	// Setup app addons
	logger.Println("init... app addons")
	if err := SetupAddons(); err != nil {
		logger.Println("[ERR] stup addons", err)
	}
	for _, addon := range addons.Addons {
		logger.Println("\t addon:", addon.Name())
		if err := addon.ExtContextPongo2(pongo2.DefaultSet.Globals); err != nil {
			logger.Printf("\t addon %q, ext. contoext err: %s", addon.Name(), err)
		}

		if err := addon.ExtTagsFiltersPongo2(
			pongo2.RegisterFilter,
			pongo2.ReplaceFilter,
			pongo2.RegisterTag,
			pongo2.ReplaceTag,
		); err != nil {
			logger.Printf("\t addon %q, ext. filters/tags err: %s", addon.Name(), err)
		}
	}

	// Routes -----------------------------------------------------

	logger.Println("init... application middlewares")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(router.VRouterMiddleware(vrouter))

	logger.Println("init... application routes")
	e.GET("*", FaderHandler)
	e.POST("*", FaderHandler)

	logger.Println("init... done")

	return nil
}
