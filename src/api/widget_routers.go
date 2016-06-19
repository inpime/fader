package api

import (
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"store"
	"sync"
	"time"
)

var router *Router
var routerMutex sync.Mutex

type Rout struct {
	Path     string   `toml:"path"`
	Handler  string   `toml:"handler"`
	Methods  []string `toml:"methods"`
	Licenses []string `toml:"licenses"`
}

type Routing struct {
	Mode  string `toml:"mode"`
	Routs []Rout `toml:"routs"`
}

func initWidgetVirtualRouts() {
	reloadWidgetVirtualRouts()

	go RefreshEvery(3*time.Second, reloadWidgetVirtualRouts)
}

func RefreshEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func reloadWidgetVirtualRouts() {
	file, err := store.LoadOrNewFile(SettingsBucketName, RoutingSettingsFileName)

	if err != nil {
		logrus.Errorf("vrouting: load error, %v", err)
		return
	}

	var routing = &Routing{}

	if _, err := toml.Decode(string(file.RawData().Bytes()), routing); err != nil {
		logrus.Errorf("vrouting: decode toml error, %v, %q", err, string(file.RawData().Bytes()))
		return
	}

	routerMutex.Lock()
	defer routerMutex.Unlock()

	router = NewRouter()

	for _, _r := range routing.Routs {

		handler := NewHandlerFromString(_r.Handler)
		handler.Licenses = _r.Licenses
		handler.Methods = _r.Methods
		handler.Path = _r.Path

		if len(_r.Methods) == 0 {
			router.Handle(_r.Path, handler).Methods("GET")
		} else {
			router.Handle(_r.Path, handler).Methods(_r.Methods...)
		}
	}
}
