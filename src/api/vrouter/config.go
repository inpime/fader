package vrouter

import (
	"api/config"
	"errors"
	"github.com/labstack/echo"
	"strings"
	"time"
)

var (
	RouteMatchCtxKey    = addonName + ".RouteMatch"
	CSRFCtxKey          = addonName + ".CSRF"
	CSRFFieldNameCtxKey = addonName + ".CSRFFieldName"

	DefaultCookieName   = "csrf"
	DefaultCookiePath   = "/"
	DefaultCookieMaxAge = 86400 // 24H

	// token lookup modes
	lookupSkip   = "skip"
	lookupHeader = "header:"
	lookupForm   = "form:"
)

func MainSettings() *Settings {

	return config.Cfgx.Config(addonName).(*Settings)
}

// AppRouts
func AppRouts() []Rout {

	return MainSettings().Routs
}

// IsSkipOrEmpty true if mode=skip or value is empty
func IsSkipOrEmpty(mode string) bool {
	return strings.HasPrefix(mode, lookupSkip) || len(mode) == 0
}

func IsForm(mode string) bool {
	return strings.HasPrefix(mode, lookupForm)
}

func IsHeader(mode string) bool {
	return strings.HasPrefix(mode, lookupHeader)
}

func ExtractorCSRFToken(mode string) (func(echo.Context) (string, error), string) {
	if IsForm(mode) {
		paramName := strings.Replace(mode, lookupForm, "", -1)
		return func(ctx echo.Context) (string, error) {
			token := ctx.FormValue(paramName)
			if token == "" {
				return "", errors.New("empty csrf token in form param")
			}
			return token, nil
		}, paramName
	}

	if IsHeader(mode) {
		paramName := strings.Replace(mode, lookupHeader, "", -1)
		return func(ctx echo.Context) (string, error) {
			return ctx.Request().Header().Get(paramName), nil
		}, paramName
	}

	return func(ctx echo.Context) (string, error) {
		return "", errors.New("empty csrf token (not supported mode)")
	}, ""
}

// CSRFEnabled
func CSRFEnabled() bool {
	return MainSettings().CSRF.Enabled
}

// CSRFSecretKey
func CSRFSecret() string {
	return MainSettings().CSRF.Secret
}

// CSRFSecretKey
func CSRFLookup() string {
	return MainSettings().CSRF.TokenLookup
}

// Cookie

func CSRFCookieName() string {
	name := MainSettings().CSRF.Cookie.Name

	if len(strings.TrimSpace(name)) == 0 {
		name = DefaultCookieName
	}

	return name
}

func CSRFCookiePath() string {
	path := MainSettings().CSRF.Cookie.Path

	if len(path) == 0 {
		path = DefaultCookiePath
	}

	return path
}

func CSRFCookieDomain() string {
	return MainSettings().CSRF.Cookie.Domain
}

func CSRFCookieExpireDate() time.Time {
	maxage := MainSettings().CSRF.Cookie.MaxAge

	if maxage == 0 {
		maxage = time.Duration(DefaultCookieMaxAge)
	}

	return time.Now().Add(maxage * time.Second)
}

func CSRFCookieSecure() bool {
	return MainSettings().CSRF.Cookie.Secure
}

func CSRFHTTPOnly() bool {
	return MainSettings().CSRF.Cookie.HTTPOnly
}
