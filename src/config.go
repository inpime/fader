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

func initConfig() {
	flag.Parse()

	workspace := *flagWorkspacePath

	api.Cfg = &api.Config{
		Address:       "192.168.1.36:3322",
		WorkspacePath: workspace,
		Search: api.SearchStore{
			Host:      "https://es.idheap.com",
			IndexName: "fader",
		},
		Session: api.ApiSessionConfig{
			Path: "/",

			SecretKey:  "secure-key",
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
