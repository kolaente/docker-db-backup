package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func getClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func getContainers(c *client.Client) ([]types.Container, error) {
	return c.ContainerList(context.Background(), types.ContainerListOptions{All: true})
}
