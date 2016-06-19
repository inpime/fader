package api

import (
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
)

var tplsLoader = MustNewBoltdDBLoader()

// virtual templates
var tpls = pongo2.NewSet("vtpl", tplsLoader)

func initTemplates() {
	pongo2InitGlobalCustoms()
	pongo2InitAddons()

	pongo2.DefaultSet = tpls
	pongo2.FromString = tpls.FromString
	pongo2.FromFile = tpls.FromFile
	pongo2.FromCache = ExecuteFromCache
	pongo2.RenderTemplateString = tpls.RenderTemplateString
	pongo2.RenderTemplateFile = tpls.RenderTemplateFile
}
