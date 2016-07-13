package context

import (
	"fmt"
	"github.com/inpime/fader/addons/session"
	"github.com/inpime/fader/api/vrouter"
	"strings"
)

// Session get current session
func (c *Context) Session() *session.Session {
	return session.GetSession(c.Context)
}

// CurrentRoute get current route if exist
func (c *Context) CurrentRoute() *vrouter.RouteMatch {
	return vrouter.GetCurrentRoute(c.Context)
}

func (c *Context) CurrentRouteIs(routeName string) bool {
	return strings.ToLower(c.CurrentRoute().Handler.Name) == strings.ToLower(routeName)
}

func (c Context) CSRFToken() string {
	if v, exists := c.Get(vrouter.CSRFCtxKey).(string); exists {
		return v
	}

	return ""
}

func (c Context) CSRFFieldName() string {
	if v, exists := c.Get(vrouter.CSRFFieldNameCtxKey).(string); exists {
		return v
	}

	return ""
}

func (c Context) CSRFField() string {
	return fmt.Sprintf(`<input type="hidden" value="%s" name="%s" />`, c.CSRFToken(), c.CSRFFieldName())
}
