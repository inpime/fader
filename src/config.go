package main

import (
	"api"
	"flag"
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"store"
	"utils"
)

var flagWorkspacePath = flag.String("workspace", "./_workspace/", "workspace application")
var flagApiAddress = flag.String("address", ":3322", "path to the static folder")
var flagElasticSearchAddress = flag.String("es_address", "https://es.idheap.com", "url elasticsearch")
var flagElasticSearchIndex = flag.String("es_index", "fader", "elasticsearch index name")
var flagSessionSecret = flag.String("session_secret", "secure-key", "sessions secure key")

func initConfig() {
	flag.Parse()

	workspace := *flagWorkspacePath
	address := *flagApiAddress
	esAddress := *flagElasticSearchAddress
	esIndexName := *flagElasticSearchIndex
	sessionSecureKey := *flagSessionSecret

	api.Cfg = &api.Config{
		Address:       address,
		WorkspacePath: workspace,
		Search: api.SearchStore{
			Host:      esAddress,
			IndexName: esIndexName,
		},
		Session: api.ApiSessionConfig{
			Path: "/",

			SecretKey:  sessionSecureKey,
			BucketName: "sessions",

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

	logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetLevel(logrus.DebugLevel)
}
