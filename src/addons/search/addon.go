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
)

func init() {
	config.AddExtension(&SearchExtension{})
}

type SearchExtension struct {
}

func (SearchExtension) Name() string {
	return addonName
}

func (s SearchExtension) Destroy() {

}

func (s SearchExtension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*SearchExtension) RegEchoHandlers(fnReg func(string, func(ctx echo.Context) error)) {

}

func (*SearchExtension) InjectTplAddons() error {
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
