package config

func Init() {
	Cfgx = NewConfig()

	InitTpl()
}

func NewConfig() *configs {
	c := make(configs)
	c.AddConfig(sectionName, &Settings{&settings{}})

	return &c
}
