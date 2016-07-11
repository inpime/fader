package api

import (
	"addons/standard"
	"api/config"
	"api/context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"net/http"
)

// AppEntryPointHandler the entry point for the application
// 	* Check access licenses
//	* Special handler if exist
// 	* Execute template
// 	* Flush session data
// 	* Custom repsone type
func AppEntryPointHandler(ctx echo.Context) error {
	uri := ctx.Request().URI()

	logrus.WithFields(logrus.Fields{
		"_breakpoint": "AppEntryPointHandler",
		"uri":         uri,
	}).Debug("trace")

	var _ctx = context.NewContext(ctx)
	var match = _ctx.CurrentRoute()

	// ------------------------
	// Check access licenses
	// ------------------------

	logrus.WithFields(logrus.Fields{"_breakpoint": "HasOneLicense", "uri": uri}).Debug("trace")

	if len(match.Handler.Licenses) > 0 && !_ctx.Session().HasOneLicense(match.Handler.Licenses) {

		return config.ForbiddenHandler(ctx)
	}

	// ------------------------
	// Special handler if exist
	// ------------------------

	logrus.WithFields(logrus.Fields{"_breakpoint": "GetSpecialHandler", "uri": uri}).Debug("trace")

	if specialHandler, err := GetSpecialHandler(match.Handler.SpecialHandler); err == nil {
		return specialHandler(ctx)
	}

	// ------------------------
	// Execute template page
	// ------------------------

	var tpl *pongo2.Template

	logrus.WithFields(logrus.Fields{"_breakpoint": "get MainSettings.TplCache", "uri": uri}).Debug("trace")

	pongo2.DefaultSet.Debug = !standard.MainSettings().TplCache

	logrus.WithFields(logrus.Fields{"_breakpoint": "get file", "uri": uri, "_file": match.Handler.Bucket + "/" + match.Handler.File}).Debug("trace")

	// if Debug true then recompile tpl on any request
	tpl, err := pongo2.FromCache(match.Handler.Bucket + "/" + match.Handler.File)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"_service": "api",
			"handler":  match.Handler.String(),
		}).WithError(err).Error("get tempalte file")

		return config.InternalErrorHandler(ctx)
	}

	logrus.WithFields(logrus.Fields{"_breakpoint": "before execute file", "uri": uri, "_file": match.Handler.Bucket + "/" + match.Handler.File}).Debug("trace")

	res, err := tpl.Execute(pongo2.Context{
		"ctx": _ctx,
	})

	logrus.WithFields(logrus.Fields{"_breakpoint": "after execute file", "uri": uri, "_file": match.Handler.Bucket + "/" + match.Handler.File}).Debug("trace")

	if err != nil {
		// TODO: Custom error
		logrus.WithFields(logrus.Fields{
			"_service": "api",
			"handler":  match.Handler.String(),
		}).WithError(err).Error("execute template")

		return config.InternalErrorHandler(ctx)
	}

	// ------------------------
	// Flush session data
	// ------------------------

	logrus.WithFields(logrus.Fields{"_breakpoint": "update session", "uri": uri}).Debug("trace")

	if err := _ctx.Session().Save(); err != nil {
		logrus.WithFields(logrus.Fields{
			"_service": "api",
		}).WithError(err).Error("save session after get the flash messages")
	}

	// redirect if specified

	logrus.WithFields(logrus.Fields{"_breakpoint": "redirect if exist", "uri": uri}).Debug("trace")

	if _ctx.IsRedirect() {

		return _ctx.Redirect(http.StatusFound, _ctx.RedirectTo())
	}

	// ------------------------
	// Custom repsone type
	// 	* JSON
	// 	* Text
	// 	* Byte
	// ------------------------

	// TODO: Custom header

	logrus.WithFields(logrus.Fields{"_breakpoint": "response", "uri": uri}).Debug("trace")

	return ctx.HTML(http.StatusOK, res)
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
		return config.NotFoundHandler, fmt.Errorf("not_found")
	}

	return fn, nil
}
