package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestApiStratey_simple(t *testing.T) {

	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	setupTestData(
		`c = ctx()
	c:Set("name", "fader")`,
		`Hello {{ ctx.Get("name") }} {{ ctx.Get("id") }} {{ ctx.QueryParam("c") }} !`,
	)

	s, b := request(echo.GET, "/fc/abc-def.123_456?a=b&c=d'd'd;", e)
	assert.Equal(t, http.StatusOK, s)
	assert.Equal(t, []byte(`Hello fader abc-def.123_456 d&#39;d&#39;d !`), b)
}
