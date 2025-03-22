package plugins

import (
	"context"
	"reflect"
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
			name = container.Names[0]
		}

		results[name] = map[string]interface{}{
			"state":  container.State,
			"uptime": now.Sub(time.Unix(container.Created, 0)).Round(time.Second).String(),
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
