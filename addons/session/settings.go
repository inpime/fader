package session

import "github.com/inpime/sdata"

type Settings struct {
	// same value as `addonName`
	*settings `toml:"session" json:"session"`
}

func (s Settings) Merge(cfg interface{}) error {
	return sdata.Mergex(s.settings, cfg.(*Settings).settings)
}

type settings struct {
}
