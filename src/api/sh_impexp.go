package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"store"
	"strconv"
	"time"
	"utils"
)

var (
	ImportExportSecionNameKey   = "importexport"
	ImportExportSysGroupNameKey = "sys"

	ImportExportFileNameZipArchive = ".faderdata"

	ImportExportImportSpecialHandlerName = "importexport.import"
	ImportExportExportSpecialHandlerName = "importexport.export"

	ImportExportImportRouteName = "AppImport"
	ImportExportExportRouteName = "AppExport"
)

//

// IsSystemBucketFromImportExport является ли файл системным согласно настройкам
func IsSystemBucketFromImportExport(bucketName string) bool {
	config := appSettings.M(ImportExportSecionNameKey).M(ImportExportSysGroupNameKey)

	if !config.Include(bucketName) {
		return false
	}

	bucketConfig := config.M(bucketName)

	if bucketConfig.Bool("all") {
		return true
	}

	bucketFiles := utils.NewA(bucketConfig.Strings("files"))

	return bucketFiles.Len() > 0
}

// IsSystemFileFromImportExport является ли файл системным согласно настройкам
func IsSystemFileFromImportExport(bucketName, fileName string) bool {
	if !IsSystemBucketFromImportExport(bucketName) {
		return false
	}

	config := appSettings.M(ImportExportSecionNameKey).M(ImportExportSysGroupNameKey)

	bucketConfig := config.M(bucketName)

	if bucketConfig.Bool("all") {
		return true
	}

	bucketFiles := utils.NewA(bucketConfig.Strings("files"))

	return bucketFiles.Include(fileName)
}

// ----------------------------
// Import export data application
// ----------------------------

type fileArchive struct {
	ID     string
	Name   string
	Bucket string
	Data   []byte
}

func newArchiveFileFromFile(file *store.File) fileArchive {
	b, _ := file.Export()
	return fileArchive{
		ID:     file.ID(),
		Name:   file.Name(),
		Bucket: file.Bucket(),
		Data:   b,
	}
}

type archivePackage struct {
	CreatedAt time.Time

	GroupName  string
	AppVersion string

	Buckets []fileArchive
	Files   []fileArchive
}

func (a archivePackage) FileName() string {
	prefix := "FADER"
	if len(a.GroupName) > 0 {
		prefix += "(" + a.GroupName + ")"
	}
	return prefix + "." + a.AppVersion + "." + time.Now().Format("20060102_150405") + ".zip"
}

func (a archivePackage) Export() ([]byte, error) {
	b, err := json.Marshal(a)

	if err != nil {
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	f, err := w.Create(ImportExportFileNameZipArchive)
	f.Write(b)
	w.Close()

	return buf.Bytes(), nil
}

func (a *archivePackage) Import(b []byte) error {
	importFileReader := bytes.NewReader(b)
	r, err := zip.NewReader(importFileReader, importFileReader.Size())

	if err != nil {
		return err
	}

	if len(r.File) == 0 {
		return fmt.Errorf("empty archive")
	}

	if r.File[0].Name != ImportExportFileNameZipArchive {
		return fmt.Errorf("not valid archive")
	}

	zf, err := r.File[0].Open()

	if err != nil {
		return err
	}
	defer zf.Close()

	archiveBuff := bytes.NewBuffer([]byte{})
	io.Copy(archiveBuff, zf)

	return json.Unmarshal(archiveBuff.Bytes(), a)
}

func newArchivePkg() *archivePackage {
	return &archivePackage{
		CreatedAt:  time.Now(),
		AppVersion: AppVersion,
	}
}

func getAllBuckets() []*store.File {
	// all buckets
	filter := store.NewSearchFilter("buckets")
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(100)

	queryRaw := buildSearchQueryFilesByBycket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return makeSearch(filter).GetFiles()
}

func getAllFiles(bucket string) []*store.File {
	// all buckets
	filter := store.NewSearchFilter(bucket)
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(1000)

	queryRaw := buildSearchQueryFilesByBycket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return makeSearch(filter).GetFiles()
}

func AppImport_SpecialHandler(ctx *ContextWrap) error {
	fileData := ctx.FormFileData("BinData")

	archive := newArchivePkg()
	err := archive.Import(fileData.Data)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	for _, _bucket := range archive.Buckets {
		bucket, err := store.BucketByName(_bucket.Name)

		if err == dbox.ErrNotFound {
			bucket.SetID(_bucket.ID)
			bucket.SetName(_bucket.Name)
			bucket.SetBucket(_bucket.Bucket)
			bucket.Import(_bucket.Data)
			bucket.Sync()
		}
	}

	time.Sleep(time.Second * 1)

	for _, _file := range archive.Files {
		file, err := store.LoadOrNewFile(_file.Bucket, _file.Name)

		if err == dbox.ErrNotFound {
			file.SetID(_file.ID)
			file.SetName(_file.Name)
			file.SetBucket(_file.Bucket)
		}

		file.Import(_file.Data)
		file.Sync()
	}

	return ctx.NoContent(http.StatusOK)
}

func AppExport_SpecialHandler(ctx *ContextWrap) error {
	archive := newArchivePkg()
	onlySystemFiles, _ := strconv.ParseBool(ctx.QueryParam("sys"))

	if onlySystemFiles {
		archive.GroupName = "sys"
	}

	for _, bucket := range getAllBuckets() {
		if onlySystemFiles && !IsSystemBucketFromImportExport(bucket.Name()) {
			continue
		}

		logrus.Infof("export: bucket %q", bucket.Name())

		archive.Buckets = append(archive.Buckets, newArchiveFileFromFile(bucket))

		for _, file := range getAllFiles(bucket.Name()) {

			if onlySystemFiles && !IsSystemFileFromImportExport(bucket.Name(), file.Name()) {
				continue
			}

			logrus.Infof("export: \t file %q", file.Name())

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
