package connpool

import (
	"log"
	"context"
	"testing"
	"net"
	"time"
)

func TestPool(t *testing.T) {
	go startTcpServer()
	time.Sleep(time.Second * 3)
	pool := GetPool("default")
	_, err := pool.Get(context.Background(), "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
}

func startTcpServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

}