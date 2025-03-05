package plugins

import (
	"fmt"
	"project/internal/utils"

	"github.com/godbus/dbus/v5"
)

type ServiceStatusPlugin struct {
	Services []interface{}
}

func (v ServiceStatusPlugin) Name() string {
	return "ServiceStatusPlugin"
}

func (v ServiceStatusPlugin) Collect() (map[string]interface{}, error) {
	var servicesToCheck = utils.ConvertInterfaceArrayToStringArray(v.Services)

	results := make(map[string]interface{})

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к D-Bus: %v", err)
	}
	defer conn.Close()

	for _, service := range servicesToCheck {
		status, err := getServiceStatusDBus(conn, service)
		if err != nil {
			results[service] = fmt.Sprintf("Ошибка: %v", err)
		} else {
			results[service] = status
		}
	}

	return results, nil
}

func getServiceStatusDBus(conn *dbus.Conn, service string) (string, error) {
	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var unitPath dbus.ObjectPath
	err := obj.Call("org.freedesktop.systemd1.Manager.GetUnit", 0, service+".service").Store(&unitPath)
	if err != nil {
		return "", fmt.Errorf("сервис %s не найден", service)
	}

	unit := conn.Object("org.freedesktop.systemd1", unitPath)
	variant, err := unit.GetProperty("org.freedesktop.systemd1.Unit.ActiveState")
	if err != nil {
		return "", fmt.Errorf("не удалось получить состояние сервиса %s", service)
	}

	status, ok := variant.Value().(string)
	if !ok {
		return "", fmt.Errorf("неизвестный формат ответа для сервиса %s", service)
	}

	return status, nil
}
