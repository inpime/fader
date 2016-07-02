package search

import (
	"api/config"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"store"
	"strings"
)

var (
	addonName = "fader.addons.search"
	version   = "0.1.0"
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

func (s Extension) Destroy() {

}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*Extension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (*Extension) InjectTplAddons() error {
	pongo2.DefaultSet.Globals["SearchFiles"] = func(
		bucketName,
		queryStr,
		page,
		perpage *pongo2.Value,
	) *pongo2.Value {

		filter := store.NewSearchFilter(strings.ToLower(bucketName.String()))
		filter.SetQueryString(queryStr.String())
		filter.SetPage(page.Integer())
		filter.SetPerPage(perpage.Integer())

		queryRaw := buildSearchQueryFilesByBycket(
			strings.ToLower(bucketName.String()),
			queryStr.String(),
			page.Integer(),
			perpage.Integer(),
		)
		filter.SetQueryRaw(queryRaw)

		return pongo2.AsValue(makeSearch(filter))
	}

	return nil
}
