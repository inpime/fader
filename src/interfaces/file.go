package interfaces

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func NewFile() *File {
	return &File{
		MetaData:       make(map[string]interface{}),
		StructuralData: make(map[string]interface{}),
	}
}

type File struct {
	FileID   uuid.UUID
	BucketID uuid.UUID

	FileName string

	LuaScript []byte

	MetaData       map[string]interface{}
	StructuralData map[string]interface{}
	RawData        []byte

	ContentType string
	Owners      []uuid.UUID
	IsPrivate   bool
	IsReadOnly  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBucket() *Bucket {
	return &Bucket{
		MetaData:       make(map[string]interface{}),
		StructuralData: make(map[string]interface{}),
	}
}

type Bucket struct {
	BucketID   uuid.UUID
	BucketName string

	Owners []uuid.UUID

	LuaScript []byte

	MetaData       map[string]interface{}
	StructuralData map[string]interface{}
	RawData        []byte

	MetaDataStoreName       string
	StructuralDataStoreName string
	DataStoreName           string

	CreatedAt time.Time
	UpdatedAt time.Time
}
