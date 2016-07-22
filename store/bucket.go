package store

import (
	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/utils/sdata"
)

var (
	MetaFilesBucketMappingKey = "BucketMappingMetaFiles"
	MapFilesBucketMappingKey  = "BucketMappingMapFiles"

	BucketMapDataStoreTypeKey  = "MapDataStoreType"
	BucketRawDataStoreTypeKey  = "RawDataStoreType"
	BucketMetaDataStoreTypeKey = "MetaDataStoreType"

	BucketMapDataStoreNameWithoutPostfixKey  = "MapDataStoreNameWithoutPostfix"
	BucketRawDataStoreNameWithoutPostfixKey  = "RawDataStoreNameWithoutPostfix"
	BucketMetaDataStoreNameWithoutPostfixKey = "MetaDataStoreNameWithoutPostfix"

	BucketFileMappingDefault = sdata.NewStringMap().
					Set("FileID", sdata.NewStringMap().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("Bucket", sdata.NewStringMap().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("Name", sdata.NewStringMap().
						Set("type", "string").
						Set("index", "not_analyzed")).
					Set("TextData", sdata.NewStringMap().
						Set("type", "string")).
					Set("MapData", sdata.NewStringMap().
						Set("type", "object").
						Set("enabled", false)) // disabled mapping by default

		// Set("IsRemoved", fieldMapping{"boolean", "", ""}).
		// Set("CreatedAt", fieldMapping{"date", "", "strict_date_optional_time||epoch_millis"}).
		// Set("UpdatedAt", fieldMapping{"date", "", "strict_date_optional_time||epoch_millis"})

	FileMetaMappingDefault = sdata.NewStringMap().
				Set(dbox.MapDataStoreNameKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.RawDataStoreNameKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.MetaDataFileStoreNameKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.CreatedAtKey, sdata.NewStringMap().
					Set("type", "date").
					Set("format", "strict_date_optional_time||epoch_millis")).
				Set(dbox.UpdatedAtKey, sdata.NewStringMap().
					Set("type", "date").
					Set("format", "strict_date_optional_time||epoch_millis")).
				Set(dbox.MapDataIDMetaKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.RawDataIDMetaKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.BucketKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed")).
				Set(dbox.NameKey, sdata.NewStringMap().
					Set("type", "string").
					Set("index", "not_analyzed"))

	FileEmptyMapDataMapping = sdata.NewStringMap().
				Set("type", "object").
				Set("enabled", false)

	BucketMapMappingDefault = sdata.NewStringMap().
				Set(MetaFilesBucketMappingKey, sdata.NewStringMap().
					Set("type", "object").
					Set("enabled", false)).
				Set(MapFilesBucketMappingKey, sdata.NewStringMap().
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
	f.SetMetaDataStoreType(_type)
	f.SetMetaDataStoreNameWithoutPostfix(true)

	f.SetRawDataStoreName(key)
	f.SetRawDataStoreType(_type)
	f.SetRawDataStoreNameWithoutPostfix(true)

	f.SetMapDataStoreName(key)
	f.SetMapDataStoreType(_type)
	f.SetMapDataStoreNameWithoutPostfix(true)
}

// InitMetaDataStore
func (f Bucket) InitMetaDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "metadata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetMetaDataStoreName(key)
	f.SetMetaDataStoreType(_type)
	f.SetMapDataStoreNameWithoutPostfix(withoutPostfix)
}

// InitRawDataStore
func (f Bucket) InitRawDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "rawdata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetRawDataStoreName(key)
	f.SetRawDataStoreType(_type)
	f.SetMapDataStoreNameWithoutPostfix(withoutPostfix)
}

// InitMapDataStore
func (f Bucket) InitMapDataStore(_type dbox.StoreType, withoutPostfix bool) {
	postfix := "mapdata"
	if withoutPostfix {
		postfix = ""
	}
	key := RegistryStore(_type, f.Name(), postfix)
	f.SetMapDataStoreName(key)
	f.SetMapDataStoreType(_type)
	f.SetMapDataStoreNameWithoutPostfix(withoutPostfix)
}

// without postfix

func (b Bucket) GetMapDataStoreNameWithoutPostfix() bool {
	return sdata.NewStringMapFrom(b.Meta()).Bool(BucketMapDataStoreNameWithoutPostfixKey)
}

func (b Bucket) GetRawDataStoreNameWithoutPostfix() bool {
	return sdata.NewStringMapFrom(b.Meta()).Bool(BucketRawDataStoreNameWithoutPostfixKey)
}

func (b Bucket) GetMetaDataStoreNameWithoutPostfix() bool {
	return sdata.NewStringMapFrom(b.Meta()).Bool(BucketMetaDataStoreNameWithoutPostfixKey)
}

func (b Bucket) SetMapDataStoreNameWithoutPostfix(v bool) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketMapDataStoreNameWithoutPostfixKey, v)
}

func (b Bucket) SetRawDataStoreNameWithoutPostfix(v bool) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketRawDataStoreNameWithoutPostfixKey, v)
}

func (b Bucket) SetMetaDataStoreNameWithoutPostfix(v bool) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketMetaDataStoreNameWithoutPostfixKey, v)
}

// store type

func (b Bucket) SetMapDataStoreType(_type dbox.StoreType) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketMapDataStoreTypeKey, _type)
}

func (b Bucket) SetRawDataStoreType(_type dbox.StoreType) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketRawDataStoreTypeKey, _type)
}

func (b Bucket) SetMetaDataStoreType(_type dbox.StoreType) {
	sdata.NewStringMapFrom(b.Meta()).Set(BucketMetaDataStoreTypeKey, _type)
}

func (b Bucket) GetMapDataStoreType() dbox.StoreType {
	switch v := sdata.NewStringMapFrom(b.Meta()).GetOrNil(BucketMapDataStoreTypeKey).(type) {
	case string:
		return dbox.StoreType(v)
	case dbox.StoreType:
		return v
	default:
		logrus.Warningf("not supported type value %T, want string or dbox.StoreType", v)
	}

	return dbox.StoreType("unknown")
}

func (b Bucket) GetRawDataStoreType() dbox.StoreType {
	switch v := sdata.NewStringMapFrom(b.Meta()).GetOrNil(BucketRawDataStoreTypeKey).(type) {
	case string:
		return dbox.StoreType(v)
	case dbox.StoreType:
		return v
	default:
		logrus.Warningf("not supported type value %T, want string or dbox.StoreType", v)
	}

	return dbox.StoreType("unknown")
}

func (b Bucket) GetMetaDataStoreType() dbox.StoreType {
	switch v := sdata.NewStringMapFrom(b.Meta()).GetOrNil(BucketMetaDataStoreTypeKey).(type) {
	case string:
		return dbox.StoreType(v)
	case dbox.StoreType:
		return v
	default:
		logrus.Warningf("not supported type value %T, want string or dbox.StoreType", v)
	}

	return dbox.StoreType("unknown")
}

//

// MetaDataFilesMapping returns the mapping metadata of file of the search index (for elastic search)
func (f *Bucket) MetaDataFilesMapping() map[string]interface{} {

	return f.MMetaDataFilesMapping().ToMap()
}

// MapDataFilesMapping returns the mapping mapdata of file of the search index (for elastic search)
func (f *Bucket) MapDataFilesMapping() map[string]interface{} {

	return f.MMapDataFilesMapping().ToMap()
}

func (f *Bucket) MMetaDataFilesMapping() *sdata.StringMap {
	return sdata.NewStringMapFrom(f.MapData()).M(MetaFilesBucketMappingKey)
}

func (f *Bucket) MMapDataFilesMapping() *sdata.StringMap {
	return sdata.NewStringMapFrom(f.MapData()).M(MapFilesBucketMappingKey)
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
