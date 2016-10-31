package api

import (
	"interfaces"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestImporter_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	bucket1 := uuid.NewV4()
	bucket2 := uuid.NewV4()

	buckets := []*interfaces.Bucket{
		&interfaces.Bucket{
			BucketID:       bucket1,
			BucketName:     "bucket1",
			MetaData:       make(map[string]interface{}),
			StructuralData: make(map[string]interface{}),
			RawData: []byte(`1
            2
            b`),
		},
		&interfaces.Bucket{
			BucketID:       bucket2,
			BucketName:     "bucket2",
			MetaData:       make(map[string]interface{}),
			StructuralData: make(map[string]interface{}),
			RawData: []byte(`1
            2
            b`),
		},
	}

	files := []*interfaces.File{
		&interfaces.File{
			FileID:         uuid.NewV4(),
			BucketID:       bucket1,
			FileName:       "file1-1",
			MetaData:       make(map[string]interface{}),
			StructuralData: make(map[string]interface{}),
			RawData: []byte(`1
            2
            b`),
		},
		&interfaces.File{
			FileID:         uuid.NewV4(),
			BucketID:       bucket1,
			FileName:       "file1-2",
			MetaData:       make(map[string]interface{}),
			StructuralData: make(map[string]interface{}),
			RawData: []byte(`1
            2
            b`),
		},
		&interfaces.File{
			FileID:         uuid.NewV4(),
			BucketID:       bucket2,
			FileName:       "file2-1",
			MetaData:       make(map[string]interface{}),
			StructuralData: make(map[string]interface{}),
			RawData: []byte(`1
            2
            b`),
		},
	}

	//
	importManager := interfaces.NewImportManager(
		bucketManager,
		fileManager,
	)

	for _, bucket := range buckets {
		err := bucketManager.CreateBucket(bucket)
		assert.NoError(t, err, "Create bucket", bucket.BucketID)
	}

	for _, file := range files {
		err := fileManager.CreateFile(file)
		assert.NoError(t, err, "Create file", file.FileID)
	}
	description := `descrption
    multiline
	multiline
	multiline
	multiline`
	data, err := importManager.Export("vDEV", "authrorDEV", description)
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	// Clear

	err = db.Close()
	assert.NoError(t, err)
	err = os.RemoveAll(settings.DatabasePath)
	assert.NoError(t, err)

	// Check

	err = Setup(e, &Settings{})
	assert.NoError(t, err)
	importManager = interfaces.NewImportManager(
		bucketManager,
		fileManager,
	)

	info, err := importManager.Import(data)
	assert.NoError(t, err, "inport")
	assert.EqualValues(t, info.Description(), description)

	fileNames := []string{
		"bucket1", "file1-1",
		"bucket1", "file1-2",
		"bucket2", "file2-1",
	}

	for i := 0; i < len(fileNames); i += 2 {
		t.Log(i, fileNames[i], fileNames[i+1])
		bucketName := fileNames[i]
		fileName := fileNames[i+1]

		file, err := fileManager.FindFileByName(
			bucketName,
			fileName,
			interfaces.FullFile,
		)

		assert.NoError(t, err, "find file: %v, %v", bucketName, fileName)
		assert.Equal(t, file.FileName, fileName)
		assert.Equal(t, file.RawData, []byte(`1
            2
            b`))
	}
}
