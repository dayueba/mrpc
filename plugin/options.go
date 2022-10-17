package plugin


type Options struct {
	SvrAddr string     // server address
	Services []string   // service arrays
	SelectorSvrAddr string  // server discovery address ，e.g. consul server address
}

// Option provides operations on Options
type Option func(*Options)

// WithSvrAddr allows you to set SvrAddr of Options
func WithSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SvrAddr = addr
	}
}

// WithSvrAddr allows you to set Services of Options
func WithServices(services []string) Option {
	return func(o *Options) {
		o.Services = services
	}
}

// WithSvrAddr allows you to set SelectorSvrAddr of Options
func WithSelectorSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SelectorSvrAddr = addr
	}
}