package api

import (
	"api/router"
	"interfaces"
	"io"
	"log"
	"time"

	"net/http"

	"bytes"

	"github.com/labstack/echo"
	lua "github.com/yuin/gopher-lua"
)

var contextMethods = map[string]lua.LGFunction{
	// "URI":        contextGetURI,
	"QueryParam": contextGetQueryParam,
	"NoContent":  contextNoContent,
	"Redirect":   contextRedirect,
	"JSON":       contextRenderJSON,
	// FileContent рендер переданного файла (RawData)
	"FileContent": contextMethodFileContnet,
	"IsGET":       contextMethodIsGET,
	"IsPOST":      contextMethodIsPOST,
	"Set":         contextMethodSet,
	"Get":         contextMethodGet,
	"FormValue":   contextMethodFormValue,
	"FormFile":    contextMethodFormFile,
	// alias IsCurrentRoute
	"Route": contextRoute,
	// "Get":        contextMethodGet,

	"AppExport": contextAppExport,
	"AppImport": contextAppImport,

	"CurrentFile":    contextGetCurrentFile,
	"MiddlewareFile": contextGetMiddlewareFile,
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

func contextNoContent(L *lua.LState) int {
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

// contextMethodFileContnet рендер файла RawData
func contextMethodFileContnet(L *lua.LState) int {
	c := checkContext(L)
	status := L.CheckInt(2)

	ud := L.CheckUserData(3)
	file, ok := ud.Value.(*luaFile)
	if !ok {
		L.ArgError(2, "file expected")
		return 0
	}

	c.Err = c.echoCtx.Blob(
		status,
		file.ContentType,
		file.RawData,
	)
	c.Rendered = true

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

func contextMethodFormFile(L *lua.LState) int {
	c := checkContext(L)

	f, err := c.echoCtx.FormFile(L.CheckString(2))

	if err != nil {
		log.Println("FormFile: ", err)
		L.Push(lua.LBool(false))
		return 1
	}

	of, err := f.Open()

	if err != nil {
		log.Println("FormFile: open file,", err)
		L.Push(lua.LBool(false))
		return 1
	}
	defer of.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, of)

	if err != nil {
		log.Println("FormFile: copy file,", err)
		L.Push(lua.LBool(false))
		return 1
	}

	return newLuaFormFile(
		f.Filename,
		f.Header.Get("Content-Type"),
		buf.Bytes(),
	)(L)
}

////////////////////////////////////////////////////////////////////////////////
// import export
////////////////////////////////////////////////////////////////////////////////

func contextAppExport(L *lua.LState) int {
	c := checkContext(L)

	importer := interfaces.NewImportManager(
		bucketManager,
		fileManager,
	)

	data, _ := importer.Export(
		"vDEV-"+time.Now().String(),
		"Fader",
		time.Now().String(),
	)
	fileName := "Fader.vDEV-" + time.Now().String() + ".txt"

	buf := bytes.NewReader(data)

	c.Err = c.echoCtx.Attachment(buf, fileName)
	c.Rendered = true

	return 0
}

func contextAppImport(L *lua.LState) int {
	c := checkContext(L)

	importer := interfaces.NewImportManager(
		bucketManager,
		fileManager,
	)

	f, err := c.echoCtx.FormFile("file")

	if err != nil {
		log.Println("AppImport: ", err)
		L.Push(lua.LBool(false))
		return 1
	}

	of, err := f.Open()

	if err != nil {
		log.Println("AppImport: open file,", err)
		L.Push(lua.LBool(false))
		return 1
	}
	defer of.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, of)
	if err != nil {
		log.Println("AppImport: io copy,", err)
		L.Push(lua.LBool(false))
		return 1
	}

	info, err := importer.Import(buf.Bytes())
	if err != nil {
		log.Println("AppImport: import,", err)
		L.Push(lua.LBool(false))
		return 1
	}

	log.Println("AppImport: success, ", info)

	L.Push(lua.LBool(true))
	return 1
}

func contextGetCurrentFile(L *lua.LState) int {
	c := checkContext(L)

	return newLuaFile(c.CurrentFile)(L)
}
func contextGetMiddlewareFile(L *lua.LState) int {
	c := checkContext(L)

	return newLuaFile(c.MiddlewareFile)(L)
}
