package config

import (
	apiutils "api/utils"
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"store"
	"time"
	"utils"
)

var appSettings utils.M

var (
	// Section names of file settings@main

	PageCachingKey = "pageCaching"
	IncludeKey     = "include"
)

// AppSettings
func AppSettings() utils.M {
	if appSettings == nil {
		// TODO: default setting (safe mode)
		appSettings = utils.Map()
	}

	return appSettings
}

// IsPageCaching
func IsPageCaching() bool {
	return AppSettings().Bool(PageCachingKey)
}

func AppSettingsIncludeFiles() []string {
	return AppSettings().Strings(IncludeKey)
}

// ReloadAppSettings reload app settings
func ReloadAppSettings() {
	file, err := store.LoadOrNewFile(SettingsBucketName, MainSettingsFileName)

	_filename := SettingsBucketName + "@" + MainSettingsFileName

	if err != nil {
		logrus.WithField("_service", loggerKey).Errorf("load error file=%q, %v", _filename, err)
		return
	}

	newAppSettings := utils.Map() // TODO: default setting (safe mode)

	if _, err := toml.Decode(string(file.RawData().Bytes()), &newAppSettings); err != nil {
		logrus.Errorf("main settings: decode toml error, %v, %q", err, string(file.RawData().Bytes()))
		return
	}

	for _, includeFileName := range AppSettingsIncludeFiles() {
		includeFile, err := store.LoadOrNewFile(SettingsBucketName, includeFileName)

		if err != nil {
			logrus.WithField("_service", loggerKey).WithError(err).Info("find include file %q", includeFileName)
			continue
		}

		if _, err := toml.Decode(string(includeFile.RawData().Bytes()), &newAppSettings); err != nil {
			logrus.WithField("_service", loggerKey).Errorf("decode toml file=%q error, %v, %q", includeFileName, err, string(includeFile.RawData().Bytes()))
			return
		}
	}

	appSettings = newAppSettings
}

// InitApp
func InitApp() {
	ReloadAppSettings()

	// TODO: synchronization with the previous launch
	go apiutils.RefreshEvery(3*time.Second, ReloadAppSettings)
}
