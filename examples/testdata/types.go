package testdata

type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}

type Request struct {	
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}

type CountResponse struct {
	Count int64 `mapstructure:"count" msgpack:"count"`
}
