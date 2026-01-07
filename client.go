package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func getClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func getContainers(c *client.Client) ([]container.Summary, error) {
	return c.ContainerList(context.Background(), container.ListOptions{All: true})
}
