package sre

import (
	"math/rand"
	"sync"
	"time"
	"math"

	"github.com/dayueba/mrpc/interceptor/breaker"

	"github.com/go-kratos/aegis/pkg/window"
)

var _ breaker.Breaker = (*Breaker)(nil)

type Breaker struct {
	stat window.RollingCounter
	r    *rand.Rand
	// rand.New(...) 不是并发安全的
	randLock sync.Mutex

	// Reducing the k will make adaptive throttling behave more aggressively,
	// Increasing the k will make adaptive throttling behave less aggressively.
	k       float64
	request int64

	state int32
}

func NewBreaker(opts ...Option) breaker.Breaker {
	opt := options{
		success: 0.6,
		request: 100,
		bucket:  10,
		window:  3 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}
	counterOpts := window.RollingCounterOpts{
		Size:           opt.bucket,
		BucketDuration: time.Duration(int64(opt.window) / int64(opt.bucket)),
	}
	stat := window.NewRollingCounter(counterOpts)
	return &Breaker{
		stat:    stat,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		request: opt.request,
		k:       1 / opt.success,
	}
}

func (b *Breaker) summary() (success int64, total int64) {
	b.stat.Reduce(func(iterator window.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

// Allow request if error returns nil.
func (b *Breaker) Allow() error {
	accepts, total := b.summary()
	requests := b.k * float64(accepts)
	if total < b.request || float64(total) < requests {
		return nil
	}
	dr := math.Max(0, (float64(total)-requests)/float64(total+1))
	drop := b.trueOnProba(dr)
	if drop {
		return breaker.ErrNotAllowed
	}
	return nil
}

// MarkSuccess mark request is success.
func (b *Breaker) MarkSuccess() {
	b.stat.Add(1)
}

// MarkFailed mark request is failed.
func (b *Breaker) MarkFailed() {
	// NOTE: when client reject request locally, continue to add counter let the
	// drop ratio higher.
	b.stat.Add(0)
}

func (b *Breaker) trueOnProba(proba float64) (truth bool) {
	b.randLock.Lock()
	truth = b.r.Float64() < proba
	b.randLock.Unlock()
	return
}
