package connpool

import (
	"context"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Pool interface {
	Get(ctx context.Context, address string) (*Conn, error)
}

// 全局连接池
type pool struct {
	opts      *Options
	connPools *sync.Map
	sg        singleflight.Group
}

var poolMap = make(map[string]Pool)

func init() {
	registorPool("default", DefaultPool)
}

func registorPool(poolName string, pool Pool) {
	poolMap[poolName] = pool
}

var DefaultPool = NewPool()

func NewPool(opt ...Option) *pool {
	// default options
	opts := &Options{
		maxCap:      100,
		initCap:     1,
		maxIdle: 10,
		idleTimeout: 1 * time.Minute,
		dialTimeout: 200 * time.Millisecond,
	}
	m := &sync.Map{}

	p := &pool{
		connPools: m,
		opts:      opts,
	}
	for _, o := range opt {
		o(p.opts)
	}

	return p
}

func GetPool(poolName string) Pool {
	if v, ok := poolMap[poolName]; ok {
		return v
	}
	return DefaultPool
}

func (p *pool) Get(ctx context.Context, address string) (*Conn, error) {
	// 先从 map 中尝试获取 key 为 address 的子连接池
	if value, ok := p.connPools.Load(address); ok {
		if cp, ok := value.(*channelPool); ok {
			conn, err := cp.Get(ctx)
			return conn, err
		}
	}

	// 创建新的连接池
	v, err, _ := p.sg.Do(address, func() (interface{}, error) {
		cp, err := p.NewConnPool(ctx, address)
		if err != nil {
			return nil, err
		}
		return cp, nil
	})
	if err != nil {
		return nil, err
	}
	cp := v.(*channelPool)

	p.connPools.Store(address, cp)

	return cp.Get(ctx)
}

func (p *pool) NewConnPool(ctx context.Context, address string) (*channelPool, error) {
	c := &channelPool{
		initCap: p.opts.initCap,
		maxCap:  p.opts.maxCap,
		maxIdle: p.opts.maxIdle,
		Dial: func(ctx context.Context) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := p.opts.dialTimeout
			if t, ok := ctx.Deadline(); ok {
				timeout = time.Until(t)
			}

			return net.DialTimeout("tcp", address, timeout)
		},
		conns:       make(chan *Conn, p.opts.maxIdle),
		idleTimeout: p.opts.idleTimeout,
		dialTimeout: p.opts.dialTimeout,
	}

	for i := 0; i < p.opts.initCap; i++ {
		conn, err := c.Dial(ctx)
		if err != nil {
			return nil, err
		}
		c.Put(c.wrapConn(conn))
	}

	c.RegisterChecker(3*time.Second, c.Checker)
	return c, nil
}
