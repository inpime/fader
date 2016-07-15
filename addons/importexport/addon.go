package importexport

import (
	"github.com/flosch/pongo2"
	"github.com/inpime/fader/api/addons"
	"github.com/labstack/echo"
)

const (
	NAME    = "importexport"
	VERSION = "0.1.0"
)

var (
	version        = VERSION
	addonName      = NAME
	filenamePrefix = "FADER"

	ArchiveFaderDataFileName = ".faderdata"

	ImportKey = addonName + ".import"
	ExportKey = addonName + ".export"

	// You can override
	ArchiveURLLatestVersion = "https://s3.eu-central-1.amazonaws.com/releases.fader.inpime.com/archives/FADER(console).dev.latest.zip"
)

func init() {
	addons.AddAddon(&Extension{})
}

type Extension struct {
}

func (Extension) Name() string {
	return addonName
}

func (Extension) Version() string {
	return version
}

func (Extension) Destroy() {

}

func (Extension) Setup() {
}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s Extension) TemplateSettings() addons.Configuration {
	return &Settings{&settings{
		Groups: []GroupSettings{},
	}}
}

func (e Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {
	fnReg(ImportKey, ImportHandler)
	fnReg(ExportKey, ExportHandler)
}

func (e Extension) InjectTplAddons() error {
	// ListGroupsImportExport возвращает список групп указанных в настройках приложения
	pongo2.DefaultSet.Globals["ListGroupsImportExport"] = func() *pongo2.Value {
		return pongo2.AsValue(ListGroupsImportExport())
	}

	return nil
}
