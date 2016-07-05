package vrouter

import (
	"api/config"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"net/url"
)

var appRouter = NewRouter()

func AppRouter() *Router {

	return appRouter
}

func RouterMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var match RouteMatch

			_url, _ := url.Parse(ctx.Request().URI())

			logrus.WithField("_service", addonName).Debugf("count routs %d", len(AppRouter().routes))

			if AppRouter().Match(&Request{_url, ctx.Request().Method()}, &match) {
				for key, value := range match.Vars {
					ctx.Set(key, value)
				}

				ctx.Set(RouteMatchCtxKey, &match)

				// -----------
				// CSRF
				// -----------

				if !CSRFEnabled() {
					logrus.WithField("_service", addonName).
						Debug("csrf disabled, next")
					return next(ctx)
				}

				// cookie

				salt, err := generateSalt(8)

				if err != nil {
					logrus.WithError(err).
						WithField("_service", addonName).
						Error("generate csrf token salt")

					return config.InternalErrorHandler(ctx)
				}

				token := generateCSRFToken([]byte(CSRFSecret()), salt)
				ctx.Set(CSRFCtxKey, token)
				ctx.SetCookie(csrfCookie(token))

				logrus.WithField("uri", _url.String()).Infof("set cookie")

				// extract and check

				if !match.Handler.CSRF {
					return next(ctx)
				}

				tokenLookup := CSRFLookup()

				// custom lookup
				if len(match.Handler.CSRFTokenLookup) > 0 {
					tokenLookup = match.Handler.CSRFTokenLookup
				}

				extractor, fieldName := ExtractorCSRFToken(tokenLookup)
				ctx.Set(CSRFFieldNameCtxKey, fieldName)

				switch ctx.Request().Method() {
				case echo.GET, echo.HEAD, echo.OPTIONS, echo.TRACE:
				default:
					token, err := extractor(ctx)
					if err != nil {
						logrus.WithField("_service", addonName).
							WithError(err).
							Error("extract token")

						return config.InternalErrorHandler(ctx)
					}
					ok, err := validateCSRFToken(token, []byte(CSRFSecret()))
					if err != nil {
						logrus.WithField("_service", addonName).
							WithError(err).
							Error("validator token")
						return config.InternalErrorHandler(ctx)
					}
					if !ok {
						return config.ForbiddenHandler(ctx)
					}
				}

				return next(ctx)
			}

			logrus.WithField("_service", addonName).Infof("NOTFOUND. count routs %d", len(AppRouter().routes))

			return config.NotFoundHandler(ctx)
		}
	}
}

// Cookie CSRF

func csrfCookie(token string) *echo.Cookie {
	cookie := new(echo.Cookie)
	cookie.SetName(CSRFCookieName())
	cookie.SetValue(token)
	_path := CSRFCookiePath()
	if _path != "" {
		cookie.SetPath(_path)
	}
	_domain := CSRFCookieDomain()
	if _domain != "" {
		cookie.SetDomain(_domain)
	}
	cookie.SetExpires(CSRFCookieExpireDate())
	cookie.SetSecure(CSRFCookieSecure())
	cookie.SetHTTPOnly(CSRFHTTPOnly())

	return cookie
}

// Config routs

func ReloadAppRouts() {
	router := NewRouter()

	refreshRouter(router)

	appRouter = router
}

// refreshRouter returns the current routing based on app settings
func refreshRouter(router *Router) {

	if len(AppRouts()) == 0 {
		logrus.WithField("_serivce", addonName).Debugf("app config empty routs %d", len(AppRouts()))
	}

	for _, _r := range AppRouts() {

		handler := NewHandlerFromRoute(_r)

		logrus.WithFields(logrus.Fields{
			"_service": addonName,
			"handler":  handler,
		}).Debug("Add route")

		if len(_r.Methods) == 0 {
			router.Handle(_r.Path, handler).Methods("GET").Name(_r.Name)
		} else {
			router.Handle(_r.Path, handler).Methods(_r.Methods...).Name(_r.Name)
		}
	}
}

func NewHandlerFromRoute(r Rout) Handler {
	var h Handler
	if r.IsSpecial {
		h = Handler{
			Bucket:         "",
			File:           "",
			SpecialHandler: r.Handler,
		}
	} else {
		h = NewHandlerFromString(r.Handler)
	}

	h.Licenses = r.Licenses
	h.Path = r.Path
	h.Methods = r.Methods
	h.CSRF = r.CSRF
	h.CSRFTokenLookup = r.CSRFTokenLookup

	return h
}
