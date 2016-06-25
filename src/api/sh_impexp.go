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
	"io/ioutil"
	"net/http"
	"store"
	"time"
	"utils"
)

var (
	ImportExportSecionNameKey = "importexport"

	ImportExportFileNameZipArchive = ".faderdata"

	ImportExportImportSpecialHandlerName = "importexport.import"
	ImportExportExportSpecialHandlerName = "importexport.export"

	ImportExportImportRouteName = "AppImport"
	ImportExportExportRouteName = "AppExport"

	ImportExportLatestVersionArchiveURL = "https://s3.eu-central-1.amazonaws.com/releases.fader.inpime.com/archives/FADER(sys).dev.latest.zip"
)

//

// IsIncludeInGroupBucketImportExport является ли файл системным согласно настройкам
func IsIncludeInGroupBucketImportExport(groupName, bucketName string) bool {
	config := appSettings.M(ImportExportSecionNameKey).M(groupName)

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

// IsIncludeInGroupFileImportExport является ли файл системным согласно настройкам
func IsIncludeInGroupFileImportExport(groupName, bucketName, fileName string) bool {
	if !IsIncludeInGroupBucketImportExport(groupName, bucketName) {
		return false
	}

	config := appSettings.M(ImportExportSecionNameKey).M(groupName)

	bucketConfig := config.M(bucketName)

	if bucketConfig.Bool("all") {
		return true
	}

	bucketFiles := utils.NewA(bucketConfig.Strings("files"))

	return bucketFiles.Include(fileName)
}

// ListGroupsImportExport возвращает список групп указанных в настройках приложения
func ListGroupsImportExport() []string {
	return appSettings.Keys(ImportExportSecionNameKey)
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

// AppImportFromLastArchive загрузить последнюю версию консоль панели и загрузить ее
func AppImportFromLastArchive() error {

	logrus.Infof("Downloading on link %q...", ImportExportLatestVersionArchiveURL)

	data, err := loadLatestArchive(ImportExportLatestVersionArchiveURL)

	if err != nil {
		logrus.Errorf("Error download, more details: %q", err)
		return err
	}

	logrus.Info("OK. Importing...")

	return makeImportImportExport(data)
}

// makeImportImportExport выполнить импорт из архива
func makeImportImportExport(data []byte) error {
	archive := newArchivePkg()
	err := archive.Import(data)

	if err != nil {
		return err
	}

	for _, _bucket := range archive.Buckets {
		bucket, err := store.BucketByName(_bucket.Name)

		if err == dbox.ErrNotFound {
			logrus.Infof("import: Create a bucket %q", _bucket.Name)

			bucket.SetID(_bucket.ID)
			bucket.SetName(_bucket.Name)
			bucket.SetBucket(_bucket.Bucket)
			bucket.Import(_bucket.Data)

			bucket.InitRawDataStore(bucket.GetRawDataStoreType(),
				bucket.GetRawDataStoreNameWithoutPostfix())
			bucket.InitMetaDataStore(bucket.GetMetaDataStoreType(),
				bucket.GetMetaDataStoreNameWithoutPostfix())
			bucket.InitMapDataStore(bucket.GetMapDataStoreType(),
				bucket.GetMapDataStoreNameWithoutPostfix())

			bucket.UpdateMapping()
			bucket.Sync()
		}
	}

	time.Sleep(time.Second * 3)

	for _, _file := range archive.Files {
		file, err := store.LoadOrNewFile(_file.Bucket, _file.Name)

		if err == dbox.ErrNotFound {
			file.SetID(_file.ID)
			file.SetName(_file.Name)
			file.SetBucket(_file.Bucket)
		}

		logrus.Infof("import: Upsert a file %q@%q", _file.Bucket, _file.Name)

		file.Import(_file.Data)
		file.Sync()
	}

	return nil
}

func AppImport_SpecialHandler(ctx *ContextWrap) error {
	fileData := ctx.FormFileData("BinData")

	err := makeImportImportExport(fileData.Data)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func AppExport_SpecialHandler(ctx *ContextWrap) error {
	archive := newArchivePkg()
	archive.GroupName = ctx.QueryParam("group")
	byGroupName := len(archive.GroupName) > 0

	for _, bucket := range getAllBuckets() {
		if byGroupName && !IsIncludeInGroupBucketImportExport(archive.GroupName, bucket.Name()) {
			continue
		}

		logrus.Infof("export: bucket %q", bucket.Name())

		archive.Buckets = append(archive.Buckets, newArchiveFileFromFile(bucket))

		for _, file := range getAllFiles(bucket.Name()) {

			if byGroupName && !IsIncludeInGroupFileImportExport(archive.GroupName, bucket.Name(), file.Name()) {
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

// helper functions

// loadLatestArchive load latest version archive
func loadLatestArchive(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("not successful request, got %q, want `200 OK`", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
