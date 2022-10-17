package plugin

import "github.com/opentracing/opentracing-go"

type Plugin interface {
}

var PluginMap = make(map[string]Plugin)

func Register(name string, plugin Plugin) {
   if PluginMap == nil {
      PluginMap = make(map[string]Plugin)
   }
   PluginMap[name] = plugin
}

type ResolverPlugin interface {
	Init(...Option) error
}

type TracingPlugin interface {
   Init(...Option) (opentracing.Tracer, error)
}
