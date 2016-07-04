package utils

import (
	"time"
)

func RefreshEvery(d time.Duration, f func()) {

	for {
		select {
		case <-time.After(d):
			f()
		}
	}

}
