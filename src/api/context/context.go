package context

import (
	"github.com/labstack/echo"
	"net/http"
)

// NewContext wrap echo context
func NewContext(ctx echo.Context) *Context {
	wrapCtx := &Context{
		Context: ctx,
		// Props:   utils.Map(map[string]interface{}{}),
	}

	// for _, key := range ctx.ParamNames() {
	// 	wrapCtx.Set(key, ctx.Param(key))
	// }

	return wrapCtx
}

type Context struct {
	echo.Context

	// Props utils.M
}

func (c Context) IsPut() bool {

	return c.Request().Method() == http.MethodPut
}

func (c Context) IsPost() bool {
	return c.Request().Method() == http.MethodPost
}

func (c Context) IsDelete() bool {
	return c.Request().Method() == http.MethodDelete
}

func (c Context) IsGet() bool {
	return c.Request().Method() == http.MethodGet
}
