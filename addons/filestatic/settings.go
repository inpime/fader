package filestatic

import (
	"github.com/inpime/fader/api/config"
	"github.com/inpime/sdata"
)

func MainSettings() *Settings {
	return config.Cfgx.Config(NAME).(*Settings)
}

type Settings struct {
	// same value as addon name
	*settings `toml:"filestatic" json:"filestatic"`
}

func (s Settings) Merge(cfg interface{}) error {

	return sdata.Mergex(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	// BucketSource the default value
	BucketSource string `toml:"bucket" json:"bucket"`
}
