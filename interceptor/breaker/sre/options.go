package sre

import "time"

type options struct {
	success float64
	request int64
	bucket  int
	window  time.Duration
}

type Option func(*options)

// WithSuccess with the K = 1 / Success value of sre breaker, default success is 0.5
// Reducing the K will make adaptive throttling behave more aggressively,
// Increasing the K will make adaptive throttling behave less aggressively.
func WithSuccess(s float64) Option {
	return func(c *options) {
		c.success = s
	}
}

// WithRequest with the minimum number of requests allowed.
func WithRequest(r int64) Option {
	return func(c *options) {
		c.request = r
	}
}

// WithWindow with the duration size of the statistical window.
func WithWindow(d time.Duration) Option {
	return func(c *options) {
		c.window = d
	}
}

// WithBucket set the bucket number in a window duration.
func WithBucket(b int) Option {
	return func(c *options) {
		c.bucket = b
	}
}
