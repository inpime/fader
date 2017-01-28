package api

import (
	"api/router"
	"interfaces"
	"log"

	"net/http"

	"github.com/labstack/echo"
	lua "github.com/yuin/gopher-lua"
)

var contextMethods = map[string]lua.LGFunction{
	// "URI":        contextGetURI,
	"QueryParam": contextGetQueryParam,
	"NoContext":  contextNoContext,
	"Redirect":   contextRedirect,
	"JSON":       contextRenderJSON,
	"IsGET":      contextMethodIsGET,
	"IsPOST":     contextMethodIsPOST,
	"Set":        contextMethodSet,
	"Get":        contextMethodGet,
	"FormValue":  contextMethodFormValue,
	// alias IsCurrentRoute
	"Route": contextRoute,
	// "Get":        contextMethodGet,
}

func contextGetPath(L *lua.LState) int {
	p := checkContext(L)
	L.Push(lua.LString(p.echoCtx.Path()))
	return 1
}

// func contextGetURI(L *lua.LState) int {
// 	p := checkContext(L)
// 	L.Push(lua.LString(p.echoCtx.Request().URI()))
// 	return 1
// }

func contextRoute(L *lua.LState) int {
	c := checkContext(L)
	route := router.MatchVRouteFromContext(c.echoCtx)

	if route == nil {
		// TODO: informing that an empty route, should not happen

		return 0
	}

	if L.GetTop() >= 2 {
		route = &interfaces.RouteMatch{
			Route: nil,
			Vars:  make(map[string]string),
		}

		foundRoute := vrouter.Get(L.CheckString(2))

		if foundRoute != nil {
			route.Route = foundRoute
			route.Handler = foundRoute.Options()
		}
	}

	// Push route
	newLuaRoute(route)(L)

	return 1
}

// Getter and setter for the Context#Queryparam
func contextGetQueryParam(L *lua.LState) int {
	p := checkContext(L)
	var value string
	if L.GetTop() == 2 {
		value = p.echoCtx.QueryParam(L.CheckString(2))
	}
	L.Push(lua.LString(value))
	return 1
}

func contextNoContext(L *lua.LState) int {
	p := checkContext(L)

	p.Err = p.echoCtx.NoContent(L.CheckInt(2))
	p.Rendered = true

	return 0
}

func contextRedirect(L *lua.LState) int {
	p := checkContext(L)

	p.Err = p.echoCtx.Redirect(http.StatusFound, L.CheckString(2))
	p.Rendered = true

	return 0
}

func contextRenderJSON(L *lua.LState) int {
	p := checkContext(L)
	status := L.CheckInt(2)
	table := L.CheckTable(3)

	jsonMap := make(map[string]interface{}, table.Len())

	table.ForEach(func(key, value lua.LValue) {
		var _key string
		var _value interface{}

		_key = key.String()

		switch value.Type() {
		case lua.LTNumber:
			_value = float64(value.(lua.LNumber))
		case lua.LTNil:
			_value = nil
		case lua.LTBool:
			_value = bool(value.(lua.LBool))
		case lua.LTString:
			_value = string(value.(lua.LString))
		case lua.LTUserData:
			_value = value.(*lua.LUserData).Value
		default:
			log.Printf(
				"[ERR] not expected type value, got %q, for field %q",
				value.Type(),
				_key,
			)
		}

		jsonMap[_key] = _value
	})

	p.Err = p.echoCtx.JSON(status, jsonMap)
	p.Rendered = true

	return 0
}

func contextResponseStatus(L *lua.LState) int {
	p := checkContext(L)
	status := L.CheckInt(2)
	p.ResponseStatus = status

	return 0
}

func contextMethodIsGET(L *lua.LState) int {
	p := checkContext(L)
	L.Push(lua.LBool(p.echoCtx.Request().Method() == echo.GET))
	return 1
}

func contextMethodIsPOST(L *lua.LState) int {
	p := checkContext(L)
	L.Push(lua.LBool(p.echoCtx.Request().Method() == echo.POST))
	return 1
}

func contextMethodSet(L *lua.LState) int {
	p := checkContext(L)
	k := L.CheckString(2)
	lv := L.CheckAny(3)

	v := ToValueFromLValue(lv)
	if v == nil {
		log.Printf("ctx.Set: not supported type, got %T, key %s", lv, k)
		return 0
	}
	p.echoCtx.Set(k, v)

	return 0
}

// contextMethodGet
// Supported types: int, float, string, bool, nil
func contextMethodGet(L *lua.LState) int {
	p := checkContext(L)
	k := L.CheckString(2)
	v := p.echoCtx.Get(k)

	lv := ToLValueOrNil(v, L)
	if lv == nil {
		log.Printf("ctx.Get: not supported type, got %T, key %s", v, k)
		return 0
	}

	L.Push(lv)
	return 1
}

func contextMethodFormValue(L *lua.LState) int {
	c := checkContext(L)

	L.Push(lua.LString(c.echoCtx.FormValue(L.CheckString(2))))

	return 1
}
