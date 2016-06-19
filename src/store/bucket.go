package store

import (
	"github.com/inpime/dbox"
	"utils"
)

var (
	MetaFilesBucketMappingKey = "BucketMappingMetaFiles"
	MapFilesBucketMappingKey  = "BucketMappingMapFiles"

	BucketFileMappingDefault = utils.Map().
					Set("FileID", utils.Map().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("Bucket", utils.Map().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("Name", utils.Map().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("TextData", utils.Map().
						Set("type", "string")).
					Set("MapData", utils.Map().
						Set("type", "object").
						Set("enabled", false)) // disabled mapping by default

		// Set("IsRemoved", fieldMapping{"boolean", "", ""}).
		// Set("CreatedAt", fieldMapping{"date", "", "strict_date_optional_time||epoch_millis"}).
		// Set("UpdatedAt", fieldMapping{"date", "", "strict_date_optional_time||epoch_millis"})

	FileMetaMappingDefault = utils.Map().
				Set(dbox.MapDataStoreNameKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.RawDataStoreNameKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.MetaDataFileStoreNameKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.CreatedAtKey, utils.Map().
					Set("type", "date").
					Set("format", "strict_date_optional_time||epoch_millis")).
				Set(dbox.UpdatedAtKey, utils.Map().
					Set("type", "date").
					Set("format", "strict_date_optional_time||epoch_millis")).
				Set(dbox.MapDataIDMetaKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.RawDataIDMetaKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.BucketKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.NameKey, utils.Map().
					Set("type", "string").
					Set("index", "not_analyzed"))

	FileEmptyMapDataMapping = utils.Map().
				Set("type", "object").
				Set("enabled", false)

	BucketMapMappingDefault = utils.Map().
				Set(MetaFilesBucketMappingKey, utils.Map().
					Set("type", "object").
					Set("enabled", false)).
				Set(MapFilesBucketMappingKey, utils.Map().
					Set("type", "object").
					Set("enabled", false))
)

var _ dbox.Object = (*Bucket)(nil)

type Bucket struct {
	*dbox.Bucket
}

func NewBucket() *Bucket {
	bucket := &Bucket{
		Bucket: dbox.NewBucket(),
	}

	bucket.SetBucket("buckets")

	return bucket
}

// StoreKeyName returns the storage key
func (f Bucket) StoreKeyName(_type dbox.StoreType, args ...string) string {
	return FormatStoreName(_type, f.Name(), args...)
}

// InitInOneStore all storage bucket in one store
func (f Bucket) InitInOneStore(_type dbox.StoreType) {
	key := RegistryStore(_type, f.Name())
	f.SetMetaDataStoreName(key)
	f.SetRawDataStoreName(key)
	f.SetMapDataStoreName(key)
}

// InitMetaDataStore
func (f Bucket) InitMetaDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "mapdata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetMetaDataStoreName(key)
}

// InitRawDataStore
func (f Bucket) InitRawDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "mapdata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetRawDataStoreName(key)
}

// InitMapDataStore
func (f Bucket) InitMapDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "mapdata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetMapDataStoreName(key)
}

// MetaDataFilesMapping returns the mapping metadata of file of the search index (for elastic search)
func (f Bucket) MetaDataFilesMapping() map[string]interface{} {

	return utils.M(f.MapData()).Map(MetaFilesBucketMappingKey)
}

// MapDataFilesMapping returns the mapping mapdata of file of the search index (for elastic search)
func (f Bucket) MapDataFilesMapping() map[string]interface{} {

	return utils.M(f.MapData()).Map(MapFilesBucketMappingKey)
}

func (f *Bucket) MMetaDataFilesMapping() utils.M {
	return utils.M(f.MapData()).M(MetaFilesBucketMappingKey)
}

func (f *Bucket) MMapDataFilesMapping() utils.M {
	return utils.M(f.MapData()).M(MapFilesBucketMappingKey)
}

func (f *Bucket) UpdateMapping() error {
	return UpdateSearchMapping(f)
}

func (f *Bucket) Sync() error {

	if err := f.Bucket.Sync(); err != nil {
		return err
	}

	return UpdateSearchDocument("buckets", f.ID(), FileSearchFromFile(&File{&f.Bucket.File}))
}
