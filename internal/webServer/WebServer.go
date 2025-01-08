package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"project/internal/pluginDispatcher"
	"project/internal/utils"
	"text/template"
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
	http.HandleFunc("/", IndexPageHandler)

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

func getTemplatePath() string {
	absPath, err := filepath.Abs("./internal/webserver/index.html")
	if err != nil {
		panic(err)
	}
	return absPath
}

func IndexPageHandler(responseWriter http.ResponseWriter, r *http.Request) {
	pluginResultCollection := pluginDispatcher.CollectAll()

	cleanResults := utils.MapDereference(pluginResultCollection)
	pageTemplate := template.Must(template.ParseFiles(getTemplatePath()))

	pageData := struct {
		Title   string
		Plugins map[string]interface{}
	}{
		Title:   "Plugin Data",
		Plugins: cleanResults,
	}

	responseWriter.Header().Set("Content-Type", "text/html")
	if err := pageTemplate.Execute(responseWriter, pageData); err != nil {
		http.Error(responseWriter, "Error rendering template", http.StatusInternalServerError)
	}

}
