package webServer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"project/internal/DTO"
	"project/internal/agentDispatcher"
	"project/internal/applicationConfigurationDispatcher"
	"project/internal/pluginDispatcher"
	"project/internal/utils"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var ServerInfo = applicationConfigurationDispatcher.ServerInfo{}

// Initializes endpoints and starts the webserver
func RunServer(webServerConfig applicationConfigurationDispatcher.WebServerConfig, serverInfo applicationConfigurationDispatcher.ServerInfo) {
	ServerInfo = serverInfo

	mux := InitEndpoints()

	err := StartServer(mux, webServerConfig)
	if err != nil {
		log.Fatalf("Error while webserver starting %v", err)
	}
}

// Defines enpoint handlers
func InitEndpoints() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/plugins/results", basicAuthMiddleware(GetPluginResultsHandler))

	if ServerInfo.AgentMode {
		return mux
	}

	mux.HandleFunc("/", IndexPageHandler)
	mux.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	mux.HandleFunc("/events", sseHandler)
	mux.HandleFunc("/styles.css", serveMergedCSS)

	return mux
}

// Fetches data from agents
// and responds data to client in a JSON form
func GetPluginResultsHandler(responseWriter http.ResponseWriter, r *http.Request) {
	pluginResultCollection := pluginDispatcher.GetPluginsJsonData()

	responseWriter.Header().Set("Content-Type", "application/json")
	_, err := responseWriter.Write([]byte(string(pluginResultCollection)))
	if err != nil {
		log.Printf("Error while getting plugins data: %v", err)
	}
}

// Starts webserver
func StartServer(mux *http.ServeMux, cfg applicationConfigurationDispatcher.WebServerConfig) error {
	if mux == nil {
		return errors.New("mux must not be nil")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}

	addr := net.JoinHostPort("", strconv.Itoa(cfg.Port))

	srv := &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: log.New(os.Stdout, "http: ", log.LstdFlags),
	}

	OnServerStopping(srv)

	log.Printf("Server starting on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not start server: %w", err)
	}

	log.Printf("Server stopped")
	return nil
}

// Shuts down the server when termination signal recieved
func OnServerStopping(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		log.Printf("Shutdown signal received, stopping server…")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()
}

// Serves the main page or a 404 page
func IndexPageHandler(responseWriter http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		HandleNotFoundPage(responseWriter)
		return
	}
	totalResults := make(map[string]interface{})

	pageTemplate := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"serverInfo": getServerName,
	}).ParseFiles("templates/index.html"))

	responseWriter.Header().Set("Content-Type", "text/html")
	if err := pageTemplate.Execute(responseWriter, totalResults); err != nil {
		http.Error(responseWriter, "Error rendering template", http.StatusInternalServerError)
	}
}

// Hanlde NotFound response
func HandleNotFoundPage(responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(http.StatusNotFound)
	tmpl404 := template.Must(template.
		New("NotFound.html").
		Funcs(template.FuncMap{"serverInfo": getServerName}).
		ParseFiles("templates/NotFound.html"))
	err := tmpl404.Execute(responseWriter, ServerInfo)
	if err != nil {
		log.Printf("Error while handling 404 error : %v", err)
	}
}

// Establishes a Server-Sent Events (SSE) connection for sending real-time data
func sseHandler(responseWriter http.ResponseWriter, r *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")

	HandleAgents(responseWriter)
}

// Fetches metrics from active agents and sends them as real-time data
func HandleAgents(responseWriter http.ResponseWriter) {
	agents := agentDispatcher.GetAgents()

	resultsChannel := make(chan DTO.AgentDataChunk, len(agents))

	agentDispatcher.GetMetricsFromAgentsAsync(agents, resultsChannel)

	for _, agent := range agents {
		if agent.Active {
			res := <-resultsChannel

			jsonData, _ := json.Marshal(res)

			_, err := fmt.Fprintf(responseWriter, "data: %s\n\n", jsonData)
			if err != nil {
				log.Printf("Error while responsing to a client : %v", err)
			}

			responseWriter.(http.Flusher).Flush()
		}
	}
}

// Returns the server's name formatted for HTML templates
func getServerName() template.HTML {
	return template.HTML(ServerInfo.Name)
}

// Serves combined CSS files from the specified directory
func serveMergedCSS(w http.ResponseWriter, r *http.Request) {
	css, err := utils.MergeFilesWithExtension(".", ".css")
	if err != nil {
		http.Error(w, "CSS file can't loaded correctly ", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	_, err = w.Write(css)
	if err != nil {
		log.Printf("Error while merging css files : %v", err)
	}
}

// Provides HTTP Basic Authentication for endpoints
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
