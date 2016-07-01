package filestatic

import (
	"api/config"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	addonName = "fader.addons.filestatic"
	// secion name of file settings@main
	FileContentSectionNameKey = "filecontent"
	// bucket name
	FileContentBucketNameKey = "bucket"

	ByNameHandlerName = "FileContentByName"
	ByIDHandlerName   = "FileContentByID"
)

func init() {
	config.AddExtension(&FileStatisExtension{})
}

type FileStatisExtension struct {
}

func (FileStatisExtension) Name() string {
	return addonName
}

func (s FileStatisExtension) Destroy() {

}

func (s FileStatisExtension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*FileStatisExtension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(addonName+".filtestatic_byname", FileContentByNameHandler)
	fnReg(addonName+".filtestatic_byid", FileContentByNameHandler)
}

func (*FileStatisExtension) InjectTplAddons() error {
	pongo2.RegisterFilter("fc", filterUrlFileByName)
	pongo2.RegisterFilter("filecontenturl", filterUrlFileByName)
	pongo2.RegisterFilter("urlfile", filterUrlFileByName) // OLD
	return nil
}
