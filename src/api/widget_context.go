package api

import (
	"github.com/labstack/echo"
	"strings"
	"utils"
)

func NewContextWrap(c echo.Context) *ContextWrap {
	ctx := &ContextWrap{
		Context: c,
		Props:   utils.Map(map[string]interface{}{}),
	}

	for _, key := range c.ParamNames() {
		ctx.Set(key, c.Param(key))
	}

	// TODO: not best solution

	// copy ctx store
	ctx.Set("session", c.Get("session"))

	return ctx
}

type ContextWrap struct {
	echo.Context

	Props utils.M
}

func (c ContextWrap) Set(key string, value interface{}) {
	c.Props.Set(key, value)
}

func (c ContextWrap) Get(key string) interface{} {
	return c.Props.Get(key)
}

func (c ContextWrap) IsPut() bool {
	return strings.ToLower(c.Request().Method()) == "put"
}

func (c ContextWrap) IsPost() bool {
	return strings.ToLower(c.Request().Method()) == "post"
}

func (c ContextWrap) IsDelete() bool {
	return strings.ToLower(c.Request().Method()) == "delete"
}

func (c ContextWrap) IsGet() bool {
	return strings.ToLower(c.Request().Method()) == "get"
}

func (c *ContextWrap) Session() *Session {
	return currentSessionFromContext(c)
}

func (c ContextWrap) CurrentUser() *User {
	return c.Session().User()
}
