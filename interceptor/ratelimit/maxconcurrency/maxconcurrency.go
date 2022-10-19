package maxconcurrency

import (
	"sync"
	"sync/atomic"

	"github.com/dayueba/mrpc/interceptor/ratelimit"
)

var _ ratelimit.Limiter = (*Maxconcurrency)(nil)

type Maxconcurrency struct {
	// 当前系统中的请求数，数据得来方法是：中间件原理在处理前+1，处理handle之后不管成功失败都减去1
	inFlight int64
	maxInFlight int64
	mu sync.Mutex
}

func (l *Maxconcurrency) Allow() (ratelimit.DoneFunc, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.maxInFlight > 0 && l.maxInFlight < l.inFlight {
		return nil, ratelimit.ErrLimitExceed
	}

	atomic.AddInt64(&l.inFlight, 1)

	return func(ratelimit.DoneInfo) {
		atomic.AddInt64(&l.inFlight, -1)
	}, nil
}

type Option func(*options)

type options struct {
	Maxconcurrency int64
}

func WithMaxconcurrency(max int64) Option {
	return func(o *options) {
		o.Maxconcurrency = max
	}
}


func NewLimiter(opts ...Option) *Maxconcurrency {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	limiter := &Maxconcurrency{
		maxInFlight: opt.Maxconcurrency,
	}

	return limiter
}