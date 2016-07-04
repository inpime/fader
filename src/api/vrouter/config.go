package vrouter

import (
	"api/config"
	"errors"
	"github.com/labstack/echo"
	"strings"
	"time"
	"utils"
)

var (
	RoutsKey = "routs"

	CSRFKey         = "csrf"
	CSRFEnabledKey  = "enabled" //boolean
	CSRFSecretKey   = "secret"  // string
	CSRFTokenLookup = "lookup"  // string
	// header:"X-CSRF-Token"
	// form:"csrf"
	// json:"csrf"
	// skip only in route

	// Cookie section
	CSRFCookieSectionKey = "cookie"
	CookieNameKey        = "name"
	CookiePathKey        = "path"     // string
	CookieDomainKey      = "domain"   // string
	CookieAgeKey         = "age"      // int
	CookieSecureKey      = "secure"   // bool
	CookieHTTPOnlyKey    = "httponly" // bool

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

// AppRouts
func AppRouts() []string {
	return config.AppSettings().Strings(RoutsKey)
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

//

func CSRFConfig() utils.M {
	return config.AppSettings().M(CSRFKey)
}

func CSRFCookieConfig() utils.M {
	return CSRFConfig().M(CSRFCookieSectionKey)
}

// CSRFEnabled
func CSRFEnabled() bool {
	return CSRFConfig().Bool(CSRFEnabledKey)
}

// CSRFSecretKey
func CSRFSecret() string {
	return CSRFConfig().String(CSRFSecretKey)
}

// CSRFSecretKey
func CSRFLookup() string {
	return CSRFConfig().String(CSRFTokenLookup)
}

// Cookie

func CSRFCookieName() string {
	name := CSRFCookieConfig().String(CookieNameKey)

	if len(strings.TrimSpace(name)) == 0 {
		name = DefaultCookieName
	}

	return name
}

func CSRFCookiePath() string {
	path := CSRFCookieConfig().String(CookiePathKey)

	if len(path) == 0 {
		path = DefaultCookiePath
	}

	return path
}

func CSRFCookieDomain() string {
	return CSRFCookieConfig().String(CookieDomainKey)
}

func CSRFCookieExpireDate() time.Time {
	maxage := CSRFCookieConfig().Int(CookieAgeKey)

	if maxage == 0 {
		maxage = DefaultCookieMaxAge
	}

	return time.Now().Add(time.Duration(maxage) * time.Second)
}

func CSRFCookieSecure() bool {
	return CSRFCookieConfig().Bool(CookieSecureKey)
}

func CSRFHTTPOnly() bool {
	return CSRFCookieConfig().Bool(CookieHTTPOnlyKey)
}
