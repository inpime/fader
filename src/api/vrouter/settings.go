package vrouter

import (
	// "api/addons"
	"api/utils"
	"time"
)

type Settings struct {
	// same value as `addonName`
	*settings `toml:"routing" json:"routing"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	Routs []Rout       `toml:"routs" json:"routs"`
	CSRF  CSRFSettings `toml:"csrf" json:"csrf"`
}

type Rout struct {
	Path      string   `toml:"path" json:"path"`
	Name      string   `toml:"name" json:"name"`
	Handler   string   `toml:"handler" json:"handler"`
	Methods   []string `toml:"methods" json:"methods"`
	Licenses  []string `toml:"licenses" json:"licenses"`
	IsSpecial bool     `toml:"special" json:"special"`
	CSRF      bool     `toml:"csrf" json:"csrf"`
	// CSRFTokenLookup string   `toml:"csrflookup" json:"csrflookup"`
}

type CSRFSettings struct {
	Enabled     bool   `toml:"enabled" json:"enabled"`
	Secret      string `toml:"secret" json:"secret"`
	TokenLookup string `toml:"tokenlookup" json:"tokenlookup"`
	// header:"X-CSRF-Token"
	// form:"csrf"

	Cookie CSRFCookieSettings `toml:"cookie" json:"cookie"`
}

type CSRFCookieSettings struct {
	Name     string        `toml:"name" json:"name"`
	Path     string        `toml:"path" json:"path"`
	MaxAge   time.Duration `toml:"age" json:"age"`
	Domain   string        `toml:"domain" json:"domain"`
	Secure   bool          `toml:"secure" json:"secure"`
	HTTPOnly bool          `toml:"httponly" json:"httponly"`
}
