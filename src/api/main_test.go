package api

import (
	"interfaces"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
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
