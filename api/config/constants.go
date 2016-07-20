package config

const (
	SpecialHandlerArgsKey = "_special_handler_args"
)

var (

	// ---------------------
	// Bucket names
	// ---------------------

	BucketsBucketName  = "buckets"
	UsersBucketName    = "users"
	SettingsBucketName = "settings"
	StaticBucketName   = "static"
	PagesBucketName    = "pages"
	// ConsoleBucketName extension addons/console
	ConsoleBucketName = "console"

	// ---------------------
	// Store names
	// ---------------------

	UsersStoreName = "boltdb.users"

	// ---------------------
	// System file names
	// ---------------------

	GuestUserFileName       = "guestuser"
	RoutingSettingsFileName = "routing"
	MainSettingsFileName    = "main"

	// ---------------------
	// Settings secions
	// ---------------------
)
