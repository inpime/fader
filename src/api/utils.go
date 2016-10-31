package api

import "interfaces"

// IsFirstStart the app first launched
func IsFirstStart() bool {
	count := 0
	bucketManager.(interfaces.BucketImportManager).
		EachBucket(func(bucket *interfaces.Bucket) error {
			count++
			return nil
		})

	return count == 0
}
