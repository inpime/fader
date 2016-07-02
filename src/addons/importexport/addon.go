package importexport

import (
	"api/config"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	version        = "0.1.0"
	addonName      = "fader.addons.importexport"
	filenamePrefix = "FADER"

	SettingsSectionNameKey = "importexport"

	ArchiveFaderDataFileName = ".faderdata"

	ImportKey = addonName + ".import"
	ExportKey = addonName + ".export"

	ArchiveURLLatestVersion = "https://s3.eu-central-1.amazonaws.com/releases.fader.inpime.com/archives/FADER(sys).dev.latest.zip"
)

func init() {
	config.AddExtension(&Extension{})
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

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
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
