package plugins

// Common interface for plugin
// Defines the name of the plugin and
// and its function for retrieving data
type Plugin interface {
	Name() string
	Collect() (map[string]interface{}, error)
}
