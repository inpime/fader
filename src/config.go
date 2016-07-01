package main

import (
	"api"
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

	api.Cfg = &api.Config{
		Address:       address,
		WorkspacePath: workspace,
		Search: api.SearchStore{
			Host:      esAddress,
			IndexName: esIndexName,
		},
		Session: api.ApiSessionConfig{
			Path: "/",

			SecretKey:   sessionSecureKey,
			BucketName:  "sessions",
			SessionName: "fds",

			Store: api.StoreConfig{
				Provider:       "boltdb",
				BoltDBFilePath: filepath.Clean(workspace + string(filepath.Separator) + "session.db"),
			},
		},

		Store: api.AppStoreConfig{
			StoreConfig: api.StoreConfig{
				Provider:       "boltdb",
				BoltDBFilePath: filepath.Clean(workspace + string(filepath.Separator) + "app.db"),
			},
			StaticPath: filepath.Clean(workspace + string(filepath.Separator) + "static/"),
		},
	}

	store.ElasticSearchIndexName = api.Cfg.Search.IndexName
	store.WorkspacePath = api.Cfg.WorkspacePath

	utils.EnsureDir(api.Cfg.WorkspacePath)
	utils.EnsureDir(api.Cfg.Session.Store.BoltDBFilePath)

	utils.EnsureDir(api.Cfg.Store.StaticPath)

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
