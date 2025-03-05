package pluginFactory

import (
	"fmt"
	"reflect"

	"project/internal/plugins"
	postgreSqlQueryPlugin "project/internal/plugins/PostgreSqlQueryPluginDirectory"
)

// Registers all available plugins
func init() {
	registeredPlugins["VersionPlugin"] = func() plugins.Plugin { return &plugins.VersionPlugin{} }
	registeredPlugins["HardDriveFreeSpacePlugin"] = func() plugins.Plugin { return &plugins.HardDriveFreeSpacePlugin{} }
	registeredPlugins["PostgreSqlQueryPlugin"] = func() plugins.Plugin { return &postgreSqlQueryPlugin.PostgreSqlQueryPlugin{} }
	registeredPlugins["FolderSizePlugin"] = func() plugins.Plugin { return &plugins.FolderSizePlugin{} }
	registeredPlugins["ServerTimePlugin"] = func() plugins.Plugin { return &plugins.ServerTimePlugin{} }
	registeredPlugins["ServiceStatusPlugin"] = func() plugins.Plugin { return &plugins.ServiceStatusPlugin{} }
}

// Collection of registered plugins
// Key: plugin name, Value: function that creates a new plugin instance.
var registeredPlugins = make(map[string]func() plugins.Plugin)

// Creates a new plugin instance by its name and fills its fields with values from recieved parameters.
//
// name - the registered type name of the plugin.
// params - a map of parameters, where keys match the field names of the plugin's struct.
//
// Returns:
// - A plugin instance (`plugins.Plugin`) if created successfully.
// - An error (`error`) if:
//   - The plugin name is not registered.
//   - A field specified in `params` does not exist in the plugin struct.
//   - The type of a value in `params` does not match the type of the corresponding struct field.
func CreatePlugin(name string, params map[string]interface{}) (plugins.Plugin, error) {
	constructor, exists := registeredPlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin type '%s' not registered", name)
	}

	plugin := constructor()

	pluginValue := reflect.ValueOf(plugin).Elem()
	for key, value := range params {
		field := pluginValue.FieldByName(key)
		if !field.IsValid() {
			return nil, fmt.Errorf("field '%s' not found in plugin '%s'", key, name)
		}

		if !field.CanSet() {
			return nil, fmt.Errorf("field '%s' in plugin '%s' is not settable", key, name)
		}

		fieldValue := reflect.ValueOf(value)
		if field.Type() != fieldValue.Type() {
			return nil, fmt.Errorf("type mismatch for field '%s' in plugin '%s': expected '%s', got '%s'",
				key, name, field.Type(), fieldValue.Type())
		}

		field.Set(fieldValue)
	}

	return plugin, nil
}
