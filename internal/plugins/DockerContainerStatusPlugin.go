package plugins

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerContainerStatusPlugin struct{}

func (d DockerContainerStatusPlugin) Name() string {
	return "Docker Container Status Plugin"
}

func (d DockerContainerStatusPlugin) Collect() (map[string]interface{}, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	now := time.Now()

	for _, container := range containers {
		name := ""
		if len(container.Names) > 0 {
			name = container.Names[0]
		}

		uptime := now.Sub(time.Unix(container.Created, 0)).String()

		results = append(results, map[string]interface{}{
			"name":   name,
			"uptime": uptime,
		})
	}

	return map[string]interface{}{
		"containers": results,
	}, nil
}
