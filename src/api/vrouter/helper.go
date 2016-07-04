package vrouter

import (
	"github.com/labstack/echo"
)

func GetCurrentRoute(ctx echo.Context) *RouteMatch {
	if r, exists := ctx.Get(RouteMatchCtxKey).(*RouteMatch); exists {
		return r
	}

	return nil
}
