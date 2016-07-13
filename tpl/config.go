package tpl

import (
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
)

var TplDefaultLoader = MustNewBoltdDBLoader()

// InitTpl
func InitTpl() {
	pongo2.DefaultSet = pongo2.NewSet("virtual tpls", TplDefaultLoader)
	pongo2.FromString = pongo2.DefaultSet.FromString
	pongo2.FromFile = pongo2.DefaultSet.FromFile
	pongo2.FromCache = ExecuteFromMemCache
	pongo2.RenderTemplateString = pongo2.DefaultSet.RenderTemplateString
	pongo2.RenderTemplateFile = pongo2.DefaultSet.RenderTemplateFile

	// pongo2InitGlobalCustoms()
	// pongo2InitAddons()
}
