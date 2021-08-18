package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"sync"
)

var (
	store map[string]Dumper
	lock  sync.Mutex
)

func init() {
	store = make(map[string]Dumper)
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

		dumper := NewDumperFromContainer(&info)
		if dumper == nil {
			continue
		}

		log.Printf("Found container %s\n", container.Names)

		store[container.ID] = dumper
	}
}
