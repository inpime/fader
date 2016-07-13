package store

import (
	"github.com/inpime/dbox"
)

var _ dbox.Store = (*LocalStore)(nil)
var _ dbox.FileStore = (*LocalStore)(nil)

type LocalStore struct {
	store dbox.Store
}

func NewLocalStore(s interface{}, path string) *LocalStore {
	return &LocalStore{
		store: dbox.NewLocalStore(path),
	}
}

func (s LocalStore) GetFile(id string, f *File) error {

	return nil
}

func (s LocalStore) GetByNameFile(name string, f *File) error {

	return nil
}

func (s LocalStore) SaveFile(f *File) error {
	return nil
}

func (s LocalStore) DeleteFile(f *File) error {
	return nil
}

// dbox.Store implement interface

func (s LocalStore) GetByName(id string, obj dbox.Object) error {
	return s.store.(dbox.FileStore).GetByName(id, obj)
}

func (s LocalStore) Get(id string, obj dbox.Object) error {
	return s.store.Get(id, obj)
}

func (s LocalStore) Save(obj dbox.Object) error {
	// save to elastic search

	return s.store.Save(obj)
}

func (s LocalStore) Delete(obj dbox.Object) error {
	// save to elastic search

	return s.store.Delete(obj)
}

func (s LocalStore) Type() dbox.StoreType {
	return s.store.Type()
}
