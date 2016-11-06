package api

import (
	"interfaces"
	"time"
)

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

func RefreshEvery(d time.Duration, f func() error) {

	for {
		select {
		case <-time.After(d):
			if err := f(); err != nil {
				logger.Println("[RefreshEvery] abort, err:", err)
				return
			}
		}
	}

}
