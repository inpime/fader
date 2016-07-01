package api

import (
	"api/config"
	"github.com/labstack/echo"
	"net/http"
	"store"
	"time"
)

var (
	FileContentSectionNameKey = "filecontent"
	FileContentBucketNameKey  = "bucket"

	FileContentByNameSpecialHandlerName = "filecontent.byname"
	FileContentByIDSpecialHandlerName   = "filecontent.byid"

	FileContentByNameRouteName = "FileContentByName"
	FileContentByIDRouteName   = "FileContentByID"
)

// FileContentByNameHandler returns the file content by name (raw data file) without access checks
//
// special handler name ''
func FileContentByName_SpecialHandler(ctx *ContextWrap) error {

	fileName, isValid := ctx.Get("file").(string)

	if !isValid {
		return ctx.NoContent(http.StatusNotFound)
	}

	bucketName := config.AppSettings().M(FileContentSectionNameKey).String(FileContentBucketNameKey)

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
func FileContentByID_SpecialHandler(ctx *ContextWrap) error {

	fileId, isValid := ctx.Get("file").(string)

	if !isValid {
		return ctx.NoContent(http.StatusNotFound)
	}

	bucketName := config.AppSettings().M(FileContentSectionNameKey).String(FileContentBucketNameKey)

	file, err := store.LoadOrNewFileID(bucketName, fileId)

	if err != nil || file.IsNew() {
		return ctx.NoContent(http.StatusNotFound)
	}

	// TODO: check current session access

	return responseFileContentWithLastModifiedHeader(ctx, file)
}

func responseFileContentWithLastModifiedHeader(ctx *ContextWrap, file *store.File) error {
	// TODO: check current session access

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
