package api

import (
	"interfaces"
	"sync"
	"time"
)

var (
	DefaultDatabasePath = "./_app.db"
	DefaultLogLevel     = DebugLevel
	DefaultBoydLimit    = "2MB"

	DefaultApiHost = ""
	DefaultApiPort = "1323"

	// DefaultTimeoutFileProvider timeout handler of load file
	DefaultTimeoutFileProvider = time.Millisecond * 100
)

// config

type Settings struct {
	ApiPort      string
	ApiHost      string
	DatabasePath string
	BodyLimit    string
	LogLevel     Level

	InitFile string // path to the file

	TimeoutFileProvider time.Duration
}

func SettingsOrDefault(config *Settings) *Settings {
	if len(config.ApiPort) == 0 {
		config.ApiPort = DefaultApiPort
	}

	if len(config.ApiHost) == 0 {
		config.ApiHost = DefaultApiHost
	}

	if len(config.DatabasePath) == 0 {
		config.DatabasePath = DefaultDatabasePath
	}

	if len(config.DatabasePath) == 0 {
		config.DatabasePath = DefaultDatabasePath
	}

	if len(config.BodyLimit) == 0 {
		config.BodyLimit = DefaultBoydLimit
	}

	if config.LogLevel == 0 {
		config.LogLevel = DefaultLogLevel
	}

	if config.TimeoutFileProvider == 0 {
		config.TimeoutFileProvider = DefaultTimeoutFileProvider
	}

	return config
}

// log Level

type Level uint

const (
	_                = iota
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnintLevel
	ErrorLevel
)

func LogLevelFrom(name string) Level {
	switch name {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnintLevel
	case "err":
		return ErrorLevel
	}

	return DefaultLogLevel
}

// config

func newConfig() *Config {
	return &Config{
		Addons: make(map[string]interface{}),
	}
}

type Config struct {
	sync.RWMutex

	Main    ConfigMain             `toml:"main"`
	Routing ConfigRouting          `toml:"routing"`
	Addons  map[string]interface{} `toml:"addons"`
}

func (c *Config) Config() *Config {
	c.RLock()
	defer c.RUnlock()

	return c
}

func (c *Config) Merge(src *Config) error {
	c.Lock()
	defer c.Unlock()

	if src.Main.Use {
		if err := Merge(&c.Main, src.Main); err != nil {
			return err
		}
	}

	if src.Routing.CSRF.Use {
		if err := Merge(&c.Routing.CSRF, src.Routing.CSRF); err != nil {
			return err
		}
	}

	if src.Routing.Use {
		if err := Merge(&c.Routing, src.Routing); err != nil {
			return err
		}
	}

	if len(src.Routing.Routs) > 0 {
		if err := Merge(&c.Routing.Routs, src.Routing.Routs); err != nil {
			return err
		}
	}

	return Merge(&c.Addons, src.Addons)
}

type ConfigMain struct {
	UseMerge

	TplCache bool     `toml:"tplcache"`
	Include  []string `toml:"include"`
}

type ConfigRouting struct {
	UseMerge

	CSRF  ConfigRoutingCSRF           `toml:"csrf"`
	Routs []interfaces.RequestHandler `toml:"routs"`
}

type ConfigRoutingCSRF struct {
	UseMerge

	Enabled      bool                `toml:"enabled"`
	Secret       string              `toml:"secret"`      // "secret"
	TokentLookup string              `toml:"tokenlookup"` // "form:csrf"
	Cookie       ConfigRoutingCookie `toml:"cookie"`
}

type ConfigRoutingCookie struct {
	Name string        `toml:"name"` // "csrf"
	Path string        `toml:"path"` // "/"
	Age  time.Duration `toml:"age"`  // 86400 # 24H
}

type UseMerge struct {
	Use bool `toml:"use"`
}
