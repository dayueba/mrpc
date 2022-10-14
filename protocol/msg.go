package protocol

type Request struct {
	MsgId string
	Type  string
  Method string
  Params interface{}
}

type RpcError struct {
	ErrorCode int `msgpack:"error_code"`
	Method string `msgpack:"method"`
	Message string `msgpack:"message"`
}

func (r RpcError) Error() string {
	return r.Message
}

type Response struct {
	// MsgId uint32
	// Status string
	// Result interface{}
	
}