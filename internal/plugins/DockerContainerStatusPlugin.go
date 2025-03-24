package plugins

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerContainerStatusPlugin struct{}

func (plugin DockerContainerStatusPlugin) Name() string {
	return "Docker Container Status Plugin"
}

func (plugin DockerContainerStatusPlugin) Collect() (map[string]interface{}, error) {

	results := make(map[string]interface{})
	results["Type"] = reflect.TypeOf(plugin).Name()

	now := time.Now()

	containers, err := GetContainers()
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		name := ""
		if len(container.Names) > 0 {
			name = strings.Replace(container.Names[0], "/", "", 1)
		}

		results[name] = map[string]interface{}{
			"State":        container.State,
			"Image":        container.Image,
			"Uptime":       now.Sub(time.Unix(container.Created, 0)).Round(time.Second).String(),
			"Container Id": container.ID[:10],
			"Ports":        GetPortsForContainer(container),
		}
	}

	return map[string]interface{}{
		"containers": results,
	}, nil
}

func GetContainers() ([]container.Summary, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Формат вывода можно настроить под себя. Например, "IP:Public->Private/Type"
func GetPortsForContainer(containerInfo container.Summary) string {
	seen := make(map[string]struct{})
	var ports []string
	for _, p := range containerInfo.Ports {
		portStr := fmt.Sprintf("%d->%d", p.PublicPort, p.PrivatePort)
		if _, exists := seen[portStr]; !exists {
			seen[portStr] = struct{}{}
			ports = append(ports, portStr)
		}
	}
	return strings.Join(ports, ", ")
}
