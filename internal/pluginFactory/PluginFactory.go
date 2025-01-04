package pluginFactory

import (
	"fmt"
	"reflect"

	"project/internal/plugins"
)

func init() {
	RegisterPluginType("VersionPlugin", func() plugins.Plugin { return &plugins.VersionPlugin{} })
	RegisterPluginType("HardDriveFreeSpacePlugin", func() plugins.Plugin { return &plugins.HardDriveFreeSpacePlugin{} })
	RegisterPluginType("PostgreSqlQueryPlugin", func() plugins.Plugin { return &plugins.PostgreSqlQueryPlugin{} })
}

var registeredPlugins = make(map[string]func() plugins.Plugin)

func RegisterPluginType(name string, constructor func() plugins.Plugin) {
	registeredPlugins[name] = constructor
}

func CreatePlugin(name string, params map[string]interface{}) (plugins.Plugin, error) {
	constructor, exists := registeredPlugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin type '%s' not registered", name)
	}

	// Создаём новый экземпляр плагина через конструктор
	plugin := constructor()

	// Заполняем поля структуры из params
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
