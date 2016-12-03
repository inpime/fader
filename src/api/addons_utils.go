package api

import (
	"interfaces"

	uuid "github.com/satori/go.uuid"
)

func filesByBucketID(id uuid.UUID) (res []*interfaces.File) {
	fileManager.(interfaces.FileImportManager).
		EachFile(func(item *interfaces.File) error {
			if uuid.Equal(id, item.BucketID) {
				res = append(res, item)
			}
			return nil
		})
	return
}

func listOfBuckets() (res []*interfaces.Bucket) {
	bucketManager.(interfaces.BucketImportManager).
		EachBucket(func(item *interfaces.Bucket) error {
			res = append(res, item)
			return nil
		})
	return
}
