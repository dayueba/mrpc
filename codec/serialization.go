package codec

import (
	"errors"
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

const (
	MsgPack = "msgpack" // msgpack
)

var serializationMap = make(map[string]Serialization)

var DefaultSerialization = NewSerialization()

var NewSerialization = func() Serialization {
	return &MsgpackSerialization{}
}

func init() {
	registerSerialization("msgpack", DefaultSerialization)
}

func registerSerialization(name string, serialization Serialization) {
	if serializationMap == nil {
		serializationMap = make(map[string]Serialization)
	}
	serializationMap[name] = serialization
}

func GetSerialization(name string) Serialization {
	if v, ok := serializationMap[name]; ok {
		return v
	}
	return DefaultSerialization
}

type MsgpackSerialization struct {}


func (c *MsgpackSerialization) Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("marshal nil interface{}")
	}

	return msgpack.Marshal(v)
}

func (c *MsgpackSerialization) Unmarshal(data []byte, v interface{}) error {
	if data == nil || len(data) == 0 {
		return errors.New("unmarshal nil or empty bytes")
	}

	decoder := msgpack.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return err
}
