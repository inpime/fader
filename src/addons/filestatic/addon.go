package filestatic

import (
	"api/config"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	addonName = "fader.addons.filestatic"
	version   = "0.1.0"
	// secion name of file settings@main
	FileContentSectionNameKey = "filecontent"
	// bucket name
	FileContentBucketNameKey = "bucket"

	ByNameHandlerName = "FileContentByName"
	ByIDHandlerName   = "FileContentByID"
)

func init() {
	config.AddExtension(&Extension{})
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

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(addonName+".filtestatic_byname", FileContentByNameHandler)
	fnReg(addonName+".filtestatic_byid", FileContentByNameHandler)
}

func (*Extension) InjectTplAddons() error {
	pongo2.RegisterFilter("fc", filterUrlFileByName)
	pongo2.RegisterFilter("filecontenturl", filterUrlFileByName)
	pongo2.RegisterFilter("urlfile", filterUrlFileByName) // OLD
	return nil
}
