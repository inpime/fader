package standard

import (
	"api/config"
	"api/utils"
	"fmt"
	"github.com/Sirupsen/logrus"
	gutils "utils"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"basic" json:"basic"`
}

func (s Settings) Merge(cfg interface{}) error {
	if s.settings == nil {
		logrus.WithField("_service", addonName).Error("null settings")
		return fmt.Errorf("invalid config")
	}

	if _, ok := cfg.(Settings); !ok {

		logrus.WithField("_service", addonName).Errorf("not supported settings %T", cfg)

		return fmt.Errorf("invalid config")
	}

	if cfg.(Settings).settings == nil {

		logrus.WithField("_service", addonName).Error("null merged settings")
		return fmt.Errorf("invalid config")
	}

	return utils.AppendOrReplace(s.settings, *cfg.(Settings).settings)
}

type settings struct {
	TplCache       bool     `toml:"tplcache" json:"tplcache"`
	MailerProvider string   `toml:"mailerprovider" json:"mailerprovider"`
	Config         gutils.M `toml:"config" json:"config"`
}

func MainSettings() Settings {
	return config.Cfgx.Config(addonName).(Settings)
}
