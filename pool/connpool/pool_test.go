package connpool

import (
	"context"
	"sync"
	"testing"
)

func TestConnPool(t *testing.T) {
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	pool := GetPool("default")
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()

			conn, err := pool.Get(context.Background(), "localhost:8000")
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
		}()
	}
	wg.Wait()
}

func BenchmarkPool_Get(b *testing.B) {
	pool := GetPool("default")
	for i := 0; i < b.N; i++ {
		conn, err := pool.Get(context.Background(), "localhost:8000")
		if err != nil {
			b.Fatal(err)
		}
		defer conn.Close()
	}
}
