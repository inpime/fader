package shop

import (
	"github.com/inpime/fader/api/config"
	"github.com/inpime/sdata"
)

func MainSettings() *Settings {
	return config.Cfgx.Config(NAME).(*Settings)
}

type Settings struct {
	// same value as addon name
	*settings `toml:"shop" json:"shop"`
}

func (s Settings) Merge(cfg interface{}) error {
	return sdata.Mergex(s.settings, cfg.(*Settings).settings)
}

type settings struct{}
