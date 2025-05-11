package applicationConfigurationDispatcher

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"project/internal/pluginDispatcher"
)

// Defines the local agent name
type ServerInfo struct {
	Name      string
	Login     string
	Password  string
	AgentMode bool
}

// Parameters for the web interface
// and network interaction with external agents
type WebServerConfig struct {
	Port                int
	AgentPollingTimeout int
}

// Remote agent parameters
type AgentConfig struct {
	Address  string
	Name     string
	Active   bool
	Login    string
	Password string
}

// Overall set of application parameters
type ApplicationConfiguration struct {
	Server    ServerInfo
	WebServer WebServerConfig
	Plugins   []pluginDispatcher.PluginConfig
	Agents    []AgentConfig
}

// Reads application's configuration
// from file located at the path
func GetConfigFromFile(configFilePath string) (ApplicationConfiguration, bool) {
	pluginsConfig, err := loadConfigFromFile(configFilePath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return ApplicationConfiguration{}, true
	}
	return pluginsConfig, false
}

// Handles operations with file.
// Tries to search and decode JSON object to
// Application Configuration
func loadConfigFromFile(configFilePath string) (ApplicationConfiguration, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Cannot read config file %s: %v", configFilePath, err)
		return ApplicationConfiguration{}, fmt.Errorf("error opening config file: %v", err)
	}

	defer file.Close()

	var data ApplicationConfiguration
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return ApplicationConfiguration{}, fmt.Errorf("error decoding config file: %v", err)
	}

	return data, nil
}
