package api

import (
	"api/config"
	"api/context"
	"api/vrouter"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

var (
	NotFoundHandler = func(c echo.Context) error {

		return c.String(http.StatusNotFound, "Not Found")
	}

	ForbiddenHandler = func(c echo.Context) error {

		return c.String(http.StatusForbidden, "Forbidden")
	}

	MaintenanceHandler = func(c echo.Context) error {
		c.Response().Header().Set("Retry-After", "3600") // retry after 1 hourse
		return c.String(http.StatusServiceUnavailable, "Service Unavailable")
	}

	InternalErrorHandler = func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
)

func AppEntryPointHandler(ctx echo.Context) error {
	var match vrouter.RouteMatch

	_url, _ := url.Parse(ctx.Request().URI())

	if config.Router.Match(&vrouter.Request{_url, ctx.Request().Method()}, &match) {

		for key, value := range match.Vars {
			ctx.Set(key, value)
		}

		// If special handler transfer control him

		if specialHandler, err := GetSpecialHandler(match.Handler.SpecialHandler); err == nil {
			return specialHandler(ctx)
		}

		// Template page

		var tpl *pongo2.Template
		var _ctx = context.NewContext(ctx)

		pongo2.DefaultSet.Debug = !config.IsPageCaching()

		// if Debug true then recompile tpl on any request
		tpl, err := pongo2.FromCache(match.Handler.Bucket + "/" + match.Handler.File)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"_service": "api",
				"_target":  "specialhandler",
				"handler":  match.Handler.String(),
			}).WithError(err).Error("get tempalte file")

			return InternalErrorHandler(ctx)
		}

		res, err := tpl.Execute(pongo2.Context{
			"ctx": _ctx,
		})

		if err != nil {
			// TODO: Custom error
			logrus.WithFields(logrus.Fields{
				"_service": "api",
				"_target":  "specialhandler",
				"handler":  match.Handler.String(),
			}).WithError(err).Error("execute template")
			return err
		}

		// Custom response if exist
		// TODO:

		// TODO: Custom header

		return ctx.HTML(http.StatusOK, res)
	}

	// TODO: HandlerNotFound

	return NotFoundHandler(ctx)
}

// -----------------------------------
// Special handlers
// -----------------------------------

type SpecialHandler echo.HandlerFunc

var registredSpecialHandlers = map[string]SpecialHandler{}

func AddSpecialHandler(name string, fn func(ctx echo.Context) error) {
	logrus.WithFields(logrus.Fields{
		"_service": "api",
		"_target":  "initspecialhandler",
		"handler":  name,
	}).Info("add special handler")

	registredSpecialHandlers[name] = fn
}

func GetSpecialHandler(name string) (SpecialHandler, error) {
	fn, exists := registredSpecialHandlers[name]

	if !exists {
		return NotFoundHandler, fmt.Errorf("not_found")
	}

	return fn, nil
}
