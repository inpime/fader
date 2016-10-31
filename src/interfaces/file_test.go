package interfaces

import (
	"bytes"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestFileSerializer_simple(t *testing.T) {
	file := &File{
		FileID:   uuid.NewV4(),
		BucketID: uuid.NewV4(),

		FileName: "filename",

		ContentType: "text/plain",
		Owners: []uuid.UUID{
			uuid.NewV4(),
		},
		IsPrivate:  true,
		IsReadOnly: false,

		MetaData: map[string]interface{}{
			"a": "b1",
		},
		StructuralData: map[string]interface{}{
			"a": "b2",
		},
		RawData: []byte("data data"),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	network, err := file.MarshalMsgpack()
	assert.NoError(t, err)
	assert.True(t, len(network) > 0)

	// unmarshal

	gotFile := &File{}
	err = gotFile.UnmarshalMsgpack(network)
	assert.NoError(t, err)
	assert.True(t, uuid.Equal(file.FileID, gotFile.FileID))
	assert.True(t, uuid.Equal(file.BucketID, gotFile.BucketID))
	assert.True(t, uuid.Equal(file.Owners[0], gotFile.Owners[0]))
	assert.EqualValues(t, file.FileName, gotFile.FileName)
	assert.EqualValues(t, file.ContentType, gotFile.ContentType)
	assert.EqualValues(t, file.IsPrivate, gotFile.IsPrivate)
	assert.EqualValues(t, file.IsReadOnly, gotFile.IsReadOnly)
	assert.EqualValues(t, file.IsReadOnly, gotFile.IsReadOnly)
	assert.EqualValues(t, file.MetaData["a"], gotFile.MetaData["a"])
	assert.EqualValues(t, file.StructuralData["a"], gotFile.StructuralData["a"])
	assert.True(t, bytes.Equal(file.RawData, gotFile.RawData))
}
