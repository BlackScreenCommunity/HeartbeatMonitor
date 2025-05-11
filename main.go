package main

import (
	"flag"
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"project/internal/webServer"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "configFilePath", "appsettings.json", "Path to config file")

	flag.Parse()

	applicationConfiguration, shouldReturn := applicationConfigurationDispatcher.GetConfigFromFile(configFilePath)
	if shouldReturn {
		return
	}

	pluginDispatcher.InitializePlugins(applicationConfiguration.Plugins)
	agentDispatcher.InitializePlugins(applicationConfiguration.Agents, applicationConfiguration.WebServer.AgentPollingTimeout)

	webServer.RunServer(applicationConfiguration.WebServer, applicationConfiguration.Server)
}
