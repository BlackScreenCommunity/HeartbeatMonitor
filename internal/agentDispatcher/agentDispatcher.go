package agentDispathcer

import (
	"encoding/json"
	"fmt"
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

	for name, agent := range agents {

		url := agent.Address + "/plugins/results"

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making GET request: %v\n", err)
		}
		defer resp.Body.Close()

		if err != nil {
			fmt.Printf("Error collecting data from plugin %s: %v\n", name, err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
			fmt.Printf("Response body: %s\n", string(body))
		}

		var results map[string]interface{}
		if err := json.Unmarshal(body, &results); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
		}

		agentResultCollection[agent.Address] = results
	}
	return utils.MapDereference(agentResultCollection)
}
