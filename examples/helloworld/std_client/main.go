package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", "127.0.0.1:8000")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			err = conn.SetReadDeadline(time.Now().Add(time.Second))
			if err != nil {
				fmt.Println(err)
			}

			conn.Write([]byte("message"))
		}()
	}
	wg.Wait()
}
