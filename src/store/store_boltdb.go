package store

import (
	"github.com/boltdb/bolt"
	"github.com/inpime/dbox"
)

var _ dbox.Store = (*BoltDBStore)(nil)
var _ dbox.FileStore = (*BoltDBStore)(nil)

type BoltDBStore struct {
	store dbox.Store
}

func NewBoltDBStore(db *bolt.DB, bucketname string) *BoltDBStore {
	return &BoltDBStore{
		store: dbox.NewBoltDBStore(db, bucketname),
	}
}

func (s BoltDBStore) GetFile(id string, f *File) error {

	return nil
}

func (s BoltDBStore) GetByNameFile(name string, f *File) error {

	return nil
}

func (s BoltDBStore) SaveFile(f *File) error {
	return nil
}

func (s BoltDBStore) DeleteFile(f *File) error {
	return nil
}

// dbox.Store implement interface

func (s BoltDBStore) GetByName(id string, obj dbox.Object) error {
	return s.store.(dbox.FileStore).GetByName(id, obj)
}

func (s BoltDBStore) Get(id string, obj dbox.Object) error {
	return s.store.Get(id, obj)
}

func (s BoltDBStore) Save(obj dbox.Object) error {
	// save to elastic search

	return s.store.Save(obj)
}

func (s BoltDBStore) Delete(obj dbox.Object) error {
	// save to elastic search

	return s.store.Delete(obj)
}

func (s BoltDBStore) Type() dbox.StoreType {
	return s.store.Type()
}
