package context

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
