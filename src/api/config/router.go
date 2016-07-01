package config

import (
	"api/vrouter"
	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"store"
	"sync"
	"time"
)

var Router *vrouter.Router
var routerMutex sync.Mutex

type Rout struct {
	Path      string   `toml:"path"`
	Name      string   `toml:"name"`
	Handler   string   `toml:"handler"`
	Methods   []string `toml:"methods"`
	Licenses  []string `toml:"licenses"`
	IsSpecial bool     `toml:"special"`
}

// Routing settings@main#routs
type Routing struct {
	Routs []Rout `toml:"routs"`
}

func ReloadAppRouts() {
	routerMutex.Lock()
	defer routerMutex.Unlock()

	Router = vrouter.NewRouter()

	for _, fileName := range AppRouts() {
		routeUpdate(fileName)
	}
}

func routeUpdate(fileName string) {

	file, err := store.LoadOrNewFile("settings", fileName)

	if err != nil {
		logrus.WithField("_service", loggerKey).Errorf("load %q error, %v", fileName, err)
		return
	}

	var routing = &Routing{}

	if _, err := toml.Decode(string(file.RawData().Bytes()), routing); err != nil {
		logrus.WithField("_service", loggerKey).Errorf("%q decode toml error, %v, %q", fileName, err, string(file.RawData().Bytes()))
		return
	}

	for _, _r := range routing.Routs {

		handler := NewHandlerFromRoute(_r)

		if len(_r.Methods) == 0 {
			Router.Handle(_r.Path, handler).Methods("GET").Name(_r.Name)
		} else {
			Router.Handle(_r.Path, handler).Methods(_r.Methods...).Name(_r.Name)
		}
	}
}

func NewHandlerFromRoute(r Rout) vrouter.Handler {
	var h vrouter.Handler
	if r.IsSpecial {
		h = vrouter.Handler{
			Bucket:         "",
			File:           "",
			SpecialHandler: r.Handler,
		}
	} else {
		h = vrouter.NewHandlerFromString(r.Handler)
	}

	h.Licenses = r.Licenses
	h.Path = r.Path
	h.Methods = r.Methods

	return h
}

// InitRoute
func InitRoute() {
	ReloadAppRouts()

	// TODO: synchronization with the previous launch
	go refreshEvery(3*time.Second, ReloadAppRouts)
}
