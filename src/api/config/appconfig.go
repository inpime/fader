package config

import (
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"store"
	"sync"
	"time"
	"utils"
)

var appSettings utils.M
var appSettingsMutex sync.Mutex

// settings@main
// settings@routs

var (
	PageCachingKey = "pageCaching"
	RoutsKey       = "routs"
	IncludeKey     = "include"
)

// AppSettings
func AppSettings() utils.M {
	return utils.M(appSettings)
}

// IsPageCaching
func IsPageCaching() bool {
	return AppSettings().Bool(PageCachingKey)
}

// AppRouts
func AppRouts() []string {
	return AppSettings().Strings(RoutsKey)
}

func AppSettingsIncludeFiles() []string {
	return AppSettings().Strings(IncludeKey)
}

// ReloadAppSettings reload app settings
func ReloadAppSettings() {
	file, err := store.LoadOrNewFile("settings", "main")

	_filename := "settings@main"

	if err != nil {
		logrus.WithField("_service", loggerKey).Errorf("load error file=%q, %v", _filename, err)
		return
	}

	routerMutex.Lock()
	appSettings = utils.Map() // clear the previous values
	defer routerMutex.Unlock()

	if _, err := toml.Decode(string(file.RawData().Bytes()), &appSettings); err != nil {
		logrus.Errorf("main settings: decode toml error, %v, %q", err, string(file.RawData().Bytes()))
		return
	}

	for _, includeFileName := range AppSettingsIncludeFiles() {
		includeFile, err := store.LoadOrNewFile("settings", includeFileName)

		if err != nil {
			logrus.WithField("_service", loggerKey).WithError(err).Info("find include file %q", includeFileName)
			continue
		}

		if _, err := toml.Decode(string(includeFile.RawData().Bytes()), &appSettings); err != nil {
			logrus.WithField("_service", loggerKey).Errorf("decode toml file=%q error, %v, %q", includeFileName, err, string(includeFile.RawData().Bytes()))
			return
		}
	}
}

// InitApp
func InitApp() {
	ReloadAppSettings()

	// TODO: synchronization with the previous launch
	go refreshEvery(3*time.Second, ReloadAppSettings)
}
