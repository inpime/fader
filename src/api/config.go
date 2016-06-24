package api

var Cfg *Config
var AppVersion string = "dev"

var (
	BucketsBucketName  = "buckets"
	UsersBucketName    = "users"
	SettingsBucketName = "settings"
	StaticBucketName   = "static"
	PagesBucketName    = "pages"

	ConsoleBucketName = "console"

	GuestUserFileName       = "guestuser"
	RoutingSettingsFileName = "routing"

	// system store names
	// BucketsStoreName  = "boltdb.buckets"
	UsersStoreName = "boltdb.users"
	// SettingsStoreName = "boltdb.settings"

	// System file names
	MainSettingsFileName = "main"
)

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

	SecretKey  string
	BucketName string

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
