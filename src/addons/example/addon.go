package example

import (
	"addons"

	"github.com/flosch/pongo2"
	"github.com/yuin/gopher-lua"
)

const (
	VERSION     = "0.1"
	NAME        = "example"
	AUTHOR      = "Fader"
	DESCRIPTION = `Example of an addon for learning`
)

var (
	_ addons.Addon = (*Addon)(nil)
)

func NewAddon() *Addon {
	return &Addon{}
}

type Addon struct {
}

func (a Addon) Version() string {
	return VERSION
}

func (a Addon) Name() string {
	return NAME
}

func (a Addon) Author() string {
	return AUTHOR
}

func (a Addon) Description() string {
	return DESCRIPTION
}

func (a *Addon) LuaLoader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

func (a *Addon) ExtContextPongo2(_ctx pongo2.Context) error {
	ctx := make(pongo2.Context)
	ctx["ContextFunction"] = func() *pongo2.Value {
		return pongo2.AsValue("check extension 'example'")
	}
	_ctx.Update(ctx)
	return nil
}

func (a *Addon) ExtTagsFiltersPongo2(
	addf addons.RegisterPongo2Filters,
	repf addons.RegisterPongo2Filters,
	addt addons.RegisterPongo2Tags,
	rapt addons.RegisterPongo2Tags,
) error {
	return nil
}
