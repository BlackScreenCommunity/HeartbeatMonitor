package main

import (
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	webserver "project/internal/webServer"
)

func main() {
	applicationConfiguration, shouldReturn := applicationConfigurationDispatcher.GetConfigFromFile()
	if shouldReturn {
		return
	}

	pluginDispatcher.InitializePlugins(applicationConfiguration.Plugins)
	webserver.RunServer(applicationConfiguration.Webserver)
}
