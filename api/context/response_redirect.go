package context

var HasRedirectKey = "_HasRedirect"

func (c Context) redirectTo() string {
	if url, exists := c.Get(HasRedirectKey).(string); exists {
		return url
	}

	return ""
}

func (c Context) IsRedirect() bool {
	if len(c.redirectTo()) == 0 {
		return false
	}

	return true
}

// RedirectTo return the URL redirect
func (c Context) RedirectTo() string {
	return c.redirectTo()
}

// Redirect302 redirect to url. Statuc code 302
func (c *Context) Redirect302(url string) *Context {
	c.Set(HasRedirectKey, url)
	return c
}
