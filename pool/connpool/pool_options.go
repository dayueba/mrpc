package connpool

import "time"

type Options struct {
	initCap     int // initial capacity
	maxCap      int // max capacity
	maxIdle int
	idleTimeout time.Duration
	dialTimeout time.Duration // dial timeout
}

type Option func(*Options)

func WithInitialCap(initialCap int) Option {
	return func(o *Options) {
		o.initCap = initialCap
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}

func WithMaxCap(maxCap int) Option {
	return func(o *Options) {
		o.maxCap = maxCap
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = idleTimeout
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = dialTimeout
	}
}
