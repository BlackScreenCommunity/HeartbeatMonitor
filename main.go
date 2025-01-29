package main

import (
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"project/internal/webServer"
)

func main() {
	applicationConfiguration, shouldReturn := applicationConfigurationDispatcher.GetConfigFromFile()
	if shouldReturn {
		return
	}

	pluginDispatcher.InitializePlugins(applicationConfiguration.Plugins)
	agentDispatcher.InitializePlugins(applicationConfiguration.Agents, applicationConfiguration.WebServer.AgentPollingTimeout)

	webServer.RunServer(applicationConfiguration.WebServer, applicationConfiguration.Server)
}
