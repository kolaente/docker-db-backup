package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"sync"
)

var (
	store map[string]*types.ContainerJSON
	lock  sync.Mutex
)

func init() {
	store = make(map[string]*types.ContainerJSON)
}

func storeContainers(c *client.Client, containers []types.Container) {
	lock.Lock()
	defer lock.Unlock()

	for _, container := range containers {
		if container.State != "running" {
			continue
		}

		info, err := c.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Fatalf("Could not get Container info: %s", err)
		}

		store[container.ID] = &info

		fmt.Printf("Container: %s, %s image: %s, labels: %v, ip: %v, env: %v\n", info.Name, info.State.Status, info.Image, info.Config.Labels, info.NetworkSettings.Networks, info.Config.Env)
	}
}
