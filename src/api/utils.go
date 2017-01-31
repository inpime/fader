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

////////////////////////////////////////////////////////////////////////////////
// Types
////////////////////////////////////////////////////////////////////////////////

func ToFloat64(v interface{}) (f float64) {
	switch _v := v.(type) {
	case int:
		f = float64(_v)
	case int16:
		f = float64(_v)
	case int32:
		f = float64(_v)
	case int64:
		f = float64(_v)
	case int8:
		f = float64(_v)
	case float32:
		f = float64(_v)
	case float64:
		f = float64(_v)
	case uint:
		f = float64(_v)
	case uint16:
		f = float64(_v)
	case uint32:
		f = float64(_v)
	case uint64:
		f = float64(_v)
	case uint8:
		f = float64(_v)
	default:
		f = 0
	}

	return
}