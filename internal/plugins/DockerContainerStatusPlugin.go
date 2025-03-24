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

// Plugin that collects information of Docker containers
type DockerContainerStatusPlugin struct{}

// Returns the name of plugin instance
func (plugin DockerContainerStatusPlugin) Name() string {
	return "Docker Container Status Plugin"
}

// Collect information for each Docker container
// Returns a map that includes container's
// state, image, uptime, container ID and port mappings
func (plugin DockerContainerStatusPlugin) Collect() (map[string]interface{}, error) {

	results := make(map[string]interface{})
	results["Type"] = reflect.TypeOf(plugin).Name()

	containers, err := GetContainers()
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		name := ""
		if len(container.Names) > 0 {
			// Remove the leading slash from the container name
			name = strings.Replace(container.Names[0], "/", "", 1)
		}

		results[name] = map[string]interface{}{
			"State":        container.State,
			"Image":        container.Image,
			"Uptime":       time.Since(time.Unix(container.Created, 0)).Round(time.Second).String(),
			"Container Id": container.ID[:10],
			"Ports":        GetPortsForContainer(container),
		}
	}

	return map[string]interface{}{
		"containers": results,
	}, nil
}

// Connects to Docker and retrieves a list of containers
// Returns a list of containers info
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

// Returns the port mappings for a container
// in "PublicPort->PrivatePort" format
func GetPortsForContainer(containerInfo container.Summary) string {
	duplicatesPorts := make(map[string]struct{})
	var ports []string
	for _, p := range containerInfo.Ports {
		formattedPort := fmt.Sprintf("%d->%d", p.PublicPort, p.PrivatePort)

		if _, exists := duplicatesPorts[formattedPort]; !exists {
			duplicatesPorts[formattedPort] = struct{}{}
			ports = append(ports, formattedPort)
		}
	}
	return strings.Join(ports, ", ")
}
