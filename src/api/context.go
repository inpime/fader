package api

import (
	"interfaces"

	"github.com/labstack/echo"
	"github.com/yuin/gopher-lua"
)

const luaContextTypeName = "context"

func ContextLuaExecutor(
	L *lua.LState,
	ctx echo.Context,
) *Context {
	mt := L.NewTypeMetatable(luaContextTypeName)
	_ctx := &Context{
		echoCtx: ctx,
	}

	L.SetGlobal("ctx", L.NewFunction(newContext(_ctx)))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), contextMethods))

	return _ctx
}

type Context struct {
	echoCtx echo.Context

	Err      error // in lua script was error
	Rendered bool  // in lua script was executed render

	CurrentFile    *interfaces.File
	MiddlewareFile *interfaces.File

	ResponseStatus int // in lua script set response status
}

func (c Context) EchoCtx() echo.Context {
	return c.echoCtx
}

// internal functions

func checkContext(L *lua.LState) *Context {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Context); ok {
		return v
	}
	L.ArgError(1, "context expected")
	return nil
}

func newContext(ctx *Context) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = ctx
		L.SetMetatable(ud, L.GetTypeMetatable(luaContextTypeName))
		L.Push(ud)
		return 1
	}
}
