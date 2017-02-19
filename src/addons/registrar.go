package addons

import (
	"github.com/flosch/pongo2"
	"github.com/yuin/gopher-lua"
)

// Регистрация расшриений
var Addons = make(map[string]Addon)

type RegisterPongo2Filters func(name string, fn pongo2.FilterFunction) error
type RegisterPongo2Tags func(name string, fn pongo2.TagParser) error

var _ RegisterPongo2Filters = pongo2.RegisterFilter
var _ RegisterPongo2Tags = pongo2.RegisterTag

type Addon interface {
	Version() string // MAJOR.MINOR.PATCH
	// utils github.com/hashicorp/go-version

	Name() string
	Author() string
	Description() string

	LuaLoader(L *lua.LState) int
	ExtContextPongo2(c pongo2.Context) error
	ExtTagsFiltersPongo2(
		addf RegisterPongo2Filters,
		repf RegisterPongo2Filters,
		addt RegisterPongo2Tags,
		rapt RegisterPongo2Tags,
	) error
}
