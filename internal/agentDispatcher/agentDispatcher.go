package agentDispatcher

import (
	"encoding/json"
	"io"
	"net/http"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/utils"
	"time"
)

var agents = make([]applicationConfigurationDispatcher.AgentConfig, 0)

func InitializePlugins(agentsConfigCollection []applicationConfigurationDispatcher.AgentConfig) {
	agents = agentsConfigCollection
}

func GetMetricsFromAgents() map[string]interface{} {
	agentResultCollection := make(map[string]interface{})

	for _, agent := range agents {
		agentResultCollection[agent.Name] = GetMetricsFromSingleAgent(agent)
	}
	return utils.MapDereference(agentResultCollection)
}

func GetMetricsFromSingleAgent(agent applicationConfigurationDispatcher.AgentConfig) map[string]interface{} {
	results := make(map[string]interface{})

	transport := &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get(agent.Address + "/plugins/results")
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
