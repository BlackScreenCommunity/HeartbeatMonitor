package main

import (
	"project/internal/agentDispatcher"
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
	agentDispatcher.InitializePlugins(applicationConfiguration.Agents)

	webserver.RunServer(applicationConfiguration.Webserver, applicationConfiguration.Server)
}
