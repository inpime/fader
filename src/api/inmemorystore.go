package api

import (
	"interfaces"

	"github.com/boltdb/bolt"
)

var InMomryStoreBucketName = []byte("InMemoryStore")

func NewInMemoryStore(db *bolt.DB) *InMemoryStore {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(InMomryStoreBucketName)
		return err
	})

	return &InMemoryStore{
		db: db,
	}
}

type InMemoryStore struct {
	db *bolt.DB
}

func (s *InMemoryStore) Set(
	key string,
	obj interfaces.MsgpackMarshaller,
) error {
	data, err := obj.MarshalMsgpack()
	if err != nil {
		return err
	}

	logger.Println("InMemoryStore: set, length data", len(data))

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InMomryStoreBucketName)
		return b.Put([]byte(key), data)
	})
}

func (s *InMemoryStore) Get(
	key string,
	obj interfaces.MsgpackMarshaller,
) error {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InMomryStoreBucketName)
		data = b.Get([]byte(key))
		return nil
	})

	if len(data) == 0 {
		return nil
	}

	logger.Println("InMemoryStore: get, length data", len(data))

	if err != nil {
		return err
	}

	return obj.UnmarshalMsgpack(data)
}

func (s *InMemoryStore) Del(
	key string,
) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(InMomryStoreBucketName)
		return b.Delete([]byte(key))
	})
}
