package agentDispatcher

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"project/internal/applicationConfigurationDispatcher"
	"time"
)

// Collection of registered agents
var agents = make([]applicationConfigurationDispatcher.AgentConfig, 0)

// The default timeout for waiting
// for a response from the agent
// Can be overridden in the agent's
// configuration in appsettings.json
var pollingTimeout = 30

// TODO Rename function
func InitializePlugins(agentsConfigCollection []applicationConfigurationDispatcher.AgentConfig, agentsPollingTimeout int) {
	agents = agentsConfigCollection
	pollingTimeout = agentsPollingTimeout
}

// Fetchs data from an external agent via API
func GetMetricsFromSingleAgent(agent applicationConfigurationDispatcher.AgentConfig) map[string]interface{} {
	results := make(map[string]interface{})

	transport := &http.Transport{
		ResponseHeaderTimeout: time.Duration(pollingTimeout) * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", agent.Address+"/plugins/results", nil)
	if err != nil {
		results["Error"] = err.Error()
		return results
	}
	req.SetBasicAuth(agent.Login, agent.Password)

	resp, err := client.Do(req)
	if err != nil {
		results["Error"] = err.Error()
		return results
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Agent %s doesn't response correctly: %v", agent.Name, err)
		}
	}()

	return ParseResponseFromAgent(resp)

}

// Read the response from an agent
// Convert the response body to a map
// If there is an error during reading, parsing, or if the HTTP status code is not OK,
// Return "Error" if there is an error during reading, parsing, or if the HTTP status code is not 200

func ParseResponseFromAgent(resp *http.Response) map[string]interface{} {

	results := make(map[string]interface{})

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		results["Error"] = err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		results["Error"] = string(body)
	}

	if err := json.Unmarshal(body, &results); err != nil {
		results["Error"] = err.Error()
	}

	return results
}

//	Provides an interface for
//
// accessing the list of registered agents
func GetAgents() []applicationConfigurationDispatcher.AgentConfig {
	return agents
}
