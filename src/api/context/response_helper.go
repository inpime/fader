package context

import (
	"net/http"
)

func (c *Context) SetStatusOK() *Context {
	return c
}

func (c *Context) SetStatusNotFound() *Context {
	return c
}

func (c *Context) SetStatusBadRequest() *Context {
	return c
}

func (c *Context) SetStatusForbidden() *Context {
	return c
}

func (c *Context) ResponseJSON() *Context {
	return c
}

func (c *Context) ResponseHTML() *Context {
	return c
}

func (c *Context) SetResponseData(data interface{}) *Context {
	return c
}

// Redirect302 redirect to url. Statuc code 302
func (c Context) Redirect302(url string) error {
	return c.Redirect(http.StatusFound, url)
}
