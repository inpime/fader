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

func testBucket() *interfaces.Bucket {
	return &interfaces.Bucket{
		BucketID:   uuid.NewV4(),
		BucketName: "bucket name",

		LuaScript: []byte("-- script"),

		MetaData: map[string]interface{}{
			"a": "b",
		},
		StructuralData: map[string]interface{}{
			"c": "d",
		},
		RawData: []byte("123456"),
		Owners: []uuid.UUID{
			uuid.NewV4(),
			uuid.NewV4(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		MetaDataStoreName:       "metadatastorname",
		StructuralDataStoreName: "structuraldatastorename",
		DataStoreName:           "rawdatastorename",
	}
}

func TestBucketManager_createfulldata_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewBucketManager(db)

	expected := testBucket()

	err = m.CreateBucket(expected)
	assert.NoError(t, err)

	// check

	got, err := m.FindBucket(expected.BucketID, interfaces.FullBucket)
	assert.NoError(t, err)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.BucketName, got.BucketName)
	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.Equal(t, expected.UpdatedAt.UnixNano(), got.UpdatedAt.UnixNano())
	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.MetaDataStoreName, got.MetaDataStoreName)
	assert.Equal(t, expected.StructuralDataStoreName, got.StructuralDataStoreName)
	assert.Equal(t, expected.DataStoreName, got.DataStoreName)
	assert.Equal(t, expected.LuaScript, got.LuaScript)
}

func TestBucketManager_update_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewBucketManager(db)

	expected := testBucket()

	err = m.CreateBucket(expected)
	assert.NoError(t, err)

	// update
	time.Sleep(time.Microsecond * 1)

	expected.BucketName = "new bucket name"
	expected.MetaData["b"] = "new value"
	err = m.UpdateBucket(expected, interfaces.PrimaryNamesData|interfaces.MetaData)
	assert.NoError(t, err)

	// check

	got, err := m.FindBucket(expected.BucketID, interfaces.FullBucket)
	assert.NoError(t, err)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.BucketName, got.BucketName)
	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.NotEqual(t, expected.CreatedAt.UnixNano(), got.UpdatedAt.UnixNano())

	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.MetaDataStoreName, got.MetaDataStoreName)
	assert.Equal(t, expected.StructuralDataStoreName, got.StructuralDataStoreName)
	assert.Equal(t, expected.DataStoreName, got.DataStoreName)
	assert.Equal(t, expected.LuaScript, got.LuaScript)
}

func TestBucketManager_findbyname_simple(t *testing.T) {
	db, err := bolt.Open("./_testdb.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	defer func() {
		os.RemoveAll("./_testdb.db")
	}()

	assert.NoError(t, err)
	m := NewBucketManager(db)

	expected := testBucket()

	err = m.CreateBucket(expected)
	assert.NoError(t, err)

	// find

	got, err := m.FindBucketByName(
		"bucket name",
		interfaces.FullBucket,
	)

	assert.NoError(t, err)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.BucketName, got.BucketName)
	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.Equal(t, expected.UpdatedAt.UnixNano(), got.UpdatedAt.UnixNano())
	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.MetaDataStoreName, got.MetaDataStoreName)
	assert.Equal(t, expected.StructuralDataStoreName, got.StructuralDataStoreName)
	assert.Equal(t, expected.DataStoreName, got.DataStoreName)
	assert.Equal(t, expected.LuaScript, got.LuaScript)

	// update

	time.Sleep(time.Microsecond * 1)

	expected.BucketName = "new bucket name"
	expected.MetaData["b"] = "new value"
	err = m.UpdateBucket(expected, interfaces.PrimaryNamesData|interfaces.MetaData)
	assert.NoError(t, err)

	// find not existing bucket

	got, err = m.FindBucketByName(
		"bucket name",
		interfaces.FullBucket,
	)

	assert.EqualError(t, err, interfaces.ErrNotFound.Error())

	// find existing bucket

	got, err = m.FindBucketByName(
		"new bucket name",
		interfaces.FullBucket,
	)

	assert.NoError(t, err)
	assert.Equal(t, expected.BucketID, got.BucketID)
	assert.Equal(t, expected.BucketName, got.BucketName)
	assert.Equal(t, expected.RawData, got.RawData)
	assert.Equal(t, expected.CreatedAt.UnixNano(), got.CreatedAt.UnixNano())
	assert.Equal(t, expected.UpdatedAt.UnixNano(), got.UpdatedAt.UnixNano())
	assert.Equal(t, expected.Owners[0], got.Owners[0])
	assert.Equal(t, expected.Owners[1], got.Owners[1])
	assert.Equal(t, expected.MetaData["a"], got.MetaData["a"])
	assert.Equal(t, expected.MetaData["b"], got.MetaData["b"])
	assert.Equal(t, expected.StructuralData["c"], got.StructuralData["c"])
	assert.Equal(t, expected.StructuralData["d"], got.StructuralData["d"])
	assert.Equal(t, expected.MetaDataStoreName, got.MetaDataStoreName)
	assert.Equal(t, expected.StructuralDataStoreName, got.StructuralDataStoreName)
	assert.Equal(t, expected.DataStoreName, got.DataStoreName)
	assert.Equal(t, expected.LuaScript, got.LuaScript)
}
