package context

import (
	"net/http"

	"github.com/inpime/sdata"
	"github.com/labstack/echo"
)

// NewContext wrap echo context
func NewContext(ctx echo.Context) *Context {
	wrapCtx := &Context{
		Context: ctx,
		JSON:    sdata.NewStringMap(),
		// Props:   utils.Map(map[string]interface{}{}),
	}

	// for _, key := range ctx.ParamNames() {
	// 	wrapCtx.Set(key, ctx.Param(key))
	// }

	return wrapCtx
}

type Context struct {
	echo.Context

	// Route *vrouter.RouteMatch

	JSON *sdata.StringMap
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

// Set saves data in the context.
func (c Context) Set(key string, v interface{}) Context {
	c.Context.Set(key, v)
	return c
}
