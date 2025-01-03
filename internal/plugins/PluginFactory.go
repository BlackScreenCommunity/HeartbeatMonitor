package plugins

import (
	"fmt"
)

func init() {
	RegisterPluginType("VersionPlugin", func() Plugin { return &VersionPlugin{} })
	RegisterPluginType("HardDriveFreeSpacePlugin", func() Plugin { return &HardDriveFreeSpacePlugin{} })
}

var registeredPlugins = make(map[string]func() Plugin)

func RegisterPluginType(name string, constructor func() Plugin) {
	registeredPlugins[name] = constructor
}

func CreatePlugin(name string) (Plugin, error) {
	constructor, exists := registeredPlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin type '%s' not registered", name)
	}
	return constructor(), nil
}
