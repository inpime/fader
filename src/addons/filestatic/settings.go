package filestatic

import (
	"api/config"
	"api/utils"
)

func MainSettings() Settings {
	return config.Cfgx.Config(NAME).(Settings)
}

type Settings struct {
	// same value as addon name
	*settings `toml:"filestatic" json:"filestatic"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, *cfg.(Settings).settings)
}

type settings struct {
	BucketSource string `toml:"bucket" json:"bucket"`
}
