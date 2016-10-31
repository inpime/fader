package interfaces

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// DataUsed strategy of data
type DataUsed uint32

const (
	// Primary Data
	PrimaryIDsData DataUsed = 1 << (32 - 1 - iota)
	PrimaryNamesData

	ContentTypeData
	OwnersData
	AccessStatusData

	// Data
	LuaScript
	MetaData
	StructuralData
	RawData

	// For buckets
	BucketStoreNames

	// 32-9
)

const (
	DataFile = LuaScript | MetaData | StructuralData | RawData

	FileWithoutRawData = LuaScript | MetaData | StructuralData |
		PrimaryIDsData | PrimaryNamesData | ContentTypeData |
		OwnersData | AccessStatusData

	FullFile = PrimaryIDsData | PrimaryNamesData | ContentTypeData |
		OwnersData | AccessStatusData | DataFile
	FullBucket = PrimaryIDsData | PrimaryNamesData | OwnersData | DataFile |
		BucketStoreNames
)

type FileManager interface {
	FindFileByName(
		bucketName, fileName string,
		used DataUsed,
	) (*File, error)

	FindFile(
		fileID uuid.UUID,
		used DataUsed,
	) (*File, error)

	CreateFile(*File) error
	CreateFileFrom(*File, DataUsed) error
	UpdateFileFrom(*File, DataUsed) error

	DeleteFile(fileID uuid.UUID) error
}

type BucketManager interface {
	FindBucketByName(
		name string,
		used DataUsed,
	) (*Bucket, error)
	FindBucket(
		bucketID uuid.UUID,
		used DataUsed,
	) (*Bucket, error)

	CreateBucket(*Bucket) error
	CreateBucketFrom(*Bucket, DataUsed) error
	UpdateBucket(*Bucket, DataUsed) error
}

type BucketImportManager interface {
	EachBucket(func(*Bucket) error) error
}

type FileImportManager interface {
	EachFile(func(*File) error) error
}

type FileLoader interface {
	File(
		ctx context.Context,
		bucketName, fileName string,
	) (*File, error)
}
