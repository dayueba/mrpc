package connpool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dayueba/mrpc/log"
)

var oneByte = make([]byte, 1)

type channelPool struct {
	initCap     int
	maxCap      int
	maxIdle     int
	idleTimeout time.Duration
	dialTimeout time.Duration
	Dial        func(context.Context) (net.Conn, error)
	conns       chan *Conn
	mu          sync.Mutex
	inflight    int32
}

func (c *channelPool) Get(ctx context.Context) (*Conn, error) {
	if c.conns == nil {
		return nil, ErrConnClosed
	}
	select {
	case conn := <-c.conns:
		if conn == nil {
			return nil, ErrConnClosed
		}

		if conn.unusable {
			return nil, ErrConnClosed // 这里出错了没有自动重试，调用方根据错误类型来决定是否重试
		}

		return conn, nil
	case <-ctx.Done(): // context取消或超时，则退出
		return nil, ctx.Err()
	default:
		if c.inflight > int32(c.maxCap) {
			select {
			case conn := <-c.conns:
				if conn == nil {
					return nil, ErrConnClosed
				}

				if conn.unusable {
					return nil, ErrConnClosed // 这里出错了没有自动重试，调用方根据错误类型来决定是否重试
				}

				return conn, nil
			case <-ctx.Done(): // context取消或超时，则退出
				return nil, ctx.Err()
			}
		}

		conn, err := c.Dial(ctx)
		if err != nil {
			return nil, err
		}
		atomic.AddInt32(&c.inflight, 1)
		return c.wrapConn(conn), nil
	}
}

func (c *channelPool) Put(conn *Conn) error {
	if conn == nil {
		return errors.New("connection closed")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conns == nil {
		conn.MarkUnusable()
		conn.Close()
	}

	select {
	case c.conns <- conn:
		return nil
	default:
		return conn.Conn.Close()
	}
}

func (c *channelPool) wrapConn(conn net.Conn) *Conn {
	p := &Conn{
		c:           c,
		t:           time.Now(),
		dialTimeout: c.dialTimeout,
	}
	p.Conn = conn
	return p
}

func (c *channelPool) RegisterChecker(internal time.Duration, checker func(conn *Conn) bool) {
	if internal <= 0 || checker == nil {
		return
	}

	go func() {
		for {
			time.Sleep(internal)
			length := len(c.conns)
			for i := 0; i < length; i++ {
				select {
				case pc := <-c.conns:
					if !checker(pc) {
						pc.MarkUnusable()
						pc.Close()
					} else {
						err := c.Put(pc)
						if err != nil {
							log.Infof("put err: %v", err)
						}
					}
				default:
				}
			}
		}
	}()
}

// 负责校验连接是否存活
func (c *channelPool) Checker(pc *Conn) bool {
	// check timeout
	if pc.t.Add(c.idleTimeout).Before(time.Now()) {
		return false
	}

	// check conn is alive or not
	if !isConnAlive(pc.Conn) {
		return false
	}

	return true
}

func isConnAlive(conn net.Conn) bool {
	_ = conn.SetReadDeadline(time.Now().Add(time.Millisecond))

	if n, err := conn.Read(oneByte); n > 0 || err == io.EOF {
		return false
	}

	_ = conn.SetReadDeadline(time.Time{})
	return true
}
