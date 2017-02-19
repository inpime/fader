package boltdb

import (
	"interfaces"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func testFile(bucketID uuid.UUID) *interfaces.File {
	return &interfaces.File{
		FileID:   uuid.NewV4(),
		BucketID: bucketID,
		FileName: "file name",

		LuaScript: []byte("-- lua script"),

		MetaData: map[string]interface{}{
			"a":     "b",
			"arr":   []int64{1, 2, 3},
			"null":  nil,
			"int":   1,
			"float": 0.5,
			"bool":  true,
			"obj": map[string]interface{}{
				"foo": 1,
				"bar": 0.5,
			},
		},
		StructuralData: map[string]interface{}{
			"c":     "d",
			"a":     "b",
			"arr":   []int64{1, 2, 3},
			"null":  nil,
			"int":   1,
			"float": 0.5,
			"bool":  true,
			"obj": map[string]interface{}{
				"foo": 1,
				"bar": 0.5,
			},
		},
		RawData:     []byte("123456"),
		ContentType: "text/plain",
		Owners: []uuid.UUID{
			uuid.NewV4(),
			uuid.NewV4(),
		},
		IsPrivate:  false,
		IsReadOnly: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func TestFileManager_createfulldata_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewFileManager(db)
	//
	expected := testFile(uuid.NewV4())

	err = m.CreateFile(expected)
	assert.NoError(t, err)

	// check

	got, err := m.FindFile(expected.FileID, interfaces.FullFile)
	assert.NoError(t, err)
	assert.Equal(t, expected.FileID, got.FileID)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.FileName, got.FileName)
	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.ContentType, got.ContentType)
	assert.Equal(t, expected.IsPrivate, got.IsPrivate)
	assert.Equal(t, expected.IsReadOnly, got.IsReadOnly)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.Equal(t, expected.UpdatedAt.UnixNano(), got.UpdatedAt.UnixNano())
	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.StructuralData["null"], got.StructuralData["null"])
	assert.EqualValues(t, expected.StructuralData["int"], got.StructuralData["int"])
	assert.EqualValues(t, expected.StructuralData["float"], got.StructuralData["float"])
	assert.Equal(t, expected.StructuralData["bool"], got.StructuralData["bool"])
	// assert.EqualValues(t, expected.StructuralData["obj"], got.StructuralData["obj"])
	// assert.Equal(t, expected.StructuralData["arr"], got.StructuralData["arr"])
	assert.Equal(t, expected.LuaScript, got.LuaScript)
}

func TestFileManager_update_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewFileManager(db)
	//
	expected := testFile(uuid.NewV4())

	err = m.CreateFile(expected)
	assert.NoError(t, err)

	// update
	time.Sleep(time.Microsecond * 1)

	expected.FileName = "new file name"
	expected.MetaData["b"] = "new value"
	err = m.UpdateFileFrom(expected, interfaces.PrimaryNamesData|interfaces.MetaData)
	assert.NoError(t, err)

	// check

	got, err := m.FindFile(expected.FileID, interfaces.FullFile)
	assert.NoError(t, err)
	assert.Equal(t, expected.FileID, got.FileID)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.FileName, got.FileName)
	assert.NotEqual(t, expected.CreatedAt.UnixNano(), got.UpdatedAt.UnixNano())

	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.ContentType, got.ContentType)
	assert.Equal(t, expected.IsPrivate, got.IsPrivate)
	assert.Equal(t, expected.IsReadOnly, got.IsReadOnly)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.StructuralData["null"], got.StructuralData["null"])
	assert.EqualValues(t, expected.StructuralData["int"], got.StructuralData["int"])
	assert.EqualValues(t, expected.StructuralData["float"], got.StructuralData["float"])
	assert.Equal(t, expected.StructuralData["bool"], got.StructuralData["bool"])
	// assert.EqualValues(t, expected.StructuralData["obj"], got.StructuralData["obj"])
	// assert.Equal(t, expected.StructuralData["arr"], got.StructuralData["arr"])
	assert.Equal(t, expected.LuaScript, got.LuaScript)
}

func TestFileManager_findbyname_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewFileManager(db)

	// related bucket
	bucket := testBucket()
	err = m.buckets.CreateBucket(bucket)
	assert.NoError(t, err)

	// file
	expected := testFile(bucket.BucketID)

	err = m.CreateFile(expected)
	assert.NoError(t, err)

	// find
	got, err := m.FindFileByName(
		bucket.BucketName,
		"file name",
		interfaces.FullFile,
	)

	if assert.NoError(t, err) {
		assert.Equal(t, expected.FileID, got.FileID)
		assert.Equal(t, expected.BucketID, got.BucketID)
		assert.Equal(t, expected.FileName, got.FileName)
		assert.Equal(t, expected.RawData, got.RawData)
		assert.Equal(t, expected.ContentType, got.ContentType)
		assert.Equal(t, expected.IsPrivate, got.IsPrivate)
		assert.Equal(t, expected.IsReadOnly, got.IsReadOnly)
		assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
		assert.Equal(t, expected.UpdatedAt.UnixNano(), got.UpdatedAt.UnixNano())
		assert.Equal(t, expected.Owners[0], got.Owners[0])
		assert.Equal(t, expected.Owners[1], got.Owners[1])
		assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
		assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
		assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
		assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
		assert.Equal(t, expected.StructuralData["null"], got.StructuralData["null"])
		assert.EqualValues(t, expected.StructuralData["int"], got.StructuralData["int"])
		assert.EqualValues(t, expected.StructuralData["float"], got.StructuralData["float"])
		assert.Equal(t, expected.StructuralData["bool"], got.StructuralData["bool"])
		// assert.EqualValues(t, expected.StructuralData["obj"], got.StructuralData["obj"])
		// assert.EqualValues(t, expected.StructuralData["arr"], got.StructuralData["arr"])
		assert.Equal(t, expected.LuaScript, got.LuaScript)
	}

	// update
	time.Sleep(time.Microsecond * 1)

	expected.FileName = "new file name"
	expected.MetaData["b"] = "new value"
	err = m.UpdateFileFrom(expected, interfaces.PrimaryNamesData|interfaces.MetaData)
	assert.NoError(t, err)

	// find
	got, err = m.FindFileByName(
		bucket.BucketName,
		"file name",
		interfaces.FullFile,
	)
	assert.EqualError(t, err, interfaces.ErrNotFound.Error())

	// find
	got, err = m.FindFileByName(
		bucket.BucketName,
		"new file name",
		interfaces.FullFile,
	)

	if assert.NoError(t, err) {
		assert.Equal(t, expected.FileID, got.FileID)
		assert.Equal(t, expected.BucketID, got.BucketID)
		assert.Equal(t, expected.FileName, got.FileName)
		assert.NotEqual(t, expected.CreatedAt.UnixNano(), got.UpdatedAt.UnixNano())

		assert.Equal(t, expected.RawData, got.RawData)
		assert.Equal(t, expected.ContentType, got.ContentType)
		assert.Equal(t, expected.IsPrivate, got.IsPrivate)
		assert.Equal(t, expected.IsReadOnly, got.IsReadOnly)
		assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
		assert.Equal(t, expected.Owners[0], got.Owners[0])
		assert.Equal(t, expected.Owners[1], got.Owners[1])
		assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
		assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
		assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
		assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
		assert.Equal(t, expected.StructuralData["null"], got.StructuralData["null"])
		assert.EqualValues(t, expected.StructuralData["int"], got.StructuralData["int"])
		assert.EqualValues(t, expected.StructuralData["float"], got.StructuralData["float"])
		assert.Equal(t, expected.StructuralData["bool"], got.StructuralData["bool"])
		// assert.EqualValues(t, expected.StructuralData["obj"], got.StructuralData["obj"])
		// assert.Equal(t, expected.StructuralData["arr"], got.StructuralData["arr"])
		assert.Equal(t, expected.LuaScript, got.LuaScript)
	}
}

func TestDeleteFile_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewFileManager(db)
	//
	expected := testFile(uuid.NewV4())

	err = m.CreateFile(expected)
	assert.NoError(t, err)

	// check

	err = m.DeleteFile(expected.FileID)
	assert.NoError(t, err)

	//

	_, err = m.FindFile(expected.FileID, interfaces.FullFile)
	assert.EqualError(t, err, interfaces.ErrNotFound.Error())
}
