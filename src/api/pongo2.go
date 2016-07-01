package api

// import (
// 	"github.com/flosch/pongo2"
// 	_ "github.com/flosch/pongo2-addons"
// )

// var DefaultLoader = MustNewBoltdDBLoader()

// // virtual templates
// // var tpls = pongo2.NewSet("vtpl", tplsLoader)

// func initTemplates() {
// 	pongo2.DefaultSet = pongo2.NewSet("vtpl", DefaultLoader)
// 	pongo2.FromString = pongo2.DefaultSet.FromString
// 	pongo2.FromFile = pongo2.DefaultSet.FromFile
// 	pongo2.FromCache = ExecuteFromCache
// 	pongo2.RenderTemplateString = pongo2.DefaultSet.RenderTemplateString
// 	pongo2.RenderTemplateFile = pongo2.DefaultSet.RenderTemplateFile

// 	pongo2InitGlobalCustoms()
// 	pongo2InitAddons()
// }
