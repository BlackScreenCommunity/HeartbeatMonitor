package webServer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"reflect"
	"strconv"
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
	pluginResultCollection := pluginDispatcher.GetPluginsJsonData()

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(pluginResultCollection), &data); err != nil {
		panic(err)
	}

	agentResultCollection := agentDispatcher.GetMetricsFromAgents()

	totalResults := make(map[string]interface{})
	totalResults[ServerInfo.Name] = data

	for k, v := range agentResultCollection {
		totalResults[k] = v
	}

	pageTemplate := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"renderList": renderList,
		"serverInfo": getServerName,
	}).ParseFiles("templates/index.html"))

	responseWriter.Header().Set("Content-Type", "text/html")
	if err := pageTemplate.Execute(responseWriter, totalResults); err != nil {
		http.Error(responseWriter, "Error rendering template", http.StatusInternalServerError)
	}
}

func getServerName() template.HTML {
	return template.HTML(ServerInfo.Name)
}

func renderList(data interface{}) template.HTML {
	switch reflect.TypeOf(data).Kind() {

	case reflect.Map:
		html := ""

		dataMap := data.(map[string]interface{})

		isWarning := false
		value, exists := dataMap["isWarning"]
		if exists {
			isWarning = value.(bool)
			delete(dataMap, "isWarning")
		}

		for key, value := range dataMap {
			if !isWarning {
				html += "<div class='widget'>"
			} else {
				html += "<div class='widget warning'>"
			}
			isWarning = false
			html += "<div class='widget-title'>" + key + ":</div> "
			html += string(renderList(value)) + ""
			html += "</div>"
		}
		return template.HTML(html)

	case reflect.Slice:
		html := "<div class='data_array'>"
		for _, value := range data.([]interface{}) {
			html += "<div class='data_array_element'>" + string(renderList(value)) + "</div>"
		}
		html += "</div>"
		return template.HTML(html)
	default:
		return template.HTML("<div class='widget-data'>" + template.HTMLEscapeString(fmt.Sprintf("%v", data)) + "</div>")
	}
}
