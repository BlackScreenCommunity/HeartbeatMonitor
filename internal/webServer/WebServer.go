package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"project/internal/pluginDispatcher"
	"project/internal/utils"
	"reflect"
)

func renderList(data interface{}) template.HTML {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		html := "<ul>"
		for key, value := range data.(map[string]interface{}) {
			html += "<li><strong>" + key + ":</strong> " + string(renderList(value)) + "</li>"
		}
		html += "</ul>"
		return template.HTML(html)
	case reflect.Slice:
		html := "<ul>"
		for _, value := range data.([]interface{}) {
			html += "<li>" + string(renderList(value)) + "</li>"
		}
		html += "</ul>"
		return template.HTML(html)
	default:
		// return template.HTML(template.HTMLEscapeString(reflect.ValueOf(data).String()))
		return template.HTML(template.HTMLEscapeString(fmt.Sprintf("%v", data)))
	}
}

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
	absPath, err := filepath.Abs("./templates/index.html")
	if err != nil {
		panic(err)
	}
	return absPath
}

func IndexPageHandler(responseWriter http.ResponseWriter, r *http.Request) {
	pluginResultCollection := pluginDispatcher.CollectAll()

	cleanResults := utils.MapDereference(pluginResultCollection)

	jsonData, err := json.MarshalIndent(cleanResults, "", "  ")

	if err != nil {
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		panic(err)
	}

	pageTemplate := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"renderList": renderList,
	}).ParseFiles("templates/index.html"))

	responseWriter.Header().Set("Content-Type", "text/html")
	if err := pageTemplate.Execute(responseWriter, data); err != nil {
		http.Error(responseWriter, "Error rendering template", http.StatusInternalServerError)
	}

}
