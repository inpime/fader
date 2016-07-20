package filestatic

import (
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/store"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func getBucketNameFromHandler(ctx echo.Context) string {
	if bucketName, ok := ctx.Get(config.SpecialHandlerArgsKey).(string); ok {

		return bucketName
	}

	return ""
}

func storeBucketName(ctx echo.Context) string {
	bucketName := getBucketNameFromHandler(ctx)

	if len(bucketName) == 0 {
		bucketName = MainSettings().BucketSource
	}

	return bucketName
}

// FileContentByNameHandler returns the file content by name (raw data file) without access checks
//
// special handler name ''
func FileContentByNameHandler(ctx echo.Context) error {

	fileName, isValid := ctx.Get("file").(string)

	if !isValid {
		return ctx.NoContent(http.StatusNotFound)
	}

	bucketName := storeBucketName(ctx)

	file, err := store.LoadOrNewFile(bucketName, fileName)

	if err != nil || file.IsNew() {
		return ctx.NoContent(http.StatusNotFound)
	}

	// TODO: check current session access

	return responseFileContentWithLastModifiedHeader(ctx, file)
}

// FileContentByIDHandler returns the file content by id (raw data file) without access checks
//
// special handler name ''
func FileContentByIDHandler(ctx echo.Context) error {

	fileId, isValid := ctx.Get("file").(string)

	if !isValid {
		return ctx.NoContent(http.StatusNotFound)
	}

	bucketName := storeBucketName(ctx)

	file, err := store.LoadOrNewFileID(bucketName, fileId)

	if err != nil || file.IsNew() {
		return ctx.NoContent(http.StatusNotFound)
	}

	// TODO: check current session access

	return responseFileContentWithLastModifiedHeader(ctx, file)
}

func responseFileContentWithLastModifiedHeader(ctx echo.Context, file *store.File) error {

	// cache
	if t, err := time.Parse(http.TimeFormat, ctx.Request().Header().Get(echo.HeaderIfModifiedSince)); err == nil && file.UpdatedAt().Before(t.Add(1*time.Second)) {
		ctx.Response().Header().Del(echo.HeaderContentType)
		ctx.Response().Header().Del(echo.HeaderContentLength)
		return ctx.NoContent(http.StatusNotModified)
	}

	ctx.Response().Header().Set(echo.HeaderLastModified, file.UpdatedAt().UTC().Format(http.TimeFormat))
	ctx.Response().Header().Add(echo.HeaderContentType, file.ContentType())
	ctx.Response().WriteHeader(http.StatusOK)

	_, err := ctx.Response().Write(file.RawData().Bytes())

	return err
}
