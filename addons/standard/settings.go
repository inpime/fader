package standard

import (
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/utils"
	// gutils "utils"
	"github.com/inpime/fader/utils/sdata"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"basic" json:"basic"`
}

func (s Settings) Merge(cfg interface{}) error {

	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	TplCache       bool             `toml:"tplcache" json:"tplcache"`
	MailerProvider string           `toml:"mailerprovider" json:"mailerprovider"`
	Config         *sdata.StringMap `toml:"config" json:"config"`
}

func MainSettings() *Settings {
	return config.Cfgx.Config(addonName).(*Settings)
}