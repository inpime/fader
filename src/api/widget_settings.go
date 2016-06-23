package api

import (
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"store"
	"sync"
	"time"
	"utils"
)

var router *Router
var routerMutex sync.Mutex

var appSettings utils.M
var appSettingsMutex sync.Mutex

// IsPageCaching
func IsPageCaching() bool {
	return appSettings.Bool("pageCaching")
}

// AppRouts
func AppRouts() []string {
	return appSettings.Strings("routs")
}

type Rout struct {
	Path      string   `toml:"path"`
	Name      string   `toml:"name"`
	Handler   string   `toml:"handler"`
	Methods   []string `toml:"methods"`
	Licenses  []string `toml:"licenses"`
	IsSpecial bool     `toml:"special"`
}

type Routing struct {
	Mode  string `toml:"mode"`
	Routs []Rout `toml:"routs"`
}

func initWidgetVirtualRouts() {
	reloadAppSettings()
	reloadAppRouts()

	RegistedSpecialHandler(FileContentByNameSpecialHandlerName, FileContentByName_SpecialHandler)
	RegistedSpecialHandler(FileContentByIDSpecialHandlerName, FileContentByID_SpecialHandler)

	RegistedSpecialHandler("exportapp", ExportAppHandler)
	RegistedSpecialHandler("importapp", ImportAppHandler)

	go RefreshEvery(3*time.Second, reloadAppSettings)
	go RefreshEvery(3*time.Second, reloadAppRouts)
}

func RefreshEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func reloadAppSettings() {
	file, err := store.LoadOrNewFile(SettingsBucketName, MainSettingsFileName)

	if err != nil {
		logrus.Errorf("main settings: load error, %v", err)
		return
	}

	routerMutex.Lock()
	appSettings = utils.Map() // clear the previous values
	defer routerMutex.Unlock()

	if _, err := toml.Decode(string(file.RawData().Bytes()), &appSettings); err != nil {
		logrus.Errorf("main settings: decode toml error, %v, %q", err, string(file.RawData().Bytes()))
		return
	}

	// logrus.WithField("settings", appSettings).Info("main settings")
}

func reloadAppRouts() {
	routerMutex.Lock()
	defer routerMutex.Unlock()

	router = NewRouter()

	for _, fileName := range AppRouts() {
		updateAppRoutes(fileName)
	}
}

func updateAppRoutes(fileName string) {

	file, err := store.LoadOrNewFile(SettingsBucketName, fileName)

	if err != nil {
		logrus.Errorf("vrouting: load %q error, %v", fileName, err)
		return
	}

	var routing = &Routing{}

	if _, err := toml.Decode(string(file.RawData().Bytes()), routing); err != nil {
		logrus.Errorf("vrouting: %q decode toml error, %v, %q", fileName, err, string(file.RawData().Bytes()))
		return
	}

	for _, _r := range routing.Routs {

		handler := NewHandlerFromRoute(_r)

		if len(_r.Methods) == 0 {
			router.Handle(_r.Path, handler).Methods("GET").Name(_r.Name)
		} else {
			router.Handle(_r.Path, handler).Methods(_r.Methods...).Name(_r.Name)
		}
	}
}
