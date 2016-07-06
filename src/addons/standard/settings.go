package standard

import (
	"api/config"
	"api/utils"
	gutils "utils"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"basic" json:"basic"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, *cfg.(Settings).settings)
}

type settings struct {
	TplCache bool     `toml:"tplcache" json:"tplcache"`
	Config   gutils.M `toml:"config" json:"config"`
}

func MainSettings() Settings {
	return config.Cfgx.Config(addonName).(Settings)
}
