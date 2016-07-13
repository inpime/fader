package config

var loggerKey = "app_settings"

var Cfg *Config

type Config struct {
	AppVersion   string
	AppBuildHash string
	AppBuildDate string

	Address       string
	WorkspacePath string

	Session ApiSessionConfig
	Store   AppStoreConfig
	Search  SearchStore
}

func (c Config) Version() string {
	return c.AppVersion + "." + c.BuildHash() + "." + c.AppBuildDate
}

func (c Config) BuildHash() string {
	if len(c.AppBuildHash) > 7 {
		return c.AppBuildHash[:7]
	}

	return "_"
}

type ApiSessionConfig struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool

	SecretKey   string
	BucketName  string
	SessionName string

	Store StoreConfig
}

type StoreConfig struct {
	Provider       string // boltdb
	BoltDBFilePath string
}

type AppStoreConfig struct {
	StoreConfig

	StaticPath string
}

type SearchStore struct {
	Host      string
	IndexName string
}
