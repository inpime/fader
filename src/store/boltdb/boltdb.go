package boltdb

var (
	DefaultBucketName_BucketIDs   = "fader.buckets.from_ids"
	DefaultBucketName_BucektNames = "fader.buckets.from_names"

	DefaultBucketName_FileIDs   = "fader.files.from_ids"
	DefaultBucketName_FileNames = "fader.files.from_names"
)

const (
	_              byte = iota
	PrimaryIDsData byte = iota
	PrimaryNamesData

	ContentTypeData
	OwnersData
	AccessStatusData

	CreatedAtData // required data times
	UpdatedAtData // required data times

	// Data
	LuaScript

	MetaData
	StructuralData
	RawData

	// For buckets
	BucketStoreNames
)
