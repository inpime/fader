package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/inpime/fader/api/addons"
	"github.com/inpime/fader/store"
	"strings"
	"sync"
)

var Cfgx *configs
var cfgxMutex sync.RWMutex

var (
	ErrConfigNameUse = errors.New("that config name is in use")
)

type configs map[string]addons.Configuration

func (c *configs) AddConfig(name string, conf addons.Configuration) error {
	cfgxMutex.Lock()
	defer cfgxMutex.Unlock()

	if _, exists := (*c)[name]; exists {
		logrus.WithField("_service", loggerKey).
			Warningf("replace config %q settings", name)
	}

	(*c)[name] = conf

	return nil
}

func (c configs) MergeConfig(name string, _cfg addons.Configuration) error {
	cfg := c.Config(name)

	if cfg == nil {
		logrus.WithField("_service", loggerKey).
			Warningf("not found config %q", name)

		return fmt.Errorf("not found config %q", name)
	}

	return cfg.Merge(_cfg)
}

func (c configs) Config(name string) addons.Configuration {

	cfgxMutex.RLock()
	defer cfgxMutex.RUnlock()

	return c[name]
}

func (c *configs) Reset(_c *configs) {
	cfgxMutex.RLock()
	defer cfgxMutex.RUnlock()
	*c = *_c
}

// Reload update app settings
func Reload() {

	newconfig := NewConfig()

	// init extensions settings

	for addonName := range *Cfgx {

		if addonName == sectionName {
			continue
		}

		newconfig.AddConfig(addonName, addons.GetAddon(addonName).TemplateSettings())
	}

	fileName := MainSettingsFileName
	file, err := store.LoadOrNewFile(SettingsBucketName, fileName)

	src := string(file.RawData().Bytes())
	// src = `[main]

	// include = ["console.route"]
	// `

	logrus.WithField("_service", loggerKey).
		Debugf("settings bucket %q", SettingsBucketName)
	logrus.WithField("_service", loggerKey).
		Debugf("source app settings:\n############\n%v\n############\n", src)

	if err != nil {
		logrus.WithError(err).
			WithField("_service", loggerKey).
			Errorf("load main settings %q", fileName)
		return
	}

	hydrateAllAppConfigs(newconfig, src, fileName)

	includeFiles := newconfig.Config(sectionName).(addons.MainConfiguration).Include()

	for _, fileName := range includeFiles {
		// include
		logrus.WithField("_service", loggerKey).
			Debug("load include settings from file %q", fileName)

		file, err := store.LoadOrNewFile(SettingsBucketName, fileName)
		src := strings.Replace(string(file.RawData().Bytes()), "[[routs]]", "[[routing.routs]]", -1)

		logrus.WithField("_service", loggerKey).
			Debugf("source include %q settings:\n############\n%v\n############\n", fileName, src)

		if err != nil {
			logrus.WithError(err).
				WithField("_service", loggerKey).
				Errorf("load include file %q for addon settings", fileName)
			continue
		}

		mergeAllAppConfigs(newconfig, src, fileName)
	}

	Cfgx.Reset(newconfig)
}

func hydrateAllAppConfigs(c *configs, src, fileName string) {
	for addonName := range *c {
		logrus.WithField("_service", loggerKey).
			Debugf("settings for addon %q", addonName)

		if err := hydratorTOML(src, c.Config(addonName)); err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"_service":  loggerKey,
					"_filename": fileName,
				}).
				Errorf("decode settings addon %q", addonName)

			continue
		}
	}
}

func mergeAllAppConfigs(c *configs, src, fileName string) {
	for addonName := range *c {
		logrus.WithField("_service", loggerKey).
			Debugf("merge settings for addon %q", addonName)

		var addonConfig addons.Configuration

		if addonName == sectionName {
			// TODO: refactoring main config
			addonConfig = &Settings{&settings{}}
		} else {
			addonConfig = addons.GetAddon(addonName).TemplateSettings()
		}

		if err := hydratorTOML(src, addonConfig); err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"_service":  loggerKey,
					"_filename": fileName,
				}).
				Errorf("decode settings addon %q for merge", addonName)

			continue
		}

		if err := c.MergeConfig(addonName, addonConfig); err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"_service":  loggerKey,
					"_filename": fileName,
				}).
				Errorf("merge settings addon %q", addonName)
		}

	}
}

func hydratorTOML(src string, i interface{}) error {
	if _, err := toml.Decode(src, i); err != nil {
		return err
	}

	return nil
}

// func hydratorJSON(src string, i interface{}) error {
// 	return nil
// }
