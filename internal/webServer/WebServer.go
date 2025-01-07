package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/pluginDispatcher"
	"project/internal/utils"
)

func RunServer() {

	InitEndpoints()

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func InitEndpoints() {

	http.HandleFunc("/plugins/results", GetPluginResultsHandler)

}

func GetPluginResultsHandler(responseWriter http.ResponseWriter, r *http.Request) {
	pluginResultCollection := pluginDispatcher.CollectAll()

	cleanResults := utils.MapDereference(pluginResultCollection)
	pluginDispatcher.PrintPluginResult(cleanResults)

	jsonData, err := json.MarshalIndent(cleanResults, "", "  ")

	if err != nil {
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write([]byte(string(jsonData)))
}
