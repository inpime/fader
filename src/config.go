package main

import (
	"api/config"
	"flag"
	"github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"store"
	"utils"
)

var flagWorkspacePath = flag.String("workspace", "./_workspace/", "workspace application")
var flagApiAddress = flag.String("address", ":3322", "path to the static folder")
var flagElasticSearchAddress = flag.String("es_address", "https://es.idheap.com", "url elasticsearch")
var flagElasticSearchIndex = flag.String("es_index", "fader", "elasticsearch index name")
var flagSessionSecret = flag.String("session_secret", "secure-key", "sessions secure key")
var flagLoggerMode = flag.String("mode", "info", "logger mode")

func initConfig() {
	flag.Parse()

	workspace := os.Getenv("FADER_WORKSPACE")
	if len(workspace) == 0 {
		workspace = *flagWorkspacePath
	}

	address := os.Getenv("FADER_API_ADDR")
	if len(address) == 0 {
		address = *flagApiAddress
	}

	esAddress := os.Getenv("FADER_ES_ADDR")
	if len(esAddress) == 0 {
		esAddress = *flagElasticSearchAddress
	}

	esIndexName := os.Getenv("FADER_ES_INDEX")
	if len(esIndexName) == 0 {
		esIndexName = *flagElasticSearchIndex
	}

	sessionSecureKey := os.Getenv("FADER_SESSION_SECRET")
	if len(sessionSecureKey) == 0 {
		sessionSecureKey = *flagSessionSecret
	}

	// TODO: if empty values then stop

	config.Cfg = &config.Config{
		Address:       address,
		WorkspacePath: workspace,
		Search: config.SearchStore{
			Host:      esAddress,
			IndexName: esIndexName,
		},
		Session: config.ApiSessionConfig{
			Path: "/",

			SecretKey:   sessionSecureKey,
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

	store.ElasticSearchIndexName = config.Cfg.Search.IndexName
	store.WorkspacePath = config.Cfg.WorkspacePath

	utils.EnsureDir(config.Cfg.WorkspacePath)
	utils.EnsureDir(config.Cfg.Session.Store.BoltDBFilePath)

	utils.EnsureDir(config.Cfg.Store.StaticPath)

	// Config logger

	appLoggerStr := os.Getenv("FADER_MODE")
	if len(appLoggerStr) == 0 {
		appLoggerStr = *flagLoggerMode
	}

	var appLoggerLevel, err = logrus.ParseLevel(appLoggerStr)

	if err != nil {
		appLoggerLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(appLoggerLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	// logrus.SetLevel(logrus.DebugLevel)
}
