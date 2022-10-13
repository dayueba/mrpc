package protocol

type Request struct {
	MsgId string
	Type  string
  Method string
  Params interface{}
}

type Response struct {
	// MsgId uint32
	// Status string
	// Result interface{}
	
}