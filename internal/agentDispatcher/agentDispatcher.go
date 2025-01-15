package agentDispatcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/utils"
)

var agents = make([]applicationConfigurationDispatcher.AgentConfig, 0)
var agentsMetricsCollection = make(map[string]interface{})

func InitializePlugins(agentsConfigCollection []applicationConfigurationDispatcher.AgentConfig) {
	agents = agentsConfigCollection
}

func CollectAll() map[string]interface{} {

	agentResultCollection := make(map[string]interface{})

	for _, agent := range agents {
		var results map[string]interface{}

		url := agent.Address + "/plugins/results"

		resp, err := http.Get(url)
		if err != nil {
			agentResultCollection[agent.Address] = err.Error()
			continue
		}

		defer resp.Body.Close()

		if err != nil {
			agentResultCollection[agent.Address] = err.Error()
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			agentResultCollection[agent.Address] = err.Error()
		}

		if resp.StatusCode != http.StatusOK {
			agentResultCollection[agent.Address] = string(body)
		}

		if err := json.Unmarshal(body, &results); err != nil {
			agentResultCollection[agent.Address] = err.Error()
		}

		agentResultCollection[agent.Address] = results
	}
	return utils.MapDereference(agentResultCollection)
}
