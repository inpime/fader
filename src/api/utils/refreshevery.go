package utils

import (
	"time"
)

func RefreshEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}
