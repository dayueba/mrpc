package breaker

import (
	"context"
	"errors"

	"github.com/dayueba/mrpc/interceptor"
	"github.com/dayueba/mrpc/interceptor/breaker"
)

var ErrNotAllowed = errors.New("request failed due to circuit breaker triggered")

type Option func(*options)

func WithBreaker(limiter breaker.Breaker) Option {
	return func(o *options) {
		o.breaker = limiter
	}
}

type options struct {
	breaker breaker.Breaker
}

// Server ratelimiter middleware
func Client(opts ...Option) interceptor.ClientInterceptor {
	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return func(ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {
		if err := options.breaker.Allow(); err != nil {
			// rejected
			// NOTE: when client reject requests locally,
			// continue to add counter let the drop ratio higher.
			options.breaker.MarkFailed()
			return ErrNotAllowed
		}

		err := ivk(ctx, req, rsp)
		if err != nil { // todo check res error
			options.breaker.MarkFailed()
		} else {
			options.breaker.MarkSuccess()
		}

		return err
	}
}
