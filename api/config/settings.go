package config

import "github.com/inpime/sdata"

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

	return sdata.Mergex(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	IncludeFiles []string `toml:"include" json:"include"`
}

func (s Settings) Include() []string {
	return s.IncludeFiles
}
