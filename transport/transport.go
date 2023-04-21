package transport

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/dayueba/mrpc/codec"
	"github.com/dayueba/mrpc/codes"
)

const DefaultPayloadLength = 1024
const MaxPayloadLength = 4 * 1024 * 1024

type ServerTransport interface {
	// 监听和处理请求
	ListenAndServe(context.Context, ...ServerTransportOption) error
}

type ClientTransport interface {
	// 发送请求
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}

type Framer interface {
	ReadFrame(net.Conn) ([]byte, error)
}

type framer struct {
	buffer  []byte
	counter int // to prevent the dead loop
}

func NewFramer() Framer {
	return &framer{
		buffer: make([]byte, DefaultPayloadLength),
	}
}

func (f *framer) Resize() {
	f.buffer = make([]byte, len(f.buffer)*2)
}

func (f *framer) ReadFrame(conn net.Conn) ([]byte, error) {
	frameHeader := make([]byte, codec.FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != codec.FrameHeadLen || err != nil {
		//fmt.Println("err: ", err)
		return nil, err
	}

	length := binary.BigEndian.Uint32(frameHeader) // 目前header里只有length

	if length > MaxPayloadLength {
		return nil, codes.NewFrameworkError(codes.ClientMsgErrorCode, "payload too large...")
	}

	for uint32(len(f.buffer)) < length && f.counter <= 12 {
		f.buffer = make([]byte, len(f.buffer)*2)
		f.counter++
	}

	if num, err := io.ReadFull(conn, f.buffer[:length]); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, f.buffer[:length]...), nil
}
