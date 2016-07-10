package config

import (
	"api/utils"
)

var sectionName = "main"

// MainSettings
func MainSettings() *Settings {

	return Cfgx.Config(sectionName).(*Settings)
}

// Wraper setting
type Settings struct {
	*settings `toml:"main" json:"main"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	IncludeFiles []string `toml:"include" json:"include"`
}

func (s Settings) Include() []string {
	return s.IncludeFiles
}
