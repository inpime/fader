package router

import (
	"interfaces"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

var (
	NotFoundHandler echo.HandlerFunc
)

func VRouterMiddleware(router interfaces.RouteMatcher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var match = &interfaces.RouteMatch{}

			_url, _ := url.Parse(ctx.Request().URL().Path())
			_url.RawQuery = ctx.Request().URL().QueryString()

			requestParams := interfaces.RequestParams{
				URL:    _url,
				Method: ctx.Request().Method(),
			}

			log.Println("[DEBUG]", requestParams)

			if router.Match(requestParams, match) {
				for key, value := range match.Vars {
					ctx.Set(key, value)
				}

				SetMatchedVRouteContext(ctx, match)

				return next(ctx)
			}

			if NotFoundHandler == nil {
				return ctx.NoContent(http.StatusBadRequest)
			}

			return NotFoundHandler(ctx)
		}
	}
}
