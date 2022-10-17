package plugin

type Plugin interface {
	Init(...Option) error
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
