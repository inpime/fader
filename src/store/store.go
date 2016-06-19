package store

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/inpime/dbox"
)

var (
	BoltDBClient  *bolt.DB
	WorkspacePath string
)

// FormatStoreName formats the storage key
func FormatStoreName(_type dbox.StoreType, bucketname string, args ...string) string {
	key := string(_type) + "." + bucketname

	if len(args) > 0 && len(args[0]) > 0 {
		key += "." + args[0]
	}

	return key
}

// RegistryStore registration on a key storage
func RegistryStore(_type dbox.StoreType, bucketname string, args ...string) string {
	key := FormatStoreName(_type, bucketname, args...)

	switch _type {
	case dbox.LocalStoreType:
		dbox.RegistryStore(key, NewLocalStore(nil, WorkspacePath+"/."+bucketname+"/"))
	case dbox.BoltDBStoreType:
		dbox.RegistryStore(key, NewBoltDBStore(BoltDBClient, bucketname))
	case dbox.MemoryStoreType:
		dbox.RegistryStore(key, dbox.NewMemoryStore())
	default:
		logrus.Errorf("registry store: not supported type %q", _type)
	}

	return key
}

// BucketStore bucket storage
// NOTICE: do not forget to set dbox.BucketStore

// BucketByName return bucket from name
// if not exist file, file accepts values nil, err accepts values ErrNotFound
func BucketByName(name string) (file *Bucket, err error) {

	file = NewBucket()
	err = dbox.BucketStore.(dbox.FileStore).GetByName(name, file)

	if err == dbox.ErrNotFound || err == dbox.ErrNotFoundBucket {
		file.SetName(name)
	}

	return
}

// LoadOrNewFileIDViaStores загрузка файла через указанные store name
// Важно соблюдать порядок: metadata file, mapdata file, rawdata file
// Если указан только один store - он применяется для mapdata, rawdata
// bucketname используется как есть. Это значит что он должен быть валидным
// Данную функцию следует использовать в случае множественной выборки
// файлов из одного и того же бакета
// WARNING: Важно указать верные store для файла что бы не было ошибок
// TODO: после реализации кеша для функций LoadOrNewFile и LoadOrNewFileID
// данная функция теряет смысл
func LoadOrNewFileIDViaStores(bucketName, fileId string, storeNames ...string) (*File, error) {
	if len(storeNames) == 0 {
		return nil, dbox.ErrInvalidData
	}

	var mapdataStore,
		rawdataStore,
		metadatStore = dbox.MustStore(storeNames[0]),
		dbox.MustStore(storeNames[0]),
		dbox.MustStore(storeNames[0])

	if len(storeNames) == 3 {
		mapdataStore = dbox.MustStore(storeNames[1])
		rawdataStore = dbox.MustStore(storeNames[2])
	}

	file, err := dbox.NewFileID(fileId, metadatStore)
	file.SetMapDataStore(mapdataStore)
	file.SetRawDataStore(rawdataStore)

	if err == dbox.ErrNotFound {
		file.SetID("")
		file.SetBucket(bucketName)
	}

	return MustFile(file), nil
}

// LoadOrNewFile return file by bucket name and filename.
// If not exist bucket, file accepts values nil, err accepts values ErrNotFoundBucket.
// If not exist file, file accepts values nil, err accepts values ErrNotFound.
// In cases where file not exist, file name and file bucket name accepts values from arguments.
func LoadOrNewFile(bucketName string, fileName string) (*File, error) {
	bucket, err := BucketByName(bucketName)
	if err != nil {

		if err == dbox.ErrNotFound {
			return nil, dbox.ErrNotFoundBucket
		}

		return nil, err
	}

	file, err := dbox.NewFileName(fileName, dbox.MustStore(bucket.MetaDataStoreName()))
	file.SetMapDataStore(dbox.MustStore(bucket.MapDataStoreName()))
	file.SetRawDataStore(dbox.MustStore(bucket.RawDataStoreName()))

	if err == dbox.ErrNotFound {
		file.SetName(fileName)
		file.SetBucket(bucketName)
	}

	return &File{File: file}, err
}

// LoadOrNewFileID return file by bucket name and file id.
// If not exist bucket, file accepts values nil, err accepts values ErrNotFoundBucket.
// If not exist file, file accepts values nil, err accepts values ErrNotFound.
// In cases where file not exist, file id and file bucket name accepts values from arguments.
func LoadOrNewFileID(bucketName string, fileId string) (*File, error) {
	bucket, err := BucketByName(bucketName)
	if err != nil {

		if err == dbox.ErrNotFound {
			return nil, dbox.ErrNotFoundBucket
		}

		return nil, err
	}

	file, err := dbox.NewFileID(fileId, dbox.MustStore(bucket.MetaDataStoreName()))
	file.SetMapDataStore(dbox.MustStore(bucket.MapDataStoreName()))
	file.SetRawDataStore(dbox.MustStore(bucket.RawDataStoreName()))

	if err == dbox.ErrNotFound {
		file.SetID("")
		file.SetBucket(bucketName)
	}

	return &File{File: file}, nil
}
