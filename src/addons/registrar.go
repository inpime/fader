package addons

import (
	"github.com/flosch/pongo2"
	"github.com/yuin/gopher-lua"
)

// Регистрация расшриений
var Addons = make(map[string]Addon)

type Addon interface {
	Version() string // MAJOR.MINOR.PATCH
	// utils github.com/hashicorp/go-version

	Name() string
	Author() string
	Description() string

	LuaLoader(L *lua.LState) int
	ExtContextPongo2(c pongo2.Context) error
	ExtTagsFiltersPongo2(
		addf func(name string, fn pongo2.FilterFunction),
		repf func(name string, fn pongo2.FilterFunction),
		addt func(name string, fn pongo2.TagParser),
		rapt func(name string, fn pongo2.TagParser),
	) error
}
