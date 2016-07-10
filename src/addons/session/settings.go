package session

import (
	"api/utils"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"session" json:"session"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
}
