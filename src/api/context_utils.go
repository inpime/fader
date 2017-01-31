package api

import (
	"interfaces"
	"time"

	lua "github.com/yuin/gopher-lua"
)

var luaFormFileTypeName = "form_file"

type luaFormFile struct {
	FileName    string
	ContentType string
	Data        []byte
	CreatedAt   time.Time
}

func newLuaFormFile(
	name,
	ct string,
	data []byte,
) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = &luaFormFile{
			FileName:    name,
			Data:        data,
			ContentType: ct,
			CreatedAt:   time.Now(),
		}
		L.SetMetatable(ud, L.GetTypeMetatable(luaFormFileTypeName))
		L.Push(ud)
		return 1
	}
}

func checkFormFile(L *lua.LState) *luaFormFile {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaFormFile); ok {
		return v
	}
	L.ArgError(1, "form file expected")
	return nil
}

var formFileMethods = map[string]lua.LGFunction{
	"Name":        formFileName,
	"Type":        formFileType,
	"ContentType": formFileContentType,
}

func formFileName(L *lua.LState) int {
	o := checkFormFile(L)
	L.Push(lua.LString(o.FileName))
	return 1
}

func formFileContentType(L *lua.LState) int {
	o := checkFormFile(L)
	L.Push(lua.LString(o.ContentType))
	return 1
}

func formFileType(L *lua.LState) int {
	o := checkFormFile(L)
	L.Push(lua.LString(interfaces.GetUserTypeFromContentType(o.ContentType)))
	return 1
}
