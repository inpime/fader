package boltdb

import (
	"errors"
	"interfaces"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func NewBucketManager(db *bolt.DB) *BucketManager {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(DefaultBucketName_BucketIDs))
		tx.CreateBucketIfNotExists([]byte(DefaultBucketName_BucektNames))
		return nil
	})

	return &BucketManager{
		db:     db,
		logger: log.New(os.Stdout, "[BUCKET_MANAGER]", -1),
	}
}

type BucketManager struct {
	db     *bolt.DB
	logger *log.Logger
}

func (m *BucketManager) CreateBucket(bucket *interfaces.Bucket) error {
	return m.CreateBucketFrom(bucket, interfaces.FullBucket)
}
func (m *BucketManager) CreateBucketFrom(
	bucket *interfaces.Bucket,
	used interfaces.DataUsed,
) error {
	bucket.UpdatedAt = time.Now()
	bucket.CreatedAt = time.Now()

	if uuid.Equal(uuid.Nil, bucket.BucketID) {
		m.logger.Println("[ERR] empty bucket ID")
		return errors.New("invalid data")
	}

	err := m.db.Update(func(tx *bolt.Tx) error {
		_bucket := tx.Bucket([]byte(DefaultBucketName_BucketIDs))

		superID := make([]byte, 17)
		copy(superID, bucket.BucketID.Bytes())

		if used&interfaces.PrimaryIDsData != 0 {
			superID[16] = PrimaryIDsData
			if err := _bucket.Put(superID, bucket.BucketID.Bytes()); err != nil {
				m.logger.Println("[ERR] put data, PrimaryIDsData", superID)
				return err
			}
		}

		if used&interfaces.PrimaryNamesData != 0 {
			data, err := msgpack.Marshal(
				bucket.BucketName,
			)
			if err != nil {
				m.logger.Println("[ERR] marshal data, PrimaryNamesData", bucket.BucketID)
				return err
			}
			superID[16] = PrimaryNamesData
			if err := _bucket.Put(superID, data); err != nil {
				m.logger.Println("[ERR] put data, PrimaryIDsData", superID)
				return err
			}
		}

		if used&interfaces.PrimaryIDsData != 0 &&
			used&interfaces.PrimaryNamesData != 0 {
			bucketFromNames := tx.Bucket([]byte(DefaultBucketName_BucektNames))

			if err := bucketFromNames.
				Put(
					SHA1(bucket.BucketName),
					bucket.BucketID.Bytes(),
				); err != nil {
				m.logger.Println("[ERR] put referance ID", bucket.BucketID)
				return err
			}
		}

		if err := putBucketData(
			_bucket,
			bucket,
			used,
			m.logger,
			true,
		); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (m *BucketManager) UpdateBucket(
	bucket *interfaces.Bucket,
	used interfaces.DataUsed,
) error {
	bucket.UpdatedAt = time.Now()

	if uuid.Equal(uuid.Nil, bucket.BucketID) {
		m.logger.Println("[ERR] empty bucket ID")
		return errors.New("invalid data")
	}

	previousBucket, err := m.FindBucket(
		bucket.BucketID,
		interfaces.PrimaryNamesData,
	)

	if err != nil {
		m.logger.Println("[ERR] find previous bucket,", bucket.BucketID)
		return err
	}

	err = m.db.Update(func(tx *bolt.Tx) error {
		_bucket := tx.Bucket([]byte(DefaultBucketName_BucketIDs))

		superID := make([]byte, 17)
		copy(superID, bucket.BucketID.Bytes())

		if used&interfaces.PrimaryNamesData != 0 {

			data, err := msgpack.Marshal(
				bucket.BucketName,
			)
			if err != nil {
				m.logger.Println("[ERR] marshal data, PrimaryNamesData", bucket.BucketID)
				return err
			}
			superID[16] = PrimaryNamesData
			if err := _bucket.Put(superID, data); err != nil {
				m.logger.Println("[ERR] put data, PrimaryNamesData", superID)
				return err
			}

			// update related informations
			bucketFromNames := tx.Bucket([]byte(DefaultBucketName_BucektNames))

			// remove previous related info
			if err := bucketFromNames.Delete(SHA1(previousBucket.BucketName)); err != nil {
				m.logger.Println("[ERR] remove previous related info, bucketID", bucket.BucketID)
				return err
			}

			// update new related info
			if err := bucketFromNames.
				Put(
					SHA1(bucket.BucketName),
					bucket.BucketID.Bytes(),
				); err != nil {
				m.logger.Println("[ERR] put related info", bucket.BucketID)
				return err
			}
		}

		if err := putBucketData(
			_bucket,
			bucket,
			used,
			m.logger,
			false,
		); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (m *BucketManager) FindBucketByName(
	name string,
	used interfaces.DataUsed,
) (*interfaces.Bucket, error) {

	// 3. get relatedInfo
	var bucketID uuid.UUID

	err := m.db.View(func(tx *bolt.Tx) error {
		bucketName := tx.Bucket([]byte(DefaultBucketName_BucektNames))

		return bucketID.UnmarshalBinary(bucketName.Get(SHA1(name)))
	})

	if uuid.Equal(uuid.Nil, bucketID) {
		m.logger.Println("[DEBUG] not found bucket,", name)
		return nil, interfaces.ErrNotFound
	}

	if err != nil {
		m.logger.Printf("[ERR] find bucket by name %q, %q\n", name, err)
		return nil, interfaces.ErrInternal
	}

	// 4. get bucket by ID
	return m.FindBucket(bucketID, used)
}
func (m *BucketManager) FindBucket(
	bucketID uuid.UUID,
	used interfaces.DataUsed,
) (*interfaces.Bucket, error) {
	if uuid.Equal(uuid.Nil, bucketID) {
		m.logger.Println("[ERR] empty bucket ID")
		return nil, errors.New("invalid data")
	}

	bucket := interfaces.NewBucket()

	err := m.db.View(func(tx *bolt.Tx) error {
		_bucket := tx.Bucket([]byte(DefaultBucketName_BucketIDs))

		superID := make([]byte, 17)
		copy(superID, bucketID.Bytes())

		// Primary IDs
		superID[16] = PrimaryIDsData
		data := _bucket.Get(superID)
		if len(data) == 0 {
			return interfaces.ErrNotFound
		}
		err := bucket.BucketID.UnmarshalBinary(data)

		if err != nil {
			return err
		}

		// Primary Names
		superID[16] = PrimaryNamesData
		err = msgpack.Unmarshal(
			_bucket.Get(superID),
			&bucket.BucketName,
		)

		if err != nil {
			return err
		}

		if err := getBucketData(
			_bucket,
			bucket,
			used,
			m.logger,
		); err != nil {
			return err
		}

		return nil
	})

	return bucket, err
}

// internale

func getBucketData(
	_bucket *bolt.Bucket,
	bucket *interfaces.Bucket,
	used interfaces.DataUsed,
	logger *log.Logger,
) error {
	superID := make([]byte, 17)
	copy(superID, bucket.BucketID.Bytes())

	{
		superID[16] = CreatedAtData
		err := bucket.CreatedAt.UnmarshalBinary(_bucket.Get(superID))

		if err != nil {
			logger.Println("[ERR] unmarshal data, CreatedAtData", bucket.BucketID)
			return err
		}
	}

	{
		superID[16] = UpdatedAtData
		err := bucket.UpdatedAt.UnmarshalBinary(_bucket.Get(superID))

		if err != nil {
			logger.Println("[ERR] unmarshal data, UpdatedAtData", bucket.BucketID)
			return err
		}
	}

	// used data

	if used&interfaces.OwnersData != 0 {
		superID[16] = OwnersData

		err := msgpack.Unmarshal(
			_bucket.Get(superID),
			&bucket.Owners,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, OwnersData", bucket.BucketID)
			return err
		}
	}

	if used&interfaces.LuaScript != 0 {
		superID[16] = LuaScript

		bucket.LuaScript = _bucket.Get(superID)
	}

	if used&interfaces.MetaData != 0 {
		superID[16] = MetaData

		err := msgpack.Unmarshal(
			_bucket.Get(superID),
			&bucket.MetaData,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, MetaData", bucket.BucketID)
			return err
		}
	}

	if used&interfaces.StructuralData != 0 {
		superID[16] = StructuralData

		err := msgpack.Unmarshal(
			_bucket.Get(superID),
			&bucket.StructuralData,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, StructuralData", bucket.BucketID)
			return err
		}
	}

	if used&interfaces.RawData != 0 {
		superID[16] = RawData
		bucket.RawData = _bucket.Get(superID)
	}

	if used&interfaces.BucketStoreNames != 0 {

		superID[16] = BucketStoreNames

		err := msgpack.Unmarshal(
			_bucket.Get(superID),
			&bucket.MetaDataStoreName,
			&bucket.StructuralDataStoreName,
			&bucket.DataStoreName,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, BucketStoreNames", bucket.BucketID)
			return err
		}
	}

	return nil
}

func putBucketData(
	_bucket *bolt.Bucket,
	bucket *interfaces.Bucket,
	used interfaces.DataUsed,
	logger *log.Logger,
	isCreated bool,
) error {
	superID := make([]byte, 17)
	copy(superID, bucket.BucketID.Bytes())

	if isCreated {
		data, err := bucket.CreatedAt.MarshalBinary()
		if err != nil {
			logger.Println("[ERR] marshal data, CreatedAtData", bucket.BucketID)
			return err
		}
		superID[16] = CreatedAtData
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, CreatedAtData", superID)
			return err
		}
	}

	{
		data, err := bucket.UpdatedAt.MarshalBinary()
		if err != nil {
			logger.Println("[ERR] marshal data, UpdatedAtData", bucket.BucketID)
			return err
		}
		superID[16] = UpdatedAtData
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, UpdatedAtData", superID)
			return err
		}
	}

	// used  data

	if used&interfaces.OwnersData != 0 {
		data, err := msgpack.Marshal(bucket.Owners)
		if err != nil {
			logger.Println("[ERR] marshal data, OwnersData", bucket.BucketID)
			return err
		}
		superID[16] = OwnersData
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, OwnersData", superID)
			return err
		}
	}

	if used&interfaces.LuaScript != 0 {
		superID[16] = LuaScript
		if err := _bucket.Put(superID, bucket.LuaScript); err != nil {
			logger.Println("[ERR] put data, LuaScript", superID)
			return err
		}
	}

	if used&interfaces.MetaData != 0 {
		data, err := msgpack.Marshal(
			bucket.MetaData,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, MetaData", bucket.BucketID)
			return err
		}
		superID[16] = MetaData
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, MetaData", superID)
			return err
		}
	}

	if used&interfaces.StructuralData != 0 {
		data, err := msgpack.Marshal(
			bucket.StructuralData,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, StructuralData", bucket.BucketID)
			return err
		}
		superID[16] = StructuralData
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, StructuralData", superID)
			return err
		}
	}

	if used&interfaces.RawData != 0 {
		superID[16] = RawData
		if err := _bucket.Put(superID, bucket.RawData); err != nil {
			logger.Println("[ERR] put data, RawData", superID)
			return err
		}
	}

	if used&interfaces.BucketStoreNames != 0 {
		data, err := msgpack.Marshal(
			bucket.MetaDataStoreName,
			bucket.StructuralDataStoreName,
			bucket.DataStoreName,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, BucketStoreNames", bucket.BucketID)
			return err
		}
		superID[16] = BucketStoreNames
		if err := _bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, BucketStoreNames", superID)
			return err
		}
	}

	return nil
}
