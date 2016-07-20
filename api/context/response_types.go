package context

import (
	"net/http"
)

const (
	ContentTypeKey     = "_ResponseContentType"
	DefaultContentType = "html"

	StatusKey     = "_ResponseStatus"
	DefaultStatus = http.StatusOK
)

func (c *Context) SetStatusOK() *Context {
	c.Set(StatusKey, http.StatusOK)
	return c
}

func (c *Context) SetStatusNotFound() *Context {
	c.Set(StatusKey, http.StatusNotFound)
	return c
}

func (c *Context) SetStatusBadRequest() *Context {
	c.Set(StatusKey, http.StatusBadRequest)
	return c
}

func (c *Context) SetStatusForbidden() *Context {
	c.Set(StatusKey, http.StatusForbidden)
	return c
}

func (c *Context) ResponseJSON() *Context {
	c.Set(ContentTypeKey, "json")
	return c
}

func (c *Context) ResponseHTML() *Context {
	c.Set(ContentTypeKey, "html")
	return c
}

func (c *Context) SetResponseJSON(data interface{}) *Context {
	c.ResponseJSON()
	c.JSON.LoadFrom(data)
	return c
}

func (c *Context) ResponseStatus() int {
	if status, ok := c.Get(StatusKey).(int); ok && status > 0 {
		return status
	}

	return DefaultStatus
}

func (c *Context) ResponseType() string {
	if _type, ok := c.Get(ContentTypeKey).(string); ok && len(_type) > 0 {
		return _type
	}

	return DefaultContentType
}
