package api

import (
	"api/router"
	"api/templates"
	"interfaces"
	"log"
	"os"
	"store/boltdb"
	"time"

	"github.com/boltdb/bolt"
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

	db *bolt.DB
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

	// Components -------------------------------------------------------------

	bucketManager = boltdb.NewBucketManager(db)
	fileManager = boltdb.NewFileManager(db)

	templates.DefaultTemplatesLoader = interfaces.NewTemplatesStore(fileManager)
	templates.SetupSettings()

	// for routings: file controller
	// only used lua script and meta informations
	fileLoaderForRouting = NewFileProvider(
		fileManager,
		interfaces.FileWithoutRawData,
	)

	// Application compoenents ------------------------------------------------

	logger.Println("init... manager routes")
	vrouter = router.NewRouter()

	// Setup init settings
	logger.Println("init... app check")
	if err := InitFirstRunIfNeed(); err != nil {
		logger.Fatalln("[FAIL] installation of first run:", err)
	}

	// TODO: setup app config
	config = newConfig()
	/*
		1. Routers
	*/

	// Application routes -----------------------------------------------------

	logger.Println("init... application middlewares")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(router.VRouterMiddleware(vrouter))

	logger.Println("init... application routes")
	e.Get("*", FaderHandler)

	logger.Println("init... done")

	return nil
}
