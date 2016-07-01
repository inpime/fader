package config

import (
	"time"
)

func refreshEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}
