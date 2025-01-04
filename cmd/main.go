package main

import (
	"project/internal/pluginConfigurationLoader"
	"project/internal/pluginDispatcher"
)

func main() {

	pluginsConfig, shouldReturn := pluginConfigurationLoader.GetPluginsConfiguration()
	if shouldReturn {
		return
	}

	pluginDispatcher.InitializePlugins(pluginsConfig)
	pluginDispatcher.CollectAll()
}
