package search

import (
	"github.com/inpime/fader/utils"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"search" json:"search"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
}
