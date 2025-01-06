package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/pluginConfigurationLoader"
	"project/internal/pluginDispatcher"
)

var ResponseData = ""

func main() {

	pluginsConfig, shouldReturn := pluginConfigurationLoader.GetPluginsConfiguration()
	if shouldReturn {
		return
	}

	pluginDispatcher.InitializePlugins(pluginsConfig)
	pluginResultCollection := pluginDispatcher.CollectAll()

	cleanResults := MapDereference(pluginResultCollection)
	pluginDispatcher.PrintPluginResult(cleanResults)

	jsonData, err := json.MarshalIndent(cleanResults, "", "  ")

	if err != nil {
		return
	}

	ResponseData = string(jsonData)
	http.HandleFunc("/plugins/results", GetPluginResultsHandler)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func GetPluginResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ResponseData))
}

func MapDereference(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			data[key] = MapDereference(v)
		case *interface{}:
			data[key] = *v
		}
	}
	return data
}
