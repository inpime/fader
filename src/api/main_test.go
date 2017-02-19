package api

import (
	"addons"
	"interfaces"
	"io"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	e = echo.New()
	TESTING = true

	os.Exit(m.Run())
}

func TestCountBuckets_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	count := 0
	bucketManager.(interfaces.BucketImportManager).
		EachBucket(func(b *interfaces.Bucket) error {
			logger.Println("[INFO] bucket name:", b.BucketName)
			count++
			return nil
		})

	assert.Equal(t, 0, count)

	tmpbucket := interfaces.NewBucket()
	tmpbucket.BucketID = uuid.NewV4()
	tmpbucket.BucketName = "a"
	if err := bucketManager.CreateBucket(tmpbucket); err != nil {
		logger.Panicln("[FAIL] create bucket", err)
	}

	count = 0
	bucketManager.(interfaces.BucketImportManager).
		EachBucket(func(b *interfaces.Bucket) error {
			logger.Println("[INFO] bucket name:", b.BucketName)
			count++
			return nil
		})

	assert.Equal(t, 1, count)
}

func request(method, path string, e *echo.Echo) (int, []byte) {
	req := test.NewRequest(method, path, nil)
	rec := test.NewResponseRecorder()
	e.ServeHTTP(req, rec)
	return rec.Status(), rec.Body.Bytes()
}

// internal

func setupLuaContext(method, url string, d io.Reader, L *lua.LState) *Context {
	// strings.NewReader(userJSON)
	e := echo.New()
	req := test.NewRequest(method, url, d)
	rec := test.NewResponseRecorder()
	ctx := e.NewContext(req, rec)

	return ContextLuaExecutor(L, ctx)
}

func setupLuaModules(L *lua.LState) {
	for _, addon := range addons.Addons {
		L.PreloadModule(addon.Name(), addon.LuaLoader)
	}
}
