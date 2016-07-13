package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/inpime/fader/api"
	"github.com/inpime/fader/config"
)

func main() {
	flag.Parse()

	logrus.WithField("ref", "main").Info("Start")

	logrus.WithField("ref", "main").Info("Init config")
	config.Init()

	logrus.WithField("ref", "main").Info("Init store")
	if err := api.Init(); err != nil {
		logrus.WithField("ref", "main").Fatal(err)
	}

	logrus.WithField("ref", "main").Info("Run")
	if err := api.Run(); err != nil {
		logrus.WithField("ref", "main").Fatal(err)
	}

	logrus.WithField("ref", "main").Info("Stop")
}
