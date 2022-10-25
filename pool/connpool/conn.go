package connpool

import (
	"errors"
	"net"
	"sync"
	"time"
)

var ErrConnClosed = errors.New("connection closed")

var _ net.Conn = (*Conn)(nil)

type Conn struct {
	net.Conn
	c           *channelPool
	unusable    bool
	mu          sync.Mutex
	t           time.Time     // 该连接的空闲时间
	dialTimeout time.Duration // connection timeout duration
}

func (p *Conn) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

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

func (p *Conn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (p *Conn) Read(b []byte) (int, error) {
	// 判断该连接状态
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

func (p *Conn) Write(b []byte) (int, error) {
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
