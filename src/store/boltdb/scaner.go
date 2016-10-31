package boltdb

import (
	"interfaces"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

func (m *BucketManager) EachBucket(fn func(*interfaces.Bucket) error) error {
	ids := []uuid.UUID{}
	err := m.db.View(func(tx *bolt.Tx) error {
		store := tx.Bucket([]byte(DefaultBucketName_BucketIDs))
		return store.ForEach(func(k, v []byte) error {
			if len(k) != 17 {
				// skip
				return nil
			}

			if k[16] != PrimaryIDsData {
				// skip
				return nil
			}

			bucketID := uuid.FromBytesOrNil(k[:16])
			if uuid.Equal(uuid.Nil, bucketID) {
				m.logger.Printf("[ERR] empty bucket ID, key hex: %x", k)
				// return errors.New("invalid bucket ID")
				return nil
			}

			ids = append(ids, bucketID)

			return nil
		})
	})

	if err != nil {
		return err
	}

	for _, bucketID := range ids {
		bucket, err := m.FindBucket(bucketID, interfaces.FullBucket)
		if err != nil {
			m.logger.Println("[ERR] find bucket ", bucketID, ":", err)
			return err
		}
		if err := fn(bucket); err != nil {
			return err
		}
	}

	return nil
}

func (m *FileManager) EachFile(fn func(*interfaces.File) error) error {
	ids := []uuid.UUID{}

	err := m.db.View(func(tx *bolt.Tx) error {
		store := tx.Bucket([]byte(DefaultBucketName_FileIDs))

		return store.ForEach(func(k, v []byte) error {
			if len(k) != 17 {
				// skip
				return nil
			}

			if k[16] != PrimaryIDsData {
				// skip
				return nil
			}

			fileID := uuid.FromBytesOrNil(k[:16])
			if uuid.Equal(uuid.Nil, fileID) {
				m.logger.Printf("[ERR] empty file ID, key hex: %x", k)
				// return errors.New("invalid bucket ID")
				return nil
			}

			ids = append(ids, fileID)

			return nil
		})
	})

	if err != nil {
		return err
	}

	for _, fileID := range ids {
		file, err := m.FindFile(fileID, interfaces.FullFile)
		if err != nil {
			m.logger.Println("[ERR] find file ", fileID, ":", err)
			return err
		}
		if err := fn(file); err != nil {
			return err
		}
	}

	return nil
}
