package connpool

import (
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrConnClosed = errors.New("connection closed ...")
)

// PoolConn 通过装饰者模式对原生连接 net.Conn 进行了修饰。这里也是通过互斥锁来保证并发安全，只不过这里粒度更细，用了读写锁 sync.RWMutex。
type PoolConn struct {
	net.Conn
	c *channelPool
	unusable bool		// if unusable is true, the conn should be closed
	mu sync.RWMutex
	t time.Time  // connection idle time
	dialTimeout time.Duration // connection timeout duration
}

// overwrite conn Close for connection reuse
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}

	// reset connection deadline
	p.Conn.SetDeadline(time.Time{})

	// 如果连接正常 则放回连接池
	return p.c.Put(p)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (p *PoolConn) Read(b []byte) (int, error) {
	// 判断连接池是否为不可读状态
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Read(b)
	if err != nil {
		p.MarkUnusable()
		// 关闭连接
		p.Conn.Close()
	}
	return n, err
}

func (p *PoolConn) Write(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Write(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}

func (c *channelPool) wrapConn(conn net.Conn) *PoolConn {
	p := &PoolConn {
		c : c,
		t : time.Now(),
		dialTimeout: c.dialTimeout,
	}
	p.Conn = conn
	return p
}

