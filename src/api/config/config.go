package config

var loggerKey = "config"

var Cfg *Config
var AppVersion string = "dev"

type Config struct {
	AppVersion string

	Address       string
	WorkspacePath string

	Session ApiSessionConfig
	Store   AppStoreConfig
	Search  SearchStore
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
