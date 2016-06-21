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

// ----------------------------
// Import export application
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

	AppVersion     string
	ConsoleVersion string
	LicenseName    string

	Buckets []fileArchive
	Files   []fileArchive
}

func (a archivePackage) FileName() string {
	versions := a.AppVersion + "." + a.ConsoleVersion + "." + a.LicenseName
	return "F." + versions + "." + time.Now().Format("2006_01_02_15_04_05") + ".zip"
}

func (a archivePackage) Export() ([]byte, error) {
	b, err := json.Marshal(a)

	if err != nil {
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	f, err := w.Create(".faderdata")
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

	if r.File[0].Name != ".faderdata" {
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
		CreatedAt: time.Now(),
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

func ImportAppHandler(ctx *ContextWrap) error {
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

func ExportAppHandler(ctx *ContextWrap) error {
	archive := newArchivePkg()

	for _, bucket := range getAllBuckets() {
		archive.Buckets = append(archive.Buckets, newArchiveFileFromFile(bucket))

		logrus.Infof("bucket name %q", bucket.Name())
		if bucket.Name() == "full_fs" {
			continue
		}
		for _, file := range getAllFiles(bucket.Name()) {
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
