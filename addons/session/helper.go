package session

import (
	"github.com/labstack/echo"
)

// GetSession get session object of the context
func GetSession(ctx echo.Context) *Session {
	_session := ctx.Get(SessionNameContextKey)

	if _session, ok := _session.(*Session); ok {
		return _session
	}

	return nil
}
