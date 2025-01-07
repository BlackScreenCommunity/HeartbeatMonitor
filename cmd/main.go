package main

import (
	"project/internal/pluginConfigurationLoader"
	"project/internal/pluginDispatcher"
	webserver "project/internal/webServer"
)

func main() {

	pluginsConfig, shouldReturn := pluginConfigurationLoader.GetPluginsConfiguration()
	if shouldReturn {
		return
	}
	pluginDispatcher.InitializePlugins(pluginsConfig)

	webserver.RunServer()
}
