package store

import (
	"github.com/inpime/dbox"
)

var (
	ContentTypeKey = "ContentType"
)

var _ dbox.Object = (*File)(nil)

type File struct {
	*dbox.File
}

func MustFile(file *dbox.File) *File {
	return &File{file}
}

func NewFile(store dbox.Store) *File {
	return &File{
		File: dbox.NewFile(store),
	}
}

func (f *File) Sync() error {
	if err := f.File.Sync(); err != nil {
		return err
	}

	return UpdateSearchDocument(f.Bucket(), f.ID(), FileSearchFromFile(f))
}

func (f *File) Delete() error {
	if err := f.File.Delete(); err != nil {
		return err
	}

	return DeleteSearchDocument(f.Bucket(), f.ID())
}
