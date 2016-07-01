package config

var loggerKey = "config"

func Init() {
	InitApp()

	InitTpl()

	InitRoute()
}
