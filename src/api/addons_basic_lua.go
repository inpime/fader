package api

import (
	"interfaces"
	"log"

	uuid "github.com/satori/go.uuid"
	lua "github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"ListBuckets":         basicFn_ListBuckets,
	"ListFilesByBucketID": basicFn_listFilesFromBucketID,

	// file manager
	"FindFileByName": func(L *lua.LState) int { return 0 },
	"FindFile":       func(L *lua.LState) int { return 0 },
	"CreateFile":     func(L *lua.LState) int { return 0 },
	"CreateFileFrom": func(L *lua.LState) int { return 0 },
	"UpdateFileFrom": func(L *lua.LState) int { return 0 },
	"DeleteFile":     func(L *lua.LState) int { return 0 },

	// bucket manager
	"FindBucketByName": func(L *lua.LState) int { return 0 },
	"FindBucket":       func(L *lua.LState) int { return 0 },
	"CreateBucket":     func(L *lua.LState) int { return 0 },
	"CreateBucketFrom": func(L *lua.LState) int { return 0 },
	"UpdateBucketFrom": func(L *lua.LState) int { return 0 },
}

////////////////////////////////////////////////////////////////////////////////
// Lua interfaces.DataUsed
////////////////////////////////////////////////////////////////////////////////

var luaDataUsed = "DataUsed"

func newDataUsed(v interfaces.DataUsed) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = v
		L.SetMetatable(ud, L.GetTypeMetatable(luaDataUsed))
		L.Push(ud)
		return 1
	}
}

func checkDataUsed(L *lua.LState) interfaces.DataUsed {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(interfaces.DataUsed); ok {
		return v
	}
	L.ArgError(1, "interfaces.DataUsed expected")
	return interfaces.DataUsed(0)
}

// types

var dataUsedMethods = map[string]lua.LGFunction{
	"Has": func(L *lua.LState) int {
		// TODO:
	},
	"Add": func(L *lua.LState) int {
		var self *lua.LUserData
		var v interfaces.DataUsed

		for i := 1; i <= L.GetTop(); i++ {
			ud := L.CheckUserData(i)
			if ud == nil {
				L.ArgError(i, "interfaces.DataUsed expected, got nil")
				continue
			}
			if i == 1 {
				self = ud
			}
			_v, ok := ud.Value.(interfaces.DataUsed)
			if !ok {
				L.ArgError(i, "interfaces.DataUsed expected")
				continue
			}
			v = v | _v
		}

		self.Value = v
		L.Push(self)

		return 1
	},
}

////////////////////////////////////////////////////////////////////////////////
// luaRoute
////////////////////////////////////////////////////////////////////////////////

var luaRouteTypeName = "route"

func newLuaRoute(route *interfaces.RouteMatch) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		ud := L.NewUserData()
		ud.Value = &luaRoute{
			Name:        route.Handler.Name,
			Path:        route.Handler.Path,
			Bucket:      route.Handler.Bucket,
			File:        route.Handler.File,
			HandlerArgs: route.Handler.HandlerArgs,
			route:       route.Route,
		}
		L.SetMetatable(ud, L.GetTypeMetatable(luaRouteTypeName))
		L.Push(ud)
		return 1
	}
}

type luaRoute struct {
	Name   string
	Path   string
	Bucket string
	File   string

	AllowedAlicenses []string
	AllowedMethods   []string
	HandlerArgs      string

	route interfaces.Route
}

func checkRoute(L *lua.LState) *luaRoute {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaRoute); ok {
		return v
	}
	L.ArgError(1, "route expected")
	return nil
}

// luaRoute methods

var routeMethods = map[string]lua.LGFunction{
	"Name":   rotueGetName,
	"Path":   rotueGetPath,
	"Bucket": rotueGetBucket,
	"File":   rotueGetFile,
	"Args":   rotueGetHandlerArgs,
	"Has":    routeHasRoute,

	// generate URL of the current routes in the parameters
	// TODO: renate to URLPath
	"URL":     routeGetURLFromParams,
	"URLPath": routeGetURLFromParams,
}

func routeHasRoute(L *lua.LState) int {
	r := checkRoute(L)

	if L.GetTop() != 2 {
		// if exists current route then return true
		L.Push(lua.LBool(r.route != nil))
		return 1
	}

	routeName := L.CheckString(2)
	L.Push(lua.LBool(routeName == r.Name))
	return 1
}

func rotueGetName(L *lua.LState) int {
	r := checkRoute(L)
	L.Push(lua.LString(r.Name))
	return 1
}

func rotueGetPath(L *lua.LState) int {
	r := checkRoute(L)
	L.Push(lua.LString(r.Path))
	return 1
}

func rotueGetBucket(L *lua.LState) int {
	r := checkRoute(L)
	L.Push(lua.LString(r.Bucket))
	return 1
}

func rotueGetFile(L *lua.LState) int {
	r := checkRoute(L)
	L.Push(lua.LString(r.File))
	return 1
}

func rotueGetHandlerArgs(L *lua.LState) int {
	r := checkRoute(L)
	L.Push(lua.LString(r.HandlerArgs))
	return 1
}

func routeGetURLFromParams(L *lua.LState) int {
	r := checkRoute(L)

	if r.route == nil {
		// TODO: error
		log.Println("empty router")
		return 0
	}

	var args []string

	if L.GetTop() > 1 {
		args = make([]string, L.GetTop()-1)
		for i := 2; i <= L.GetTop(); i++ {
			args[i-2] = L.CheckString(i)
		}
	}

	url, err := r.route.URLPath(args...)
	// TODO: URL as custom object
	if err != nil {
		// TODO: error
		log.Println("build url", err)
		return 0
	}
	L.Push(lua.LString(url.String()))
	return 1
}

////////////////////////////////////////////////////////////////////////////////
// Bucket and file utils
////////////////////////////////////////////////////////////////////////////////

func basicFn_ListBuckets(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = listOfBuckets()
	L.Push(ud)
	return 1
}

////////////////////////////////////////////////////////////////////////////////
// File type
////////////////////////////////////////////////////////////////////////////////

var luaFileTypeName = "file"

func FileAsLuaFile(L *lua.LState, file *interfaces.File) *lua.LUserData {
	ud := L.NewUserData()
	if file == nil {
		_f := interfaces.NewFile()
		_f.FileID = uuid.NewV4()
		ud.Value = &luaFile{_f}
	} else {
		ud.Value = &luaFile{file}
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaFileTypeName))
	return ud
}

func newLuaFile(file *interfaces.File) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		L.Push(FileAsLuaFile(L, file))
		return 1
	}

}

type luaFile struct {
	*interfaces.File
}

func checkFile(L *lua.LState) *luaFile {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaFile); ok {
		return v
	}
	L.ArgError(1, "route expected")
	return nil
}

// luaRoute methods

var fileMethods = map[string]lua.LGFunction{
	"SetFileName":    func(*lua.LState) int { return 0 },
	"SetBucketID":    func(*lua.LState) int { return 0 },
	"SetLuaScript":   func(*lua.LState) int { return 0 },
	"MetaData":       func(*lua.LState) int { return 0 },
	"StructuralData": func(*lua.LState) int { return 0 },
	"SetRawData":     func(*lua.LState) int { return 0 },
	"SetContentType": func(*lua.LState) int { return 0 },
	"SetOwners":      func(*lua.LState) int { return 0 },
	"SetPrivate":     func(*lua.LState) int { return 0 },
	"SetPublic":      func(*lua.LState) int { return 0 },
	"SetReadOnly":    func(*lua.LState) int { return 0 },
}

////////////////////////////////////////////////////////////////////////////////
// File manager
////////////////////////////////////////////////////////////////////////////////

// список файлов в бакете
func basicFn_listFilesFromBucketID(L *lua.LState) int {
	var bid uuid.UUID
	if L.GetTop() == 2 {
		switch v := L.CheckUserData(2).Value.(type) {
		case uuid.UUID:
			bid = v
		case string:
			bid = uuid.FromStringOrNil(v)
		}
	}

	ud := L.NewUserData()
	ud.Value = filesByBucketID(bid)
	L.Push(ud)
	return 1
}

// найти файл по имени бакета и имени файла
func basicFn_FindFileByName(L *lua.LState) int {
	/*
		bucketName, fileName string,
		used DataUsed,
	*/
	// var bucketName, fileName string
	// var used interfaces.DataUsed

	// bucketName = L.CheckString(2)
	// fileName = L.CheckString(3)

	return 0
}

func basicFn_FindFile(L *lua.LState) int {
	return 0
}

func basicFn_CreateFileFrom(L *lua.LState) int {
	return 0
}

func basicFn_UpdateFileFrom(L *lua.LState) int {
	return 0
}

func basicFn_DeleteFile(L *lua.LState) int {
	return 0
}
