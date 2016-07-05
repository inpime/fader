package filestatic

import (
	"api/config"
	"api/utils"
)

func MainSettings() Settings {
	return config.Cfgx.Config(addonName).(Settings)
}

type Settings struct {
	// same value as `addonName`
	*settings `toml:"filestatic" json:"filestatic"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, *cfg.(Settings).settings)
}

type settings struct {
	BucketSource string `toml:"bucket" json:"bucket"`
}
