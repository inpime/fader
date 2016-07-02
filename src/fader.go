package main

import (
	"api"
	"api/config"
	"github.com/Sirupsen/logrus"
)

func main() {
	logrus.Info("Fader starting...")

	logrus.Info("Init config...")
	initConfig()

	logrus.Info("Init elasticsearch...")
	initElasticSearch()

	logrus.Debugf("Current config: %#v", config.Cfg)

	logrus.Info("Init stores...")
	initStroes()

	logrus.Info("Api...")
	api.Run()
}
