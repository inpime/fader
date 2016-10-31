package templates

import (
	"io"

	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
)

var DefaultTemplatesLoader TemplatesLoader

func SetupSettings() {
	pongo2.DefaultSet = pongo2.NewSet("virtual tpls", DefaultTemplatesLoader)
	pongo2.FromString = pongo2.DefaultSet.FromString
	pongo2.FromFile = pongo2.DefaultSet.FromFile
	pongo2.FromCache = ExecuteFromMemCache
	pongo2.RenderTemplateString = pongo2.DefaultSet.RenderTemplateString
	pongo2.RenderTemplateFile = pongo2.DefaultSet.RenderTemplateFile
}

type TemplatesLoader interface {
	Abs(base, name string) string
	Get(path string) (io.Reader, error)
}
