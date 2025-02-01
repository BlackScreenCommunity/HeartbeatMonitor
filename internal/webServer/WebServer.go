package webServer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"strconv"
	"time"
)

var ServerInfo = applicationConfigurationDispatcher.ServerInfo{}

func RunServer(webServerConfig applicationConfigurationDispatcher.WebServerConfig, serverInfo applicationConfigurationDispatcher.ServerInfo) {
	ServerInfo = serverInfo

	InitEndpoints()
	StartServer(webServerConfig)
}

func InitEndpoints() {
	http.HandleFunc("/plugins/results", GetPluginResultsHandler)
	http.HandleFunc("/", IndexPageHandler)
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	http.HandleFunc("/events", sseHandler)
}

func GetPluginResultsHandler(responseWriter http.ResponseWriter, r *http.Request) {
	pluginResultCollection := pluginDispatcher.GetPluginsJsonData()

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write([]byte(string(pluginResultCollection)))
}

func StartServer(config applicationConfigurationDispatcher.WebServerConfig) {
	fmt.Println("Server is running on port " + strconv.Itoa(config.Port))
	if err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func IndexPageHandler(responseWriter http.ResponseWriter, r *http.Request) {
	totalResults := make(map[string]interface{})

	pageTemplate := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"serverInfo": getServerName,
	}).ParseFiles("templates/index.html"))

	responseWriter.Header().Set("Content-Type", "text/html")
	if err := pageTemplate.Execute(responseWriter, totalResults); err != nil {
		http.Error(responseWriter, "Error rendering template", http.StatusInternalServerError)
	}
}

func sseHandler(responseWriter http.ResponseWriter, r *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")

	HandleAgents(responseWriter)
}

func HandlePlugins(responseWriter http.ResponseWriter) {
	for name, plugin := range pluginDispatcher.GetPlugins() {
		data, err := plugin.Collect()
		if err != nil {
			fmt.Printf("Error collecting data from plugin %s: %v\n", name, err)
			continue
		}

		type DataChunk struct {
			PluginName string                 `json:"plugin_name"`
			Data       map[string]interface{} `json:"data"`
		}

		pluginData := DataChunk{
			PluginName: plugin.Name(),
			Data:       data,
		}

		jsonData, _ := json.Marshal(pluginData)

		fmt.Fprintf(responseWriter, "data: %s\n\n", jsonData)
		responseWriter.(http.Flusher).Flush()

	}
}

func HandleAgents(responseWriter http.ResponseWriter) {
	agents := agentDispatcher.GetAgents()
	resultsChannel := make(chan struct {
		Key      string
		Result   map[string]interface{}
		Duration float64
	}, len(agents))

	for i, agent := range agents {
		go func(i int, agent applicationConfigurationDispatcher.AgentConfig) {
			start := time.Now()
			result := agentDispatcher.GetMetricsFromSingleAgent(agent)
			resultsChannel <- struct {
				Key      string
				Result   map[string]interface{}
				Duration float64
			}{
				Key:      strconv.Itoa(i+1) + ". " + agent.Name,
				Result:   result,
				Duration: time.Duration(time.Since(start)).Seconds(),
			}
		}(i, agent)

	}

	type AgentDataChunk struct {
		AgentName string                 `json:"agent_name"`
		Data      map[string]interface{} `json:"data"`
		Duration  float64                `json:"duration"`
	}

	for range agents {
		res := <-resultsChannel

		agentDataChunk := AgentDataChunk{
			AgentName: res.Key,
			Data:      res.Result,
			Duration:  res.Duration,
		}

		jsonData, _ := json.Marshal(agentDataChunk)

		fmt.Fprintf(responseWriter, "data: %s\n\n", jsonData)
		responseWriter.(http.Flusher).Flush()

	}
}

func getServerName() template.HTML {
	return template.HTML(ServerInfo.Name)
}
