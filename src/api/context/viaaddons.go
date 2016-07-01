package context

import (
	"addons/session"
)

// Session get current session
func (c *Context) Session() *session.Session {
	return session.GetSession(c.Context)
}
