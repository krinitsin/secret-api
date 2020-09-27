package secret_api

import (
	"time"
)

// Options опции
type Options struct {
	Timeout        time.Duration
	PolligInterval time.Duration
	PollingDelay   time.Duration
}

var (
	options *Options
)

func Init(o *Options) {
	options = o
}
