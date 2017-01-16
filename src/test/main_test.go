package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestContext(t *testing.T) {
	var L = lua.NewState()
	defer L.Close()
	ctx := &testcontext{}
	ctx.Set("a", "b")
	regContext(L)
	L.SetGlobal("ctx", L.NewFunction(newContext(ctx)))

	err := L.DoString(`
ctx():Set("res", ctx():Get("a") == "b")
ctx():Set("a", "c")
`)
	assert.NoError(t, err)
	assert.Equal(t, "c", ctx.Get("a"))
	assert.True(t, ctx.Get("res").(bool))
}

////////////////////////////////////////////////////////////////////////////////
// lua test context
////////////////////////////////////////////////////////////////////////////////

var luaContext = "context"

func regContext(L *lua.LState) {
	mt := L.NewTypeMetatable(luaContext)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), ctxMethods))
}

func newContext(ctx *testcontext) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = ctx
		L.SetMetatable(
			ud,
			L.GetTypeMetatable(luaContext),
		)
		L.Push(ud)
		return 1
	}
}

func checkContext(L *lua.LState) *testcontext {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*testcontext); ok {
		return v
	}
	L.ArgError(1, "expected *testcontext")
	return nil
}

var ctxMethods = map[string]lua.LGFunction{
	"Set": func(L *lua.LState) int {
		ctx := checkContext(L)
		if L.GetTop() != 3 {
			L.ArgError(1, "expected 2 args")
			return 0
		}
		key := L.CheckString(2)
		v := L.Get(3)
		switch _v := v.(type) {
		case lua.LString:
			ctx.Set(key, string(_v))
		case lua.LBool:
			ctx.Set(key, bool(_v))
		case lua.LNumber:
			ctx.Set(key, float64(_v))
		}
		return 0
	},
	"Get": func(L *lua.LState) int {
		ctx := checkContext(L)
		if L.GetTop() != 2 {
			L.ArgError(1, "expected 1 args")
			return 0
		}
		key := L.CheckString(2)
		v := ctx.Get(key)

		switch _v := v.(type) {
		case string:
			L.Push(lua.LString(_v))
		case bool:
			L.Push(lua.LBool(_v))
		case float64:
			L.Push(lua.LNumber(_v))
		default:
			reason := fmt.Sprintf("not supported type %T", v)
			L.ArgError(2, reason)
			return 0
		}

		return 1
	},
}
