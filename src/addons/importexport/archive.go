package importexport

import (
	"api/config"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"store"
	"time"
)

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

	f, err := w.Create(ArchiveFaderDataFileName)
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

	if r.File[0].Name != ArchiveFaderDataFileName {
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
		AppVersion: config.AppVersion,
	}
}
