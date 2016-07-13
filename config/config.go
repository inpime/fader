package config

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/inpime/fader/addons/importexport"
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/store"
	"github.com/inpime/fader/utils"
	"os"
	"path/filepath"
)

var BuildDate = ""
var BuildHash = ""
var Version = ""

var (
	defaultWorkspace           = "./_workspace"
	defaultHttpListenerAddress = ":3322"
	defaultESAddress           = ":9200"
	defaultESIndexName         = "fader"
	defaultSessionSecretKey    = "very secret"
	defaultLoggerMode          = "error"

	flagWorkspacePath = flag.String("workspace",
		defaultWorkspace, "")
	flagHttpListener = flag.String("listener",
		defaultHttpListenerAddress, "")
	flagESAddress = flag.String("es_address",
		defaultESAddress, "")
	flagESIndex = flag.String("es_index",
		defaultESIndexName, "")
	flagESSessionSecret = flag.String("session_secret",
		defaultSessionSecretKey, "")
	flagLoggerMode = flag.String("logger",
		defaultLoggerMode, "")
	flagAppFaderSetup = flag.String("startapp",
		"", "")
)

const (
	ENV_WORKSPACE             = "FADER_WORKSPACE"
	ENV_HTTPLISTENER          = "FADER_HTTPLISTENER"
	ENV_ESADDRESS             = "FADER_ESADDRESS"
	ENV_ESINDEX               = "FADER_ESINDEX"
	ENV_SESSIONSECRET         = "FADER_SESSIONSECRET"
	ENV_LOGGERMODE            = "FADER_LOGGERMODE"
	ENV_FADER_ESADDRESSDOCKER = "FADER_ESADDRESSDOCKER"
	ENV_FADER_APPSETUP        = "FADER_APPSETUP"
)

func Init() {
	workspace := os.Getenv(ENV_WORKSPACE)
	if len(workspace) == 0 {
		workspace = *flagWorkspacePath
	}

	httpListenerAddress := os.Getenv(ENV_HTTPLISTENER)
	if len(httpListenerAddress) == 0 {
		httpListenerAddress = *flagHttpListener
	}

	esAddress := os.Getenv(ENV_ESADDRESS)
	if len(esAddress) == 0 {
		esAddress = *flagESAddress
	}

	esIndexName := os.Getenv(ENV_ESINDEX)
	if len(esIndexName) == 0 {
		esIndexName = *flagESIndex
	}

	sessionSecretKey := os.Getenv(ENV_SESSIONSECRET)
	if len(sessionSecretKey) == 0 {
		sessionSecretKey = *flagESSessionSecret
	}

	loggerModeStr := os.Getenv(ENV_LOGGERMODE)
	if len(loggerModeStr) == 0 {
		loggerModeStr = *flagLoggerMode
	}

	AppFaderSetup := os.Getenv(ENV_FADER_APPSETUP)
	if len(AppFaderSetup) == 0 {
		AppFaderSetup = *flagAppFaderSetup
	}

	// override
	importexport.ArchiveURLLatestVersion = AppFaderSetup

	config.Cfg = &config.Config{
		AppVersion:   Version,
		AppBuildDate: BuildDate,
		AppBuildHash: BuildHash,

		Address:       httpListenerAddress,
		WorkspacePath: workspace,
		Search: config.SearchStore{
			Host:      esAddress,
			IndexName: esIndexName,
		},
		Session: config.ApiSessionConfig{
			Path: "/",

			SecretKey:   sessionSecretKey,
			BucketName:  "sessions",
			SessionName: "fds",

			Store: config.StoreConfig{
				Provider:       "boltdb",
				BoltDBFilePath: filepath.Clean(workspace + string(filepath.Separator) + "session.db"),
			},
		},

		Store: config.AppStoreConfig{
			StoreConfig: config.StoreConfig{
				Provider:       "boltdb",
				BoltDBFilePath: filepath.Clean(workspace + string(filepath.Separator) + "app.db"),
			},
			StaticPath: filepath.Clean(workspace + string(filepath.Separator) + "static/"),
		},
	}

	//
	pwd, _ := os.Getwd()
	logrus.WithField("ref", "config").Info("Version:", config.Cfg.Version())
	logrus.WithField("ref", "config").Info("Work dir:", pwd)

	if err := utils.EnsureDir(config.Cfg.WorkspacePath); err != nil {
		logrus.WithField("ref", "config").Fatal(err)
	}
	if err := utils.EnsureDir(config.Cfg.Session.Store.BoltDBFilePath); err != nil {
		logrus.WithField("ref", "config").Fatal(err)
	}
	if err := utils.EnsureDir(config.Cfg.Store.StaticPath); err != nil {
		logrus.WithField("ref", "config").Fatal(err)
	}

	// store enviroment
	store.ElasticSearchIndexName = config.Cfg.Search.IndexName
	store.WorkspacePath = config.Cfg.WorkspacePath

	// logger
	var loggerLevel, err = logrus.ParseLevel(loggerModeStr)

	if err != nil {
		logrus.WithField("ref", "config").
			Warning("error parse level logger: %s", err)
		loggerLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(loggerLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
}
