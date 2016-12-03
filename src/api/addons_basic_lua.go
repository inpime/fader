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
	"URL": routeGetURLFromParams,
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

	url, err := r.route.URLPath()
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
