package search

import (
	"strings"

	"github.com/flosch/pongo2"
	"github.com/inpime/fader/api/addons"
	"github.com/inpime/fader/store"
	"github.com/labstack/echo"
)

const NAME = "search"

var (
	addonName = NAME
	version   = "0.1.0"
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

func (s Extension) Destroy() {

}

func (s Extension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (Extension) Setup() {
}

func (Extension) TemplateSettings() addons.Configuration {
	return &Settings{&settings{}}
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

	pongo2.DefaultSet.Globals["Search"] = func(bucketName,
		queryString,
		fields,
		page,
		perpage,
		filter,
		sortFields *pongo2.Value) *pongo2.Value {

		_filter := store.NewSearchFilter(strings.ToLower(bucketName.String()))
		_filter.SetQueryString(queryString.String())
		_filter.SetPage(page.Integer())
		_filter.SetPerPage(perpage.Integer())

		queryRaw := BuildAdvancedQuery(queryString.String(),
			fields.String(), // field1 (field2 ...)
			page.Integer(),
			perpage.Integer(),
			filter.String(),     // prefix field_name field_value ([operation_name field_name field_value], ...)
			sortFields.String()) // field asc|desc ([field asc|desc], ...)
		_filter.SetQueryRaw(queryRaw)

		return pongo2.AsValue(MakeSearch(_filter))
	}

	return nil
}
