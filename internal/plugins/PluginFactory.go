package plugins

import (
	"fmt"
	"reflect"
)

func init() {
	RegisterPluginType("VersionPlugin", func() Plugin { return &VersionPlugin{} })
	RegisterPluginType("HardDriveFreeSpacePlugin", func() Plugin { return &HardDriveFreeSpacePlugin{} })
	RegisterPluginType("PostgreSqlQueryPlugin", func() Plugin { return &PostgreSqlQueryPlugin{} })
}

var registeredPlugins = make(map[string]func() Plugin)

func RegisterPluginType(name string, constructor func() Plugin) {
	registeredPlugins[name] = constructor
}

func CreatePlugin(name string, params map[string]interface{}) (Plugin, error) {
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
