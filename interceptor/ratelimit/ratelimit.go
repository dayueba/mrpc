package ratelimit

import (
	"errors"
)

var (
	ErrLimitExceed = errors.New("rate limit exceeded")
)

type DoneFunc func(DoneInfo)

type DoneInfo struct {
	Err error
}

// 限流器只需要一个 Allow 接口
type Limiter interface {
	Allow() (DoneFunc, error)
}
