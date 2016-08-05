package search

import "github.com/inpime/sdata"

type Settings struct {
	// same value as `addonName`
	*settings `toml:"search" json:"search"`
}

func (s Settings) Merge(cfg interface{}) error {
	return sdata.Mergex(s.settings, cfg.(*Settings).settings)
}

type settings struct {
}
