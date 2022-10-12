package codec

import (
	"bytes"
	"encoding/binary"
	"math"

	// "github.com/golang/protobuf/proto"
)

type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

const FrameHeadLen = 4

type FrameHeader struct {
	Length uint32  	// total packet length
}

// GetCodec get a Codec by a codec name
func GetCodec(name string) Codec {
	if codec, ok := codecMap[name]; ok {
		return codec
	}
	return DefaultCodec
}

var codecMap = make(map[string]Codec)

var DefaultCodec = NewCodec()

var NewCodec = 	func () Codec {
	return &defaultCodec{}
}

func init() {
	RegisterCodec("proto", DefaultCodec)
}

// RegisterCodec registers a codec, which will be added to codecMap
func RegisterCodec(name string, codec Codec) {
	if codecMap == nil {
		codecMap = make(map[string]Codec)
	}
	codecMap[name] = codec
}

func (c *defaultCodec) Encode(data []byte) ([]byte, error) {
	totalLen := FrameHeadLen + len(data)
	buffer := bytes.NewBuffer(make([]byte, 0, totalLen))

	frame := FrameHeader{
		Length: uint32(len(data)),
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}


func (c *defaultCodec) Decode(frame []byte) ([]byte,error) {
	return frame[FrameHeadLen:], nil
}

type defaultCodec struct{}

func upperLimit(val int) uint32 {
	if val > math.MaxInt32 {
		return uint32(math.MaxInt32)
	}
	return uint32(val)
}