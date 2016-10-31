package boltdb

import (
	"errors"
	"interfaces"
	"log"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"os"

	"gopkg.in/vmihailenco/msgpack.v2"
)

var _ interfaces.FileManager = (*FileManager)(nil)

func NewFileManager(db *bolt.DB) *FileManager {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(DefaultBucketName_FileIDs))
		tx.CreateBucketIfNotExists([]byte(DefaultBucketName_FileNames))
		return nil
	})

	return &FileManager{
		db:      db,
		buckets: NewBucketManager(db),
		logger:  log.New(os.Stdout, "[FILE_MANAGER]", -1),
	}
}

type FileManager struct {
	db      *bolt.DB
	buckets interfaces.BucketManager
	logger  *log.Logger
}

func (m *FileManager) CreateFile(file *interfaces.File) error {
	return m.CreateFileFrom(file, interfaces.FullFile)
}

func (m *FileManager) CreateFileFrom(
	file *interfaces.File,
	used interfaces.DataUsed,
) error {
	file.UpdatedAt = time.Now()
	file.CreatedAt = time.Now()

	if uuid.Equal(uuid.Nil, file.FileID) {
		m.logger.Println("[ERR] empty file ID")
		return errors.New("empty file ID")
	}
	if uuid.Equal(uuid.Nil, file.BucketID) {
		m.logger.Println("[ERR] empty bucket ID")
		return errors.New("empty bucket ID")
	}

	err := m.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DefaultBucketName_FileIDs))

		superID := make([]byte, 17)
		copy(superID, file.FileID.Bytes())

		if used&interfaces.PrimaryIDsData != 0 {
			data, err := msgpack.Marshal(
				file.FileID,
				file.BucketID,
			)
			if err != nil {
				m.logger.Println("[ERR] marshal data, PrimaryIDsData", file.FileID)
				return err
			}
			superID[16] = PrimaryIDsData
			if err := bucket.Put(superID, data); err != nil {
				m.logger.Println("[ERR] put data, PrimaryIDsData", superID)
				return err
			}
		}

		if used&interfaces.PrimaryNamesData != 0 {
			data, err := msgpack.Marshal(
				file.FileName,
			)
			if err != nil {
				m.logger.Println("[ERR] marshal data, PrimaryNamesData", file.FileID)
				return err
			}
			superID[16] = PrimaryNamesData
			if err := bucket.Put(superID, data); err != nil {
				m.logger.Println("[ERR] put data, PrimaryIDsData", superID)
				return err
			}
		}

		if used&interfaces.PrimaryIDsData != 0 &&
			used&interfaces.PrimaryNamesData != 0 {
			bucketFromNames := tx.Bucket([]byte(DefaultBucketName_FileNames))

			var relatedInfo = make([]byte, 16*2)
			copy(relatedInfo[:16], file.BucketID[:])
			copy(relatedInfo[16:], file.FileID[:])

			if err := bucketFromNames.
				Put(
					hashFromFile(file.BucketID, file.FileName),
					relatedInfo,
				); err != nil {
				m.logger.Println("[ERR] put referance ID", relatedInfo)
				return err
			}
		}

		if err := putFileData(
			bucket,
			file,
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

func (m *FileManager) UpdateFileFrom(
	file *interfaces.File,
	used interfaces.DataUsed,
) error {
	file.UpdatedAt = time.Now()

	if uuid.Equal(uuid.Nil, file.FileID) {
		m.logger.Println("[ERR] empty file ID")
		return errors.New("empty file ID")
	}

	if used&interfaces.PrimaryNamesData != 0 &&
		uuid.Equal(uuid.Nil, file.BucketID) {
		return errors.New("required bucket ID, because it requires strategy updates")
	}

	previousFile, err := m.FindFile(
		file.FileID,
		interfaces.PrimaryNamesData,
	)

	if err != nil {
		m.logger.Println("[ERR] find previous file,", file.FileID)
		return err
	}

	err = m.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DefaultBucketName_FileIDs))

		superID := make([]byte, 17)
		copy(superID, file.FileID.Bytes())

		if used&interfaces.PrimaryNamesData != 0 {

			// TODO: remove related info by previous name

			data, err := msgpack.Marshal(
				file.FileName,
			)
			if err != nil {
				m.logger.Println("[ERR] marshal data, PrimaryNamesData", file.FileID)
				return err
			}
			superID[16] = PrimaryNamesData
			if err := bucket.Put(superID, data); err != nil {
				m.logger.Println("[ERR] put data, PrimaryIDsData", superID)
				return err
			}

			// update related informations
			bucketFromNames := tx.Bucket([]byte(DefaultBucketName_FileNames))

			// remove previous related info
			if err := bucketFromNames.Delete(SHA1(file.BucketID.String() + previousFile.FileName)); err != nil {
				m.logger.Println("[ERR] remove previous related info, fileID", file.FileID)
				return err
			}

			// update new related info
			var relatedInfo = make([]byte, 16*2)
			copy(relatedInfo[:16], file.BucketID[:])
			copy(relatedInfo[16:], file.FileID[:])

			if err := bucketFromNames.
				Put(
					SHA1(file.BucketID.String()+file.FileName),
					relatedInfo,
				); err != nil {
				m.logger.Println("[ERR] put referance ID", relatedInfo)
				return err
			}
		}

		if err := putFileData(
			bucket,
			file,
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

func (m *FileManager) FindFileByName(
	bucketName,
	fileName string,
	used interfaces.DataUsed,
) (*interfaces.File, error) {
	// 1. find bucket by name
	bucket, err := m.buckets.FindBucketByName(bucketName, interfaces.PrimaryIDsData)
	if err != nil {
		m.logger.Println("[ERR] find bucket by name, ", bucketName)
		return nil, err
	}

	// 2. get hash from bucketID+fileName
	relatedID := hashFromFile(
		bucket.BucketID,
		fileName,
	)

	// 3. get relatedInfo
	var bucketID, fileID uuid.UUID

	err = m.db.View(func(tx *bolt.Tx) error {
		bucketName := tx.Bucket([]byte(DefaultBucketName_FileNames))

		var relatedInfo = make([]byte, 16*2)
		copy(relatedInfo, bucketName.Get(relatedID))
		copy(bucketID[:], relatedInfo[:16])
		copy(fileID[:], relatedInfo[16:])

		return nil
	})

	if uuid.Equal(uuid.Nil, fileID) {
		m.logger.Println("[DEBUG] not found file,", bucketName, fileName)
		return nil, interfaces.ErrNotFound
	}

	if err != nil {
		m.logger.Printf("[ERR] find file by name %q, %q, %q\n", bucketName, fileName, err)
		return nil, interfaces.ErrInternal
	}

	// 4. get file by fileID
	return m.FindFile(fileID, used)
}

func (m *FileManager) FindFile(
	fileID uuid.UUID,
	used interfaces.DataUsed,
) (*interfaces.File, error) {
	if uuid.Equal(uuid.Nil, fileID) {
		m.logger.Println("[ERR] empty file ID")
		return nil, errors.New("empty file ID")
	}

	file := interfaces.NewFile()
	superID := make([]byte, 17)
	copy(superID, fileID.Bytes())

	err := m.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DefaultBucketName_FileIDs))

		// Primary IDs
		superID[16] = PrimaryIDsData
		data := bucket.Get(superID)
		if len(data) == 0 {
			return interfaces.ErrNotFound
		}
		err := msgpack.Unmarshal(
			data,
			&file.FileID,
			&file.BucketID,
		)

		if err != nil {
			m.logger.Println("[ERR] unmarshal PrimaryIDsData data,", err)
			return err
		}

		// Primary Names
		superID[16] = PrimaryNamesData
		err = msgpack.Unmarshal(
			bucket.Get(superID),
			&file.FileName,
		)

		if err != nil {
			m.logger.Println("[ERR] unmarshal PrimaryNamesData data,", err)
			return err
		}

		if err := getFileData(
			bucket,
			file,
			used,
			m.logger,
		); err != nil {
			return err
		}

		return nil
	})

	return file, err
}

func (m *FileManager) DeleteFile(fileID uuid.UUID) error {
	superID := make([]byte, 17)
	copy(superID, fileID.Bytes())

	if uuid.Equal(uuid.Nil, fileID) {
		m.logger.Println("[ERR] delete, empty file ID")
		return errors.New("empty file ID")
	}

	file, err := m.FindFile(fileID, interfaces.PrimaryIDsData|interfaces.PrimaryNamesData)

	if err != nil {
		return err
	}

	err = m.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(DefaultBucketName_FileIDs))

		flags := []byte{
			PrimaryIDsData,
			PrimaryNamesData,
			CreatedAtData,
			UpdatedAtData,
			ContentTypeData,
			OwnersData,
			AccessStatusData,
			MetaData,
			StructuralData,
			RawData,
		}
		for _, flag := range flags {
			superID[16] = flag
			if err := bucket.Delete(superID); err != nil {
				m.logger.Println("[ERR] delete, from flag '%v'", superID, ":", err)
				return err
			}
		}

		bucket = tx.Bucket([]byte(DefaultBucketName_FileNames))
		if err := bucket.Delete(hashFromFile(file.BucketID, file.FileName)); err != nil {
			m.logger.Println("[ERR] delete, related info:", err)
			return err
		}

		return nil
	})

	return err
}

// internal

func getFileData(
	bucket *bolt.Bucket,
	file *interfaces.File,
	used interfaces.DataUsed,
	logger *log.Logger,
) error {
	superID := make([]byte, 17)
	copy(superID, file.FileID.Bytes())

	{
		superID[16] = CreatedAtData
		err := file.CreatedAt.UnmarshalBinary(bucket.Get(superID))

		if err != nil {
			logger.Println("[ERR] unmarshal data, CreatedAtData", file.FileID)
			return err
		}
	}

	{
		superID[16] = UpdatedAtData
		err := file.UpdatedAt.UnmarshalBinary(bucket.Get(superID))

		if err != nil {
			logger.Println("[ERR] unmarshal data, UpdatedAtData", file.FileID)
			return err
		}
	}

	// used data

	if used&interfaces.ContentTypeData != 0 {
		superID[16] = ContentTypeData

		err := msgpack.Unmarshal(
			bucket.Get(superID),
			&file.ContentType,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, ContentTypeData", file.FileID)
			return err
		}
	}

	if used&interfaces.OwnersData != 0 {
		superID[16] = OwnersData

		err := msgpack.Unmarshal(
			bucket.Get(superID),
			&file.Owners,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, OwnersData", file.FileID)
			return err
		}
	}

	if used&interfaces.AccessStatusData != 0 {
		superID[16] = AccessStatusData

		err := msgpack.Unmarshal(
			bucket.Get(superID),
			&file.IsPrivate,
			&file.IsReadOnly,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, AccessStatusData", file.FileID)
			return err
		}
	}

	if used&interfaces.LuaScript != 0 {
		superID[16] = LuaScript
		file.LuaScript = bucket.Get(superID)
	}

	if used&interfaces.MetaData != 0 {
		superID[16] = MetaData

		err := msgpack.Unmarshal(
			bucket.Get(superID),
			&file.MetaData,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, MetaData", file.FileID)
			return err
		}
	}

	if used&interfaces.StructuralData != 0 {
		superID[16] = StructuralData

		err := msgpack.Unmarshal(
			bucket.Get(superID),
			&file.StructuralData,
		)

		if err != nil {
			logger.Println("[ERR] unmarshal data, StructuralData", file.FileID)
			return err
		}
	}

	if used&interfaces.RawData != 0 {

		superID[16] = RawData
		file.RawData = bucket.Get(superID)
	}

	return nil
}

func putFileData(
	bucket *bolt.Bucket,
	file *interfaces.File,
	used interfaces.DataUsed,
	logger *log.Logger,
	isCreated bool,
) error {
	superID := make([]byte, 17)
	copy(superID, file.FileID.Bytes())

	if isCreated {
		data, err := file.CreatedAt.MarshalBinary()
		if err != nil {
			logger.Println("[ERR] marshal data, CreatedAtData", file.FileID)
			return err
		}
		superID[16] = CreatedAtData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, CreatedAtData", superID)
			return err
		}
	}

	{
		data, err := file.UpdatedAt.MarshalBinary()
		if err != nil {
			logger.Println("[ERR] marshal data, UpdatedAtData", file.FileID)
			return err
		}
		superID[16] = UpdatedAtData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, UpdatedAtData", superID)
			return err
		}
	}

	// used  data

	if used&interfaces.ContentTypeData != 0 {
		data, err := msgpack.Marshal(file.ContentType)
		if err != nil {
			logger.Println("[ERR] marshal data, ContentTypeData", file.FileID)
			return err
		}
		superID[16] = ContentTypeData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, ContentTypeData", superID)
			return err
		}
	}

	if used&interfaces.OwnersData != 0 {
		data, err := msgpack.Marshal(file.Owners)
		if err != nil {
			logger.Println("[ERR] marshal data, OwnersData", file.FileID)
			return err
		}
		superID[16] = OwnersData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, OwnersData", superID)
			return err
		}
	}

	if used&interfaces.AccessStatusData != 0 {
		data, err := msgpack.Marshal(
			file.IsPrivate,
			file.IsReadOnly,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, AccessStatusData", file.FileID)
			return err
		}
		superID[16] = AccessStatusData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, AccessStatusData", superID)
			return err
		}
	}

	if used&interfaces.LuaScript != 0 {
		superID[16] = LuaScript
		if err := bucket.Put(superID, file.LuaScript); err != nil {
			logger.Println("[ERR] put data, LuaScript", superID)
			return err
		}
	}

	if used&interfaces.MetaData != 0 {
		data, err := msgpack.Marshal(
			file.MetaData,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, MetaData", file.FileID)
			return err
		}
		superID[16] = MetaData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, MetaData", superID)
			return err
		}
	}

	if used&interfaces.StructuralData != 0 {
		data, err := msgpack.Marshal(
			file.StructuralData,
		)
		if err != nil {
			logger.Println("[ERR] marshal data, StructuralData", file.FileID)
			return err
		}
		superID[16] = StructuralData
		if err := bucket.Put(superID, data); err != nil {
			logger.Println("[ERR] put data, StructuralData", superID)
			return err
		}
	}

	if used&interfaces.RawData != 0 {
		superID[16] = RawData
		if err := bucket.Put(superID, file.RawData); err != nil {
			logger.Println("[ERR] put data, RawData", superID)
			return err
		}
	}

	return nil
}
