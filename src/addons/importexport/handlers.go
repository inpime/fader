package importexport

import (
	"addons/search"
	"api/context"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"net/http"
)

func ImportHandler(ctx echo.Context) error {
	fileData := context.NewContext(ctx).FormFileData("BinData")

	err := makeImportImportExport(fileData.Data)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func ExportHandler(ctx echo.Context) error {
	archive := newArchivePkg()
	archive.GroupName = ctx.QueryParam("group")
	byGroupName := len(archive.GroupName) > 0

	for _, bucket := range search.GetAllBuckets() {
		if byGroupName && !IsIncludeInGroupBucketImportExport(archive.GroupName, bucket.Name()) {
			continue
		}

		logrus.WithFields(logrus.Fields{
			"_service": addonName,
		}).Infof("export: bucket %q", bucket.Name())

		archive.Buckets = append(archive.Buckets, newArchiveFileFromFile(bucket))

		for _, file := range search.GetAllFiles(bucket.Name()) {

			if byGroupName && !IsIncludeInGroupFileImportExport(archive.GroupName, bucket.Name(), file.Name()) {
				continue
			}

			logrus.WithFields(logrus.Fields{
				"_service": addonName,
			}).Infof("export: \t file %q", file.Name())

			archive.Files = append(archive.Files, newArchiveFileFromFile(file))
		}
	}

	ctx.Response().Header().Add(echo.HeaderContentType, "application/zip")
	ctx.Response().Header().Add(echo.HeaderContentDisposition, "attachment; filename="+archive.FileName())
	ctx.Response().Header().Add("Content-Transfer-Encoding", "binary")
	ctx.Response().Header().Add("Expires", "0")
	ctx.Response().WriteHeader(http.StatusOK)
	b, err := archive.Export()

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	_, err = ctx.Response().Write(b)

	return err
}
