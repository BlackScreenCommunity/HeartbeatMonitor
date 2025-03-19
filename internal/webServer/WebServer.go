package webServer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"strconv"
	"strings"
	"time"
)

var ServerInfo = applicationConfigurationDispatcher.ServerInfo{}

func RunServer(webServerConfig applicationConfigurationDispatcher.WebServerConfig, serverInfo applicationConfigurationDispatcher.ServerInfo) {
	ServerInfo = serverInfo

	InitEndpoints()
	StartServer(webServerConfig)
}

// Defines enpoint handlers
func InitEndpoints() {
	http.HandleFunc("/plugins/results", basicAuthMiddleware(GetPluginResultsHandler))
	http.HandleFunc("/", IndexPageHandler)
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/styles.css", serveMergedCSS)
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
		if agent.Active {
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
					Duration: math.Floor((time.Duration(time.Since(start)).Seconds())*100) / 100,
				}
			}(i, agent)
		}
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

func serveMergedCSS(w http.ResponseWriter, r *http.Request) {
	css, err := mergeCSSFiles(".")
	if err != nil {
		http.Error(w, "Ошибка загрузки CSS", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(css)
}

func mergeCSSFiles(dir string) ([]byte, error) {
	var buffer bytes.Buffer

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".css" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(&buffer, file)
			if err != nil {
				return err
			}
			buffer.WriteString("\n")
		}
		return nil
	})

	return buffer.Bytes(), err
}

// Basic auth implementation for endpoints
func basicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check header format
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Invalid authorization format", http.StatusBadRequest)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 || credentials[0] != ServerInfo.Login || credentials[1] != ServerInfo.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
