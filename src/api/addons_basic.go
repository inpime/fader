package api

import (
	"addons"

	"github.com/flosch/pongo2"
	uuid "github.com/satori/go.uuid"
	"github.com/yuin/gopher-lua"
)

const (
	ADDONS_BASIC_VERSION     = "0.1"
	ADDONS_BASIC_NAME        = "basic"
	ADDONS_BASIC_AUTHOR      = "Fader"
	ADDONS_BASIC_DESCRIPTION = `Example of an addon for learning`
)

var (
	_ addons.Addon = (*AddonBasic)(nil)
)

func NewBasicAddon() *AddonBasic {
	return &AddonBasic{}
}

type AddonBasic struct {
}

func (a AddonBasic) Version() string {
	return ADDONS_BASIC_VERSION
}

func (a AddonBasic) Name() string {
	return ADDONS_BASIC_NAME
}

func (a AddonBasic) Author() string {
	return ADDONS_BASIC_AUTHOR
}

func (a AddonBasic) Description() string {
	return ADDONS_BASIC_DESCRIPTION
}

func (a *AddonBasic) LuaLoader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)

	// register other stuff
	L.SetField(mod, "version", lua.LString(ADDONS_BASIC_VERSION))
	L.SetField(mod, "name", lua.LString(ADDONS_BASIC_NAME))
	L.SetField(mod, "author", lua.LString(ADDONS_BASIC_AUTHOR))
	L.SetField(mod, "description", lua.LString(ADDONS_BASIC_DESCRIPTION))

	// returns the module
	L.Push(mod)

	////////////////////////////////////
	// Custom Types
	////////////////////////////////////

	// Route
	mt := L.NewTypeMetatable(luaRouteTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), routeMethods))
	return 1
}

func (a *AddonBasic) ExtContextPongo2(_ctx pongo2.Context) error {
	ctx := make(pongo2.Context)
	// ctx["ContextFunction"] = func() *pongo2.Value {
	// 	return pongo2.AsValue("context function")
	// }

	ctx["ListBuckets"] = func(bucketID *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(listOfBuckets())
	}

	ctx["ListFilesByBucketID"] = func(bucketID *pongo2.Value) *pongo2.Value {
		var bid uuid.UUID
		switch v := bucketID.Interface().(type) {
		case uuid.UUID:
			bid = v
		case string:
			bid = uuid.FromStringOrNil(v)
		}
		return pongo2.AsValue(filesByBucketID(bid))
	}

	ctx["Route"] = func(name *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(
			&RoutePongo2{
				vrouter.Get(name.String()),
			},
		)
	}

	_ctx.Update(ctx)
	return nil
}

func (a *AddonBasic) ExtTagsFiltersPongo2(
	addf func(name string, fn pongo2.FilterFunction),
	repf func(name string, fn pongo2.FilterFunction),
	addt func(name string, fn pongo2.TagParser),
	rapt func(name string, fn pongo2.TagParser),
) error {
	return nil
}
