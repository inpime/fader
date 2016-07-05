package filestatic

import (
	"api/addons"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	addonName = "filestatic"
	version   = "0.1.0"
	// secion name of file settings@main
	FileContentSectionNameKey = "filecontent"

	ByNameRouteName = addonName + ".byname"
	ByIDRouteName   = addonName + ".byid"

	// bucket name
	FileContentBucketNameKey = "bucket"
)

func init() {
	addons.AddAddon(&Extension{})
}

type Extension struct {
}

func (Extension) Version() string {
	return version
}

func (Extension) Name() string {
	return addonName
}

func (s Extension) Destroy() {

}

func (Extension) Setup() {
}

func (s Extension) TemplateSettings() addons.Configuration {
	return Settings{&settings{}}
}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(addonName+".byname", FileContentByNameHandler)
	fnReg(addonName+".byid", FileContentByNameHandler)
}

func (*Extension) InjectTplAddons() error {
	pongo2.RegisterFilter("fc", filterUrlFileByName)
	pongo2.RegisterFilter("filecontenturl", filterUrlFileByName)
	pongo2.RegisterFilter("urlfile", filterUrlFileByName) // OLD
	return nil
}
