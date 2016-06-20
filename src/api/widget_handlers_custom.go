package api

import (
	"github.com/labstack/echo"
	"net/http"
	"store"
	"time"
)

// UserContentHandler returns the file content (raw data file) without access checks
func UserContentHandler(ctx *ContextWrap) error {

	fileId, isValid := ctx.Get("fileid").(string)

	if !isValid {
		return ctx.NoContent(http.StatusNotFound)
	}

	bucketName := appSettings.M("usercontent").String("bucket")

	file, err := store.LoadOrNewFileID(bucketName, fileId)

	if err != nil || file.IsNew() {
		return ctx.NoContent(http.StatusNotFound)
	}

	// TODO: check current session access

	if t, err := time.Parse(http.TimeFormat, ctx.Request().Header().Get(echo.HeaderIfModifiedSince)); err == nil && file.UpdatedAt().Before(t.Add(1*time.Second)) {
		ctx.Response().Header().Del(echo.HeaderContentType)
		ctx.Response().Header().Del(echo.HeaderContentLength)
		return ctx.NoContent(http.StatusNotModified)
	}

	ctx.Response().Header().Set(echo.HeaderLastModified, file.UpdatedAt().UTC().Format(http.TimeFormat))
	ctx.Response().Header().Add(echo.HeaderContentType, file.ContentType())
	ctx.Response().WriteHeader(http.StatusOK)

	_, err = ctx.Response().Write(file.RawData().Bytes())

	return err
}
