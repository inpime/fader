package router

import (
	"interfaces"

	"github.com/labstack/echo"
)

var (
	MatchRouteCtxKey = "_MatchRoute"
)

func MatchVRouteFromContext(
	ctx echo.Context,
) *interfaces.RouteMatch {

	if route, exists := ctx.Get(MatchRouteCtxKey).(*interfaces.RouteMatch); exists {
		return route
	}

	return nil
}

func SetMatchedVRouteContext(
	ctx echo.Context,
	match *interfaces.RouteMatch,
) {
	ctx.Set(MatchRouteCtxKey, match)
}
