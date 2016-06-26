package main

import (
	"api"
	"github.com/Sirupsen/logrus"
)

func main() {
	logrus.Info("Fader starting...")

	logrus.Info("Init config...")
	initConfig()

	logrus.Info("Init elasticsearch...")
	initElasticSearch()

	logrus.Debugf("Current config: %#v", api.Cfg)

	logrus.Info("Init stores...")
	initStroes()

	logrus.Info("Api...")
	api.Run()
}
