// special handlers
package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type SpecialHandler func(*ContextWrap) error

var registredSpecialHandlers = map[string]SpecialHandler{}

func RegistedSpecialHandler(name string, fn SpecialHandler) {
	registredSpecialHandlers[name] = fn
}

func ExistSpecialHandler(name string) bool {
	_, exists := registredSpecialHandlers[name]
	return exists
}

func GetSpecialHandler(name string) SpecialHandler {
	fn, exists := registredSpecialHandlers[name]

	if !exists {
		return func(c *ContextWrap) error {

			return c.NoContent(http.StatusNotFound)
		}
	}

	return fn
}

//

func EmptyHandler(c echo.Context) error {

	return c.String(200, "empty")
}
