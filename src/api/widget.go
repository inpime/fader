package api

import (
	// "github.com/Sirupsen/logrus"
	"api/config"
	"api/vrouter"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

func initWidgets() {
	initWidgetVirtualRouts()
}

func ExecuteWidget(c echo.Context) error {
	var match vrouter.RouteMatch

	_u, _ := url.Parse(c.Request().URI())

	if config.Router.Match(&vrouter.Request{_u, c.Request().Method()}, &match) {
		// Init context

		widgetCtx := NewContextWrap(c)

		for key, value := range match.Vars {
			widgetCtx.Set(key, value)
		}

		// Check access
		// userSession := widgetCtx.CurrentUser()
		// logrus.WithFields(logrus.Fields{
		// 	"user_session":    userSession,
		// 	"default_Session": DefaultGuestSession,
		// }).Info("session")

		if !widgetCtx.Session().HasOneLicense(match.Handler.Licenses) {

			return c.NoContent(http.StatusForbidden)
		}

		// Special handler

		if ExistSpecialHandler(match.Handler.SpecialHandler) {
			sHandler := GetSpecialHandler(match.Handler.SpecialHandler)

			return sHandler(widgetCtx)
		}

		// Template execute

		var tpl *pongo2.Template

		if config.IsPageCaching() {
			tpl = pongo2.Must(pongo2.FromCache(match.Handler.Bucket + "/" + match.Handler.File))
		} else {
			ClearWidgetTemplatesCache()
			tpl = pongo2.Must(pongo2.FromFile(match.Handler.Bucket + "/" + match.Handler.File))
		}

		res, err := tpl.Execute(pongo2.Context{
			"ctx": widgetCtx,
		})

		if err != nil {
			return err
		}

		// Custom response if exist

		resContentType := widgetCtx.responseContentType()
		resStatus := widgetCtx.responseStatus()

		switch resContentType {
		case JSONContentType:
			data := widgetCtx.responseData()
			if data == nil {
				data = struct{}{} // default value
			}
			return c.JSON(resStatus, data)
		}

		// TODO: custom header

		return c.HTML(resStatus, res)
	}

	// TODO: If not exist

	return nil
}
