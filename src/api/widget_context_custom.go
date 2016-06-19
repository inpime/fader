package api

import (
	"net/http"
)

var (
	HTMLContentType = "response:content_type:html"
	JSONContentType = "response:content_type:json"

	ResponseContentTypeKey = "response:content_type"
	ResponseStatusKey      = "response:status"
	ResponseDataKey        = "response:data"
)

// responseContentType
func (c ContextWrap) responseContentType() string {
	res := c.Props.String(ResponseContentTypeKey)
	if len(res) == 0 {
		res = HTMLContentType
	}

	return res
}

func (c *ContextWrap) ResponseHTML() *ContextWrap {
	c.Set(ResponseContentTypeKey, HTMLContentType)

	return c
}

func (c *ContextWrap) ResponseJSON() *ContextWrap {
	c.Set(ResponseContentTypeKey, JSONContentType)

	return c
}

func (c ContextWrap) responseStatus() int {
	res := c.Props.Int(ResponseStatusKey)
	if res == 0 {
		res = http.StatusOK
	}
	return res
}

func (c *ContextWrap) ResponseOK() *ContextWrap {
	c.Set(ResponseStatusKey, http.StatusOK)

	return c
}

func (c *ContextWrap) ResponseNotFound() *ContextWrap {
	c.Set(ResponseStatusKey, http.StatusNotFound)
	return c
}

func (c *ContextWrap) ResponseBad() *ContextWrap {
	c.Set(ResponseStatusKey, http.StatusBadRequest)
	return c
}

func (c *ContextWrap) ResponseForbidden() *ContextWrap {
	c.Set(ResponseStatusKey, http.StatusForbidden)

	return c
}

func (c ContextWrap) responseData() interface{} {
	return c.Get(ResponseDataKey)
}

func (c *ContextWrap) SetResponseData(data interface{}) *ContextWrap {
	c.Set(ResponseDataKey, data)

	return c
}

// Redirect302 redirect to url. Statuc code 302
func (c ContextWrap) Redirect302(url string) error {
	return c.Redirect(http.StatusFound, url)
}
