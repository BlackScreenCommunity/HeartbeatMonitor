package agentDispatcher

import (
	"encoding/json"
	"io"
	"net/http"
	"project/internal/applicationConfigurationDispatcher"
	"strconv"
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

// Gets metric data from a collection of agents
func GetMetricsFromAgents() map[string]interface{} {
	resultsChannel := make(chan struct {
		Key    string
		Result map[string]interface{}
	}, len(agents))

	agentResultCollection := make(map[string]interface{})

	for i, agent := range agents {
		if agent.IsActive {
			go func(i int, agent applicationConfigurationDispatcher.AgentConfig) {
				result := GetMetricsFromSingleAgent(agent)
				resultsChannel <- struct {
					Key    string
					Result map[string]interface{}
				}{
					Key:    strconv.Itoa(i+1) + ". " + agent.Name,
					Result: result,
				}
			}(i, agent)
		}
	}

	for range agents {
		res := <-resultsChannel
		agentResultCollection[res.Key] = res.Result
	}

	return agentResultCollection
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
	defer resp.Body.Close()

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
