package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestApiGlobal_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	setupSysConfigFilesCase2()

	err = appConfigUpdateFn()
	assert.NoError(t, err)

	s, b := request(echo.GET, "/route22/abc-def.123_456?a=b&c=d'd'd;", e)
	assert.Equal(t, http.StatusOK, s)
	assert.Equal(t, []byte(`Hello check extension &#39;example&#39; fader abc-def.123_456 d&#39;d&#39;d !`), b)
	//                            | addons function               | |lua| |route param|   |url param  |
}
